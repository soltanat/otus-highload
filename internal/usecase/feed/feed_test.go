package feed

import (
	"context"
	"github.com/google/uuid"
	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_mergePosts(t *testing.T) {
	id := uuid.New()

	type args struct {
		posts1 []entity.Post
		posts2 []entity.Post
	}
	tests := []struct {
		name string
		args args
		want []entity.Post
	}{
		{
			name: "test",
			args: args{
				posts1: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
				posts2: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
			},
			want: []entity.Post{
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
			},
		},
		{
			name: "test",
			args: args{
				posts1: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
				posts2: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
			},
			want: []entity.Post{
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
			},
		},
		{
			name: "test",
			args: args{
				posts1: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
				posts2: []entity.Post{
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
					{
						ID:        id,
						Text:      "test",
						CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
						AuthorID:  id,
					},
				},
			},
			want: []entity.Post{
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
				{
					ID:        id,
					Text:      "test",
					CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
					AuthorID:  id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergePosts(tt.args.posts1, tt.args.posts2)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestFeed_AllInOne(t *testing.T) {
	ctrl := gomock.NewController(t)

	usersStorager := NewMockUserStorager(ctrl)
	friendsStorager := NewMockFriendStorager(ctrl)
	postsStorager := NewMockPostStorager(ctrl)

	t.Run("get feed for user with star friend from cache", func(t *testing.T) {
		ctx := context.Background()

		userID := uuid.New()
		friendID := uuid.New()
		starFriendID := uuid.New()

		usersStorager.EXPECT().Find(gomock.Any(), nil, &entity.UserFilter{
			Limit:  ptr(100),
			Offset: ptr(0),
		}).Return(
			[]*entity.User{
				{
					ID:   userID,
					Star: false,
				},
				{
					ID:   friendID,
					Star: false,
				},
				{
					ID:   starFriendID,
					Star: true,
				},
			},
			nil,
		)

		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &userID,
			Star:   ptr(false),
		}).Return([]uuid.UUID{
			friendID,
		}, nil)

		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &friendID,
			Star:   ptr(false),
		}).Return([]uuid.UUID{}, nil)

		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &starFriendID,
			Star:   ptr(false),
		}).Return([]uuid.UUID{}, nil)

		postsStorager.EXPECT().List(gomock.Any(), nil, &entity.PostFilter{
			AuthorIDs: []uuid.UUID{friendID},
			Limit:     ptr(1000),
		}).Return([]entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
		}, nil)

		postsStorager.EXPECT().List(gomock.Any(), nil, &entity.PostFilter{
			AuthorIDs: []uuid.UUID{starFriendID},
			Limit:     ptr(1000),
		}).Return([]entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
		}, nil)

		feed := NewFeed(usersStorager, postsStorager, friendsStorager)
		err := feed.Init(ctx)
		require.NoError(t, err)

		// check cache exist feed
		assert.Equal(t, feed.cache[userID], []entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
		})

		assert.Equal(t, feed.star[starFriendID], []entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
		})

		//
		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &userID,
			Star:   ptr(true),
		}).Return([]uuid.UUID{
			starFriendID,
		}, nil)

		posts, err := feed.Get(ctx, userID, 0, 10)
		require.NoError(t, err)
		require.Equal(t, []entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
		}, posts)

		// New post from star

		usersStorager.EXPECT().Get(gomock.Any(), nil, starFriendID).Return(&entity.User{
			ID:   starFriendID,
			Star: true,
		}, nil)

		err = feed.newPost(ctx, &entity.Post{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Text:      "test",
			CreatedAt: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC),
			AuthorID:  starFriendID,
		})
		require.NoError(t, err)

		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &userID,
			Star:   ptr(true),
		}).Return([]uuid.UUID{
			starFriendID,
		}, nil)

		posts, err = feed.Get(ctx, userID, 0, 10)
		require.NoError(t, err)
		require.Equal(t, []entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
		}, posts)

		// New post from user
		usersStorager.EXPECT().Get(gomock.Any(), nil, friendID).Return(&entity.User{
			ID:   friendID,
			Star: false,
		}, nil)
		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			FriendID: &friendID,
		}).Return([]uuid.UUID{
			userID,
		}, nil)

		err = feed.newPost(ctx, &entity.Post{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			Text:      "test",
			CreatedAt: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC),
			AuthorID:  friendID,
		})
		require.NoError(t, err)

		friendsStorager.EXPECT().List(gomock.Any(), nil, &entity.FriendFilter{
			UserID: &userID,
			Star:   ptr(true),
		}).Return([]uuid.UUID{
			starFriendID,
		}, nil)

		posts, err = feed.Get(ctx, userID, 0, 10)
		require.NoError(t, err)
		require.Equal(t, []entity.Post{
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC),
				AuthorID:  starFriendID,
			},
			{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
				Text:      "test",
				CreatedAt: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC),
				AuthorID:  friendID,
			},
		}, posts)
	})
}
