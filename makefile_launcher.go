package main

import (
	"fmt"
	"os"
	"os/exec"
)

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

func launchMakefile(g *Graph, firstTarget string) {
	if firstTarget == "" {
		firstTarget = g.firstTarget
	}

	for _, value := range g.Vertices[firstTarget].dependencies {
		// Case where dependencies is a file
		// TODO
		// Case where dependencies is a a target / nothing
		launchMakefile(g, value)
	}

	// We create a bash

	for _, command := range g.Vertices[firstTarget].command {
		launchCommand(command)
	}

}
