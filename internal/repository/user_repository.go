package repository

import (
	"context"
	"errors"

	"rentalin/internal/errs"
	"rentalin/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) UserRepository {
	return &userRepository{
		db: db,
	}
}

func scanUser(row pgx.Row) (*model.User, error) {
	user := &model.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Phone,
		user.Address,
		user.Password,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	return scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	return scanUser(r.db.QueryRow(ctx, query, email))
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE phone = $1
	`

	return scanUser(r.db.QueryRow(ctx, query, phone))
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		FROM users
		WHERE username = $1
	`

	return scanUser(r.db.QueryRow(ctx, query, username))
}

func (r *userRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			phone,
			address,
			password,
			role,
			created_at,
			updated_at
		FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		user := new(model.User)

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			email = $2,
			phone = $3,
			address = $4,
			password = $5,
			role = $6,
			updated_at = $7
		WHERE id = $8
	`

	ct, err := r.db.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Phone,
		user.Address,
		user.Password,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = $1`

	var total int

	err := r.db.QueryRow(ctx, query, "customer").Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
