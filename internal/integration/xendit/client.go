package xendit

import (
	"context"
	"net/http"
	"time"
)

type Client interface {
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error)
}

type client struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client
}

func NewClient(secretKey string) Client {
	return &client{
		baseURL:   "https://api.xendit.co",
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
