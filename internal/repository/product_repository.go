package repository

import (
	"context"
	"errors"

	"rentalin/internal/errs"
	"rentalin/internal/model"

	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int) (*model.Product, error)
	GetAll(ctx context.Context) ([]*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type productRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func scanProduct(row pgx.Row) (*model.Product, error) {
	product := &model.Product{}

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.DailyPrice,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (
			name,
			daily_price,
			stock,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.Name,
		product.DailyPrice,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	query := `
		SELECT
			id,
			name,
			daily_price,
			stock,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
	`

	return scanProduct(r.db.QueryRow(ctx, query, id))
}

func (r *productRepository) GetAll(ctx context.Context) ([]*model.Product, error) {
	query := `
		SELECT
			id,
			name,
			daily_price,
			stock,
			created_at,
			updated_at
		FROM products
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product

	for rows.Next() {
		product := new(model.Product)

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.DailyPrice,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET
			name = $1,
			daily_price = $2,
			stock = $3,
			updated_at = $4
		WHERE id = $5
	`

	ct, err := r.db.Exec(
		ctx,
		query,
		product.Name,
		product.DailyPrice,
		product.Stock,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *productRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM products`

	var total int

	err := r.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
