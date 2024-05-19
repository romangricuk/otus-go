package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func copyFile(from, to string, offset, limit int64) (err error) {
	sourceFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("on open source file: %w", err)
	}
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

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

	var bytesToCopy int64

	if limit == 0 || limit-offset > sourceInfo.Size() {
		bytesToCopy = sourceInfo.Size() - offset
	} else if limit > 0 {
		bytesToCopy = limit
	}

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("on create destination file: %w", err)
	}

	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	return copyWithProgress(sourceFile, destFile, bytesToCopy)
}

func copyWithProgress(src io.Reader, dst io.Writer, bytesToCopy int64) error {
	bufSize := int64(1024)
	buf := make([]byte, bufSize)
	var totalWrite int64
	var totalRead int64
	var toByte int64

	startTime := time.Now()

	for {
		numOfReadBytes, err := src.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("on read file: %w", err)
		}

		totalRead += int64(numOfReadBytes)

		if totalWrite+int64(numOfReadBytes) > bytesToCopy {
			if bufSize >= bytesToCopy {
				toByte = bytesToCopy
			} else {
				toByte = bufSize - (totalWrite + int64(numOfReadBytes) - bytesToCopy)
			}
		} else {
			toByte = int64(numOfReadBytes)
		}

		numOfWriteBytes, err := dst.Write(buf[0:toByte])
		if err != nil {
			return fmt.Errorf("on write: %w", err)
		}

		totalWrite += int64(numOfWriteBytes)

		if toByte != int64(numOfWriteBytes) {
			return errors.New("read/write mismatch")
		}

		printProgress(totalWrite, bytesToCopy, startTime)

		if totalWrite >= bytesToCopy {
			break
		}
	}

	if totalWrite > 0 {
		fmt.Print("\r")
	}

	return nil
}

func printProgress(copied, total int64, startTime time.Time) {
	percentage := float64(copied) / float64(total) * 100
	elapsed := time.Since(startTime).Seconds()
	speed := float64(copied) / elapsed / 1024 // KB/s
	fmt.Printf("\rProgress: %.2f%% (%d/%d bytes) [%.2f KB/s]", percentage, copied, total, speed)
}
