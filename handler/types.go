package handler

type CreateUserRequest struct {
	Fullname string `json:"fullname" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreateUserResponse struct {
	UserID  int32 `json:"user_id"`
	Message string `json:"message"`
}