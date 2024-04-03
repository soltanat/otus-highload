package feed

import (
	"context"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
)

func (f *Feed) AddFriendEvent(ctx context.Context, friendEvent entity.FriendEvent) {
	f.friendsEventsChannel <- friendEvent
}

func (f *Feed) RunFriendProcessor(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-f.friendsEventsChannel:
			if event.Event == entity.FriendAdd {
				if err := f.newFriend(ctx, event.UserID, event.FriendID); err != nil {
					return err
				}
			} else if event.Event == entity.FriendDelete {
				if err := f.deleteFriend(ctx, event.UserID, event.FriendID); err != nil {
					return err
				}
			}
		}
	}
}

func (f *Feed) newFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	// Пока просто удаляем фид
	f.mu.Lock()
	delete(f.cache, userID)
	f.mu.Unlock()
	return nil
}

func (f *Feed) deleteFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	// Пока просто удаляем фид
	f.mu.Lock()
	delete(f.cache, userID)
	f.mu.Unlock()
	return nil
}
