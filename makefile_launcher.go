package main

import (
	"fmt"
	"os"
	"os/exec"
)

func launchCommand(command string) {
	fmt.Println(command)
	cmd := exec.Command("/bin/sh", "-c", command)

	// Diriger la sortie standard et les erreurs vers le terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Exécuter la commande
	err := cmd.Run()
	if err != nil {
		fmt.Println("Erreur:", err)
	}
}

func launchMakefile(g *Graph, firstTarget string) {
	if firstTarget == "" {
		firstTarget = g.firstTarget
	}

	for _, value := range g.Vertices[firstTarget].dependances {
		// Cas où la dépendance est un fichier
		// TODO
		// Cas où la dépendance est une autre cible
		launchMakefile(g, value)
	}

	for _, command := range g.Vertices[firstTarget].commmande {
		launchCommand(command)
	}

}
