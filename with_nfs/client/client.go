package main

import (
	"io/fs"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

var id int = -1

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

type PingDef struct {
	ClientId int
}

type Order struct {
	Value        byte
	Command      string
	Dependencies []FileStruct
	Name         string
	ClientId     int
}

func removeAllFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}
	for _, file := range files {
		err = os.Remove(directory + file.Name())
		if err != nil {
			// log.Println("Impossible to delete the file ", err)
		}
	}
}

func diffFiles(filesBefore []fs.DirEntry, filesAfter []fs.DirEntry) []string {
	fileMap := make(map[string]bool)
	for _, file := range filesBefore {
		fileMap[file.Name()] = true
	}

	var diff []string
	for _, file := range filesAfter {
		if !fileMap[file.Name()] {
			diff = append(diff, file.Name())
		}
	}
	return diff
}

func launchCommand(storage string, command string) []string {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = storage
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	filesBefore, err := os.ReadDir(storage)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}

	err = cmd.Run()
	if err != nil {
		log.Println("Error:", err)
	}

	filesAfter, err := os.ReadDir(storage)
	if err != nil {
		log.Fatalf("Impossible to read the directory : %e", err)
	}

	return diffFiles(filesBefore, filesAfter)
}

func ask_init(address string) {
	client, err := rpc.Dial("tcp", address)
	// Retry loop to wait for the server to be available
	cpt := 0
	for {
		client, err = rpc.Dial("tcp", address)
		if cpt == 600000 && err != nil {
			panic(err)
		}
		if err == nil {
			break // Exit the loop if connection is successful
		}
		time.Sleep(1 * time.Millisecond) // Wait before retrying
		cpt++
	}
	defer client.Close()

	args := &PingDef{ClientId: id}
	var reply FileList
	err = client.Call("MakeService.Initialization", args, &reply)
	if err != nil {
		panic(err)
	}
	if id == -1 {
		id = reply.ClientId
	}
}

func send_ping(address string) Order {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	args := &PingDef{ClientId: id}
	var reply Order
	err = client.Call("MakeService.Ping", args, &reply)
	if err != nil {
		panic(err)
	}

	return reply
}

func send_update(filename string, codeValue JobReturn, address string) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	var args *FileStruct

	if codeValue.CodeValue == 0 {
		args = &FileStruct{FileName: filename, ReturnVal: codeValue, ClientId: id}
	} else {
		args = &FileStruct{FileName: codeValue.TargetName, ReturnVal: codeValue, ClientId: id}
	}

	var reply FileStruct
	err = client.Call("MakeService.SendUpdate", args, &reply)
	if err != nil {
		panic(err)
	}
}

func main() {
	var storage string = "../commun_storage/"
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalf("Excepted 1 argument")
	}
	ask_init(args[0])

forLoop:
	for {
		log.Println("Send ping")
		o := send_ping(args[0])
		log.Println("Pong received")
		switch o.Value {
		case 0:
			log.Println("Task done, take some rest soldier")
			break forLoop
		case 1:
			log.Println("Server not ready")
			time.Sleep(1 * time.Second)
		case 2:
			log.Printf("I'm the id %d and I work on : %s", id, o.Command)
			log.Println("Launching target: ", o.Name)
			startTime := time.Now()
			filesCreated := launchCommand(storage, o.Command)
			elapsedTime := time.Since(startTime)
			codeError := 0

			if len(filesCreated) == 0 {
				log.Println("NO FILE CREATED")
				info, err := os.Stat(storage + o.Name)
				if err != nil {
					log.Printf("File doesn't exist  %s, %v\n", o.Name, err)
					send_update("", JobReturn{CodeValue: 1, TargetName: o.Name}, args[0])
					break
				}
				if info.Size() == 0 {
					log.Printf("File is empty: %s. Retrying...\n", o.Name)
					send_update("", JobReturn{CodeValue: 2, TargetName: o.Name}, args[0])
					break
				}
			}

			for _, file := range filesCreated {
				info, err := os.Stat(storage + file)
				if err != nil {
					panic("Can't access to the file")
				}
				if info.Size() == 0 {
					send_update("", JobReturn{CodeValue: 3, TargetName: o.Name}, args[0])
					break
				}

				jobReturn := JobReturn{CodeValue: codeError, TargetName: o.Name}
				log.Printf("Command %s done, execution time: %.2f seconds", o.Name, elapsedTime.Seconds())
				for _, fileName := range filesCreated {
					send_update(fileName, jobReturn, args[0])
				}
			}
		}
	}
	removeAllFiles(storage)

}
