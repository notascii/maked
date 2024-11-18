package main

import (
	"net"
	"net/rpc"
)

func main() {
	// Makefile treatment
	////////////// TEST PREMIER
	makefilePath := "../makefiles/matrix/Makefile"
	makefileDir := "../makefiles/matrix/"
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

	// Listen on TCP port 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	println("Server listening on port 1234...")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}
}
