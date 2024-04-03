package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soltanat/otus-highload/internal/entity"
	"strings"
)

type FriendStorage struct {
	conn *pgxpool.Pool
}

func NewFriendStorage(conn *pgxpool.Pool) *FriendStorage {
	return &FriendStorage{
		conn: conn,
	}
}

func (s *FriendStorage) List(ctx context.Context, tx entity.Tx, filter *entity.FriendFilter) ([]uuid.UUID, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	var rows pgx.Rows
	var err error

	query := `SELECT friend_id FROM friends`
	where := make([]string, 0)

	statementIndex := 1
	statements := make([]any, 0)
	if filter.Star != nil {
		where = append(where, fmt.Sprintf("star = $%d", statementIndex))
		statements = append(statements, *filter.Star)
		statementIndex++
	}
	if filter.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", statementIndex))
		statements = append(statements, *filter.UserID)
		statementIndex++
	}
	if filter.FriendID != nil {
		where = append(where, fmt.Sprintf("friend_id = $%d", statementIndex))
		statements = append(statements, *filter.FriendID)
		statementIndex++
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", statementIndex)
		statements = append(statements, *filter.Limit)
		statementIndex++
	}

	rows, err = conn.Query(ctx, query, statements...)
	if err != nil {
		return nil, entity.StorageError{Err: err}
	}
	defer rows.Close()

	friends := make([]uuid.UUID, 0)
	for rows.Next() {
		var friendID uuid.UUID
		if err := rows.Scan(&friendID); err != nil {
			return nil, entity.StorageError{Err: err}
		}
		friends = append(friends, friendID)
	}

	return friends, nil

}

func (s *FriendStorage) Add(ctx context.Context, tx entity.Tx, userID, friendID uuid.UUID) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	_, err := conn.Exec(ctx, "INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", userID, friendID)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	return nil
}

func (s *FriendStorage) Delete(ctx context.Context, tx entity.Tx, userID, friendID uuid.UUID) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	_, err := conn.Exec(ctx, "DELETE FROM friends WHERE user_id = $1 AND friend_id = $2", userID, friendID)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	return nil
}
