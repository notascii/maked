package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

func main() {
	// Retrieve the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Define paths
	makefilePath := homeDir + "/maked/makefiles/premier/Makefile-small"
	makefileDir := homeDir + "/maked/makefiles/premier/"

	// Parse the Makefile
	g := GraphParser(makefilePath)

	// Execute all commands in the specified directory
	commandListe := []MakeElement{}
	launchMakefile(g, "", &commandListe)

	// Initialize the MakeService
	makeService := &MakeService{
		Instructions: commandListe,
		Directory:    makefileDir,
	}

	// Register the MakeService
	rpc.Register(makeService)

	// Listen on TCP port 8090
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("Failed to listen on port 8090: %v", err)
	}
	defer listener.Close()
	log.Println("Server listening on port 8090...")

	// Channel to signal server shutdown
	shutdownChan := make(chan struct{})

	// Goroutine to monitor the command list
	go func() {
		for {
			time.Sleep(1 * time.Second) // Adjust the interval as needed
			makeService.mu.Lock()
			if len(makeService.Instructions) == 0 {
				makeService.mu.Unlock()
				log.Println("Command list is empty. Initiating server shutdown...")
				close(shutdownChan)
				return
			}
			makeService.mu.Unlock()
		}
	}()

	// Accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-shutdownChan:
					// Shutdown initiated, exit the loop
					return
				default:
					log.Printf("Failed to accept connection: %v", err)
					continue
				}
			}
			go rpc.ServeConn(conn)
		}
	}()

	// Wait for shutdown signal
	<-shutdownChan
	log.Println("Server has shut down gracefully.")
}
