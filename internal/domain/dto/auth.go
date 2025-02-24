package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5,max=15,alphanum"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=15"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}
