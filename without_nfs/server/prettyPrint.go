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

func writeResults(
	makeDuration time.Duration,
	clientJobs map[int][]Job,
	makedDuration time.Duration,
	nfsDirectory string,
	makefileName string,
	numberOfNodes string,
) {
	// Expand the home directory if present
	if strings.HasPrefix(nfsDirectory, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error fetching home directory: %v\n", err)
			return
		}
		// Remove the '~' and join with homeDir
		relPath := strings.TrimPrefix(nfsDirectory, "~")
		nfsDirectory = filepath.Join(homeDir, relPath)
	}

	// Always append "server/json_storage" to ensure a consistent directory structure
	nfsDirectory = filepath.Join(nfsDirectory, "server", "json_storage")

	// Print out the directory we are about to create
	fmt.Printf("Creating base directory: %s\n", nfsDirectory)
	if err := os.MkdirAll(nfsDirectory, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", nfsDirectory, err)
		return
	}

	makefileDir := filepath.Join(nfsDirectory, makefileName)
	fmt.Printf("Creating makefile directory: %s\n", makefileDir)
	if err := os.MkdirAll(makefileDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", makefileDir, err)
		return
	}

	// Define the file path as numberOfNodes.json inside the makefileDir
	filePath := filepath.Join(makefileDir, numberOfNodes+".json")
	fmt.Printf("Creating file: %s\n", filePath)

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

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(output); err != nil {
		fmt.Printf("Error encoding JSON to file %s: %v\n", filePath, err)
		return
	}

	fmt.Printf("Client list written successfully to %s\n", filePath)
}
