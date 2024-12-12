package main

import (
	"log"
	"os"
	"sync"
)

var storageAbs string = "./server_storage/"
var currentClientId int = 1

type MakeService struct {
	mu                     sync.Mutex // handle concurrent access
	InstructionsToDo       []MakeElement
	InstructionsInProgress []MakeElement
	InstructionsDone       []string
	Directory              string
}

type FileStruct struct {
	Data      []byte
	FileName  string
	ReturnVal JobReturn
}

type JobReturn struct {
	CodeValue  int
	TargetName string
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
	Name         string
	ClientId 	 int
}

type MakeInstruction struct {
	Command      string
	Dependencies []FileStruct
}

type PingDef struct {
	ClientId int
}

func contains(list []string, word string) bool {
	for _, w := range list {
		if w == word {
			return true
		}
	}
	return false
}

func (p *MakeService) Ping(args *PingDef, reply *Order) error {
	storage := storageAbs
	p.mu.Lock()
	defer p.mu.Unlock()
	if (args.ClientId == -1)
	{
		reply.ClientId = currentClientId
		currentClientId ++
	}
	if len(p.InstructionsToDo) > 0 {
		for i, ins := range p.InstructionsToDo {
			test := true
			var list []FileStruct
			for _, dep := range ins.Dependencies {
				// We check that dependencies are in done list
				if !contains(p.InstructionsDone, dep) {
					log.Println("Error " + dep + " Missing for target " + ins.Name)
					test = false
					break
				}
				file, err := os.ReadFile(storage + dep)
				// We check is there is a file corresponding to this dependency to add it
				if err != nil {
					// We check the repo
					file, err = os.ReadFile(p.Directory + dep)
					if err != nil {
						test = false
						log.Println("dependency file missing : ", dep)
					}

				}
				tmp := FileStruct{Data: file, FileName: dep}
				list = append(list, tmp)
			}
			if test {
				reply.Value = 2
				reply.Command = ins.Command
				reply.Name = ins.Name
				reply.Dependencies = list
				p.InstructionsInProgress = append(p.InstructionsInProgress, p.InstructionsToDo[i])
				p.InstructionsToDo = append(p.InstructionsToDo[:i], p.InstructionsToDo[i+1:]...)
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
	// Check the code error
	if args.ReturnVal.CodeValue != 0 {
		// There is an error, so we put back the instruction in the beginning of todoList
		p.mu.Lock()
		defer p.mu.Unlock()

		// Find the instruction associated with the failed file and re-add to the todo list
		for i, ins := range p.InstructionsInProgress {
			if ins.Name == args.FileName {
				p.InstructionsInProgress = append(p.InstructionsInProgress[:i], p.InstructionsInProgress[i+1:]...)
				p.InstructionsToDo = append([]MakeElement{ins}, p.InstructionsToDo...)
				log.Printf("Instruction for %s moved back to the TODO list due to error\n", args.FileName)
				panic(1)
			}
		}
	} else { // No problem, so the task can be removed from progress list and put in Done
		p.mu.Lock()
		defer p.mu.Unlock()

		// Find and remove the instruction from InstructionsInProgress
		for i, ins := range p.InstructionsInProgress {
			if ins.Name == args.FileName {
				p.InstructionsInProgress = append(p.InstructionsInProgress[:i], p.InstructionsInProgress[i+1:]...)
				// Add the instruction name to InstructionsDone
				p.InstructionsDone = append(p.InstructionsDone, ins.Name)
				log.Printf("Instruction for target %s marked as done\n", ins.Name)
			}
		}
	}

	storage := storageAbs
	// Reply with an acknowledgment byte
	log.Println("Name of the file received : ", args.FileName)

	err := os.WriteFile(storage+args.FileName, args.Data, 0644)
	if err != nil {
		return err
	}
	log.Println("File received and stored as " + args.FileName)
	return nil
}

func (p *MakeService) Initialization(args *FileStruct, reply *FileList) error {
	// Reply with an acknowledgment byte

	files, err := os.ReadDir(p.Directory)
	if err != nil {
		log.Println("Impossible to read the directory : ", err)
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
