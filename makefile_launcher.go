package main

import (
	"fmt"
	"log"
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

func launchMakefile(g *Graph, firstTarget string, directory string) {

	targetDir := directory
	err := os.Chdir(targetDir)
	if err != nil {
		log.Fatalf("Error while changing repo : %v", err)
	}

	exploreGraph(g, firstTarget)
}

func exploreGraph(g *Graph, target string) {
	if target == "" {
		target = g.firstTarget
	}

	for _, dependency := range g.Vertices[target].dependencies {
		// Case where dependencies is a file
		// TODO
		// Case where dependencies is a a target / nothing
		exploreGraph(g, dependency)
	}

	for _, command := range g.Vertices[target].command {
		launchCommand(command)
	}

}
