package server

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

type Server struct {
	logger *httplog.Logger
	
	HTTP *http.Server

	running bool
}

func NewServer(cfg *config.ServerConfig, logger *httplog.Logger) (*Server, error) {
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
	r.Use(middleware.Logger)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer(r, "/", http.Dir("./static"))

	//TODO API routes

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
