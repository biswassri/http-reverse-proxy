package main

import "log"

func main() {
	errCh := make(chan error, 2)

	go runOriginServer(errCh)
	go runReverseProxy(errCh)
	err := <-errCh
	log.Fatal("Server exited with error:", err)
}
