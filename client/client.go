package main

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
)

// fileSend holds the data sent from client to server.
type FileStruct struct {
	Data     []byte
	FileName string
}

type Message struct {
	Msg string
}

type PingDef struct {
}

type Order struct {
	Value        byte
	Command      string
	Dependencies []FileStruct
}

func launchCommand(command string) {
	fmt.Println(command)
	cmd := exec.Command("/bin/sh", "-c", command)

	// stdout and stderr directed to our terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}

}

func send_ping() Order {
	// Connect to the server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Prepare the content
	args := &PingDef{}
	var reply Order
	err = client.Call("MakeService.Ping", args, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ping send")
	fmt.Println("Order received :", reply.Value)
	return reply
}

func send_file(directory string, filename string) {
	// Connect to the server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Read the file content
	fileData, err := os.ReadFile(directory + filename)
	if err != nil {
		panic(err)
	}

	// Send the file content
	args := &FileStruct{Data: fileData, FileName: filename} // N-byte message
	var reply FileStruct
	err = client.Call("MakeService.SendFile", args, &reply)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File send\n")
}

func main() {
forLoop:
	for {
		// Say to the server "hello I'm available"
		o := send_ping()
		// Server respond with an order
		// 0 -> no work available
		// 1 -> work available but waiting some jobs to end
		// 2 -> work available
		switch o.Value {
		case 0:
			fmt.Println("Task done, take some rest soldier")
			break forLoop
		case 1:
			fmt.Println("Server not ready")
		case 2:
			fmt.Println("Ah shit here we go again")
			// download all files
			// execute the command
			launchCommand(o.Command)
			// send the created file
			send_file("./client_storage/", "test.txt")
		}
	}
}
