package main

import (
	"fmt"
	"net/rpc"
	"os"
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
func establishConn() (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Warm-up ping to establish the connection
	warmupArgs := &PingArgs{Data: []byte{1}} // Small 1-byte message
	var warmupReply PingResponse
	_ = client.Call("PingPongService.Ping", warmupArgs, &warmupReply) // Ignoring error for warm-up

	// Read the file content
	fileData, err := os.ReadFile(filename)

	// Send file content to server to calculate throughput
	start := time.Now()
	err = client.Call("PingPongService.Ping", &PingArgs{Data: fileData}, &PingResponse{})
	if err != nil {
		return 0, 0,0, err
	}
	elapsed := time.Since(start).Seconds()

	// Calculate throughput in MB/s
	throughput := float64(len(fileData)) / elapsed
	throughputMB := throughput / (1024 * 1024) // Convert to MB/s

	return len(fileData), throughputMB,elapsed/2, nil
}


// Create a file with the specified size in the client_storage directory.
func createFileWithSize(filename string, size int) error {
	// Create a file with the specified size
	data := make([]byte, size)
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func measure_perf() {
	// Establish the connection to the server
	client, err := establishConn()
	if err != nil {
		panic(err)
	}
	defer client.Close()

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

	// Create files first
	messageSize := 1024
	for messageSize <= 1048576 {
		// Create a file with the current message size
		filename := fmt.Sprintf("./client_storage/file_%d.bytes", messageSize)
		err := createFileWithSize(filename, messageSize)
		if err != nil {
			fmt.Printf("Error creating file with size %d bytes: %v\n", messageSize, err)
			continue
		}
		// Double the message size for the next iteration
		messageSize *= 2
	}

	// Now, measure throughput and latency for each file
	messageSize = 1024
	for messageSize <= 1048576 {
		filename := fmt.Sprintf("./client_storage/file_%d.bytes", messageSize)

		// Measure throughput
		_, throughput,latency,err := measurePerf(client, filename)
		if err != nil {
			fmt.Printf("Error measuring throughput for file %s: %v\n", filename, err)
			continue
		}

		// Log results to files
		_, err = fmt.Fprintf(fileThroughput, "%d: %.2f MB/s\n", messageSize, throughput)
		if err != nil {
			fmt.Printf("Error writing to throughput log: %v\n", err)
		}

		_, err = fmt.Fprintf(fileLatency, "%d: %.10f ms\n", messageSize, latency*1000) // Convert latency to milliseconds
		if err != nil {
			fmt.Printf("Error writing to latency log: %v\n", err)
		}

		// Print results to console
		fmt.Printf("File: %s, Throughput: %.2f MB/s, Latency: %.10f ms\n", messageSize, throughput, latency*1000)

		// Double the message size for the next iteration
		messageSize *= 2
	}
}

func main() {
	// Measure performance
	measure_perf()
}
