package nanohttp

import (
	"errors"
	"slices"
	"strings"
)

var (
	ErrPathAlreadyRegistered = errors.New("route already registered for this method and path")
	ErrInvalidPath           = errors.New("route path cannot be empty")
	ErrNilHandler            = errors.New("handler cannot be nil")
)

type HandlerFunc func(w *ResponseWriter, r *Request)

type Router struct {
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
	}
}
func (r *Router) Handle(method, path string, handler HandlerFunc) error {
	m := strings.ToUpper(method)
	if !slices.Contains(ValidHttpMethods, m) {
		return ErrUnsupportedHttpMethod
	}
	if path == "" || !strings.HasPrefix(path, "/") {
		return ErrInvalidPath
	}
	if handler == nil {
		return ErrNilHandler
	}
	key := m + " " + path
	if _, exists := r.routes[key]; exists {
		return ErrPathAlreadyRegistered
	}
	r.routes[key] = handler
	return nil
}

func (r *Router) GET(path string, handler HandlerFunc) error {
	return r.Handle("GET", path, handler)
}
func (r *Router) POST(path string, handler HandlerFunc) error {
	return r.Handle("POST", path, handler)
}
func (r *Router) PUT(path string, handler HandlerFunc) error {
	return r.Handle("PUT", path, handler)
}
func (r *Router) PATCH(path string, handler HandlerFunc) error {
	return r.Handle("PATCH", path, handler)
}
func (r *Router) DELETE(path string, handler HandlerFunc) error {
	return r.Handle("DELETE", path, handler)
}

func (r *Router) ServeHTTP(w *ResponseWriter, req *Request) error {
	key := strings.ToUpper(req.Method) + " " + req.Path
	handler, exists := r.routes[key]
	if !exists {
		return w.String(404, "404 page not found\n")
	}
	handler(w, req)
	return nil
}
