# Http-reverse-proxy server

This project successfully implements an HTTP reverse proxy, without the use of third party implementations or the net/http/httputil package.
The terminology used in the project can be better referenced from the below image.

![img.png](img.png)
source: https://www.cloudflare.com/learning/cdn/glossary/reverse-proxy/
## Project Structure

* origin_server.go    # Implementation of the Origin Server
* reverse_proxy.go    # Implementation of the Reverse Proxy
* main.go             # Starting point of the application which runs both the servers concurrently


### How could someone get started with your codebase?
Clone the repo locally and run the below commands:
- `go install` or `go build `
- `go run .`
  from a terminal you can run
- `curl -i http://localhost:8082/test` or access http://localhost:8082/test to see the response from the origin server


In order to explore this further one can deploy an application at http://localhost:8081/test and run the above curl
to test out the reverse-proxy

### What resources did you use to build your implementation?
Below are the resources used:
- https://go.dev/src/net/http/client.go
- https://pkg.go.dev/net/http#hdr-Clients_and_Transports
- https://www.cloudflare.com/learning/cdn/glossary/reverse-proxy/


### Explain any design decisions you made, including limitations of the system



### How would you scale this?

The scaling of a reverse-proxy server would entail increasing its capacity to handle more incoming traffic efficiently, adding redundancy and removing single-point of failures.

- Vertical Scaling:
    - Spawning Goroutines: Currently this application has goroutines per request. We better optimize the resource allocation
      during these goroutines by using fixed-size [worker pool](https://gobyexample.com/worker-pools) goroutines to process the incoming proxy requests concurrently.
      Incoming tasks will be in a buffered channel and each worker can read from this channel to process each tasks.
    - FineTuning resource limits like `writeTimeout`, `IdleTimeout`, `ReadTimeout` to prevent resource exhaustion

- Horizontal Scaling:
    - Multiple Instances & Load Balancing: Running multiple instances of the proxy server across different machines or containers. Using a dedicated load balancer like Nginx to distribute incoming traffic evenly.
      This isolates failures and allows one to scale out as needed.
    - Deploying in an Orchestration: Creating container images of the reverse proxy-server and deploying in an distributed environment like kubernetes can help scale this.
      This would also ease auto-scalling, rolling-updates and load-balancing.
    - Service Discovery: Implementing dynamic service discovery using tools like Consul or Kubernetes service discovery can help evenly distribute traffic to healthy proxy servers.

- Queuing Request Task :
    - Implementing a task queue that manages incoming requests thereby preventing the server from being overwhelmed by bursts of traffic.
- Caching and Rate Limiting:
    - Adding a cache to store the repeated request would reduce the load and handle response rates better.

### How would you make it more secure?

In order to make the server more secure we do the following:

- **Sanitization of the incoming requests:** Validating the incoming requests to ensure it does not contain malicious payload like a large string that could be used to exploit vulnerabilities (like buffer overflows)

- **HTTP Header Filtering:** Removing or sanitizing header that could be exploited.

- **Rate-limiting and Throttling the incoming requests:** A sudden large number of incoming request can cause of Denial of Service (DoS attack). In order to prevent this kind of attack we need to implement
  a request task queue with timeouts to drop and delay requests under heavy load.

- **Adding Encryption:** Using TLS to encrypt all the traffic from the client to the server and from backend to the proxy server


### What resources (including programming tool assistants) did you use to build your implementation?

The above linked resources along with the auto-code completed from GoLand IDE. 

