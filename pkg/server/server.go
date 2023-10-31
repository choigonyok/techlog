package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func New(h http.Handler) *Server {
	s := newHTTPServer(h)
	return &Server{
		httpServer: s,
	}
}

func newHTTPServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:    "localhost:8000",
		Handler: h,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
