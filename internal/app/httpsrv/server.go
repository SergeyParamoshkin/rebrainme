package httpsrv

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	http.Server

	logger *zap.Logger
	api    API
}

func New(params Params) *Server {
	return &Server{
		Server: http.Server{
			Addr:              params.Config.Addr,
			Handler:           params.API,
			ReadHeaderTimeout: params.Config.ReadHeaderTimeout,
		},
		logger: params.Logger,
		api:    params.API,
	}
}

func (s *Server) Start(context.Context) error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.logger.Info(
		"Starting HTTP server",
		zap.String("addr", s.Addr),
		zap.String("version", s.api.Version()),
	)

	go func() {
		if err := s.Serve(ln); err != nil {
			s.logger.Error("Failed to start HTTP server", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
