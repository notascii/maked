package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type MakeElement struct {
	Command      string
	Dependencies []string
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

func launchMakefile(g *Graph, firstTarget string, directory string, commandList *[]MakeElement) {

	targetDir := directory
	err := os.Chdir(targetDir)
	if err != nil {
		log.Fatalf("Error while changing repo : %v", err)
	}
	exploreGraph(g, firstTarget, commandList)
}

func exploreGraph(g *Graph, target string, commandList *[]MakeElement) {
	if target == "" {
		target = g.firstTarget
	}

	for _, dependency := range g.Vertices[target].dependencies {
		// We check if the dependency already exists
		// TODO
		exploreGraph(g, dependency, commandList)
	}

	for _, command := range g.Vertices[target].command {
		fmt.Println("Voila la commande que j'execute")
		launchCommand(command)
		ins := MakeElement{Command: command, Dependencies: g.Vertices[target].dependencies}
		*commandList = append(*commandList, ins)

	}

}
