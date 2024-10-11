package main

import (
  "fmt"
  "time"
)

/*
receive from ping channel (read only)
Then send to pong channel (write only)

*/

func pinger(ping_channel <-chan string, pong_channel chan<- string){
	for i := range ping_channel {
		fmt.Print("Received a ")
	  	fmt.Print(i)
		fmt.Println(", Sending a pong                                        ----> ")
		time.Sleep(time.Second) // Sleeping 1 sec for better display
		pong_channel <- "pong"
	}
}

/*
receive from pong channel (read  only)
Then send to ping channel (write only)
*/
func ponger(ping_channel chan<- string, pong_channel <-chan string) {
	for i := range pong_channel {
	  fmt.Print("Received a ")
	  fmt.Print(i)
	  fmt.Println(", Sending a ping                                        <---- ")
	  time.Sleep(time.Second) // Sleeping 1 sec for better display
	  ping_channel <- "ping"
	}
  }

func main() {
  ping := make(chan string)
  pong := make(chan string)

  // Player 1
  go pinger(ping, pong)
  // Player 2
  go ponger(ping, pong)

  // Player 1 starts the game
  fmt.Println("######### Sending a ping                                        <---- ")
  ping <- "ping"



  for {
  }
}
