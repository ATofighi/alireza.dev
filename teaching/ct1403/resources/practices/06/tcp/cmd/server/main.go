package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// Listen on TCP port 9999 on all interfaces.
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalf("Failed to listen on port 9999: %v", err)
	}
	defer listener.Close()

	log.Println("TCP server listening on :9999")

	for {
		// Accept blocks until a connection is received.
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle each connection in a new goroutine.
		go handleConnection(conn)
	}
}

// handleConnection reads data from the connection and echoes it back.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Log the client address
	log.Printf("Client connected: %s", conn.RemoteAddr().String())

	// Wrap the connection in a buffered reader for convenience.
	reader := bufio.NewReader(conn)

	for {
		// Read data until newline or EOF.
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Client closed the connection
				log.Printf("Client disconnected: %s", conn.RemoteAddr().String())
			} else {
				log.Printf("Error reading from client %s: %v", conn.RemoteAddr().String(), err)
			}
			return
		}

		// Echo the received message back to the client.
		log.Printf("Received from %s: %s", conn.RemoteAddr().String(), line)
		_, writeErr := conn.Write([]byte(fmt.Sprintf("Echo: %s", line)))
		if writeErr != nil {
			log.Printf("Error writing to client %s: %v", conn.RemoteAddr().String(), writeErr)
			return
		}
	}
}
