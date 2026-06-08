package model

import (
	"time"
)

type RentalStatus string

const (
	RentalPending   RentalStatus = "pending"
	RentalOngoing   RentalStatus = "ongoing"
	RentalComplete  RentalStatus = "completed"
	RentalOverdue   RentalStatus = "overdue"
	RentalCancelled RentalStatus = "cancelled"
)

type Rental struct {
	ID         int
	CustomerID int
	ProductID  int
	CreatedBy  int
	StartDate  time.Time
	EndDate    time.Time
	TotalPrice int64
	Status     RentalStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
