package nanohttp

import (
	"bufio"
	"errors"
	"log"
	"net"
)

type Server struct {
	router *Router
	addr   string
}

func NewServer(router *Router, addr string) *Server {
	return &Server{
		router: router,
		addr:   addr,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer listener.Close()
	log.Printf("server listening on %s", s.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go s.handleConnection(conn)
	}

}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in connection handler: %v\n", r)
		}
	}()

	reader := bufio.NewReader(conn)
	req, err := ParseRequest(reader)
	w := NewResponseWriter(conn)

	if err != nil {
		log.Printf("failed to parse request from %s: %v",
			conn.RemoteAddr(), err)

		if writeErr := w.String(400, "bad request\n"); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	if err := s.router.ServeHTTP(w, req); err != nil {
		log.Printf("failed to serve request from %s: %v",
			conn.RemoteAddr(), err)
	}
}
