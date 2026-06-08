package errs

import "errors"

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

// User errors
var (
	ErrInvalidAddress        = errors.New("invalid address")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidPhoneFormat    = errors.New("invalid phone number format")
	ErrInvalidUsernameFormat = errors.New("invalid username format")
	ErrUserEmailExists       = errors.New("email already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserPhoneExists       = errors.New("phone number already exists")
	ErrUserUsernameExists    = errors.New("username already exists")
	ErrUsernameContainsSpace = errors.New("username cannot contain spaces")
	ErrWeakPassword          = errors.New("password is too weak")
)

// Product errors
var (
	ErrInvalidPrice    = errors.New("price must be greater than zero")
	ErrInvalidStock    = errors.New("stock cannot be negative")
	ErrProductNotFound = errors.New("product not found")
)

// Rental errors
var (
	ErrProductOutOfStock   = errors.New("product is out of stock")
	ErrInvalidDate         = errors.New("invalid date format")
	ErrInvalidRentalPeriod = errors.New("rental period is invalid")
	ErrRentalNotFound      = errors.New("rental not found")
)

// Payment errors
var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInternal        = errors.New("please try again later")
)
