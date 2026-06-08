package xendit

type CreateInvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description,omitempty"`
}

type CreateInvoiceResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
	ExpiredAt  string `json:"expiry_date"`
}
