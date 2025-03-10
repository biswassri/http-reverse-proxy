package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

const defaultTimeout = 5 * time.Second

// check if server is reachable
func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil // Server up
		}
		time.Sleep(100 * time.Millisecond)
	}
	return http.ErrServerClosed // Timeout is reached
}

// testing origin server
func TestOriginServerResponse(t *testing.T) {
	errCh := make(chan error, 1)
	go runOriginServer(errCh)
	if err := waitForServer("http://localhost:8081", defaultTimeout); err != nil {
		t.Fatalf("Origin server did not start: %v", err)
	}

	resp, err := http.Get("http://localhost:8081/test")
	if err != nil {
		t.Fatalf("Failed to reach origin server: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	expected := "origin server response"
	if string(body) != expected {
		t.Errorf("Unexpected origin server response: got %q, want %q", string(body), expected)
	}
}

// testing proxy server
func TestReverseProxyResponse(t *testing.T) {
	errCh := make(chan error, 1)
	go runOriginServer(errCh) // Can be replaced with a mock Origin server for testing with Reverse Proxy
	go runReverseProxy(errCh)
	if err := waitForServer("http://localhost:8082", defaultTimeout); err != nil {
		t.Fatalf("Reverse proxy did not start: %v", err)
	}
	resp, err := http.Get("http://localhost:8082/test")
	if err != nil {
		t.Fatalf("Failed to reach reverse proxy: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	expected := "origin server response"
	if string(body) != expected {
		t.Errorf("Unexpected reverse proxy response: got %q, want %q", string(body), expected)
	}
}

// testing both
func TestIntegration(t *testing.T) {
	errCh := make(chan error, 2)

	go runOriginServer(errCh) // localhost:8081
	go runReverseProxy(errCh) // localhost:8082

	//waiting for origin server
	if err := waitForServer("http://localhost:8081", defaultTimeout); err != nil {
		t.Fatalf("Origin server did not start: %v", err)
	}

	//waiting for reverse proxy server
	if err := waitForServer("http://localhost:8082", defaultTimeout); err != nil {
		t.Fatalf("Reverse proxy did not start: %v", err)
	}

}
