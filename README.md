# nanohttp

A minimal HTTP server and routing library built from scratch in Go using TCP connections and the standard `net` package.

`nanohttp` is a learning project designed to explore how HTTP servers work internally, including request parsing, response construction, routing, and concurrent connection handling.

## Features

- **TCP server:** Accepts incoming connections using Go's `net` package.
- **HTTP request parsing:** Parses the request line, headers, query parameters, and request body.
- **Routing:** Registers handlers using HTTP methods and exact paths.
- **Response writer:** Provides `String()` and `JSON()` helpers for sending HTTP responses.
- **HTTP status codes:** Supports common HTTP response status codes.
- **Concurrent connections:** Handles each connection in its own goroutine.
- **Basic error handling:** Handles malformed requests and unknown routes.

## Installation

Clone the repository:

```bash
git clone https://github.com/AhmadRK94/nanohttp.git
cd nanohttp
```

## Usage

Create a router, register a handler, and start the server:

```go
package main

import (
    "log"

    "github.com/AhmadRK94/nanohttp"
)

func main() {
    router := nanohttp.NewRouter()

    router.GET("/hello", func(w *nanohttp.ResponseWriter, r *nanohttp.Request) {
        w.String(200, "Hello from nanohttp!")
    })

    server := nanohttp.NewServer(router, ":8080")

    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
```

Start the application and visit:

```text
http://localhost:8080/hello
```

You should receive:

```text
Hello from nanohttp!
```

You can also inspect the HTTP response using curl:

```bash
curl -i http://localhost:8080/hello
```

## Supported Features and Limitations

This project intentionally implements only a small subset of HTTP server functionality.

**Supported**

- HTTP request parsing with `Content-Length` body framing
- Exact method-and-path route matching
- String and JSON responses
- One request per connection, with the connection closed afterward

**Not currently supported**

- Persistent connections and HTTP keep-alive
- Chunked transfer encoding
- Dynamic path parameters
- Middleware
- Advanced routing and HTTP protocol features

## Running Tests

Run the test suite:

```bash
go test ./...
```

Format and check the code:

```bash
gofmt -w .
go vet ./...
```

## Project Goals

The goal of `nanohttp` is to understand the fundamentals of HTTP servers and backend development by implementing a minimal version rather than relying on a full-featured web framework.

This project is intended for learning and experimentation, not production use. For production applications, use Go's standard `net/http` package or an established HTTP framework.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
