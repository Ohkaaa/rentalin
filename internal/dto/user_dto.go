package dto

type UserResponse struct {
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"alice"`
	Email    string `json:"email" example:"alice@example.com"`
	Phone    string `json:"phone" example:"628123456789"`
	Address  string `json:"address" example:"Jl. Merdeka No. 123"`
	Role     string `json:"role" example:"customer"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=100"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
}
