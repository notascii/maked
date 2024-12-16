package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var storageAbs string = "../commun_storage"
var currentClientId int = 1
var firstPing bool = true
var timeStart time.Time

type MakeService struct {
	mu sync.Mutex
	// Work repartition
	InstructionsToDo       []MakeElement
	InstructionsInProgress []MakeElement
	InstructionsDone       []string
	// Work measure
	InstructionsStart map[string]time.Time
	InstructionsEnd   map[string]time.Time
	ClientList        map[int][]Job
	// Work space
	Directory string

	// Channel to receive client requests
	clientRequests chan ClientRequest
}

// A ClientRequest holds the client ID and a channel to send them back a response.
type ClientRequest struct {
	ClientId  int
	replyChan chan Order
}

type Job struct {
	Name     string
	Duration time.Duration
}

type FileStruct struct {
	FileName  string
	ReturnVal JobReturn
	ClientId  int
}

type JobReturn struct {
	CodeValue  int
	TargetName string
}

type FileList struct {
	List     []FileStruct
	ClientId int
}

type Message struct {
	Msg string
}

type Order struct {
	Value        byte
	Command      string
	Dependencies []FileStruct
	Name         string
	ClientId     int
}

type MakeInstruction struct {
	Command      string
	Dependencies []FileStruct
}

type PingDef struct {
	ClientId int
}

func importDependencies(sourceDirectory string, targetDirectory string) error {
	// Check if the source directory exists
	if _, err := os.Stat(sourceDirectory); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDirectory)
	}

	// Create the target directory if it doesn't exist
	if _, err := os.Stat(targetDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create target directory: %s, error: %v", targetDirectory, err)
		}
	}

	// Walk through the source directory
	err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate target path
		relPath, err := filepath.Rel(sourceDirectory, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDirectory, relPath)

		// If it's a directory, create it in the target
		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		// If it's a file, copy it to the target
		return copyFile(path, targetPath, info.Mode())
	})

	return err
}

// clearDirectory removes all files and subdirectories inside the given directory,
// leaving the directory itself intact.
func clearDirectory(directoryName string) {
	entries, err := os.ReadDir(directoryName)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", directoryName, err)
		return
	}

	for _, entry := range entries {
		path := filepath.Join(directoryName, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Failed to remove %s: %v\n", path, err)
		}
	}
}

// For priority queue (min-heap) based on client ID
type RequestHeap []ClientRequest

func (h RequestHeap) Len() int { return len(h) }
func (h RequestHeap) Less(i, j int) bool {
	return h[i].ClientId < h[j].ClientId
}
func (h RequestHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *RequestHeap) Push(x interface{}) {
	*h = append(*h, x.(ClientRequest))
}

func (h *RequestHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func contains(list []string, word string) bool {
	for _, w := range list {
		if w == word {
			return true
		}
	}
	return false
}

// The Ping method now just enqueues a request and waits for a response.
func (p *MakeService) Ping(args *PingDef, reply *Order) error {
	req := ClientRequest{
		ClientId:  args.ClientId,
		replyChan: make(chan Order),
	}
	// Send request to scheduler
	p.clientRequests <- req

	// Wait for the scheduler's response
	res := <-req.replyChan
	*reply = res
	return nil
}

// The scheduler will call this method to process the next request in priority order.
func (p *MakeService) processRequest() Order {
	p.mu.Lock()
	defer p.mu.Unlock()

	var reply Order
	if len(p.InstructionsToDo) > 0 {
		for i, ins := range p.InstructionsToDo {
			test := true
			var list []FileStruct
			for _, dep := range ins.Dependencies {
				// We check that dependencies are in done list
				if !contains(p.InstructionsDone, dep) {
					test = false
					break
				}
				tmp := FileStruct{FileName: dep}
				list = append(list, tmp)
			}
			if test {
				reply.Value = 2
				reply.Command = ins.Command
				reply.Name = ins.Name
				reply.Dependencies = list
				p.InstructionsStart[ins.Name] = time.Now()
				p.InstructionsInProgress = append(p.InstructionsInProgress, p.InstructionsToDo[i])
				p.InstructionsToDo = append(p.InstructionsToDo[:i], p.InstructionsToDo[i+1:]...)
				return reply
			}
		}
		// No dependencies available found...
		reply.Value = 1
		return reply
	} else {
		reply.Value = 0
	}
	return reply
}

func (p *MakeService) SendUpdate(args *FileStruct, reply *FileStruct) error {
	senderId := args.ClientId
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
				p.InstructionsEnd[ins.Name] = time.Now()
				p.ClientList[senderId] = append(p.ClientList[senderId], Job{
					Name:     ins.Name,
					Duration: time.Duration(p.InstructionsEnd[ins.Name].Sub(p.InstructionsStart[ins.Name]).Microseconds()),
				})
			}
		}
	}
	return nil
}

// Id distribution
func (p *MakeService) Initialization(args *PingDef, reply *FileList) error {
	if firstPing {
		firstPing = false
		timeStart = time.Now()
	}

	// Reply with an acknowledgment byte
	if args.ClientId == -1 {
		reply.ClientId = currentClientId
		currentClientId++
	} else {
		reply.ClientId = args.ClientId
	}
	return nil
}
