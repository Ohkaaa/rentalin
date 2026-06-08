package model

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Phone     string
	Address   string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
