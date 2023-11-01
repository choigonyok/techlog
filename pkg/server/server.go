package server

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(h http.Handler, address string) *Server {
	s := newHTTPServer(h, address)
	return &Server{
		httpServer: s,
	}
}

func newHTTPServer(h http.Handler, address string) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: h,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
