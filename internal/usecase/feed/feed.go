package feed

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
	"sync"
)

var ErrCacheMiss = fmt.Errorf("cache miss")

type PostStorager interface {
	List(ctx context.Context, tx entity.Tx, filter *entity.PostFilter) ([]entity.Post, error)
}

type FriendStorager interface {
	List(ctx context.Context, tx entity.Tx, filter *entity.FriendFilter) ([]uuid.UUID, error)
}

type UserStorager interface {
	Get(ctx context.Context, tx entity.Tx, userID uuid.UUID) (*entity.User, error)
	Find(ctx context.Context, tx entity.Tx, filter *entity.UserFilter) ([]*entity.User, error)
}

type Feed struct {
	usersStorager   UserStorager
	postsStorager   PostStorager
	friendsStorager FriendStorager

	cache map[uuid.UUID][]entity.Post
	mu    sync.RWMutex

	star   map[uuid.UUID][]entity.Post
	starMu sync.RWMutex

	postsEventsChannel   chan entity.PostEvent
	friendsEventsChannel chan entity.FriendEvent
}

func NewFeed(
	usersStorager UserStorager,
	postsStorager PostStorager,
	friendsStorager FriendStorager,
) *Feed {
	return &Feed{
		usersStorager:        usersStorager,
		postsStorager:        postsStorager,
		friendsStorager:      friendsStorager,
		cache:                make(map[uuid.UUID][]entity.Post),
		star:                 make(map[uuid.UUID][]entity.Post),
		postsEventsChannel:   make(chan entity.PostEvent, 10),
		friendsEventsChannel: make(chan entity.FriendEvent, 10),
	}
}

func (f *Feed) Get(ctx context.Context, userID uuid.UUID, offset, limit int) ([]entity.Post, error) {
	posts, err := f.getFeed(ctx, userID)
	if err != nil {
		return nil, err
	}

	posts, err = f.mergeStarsFeed(ctx, userID, posts, offset, limit)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (f *Feed) getFeed(ctx context.Context, userID uuid.UUID) ([]entity.Post, error) {
	posts, err := f.getFeedCache(ctx, userID)
	if err == nil {
		return posts, nil
	} else if !errors.Is(err, ErrCacheMiss) {
		return nil, err
	}

	friends, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
		UserID: &userID,
		Star:   ptr(false),
	})
	if err != nil {
		return nil, err
	}

	posts, err = f.postsStorager.List(ctx, nil, &entity.PostFilter{
		AuthorIDs: friends,
		Limit:     ptr(1000),
	})
	if err != nil {
		return nil, err
	}

	f.mu.Lock()
	f.cache[userID] = posts
	f.mu.Unlock()

	return posts, nil
}

func (f *Feed) getFeedCache(_ context.Context, userID uuid.UUID) ([]entity.Post, error) {
	f.mu.RLock()
	posts, ok := f.cache[userID]
	f.mu.RUnlock()

	if !ok {
		return nil, ErrCacheMiss
	}

	return posts, nil
}

func (f *Feed) mergeStarsFeed(ctx context.Context, userID uuid.UUID, posts []entity.Post, offset, limit int) ([]entity.Post, error) {
	//Получаем список его друзей звезд
	starsFriends, err := f.getUserStarFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	//Получаем посты по каждому другу звезде и соединяем их с фидом пользователя
	if len(starsFriends) > 0 {
		for _, friendID := range starsFriends {
			starPosts, err := f.getStarPosts(ctx, friendID)
			if err != nil {
				return nil, err
			}
			posts = mergePosts(posts, starPosts)
		}
	}

	if len(posts) < offset {
		return nil, nil
	}
	if len(posts) <= offset+limit {
		return posts[offset:], nil
	}
	posts = posts[offset : offset+limit]

	return posts, nil
}

func (f *Feed) getUserStarFriends(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	stars, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
		UserID: &userID,
		Star:   ptr(true),
	})
	if err != nil {
		return nil, err
	}
	return stars, nil
}

func (f *Feed) getStarPosts(ctx context.Context, userID uuid.UUID) ([]entity.Post, error) {
	f.starMu.RLock()
	posts, ok := f.star[userID]
	f.starMu.RUnlock()

	if !ok {
		var err error
		posts, err = f.postsStorager.List(ctx, nil, &entity.PostFilter{
			AuthorIDs: []uuid.UUID{userID},
			Limit:     ptr(1000),
		})
		if err != nil {
			return nil, err
		}

		f.starMu.Lock()
		f.star[userID] = posts
		f.starMu.Unlock()
	}

	return posts, nil
}

// Init
// Заполняет кэш для всех пользователей
// Для каждого пользователя получаем посты его друзей
func (f *Feed) Init(ctx context.Context) error {
	var limit = 100
	var offset int
	for {
		users, err := f.usersStorager.Find(ctx, nil, &entity.UserFilter{
			Limit:  &limit,
			Offset: &offset,
		})
		if err != nil {
			return err
		}
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			// Получим друзей пользователя, кроме звезд
			friends, err := f.friendsStorager.List(ctx, nil, &entity.FriendFilter{
				UserID: &user.ID,
				Star:   ptr(false),
			})
			if err != nil {
				return err
			}

			if len(friends) != 0 {
				// Получим посты друзей пользователя
				posts, err := f.postsStorager.List(ctx, nil, &entity.PostFilter{
					AuthorIDs: friends,
					Limit:     ptr(1000),
				})
				if err != nil {
					return err
				}

				// Сохраним в кэш фид пользователя
				f.mu.Lock()
				if _, ok := f.cache[user.ID]; !ok {
					f.cache[user.ID] = posts
				}
				f.mu.Unlock()
			}

			// Если пользователь star, то добавляем его посты в кэш
			if user.Star {
				posts, err := f.postsStorager.List(ctx, nil, &entity.PostFilter{
					AuthorIDs: []uuid.UUID{user.ID},
					Limit:     ptr(1000),
				})
				if err != nil {
					return err
				}
				f.starMu.Lock()
				f.star[user.ID] = posts
				f.starMu.Unlock()
			}
		}

		if len(users) < limit {
			break
		}

		offset += limit
	}
	return nil
}

func mergePosts(posts1, posts2 []entity.Post) []entity.Post {
	posts := make([]entity.Post, 0, len(posts1)+len(posts2))

	for i, j := 0, 0; i != len(posts1) || j != len(posts2); {
		if i == len(posts1) {
			posts = append(posts, posts2[j:]...)
			break
		}
		if j == len(posts2) {
			posts = append(posts, posts1[i:]...)
			break
		}
		if posts1[i].CreatedAt.Before(posts2[j].CreatedAt) {
			posts = append(posts, posts1[i])
			i++
		} else {
			posts = append(posts, posts2[j])
			j++
		}
	}

	return posts
}

func ptr[T any](v T) *T {
	return &v
}
