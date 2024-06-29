package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeoutArg := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeoutArg, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s: %v\n", address, err)
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "...Connection was closed by signal")
		_ = client.Close()
		os.Exit(0)
	}()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			_ = client.Close()
			os.Exit(1)
		}
	}()
	err = client.Receive()
	if err != nil {
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		_ = client.Close()
		os.Exit(1)
	}

	_ = client.Close()
}
