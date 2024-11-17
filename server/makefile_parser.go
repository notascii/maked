package server

import (
	"bufio"
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
	UNKNOWN
)

var firstFound = false

/*
* 	Add a target line inside the graph
*	example :
*	target : test1 test2
*	add target as a vertex and link the dependencies with test1 and test2
*
*
*	s : the line
*   g : the graph
 */
func targetLoad(s string, g *Graph, currentTarget *string) {
	// Treatment of the string to separate the target name and dependencies
	res := strings.Split(s, ":")
	targetName := strings.Replace(res[0], " ", "", -1)
	targetDependencies := strings.Fields(res[1])
	// Initialization of a new vertex
	var newV = Vertex{
		target:       targetName,
		command:      make([]string, 0, 256),
		dependencies: targetDependencies,
	}
	// We add it to the vertex list
	g.Vertices[targetName] = newV
	// We init the values of its adjacency list
	*currentTarget = targetName
	if !firstFound {
		g.firstTarget = targetName
		firstFound = true
	}
}

/*
*
*	Add a command to a vertex which got a target s
*
 */
func commandLoad(s string, g *Graph, currentTarget *string) {
	v := g.Vertices[*currentTarget]
	v.command = append(v.command, s)
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
		return VARIABLE
	} else if targetDefinition.MatchString(s) {
		return TARGET
	} else if commandDefinition.MatchString(s) {
		return COMMAND
	} else if emptyLine.MatchString(s) {
		return IGNORED
	} else {
		panic("Makefile incorrect")
	}
}

/** We have an automaton with three states:
*	state 0 : waiting for a new target (0 -> 1) OR a variable definition (0 -> 0)
*	state 1 : waiting for a command (target already load) (1 -> 2)
*   state 2 : waiting for a command (2 -> 2) OR a target (2 -> 1) OR a variable definition (2 -> 0)
**/

func lineTreatment(s string, g *Graph, currentState int, currentTarget *string) int {
	// Empty line
	buffer := lineType(s)
	if buffer == IGNORED {
		return currentState
	}
	switch currentState {
	case 0:
		if buffer == VARIABLE {
			return 0
		} else if buffer == TARGET {
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

/*
*	GraphParser return a graph representing the dependencies graph of the makefile in fileName
*
 */
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
