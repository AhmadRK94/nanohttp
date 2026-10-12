package nanohttp

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"testing"
)

func TestRouterHandleValidation(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		handler HandlerFunc
		wantErr error
	}{
		{
			name:    "unsupported method",
			method:  "FOOBAR",
			path:    "/hello",
			handler: func(w *ResponseWriter, r *Request) {},
			wantErr: ErrUnsupportedHttpMethod,
		},
		{
			name:    "empty path",
			method:  "GET",
			path:    "",
			handler: func(w *ResponseWriter, r *Request) {},
			wantErr: ErrInvalidPath,
		},
		{
			name:    "path without leading slash",
			method:  "GET",
			path:    "hello",
			handler: func(w *ResponseWriter, r *Request) {},
			wantErr: ErrInvalidPath,
		},
		{
			name:    "nil handler",
			method:  "GET",
			path:    "/hello",
			handler: nil,
			wantErr: ErrNilHandler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()

			err := router.Handle(
				tt.method,
				tt.path,
				tt.handler,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Handle() error = %v, want %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestRouterDuplicatePath(t *testing.T) {
	router := NewRouter()

	if err := router.GET("/test", func(w *ResponseWriter, r *Request) {}); err != nil {
		t.Fatal(err)
	}

	err := router.GET("/test", func(w *ResponseWriter, r *Request) {})
	if err == nil {
		t.Fatal("error should not be nil")
	}

	if !errors.Is(err, ErrPathAlreadyRegistered) {
		t.Fatalf("error = %v, want %v", err, ErrPathAlreadyRegistered)
	}
}

func TestRouterServerHttp(t *testing.T) {
	router := NewRouter()
	called := false
	err := router.GET("/test", func(w *ResponseWriter, r *Request) {
		called = true
	})
	if err != nil {
		t.Fatal(err)
	}
	w := NewResponseWriter(nil)
	r := &Request{
		Method: "GET",
		Path:   "/test",
	}
	err = router.ServeHTTP(w, r)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("handler not invoked.")
	}
}

// integration test
func TestRouter(t *testing.T) {
	rawRequest := "GET /hello HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n"

	reader := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := ParseRequest(reader)
	if err != nil {
		t.Fatal(err)
	}

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	w := NewResponseWriter(serverConn)

	router := NewRouter()

	err = router.GET("/hello", func(w *ResponseWriter, r *Request) {
		_ = w.String(200, "hello from handler")
	})
	if err != nil {
		t.Fatal(err)
	}
	serveHttpErr := make(chan error, 1)
	go func() {
		serveHttpErr <- router.ServeHTTP(w, req)
	}()

	statusLine, _, body := readResponse(t, clientConn)
	if statusLine != "HTTP/1.1 200 OK\r\n" {
		t.Errorf("status line = %q, want %q",
			statusLine, "HTTP/1.1 200 OK\r\n")
	}

	if string(body) != "hello from handler" {
		t.Errorf("body = %q, want %q", body, "hello from handler")
	}

	if err := <-serveHttpErr; err != nil {
		t.Fatal(err)
	}
}
