# healthcheck
[![build](https://github.com/tamcore/go-healthcheck/actions/workflows/build.yml/badge.svg)](https://github.com/tamcore/go-healthcheck/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamcore/go-healthcheck)](https://goreportcard.com/report/github.com/tamcore/go-healthcheck)
[![GoDoc](https://godoc.org/github.com/tamcore/go-healthcheck?status.svg)](https://godoc.org/github.com/tamcore/go-healthcheck)


Healthcheck is a library for implementing Kubernetes [liveness and readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) probe handlers in your Go application.

## Features

 - Integrates easily with Kubernetes. This library explicitly separates liveness vs. readiness checks instead of lumping everything into a single category of check.

 - Supports asynchronous checks, which run in a background goroutine at a fixed interval. These are useful for expensive checks that you don't want to add latency to the liveness and readiness endpoints.

 - Includes a follow useful checks:
    * DNS
    * TCP
    * HTTP
    * database
    * Go runtime.
    * GC max pause
    * MongoDb
    * Redis

 ## Usage

See the [GoDoc examples](https://godoc.org/github.com/tamcore/go-healthcheck) for more detail.

Install dependency
 
  ```bash
  go get -u github.com/tamcore/go-healthcheck
  ```

Import the package 

```go
import "github.com/tamcore/go-healthcheck"
```

Create a `healthcheck.Handler`:

```go
health := healthcheck.NewHandler()
```

Configure some application-specific liveness checks (whether the app itself is unhealthy):

```go
// Our app is not happy if we've got more than 100 goroutines running.
health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))
```

Configure some application-specific readiness checks (whether the app is ready to serve requests):

```go
// Our app is not ready if we can't resolve our upstream dependency in DNS.
health.AddReadinessCheck("upstream-dep-dns", dns.Resolve("upstream.example.com", 50*time.Millisecond))

// Our app is not ready if we can't connect to our database (`var DB *sql.DB`) in <1s.
health.AddReadinessCheck("database", db.Ping(DB, 1*time.Second))
```

Expose the `/live` and `/ready` endpoints over HTTP (on port 8086):

```go
go func() {
	httpError := http.ListenAndServe("0.0.0.0:8086", health)

	if httpError != nil {
		log.Println("While serving HTTP: ", httpError)
	}
}()
```

Configure your Kubernetes container with HTTP liveness and readiness probes see the ([Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)) for more detail:

```yaml
# this is a bare bones example
# copy and paste livenessProbe and readinessProbe as appropriate for your app
apiVersion: v1
kind: Pod
metadata:
  name: healthcheck-example
spec:
  containers:
  - name: liveness
    image: your-registry/your-container

    # define a liveness probe that checks every 5 seconds, starting after 5 seconds
    livenessProbe:
      httpGet:
        path: /live
        port: 8086
      initialDelaySeconds: 5
      periodSeconds: 5

    # define a readiness probe that checks every 5 seconds
    readinessProbe:
      httpGet:
        path: /ready
        port: 8086
      periodSeconds: 5
```

If one of your readiness checks fails, Kubernetes will stop routing traffic to that pod within a few seconds (depending on `periodSeconds` and other factors).

If one of your liveness checks fails or your app becomes totally unresponsive, Kubernetes will restart your container.

 ## HTTP Endpoints
 When you run `go http.ListenAndServe("0.0.0.0:8086", health)`, two HTTP endpoints are exposed:

  - **`/live`**: liveness endpoint (HTTP 200 if healthy, HTTP 503 if unhealthy)
  - **`/ready`**: readiness endpoint (HTTP 200 if healthy, HTTP 503 if unhealthy)

Pass the `?full=1` query parameter to see the full check results as JSON. These are omitted by default for performance.
