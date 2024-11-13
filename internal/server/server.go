package server

import (
	"fmt"
	"net/http"

	"github.com/choigonyok/techlog/internal/router"
)

type Server struct {
	*http.Server
}

const (
	listenAddress = "0.0.0.0:8080" // should be 0.0.0.0, not localhost
)

// New creates new Server object
func New() (*Server, error) {
	r, err := router.New()
	if err != nil {
		fmt.Println("ERR SETTING MIDDLEWARES:", err)
	}
	httpHandlers := r.GetHTTPHandler()

	s := newHTTPServer(httpHandlers, listenAddress)

	svr := &Server{
		Server: s,
	}
	return svr, err
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
	return s.ListenAndServe()
}
