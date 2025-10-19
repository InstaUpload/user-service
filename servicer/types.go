package servicer

type CreateUserInput struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	UserID int32 `json:"user_id"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutput struct {
	UserID      int32  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
