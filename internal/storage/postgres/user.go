package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	conn *pgxpool.Pool
}

func NewUserStorage(conn *pgxpool.Pool) *User {
	return &User{
		conn: conn,
	}
}

func (s *User) Save(ctx context.Context, tx entity.Tx, user *entity.User) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
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

func (s *User) Get(ctx context.Context, tx entity.Tx, userID uuid.UUID) (*entity.User, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	row := conn.QueryRow(
		ctx,
		`SELECT id, first_name, second_name, birthdate, biography, city, password FROM users WHERE id = $1`, userID,
	)
	var user entity.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.BirthDate, &user.Biography, &user.City, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.NotFoundError{Err: fmt.Errorf("user %s not found", userID)}
		}
		return nil, entity.StorageError{Err: err}
	}

	return &user, nil
}

func (s *User) Find(ctx context.Context, tx entity.Tx, filter *entity.UserFilter) ([]*entity.User, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	rows, err := conn.Query(
		ctx,
		`SELECT id, first_name, second_name, birthdate, biography, city, password FROM users WHERE first_name LIKE $1 || '%' AND second_name LIKE $2 || '%' OFFSET $3 LIMIT $4`,
		filter.FirstName.Like, filter.SecondName.Like, filter.Offset, filter.Limit,
	)
	if err != nil {
		return nil, entity.StorageError{Err: err}
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.BirthDate, &user.Biography, &user.City, &user.Password); err != nil {
			return nil, entity.StorageError{Err: err}
		}
		users = append(users, &user)
	}

	return users, nil
}
