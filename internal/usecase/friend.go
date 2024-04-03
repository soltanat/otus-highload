package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/usecase/feed"
)

type FriendStorager interface {
	Add(ctx context.Context, tx entity.Tx, userID, friendID uuid.UUID) error
	Delete(ctx context.Context, tx entity.Tx, userID, friendID uuid.UUID) error
}

type Friend struct {
	friendStorager FriendStorager
	feed           *feed.Feed
}

func NewFriend(friendStorager FriendStorager, feed *feed.Feed) *Friend {
	return &Friend{
		friendStorager: friendStorager,
		feed:           feed,
	}
}

func (u *Friend) Add(ctx context.Context, userID, friendID uuid.UUID) error {
	err := u.friendStorager.Add(ctx, nil, userID, friendID)
	if err != nil {
		return err
	}
	u.feed.AddFriendEvent(ctx, entity.FriendEvent{
		Event:    entity.FriendAdd,
		UserID:   userID,
		FriendID: friendID,
	})
	return nil
}

func (u *Friend) Delete(ctx context.Context, userID, friendID uuid.UUID) error {
	err := u.friendStorager.Delete(ctx, nil, userID, friendID)
	if err != nil {
		return err
	}
	u.feed.AddFriendEvent(ctx, entity.FriendEvent{
		Event:    entity.FriendDelete,
		UserID:   userID,
		FriendID: friendID,
	})
	return nil
}
