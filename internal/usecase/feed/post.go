package feed

import (
	"context"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
)

func (f *Feed) AddPostEvent(ctx context.Context, postEvent entity.PostEvent) {
	f.postsEventsChannel <- postEvent
}

func (f *Feed) RunPostProcessor(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case postEvent := <-f.postsEventsChannel:
			if postEvent.Event == entity.PostCreate {
				if err := f.newPost(ctx, &postEvent.Post); err != nil {
					return err
				}
			} else if postEvent.Event == entity.PostDelete {
				if err := f.deletePost(ctx, &postEvent.Post); err != nil {
					return err
				}
			}
		}
	}
}

func (f *Feed) newPost(ctx context.Context, post *entity.Post) error {
	user, err := f.usersStorager.Get(ctx, nil, post.AuthorID)
	if err != nil {
		return err
	}
	if user.Star {
		if _, ok := f.star[post.AuthorID]; !ok {
			posts, err := f.postsStorager.List(ctx, nil, &entity.PostFilter{
				AuthorIDs: []uuid.UUID{post.AuthorID},
				Limit:     ptr(1000),
			})
			if err != nil {
				return err
			}
			f.starMu.Lock()
			f.star[post.AuthorID] = posts
			f.starMu.Unlock()
			return nil
		}

		f.starMu.Lock()
		f.star[post.AuthorID] = append(f.star[post.AuthorID], *post)
		if len(f.star[post.AuthorID]) > 1000 {
			f.star[post.AuthorID] = f.star[post.AuthorID][:1000]
		}
		f.starMu.Unlock()
		return nil
	}

	users, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
		FriendID: &post.AuthorID,
	})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	for _, user := range users {
		if _, ok := f.cache[user]; ok {
			f.mu.Lock()
			f.cache[user] = append(f.cache[user], *post)
			if len(f.cache[user]) > 1000 {
				f.cache[user] = f.cache[user][:1000]
			}
			f.mu.Unlock()
			continue
		}

		friends, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
			UserID: &user,
		})
		if err != nil {
			return err
		}

		if len(friends) == 0 {
			continue
		}

		if len(friends) > 100 {
			continue
		}

		posts, err := f.postsStorager.List(ctx, nil, &entity.PostFilter{
			AuthorIDs: friends,
		})
		if err != nil {
			return err
		}

		f.mu.Lock()
		f.cache[user] = posts
		f.mu.Unlock()
	}
	return nil
}

func (f *Feed) deletePost(ctx context.Context, post *entity.Post) error {
	user, err := f.usersStorager.Get(ctx, nil, post.AuthorID)
	if err != nil {
		return err
	}
	if user.Star {
		f.starMu.Lock()
		for i, p := range f.star[post.AuthorID] {
			if p.ID == post.ID {
				f.star[post.AuthorID] = append(f.star[post.AuthorID][:i], f.star[post.AuthorID][i+1:]...)
				break
			}
		}
		f.starMu.Unlock()
		return nil
	}

	users, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
		FriendID: &post.AuthorID,
	})
	if err != nil {
		return err
	}

	for _, user := range users {
		f.mu.Lock()
		for i, p := range f.cache[user] {
			if p.ID == post.ID {
				f.cache[user] = append(f.cache[user][:i], f.cache[user][i+1:]...)
				break
			}
		}
		f.mu.Unlock()
	}
	return nil
}
