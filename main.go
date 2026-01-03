package main

import (
	"fmt"
	"io"
	"os"

	"github.com/giovannirossini/markdown-render/render"
)

const exitCodeSuccess = 0
const exitCodeError = 1

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "mdrender: %v\n", err)
		os.Exit(exitCodeError)
	}
	os.Exit(exitCodeSuccess)
}

func run() error {
	if len(os.Args) < 2 {
		// Read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}
		render.Render(string(content))
		return nil
	}

	input := os.Args[1]

	// Try to read as file
	content, err := os.ReadFile(input)
	if err != nil {
		// Treat as direct markdown input if file doesn't exist
		render.Render(input)
		return nil
	}

	render.Render(string(content))
	return nil
}
