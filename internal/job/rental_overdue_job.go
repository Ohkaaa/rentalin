package job

import (
	"context"
	"log"
	"time"

	"rentalin/internal/repository"
)

func StartOverdueChecker(rentalRepo repository.RentalRepository) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			err := rentalRepo.UpdateOverdueRentals(context.Background())
			if err != nil {
				log.Println("overdue job:", err)
			}
		}
	}()
}
