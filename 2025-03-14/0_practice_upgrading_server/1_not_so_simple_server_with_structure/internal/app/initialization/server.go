package initialization

import (
	"net/http"
)

type Server struct {
	HTTP *http.Server
}

func NewServer(addr string, router *http.ServeMux) *Server {
	s := &Server{HTTP: &http.Server{}}
	s.HTTP.Addr = addr
	s.HTTP.Handler = router

	return s
}
