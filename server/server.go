package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/instaUpload/user-service/handler"
	"github.com/instaUpload/user-service/servicer/tokenizer"
)

type Server interface {
	// Define server methods here
	Handler() http.Handler
	MountRoutes()
	Close(ctx context.Context, cancel context.CancelFunc)
}

type ChiServer struct {
	// Define Chi server specific fields here
	router    *chi.Mux
	handle    handler.Handler
	tokenizer *tokenizer.BasicTokenizerJWT
}

func NewChiServer(ctx context.Context) *ChiServer {
	ctx, cancel := context.WithCancel(ctx)
	router := chi.NewRouter()
	tokenizer := tokenizer.NewBasicTokenizerJWT()
	// Create a Handler.
	handle := handler.NewChiHandler(ctx)
	server := &ChiServer{
		router:    router,
		handle:    handle,
		tokenizer: tokenizer,
	}
	// Initialize Chi server specific fields here
	go server.Close(ctx, cancel)
	return server
}

func (s *ChiServer) MountRoutes() {
	// Mount routes to the Chi server
	s.router.Get("/health", s.handle.HealthCheck)
	s.router.Post("/users", s.handle.CreateUser)
	s.router.Post("/users/login", s.handle.LoginUser)

	// These are the protected endpoints that only
	// authenticated users can access
	s.router.Group(func(protected chi.Router) {
		protected.Use(tokenizer.NewBasicTokenizerJWT().Authenticate)
		protected.Get("/profile", tokenizer.Profile)
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
