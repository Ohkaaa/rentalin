package model

import (
	"encoding/json"
	"time"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
	PaymentExpired PaymentStatus = "expired"
)

type Payment struct {
	ID              int
	CustomerID      int
	RentalID        int
	ExternalID      string
	InvoiceURL      string
	Amount          int64
	PaidAmount      *int64
	Currency        string
	Method          *string
	PaymentChannel  *string
	Status          PaymentStatus
	ExpiredAt       time.Time
	PaidAt          *time.Time
	Description     *string
	CallbackPayload json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
