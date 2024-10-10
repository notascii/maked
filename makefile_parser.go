package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/** Le parser est vu comme un automate à plusieurs états :
*	état 0 : en attente d'une nouvelle cible
*	état 1 : en attente de dépendances
*	état 2 : en attand d'une commande
**/

func lineTreatment(s string, g *Graph, currentState int) {

}

// Hello returns a greeting for the named person.
func GraphParser(fileName string) *Graph {
	// First we init a graph
	g := &Graph{}
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
