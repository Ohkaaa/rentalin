package transaction

import (
	"context"
	"rentalin/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TxRunner struct {
	db *pgxpool.Pool
}

func NewTxRunner(db *pgxpool.Pool) *TxRunner {
	return &TxRunner{
		db: db,
	}
}

func (r *TxRunner) Run(ctx context.Context, fn func(tx repository.DBTX) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
