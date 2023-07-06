package httpsrv

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	http.Server

	logger *zap.Logger
	// api    API
	Config *Config
}

func NewHTTPServer(logger *zap.Logger, config *Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Context()
		fmt.Fprint(w, "Hello, World!")
	})

	return &Server{
		Server: http.Server{
			Addr:              config.Port,
			Handler:           mux,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
		},
		logger: logger,
		// api:    API,
		Config: config,
	}
}

func (s *Server) Start(context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info(
		"Starting HTTP server",
		zap.String("addr", s.Addr),
		// zap.String("version", s.api.Version()),
	)

	go func() {
		if err := s.Serve(listener); err != nil {
			s.logger.Error("Failed to start HTTP server", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return func() error {
		err := s.Shutdown(ctx)
		if err != nil {
			s.logger.Error("Failed to shutdown HTTP server", zap.Error(err))

			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}

		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}()
}
