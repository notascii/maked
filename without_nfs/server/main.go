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
	if len(args) != 1 {
		log.Fatalf("Excepted 1 argument (name of repo containing the makefile)")
	}

	makefilePath := "../../makefiles/" + args[0] + "/Makefile"
	makefileDir := "../../makefiles/" + args[0] + "/"
	// First we parse the makefile
	var g *Graph = GraphParser(makefilePath)

	// Now we execute all commands in the directory asked
	commandListe := []MakeElement{}
	launchMakefile(g, "", &commandListe)

	var dependenciesThere []string
	// Scan the makefile directory
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
		InstructionsToDo:  commandListe,
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

	// Start the scheduler goroutine
	go schedulerLoop(makeService)

	// Goroutine to monitor `commandListe` and shut down the server
	go func() {
		for {
			makeService.mu.Lock()
			if len(makeService.InstructionsToDo) == 0 && len(makeService.InstructionsInProgress) == 0 {
				fmt.Printf("No more instructions. Shutting down the server...\n")
				makeService.mu.Unlock()
				totalDuration := time.Since(timeStart)
				writeClientList(makeService.ClientList, totalDuration, nfsDirectory, args[0]+"_"+strconv.Itoa(currentClientId-1))
				listener.Close()
				os.Exit(0)
			}
			makeService.mu.Unlock()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}
}

func schedulerLoop(makeService *MakeService) {
	h := &RequestHeap{}
	heap.Init(h)

	for {
		select {
		case req := <-makeService.clientRequests:
			// Add the request to the heap
			heap.Push(h, req)
		default:
			// If heap is not empty, process the lowest ID request
			if h.Len() > 0 {
				req := heap.Pop(h).(ClientRequest)
				// Process the request
				res := makeService.processRequest()
				req.replyChan <- res
			} else {
				// No requests at the moment
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
