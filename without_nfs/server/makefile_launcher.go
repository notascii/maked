package main

type MakeElement struct {
	Command      string
	Dependencies []string
}

func launchMakefile(g *Graph, firstTarget string, commandList *[]MakeElement) {

	exploreGraph(g, firstTarget, commandList)
}

func exploreGraph(g *Graph, target string, commandList *[]MakeElement) {
	if target == "" {
		target = g.firstTarget
	}

	for _, dependency := range g.Vertices[target].dependencies {
		// We check if the dependency already exists
		// TODO
		exploreGraph(g, dependency, commandList)
	}

	for _, command := range g.Vertices[target].command {
		// launchCommand(command)
		ins := MakeElement{Command: command, Dependencies: g.Vertices[target].dependencies}
		*commandList = append(*commandList, ins)

	}

}
