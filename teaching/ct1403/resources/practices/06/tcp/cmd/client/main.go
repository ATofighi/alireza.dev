package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Connect to the TCP server at localhost:9999.
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to TCP server.")

	// Send a message to the server.
	message := "Hello from the TCP client!\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	log.Printf("Sent: %s", message)

	// Read the server's response.
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	fmt.Printf("Server replied: %s", response)

	fmt.Println("Type messages to send to the server:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text() + "\n"
		if _, err = conn.Write([]byte(text)); err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		// Read the server's echo
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading response: %v", err)
			break
		}
		fmt.Printf("Server replied: %s", response)
	}
}
