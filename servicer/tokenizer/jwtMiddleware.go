package tokenizer

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	u "github.com/instaUpload/user-service/utils"
)

func (bt *BasicTokenizerJWT) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		// Check if the header is missing
		if tokenString == "" {
			u.WriteErrorResponse(w, http.StatusUnauthorized, "Missing header for the token")
			return
		}

		// Extract token from "Bearer <token>" format
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return bt.secretKeyJWT, nil
		})

		// Check if token is valid and without parsing errors
		if err != nil || !token.Valid {
			u.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Token is valid
		next.ServeHTTP(w, r)
	})
}

// This is function tests if the authentication is working or not
// This is a protected function. It means that it can be accessed only
// if the user is authenticated
func Profile(w http.ResponseWriter, r *http.Request) {
	u.WriteResponse(w, http.StatusOK, "Hello, you are an authenticated user!")
}
