package nanohttp

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
)

func readResponse(t *testing.T, conn net.Conn) (string, map[string]string, []byte) {
	reader := bufio.NewReader(conn)
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("error reading status line: %v", err)
	}
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("reading headers: %v", err)
		}
		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(strings.TrimRight(line, "\r\n"), ":", 2)
		if len(parts) != 2 {
			t.Fatalf("invalid response header: %q", line)
		}
		headers[strings.ToLower(parts[0])] = strings.TrimSpace(parts[1])
	}
	var body []byte

	if _, exists := headers["content-length"]; exists {
		bodySize, err := strconv.Atoi(headers["content-length"])
		if err != nil {
			t.Fatal(err)
		}
		body = make([]byte, bodySize)
		if _, err := io.ReadFull(reader, body); err != nil {
			t.Fatalf("reading body: %v", err)
		}
	}

	return statusLine, headers, body
}

func TestResponseWriterString(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	w := NewResponseWriter(serverConn)
	writeErr := make(chan error, 1)
	go func() {
		writeErr <- w.String(200, "hello")
	}()

	statusLine, headers, body := readResponse(t, clientConn)

	if statusLine != "HTTP/1.1 200 OK\r\n" {
		t.Errorf("status line = %q, want %q",
			statusLine, "HTTP/1.1 200 OK\r\n")
	}

	if headers["content-type"] != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q", headers["content-type"])
	}

	if headers["content-length"] != "5" {
		t.Errorf("Content-Length = %q, want %q",
			headers["content-length"], "5")
	}
	if string(body) != "hello" {
		t.Errorf("body = %q, want %q", body, "hello")
	}
	if err := <-writeErr; err != nil {
		t.Fatalf("String() error = %v", err)
	}
}

type name struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func TestResponseWriterJSON(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	w := NewResponseWriter(serverConn)
	writeErr := make(chan error, 1)
	jsonData := name{
		Name: "ahmad",
		Age:  "32",
	}
	go func() {
		writeErr <- w.JSON(200, jsonData)
	}()

	statusLine, headers, body := readResponse(t, clientConn)
	if statusLine != "HTTP/1.1 200 OK\r\n" {
		t.Errorf("status line = %q, want %q",
			statusLine, "HTTP/1.1 200 OK\r\n")
	}

	if headers["content-type"] != "application/json" {
		t.Errorf("Content-Type = %q", headers["content-type"])
	}
	expectedBody, err := json.Marshal(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	if headers["content-length"] != strconv.Itoa(len(expectedBody)) {
		t.Errorf("Content-Length = %q, want %q",
			headers["content-length"], strconv.Itoa(len(expectedBody)))
	}

	var data name
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("error to unmarshaling data: %v", err)
	}
	if data.Name != "ahmad" || data.Age != "32" {
		t.Errorf("got name = %q, age =%q, want %q and %q", data.Name, data.Age, "ahmad", "32")
	}
	if err := <-writeErr; err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
}

func TestResponseWriterInvalidStatusCode(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	w := NewResponseWriter(serverConn)

	err := w.String(999, "hello")
	if err == nil {
		t.Fatal("String() should write to client due to invalid status code.")
	}
	if !errors.Is(err, ErrInvalidStatusCode) {
		t.Fatalf("got error %v, want %v", err, ErrInvalidStatusCode)
	}
}

func TestResponseWriter204(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	w := NewResponseWriter(serverConn)
	writeErr := make(chan error, 1)
	go func() {
		writeErr <- w.String(204, "hello")
	}()

	statusLine, headers, body := readResponse(t, clientConn)
	if statusLine != "HTTP/1.1 204 No Content\r\n" {
		t.Errorf("status line = %q, want %q",
			statusLine, "HTTP/1.1 204 No Content\r\n")
	}

	if headers["content-type"] != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q", headers["content-type"])
	}

	if _, exists := headers["content-length"]; exists {
		t.Fatal("Content-Length shouldn't be specified.")
	}
	if len(body) != 0 {
		t.Fatal("Body should be empty.")
	}
	if err := <-writeErr; err != nil {
		t.Fatalf("String() error = %v", err)
	}
}
