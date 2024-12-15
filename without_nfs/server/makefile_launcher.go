package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MakeElement struct {
	Command      string
	Dependencies []string
	Name         string
}

// Put inside the commandList all command that should be executed by doing a depth search in the dependencies graph
func launchMakefile(g *Graph, firstTarget string, commandList *[]MakeElement) {
	visited := make(map[string]bool) // To track visited vertices
	exploreGraph(g, firstTarget, commandList, visited)
}

// depth search of the graph
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
func launchClassicMake(directory string) time.Duration {
	targetDir := "./storage_for_make"

	// Remove existing target directory if it exists
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(targetDir); err != nil {
			fmt.Printf("Failed to remove existing directory %s: %v\n", targetDir, err)
			return 0
		}
	}

	// Create a fresh directory
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", targetDir, err)
		return 0
	}

	// Copy all files from the source `directory` to `./storage_for_make`
	if err := copyDir(directory, targetDir); err != nil {
		fmt.Printf("Failed to copy from %s to %s: %v\n", directory, targetDir, err)
		return 0
	}

	// Start timer
	startTime := time.Now()

	// Run `make` in `./storage_for_make`
	cmd := exec.Command("make")
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Make command failed: %v\n", err)
		return 0
	}

	// End timer and calculate duration
	duration := time.Since(startTime)

	return duration
}

// copyDir recursively copies a directory from src to dst
func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Construct target path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create the directory
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err := copyFile(path, targetPath, info.Mode()); err != nil {
				return err
			}
		}
		return nil
	})
}

// copyFile copies a single file from src to dst with given permission mode
func copyFile(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}
	return nil
}
