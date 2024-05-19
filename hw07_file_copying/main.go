package main

import (
	"errors"
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func validateFlags() error {
	if from == "" {
		return errors.New("source file path (-from) is required")
	}
	if to == "" {
		return errors.New("destination file path (-to) is required")
	}
	if offset < 0 {
		return errors.New("offset cannot be negative")
	}
	if limit < 0 {
		return errors.New("limit cannot be negative")
	}
	return nil
}

func main() {
	flag.Parse()
	// time.Sleep(20 * time.Second)
	if err := validateFlags(); err != nil {
		fmt.Println(err)
		return
	}

	if err := copyFile(from, to, offset, limit); err != nil {
		fmt.Println(err)
	}
}
