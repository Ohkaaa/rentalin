package model

import "time"

type Product struct {
	ID         int
	Name       string
	DailyPrice int64
	Stock      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
