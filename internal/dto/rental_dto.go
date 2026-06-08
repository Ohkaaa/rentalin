package dto

type CreateRentalRequest struct {
	CustomerID *int   `json:"customer_id,omitempty"` // Hanya untuk admin
	ProductID  int    `json:"product_id" validate:"required,gt=0"`
	StartDate  string `json:"start_date" validate:"required"`
	EndDate    string `json:"end_date" validate:"required"`
}
