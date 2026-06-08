package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100,alphanum" example:"alice"`
	Email    string `json:"email" validate:"required,email" example:"alice@example.com"`
	Phone    string `json:"phone" validate:"required" example:"628123456789"`
	Address  string `json:"address" validate:"required,min=5" example:"Jl. Merdeka No. 123"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"alice@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxx.yyy"`
}

type RegisterResponse struct {
	Message string       `json:"message" example:"register success"`
	Data    UserResponse `json:"data"`
}

type LoginResponse struct {
	Message string       `json:"message" example:"login success"`
	Data    AuthResponse `json:"data"`
}
