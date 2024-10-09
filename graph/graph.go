package graph

import "fmt"

type Graph struct {
	Vertices map[string][]Vertex
}

type Vertex struct {
	cible     string
	commmande string
}

func (g Graph) AddVertex(cible string, commmande string, cap int) {
	v := Vertex{
		cible:     cible,
		commmande: commmande,
	}
	g.Vertices[cible] = make([]Vertex, 0, cap)
	fmt.Println(v)

}
