package main

type Graph struct {
	Vertices    map[string]Vertex
	firstTarget string
}

type Vertex struct {
	cible       string
	dependances []string
	commmande   []string
}
