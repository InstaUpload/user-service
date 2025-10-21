package tokenizer

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	t "github.com/instaUpload/user-service/types"
)

type TokenizerJWT interface {
	GenerateJWT(email string) (string, error)
}

type Claims struct {
	Email string
	jwt.RegisteredClaims
}

type BasicTokenizerJWT struct {
	secretKeyJWT []byte
}

func NewBasicTokenizerJWT() *BasicTokenizerJWT {
	config := t.NewTokenizerConfigJWT()
	return &BasicTokenizerJWT{secretKeyJWT: config.SecretKeyJWT}
}

// GenerateJWT generates a token for a email
func (bt *BasicTokenizerJWT) GenerateJWT(email string) (string, error) {
	//token will expire in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(bt.secretKeyJWT))
}
