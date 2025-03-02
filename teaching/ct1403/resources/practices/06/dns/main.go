package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	domain := "alireza.dev"

	// LookupIP returns a slice of that domain's IPv4 and IPv6 addresses.
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Fatalf("Failed to lookup IP addresses: %v\n", err)
	}

	// Print all IP addresses associated with the domain.
	for _, ip := range ips {
		fmt.Printf("%s -> %s\n", domain, ip.String())
	}
}
