package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"
)

type User interface {
	Register(ctx context.Context, firstName string, secondName string, birthdate time.Time, biography string, city string, password string) (*uuid.UUID, error)
	Authenticate(ctx context.Context, userID uuid.UUID, password string) (*entity.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}
