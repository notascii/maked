package main

import (
	"fmt"

	"example.com/parsing"
)

func main() {
	// Get a greeting message and print it.
	message := parsing.Hello("Gladys")
	fmt.Println(message)
}
