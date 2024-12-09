package main

import (
	"io/fs"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

// fileSend holds the data sent from client to server.
type FileStruct struct {
	Data     []byte
	FileName string
}

type FileList struct {
	List []FileStruct
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

func removeAllFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}
	for _, file := range files {
		// log.Println("Deleting file ", file.Name())
		err = os.Remove(directory + file.Name())
		if err != nil {
			// log.Println("Impossible to delete the file ", err)
		}
	}
}

func diffFiles(filesBefore []fs.DirEntry, filesAfter []fs.DirEntry) []string {
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
	// Todo if the command delete a file (IMPORTANT)

	return diff
}

func createFile(storage string, fileName string, data []byte) {
	err := os.WriteFile(storage+fileName, data, 0644)
	if err != nil {
		panic(err)
	}
	err = os.Chmod(storage+fileName, 0755)
	if err != nil {
		log.Fatalf("Impossible to add permission")
	}
}

func launchCommand(storage string, command string) []string {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = storage

	// stdout and stderr directed to our terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Check all files before the command
	filesBefore, err := os.ReadDir(storage)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}
	// log.Println("Files before : ", filesBefore)

	// Execute the command
	err = cmd.Run()
	if err != nil {
		log.Println("Error:", err)
	}

	// We check all files after the command
	filesAfter, err := os.ReadDir(storage)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}
	// log.Println("Files after : ", filesAfter)

	return diffFiles(filesBefore, filesAfter)

}

func ask_init(storage string, address string) {
	// Connect to the server
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Prepare the content
	args := &PingDef{}
	var reply FileList
	err = client.Call("MakeService.Initialization", args, &reply)
	if err != nil {
		panic(err)
	}
	// log.Println("Downloading files")
	for _, file := range reply.List {
		createFile(storage, file.FileName, file.Data)
	}
}

func send_ping(address string) Order {
	// Connect to the server
	client, err := rpc.Dial("tcp", address)
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
	return reply
}

func send_file(directory string, filename string, address string) {
	// Connect to the server
	client, err := rpc.Dial("tcp", address)
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
}

func main() {
	var storage string = "./client_storage/"
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalf("Excepted 1 argument")
	}
	// we ask for essential files
	ask_init(storage, args[0])

forLoop:
	for {
		// Say to the server "hello I'm available"
		log.Println("Send ping")
		o := send_ping(args[0])
		log.Println("Pong received")
		// Server respond with an order
		// 0 -> no work available
		// 1 -> work available but waiting some jobs to end
		// 2 -> work available
		switch o.Value {
		case 0:
			log.Println("Task done, take some rest soldier")
			break forLoop
		case 1:
			log.Println("Server not ready")
		case 2:
			// log.Println("Ah shit, here we go again")
			// download all files
			// log.Println("Start of dependencies downloading")
			for _, dep := range o.Dependencies {
				createFile(storage, dep.FileName, dep.Data)
			}
			// log.Println("End of dependencies downloading")
			// execute the command
			// log.Println("Launching command")
			startTime := time.Now()
			filesCreated := launchCommand(storage, o.Command)
			elapsedTime := time.Since(startTime)
			log.Printf("Command done, execution time: %.2f seconds", elapsedTime.Seconds())
			// Send the created files
			// log.Println("Sending created files")
			for _, fileName := range filesCreated {
				send_file(storage, fileName, args[0])
			}
			// log.Println("Sended")
		}
	}
	removeAllFiles(storage)

}
