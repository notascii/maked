package main

import (
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
	// Makefile treatment

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalf("Excepted 1 argument (name of repo containing the makefile)")
	}

	////////////// TEST PREMIER
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

	// Add all file names in the directory to the dependenciesThere slice
	for _, file := range files {
		// Ensure it's not a directory and add to dependencies list
		if !file.IsDir() {
			dependenciesThere = append(dependenciesThere, file.Name())
		}
	}
	// Register the MakeService: This will allow the client to recognize it over the network
	makeService := &MakeService{
		InstructionsToDo:  commandListe,
		InstructionsDone:  dependenciesThere,
		Directory:         makefileDir,
		ClientList:        make(map[int][]Job),
		InstructionsStart: make(map[string]time.Time),
		InstructionsEnd:   make(map[string]time.Time),
	}
	rpc.Register(makeService)

	// Listen on TCP port 8090
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Server listening on port 8090...")

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
				os.Exit(0) // Exit the program
			}
			makeService.mu.Unlock()
		}
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}
}
