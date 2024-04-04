package server

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/ScreamingHawk/go-adventure/narrator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

type Server struct {
	logger *httplog.Logger
	
	HTTP *http.Server
	Narrator *narrator.Narrator

	running bool
}

func NewServer(cfg *config.ServerConfig, logger *httplog.Logger, narrator *narrator.Narrator) (*Server, error) {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout: 45 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout: 60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &Server{
		logger: logger,
		HTTP: httpServer,
		Narrator: narrator,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	if s.running {
		return fmt.Errorf("server is already running")
	}
	s.running = true

	s.logger.Info(fmt.Sprintf("Server is running on %s", s.HTTP.Addr), "op", "run")

	s.HTTP.Handler = s.handler()

	if err := s.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if !s.running {
		return fmt.Errorf("server is not running")
	}
	s.running = false

	s.logger.Info("Server is stopping", "op", "stop")

	return s.HTTP.Shutdown(ctx)
}

func (s *Server) handler() http.Handler {
	// Add file extension to support Windows
	if err := mime.AddExtensionType(".js", "text/javascript"); err != nil {
		panic(err)
	}
	// Add file extension to support Windows
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Timeout
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.ThrottleBacklog(1, 10, 60 * time.Second)) // Cheap rate limiter

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		s.addNarratorRoutes(r)
	})

	s.addFileServer(r, "/", http.Dir("./static"))

	return r
}
