package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
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
	client, err := rpc.Dial("tcp", hostname+":1234")
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Measures throughput by sending a message of size N and calculating bytes per second.
func measurePerfThroughput(client *rpc.Client, N int) (int, float64,float64,error) {

	// Send message of size N to calculate throughput
	message := make([]byte, N)
	start := time.Now()
	err := client.Call("PingPongService.Ping", &PingArgs{Data: message}, &PingResponse{})
	if err != nil {
		return 0, 0,0, err
	}
	elapsed := time.Since(start).Seconds()

	// Calculate throughput in MB/s
	throughput := (float64(N) / elapsed)
	throughputMB := throughput / (1024 * 1024)

	return N, throughputMB,elapsed/2, nil
}

func main() {
	// Establish a connection to the server
	messageSizeMaxstr := os.Args[2]
	messageSizeMax,err := strconv.Atoi(messageSizeMaxstr)
	client, err := establishConn(os.Args[1])
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

	// Run measurements for increasing message sizes
	messageSize := 1
	for i := 0; i <= messageSizeMax; i++ {
		// Measure throughput and latency
		size, throughput,latency, err := measurePerfThroughput(client, messageSize)
		if err != nil {
			fmt.Printf("Error measuring throughput for size %d: %v\n", messageSize, err)
			continue
		}

		// Log results to files
		_, err = fmt.Fprintf(fileThroughput, "%d: %.2f\n", size, throughput)
		if err != nil {
			fmt.Printf("Error writing to throughput log: %v\n", err)
		}
			println(latency)
		_, err = fmt.Fprintf(fileLatency, "%d: %.10f\n", size, latency*1000)
		if err != nil {
			fmt.Printf("Error writing to latency log: %v\n", err)
		}

		// Double the message size for the next iteration
		if messageSize == 1{
			messageSize=512
		}
		messageSize *= 2
	}
}
