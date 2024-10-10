package main

type Graph struct {
	AdjencyList map[string][]Vertex
	Vertices    map[string]Vertex
}

type Vertex struct {
	cible       string
	dependances string
	commmande   []string
}

func (g *Graph) AddVertex(cible string, commmande []string, dependances string, cap int) {
	v := Vertex{
		cible:       cible,
		commmande:   commmande,
		dependances: dependances,
	}
	g.Vertices[cible] = v
}
