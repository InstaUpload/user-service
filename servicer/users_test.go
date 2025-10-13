package servicer

import (
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	// serv := CreateTestService(ctx)
	userInput := CreateUserInput{
		Email:    "gpt.sahaj28@gmail.com",
		Fullname: "GPT Sahaj",
		Password: "password123",
	}
	t.Run("Create new user", func(t *testing.T) {
		userOutput, err := serv.CreateUser(ctx, userInput)
		if userOutput.UserID == 0 || err != nil {
			t.Errorf("expected user id to be non-zero and no error, got userID: %d, err: %v", userOutput.UserID, err)
		}
	})

	t.Run("Create user with existing email", func(t *testing.T) {
		_, err := serv.CreateUser(ctx, userInput)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
