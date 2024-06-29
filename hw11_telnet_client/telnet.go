package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return fmt.Errorf("on connecting: %w", err)
	}
	tc.conn = conn
	return nil
}

func (tc *telnetClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *telnetClient) Send() error {
	scanner := bufio.NewScanner(tc.in)
	for scanner.Scan() {
		_, err := tc.conn.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return fmt.Errorf("on sending data: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("on reading from input: %w", err)
	}
	return nil
}

func (tc *telnetClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	if err != nil {
		return fmt.Errorf("on receiving data: %w", err)
	}
	return nil
}
