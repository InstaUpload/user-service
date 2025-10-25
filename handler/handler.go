package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/instaUpload/user-service/cache"
	"github.com/instaUpload/user-service/servicer"

	t "github.com/instaUpload/user-service/types"
	u "github.com/instaUpload/user-service/utils"
)

type Handler interface {
	// Define handler methods here
	HealthCheck(http.ResponseWriter, *http.Request)
	CreateUser(http.ResponseWriter, *http.Request)
	LoginUser(http.ResponseWriter, *http.Request)
	Close(ctx context.Context, cancel context.CancelFunc)
}

type ChiHandler struct {
	servicer servicer.Servicer
}

func NewChiHandler(ctx context.Context) *ChiHandler {
	// Create a new servicer using the provided context.
	ctx, cancel := context.WithCancel(ctx)
	serv := servicer.NewService(ctx)
	// Create a new producer using the provided context.
	// prod := producer.NewProducer(ctx)
	// Return a new ChiHandler instance with the servicer and producer.
	handle := &ChiHandler{servicer: serv}
	go handle.Close(ctx, cancel)
	return handle
}

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Implement handler methods for ChiHandler here
func (h *ChiHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Implementation here
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	slog.Info("Health check endpoint hit", slog.String("source", cwd))
	version := h.servicer.GetVersion()
	response := HealthCheckResponse{
		Status:  "OK",
		Version: version,
	}
	u.WriteResponse(w, http.StatusOK, response)
}

func (h *ChiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Step 1 convert request body to CreateUserRequest struct
	var newUser t.User
	err := u.ParseJSON(r.Body, &newUser)
	if err != nil {
		u.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Step 2 validate the struct fields.
	err = u.ValidateStruct(newUser)
	if err != nil {
		u.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = h.servicer.CreateUser(r.Context(), &newUser)
	if err != nil {
		slog.Error("Failed to create user", slog.String("error", err.Error()))
		u.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	u.WriteResponse(w, http.StatusCreated, &newUser)
}

func (h *ChiHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Step 1 convert request body to LoginUserRequest struct
	var req t.User
	err := u.ParseJSON(r.Body, &req)
	if err != nil {
		u.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Step 2 validate the struct fields.
	err = u.ValidateStruct(req)
	if err != nil {
		u.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.servicer.LoginUser(r.Context(), &req)
	if err != nil {
		slog.Error("Failed to login user", slog.String("error", err.Error()))
		u.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	response := map[string]string{
		"token": token,
	}
	u.WriteResponse(w, http.StatusOK, response)
}

func (h *ChiHandler) Close(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	<-ctx.Done()
	slog.Warn("Shutting down Chi handler...")
}
