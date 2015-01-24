package main

import (
	"os"
	"testing"
)

func TestProcessStdin(t *testing.T) {
	// Create a writer to send to the function...
	processStdin(os.Stdin, os.Stdout)
}
