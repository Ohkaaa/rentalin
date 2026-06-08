package service

import (
	"context"
	"errors"
	"time"

	"rentalin/internal/errs"
	"rentalin/internal/model"
	"rentalin/internal/repository"
	"rentalin/pkg/logger"
	"rentalin/pkg/transaction"

	"github.com/sirupsen/logrus"
)

type RentalService interface {
	CreateRental(ctx context.Context, customerID, productID, createdBy int, startDateStr, endDateStr string) (*model.Rental, error)
	GetRentalByID(ctx context.Context, rentalID, userID int, role string) (*model.Rental, error)
	GetRentalsByCustomerID(ctx context.Context, customerID int) ([]*model.Rental, error)
	GetAllRentals(ctx context.Context) ([]*model.Rental, error)
	CompleteRental(ctx context.Context, rentalID int) error
	CancelRental(ctx context.Context, rentalID, userID int, role string) error
	UpdateRentalStatus(ctx context.Context, rentalID int, status model.RentalStatus) error
}

type rentalService struct {
	txRunner    *transaction.TxRunner
	rentalRepo  repository.RentalRepository
	productRepo repository.ProductRepository
}

func NewRentalService(txRunner *transaction.TxRunner, rentalRepo repository.RentalRepository, productRepo repository.ProductRepository) RentalService {
	return &rentalService{
		txRunner:    txRunner,
		rentalRepo:  rentalRepo,
		productRepo: productRepo,
	}
}

func isValidTransition(oldStatus, newStatus model.RentalStatus) bool {
	transitions := map[model.RentalStatus][]model.RentalStatus{
		model.RentalPending: {
			model.RentalOngoing,
			model.RentalCancelled,
		},
		model.RentalOngoing: {
			model.RentalComplete,
			model.RentalCancelled,
		},
	}

	allowed, ok := transitions[oldStatus]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}

	return false
}

func applyStockChange(product *model.Product, oldStatus, newStatus model.RentalStatus) error {
	switch {
	case oldStatus == model.RentalPending && newStatus == model.RentalOngoing:
		if product.Stock <= 0 {
			return errs.ErrProductOutOfStock
		}
		product.Stock--

	case oldStatus == model.RentalOngoing &&
		(newStatus == model.RentalComplete || newStatus == model.RentalCancelled):
		product.Stock++
	}

	return nil
}

func (s *rentalService) CreateRental(ctx context.Context, customerID, productID, createdBy int, startDateStr, endDateStr string) (*model.Rental, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		logger.Error(
			"failed get product",
			err,
			logrus.Fields{
				"product_id": productID,
			},
		)

		return nil, errs.WrapErr("GetByID", err)
	}

	if product.Stock <= 0 {
		logger.Warn(
			"product out of stock",
			logrus.Fields{
				"product_id": productID,
			},
		)

		return nil, errs.ErrProductOutOfStock
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		logger.Warn(
			"invalid rental start date",
			logrus.Fields{
				"start_date": startDateStr,
			},
		)

		return nil, errs.ErrInvalidDate
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		logger.Warn(
			"invalid rental end date",
			logrus.Fields{
				"end_date": endDateStr,
			},
		)

		return nil, errs.ErrInvalidDate
	}

	if !endDate.After(startDate) {
		logger.Warn(
			"invalid rental period",
			logrus.Fields{
				"start_date": startDateStr,
				"end_date":   endDateStr,
			},
		)

		return nil, errs.ErrInvalidRentalPeriod
	}

	days := int(endDate.Sub(startDate).Hours() / 24)
	totalPrice := int64(days) * product.DailyPrice

	now := time.Now()

	rental := &model.Rental{
		CustomerID: customerID,
		ProductID:  productID,
		CreatedBy:  createdBy,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
		Status:     model.RentalPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.rentalRepo.Create(ctx, rental); err != nil {
		logger.Error(
			"failed create rental",
			err,
			logrus.Fields{
				"customer_id": customerID,
				"product_id":  productID,
			},
		)

		return nil, errs.WrapErr("Create", err)
	}

	logger.Info(
		"rental created successfully",
		logrus.Fields{
			"rental_id":   rental.ID,
			"customer_id": customerID,
			"product_id":  productID,
			"total_price": totalPrice,
		},
	)

	return rental, nil
}

func (s *rentalService) GetRentalByID(ctx context.Context, rentalID, userID int, role string) (*model.Rental, error) {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"rental not found",
				logrus.Fields{
					"rental_id": rentalID,
				},
			)

			return nil, errs.ErrRentalNotFound
		}

		logger.Error(
			"failed get rental",
			err,
			logrus.Fields{
				"rental_id": rentalID,
			},
		)

		return nil, errs.WrapErr("GetByID", err)
	}

	if rental.CustomerID != userID && role != "admin" {
		logger.Warn(
			"unauthorized rental access",
			logrus.Fields{
				"rental_id": rentalID,
				"user_id":   userID,
			},
		)

		return nil, errs.ErrRentalNotFound
	}

	return rental, nil
}

func (s *rentalService) GetRentalsByCustomerID(ctx context.Context, customerID int) ([]*model.Rental, error) {
	rentals, err := s.rentalRepo.GetByUserID(ctx, customerID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"customer rentals not found",
				logrus.Fields{
					"customer_id": customerID,
				},
			)

			return nil, errs.ErrRentalNotFound
		}

		logger.Error(
			"failed get rentals by customer",
			err,
			logrus.Fields{
				"customer_id": customerID,
			},
		)

		return nil, errs.WrapErr("GetByUserID", err)
	}

	return rentals, nil
}

func (s *rentalService) GetAllRentals(ctx context.Context) ([]*model.Rental, error) {
	rentals, err := s.rentalRepo.GetAll(ctx)
	if err != nil {
		logger.Error(
			"failed get all rentals",
			err,
			nil,
		)

		return nil, errs.WrapErr("GetAll", err)
	}

	logger.Info(
		"all rentals fetched successfully",
		logrus.Fields{
			"total": len(rentals),
		},
	)

	return rentals, nil
}

func (s *rentalService) CompleteRental(ctx context.Context, rentalID int) error {
	logger.Info(
		"rental completion requested",
		logrus.Fields{
			"rental_id": rentalID,
		},
	)

	return s.UpdateRentalStatus(ctx, rentalID, model.RentalComplete)
}

func (s *rentalService) CancelRental(ctx context.Context, rentalID, userID int, role string) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"rental not found",
				logrus.Fields{
					"rental_id": rentalID,
				},
			)

			return errs.ErrRentalNotFound
		}

		logger.Error(
			"failed get rental",
			err,
			logrus.Fields{
				"rental_id": rentalID,
			},
		)

		return errs.WrapErr("GetByID", err)
	}

	if role != "admin" && rental.CustomerID != userID {
		logger.Warn(
			"forbidden rental cancellation",
			logrus.Fields{
				"rental_id": rentalID,
				"user_id":   userID,
			},
		)

		return errs.ErrForbidden
	}

	logger.Info(
		"rental cancellation requested",
		logrus.Fields{
			"rental_id": rentalID,
			"user_id":   userID,
			"role":      role,
		},
	)

	return s.UpdateRentalStatus(ctx, rentalID, model.RentalCancelled)
}

func (s *rentalService) UpdateRentalStatus(ctx context.Context, rentalID int, newStatus model.RentalStatus) error {
	return s.txRunner.Run(ctx, func(tx repository.DBTX) error {
		rentalRepo := repository.NewRentalRepository(tx)
		productRepo := repository.NewProductRepository(tx)

		rental, err := rentalRepo.GetByID(ctx, rentalID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				return errs.ErrRentalNotFound
			}

			return errs.WrapErr("GetByID", err)
		}

		oldStatus := rental.Status

		if oldStatus == model.RentalComplete ||
			oldStatus == model.RentalCancelled {
			return errs.ErrInvalidInput
		}

		if !isValidTransition(oldStatus, newStatus) {
			return errs.ErrInvalidInput
		}

		product, err := productRepo.GetByID(ctx, rental.ProductID)
		if err != nil {
			return errs.WrapErr("GetByID", err)
		}

		if err := applyStockChange(product, oldStatus, newStatus); err != nil {
			return err
		}

		if err := productRepo.Update(ctx, product); err != nil {
			return errs.WrapErr("Update", err)
		}

		if err := rentalRepo.UpdateStatus(ctx, rentalID, newStatus); err != nil {
			return errs.WrapErr("UpdateStatus", err)
		}

		return nil
	})
}
