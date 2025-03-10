package main

import "log"

func main() {
	errChOrigin := make(chan error, 1)
	errChRProxy := make(chan error, 1)

	go runOriginServer(errChOrigin)
	go runReverseProxy(errChRProxy)

	select {
	case err := <-errChOrigin:
		log.Fatal("Origin server exited ", err)
	case err := <-errChRProxy:
		log.Fatal("Reverse proxy server exited", err)
	}
}
