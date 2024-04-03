package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/usecase/feed"
)

type PostStorager interface {
	Create(ctx context.Context, tx entity.Tx, post *entity.Post) (uuid.UUID, error)
	Update(ctx context.Context, tx entity.Tx, id uuid.UUID, update *entity.Post) error
	Delete(ctx context.Context, tx entity.Tx, id uuid.UUID) error
	Get(ctx context.Context, tx entity.Tx, id uuid.UUID) (*entity.Post, error)
}

type Post struct {
	postStorager PostStorager
	feedUseCase  *feed.Feed
}

func NewPost(postStorager PostStorager, feedUseCase *feed.Feed) *Post {
	return &Post{postStorager: postStorager, feedUseCase: feedUseCase}
}

func (u *Post) Create(ctx context.Context, post *entity.Post) (uuid.UUID, error) {
	post.ID = uuid.New()
	postID, err := u.postStorager.Create(ctx, nil, post)
	if err != nil {
		return uuid.Nil, err
	}

	u.feedUseCase.AddPostEvent(ctx, entity.PostEvent{
		Post:  *post,
		Event: entity.PostCreate,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return postID, nil

}

func (u *Post) Update(ctx context.Context, userID, id uuid.UUID, update *entity.Post) error {
	return u.postStorager.Update(ctx, nil, id, update)
}

func (u *Post) Delete(ctx context.Context, userID, id uuid.UUID) error {
	err := u.postStorager.Delete(ctx, nil, id)
	if err != nil {
		return err
	}

	u.feedUseCase.AddPostEvent(ctx, entity.PostEvent{
		Post: entity.Post{
			ID:       id,
			AuthorID: userID,
		},
		Event: entity.PostDelete,
	})
	return nil
}

func (u *Post) Get(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	return u.postStorager.Get(ctx, nil, id)
}
