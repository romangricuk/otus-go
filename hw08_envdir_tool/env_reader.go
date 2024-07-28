package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (envVars Environment, err error) {
	envVars = make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("on reading directory %s %w", dir, err)
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			err = fmt.Errorf("invalid environment variable name: %s", name)
			return nil, err
		}

		path := filepath.Join(dir, name)
		value, err := readFirstLine(path)
		if err != nil {
			err = fmt.Errorf("on read file first line %s %w", path, err)
			return nil, err
		}

		envVars[name] = EnvValue{Value: value, NeedRemove: value == ""}
	}

	return envVars, nil
}

func readFirstLine(path string) (line string, err error) {
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("on open file %w", err)
		return "", err
	}

	defer func(file *os.File) {
		errOnClose := file.Close()
		if errOnClose != nil {
			errOnClose = fmt.Errorf("on close file %w", errOnClose)
			err = errors.Join(err, errOnClose)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = scanner.Text()
		line = strings.TrimRight(line, " \t")
		line = strings.ReplaceAll(line, "\x00", "\n")
		return line, err
	}
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("on scan first line: %w", err)
		return "", err
	}

	return "", err
}
