package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/interface/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	conn *pgxpool.Pool
}

func NewUserStorage(conn *pgxpool.Pool) storage.UserStorager {
	return &User{
		conn: conn,
	}
}

func (s *User) Save(ctx context.Context, tx storage.Tx, user *entity.User) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(
		ctx,
		`INSERT INTO users (id, first_name, second_name, birthdate, biography, city, password) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.FirstName, user.SecondName, user.BirthDate, user.Biography, user.City, user.Password,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ExistUserError{}
		}
		return entity.StorageError{Err: err}
	}

	return nil
}

func (s *User) Get(ctx context.Context, tx storage.Tx, userID uuid.UUID) (*entity.User, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	row := conn.QueryRow(
		ctx,
		`SELECT id, first_name, second_name, birthdate, biography, city, password FROM users WHERE id = $1`, userID,
	)
	var user entity.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.BirthDate, &user.Biography, &user.City, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.NotFoundError{}
		}
		return nil, entity.StorageError{Err: err}
	}

	return &user, nil
}
