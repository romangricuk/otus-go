package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func copyFile(from, to string, offset, limit int64) error {
	sourceFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("on open source file: %w", err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("on get source file info: %w", err)
	}

	if offset > sourceInfo.Size() {
		return errors.New("offset exceeds file size")
	}

	if _, err := sourceFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("on seek in source file: %w", err)
	}

	bytesToCopy := calculateBytesToCopy(sourceInfo.Size(), offset, limit)

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("on create destination file: %w", err)
	}
	defer destFile.Close()

	return copyWithProgress(sourceFile, destFile, bytesToCopy)
}

func calculateBytesToCopy(fileSize, offset, limit int64) int64 {
	if limit == 0 || limit > fileSize-offset {
		return fileSize - offset
	}
	return limit
}

func copyWithProgress(src io.Reader, dst io.Writer, bytesToCopy int64) error {
	const bufSize = 1024
	buf := make([]byte, bufSize)
	var totalWritten int64

	startTime := time.Now()

	for totalWritten < bytesToCopy {
		bytesToRead := min(bufSize, bytesToCopy-totalWritten)
		n, err := src.Read(buf[:bytesToRead])
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("on read file: %w", err)
		}

		if n == 0 {
			break
		}

		n, err = dst.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("on write: %w", err)
		}

		totalWritten += int64(n)
		printProgress(totalWritten, bytesToCopy, startTime)
	}

	fmt.Print("\r")

	return nil
}

func printProgress(copied, total int64, startTime time.Time) {
	percentage := float64(copied) / float64(total) * 100
	elapsed := time.Since(startTime).Seconds()
	speed := float64(copied) / elapsed / 1024 // KB/s
	fmt.Printf("\rProgress: %.2f%% (%d/%d bytes) [%.2f KB/s]", percentage, copied, total, speed)
}
