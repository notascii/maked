package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

const (
	VARIABLE = iota
	TARGET
	COMMAND
	IGNORED
	UNKNOW
)

func lineType(s string) int {
	var varDefinition = regexp.MustCompile(`^\s*[0-9a-zA-Z_&]+\s*=.*$`)
	var targetDefinition = regexp.MustCompile(`^\s*[0-9a-zA-Z_&\.]+\s*:.*$`)
	var commandDefinition = regexp.MustCompile(`^\t.*$`)
	var emptyLine = regexp.MustCompile(`^(#.*)$|^(\s*)$`)
	if varDefinition.MatchString(s) {
		println("Définition de variable")
		return VARIABLE
	} else if targetDefinition.MatchString(s) {
		println("Définition de target")
		return TARGET
	} else if commandDefinition.MatchString(s) {
		println("Définition d'une commande")
		return COMMAND
	} else if emptyLine.MatchString(s) {
		println("Ligne ignorée")
		return IGNORED
	} else {
		println("ligne inconnu")
		return UNKNOW
	}
}

/** Le parser est vu comme un automate à plusieurs états (à ajouter en fonction du type de makefile traité):
*	état 0 : en attente d'une nouvelle cible + dépendances OU d'une définition de variable
*	état 1 : en attente d'une commande (cible chargée)
*   état 2 : en attente d'une commande où d'une nouvelle cible
**/

func lineTreatment(s string, g *Graph, currentState int, currentCible string) int {
	println("Ligne traitrée : " + s)
	// Empty line
	buffer := lineType(s)

	switch currentState {
	case 0:
		// si s est une définition de variable on reste dans le cas 0
		if buffer == VARIABLE {
			// todo
			return 0
		} else if buffer == TARGET { // Sinon on va dans le cas 1
			// todo
			return 1
		} else {
			return -1
		}
	case 1:
		if buffer == COMMAND {
			return 2
		} else {
			return -1
		}
	case 2:
		if buffer == COMMAND {
			return 2
		} else if buffer == TARGET {
			return 1
		} else {
			return -1
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
		lineTreatment(Scanner.Text(), g, 0, "")
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return g
}
