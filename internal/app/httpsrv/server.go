package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	http.Server

	logger *zap.Logger
	// api    API
	config      *Config
	userService UserService
}

func NewHTTPServer(
	logger *zap.Logger,
	config *Config,
	userService UserService,
) *Server {
	s := &Server{
		logger:      logger,
		config:      config,
		userService: userService,
	}

	router := chi.NewRouter()

	router.Get(config.MetricPath, promhttp.Handler().ServeHTTP)
	router.Route("/users", func(r chi.Router) {
		// r.Get("/", s.list)
		r.Get("/{userID}", s.GetUserByID)
		// r.Get("/{id}", s.get)
		// r.Put("/{id}", s.update)
		// r.Delete("/{id}", s.delete)
	})

	s.Server = http.Server{
		Addr:              fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler:           router,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
	}

	return s
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

func (s *Server) GetUserByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	ctx := r.Context()

	user, err := s.userService.GetByID(ctx, uuid)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	payload, err := json.Marshal(user)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
