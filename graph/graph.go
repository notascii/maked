package graph

type Graph struct {
	Vertices map[Vertex][]Vertex
}

type Vertex struct {
	cible       string
	dependances string
	commmande   string
}

func (g Graph) AddVertex(cible string, commmande string, cap int) {
	v := Vertex{
		cible:       cible,
		commmande:   commmande,
		dependances: "lol",
	}
	g.Vertices[v] = make([]Vertex, 0, cap)
}
