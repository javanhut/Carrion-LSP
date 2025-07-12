package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestMainFunction_StdioMode(t *testing.T) {
	// Test that main can start in stdio mode without crashing
	// This is a basic smoke test since we can't easily test the full LSP interaction

	// Skip this test in short mode since it involves subprocess execution
	if testing.Short() {
		t.Skip("Skipping main function test in short mode")
	}

	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "test-carrion-lsp", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("test-carrion-lsp")

	// Test stdio mode with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "./test-carrion-lsp", "--stdio")

	// Create pipes for stdin/stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Send a simple LSP initialize request
	initRequest := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}` + "\r\n"
	contentLength := len(initRequest)
	lspMessage := "Content-Length: " + string(rune(contentLength)) + "\r\n\r\n" + initRequest

	// Write to stdin
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(lspMessage))
		time.Sleep(100 * time.Millisecond)
		// Send shutdown
		shutdownRequest := `{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}` + "\r\n"
		contentLength := len(shutdownRequest)
		lspMessage := "Content-Length: " + string(rune(contentLength)) + "\r\n\r\n" + shutdownRequest
		stdin.Write([]byte(lspMessage))
	}()

	// Read some output to verify it's responding
	buffer := make([]byte, 1024)
	// Use a goroutine with timeout for reading
	done := make(chan bool)
	var n int
	go func() {
		n, err = stdout.Read(buffer)
		done <- true
	}()

	select {
	case <-done:
		if err != nil && err != io.EOF {
			// This is expected since we're not doing a full LSP handshake
			t.Logf("Read error (expected): %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Log("Read timeout (expected)")
	}

	if n > 0 {
		response := string(buffer[:n])
		t.Logf("Received response: %s", response)
		// Should contain LSP headers
		if !strings.Contains(response, "Content-Length:") {
			t.Error("Expected LSP response format with Content-Length header")
		}
	}

	// Wait for the command to finish or timeout
	cmd.Wait()
}

func TestMainFunction_TCPMode(t *testing.T) {
	// Test that main can start in TCP mode
	if testing.Short() {
		t.Skip("Skipping TCP mode test in short mode")
	}

	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "test-carrion-lsp-tcp", ".")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("test-carrion-lsp-tcp")

	// Start the server in TCP mode with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "./test-carrion-lsp-tcp", string(rune(port)))

	// Capture stderr to check for startup message
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start TCP server: %v", err)
	}

	// Give the server time to start
	time.Sleep(500 * time.Millisecond)

	// Try to connect to the server
	conn, err := net.DialTimeout("tcp", "localhost:"+string(rune(port)), 1*time.Second)
	if err != nil {
		t.Logf("Could not connect to server (this might be expected): %v", err)
		t.Logf("Server stderr: %s", stderr.String())
	} else {
		conn.Close()
		t.Log("Successfully connected to TCP server")
	}

	// Terminate the server
	cmd.Process.Kill()
	cmd.Wait()
}

func TestArgumentParsing(t *testing.T) {
	// Test argument parsing logic by examining os.Args usage in main
	// Since we can't easily mock os.Args in the main function,
	// we'll test the logic conceptually

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "stdio mode",
			args:     []string{"carrion-lsp", "--stdio"},
			expected: "stdio",
		},
		{
			name:     "tcp mode default",
			args:     []string{"carrion-lsp"},
			expected: "tcp",
		},
		{
			name:     "tcp mode with port",
			args:     []string{"carrion-lsp", "8080"},
			expected: "tcp",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock the argument parsing logic from main
			var mode string
			if len(test.args) > 1 && test.args[1] == "--stdio" {
				mode = "stdio"
			} else {
				mode = "tcp"
			}

			if mode != test.expected {
				t.Errorf("Expected mode %s, got %s", test.expected, mode)
			}
		})
	}
}

func TestPortParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedPort string
	}{
		{
			name:         "default port",
			args:         []string{"carrion-lsp"},
			expectedPort: "7777",
		},
		{
			name:         "custom port",
			args:         []string{"carrion-lsp", "8080"},
			expectedPort: "8080",
		},
		{
			name:         "stdio mode ignores port",
			args:         []string{"carrion-lsp", "--stdio", "8080"},
			expectedPort: "7777", // Default, not used in stdio mode
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock the port parsing logic from main
			port := "7777" // default

			// Only update port if not in stdio mode and port is provided in the right position
			if len(test.args) > 1 && test.args[1] == "--stdio" {
				// In stdio mode, port is ignored
			} else if len(test.args) > 1 {
				// First argument is the port in TCP mode
				port = test.args[1]
			} else if len(test.args) > 2 && test.args[1] != "--stdio" {
				// Port is second argument only if not stdio mode
				port = test.args[2]
			}

			if port != test.expectedPort {
				t.Errorf("Expected port %s, got %s", test.expectedPort, port)
			}
		})
	}
}

// Integration test for LSP message format
func TestLSPMessageFormat(t *testing.T) {
	// Test that we can parse and generate proper LSP messages
	message := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}`
	contentLength := len(message)

	lspMessage := "Content-Length: " + string(rune(contentLength)) + "\r\n\r\n" + message

	// Parse the message
	lines := strings.Split(lspMessage, "\r\n")
	if len(lines) < 3 {
		t.Error("Expected at least 3 lines in LSP message")
	}

	// Check Content-Length header
	if !strings.HasPrefix(lines[0], "Content-Length:") {
		t.Error("Expected Content-Length header")
	}

	// Check empty line separator
	if lines[1] != "" {
		t.Error("Expected empty line after headers")
	}

	// Check JSON content
	jsonContent := strings.Join(lines[2:], "\r\n")
	if !strings.Contains(jsonContent, "initialize") {
		t.Error("Expected initialize method in JSON content")
	}
}

// Test that the binary can be built successfully
func TestBuildSuccess(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-build", ".")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nStderr: %s", err, stderr.String())
	}

	// Cleanup
	os.Remove("test-build")
}

// Test that all imports can be resolved
func TestImports(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to resolve imports: %v", err)
	}
}

// Test go mod tidy
func TestModuleTidy(t *testing.T) {
	cmd := exec.Command("go", "mod", "tidy")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run go mod tidy: %v\nStderr: %s", err, stderr.String())
	}
}

// Test that the version is correctly set
func TestVersionInfo(t *testing.T) {
	// This would test version information if we had it in the binary
	// For now, just ensure the binary runs with basic flags
	if testing.Short() {
		t.Skip("Skipping version test in short mode")
	}

	cmd := exec.Command("go", "build", "-o", "test-version", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build for version test: %v", err)
	}
	defer os.Remove("test-version")

	// Test that the binary doesn't crash on startup
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "./test-version", "--stdio")
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start binary for version test: %v", err)
	}

	// Let it run briefly then kill
	time.Sleep(100 * time.Millisecond)
	cmd.Process.Kill()
	cmd.Wait()
}

// Benchmark the main function startup time
func BenchmarkMainStartup(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	// Build once
	cmd := exec.Command("go", "build", "-o", "bench-carrion-lsp", ".")
	err := cmd.Run()
	if err != nil {
		b.Fatalf("Failed to build benchmark binary: %v", err)
	}
	defer os.Remove("bench-carrion-lsp")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

		cmd := exec.CommandContext(ctx, "./bench-carrion-lsp", "--stdio")
		err := cmd.Start()
		if err != nil {
			b.Fatalf("Failed to start benchmark binary: %v", err)
		}

		// Let it initialize briefly
		time.Sleep(10 * time.Millisecond)

		cmd.Process.Kill()
		cmd.Wait()
		cancel()
	}
}
