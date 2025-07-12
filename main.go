package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/javanhut/CarrionLSP/internal/server"
)

// stdioPipe combines stdin and stdout into a ReadWriteCloser
type stdioPipe struct{}

func (s *stdioPipe) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (s *stdioPipe) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (s *stdioPipe) Close() error {
	// Don't actually close stdin/stdout
	return nil
}

func main() {
	var conn *jsonrpc2.Conn
	var err error

	// Check for command line arguments
	if len(os.Args) > 1 && os.Args[1] == "--stdio" {
		// Use stdio transport
		stream := jsonrpc2.NewPlainObjectStream(&stdioPipe{})
		conn = jsonrpc2.NewConn(
			context.Background(),
			stream,
			server.NewHandler(),
		)
	} else {
		// Use TCP transport (default port 7777)
		port := "7777"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", port, err)
		}
		defer listener.Close()

		fmt.Printf("Carrion LSP server listening on port %s\n", port)

		for {
			netConn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			stream := jsonrpc2.NewPlainObjectStream(netConn)
			conn = jsonrpc2.NewConn(
				context.Background(),
				stream,
				server.NewHandler(),
			)

			// Handle one connection at a time for now
			<-conn.DisconnectNotify()
		}
	}

	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}

	// Wait for connection to close
	<-conn.DisconnectNotify()
}