package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
)

type MakeService struct {
	mu    sync.Mutex // Pour protéger l'accès concurrent
	Items []string
}

type FileStruct struct {
	Data     []byte
	FileName string
}

type Message struct {
	Msg string
}

type Order struct {
	Value        byte
	Command      string
	Dependencies []FileStruct
}

type PingDef struct {
}

func (p *MakeService) Ping(args *PingDef, reply *Order) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Items) > 0 {
		reply.Value = 2
		reply.Command = p.Items[0]
		p.Items = p.Items[1:]

	} else {
		reply.Value = 0
	}
	return nil
}

func (p *MakeService) SendFile(args *FileStruct, reply *FileStruct) error {
	// Reply with an acknowledgment byte
	filePath := filepath.Join("server_storage", args.FileName)
	err := os.WriteFile(filePath, args.Data, 0644)
	if err != nil {
		return err
	}
	fmt.Println("File received and stored as " + args.FileName)
	return nil
}

func main() {
	// Register the MakeService: This will allow the client to recognize it over the network
	makeService := &MakeService{
		Items: []string{"ls -A", "ls -A", "ls -A"},
	}
	rpc.Register(makeService)

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
