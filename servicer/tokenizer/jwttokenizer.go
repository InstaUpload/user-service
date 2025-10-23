package tokenizer

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	t "github.com/instaUpload/user-service/types"
)

type JWTTokenizer struct {
	secretKey  string
	expireTime int64 // in hours
}

type JWTClaims struct {
	UserID int32 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTTokenizer() *JWTTokenizer {
	config := t.NewTokenizerConfig()
	return &JWTTokenizer{
		secretKey:  config.SecretKey,
		expireTime: int64(config.ExpirationHours)}
}

func (jt *JWTTokenizer) GenerateToken(userID int32) (string, error) {
	duration := time.Duration(jt.expireTime)
	claims := JWTClaims{
		userID,
		jwt.RegisteredClaims{
			Subject:   "User authentication",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "instaUpload",
		},
	}
	// Implement JWT token generation logic here
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jt.secretKey))
}

func (jt *JWTTokenizer) ValidateToken(token string) (int32, error) {
	values, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jt.secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := values.Claims.(*JWTClaims); ok {
		return claims.UserID, nil
	}
	// Implement JWT token validation logic here
	return 0, fmt.Errorf("invalid token")
}
