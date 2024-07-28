package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "develop"             //nolint:unused
	buildDate = "2024-07-24T13:22:15" //nolint:unused
	gitHash   = "8118ea7"             //nolint:unused
)

func printVersion() { //nolint:unused
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
