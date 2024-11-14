package main

import (
	"fmt"
	"net/rpc"
	"time"
	"io/ioutil"
	"os"
)

// PingArgs holds the data sent from client to server.
type PingArgs struct {
	Data []byte
}

// PingResponse holds the acknowledgment data sent back from server.
type PingResponse struct {
	Ack byte
}


func measure_perf(filename string) {
	// Connect to the server
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
	fileData, err := ioutil.ReadFile(filename)

	// First ping with size 1 to calculate RTT
	start := time.Now()
	args := &PingArgs{Data: []byte{1}} // 1-byte message for RTT
	var reply PingResponse
	err = client.Call("PingPongService.Ping", args, &reply)
	if err != nil {
		panic(err)
	}
	rtt := time.Since(start).Seconds()

	// Second ping with size N to calculate throughput
	start = time.Now()
	args =  &PingArgs{Data: fileData} // N-byte message
	err = client.Call("PingPongService.Ping", args, &reply)
	if err != nil {
		panic(err)
	}
	end := time.Since(start).Seconds()

	N,err := os.Stat(filename) // get the size of the file
	if err != nil {
		panic(err)
	}
	// Calculate throughput in bytes per second
	throughput := float64(N.Size()) / end - rtt
	// Convert throughput to megabytes per second
	throughputMB := throughput / (1024 * 1024)
	
	fmt.Printf("Throughput: %f MB/second\n", throughputMB)
}


func main() {
	measure_perf("./client_storage/disk.txt")
}
