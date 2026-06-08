package job

import (
	"context"
	"log"
	"time"

	"rentalin/internal/repository"
)

func StartPaymentExpiryJob(paymentRepo repository.PaymentRepository) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			err := paymentRepo.ExpirePendingPayments(context.Background())
			if err != nil {
				log.Println("payment expiry job:", err)
			}
		}
	}()
}
