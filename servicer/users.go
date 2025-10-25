package servicer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/instaUpload/user-service/database"
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
	output := LoginUserOutput{
		UserID:      user.ID,
		AccessToken: token,
	}
	return output, nil
}
