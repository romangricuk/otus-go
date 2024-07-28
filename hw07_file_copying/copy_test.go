package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testSourceFile = ""

func TestMain(m *testing.M) {
	// Setup
	f, err := os.CreateTemp("", "source.*.txt")
	if err != nil {
		log.Fatal(err)
	}
	testSourceFile = f.Name()

	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	_, err = f.Write([]byte("hello world!"))
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestCopyFile(t *testing.T) {
	testCases := []struct {
		name       string
		from       string
		to         string
		offset     int64
		limit      int64
		expectErr  bool
		expectSize int64
	}{
		{
			name:       "Copy entire file",
			from:       testSourceFile,
			to:         "dest1.*.txt",
			offset:     0,
			limit:      0,
			expectErr:  false,
			expectSize: 12,
		},
		{
			name:       "Copy with offset",
			from:       testSourceFile,
			to:         "dest2.*.txt",
			offset:     6,
			limit:      0,
			expectErr:  false,
			expectSize: 6,
		},
		{
			name:       "Copy with limit",
			from:       testSourceFile,
			to:         "dest3.*.txt",
			offset:     0,
			limit:      6,
			expectErr:  false,
			expectSize: 6,
		},
		{
			name:      "Offset exceeds file size",
			from:      testSourceFile,
			to:        "dest4.*.txt",
			offset:    20,
			limit:     0,
			expectErr: true,
		},
		{
			name:       "Big file",
			from:       "testdata/input.txt",
			to:         "dest5.*.txt",
			offset:     0,
			limit:      0,
			expectErr:  false,
			expectSize: 6617,
		},
		{
			name:       "Big file with offset",
			from:       "testdata/input.txt",
			to:         "dest5.*.txt",
			offset:     10,
			limit:      0,
			expectErr:  false,
			expectSize: 6607,
		},
		{
			name:       "Big file with offset and limit",
			from:       "testdata/input.txt",
			to:         "dest5.*.txt",
			offset:     10,
			limit:      5000,
			expectErr:  false,
			expectSize: 5000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp("", tc.to)
			if err != nil {
				log.Fatal(err)
			}
			toFile := f.Name()

			defer func(name string) {
				_ = os.Remove(name)
			}(f.Name())

			err = copyFile(tc.from, toFile, tc.offset, tc.limit)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !tc.expectErr {
				info, err := os.Stat(toFile)
				if err != nil {
					t.Fatalf("failed to stat destination file: %v", err)
				}

				require.Equal(t, tc.expectSize, info.Size())
			}
		})
	}
}
