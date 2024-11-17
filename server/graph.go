package server

type Graph struct {
	Vertices    map[string]Vertex
	firstTarget string
}

type Vertex struct {
	target       string
	dependencies []string
	command      []string
}
