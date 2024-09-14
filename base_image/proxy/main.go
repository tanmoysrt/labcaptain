package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

// Regular expression to match "port-<dest_port>:<random_string>"
var hostPattern = regexp.MustCompile(`^port-(\d+):[a-zA-Z0-9]+$`)
// Blacklisted ports
var blacklistedPorts = map[string]bool{
	"80": true, // proxy port - high chance go in recursive loop
	"8001": true, // web terminal, already mapped so skip
	"8002": true, // code server, already mapped so skip
	"8003": true, // vnc, already mapped so skip
	"8004": true, // internal proxy, none should be allowed to directly access this
}

// Handle the incoming requests and proxy them based on the port from the Host header
func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the Host header
	host := r.Host

	// Check if the Host header matches the required format
	matches := hostPattern.FindStringSubmatch(host)
	if matches == nil {
		http.Error(w, "Host header format invalid. Expected format: port-<dest_port>:<random_string>", http.StatusExpectationFailed)
		return
	}

	// Extract the destination port from the Host header
	port := matches[1]

	// Check if the destination port is blacklisted
	if blacklistedPorts[port] {
		http.Error(w, fmt.Sprintf("Port %s is blacklisted", port), http.StatusBadRequest)
		return
	}

	// Set up the target address to forward to
	target := fmt.Sprintf("localhost:%s", port)

	// Establish a TCP connection to the target
	destConn, err := net.DialTimeout("tcp", target, 10*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to %s", target), http.StatusBadGateway)
		return
	}
	defer destConn.Close()

	// Hijack the original connection to get raw access to the underlying network connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Forward the original request to the destination
	err = r.Write(destConn)
	if err != nil {
		http.Error(w, "Failed to write request to destination", http.StatusInternalServerError)
		return
	}

	// Create a bidirectional copy: copy data from client to destination and vice versa
	go io.Copy(destConn, clientConn) // Copy data from client to destination
	io.Copy(clientConn, destConn)    // Copy data from destination to client
}

func main() {
	// Start the HTTP server and handle requests
	http.HandleFunc("/", handler)
	log.Println("Starting proxy server on :8004")
	log.Fatal(http.ListenAndServe(":8004", nil))
}
