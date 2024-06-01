package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dir := os.TempDir()
	dir = filepath.Join(dir, "TestReadDir")
	err := os.MkdirAll(dir, 0o755)
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	files := map[string]string{
		"FOO":   "123",
		"BAR":   "value",
		"EMPTY": "",
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	env, err := ReadDir(dir)
	require.NoError(t, err)

	expected := Environment{
		"FOO":   {Value: "123", NeedRemove: false},
		"BAR":   {Value: "value", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
	}

	assert.Equal(t, expected, env)
}

func TestInvalidEnvVarName(t *testing.T) {
	dir := os.TempDir()
	dir = filepath.Join(dir, "TestReadDir")
	err := os.MkdirAll(dir, 0o755)
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.WriteFile(filepath.Join(dir, "INVALID=NAME"), []byte("value"), 0o644)
	require.NoError(t, err)

	_, err = ReadDir(dir)
	assert.Error(t, err)
}
