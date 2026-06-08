package dto

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required,min=3,max=100"`
	DailyPrice int64  `json:"daily_price" validate:"required,gt=0"`
	Stock      int    `json:"stock" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name       *string `json:"name" validate:"omitempty,min=3,max=100"`
	DailyPrice *int64  `json:"daily_price" validate:"omitempty,gt=0"`
	Stock      *int    `json:"stock" validate:"omitempty,gte=0"`
}
