package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soltanat/otus-highload/internal/entity"
)

type PostStorage struct {
	conn *pgxpool.Pool
}

func NewPostStorage(conn *pgxpool.Pool) *PostStorage {
	return &PostStorage{
		conn: conn,
	}
}

func (s *PostStorage) Create(ctx context.Context, tx entity.Tx, post *entity.Post) (uuid.UUID, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	_, err := conn.Exec(ctx, "INSERT INTO posts (id, author_id, created_at, text) VALUES ($1, $2, $3, $4)", post.ID, post.AuthorID, post.CreatedAt, post.Text)
	if err != nil {
		return uuid.Nil, entity.StorageError{Err: err}
	}
	return post.ID, nil

}

func (s *PostStorage) Update(ctx context.Context, tx entity.Tx, id uuid.UUID, update *entity.Post) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	_, err := conn.Exec(ctx, "UPDATE posts SET text = $1 WHERE id = $2", update.Text, id)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	return nil
}

func (s *PostStorage) Delete(ctx context.Context, tx entity.Tx, id uuid.UUID) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	_, err := conn.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return entity.StorageError{Err: err}
	}
	return nil
}

func (s *PostStorage) Get(ctx context.Context, tx entity.Tx, id uuid.UUID) (*entity.Post, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	var post entity.Post
	err := conn.QueryRow(ctx, "SELECT id, author_id, created_at, text FROM posts WHERE id = $1", id).Scan(&post.ID, &post.AuthorID, &post.CreatedAt, &post.Text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.NotFoundError{
				Err: fmt.Errorf("post %s not found", id),
			}
		}
		return nil, entity.StorageError{Err: err}
	}
	return &post, nil
}

func (s *PostStorage) List(ctx context.Context, tx entity.Tx, filter *entity.PostFilter) ([]entity.Post, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*PgTx).conn
	}

	//fIDs := make([]string, 0)
	//for _, fID := range filter.AuthorIDs {
	//	fIDs = append(fIDs, fID.String())
	//}

	var rows pgx.Rows
	var err error
	rows, err = conn.Query(
		ctx,
		`SELECT id, author_id, created_at, text FROM posts WHERE author_id = any($1) ORDER BY created_at DESC LIMIT $2`,
		filter.AuthorIDs, filter.Limit,
	)
	if err != nil {
		return nil, entity.StorageError{Err: err}
	}
	defer rows.Close()

	posts := make([]entity.Post, 0)
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.CreatedAt, &post.Text); err != nil {
			return nil, entity.StorageError{Err: err}
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostStorage) Count(ctx context.Context, tx entity.Tx, filter *entity.PostFilter) (int, error) {
	//TODO implement me
	panic("implement me")
}
