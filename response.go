package nanohttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

var HttpStatus = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	301: "Moved Permanently",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	413: "Content Too Large",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
}

var (
	ErrInvalidStatusCode = errors.New("Invalid status code.")
)

type ResponseWriter struct {
	conn net.Conn
}

func NewResponseWriter(conn net.Conn) *ResponseWriter {
	return &ResponseWriter{conn: conn}
}

func (w *ResponseWriter) writeResponse(statusCode int, contentType string, body []byte) error {
	if _, exist := HttpStatus[statusCode]; !exist {
		return ErrInvalidStatusCode
	}
	if _, err := fmt.Fprintf(w.conn, "HTTP/1.1 %d %s\r\n", statusCode, HttpStatus[statusCode]); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w.conn, "Content-Type: %s\r\n", contentType); err != nil {
		return err
	}
	if statusCode != 204 {
		if _, err := fmt.Fprintf(w.conn, "Content-Length: %d\r\n", len(body)); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w.conn, "Connection: close\r\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(w.conn, "\r\n"); err != nil {
		return err
	}
	if statusCode == 204 {
		return nil
	}
	if _, err := w.conn.Write(body); err != nil {
		return err
	}
	return nil

}

func (w *ResponseWriter) String(statusCode int, body string) error {

	return w.writeResponse(
		statusCode,
		"text/plain; charset=utf-8",
		[]byte(body))
}

func (w *ResponseWriter) JSON(statusCode int, data any) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.writeResponse(
		statusCode,
		"application/json", d)
}
