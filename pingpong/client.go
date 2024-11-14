package main

import (
	"fmt"
	"net/rpc"
	"time"
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


func measure_perf(N int)  (int,float64){
	// Connect to the server
	client, err := rpc.Dial("tcp", "172.18.20.22:1234")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Warm-up ping to establish the connection
	warmupArgs := &PingArgs{Data: []byte{1}} // Small 1-byte message
	var warmupReply PingResponse
	_ = client.Call("PingPongService.Ping", warmupArgs, &warmupReply) // Ignoring error for warm-up

	// Build the message of size N
	message := make([]byte, N)

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
	args = &PingArgs{Data: message} // N-byte message
	err = client.Call("PingPongService.Ping", args, &reply)
	if err != nil {
		panic(err)
	}
	end := time.Since(start).Seconds()

	// Calculate throughput in bytes per second
	throughput := float64(N) / end - rtt
	throughputMB := throughput / (1024 * 1024)
	
	fmt.Printf("Message size: %d bytes\n", N)
	fmt.Printf("Throughput: %f MB/second\n", throughputMB)

	return N,throughputMB

}


func main() {
	file, err := os.OpenFile("./perf/logs/throughput.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	val:=1024
	for i:=0;i<=10;i++{
		messageSize,throughput:=measure_perf(val)
		// Write message size and throughput to the log file
		_, err = fmt.Fprintf(file, "%d: %.2f\n", messageSize, throughput)
		if err != nil {
			panic(err)
		}
		val=val*2
	}
}
