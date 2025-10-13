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
