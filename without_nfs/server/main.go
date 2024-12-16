package main

import (
	"container/heap"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

var nfsDirectory string = "~/maked/without_nfs"

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Printf("Error exepted two arguments")
		log.Fatalf("Usage: go run . <MakefileDIrectory> <O or 1>")
	}

	makefilePath := "../../makefiles/" + args[0] + "/Makefile"
	makefileDir := "../../makefiles/" + args[0] + "/"

	// First we execute the classic makefile
	var makeDuration time.Duration
	var makeLaunched bool
	if args[1] == "1" {
		makeDuration = launchClassicMake(makefileDir)
		makeLaunched = true
	} else {
		makeLaunched = false
	}

	// We clear ./server_storage
	clearDirectory(storageAbs)

	// Parse the makefile
	var g *Graph = GraphParser(makefilePath)

	// Generate the instructions list
	commandList := []MakeElement{}
	launchMakefile(g, "", &commandList)

	var dependenciesThere []string
	// Read the directory
	files, err := os.ReadDir(makefileDir)
	if err != nil {
		log.Fatalf("Failed to read makefile directory: %v", err)
	}
	for _, file := range files {
		if !file.IsDir() {
			dependenciesThere = append(dependenciesThere, file.Name())
		}
	}

	makeService := &MakeService{
		InstructionsToDo:  commandList,
		InstructionsDone:  dependenciesThere,
		Directory:         makefileDir,
		ClientList:        make(map[int][]Job),
		InstructionsStart: make(map[string]time.Time),
		InstructionsEnd:   make(map[string]time.Time),
		clientRequests:    make(chan ClientRequest, 100), // buffered channel
	}
	rpc.Register(makeService)

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Server listening on port 8090...\n")

	// Channel to signal a clean server shutdown
	done := make(chan struct{})

	// Start the scheduler goroutine
	go schedulerLoop(makeService, done)

	// Goroutine to monitor instruction completion
	go stopSignal(makeService, done, makeDuration, args[0], makeLaunched)

	// Main loop to accept connections
	for {
		select {
		case <-done:
			// Stop accepting new connections gracefully
			log.Println("Main loop received done signal. Stopping accept.")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				// If there's an error, check if we are stopping
				select {
				case <-done:
					log.Println("Stopping due to done signal after accept error.")
					return
				default:
					panic(err)
				}
			}
			go rpc.ServeConn(conn)
		}
	}
}

func schedulerLoop(makeService *MakeService, done chan struct{}) {
	h := &RequestHeap{}
	heap.Init(h)

	for {
		select {
		case req := <-makeService.clientRequests:
			heap.Push(h, req)
		case <-done:
			// The server is shutting down, exit the scheduler loop gracefully
			log.Println("Scheduler loop done signal received. Exiting schedulerLoop.")
			os.Exit(0)
		default:
			if h.Len() > 0 {
				req := heap.Pop(h).(ClientRequest)
				res := makeService.processRequest()
				req.replyChan <- res
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func stopSignal(makeService *MakeService, done chan struct{}, makeDuration time.Duration, makefileName string, makeLaunched bool) {
	for {
		makeService.mu.Lock()
		if len(makeService.InstructionsToDo) == 0 && len(makeService.InstructionsInProgress) == 0 {
			fmt.Printf("No more instructions. Shutting down the server...\n")
			makedDuration := time.Since(timeStart)
			writeResults(makeDuration, makeService.ClientList, makedDuration, nfsDirectory, makefileName, strconv.Itoa(currentClientId-1), makeLaunched)

			makeService.mu.Unlock()
			// Signal the main loop to stop
			close(done)

			// Let some time so the client stop
			time.Sleep(4 * time.Second)
			return
		}
		makeService.mu.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}
