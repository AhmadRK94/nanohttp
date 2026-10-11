package nanohttp

import (
	"bufio"
	"errors"
	"io"
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

const maxBodySize = 10 * 1024 * 1024 //10MB

var (
	ValidHttpMethods  = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	ValidHttpVersions = []string{"HTTP/1.1"}
)

var (
	ErrUnsupportedHttpMethod       = errors.New("unsupported http method")
	ErrUnsupportedHttpVersion      = errors.New("unsupported http version")
	ErrInvalidHttpRequestLine      = errors.New("invalid request line")
	ErrInvalidHttpHeader           = errors.New("invalid header")
	ErrInvalidContentLength        = errors.New("invalid Content Length")
	ErrBodyTooLarge                = errors.New("body size is too large. Maximum size 10MB")
	ErrUnsupportedTransferEncoding = errors.New("unsupported transfer encoding")
)

func ParseRequest(reader *bufio.Reader) (*Request, error) {
	req := Request{
		Headers:     make(map[string]string),
		QueryParams: make(map[string][]string),
		PathParams:  make(map[string]string),
	}
	//parsing request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	requestLine = strings.TrimSuffix(requestLine, "\r\n")
	requestLineSplits := strings.Fields(requestLine)
	if len(requestLineSplits) != 3 {
		return nil, ErrInvalidHttpRequestLine
	}
	if !slices.Contains(ValidHttpMethods, requestLineSplits[0]) {
		return nil, ErrUnsupportedHttpMethod
	}
	if !slices.Contains(ValidHttpVersions, requestLineSplits[2]) {
		return nil, ErrUnsupportedHttpVersion
	}
	req.Version = requestLineSplits[2]
	req.Method = requestLineSplits[0]
	parsedURL, err := url.ParseRequestURI(requestLineSplits[1])
	if err != nil {
		return nil, err
	}
	req.Path = parsedURL.Path
	maps.Copy(req.QueryParams, parsedURL.Query())

	//parsing headers
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if header == "\r\n" {
			break
		}
		header = strings.TrimSuffix(header, "\r\n")
		headerSplits := strings.SplitN(header, ":", 2)
		if len(headerSplits) != 2 {
			return nil, ErrInvalidHttpHeader
		}
		key := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(headerSplits[0], ":")))
		if key == "" {
			return nil, ErrInvalidHttpHeader
		}
		value := strings.TrimSpace(headerSplits[1])
		req.Headers[key] = value
	}
	//parsing body
	_, exists := req.Headers["transfer-encoding"]
	if exists {
		return nil, ErrUnsupportedTransferEncoding
	}
	clStr, exists := req.Headers["content-length"]
	if !exists {
		return &req, nil
	}
	cl, err := strconv.Atoi(clStr)
	if err != nil || cl < 0 {
		return nil, ErrInvalidContentLength
	}
	if cl > maxBodySize {
		return nil, ErrBodyTooLarge
	}
	if cl > 0 {
		buf := make([]byte, cl)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}
		req.Body = buf

	}

	return &req, nil
}
