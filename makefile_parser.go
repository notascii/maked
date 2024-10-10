package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func emptyLine(s string) bool {
	for c := range s {
		if c != '\t' && c != ' ' && c != '\n' {
			return false
		}
	}
	return true
}

/** Le parser est vu comme un automate à plusieurs états (à ajouter en fonction du type de makefile traité):
*	état 0 : en attente d'une nouvelle cible + dépendances OU d'une définition de variable
*	état 1 : en attente d'une commande (cible chargée)
*   état 2 : en attente d'une commande où d'une nouvelle cible
**/

func lineTreatment(s string, g *Graph, currentState int, currentCible string) int {
	// Comment
	if s[0] == '#' {
		return currentState
	}
	// Empty line
	if emptyLine(s) {
		return currentState
	}

	switch currentState {
	case 0:
		// si s est une définition de variable on reste dans le cas 0
		if true {
			// todo
			return 0
		} else { // Sinon on va dans le cas 1
			// todo
			return 1
		}
	case 1:
		// todo
		return 2
	case 2:
		//todo
		if /* le string est une commande*/ true {
			return 2
		} else // nouvelle cible
		{
			return 1
		}
	}

	return currentState

}

// GraphParser return a graph representing
func GraphParser(fileName string) *Graph {
	// First we init a graph
	g := &Graph{}

	// We read fileName line per line
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)

	for Scanner.Scan() {
		fmt.Println(Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return g
}
