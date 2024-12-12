package main

import (
	"fmt"
	"time"
)

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

func printClientList(clientJobs map[int][]Job, totalDuration time.Duration) {
	fmt.Println("{")

	// Print the total duration
	fmt.Printf("  \"totalDuration\": %d,\n", totalDuration.Milliseconds())

	// Print durations per client
	fmt.Println("  \"clients\": {")
	clientCount := len(clientJobs)
	currentIndex := 0

	for clientID, jobs := range clientJobs {
		currentIndex++
		fmt.Printf("    \"%d\": {\n", clientID)

		// Calculate total duration for the client
		clientDuration := time.Duration(0)
		for _, job := range jobs {
			clientDuration += time.Duration(job.Duration) * time.Millisecond
		}

		fmt.Printf("      \"totalDuration\": %d,\n", clientDuration.Milliseconds())
		fmt.Println("      \"jobs\": [")

		// Print each job for the client
		for i, job := range jobs {
			fmt.Printf("        {\"name\": \"%s\", \"duration\": %d}", job.Name, job.Duration)
			if i < len(jobs)-1 {
				fmt.Print(",") // Add a comma if not the last job
			}
			fmt.Println()
		}

		fmt.Print("      ]\n")
		fmt.Print("    }")
		if currentIndex < clientCount {
			fmt.Print(",") // Add a comma if not the last client
		}
		fmt.Println()
	}

	fmt.Println("  }")
	fmt.Println("}")
}
