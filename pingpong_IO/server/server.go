package main

import (
	"fmt"
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

// Ping is the method to handle the ping request and reply with a pong.
func (p *PingPongService) Ping(args *PingArgs, reply *PingResponse) error {
	// Reply with an acknowledgment byte
	err := os.WriteFile("./server_storage/backup_disk.txt", args.Data, 0644)
	if err != nil {
		return err
	}
	reply.Ack = 1
	fmt.Println("File received and stored as backup_log.txt")
	return nil
}

func main() {
	// Register the PingPongService: This will allow the client to recognize it over the network
	pingPongService := new(PingPongService)
	rpc.Register(pingPongService)

	// Listen on TCP port 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	println("Server listening on port 1234...")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)

	}
}
