package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var storageAbs string = "/maked/server/server_storage/"

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	storage := homeDir + storageAbs
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Instructions) > 0 {
		for i, ins := range p.Instructions {
			test := true
			var list []FileStruct
			for _, dep := range ins.Dependencies {
				file, err := os.ReadFile(storage + dep)

				if err != nil {
					// We check the repo
					file, err = os.ReadFile(p.Directory + dep)
					if err != nil {
						test = false
						fmt.Println("dependency missing : ", dep)
					}

				}
				tmp := FileStruct{Data: file, FileName: dep}
				list = append(list, tmp)
			}
			if test {
				reply.Value = 2
				reply.Command = ins.Command
				reply.Dependencies = list
				p.Instructions = append(p.Instructions[:i], p.Instructions[i+1:]...)
				return nil
			}
		}
		// No dependencies available found...
		reply.Value = 1
		return nil
	} else {
		reply.Value = 0
	}
	return nil
}

func (p *MakeService) SendFile(args *FileStruct, reply *FileStruct) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	storage := homeDir + storageAbs
	// Reply with an acknowledgment byte
	fmt.Println("Name of the file received : ", args.FileName)

	err = os.WriteFile(storage+args.FileName, args.Data, 0644)
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
