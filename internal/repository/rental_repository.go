package repository

import (
	"context"

	"rentalin/internal/model"

	"github.com/jackc/pgx/v5"
)

type RentalRepository interface {
	Create(ctx context.Context, rental *model.Rental) error
	GetByID(ctx context.Context, id int) (*model.Rental, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.Rental, error)
	GetAll(ctx context.Context) ([]*model.Rental, error)
	UpdateStatus(ctx context.Context, id int, status model.RentalStatus) error
	UpdateOverdueRentals(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	CountByStatus(ctx context.Context, status model.RentalStatus) (int, error)
}

type rentalRepository struct {
	db DBTX
}

func NewRentalRepository(db DBTX) RentalRepository {
	return &rentalRepository{
		db: db,
	}
}

func scanRental(row pgx.Row) (*model.Rental, error) {
	rental := new(model.Rental)

	err := row.Scan(
		&rental.ID,
		&rental.CustomerID,
		&rental.ProductID,
		&rental.CreatedBy,
		&rental.StartDate,
		&rental.EndDate,
		&rental.TotalPrice,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return rental, nil
}

func (r *rentalRepository) Create(ctx context.Context, rental *model.Rental) error {
	query := `
		INSERT INTO rentals (
			customer_id,
			product_id,
			created_by,
			start_date,
			end_date,
			total_price,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		rental.CustomerID,
		rental.ProductID,
		rental.CreatedBy,
		rental.StartDate,
		rental.EndDate,
		rental.TotalPrice,
		rental.Status,
		rental.CreatedAt,
		rental.UpdatedAt,
	).Scan(&rental.ID)
}

func (r *rentalRepository) GetByID(ctx context.Context, id int) (*model.Rental, error) {
	query := `
		SELECT
			id,
			customer_id,
			product_id,
			created_by,
			start_date,
			end_date,
			total_price,
			status,
			created_at,
			updated_at
		FROM rentals
		WHERE id = $1
	`

	return scanRental(r.db.QueryRow(ctx, query, id))
}

func (r *rentalRepository) GetByUserID(ctx context.Context, userID int) ([]*model.Rental, error) {
	query := `
		SELECT
			id,
			customer_id,
			product_id,
			created_by,
			start_date,
			end_date,
			total_price,
			status,
			created_at,
			updated_at
		FROM rentals
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []*model.Rental

	for rows.Next() {
		rental := new(model.Rental)

		err := rows.Scan(
			&rental.ID,
			&rental.CustomerID,
			&rental.ProductID,
			&rental.CreatedBy,
			&rental.StartDate,
			&rental.EndDate,
			&rental.TotalPrice,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rentals = append(rentals, rental)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rentals, rows.Err()
}

func (r *rentalRepository) GetAll(ctx context.Context) ([]*model.Rental, error) {
	query := `
		SELECT
			id,
			customer_id,
			product_id,
			created_by,
			start_date,
			end_date,
			total_price,
			status,
			created_at,
			updated_at
		FROM rentals
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []*model.Rental

	for rows.Next() {
		rental := new(model.Rental)

		err := rows.Scan(
			&rental.ID,
			&rental.CustomerID,
			&rental.ProductID,
			&rental.CreatedBy,
			&rental.StartDate,
			&rental.EndDate,
			&rental.TotalPrice,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rentals = append(rentals, rental)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rentals, nil
}

func (r *rentalRepository) UpdateStatus(ctx context.Context, id int, status model.RentalStatus) error {
	query := `
		UPDATE rentals
		SET 
			status = $1, 
			updated_at = NOW()
		WHERE id = $2
	`

	ct, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *rentalRepository) UpdateOverdueRentals(ctx context.Context) error {
	query := `
		UPDATE rentals
		SET
			status = 'overdue',
			updated_at = NOW()
		WHERE status = 'ongoing'
			AND end_date < CURRENT_DATE
	`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *rentalRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM rentals`

	var total int

	err := r.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *rentalRepository) CountByStatus(ctx context.Context, status model.RentalStatus) (int, error) {
	query := `SELECT COUNT(*) FROM rentals WHERE status = $1`

	var total int

	err := r.db.QueryRow(ctx, query, status).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
