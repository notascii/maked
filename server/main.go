package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

func main() {
	// Makefile treatment
	////////////// TEST PREMIER
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	makefilePath := homeDir + "/maked/makefiles/premier/Makefile"
	makefileDir := homeDir + "/maked/makefiles/premier/"
	// First we parse the makefile
	var g *Graph = GraphParser(makefilePath)
	// Now we execute all commands in the directory asked
	commandListe := []MakeElement{}
	launchMakefile(g, "", &commandListe)

	// Register the MakeService: This will allow the client to recognize it over the network
	makeService := &MakeService{
		Instructions: commandListe,
		Directory:    makefileDir,
	}
	rpc.Register(makeService)

	// Listen on TCP port 8090
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Println("Server listening on port 8090...")

	// TODO CEST LE DERNIER CLIENT QUI DIT AU SERVEUR DE SHUT DOWN
	// Goroutine to monitor `commandListe` and shut down the server
	go func() {
		for {
			time.Sleep(5 * time.Second) // Check every second
			makeService.mu.Lock()
			if len(makeService.Instructions) == 0 {
				makeService.mu.Unlock()
				log.Println("No more instructions. Shutting down the server...")
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
