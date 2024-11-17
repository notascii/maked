package server

import (
	"fmt"
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

type MakeInstruction struct {
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
