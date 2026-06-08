package repository

import (
	"context"
	"time"

	"rentalin/internal/model"

	"github.com/jackc/pgx/v5"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByID(ctx context.Context, id int) (*model.Payment, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.Payment, error)
	GetByExternalID(ctx context.Context, externalID string) (*model.Payment, error)
	GetAll(ctx context.Context) ([]*model.Payment, error)
	UpdateStatus(ctx context.Context, id int, status model.PaymentStatus, paidAmount *int64, method *string, paymentChannel *string, callbackPayload []byte, paidAt *time.Time) error
	ExpirePendingPayments(ctx context.Context) error
	CountByStatus(ctx context.Context, status model.PaymentStatus) (int, error)
	SumPaidRevenue(ctx context.Context) (int64, error)
}

type paymentRepository struct {
	db DBTX
}

func NewPaymentRepository(db DBTX) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func scanPayment(row pgx.Row) (*model.Payment, error) {
	payment := new(model.Payment)

	err := row.Scan(
		&payment.ID,
		&payment.CustomerID,
		&payment.RentalID,
		&payment.ExternalID,
		&payment.InvoiceURL,
		&payment.Amount,
		&payment.PaidAmount,
		&payment.Currency,
		&payment.Method,
		&payment.PaymentChannel,
		&payment.Status,
		&payment.ExpiredAt,
		&payment.PaidAt,
		&payment.Description,
		&payment.CallbackPayload,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (
			customer_id,
			rental_id,
			external_id,
			invoice_url,
			amount,
			paid_amount,
			currency,
			method,
			payment_channel,
			status,
			expired_at,
			paid_at,
			description,
			callback_payload,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4,	$5, $6, $7, $8,	$9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		payment.CustomerID,
		payment.RentalID,
		payment.ExternalID,
		payment.InvoiceURL,
		payment.Amount,
		payment.PaidAmount,
		payment.Currency,
		payment.Method,
		payment.PaymentChannel,
		payment.Status,
		payment.ExpiredAt,
		payment.PaidAt,
		payment.Description,
		payment.CallbackPayload,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id int) (*model.Payment, error) {
	query := `
		SELECT
			id,
			customer_id,
			rental_id,
			external_id,
			invoice_url,
			amount,
			paid_amount,
			currency,
			method,
			payment_channel,
			status,
			expired_at,
			paid_at,
			description,
			callback_payload,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	return scanPayment(r.db.QueryRow(ctx, query, id))
}

func (r *paymentRepository) GetByUserID(ctx context.Context, userID int) ([]*model.Payment, error) {
	query := `
		SELECT
    		id,
			customer_id,
			rental_id,
			external_id,
			invoice_url,
			amount,
			paid_amount,
			currency,
			COALESCE(method, '') as method,
			COALESCE(payment_channel, '') as payment_channel,
			status,
			expired_at,
			paid_at,
			description,
			callback_payload,
			created_at,
			updated_at
		FROM payments
		WHERE customer_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment

	for rows.Next() {
		payment := new(model.Payment)

		err := rows.Scan(
			&payment.ID,
			&payment.CustomerID,
			&payment.RentalID,
			&payment.ExternalID,
			&payment.InvoiceURL,
			&payment.Amount,
			&payment.PaidAmount,
			&payment.Currency,
			&payment.Method,
			&payment.PaymentChannel,
			&payment.Status,
			&payment.ExpiredAt,
			&payment.PaidAt,
			&payment.Description,
			&payment.CallbackPayload,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) GetByExternalID(ctx context.Context, externalID string) (*model.Payment, error) {
	query := `
		SELECT
			id,
			customer_id,
			rental_id,
			external_id,
			invoice_url,
			amount,
			paid_amount,
			currency,
			method,
			payment_channel,
			status,
			expired_at,
			paid_at,
			description,
			callback_payload,
			created_at,
			updated_at
		FROM payments
		WHERE external_id = $1
	`

	return scanPayment(r.db.QueryRow(ctx, query, externalID))
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]*model.Payment, error) {
	query := `
		SELECT
			id,
			customer_id,
			rental_id,
			external_id,
			invoice_url,
			amount,
			paid_amount,
			currency,
			method,
			payment_channel,
			status,
			expired_at,
			paid_at,
			description,
			callback_payload,
			created_at,
			updated_at
		FROM payments
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment

	for rows.Next() {
		payment := new(model.Payment)

		err := rows.Scan(
			&payment.ID,
			&payment.CustomerID,
			&payment.RentalID,
			&payment.ExternalID,
			&payment.InvoiceURL,
			&payment.Amount,
			&payment.PaidAmount,
			&payment.Currency,
			&payment.Method,
			&payment.PaymentChannel,
			&payment.Status,
			&payment.ExpiredAt,
			&payment.PaidAt,
			&payment.Description,
			&payment.CallbackPayload,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id int, status model.PaymentStatus, paidAmount *int64, method *string, paymentChannel *string, callbackPayload []byte, paidAt *time.Time) error {
	query := `
		UPDATE payments
		SET
			status = $1,
			paid_amount = $2,
			method = $3,
			payment_channel = $4,
			callback_payload = $5,
			paid_at = $6,
			updated_at = NOW()
		WHERE id = $7
	`

	ct, err := r.db.Exec(
		ctx,
		query,
		status,
		paidAmount,
		method,
		paymentChannel,
		callbackPayload,
		paidAt,
		id,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *paymentRepository) ExpirePendingPayments(ctx context.Context) error {
	query := `
		WITH expired_payments AS (
			UPDATE payments
			SET
				status = 'expired',
				updated_at = NOW()
			WHERE status = 'pending'
			AND expired_at < NOW()
			RETURNING rental_id
		)

		UPDATE rentals
		SET status = 'cancelled'
		WHERE id IN (
			SELECT rental_id
			FROM expired_payments
		)
		AND status = 'pending'
	`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) CountByStatus(ctx context.Context, status model.PaymentStatus) (int, error) {
	query := `SELECT COUNT(*) FROM payments WHERE status = $1`

	var total int

	err := r.db.QueryRow(
		ctx,
		query,
		status,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *paymentRepository) SumPaidRevenue(ctx context.Context) (int64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM payments
		WHERE status = 'paid'
	`

	var revenue int64

	err := r.db.QueryRow(
		ctx,
		query,
	).Scan(&revenue)

	if err != nil {
		return 0, err
	}

	return revenue, nil
}
