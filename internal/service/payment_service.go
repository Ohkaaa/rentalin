package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rentalin/internal/errs"
	"rentalin/internal/integration/xendit"
	"rentalin/internal/model"
	"rentalin/internal/repository"
	"rentalin/pkg/logger"
	"rentalin/pkg/transaction"

	"github.com/sirupsen/logrus"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, rentalID, userID int) (*model.Payment, error)
	GetPaymentsByCustomerID(ctx context.Context, userID int) ([]*model.Payment, error)
	GetAllPayments(ctx context.Context) ([]*model.Payment, error)
	HandleInvoicePaid(ctx context.Context, externalID string, paidAmount int64, method string, paymentChannel string, payload []byte) error
	HandleInvoiceExpired(ctx context.Context, externalID string, payload []byte) error
	HandleInvoiceFailed(ctx context.Context, externalID string, payload []byte) error
}

type paymentService struct {
	txRunner     *transaction.TxRunner
	paymentRepo  repository.PaymentRepository
	rentalRepo   repository.RentalRepository
	productRepo  repository.ProductRepository
	xenditClient xendit.Client
}

func NewPaymentService(txRunner *transaction.TxRunner, paymentRepo repository.PaymentRepository, rentalRepo repository.RentalRepository, productRepo repository.ProductRepository, xenditClient xendit.Client) PaymentService {
	return &paymentService{
		txRunner:     txRunner,
		paymentRepo:  paymentRepo,
		rentalRepo:   rentalRepo,
		productRepo:  productRepo,
		xenditClient: xenditClient,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, rentalID, customerID int) (*model.Payment, error) {
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

	if rental.CustomerID != customerID {
		logger.Warn(
			"forbidden payment creation",
			logrus.Fields{
				"customer_id": customerID,
				"rental_id":   rentalID,
			},
		)

		return nil, errs.ErrForbidden
	}

	if rental.Status != model.RentalPending {
		logger.Warn(
			"payment can only be created for pending rental",
			logrus.Fields{
				"rental_id": rentalID,
				"status":    rental.Status,
			},
		)

		return nil, errs.ErrInvalidInput
	}

	externalID := fmt.Sprintf(
		"RENTAL-%d-%d",
		rental.ID,
		time.Now().Unix(),
	)

	resp, err := s.xenditClient.CreateInvoice(
		ctx,
		xendit.CreateInvoiceRequest{
			ExternalID: externalID,
			Amount:     rental.TotalPrice,
			Description: fmt.Sprintf(
				"Rental #%d",
				rental.ID,
			),
		},
	)
	if err != nil {
		logger.Error(
			"failed create xendit invoice",
			err,
			logrus.Fields{
				"rental_id":   rentalID,
				"external_id": externalID,
			},
		)

		return nil, errs.WrapErr("CreateInvoice", err)
	}

	logger.Info(
		"xendit invoice created",
		logrus.Fields{
			"external_id": externalID,
			"invoice_url": resp.InvoiceURL,
		},
	)

	expiredAt, err := time.Parse(time.RFC3339, resp.ExpiredAt)
	if err != nil {
		logger.Error(
			"failed parse invoice expiration",
			err,
			logrus.Fields{
				"external_id": externalID,
			},
		)

		return nil, errs.WrapErr("ParseExpiredAt", err)
	}

	now := time.Now()

	payment := &model.Payment{
		CustomerID: customerID,
		RentalID:   rental.ID,
		ExternalID: externalID,
		InvoiceURL: resp.InvoiceURL,
		Amount:     rental.TotalPrice,
		Currency:   "IDR",
		Status:     model.PaymentPending,
		ExpiredAt:  expiredAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		logger.Error(
			"failed create payment",
			err,
			logrus.Fields{
				"external_id": externalID,
			},
		)

		return nil, errs.WrapErr("Create", err)
	}

	logger.Info(
		"payment created successfully",
		logrus.Fields{
			"payment_id":  payment.ID,
			"rental_id":   rental.ID,
			"customer_id": customerID,
			"amount":      payment.Amount,
		},
	)

	return payment, nil
}

func (s *paymentService) GetPaymentsByCustomerID(ctx context.Context, userID int) ([]*model.Payment, error) {
	payments, err := s.paymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"payments not found",
				logrus.Fields{
					"user_id": userID,
				},
			)

			return nil, errs.ErrPaymentNotFound
		}

		logger.Error(
			"failed get payments",
			err,
			logrus.Fields{
				"user_id": userID,
			},
		)

		return nil, errs.WrapErr("GetByUserID", err)
	}

	return payments, nil
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]*model.Payment, error) {
	payments, err := s.paymentRepo.GetAll(ctx)
	if err != nil {
		logger.Error(
			"failed get all payments",
			err,
			nil,
		)

		return nil, errs.WrapErr("GetAll", err)
	}

	logger.Info(
		"all payments fetched",
		logrus.Fields{
			"total": len(payments),
		},
	)

	return payments, nil
}

func (s *paymentService) HandleInvoicePaid(ctx context.Context, externalID string, paidAmount int64, method string, paymentChannel string, payload []byte) error {
	logger.Info(
		"invoice paid webhook received",
		logrus.Fields{
			"external_id":     externalID,
			"paid_amount":     paidAmount,
			"method":          method,
			"payment_channel": paymentChannel,
		},
	)

	return s.txRunner.Run(ctx, func(tx repository.DBTX) error {
		paymentRepo := repository.NewPaymentRepository(tx)
		rentalRepo := repository.NewRentalRepository(tx)
		productRepo := repository.NewProductRepository(tx)

		payment, err := paymentRepo.GetByExternalID(ctx, externalID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				logger.Warn(
					"payment not found",
					logrus.Fields{
						"external_id": externalID,
					},
				)

				return errs.ErrPaymentNotFound
			}

			logger.Error(
				"failed get payment by external id",
				err,
				logrus.Fields{
					"external_id": externalID,
				},
			)

			return errs.WrapErr("GetByExternalID", err)
		}

		if payment.Status == model.PaymentPaid {
			logger.Info(
				"payment already processed",
				logrus.Fields{
					"payment_id":  payment.ID,
					"external_id": externalID,
				},
			)

			return nil
		}

		rental, err := rentalRepo.GetByID(ctx, payment.RentalID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				logger.Warn(
					"rental not found",
					logrus.Fields{
						"rental_id": payment.RentalID,
					},
				)

				return errs.ErrRentalNotFound
			}

			logger.Error(
				"failed get rental",
				err,
				logrus.Fields{
					"rental_id": payment.RentalID,
				},
			)

			return errs.WrapErr("GetByID", err)
		}

		product, err := productRepo.GetByID(ctx, rental.ProductID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				logger.Warn(
					"product not found",
					logrus.Fields{
						"product_id": rental.ProductID,
					},
				)

				return errs.ErrProductNotFound
			}

			logger.Error(
				"failed get product",
				err,
				logrus.Fields{
					"product_id": rental.ProductID,
				},
			)

			return errs.WrapErr("GetByID", err)
		}

		if product.Stock <= 0 {
			logger.Warn(
				"product out of stock while processing payment",
				logrus.Fields{
					"payment_id": payment.ID,
					"product_id": product.ID,
				},
			)

			return errs.ErrProductOutOfStock
		}

		product.Stock--

		now := time.Now()

		if err := paymentRepo.UpdateStatus(
			ctx,
			payment.ID,
			model.PaymentPaid,
			&paidAmount,
			&method,
			&paymentChannel,
			payload,
			&now,
		); err != nil {
			logger.Error(
				"failed update payment status",
				err,
				logrus.Fields{
					"payment_id": payment.ID,
				},
			)

			return errs.WrapErr("UpdateStatus", err)
		}

		if err := rentalRepo.UpdateStatus(
			ctx,
			rental.ID,
			model.RentalOngoing,
		); err != nil {
			logger.Error(
				"failed update rental status",
				err,
				logrus.Fields{
					"rental_id": rental.ID,
				},
			)

			return errs.WrapErr("UpdateStatus", err)
		}

		if err := productRepo.Update(ctx, product); err != nil {
			logger.Error(
				"failed update product stock",
				err,
				logrus.Fields{
					"product_id": product.ID,
				},
			)

			return errs.WrapErr("Update", err)
		}

		logger.Info(
			"payment processed successfully",
			logrus.Fields{
				"payment_id":      payment.ID,
				"rental_id":       rental.ID,
				"product_id":      product.ID,
				"paid_amount":     paidAmount,
				"payment_channel": paymentChannel,
				"remaining_stock": product.Stock,
			},
		)

		return nil
	})
}

func (s *paymentService) HandleInvoiceExpired(ctx context.Context, externalID string, payload []byte) error {
	payment, err := s.paymentRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"payment not found",
				logrus.Fields{
					"external_id": externalID,
				},
			)

			return errs.ErrPaymentNotFound
		}

		logger.Error(
			"failed get payment",
			err,
			logrus.Fields{
				"external_id": externalID,
			},
		)

		return errs.WrapErr("GetByExternalID", err)
	}

	if payment.Status != model.PaymentPending {
		logger.Info(
			"ignore expired webhook because payment already processed",
			logrus.Fields{
				"payment_id": payment.ID,
				"status":     payment.Status,
			},
		)

		return nil
	}

	if err := s.paymentRepo.UpdateStatus(
		ctx,
		payment.ID,
		model.PaymentExpired,
		nil,
		nil,
		nil,
		payload,
		nil,
	); err != nil {
		logger.Error(
			"failed update payment expired",
			err,
			logrus.Fields{
				"payment_id": payment.ID,
			},
		)

		return errs.WrapErr("UpdateStatus", err)
	}

	logger.Info(
		"payment expired successfully",
		logrus.Fields{
			"payment_id":  payment.ID,
			"external_id": externalID,
		},
	)

	return nil
}

func (s *paymentService) HandleInvoiceFailed(ctx context.Context, externalID string, payload []byte) error {
	payment, err := s.paymentRepo.GetByExternalID(ctx,externalID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"payment not found",
				logrus.Fields{
					"external_id": externalID,
				},
			)

			return errs.ErrPaymentNotFound
		}

		logger.Error(
			"failed get payment",
			err,
			logrus.Fields{
				"external_id": externalID,
			},
		)

		return errs.WrapErr("GetByExternalID", err)
	}

	if payment.Status != model.PaymentPending {
		logger.Info(
			"ignore failed webhook because payment already processed",
			logrus.Fields{
				"payment_id": payment.ID,
				"status":     payment.Status,
			},
		)

		return nil
	}

	if err := s.paymentRepo.UpdateStatus(
		ctx,
		payment.ID,
		model.PaymentFailed,
		nil,
		nil,
		nil,
		payload,
		nil,
	); err != nil {
		logger.Error(
			"failed update payment failed",
			err,
			logrus.Fields{
				"payment_id": payment.ID,
			},
		)

		return errs.WrapErr("UpdateStatus", err)
	}

	logger.Info(
		"payment failed successfully",
		logrus.Fields{
			"payment_id":  payment.ID,
			"external_id": externalID,
		},
	)

	return nil
}
