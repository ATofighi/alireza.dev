package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Create a UDP address on port 9999 (listens on all interfaces).
	addr := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("0.0.0.0"),
	}

	// Listen for incoming UDP packets at the specified address.
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Error listening on UDP port 9999: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP server listening on %s\n", addr.String())

	for {
		// Prepare a buffer to hold incoming data.
		buffer := make([]byte, 1024)

		// ReadFromUDP blocks until a UDP message arrives.
		// n is the number of bytes read, addr is the source address.
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error receiving UDP packet: %v\n", err)
			continue
		}

		// Log received message.
		message := string(buffer[:n])
		log.Printf("Received: %s from %s\n", message, clientAddr.String())

		// Send a response back to the client.
		response := fmt.Sprintf("Echo: %s", message)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			log.Printf("Error sending UDP response: %v\n", err)
		}
	}
}
