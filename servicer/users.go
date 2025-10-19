package servicer

import (
	"context"
	"fmt"

	"github.com/instaUpload/user-service/database"
	u "github.com/instaUpload/user-service/utils"
	"github.com/jackc/pgx/v5"
)

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	// Step 1: Check if user with same email already exists
	user, err := s.db.GetUserByEmail(ctx, input.Email)
	if err != nil && err != pgx.ErrNoRows {
		return CreateUserOutput{}, err
	}
	if user.Email != "" {
		return CreateUserOutput{}, fmt.Errorf("user with email %s already exists", input.Email)
	}
	// Step 2: Hash the password
	hashedPassword, err := u.HashPassword(input.Password)
	if err != nil {
		return CreateUserOutput{}, err
	}
	// Step 3: Insert user into database
	userToBeCreated := database.CreateUserParams{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: hashedPassword,
	}
	createdUser, err := s.db.CreateUser(ctx, userToBeCreated)
	if err != nil {
		return CreateUserOutput{}, err
	}
	// Step 4: Return user ID and success message
	output := CreateUserOutput{
		UserID: createdUser.ID,
	}
	return output, nil
}

func (s *Service) LoginUser(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	// Step 1: Fetch user by email
	user, err := s.db.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return LoginUserOutput{}, fmt.Errorf("invalid email or password")
		}
		return LoginUserOutput{}, err
	}
	// Step 2: Compare passwords
	err = u.CheckPasswordHash(input.Password, user.Password)
	if err != nil {
		return LoginUserOutput{}, fmt.Errorf("invalid email or password")
	}
	// Step 3: Generate JWT token
	token, err := s.tokenizer.GenerateToken(user.ID)
	if err != nil {
		return LoginUserOutput{}, err
	}
	// Step 4: Return user ID and access token
	output := LoginUserOutput{
		UserID:      user.ID,
		AccessToken: token,
	}
	return output, nil
}
