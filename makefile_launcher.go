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
		// We check if the dependency already exists
		// TODO
		exploreGraph(g, dependency)
	}

	for _, command := range g.Vertices[target].command {
		fmt.Println("Voila la commande que j'execute")
		for _, dependence := range g.Vertices[target].dependencies {
			fmt.Println("Voila une dep " + dependence)
		}
		launchCommand(command)

	}

}
