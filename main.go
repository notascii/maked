package main

func main() {

	var g *Graph = GraphParser("./resources/Makefile")
	printVertices(*g)
}
