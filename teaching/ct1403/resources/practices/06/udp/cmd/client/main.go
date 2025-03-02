package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Remote server address (replace IP/port as needed).
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9999,
	}

	// DialUDP creates a UDP “connection” for sending data to the server.
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		log.Fatalf("Error dialing UDP server: %v", err)
	}
	defer conn.Close()

	for {
		// Send a message to the server.
		message := make([]byte, 1000)
		_, err := fmt.Scanln(&message)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		_, err = conn.Write(message)
		if err != nil {
			log.Fatalf("Error sending UDP message: %v", err)
		}
		fmt.Printf("Sent: %s\n", string(message))

		// Prepare a buffer to read server response.
		responseBuf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(responseBuf)
		if err != nil {
			log.Fatalf("Error reading UDP response: %v", err)
		}

		// Print server response.
		fmt.Printf("Server reply: %s\n", string(responseBuf[:n]))
	}
}
