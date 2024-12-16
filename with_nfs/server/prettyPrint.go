package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func writeResults(makeDuration time.Duration, clientJobs map[int][]Job, makedDuration time.Duration, nfsDirectory string, fileName string) {
	// Expand the home directory if present
	if strings.HasPrefix(nfsDirectory, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error fetching home directory: %v\n", err)
			return
		}
		nfsDirectory = filepath.Join(homeDir, nfsDirectory[1:], "server/json_storage/")
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(nfsDirectory, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", nfsDirectory, err)
		return
	}

	// Prepare the output file path
	filePath := filepath.Join(nfsDirectory, fileName)

	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// Prepare the data structure for JSON encoding
	output := make(map[string]interface{})
	output["makeDuration"] = makeDuration.Microseconds()
	output["makedDuration"] = makedDuration.Microseconds()
	output["clients"] = map[string]interface{}{}

	clients := output["clients"].(map[string]interface{})

	for clientID, jobs := range clientJobs {
		clientData := map[string]interface{}{
			"totalDuration": 0,
			"jobs":          []map[string]interface{}{},
		}

		clientDuration := time.Duration(0)
		jobList := clientData["jobs"].([]map[string]interface{})

		for _, job := range jobs {
			clientDuration += time.Duration(job.Duration) * time.Microsecond
			jobList = append(jobList, map[string]interface{}{
				"name":     job.Name,
				"duration": job.Duration,
			})
		}

		clientData["totalDuration"] = clientDuration.Microseconds()
		clientData["jobs"] = jobList
		clients[fmt.Sprintf("%d", clientID)] = clientData
	}

	// Write the JSON data to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indentation
	if err := encoder.Encode(output); err != nil {
		fmt.Printf("Error encoding JSON to file %s: %v\n", filePath, err)
		return
	}

	fmt.Printf("Client list written successfully to %s\n", filePath)
}