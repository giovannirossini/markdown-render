package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// Read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		Render(string(content))
		return
	}

	input := os.Args[1]

	// Try to read as file
	content, err := os.ReadFile(input)
	if err != nil {
		// Treat as direct markdown input
		Render(input)
		return
	}

	Render(string(content))
}
