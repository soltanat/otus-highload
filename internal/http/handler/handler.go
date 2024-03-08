package handler

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/http/api"
)

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
}

type User interface {
	Register(ctx context.Context, user *entity.RegisterUser) (*uuid.UUID, error)
	Authenticate(ctx context.Context, userID uuid.UUID, password string) (*entity.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}

type Handler struct {
	userUseCase   User
	tokenProvider TokenProvider
}

func New(userUseCase User, tokenProvider TokenProvider) *Handler {
	if userUseCase == nil || tokenProvider == nil {
		return nil
	}
	return &Handler{
		userUseCase:   userUseCase,
		tokenProvider: tokenProvider,
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
			return api.Login400Response{}, err
		}
		log.Errorf("failed to authenticate user: %v", err)
		return api.Login500JSONResponse{}, err
	}

	token, err := h.tokenProvider.GenerateToken(user.ID.String())
	if err != nil {
		log.Errorf("failed to generate token: %v", err)
		return api.Login500JSONResponse{}, err
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
			return api.GetUser400Response{}, err
		}
		log.Errorf("failed to get user: %v", err)
		return api.GetUser500JSONResponse{}, err
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

func (h *Handler) GetDialog(ctx context.Context, request api.GetDialogRequestObject) (api.GetDialogResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) SendDialogMessage(ctx context.Context, request api.SendDialogMessageRequestObject) (api.SendDialogMessageResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteFriend(ctx context.Context, request api.DeleteFriendRequestObject) (api.DeleteFriendResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) SetFriend(ctx context.Context, request api.SetFriendRequestObject) (api.SetFriendResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CreatePost(ctx context.Context, request api.CreatePostRequestObject) (api.CreatePostResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeletePost(ctx context.Context, request api.DeletePostRequestObject) (api.DeletePostResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetFeed(ctx context.Context, request api.GetFeedRequestObject) (api.GetFeedResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetPost(ctx context.Context, request api.GetPostRequestObject) (api.GetPostResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) UpdatePost(ctx context.Context, request api.UpdatePostRequestObject) (api.UpdatePostResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) SearchUser(ctx context.Context, request api.SearchUserRequestObject) (api.SearchUserResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func stringPtr(v string) *string {
	return &v
}
