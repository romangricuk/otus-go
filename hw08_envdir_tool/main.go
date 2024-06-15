package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2:]
	envVars, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Error reading env dir: %v\n", err)
		os.Exit(1)
	}

	exitCode := RunCmd(command, envVars)
	os.Exit(exitCode)
}
