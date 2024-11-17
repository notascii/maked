package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/rpc"
	"os"
	"os/exec"
)

var path string = "./client_storage/"

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

func DiffFiles(filesBefore []fs.DirEntry, filesAfter []fs.DirEntry) []string {
	// Map on files2 name
	fileMap := make(map[string]bool)
	for _, file := range filesBefore {
		fileMap[file.Name()] = true
	}

	// Now we check the differences
	var diff []string
	for _, file := range filesAfter {
		if !fileMap[file.Name()] {
			diff = append(diff, file.Name())
		}
	}
	// Todo if the command delete a file (idk if it even supposed to occur)

	return diff
}

func createFile(fileName string, data []byte) {
	err := os.WriteFile(path+fileName, data, 0644)
	if err != nil {
		panic(err)
	}
}

func launchCommand(command string) []string {

	fmt.Println(command)
	cmd := exec.Command("/bin/sh", "-c", command)

	// stdout and stderr directed to our terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Check all files before the command
	filesBefore, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Impossible to read the directory : ", err)
	}
	fmt.Println("Files before : ", filesBefore)

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// We check all files after the command
	filesAfter, err := os.ReadDir(".")
	if err != nil {
		log.Println("Impossible to read the directory : ", err)
	}
	fmt.Println("Files after : ", filesAfter)

	return DiffFiles(filesBefore, filesAfter)

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
	fileData, err := os.ReadFile(filename)
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
	// change dir
	err := os.Chdir(path)
	if err != nil {
		log.Fatalf("Error while changing repo : %v", err)
	}
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
			fmt.Println("Ah shit, here we go again")
			// download all files
			for _, dep := range o.Dependencies {
				createFile(dep.FileName, dep.Data)
			}
			// execute the command
			filesCreated := launchCommand(o.Command)
			fmt.Println("Command done")
			// Send the created files
			fmt.Println("Created files : ", filesCreated)
			for _, fileName := range filesCreated {
				send_file(path, fileName)
			}
			fmt.Println("File sended")

		}
	}
}
