package main

import (
	"fmt"
	"time"
)

func main() {
	debut := time.Now()
	var g *Graph = GraphParser("./makefiles/premier/Makefile")
	duree := time.Since(debut)
	fmt.Println("Temps du parsing : ", duree.Seconds())
	// printVertices(*g)
	debut = time.Now()
	launchMakefile(g, "")
	duree = time.Since(debut)
	fmt.Println("Temps d'execution du graphe : ", duree.Seconds())

}
