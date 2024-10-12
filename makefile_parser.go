package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	VARIABLE = iota
	TARGET
	COMMAND
	IGNORED
	UNKNOW
)

/**
* 	Add a target line inside the graph
*	example :
*	target : test1 test2
*	add target as a vertex and link the dependances with test1 and test2
*
*
*	s : the line
*   *g : the graph
**/
func targetLoad(s string, g *Graph, currentTarget *string) {
	// Treatment of the string to separate the target name and dependencies
	res := strings.Split(s, ":")
	targetName := strings.Replace(res[0], " ", "", -1)
	targetDependencies := strings.Split(res[1], " ")
	// Initialization of a new vertex
	var newV = Vertex{
		cible:       targetName,
		commmande:   make([]string, 0, 256),
		dependances: targetDependencies,
	}
	// We add it to the vertex list
	g.Vertices[targetName] = newV
	// We init the values of its adjency list
	*currentTarget = targetName
}

/*
*
*	Add a command to a vertex which got a target s
*
 */
func commandLoad(s string, g *Graph, currentTarget *string) {
	fmt.Println("Voici la commande que je charge : " + s)
	v := g.Vertices[*currentTarget]
	v.commmande = append(v.commmande, s)
	g.Vertices[*currentTarget] = v
}

/*
*
*	return the type of a string.
*
 */
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
		panic("Makefile incorrect")
	}
}

/** Le parser est vu comme un automate à plusieurs états (à ajouter en fonction du type de makefile traité):
*	état 0 : en attente d'une nouvelle cible + dépendances OU d'une définition de variable
*	état 1 : en attente d'une commande (cible chargée)
*   état 2 : en attente d'une commande où d'une nouvelle cible ou d'une definition de variable
**/

func lineTreatment(s string, g *Graph, currentState int, currentTarget *string) int {
	println("Ligne traitrée : " + s)
	// Empty line
	buffer := lineType(s)
	if buffer == IGNORED {
		return currentState
	}
	switch currentState {
	case 0:
		// si s est une définition de variable on reste dans le cas 0
		if buffer == VARIABLE {
			return 0
		} else if buffer == TARGET { // Sinon on va dans le cas 1
			targetLoad(s, g, currentTarget)
			return 1
		} else {
			return -1
		}
	case 1:
		if buffer == COMMAND {
			commandLoad(s, g, currentTarget)
			return 2
		} else {
			return -1
		}
	case 2:
		if buffer == TARGET {
			targetLoad(s, g, currentTarget)
			return 1
		} else if buffer == VARIABLE {
			return 0
		} else if buffer == COMMAND {
			commandLoad(s, g, currentTarget)
			return 2
		}
	}

	return currentState

}

// GraphParser return a graph representing
func GraphParser(fileName string) *Graph {
	// First we init a graph
	g := &Graph{
		Vertices: make(map[string]Vertex),
	}

	// We read fileName line per line
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanLines)
	currentTargetVal := ""
	currentTarget := &currentTargetVal
	currentState := 0
	for Scanner.Scan() {
		currentState = lineTreatment(Scanner.Text(), g, currentState, currentTarget)
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return g
}
