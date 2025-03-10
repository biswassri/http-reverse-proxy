package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func runOriginServer(errChan chan<- error) {

	originServerHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("[origin server] received request at:%s originServerHandler %s %s\n ", time.Now(), req.URL.Path, req.Method)
		//Record response from Server
		_, err := fmt.Fprint(w, "origin server response")
		if err != nil {
			log.Fatal(err)
		}
	})
	//Custom server for monitoring incoming requests
	s := &http.Server{
		Addr:         ":8081",
		Handler:      originServerHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	//Origin Server running on port:8081
	log.Printf("Starting origin server on localhost port %s ", s.Addr)
	err := s.ListenAndServe()
	log.Printf("Server on %s exited: %v", s.Addr, err)
	errChan <- err

}
