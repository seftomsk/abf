package web

import (
	"context"
	"net"
	"net/http"

	"github.com/seftomsk/abf/internal/access"
	"github.com/seftomsk/abf/internal/limiter"
)

type ResponseDTO struct {
	Status string
	Code   int
	Msg    string
}

type RequestDTO struct {
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Server struct {
	host    string
	port    string
	server  *http.Server
	limiter *limiter.MultiLimiter
	access  *access.IPAccess
}

type ServerConfig interface {
	Host() string
	Port() string
}

func NewServer(
	l *limiter.MultiLimiter,
	a *access.IPAccess,
	cfg ServerConfig) *Server {
	return &Server{
		host:    cfg.Host(),
		port:    cfg.Port(),
		limiter: l,
		access:  a,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error

	go func() {
		server := &http.Server{
			Addr:    net.JoinHostPort(s.host, s.port),
			Handler: nil,
		}
		s.server = server

		http.HandleFunc("/auth", getAuthHandler(s.access, s.limiter))

		err = server.ListenAndServe()
	}()

	<-ctx.Done()

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
