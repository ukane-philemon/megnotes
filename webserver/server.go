package webserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// WebServer handles all routing and server logic.
type WebServer struct {
	mux *chi.Mux
	log *slog.Logger
}

// New returns a new instance of *WebServer.
func New(log *slog.Logger) *WebServer {
	chiMux := chi.NewMux()
	chiMux.Use(middleware.Logger)

	server := &WebServer{
		mux: chiMux,
		log: log,
	}

	server.registerRoutes()
	return server
}

// Start starts the server and blocks until the server is stopped.
func (s *WebServer) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	var serverError error
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverError = err
		}
	}()

	s.log.Info("Megnotes server has started on -> ", "addr", server.Addr)

	// Wait for application shutdown.
	<-ctx.Done()

	err := server.Shutdown(ctx)
	if err != nil {
		s.log.Error("server.Shutdown error: ", "msg", err)
	}

	return serverError
}

// registerRoutes registers all the required routes on s.mux.
func (s *WebServer) registerRoutes() {
	s.mux.Get("/", s.handleHome)
}

// handleHome handles the "GET /" endpoint and returns a server message.
func (s *WebServer) handleHome(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(res, "Welcome to the Megnotes API")
	if err != nil {
		s.log.Error("failed to write response: ", "error", err)
	}
}
