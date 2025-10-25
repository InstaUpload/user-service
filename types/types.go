package types

type User struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname" validate:"max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}
