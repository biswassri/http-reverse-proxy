package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func runReverseProxy() {

	originServerURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		return
	}

	reverseProxyServerHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("[reverse proxy] received request at:%s reverProxyServerHandler %s %s\n ", time.Now(), req.URL.Path, req.Method)
		//modifying the original request
		req.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.URL.Host = originServerURL.Host
		req.RequestURI = ""

		/*_, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}*/

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = fmt.Fprint(w, err)
			return
		}

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
	log.Printf("Starting reverse proxy server on %s ", s.Addr)
	err = s.ListenAndServe()
	//
	log.Printf("Server on %s exited: %v", s.Addr, err)

}

func main() {

	go runOriginServer()
	go runReverseProxy()
}
