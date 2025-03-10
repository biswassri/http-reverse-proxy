package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const baseHost = "http://localhost"

func runReverseProxy(errChan chan<- error) {

	originServerURL, err := url.Parse(fmt.Sprintf("%s:8081", baseHost))
	if err != nil {
		errChan <- fmt.Errorf("invalid Origin Server URL: %v", err)
		return
	}

	reverseProxyServerHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("[reverse proxy] received request at:%s reverProxyServerHandler %s %s\n ", time.Now(), req.URL.Path, req.Method)
		//Creating a proxy request with url path from origin server

		proxyReq, err := http.NewRequest(req.Method, originServerURL.String()+req.URL.Path, req.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		proxyReq.Header = req.Header.Clone()

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = fmt.Fprint(w, err)
			return
		}
		//close response
		defer resp.Body.Close()

		// Copy headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}

	})
	s := &http.Server{
		Addr:         ":8082",
		Handler:      reverseProxyServerHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	//Reverse Proxy server is running on: 8082
	log.Printf("Starting reverse proxy server on localhost port %s ", s.Addr)
	err = s.ListenAndServe()
	//
	log.Printf("Server on %s exited: %v", s.Addr, err)
	errChan <- err
}
