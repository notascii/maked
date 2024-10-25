package main

import (
  "fmt"
  "time"
  "os"
)

/*
receive from ping channel (read only)
Then send to pong channel (write only)

*/

func pinger(ping_channel <-chan string, pong_channel chan<- string,file *os.File){
	for i := range ping_channel {
		timestamp := time.Now()
		log := fmt.Sprintf("Received a %s, Sending a pong ----> [%s]\n", i, timestamp)
		file.WriteString(log)
		time.Sleep(time.Second) // Sleeping 1 sec for better display
		pong_channel <- "pong"
	}
}

/*
receive from pong channel (read  only)
Then send to ping channel (write only)
*/
func ponger(ping_channel chan<- string, pong_channel <-chan string,file *os.File) {
	for i := range pong_channel {
		timestamp := time.Now()
		log := fmt.Sprintf("Received a %s, Sending a ping ----> [%s]\n", i, timestamp)
		file.WriteString(log)
	    time.Sleep(time.Second) // Sleeping 1 sec for better display
	    ping_channel <- "ping"
	}
  }

func main() {
  file,err := os.Create("pingpong_log.txt")
  if err != nil {
	fmt.Println("Error creating file:", err)
	return
}
  ping := make(chan string)
  pong := make(chan string)

  defer file.Close()

  // Player 1
  go pinger(ping, pong,file)
  // Player 2
  go ponger(ping, pong,file)

  // Player 1 starts the game
  file.WriteString("######### Sending a ping <---- [" + time.Now().Format("2006-01-02 15:04:05") + "]\n")
  ping <- "ping"



  for {
  }
}
