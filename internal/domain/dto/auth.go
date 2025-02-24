package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5,max=15,alphanum"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=15"`
}
