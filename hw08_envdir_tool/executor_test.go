package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": {Value: "123", NeedRemove: false},
		"BAR": {Value: "value", NeedRemove: false},
	}

	cmd := []string{"sh", "-c", "echo $FOO $BAR"}

	// Создаём временный файл для захвата Stdout
	tmpfile, err := os.CreateTemp("", "stdout")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Сохраняем исходный Stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Переводим stdout в наш временный файл
	os.Stdout = tmpfile

	returnCode := RunCmd(cmd, env)

	// Восстанавливаем Stdout
	os.Stdout = originalStdout

	err = tmpfile.Close()
	require.NoError(t, err)

	// Читаем то что попало во временный "stdout"
	output, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)

	// Проверяем ожидаемые значения
	assert.Equal(t, 0, returnCode)

	outputStr := strings.TrimSpace(string(output))
	assert.Equal(t, "123 value", outputStr)
}
