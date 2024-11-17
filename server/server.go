package main

import (
	"fmt"
	"os"
	"sync"
)

type MakeService struct {
	mu           sync.Mutex // handle concurrent access
	Instructions []MakeElement
	Directory    string
}

type FileStruct struct {
	Data     []byte
	FileName string
}

type FileList struct {
	List []FileStruct
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

	if len(p.Instructions) > 0 {
		reply.Value = 2
		reply.Command = p.Instructions[0].Command
		p.Instructions = p.Instructions[1:]

	} else {
		reply.Value = 0
	}
	return nil
}

func (p *MakeService) SendFile(args *FileStruct, reply *FileStruct) error {
	// Reply with an acknowledgment byte
	fmt.Println("Name of the file received : ", args.FileName)

	err := os.WriteFile("./server_storage/"+args.FileName, args.Data, 0644)
	if err != nil {
		return err
	}
	fmt.Println("File received and stored as " + args.FileName)
	return nil
}

func (p *MakeService) Initialization(args *FileStruct, reply *FileList) error {
	// Reply with an acknowledgment byte

	files, err := os.ReadDir(p.Directory)
	if err != nil {
		fmt.Println("Impossible to read the directory : ", err)
	}
	var list []FileStruct
	var tmp FileStruct
	for _, file := range files {
		tmp.Data, err = os.ReadFile(p.Directory + file.Name())
		if err != nil {
			panic(err)
		}
		tmp.FileName = file.Name()
		list = append(list, tmp)
	}
	reply.List = list
	return nil
}
