package servicer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/instaUpload/user-service/database"
	t "github.com/instaUpload/user-service/types"
	u "github.com/instaUpload/user-service/utils"
	"github.com/jackc/pgx/v5"
)

func (s *Service) CreateUser(ctx context.Context, newUser *t.User) error {
	// Step 1: Check if user with same email already exists
	user, err := s.db.GetUserByEmail(ctx, newUser.Email)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}
	if user.Email != "" {
		return fmt.Errorf("user with email %s already exists", newUser.Email)
	}
	// Step 2: Hash the password
	hashedPassword, err := u.HashPassword(newUser.Password)
	if err != nil {
		return err
	}
	// Step 3: Insert user into database
	userToBeCreated := database.CreateUserParams{
		Fullname: newUser.Fullname,
		Email:    newUser.Email,
		Password: hashedPassword,
	}
	createdUser, err := s.db.CreateUser(ctx, userToBeCreated)
	if err != nil {
		return err
	}
	// Step 4: Return user ID and success message
	newUser.ID = int64(createdUser.ID)
	newUser.Password = ""
	return nil
}

func (s *Service) LoginUser(ctx context.Context, loginUser *t.User) (string, error) {
	// Step 1: Fetch user by email
	user, err := s.db.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("invalid email or password")
		}
		return "", err
	}
	// Step 2: Compare passwords
	err = u.CheckPasswordHash(loginUser.Password, user.Password)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}
	// Step 3: Generate JWT token
	token, err := s.tokenizer.GenerateToken(user.ID)
	if err != nil {
		slog.Error("Error generating token", "error", err)
		return "", err
	}
	// Step 4: Return user ID and access token
	return token, nil
}

func (s *Service) AuthenticateToken(ctx context.Context, token string) (int64, error) {
	// Step 1: Validate and parse the token
	userID, err := s.tokenizer.ValidateToken(token)
	if err != nil {
		slog.Error("Error validating token", "error", err)
		return 0, fmt.Errorf("invalid or expired token")
	}
	// Step 2: Return user ID
	return int64(userID), nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int64) (t.User, error) {
	// Step 1 get user by Id from store
	user, err := s.db.GetUserByID(ctx, int32(userID))
	if err != nil {
		return t.User{}, err
	}

	return t.User{
		ID:       int64(user.ID),
		Email:    user.Email,
		Fullname: user.Fullname,
	}, nil
}
