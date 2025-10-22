package tokenizer

type Tokenizer interface {
	GenerateToken(userID int32) (string, error)
	ValidateToken(token string) (int32, error)
}
