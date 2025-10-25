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
	AuthenticateUser(http.Handler) http.Handler
	GetUserByID(http.ResponseWriter, *http.Request)
	Close(ctx context.Context, cancel context.CancelFunc)
}

type ChiHandler struct {
	servicer servicer.Servicer
	cache    cache.Cacher
}

func NewChiHandler(ctx context.Context) *ChiHandler {
	// Create a new servicer using the provided context.
	ctx, cancel := context.WithCancel(ctx)
	serv := servicer.NewService(ctx)
	cache := cache.NewRedisCache(ctx)
	// Create a new producer using the provided context.
	// prod := producer.NewProducer(ctx)
	// Return a new ChiHandler instance with the servicer and producer.
	handle := &ChiHandler{servicer: serv, cache: cache}
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

func (h *ChiHandler) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implementation here
		token := r.Header.Get("Authorization")
		if token == "" {
			u.WriteErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}
		userID, err := h.servicer.AuthenticateToken(r.Context(), token)
		if err != nil {
			slog.Error("Failed to authenticate token", slog.String("error", err.Error()))
			u.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		var user t.User
		// Check if user exists in cache
		user, err = h.cache.RetrieveUser(userID)
		if err != nil {
		}
		if user.ID == 0 {
			// User not in cache, fetch from database
			user, err = h.servicer.GetUserByID(r.Context(), userID)
			if err != nil {
				slog.Error("Failed to fetch user from database", slog.String("error", err.Error()))
				u.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
				return
			}
			// Store user in cache
			err = h.cache.StoreUser(user)
			if err != nil {
				slog.Error("Failed to store user in cache", slog.String("error", err.Error()))
				u.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
				return
			}
			slog.Info("User fetched from database and stored in cache", slog.Int64("userID", user.ID))
		}
		slog.Info("User retrived from cache", slog.Int64("userID", user.ID))
		// Store user in request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *ChiHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Implementation here
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		u.WriteErrorResponse(w, http.StatusBadRequest, "Missing user ID")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		u.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, err := h.cache.RetrieveUser(userID)
	if err != nil {
		slog.Error("Failed to fetch user", slog.String("error", err.Error()))
		u.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}
	u.WriteResponse(w, http.StatusOK, user)
}

func (h *ChiHandler) Close(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	<-ctx.Done()
	slog.Warn("Shutting down Chi handler...")
}
