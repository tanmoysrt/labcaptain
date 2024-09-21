package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var (
	hostPattern      = regexp.MustCompile(`^port-(\d+).\S+$`)
	blacklistedPorts = map[string]bool{
		"80": true, "8001": true, "8002": true, "8003": true, "8004": true,
	}
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 32*1024)
		},
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	matches := hostPattern.FindStringSubmatch(r.Host)
	if matches == nil {
		http.Error(w, "Host header format invalid.", http.StatusExpectationFailed)
		return
	}

	port := matches[1]
	if blacklistedPorts[port] {
		http.Error(w, fmt.Sprintf("Port %s is blacklisted", port), http.StatusBadRequest)
		return
	}

	destConn, err := net.DialTimeout("tcp", "localhost:"+port, 10*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to localhost:%s", port), http.StatusBadGateway)
		return
	}
	defer destConn.Close()

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

	if err := r.Write(destConn); err != nil {
		log.Printf("Failed to write request to destination: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	copyConn := func(dst io.Writer, src io.Reader) {
		defer wg.Done()
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf)
		io.CopyBuffer(dst, src, buf)
	}

	go copyConn(destConn, clientConn)
	go copyConn(clientConn, destConn)

	wg.Wait()
}

func main() {
	server := &http.Server{
		Addr:         ":8004",
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Starting proxy server on :8004")
	log.Fatal(server.ListenAndServe())
}
