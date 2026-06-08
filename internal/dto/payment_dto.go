package dto

type CreatePaymentRequest struct {
	RentalID int `json:"rental_id"`
}

type XenditWebhookPayload struct {
	ExternalID     string `json:"external_id"`
	Status         string `json:"status"`
	PaidAmount     int64  `json:"paid_amount"`
	PaymentMethod  string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
}
