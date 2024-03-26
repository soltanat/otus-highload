package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soltanat/otus-highload/internal/entity"
)

type PgTx struct {
	conn *pgxpool.Pool
	tx   *pgx.Tx
}

func (t *PgTx) Begin(ctx context.Context) error {
	if t.tx != nil {
		return nil
	}
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return err
	}
	t.tx = &tx
	return nil
}

func (t *PgTx) Commit(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("tx is nil")
	}
	err := (*t.tx).Commit(ctx)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	t.tx = nil
	return nil
}

func (t *PgTx) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return fmt.Errorf("tx is nil")
	}
	err := (*t.tx).Rollback(ctx)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	t.tx = nil
	return nil
}
