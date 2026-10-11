package nanohttp

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	rawRequest := "POST /users?id=42 HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: 5\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"hello"

	reader := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := ParseRequest(reader)
	if err != nil {
		t.Fatalf("ParseRequest() error = %v", err)
	}
	if req.Method != "POST" {
		t.Errorf("Method = %q, want %q", req.Method, "POST")
	}
	if req.Path != "/users" {
		t.Errorf("Path = %q, want %q", req.Path, "/users")
	}
	if req.Version != "HTTP/1.1" {
		t.Errorf("Version = %q, want %q", req.Version, "HTTP/1.1")
	}
	if string(req.Body) != "hello" {
		t.Errorf("Body = %q, want %q", req.Body, "hello")
	}
	if req.Headers["content-length"] != "5" {
		t.Errorf("Content-Length = %q, want %q", req.Headers["content-length"], "5")
	}
	if req.Headers["host"] != "localhost" {
		t.Errorf("Host= %q, want %q", req.Headers["host"], "localhost")
	}
	if req.Headers["content-type"] != "text/plain" {
		t.Errorf("Content-Length = %q, want %q", req.Headers["content-type"], "text/plain")
	}
	if got := req.QueryParams["id"]; len(got) != 1 || got[0] != "42" {
		t.Errorf("QueryParams[id] = %q, want 42", got[0])
	}

}
