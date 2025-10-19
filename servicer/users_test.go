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

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	// serv := CreateTestService(ctx)
	userInput := LoginUserInput{
		Email:    "gpt.sahaj28@gmail.com",
		Password: "password123",
	}
	t.Run("Login with correct credentials", func(t *testing.T) {
		userOutput, err := serv.LoginUser(ctx, userInput)
		if userOutput.AccessToken == "" || err != nil {
			t.Errorf("expected valid token and no error, got token: %s, err: %v", userOutput.AccessToken, err)
		}
	})
	t.Run("Login with incorrect password", func(t *testing.T) {
		invalidInput := LoginUserInput{
			Email:    "gpt.sahaj28@gmail.com",
			Password: "wrongpassword",
		}
		_, err := serv.LoginUser(ctx, invalidInput)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
	t.Run("Login with non-existing email", func(t *testing.T) {
		invalidInput := LoginUserInput{
			Email:    "nonexisting@emmail.com",
			Password: "password123",
		}
		_, err := serv.LoginUser(ctx, invalidInput)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
