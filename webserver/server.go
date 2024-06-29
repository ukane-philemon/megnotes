package webserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ukane-philemon/megtask/jwt"
)

// WebServer handles all routing and server logic.
type WebServer struct {
	mux    *chi.Mux
	log    *slog.Logger
	taskDB TaskDatabase

	jwtManager *jwt.Manager
}

// New returns a new instance of *WebServer.
func New(db TaskDatabase, logger *slog.Logger) (*WebServer, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	jwtManager, err := jwt.NewJWTManager()
	if err != nil {
		return nil, fmt.Errorf("jwt.NewJWTManager error: %w", err)
	}

	chiMux := chi.NewMux()
	chiMux.Use(middleware.Logger)
	chiMux.Use(middleware.AllowContentType("application/json"))

	server := &WebServer{
		mux:        chiMux,
		log:        logger,
		taskDB:     db,
		jwtManager: jwtManager,
	}

	server.registerRoutes()

	return server, nil
}

// Start starts the server and blocks until the server is stopped. All resources
// used by the server (e.g TaskDatabase) will be shutdown after server has been
// shutdown successfully.
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

	s.log.Info("Megtask server has started on -> ", "addr", server.Addr)

	// Wait for application shutdown.
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.log.Info("Gracefully shutting down the HTTP webserver....")

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		s.log.Error("server.Shutdown error: ", "msg", err)
	}

	dbShutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.log.Info("Gracefully shutting down TaskDatabase....")

	err = s.taskDB.Shutdown(dbShutdownCtx)
	if err != nil {
		s.log.Error("taskDB.Shutdown error: ", "msg", err)
	}

	return serverError
}

// registerRoutes registers all the required routes on s.mux.
func (s *WebServer) registerRoutes() {
	s.mux.Get("/", s.handleHome)

	s.mux.Post("/create-account", s.handleCreateAccount)
	s.mux.Post("/login", s.handleLogin)

	// Endpoints the require authentication.
	s.mux.Group(func(authedMux chi.Router) {
		authedMux.Use(s.authMiddleware)

		authedMux.Post("/task", s.handleCreateTask)
	})
}

// handleHome handles the "GET /" endpoint and returns a server message.
func (s *WebServer) handleHome(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(res, "Welcome to the Megtask API")
	if err != nil {
		s.log.Error("failed to write response: ", "error", err)
	}
}
