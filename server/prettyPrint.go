package server

import "fmt"

func printVertex(v Vertex) {
	fmt.Printf("%s : ", v.target)
	for _, value := range v.dependencies {
		fmt.Printf("%s ", value)
	}
	fmt.Printf("\n\t")
	for _, value := range v.command {
		fmt.Printf("%s ", value)
	}
	fmt.Printf("\n")
}

func printVertices(g Graph) {
	for _, value := range g.Vertices {
		// fmt.Printf("clé : %s\n", key)
		printVertex(value)
	}
}
