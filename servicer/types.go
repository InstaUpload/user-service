package servicer

type CreateUserInput struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	UserID int32 `json:"user_id"`
}
