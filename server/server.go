package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/instaUpload/user-service/handler"
)

type Server interface {
	// Define server methods here
	Handler() http.Handler
	MountRoutes()
	Close(ctx context.Context, cancel context.CancelFunc)
}

type ChiServer struct {
	// Define Chi server specific fields here
	router *chi.Mux
	handle handler.Handler
}

func NewChiServer(ctx context.Context) *ChiServer {
	ctx, cancel := context.WithCancel(ctx)
	router := chi.NewRouter()
	// Create a Handler.
	handle := handler.NewChiHandler(ctx)
	server := &ChiServer{router: router, handle: handle}
	// Initialize Chi server specific fields here
	go server.Close(ctx, cancel)
	return server
}

func (s *ChiServer) MountRoutes() {
	// Mount routes to the Chi server
	s.router.Get("/health", s.handle.HealthCheck)
	s.router.Post("/users", s.handle.CreateUser)
	s.router.Post("/users/login", s.handle.LoginUser)
	s.router.Group(func(r chi.Router) {
		r.Use(s.handle.AuthenticateUser)
		// Protected routes go here
		r.Get("/users", s.handle.GetUserByID)
	})
	// Add more routes as needed
}

func (s *ChiServer) Handler() http.Handler {
	// Return the HTTP handler for the Chi server
	return s.router
}

func (s *ChiServer) Close(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	<-ctx.Done()
	slog.Warn("Shutting down Chi server...")
	// s.router.Shutdown(ctx)
}
