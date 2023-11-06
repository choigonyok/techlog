package server

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// New creates new Server object
func New(h http.Handler, address string) *Server {
	s := newHTTPServer(h, address)
	return &Server{
		httpServer: s,
	}
}

// newHTTPServer creates new http.Server with specified address and handler
func newHTTPServer(h http.Handler, address string) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: h,
	}
}

// Start starts http.Server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
