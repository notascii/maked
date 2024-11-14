package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// PingPongService is a struct with no fields for our RPC service.
type PingPongService struct{}

// PingArgs holds the data sent from client to server.
type PingArgs struct {
	Data []byte
}

// PingResponse holds the acknowledgment data sent back to the client.
type PingResponse struct {
	Ack byte
}

// Ping handles the ping request and replies with an acknowledgment.
func (p *PingPongService) Ping(args *PingArgs, reply *PingResponse) error {
	// Log received message size
	log.Printf("Received a message of size %d\n", len(args.Data))
	reply.Ack = 1
	return nil
}

// HealthCheck handles a simple ping to verify server health.
func (p *PingPongService) HealthCheck(args *PingArgs, reply *PingResponse) error {
	reply.Ack = 1 // Just return an acknowledgment to indicate server is up
	return nil
}

func main() {
	// Create a log file
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Server starting...")

	// Register the PingPongService
	pingPongService := new(PingPongService)
	rpc.Register(pingPongService)

	// Listen on TCP port 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to listen on port 1234: %v", err)
		return
	}
	defer listener.Close()
	log.Println("Server listening on port 1234...")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go func() {
			log.Println("Accepted new connection")
			defer log.Println("Connection closed")
			rpc.ServeConn(conn)
		}()
	}
}
