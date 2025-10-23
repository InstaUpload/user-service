package tokenizer

import (
	b64 "encoding/base64"
	"fmt"

	t "github.com/instaUpload/user-service/types"
)

type BasicTokenizer struct {
	secretKey string
}

func NewBasicTokenizer() *BasicTokenizer {
	config := t.NewTokenizerConfig()
	return &BasicTokenizer{secretKey: config.SecretKey}
}

func (bt *BasicTokenizer) GenerateToken(userID int32) (string, error) {
	// Implement token generation logic here
	tokenData := fmt.Sprintf("%d:%s", userID, bt.secretKey)
	token := b64.StdEncoding.EncodeToString([]byte(tokenData))
	return token, nil
}

func (bt *BasicTokenizer) ValidateToken(token string) (int32, error) {
	// Implement token validation logic here
	decodedData, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	var userID int32
	var secret string
	_, err = fmt.Sscanf(string(decodedData), "%d:%s", &userID, &secret)
	if err != nil {
		return 0, err
	}
	if secret != bt.secretKey {
		return 0, fmt.Errorf("invalid token")
	}
	return userID, nil
}
