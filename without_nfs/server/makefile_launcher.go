package main

type MakeElement struct {
	Command      string
	Dependencies []string
	Name         string
}

func launchMakefile(g *Graph, firstTarget string, commandList *[]MakeElement) {
	visited := make(map[string]bool) // To track visited vertices
	exploreGraph(g, firstTarget, commandList, visited)
}

func exploreGraph(g *Graph, target string, commandList *[]MakeElement, visited map[string]bool) {
	if target == "" {
		target = g.firstTarget
	}

	// If the vertex has already been visited, skip it
	if visited[target] {
		return
	}
	visited[target] = true // Mark the current vertex as visited

	for _, dependency := range g.Vertices[target].dependencies {
		// Explore each dependency
		exploreGraph(g, dependency, commandList, visited)
	}

	for _, command := range g.Vertices[target].command {
		// Add the command to the list
		ins := MakeElement{Command: command, Dependencies: g.Vertices[target].dependencies, Name: g.Vertices[target].target}
		*commandList = append(*commandList, ins)
	}
}
