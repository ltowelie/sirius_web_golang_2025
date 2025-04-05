package initialization

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	HTTP *http.Server
}

func NewServer(addr string, router *gin.Engine) *Server {
	s := &Server{HTTP: &http.Server{}}
	s.HTTP.Addr = addr
	s.HTTP.Handler = router

	return s
}

func (s *Server) Close(ctx context.Context) error {
	err := s.HTTP.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
