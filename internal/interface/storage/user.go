package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"
)

type Tx interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UserStorager interface {
	Save(ctx context.Context, tx Tx, user *entity.User) error
	Get(ctx context.Context, tx Tx, userID uuid.UUID) (*entity.User, error)
}
