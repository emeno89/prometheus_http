prometheus_http
======

[![Build&test](https://github.com/emeno89/prometheus_http/actions/workflows/go.yml/badge.svg)](https://github.com/emeno89/prometheus_http/actions/workflows/go.yml)

prometheus_http implements simple middleware for HTTP queries, which collects for prometheus request time, request and response counts.

# Installation
```
go get github.com/emeno89/prometheus_http
```

# Metrics

- *http_duration_seconds* - value (in seconds) of duration your server spent for execution;
- *http_requests_total* - how many requests your server received;
- *http_responses_total* - how many responses your server returned.

# Quick Start

```go
package main

import (
	"fmt"
	"github.com/emeno89/prometheus_http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	promMiddleware := middleware.NewHandler("simple_project", func(r *http.Request) string {
		return r.RequestURI
	})

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":9090", mux)
	}()

	mux := http.NewServeMux()
	mux.Handle("/products", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Products list response for %s", r.RequestURI)
	}))

	if err := http.ListenAndServe(":80", promMiddleware.Middleware(mux)); err != nil {
		log.Println(err)
	}
}
````

In this example we create simple HTTP server (named "simple project"), implement path collecting function and use middleware with server mux.