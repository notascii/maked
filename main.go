package main

func main() {
	////////////// TEST PREMIER
	makefilePath := "./makefiles/premier/Makefile"
	makefileDirectory := "./makefiles/premier"
	// First we parse the makefile
	var g *Graph = GraphParser(makefilePath)
	// Now we execute all commands in the directory asked
	launchMakefile(g, "", makefileDirectory)

}
