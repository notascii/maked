package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"time"
)

// PingArgs holds the data sent from client to server.
type PingArgs struct {
	Data []byte
}

// PingResponse holds the acknowledgment data sent back from server.
type PingResponse struct {
	Ack byte
}

// Establishes and returns an RPC connection to the server.
func establishConn(hostname string) (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", hostname+":1236")
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Reads a file and sends its content to the server, measuring throughput and latency.
func measurePerf(client *rpc.Client, filename string) (int, float64, float64, error) {
	// Warm-up ping to establish the connection
	warmupArgs := &PingArgs{Data: []byte{1}} // Small 1-byte message
	var warmupReply PingResponse
	_ = client.Call("PingPongService.Ping", warmupArgs, &warmupReply) // Ignoring error for warm-up

	// Read the file content
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, 0, 0, err
	}

	// Send file content to server to calculate throughput
	start := time.Now()
	err = client.Call("PingPongService.Ping", &PingArgs{Data: fileData}, &PingResponse{})
	if err != nil {
		return 0, 0, 0, err
	}
	elapsed := time.Since(start).Seconds()

	// Calculate throughput in MB/s
	throughput := float64(len(fileData)) / elapsed
	throughputMB := throughput / (1024 * 1024) // Convert to MB/s

	return len(fileData), throughputMB, elapsed / 2, nil
}

// Measures performance for files in a directory.
func measurePerfForFiles(client *rpc.Client, dirPath string) {
	// Open log files for throughput and latency
	fileThroughput, err := os.OpenFile("./perf/logs/throughput.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer fileThroughput.Close()

	fileLatency, err := os.OpenFile("./perf/logs/latency.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer fileLatency.Close()

	// Read all files from the directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(fmt.Errorf("error reading directory: %v", err))
	}

	for _, file := range files {
		filename := filepath.Join(dirPath, file.Name())

		// Measure throughput and latency
		fileSize, throughput, latency, err := measurePerf(client, filename)
		if err != nil {
			fmt.Printf("Error measuring throughput for file %s: %v\n", filename, err)
			continue
		}

		// Log results to files
		_, err = fmt.Fprintf(fileThroughput, "%d: %.2f\n", fileSize, throughput)
		if err != nil {
			fmt.Printf("Error writing to throughput log: %v\n", err)
		}

		_, err = fmt.Fprintf(fileLatency, "%d: %.10f\n", fileSize, latency*1000) // Convert latency to milliseconds
		if err != nil {
			fmt.Printf("Error writing to latency log: %v\n", err)
		}

		// Print results to console
		fmt.Printf("File: %s, Throughput: %.2f MB/s, Latency: %.10f ms\n", file.Name(), throughput, latency*1000)
	}
}

func main() {

	hostname := os.Args[1]
	dirPath := "./client/disk"

	client, err := establishConn(hostname)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	measurePerfForFiles(client, dirPath)
}
