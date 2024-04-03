package handler

import (
	"context"
	"errors"
	"github.com/soltanat/otus-highload/internal/http/middleware"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/http/api"
)

var userIDKey = "subject"

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
}

type User interface {
	Register(ctx context.Context, user *entity.RegisterUser) (*uuid.UUID, error)
	Authenticate(ctx context.Context, userID uuid.UUID, password string) (*entity.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	Search(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error)
}

type Feed interface {
	Get(ctx context.Context, userID uuid.UUID, offset, limit int) ([]entity.Post, error)
}

type Friend interface {
	Add(ctx context.Context, userID, friendID uuid.UUID) error
	Delete(ctx context.Context, userID, friendID uuid.UUID) error
}

type Post interface {
	Create(ctx context.Context, post *entity.Post) (uuid.UUID, error)
	Update(ctx context.Context, userID, id uuid.UUID, update *entity.Post) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Post, error)
}

type Handler struct {
	userUseCase   User
	tokenProvider TokenProvider
	feedUseCase   Feed
	friendUseCase Friend
	postUseCase   Post
}

func New(userUseCase User, tokenProvider TokenProvider, feedUseCase Feed, friendUseCase Friend, postUseCase Post) *Handler {
	if userUseCase == nil || tokenProvider == nil || feedUseCase == nil || friendUseCase == nil || postUseCase == nil {
		return nil
	}

	return &Handler{
		userUseCase:   userUseCase,
		tokenProvider: tokenProvider,
		feedUseCase:   feedUseCase,
		friendUseCase: friendUseCase,
		postUseCase:   postUseCase,
	}
}

func (h *Handler) Register(ctx context.Context, request api.RegisterRequestObject) (api.RegisterResponseObject, error) {
	if request.Body == nil {
		return api.Register400Response{}, nil
	}

	var birthDate *time.Time
	if request.Body.Birthdate != nil {
		birthDate = &request.Body.Birthdate.Time
	}

	if request.Body.Password == nil {
		return api.Register400Response{}, nil
	}

	userID, err := h.userUseCase.Register(
		ctx,
		&entity.RegisterUser{
			FirstName:  request.Body.FirstName,
			SecondName: request.Body.SecondName,
			BirthDate:  birthDate,
			Biography:  request.Body.Biography,
			City:       request.Body.City,
			Password:   *request.Body.Password,
		},
	)

	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.Register400Response{}, err
		}
		log.Errorf("failed to register user: %v", err)
		return api.Register500JSONResponse{}, err
	}

	return api.Register200JSONResponse{
		UserId: stringPtr(userID.String()),
	}, nil

}

func (h *Handler) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	if request.Body == nil {
		return api.Login400Response{}, nil
	}

	if request.Body.Password == nil {
		return api.Login400Response{}, nil
	}

	if request.Body.Id == nil {
		return api.Login400Response{}, nil
	}

	user, err := h.userUseCase.Authenticate(
		ctx,
		uuid.MustParse(*request.Body.Id),
		*request.Body.Password,
	)
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.Login400Response{}, nil
		}
		log.Errorf("failed to authenticate user: %v", err)
		return api.Login500JSONResponse{}, nil
	}

	token, err := h.tokenProvider.GenerateToken(user.ID.String())
	if err != nil {
		log.Errorf("failed to generate token: %v", err)
		return api.Login500JSONResponse{}, nil
	}

	return api.Login200JSONResponse{
		Token: &token,
	}, nil
}

func (h *Handler) GetUser(ctx context.Context, request api.GetUserRequestObject) (api.GetUserResponseObject, error) {
	userID, err := uuid.Parse(request.Id)
	if err != nil {
		return api.GetUser400Response{}, nil
	}

	user, err := h.userUseCase.GetUser(ctx, userID)
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.GetUser400Response{}, nil
		}
		notFoundErr := entity.NotFoundError{}
		if errors.As(err, &notFoundErr) {
			return api.GetUser404Response{}, nil
		}
		log.Errorf("failed to get user: %v", err)
		return api.GetUser500JSONResponse{}, nil
	}

	var birthDate *openapi_types.Date
	if user.BirthDate != nil {
		birthDate = &openapi_types.Date{Time: *user.BirthDate}
	}

	return api.GetUser200JSONResponse{
		Id:         stringPtr(user.ID.String()),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Biography:  user.Biography,
		Birthdate:  birthDate,
		City:       user.City,
	}, nil
}

func (h *Handler) SearchUser(ctx context.Context, request api.SearchUserRequestObject) (api.SearchUserResponseObject, error) {
	if request.Params.FirstName == "" || request.Params.LastName == "" {
		return api.SearchUser400Response{}, nil
	}

	filter := &entity.UserFilter{
		FirstName:  &entity.Filter{Like: request.Params.FirstName},
		SecondName: &entity.Filter{Like: request.Params.LastName},
	}

	users, err := h.userUseCase.Search(ctx, filter)
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.SearchUser400Response{}, err
		}
		log.Errorf("failed to search user: %v", err)
		return api.SearchUser500JSONResponse{}, nil
	}

	response := make([]api.User, len(users))
	for i, user := range users {
		response[i] = api.User{
			Id:         stringPtr(user.ID.String()),
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			Biography:  user.Biography,
			City:       user.City,
		}
	}

	return api.SearchUser200JSONResponse(response), nil
}

func (h *Handler) GetDialog(ctx context.Context, request api.GetDialogRequestObject) (api.GetDialogResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) SendDialogMessage(ctx context.Context, request api.SendDialogMessageRequestObject) (api.SendDialogMessageResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteFriend(ctx context.Context, request api.DeleteFriendRequestObject) (api.DeleteFriendResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	err = h.friendUseCase.Delete(ctx, userID, uuid.MustParse(request.UserId))
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.DeleteFriend400Response{}, nil
		}
		log.Errorf("failed to delete friend: %v", err)
		return api.DeleteFriend500JSONResponse{}, nil
	}

	return api.DeleteFriend200Response{}, nil
}

func (h *Handler) SetFriend(ctx context.Context, request api.SetFriendRequestObject) (api.SetFriendResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	err = h.friendUseCase.Add(ctx, userID, uuid.MustParse(request.UserId))
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.SetFriend400Response{}, nil
		}
		log.Errorf("failed to set friend: %v", err)
		return api.SetFriend500JSONResponse{}, nil
	}

	return api.SetFriend200Response{}, nil
}

func (h *Handler) CreatePost(ctx context.Context, request api.CreatePostRequestObject) (api.CreatePostResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	postID, err := h.postUseCase.Create(ctx, &entity.Post{
		Text:     request.Body.Text,
		AuthorID: userID,
	})
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.CreatePost400Response{}, nil
		}
		log.Errorf("failed to create post: %v", err)
		return api.CreatePost500JSONResponse{}, nil
	}

	return api.CreatePost200JSONResponse(postID.String()), nil
}

func (h *Handler) DeletePost(ctx context.Context, request api.DeletePostRequestObject) (api.DeletePostResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	err = h.postUseCase.Delete(ctx, userID, uuid.MustParse(request.Id))
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.DeletePost400Response{}, nil
		}
		log.Errorf("failed to delete post: %v", err)
		return api.DeletePost500JSONResponse{}, nil
	}

	return api.DeletePost200Response{}, nil
}

func (h *Handler) GetFeed(ctx context.Context, request api.GetFeedRequestObject) (api.GetFeedResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	var offset, limit int
	if request.Params.Offset == nil {
		offset = 0
	} else {
		offset = int(*request.Params.Offset)
	}

	if request.Params.Limit == nil {
		limit = 20
	} else {
		limit = int(*request.Params.Limit)
	}

	feed, err := h.feedUseCase.Get(ctx, userID, offset, limit)
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.GetFeed400Response{}, nil
		}
		log.Errorf("failed to get feed: %v", err)
		return api.GetFeed500JSONResponse{}, nil
	}

	response := make([]api.Post, len(feed))
	for i, post := range feed {
		response[i] = api.Post{
			Id:           stringPtr(post.ID.String()),
			Text:         &post.Text,
			AuthorUserId: stringPtr(post.AuthorID.String()),
		}
	}

	return api.GetFeed200JSONResponse(response), nil
}

func (h *Handler) GetPost(ctx context.Context, request api.GetPostRequestObject) (api.GetPostResponseObject, error) {
	post, err := h.postUseCase.Get(ctx, uuid.MustParse(request.Id))
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.GetPost400Response{}, nil
		}
		log.Errorf("failed to get post: %v", err)
		return api.GetPost500JSONResponse{}, nil
	}

	return api.GetPost200JSONResponse(api.Post{
		Id:           stringPtr(post.ID.String()),
		Text:         &post.Text,
		AuthorUserId: stringPtr(post.AuthorID.String()),
	}), nil
}

func (h *Handler) UpdatePost(ctx context.Context, request api.UpdatePostRequestObject) (api.UpdatePostResponseObject, error) {
	ctxUserID := ctx.Value(middleware.UserIDKeyStruct)

	if ctxUserID == nil {
		return api.GetFeed401Response{}, nil
	}

	userID, err := uuid.Parse(ctxUserID.(string))
	if err != nil {
		return api.GetFeed400Response{}, nil
	}

	err = h.postUseCase.Update(ctx, userID, uuid.MustParse(request.Body.Id), &entity.Post{
		Text: request.Body.Text,
	})
	if err != nil {
		validationErr := entity.ValidationError{}
		if errors.As(err, &validationErr) {
			return api.UpdatePost400Response{}, nil
		}
		log.Errorf("failed to update post: %v", err)
		return api.UpdatePost500JSONResponse{}, nil
	}

	return api.UpdatePost200Response{}, nil
}

func stringPtr(v string) *string {
	return &v
}
