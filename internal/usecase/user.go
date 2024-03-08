package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/interface/storage"
	interfaces "github.com/soltanat/otus-highload/internal/interface/usecase"
)

type User struct {
	userStorager storage.UserStorager
	hasher       interfaces.PasswordHasher
}

func NewUser(userStorager storage.UserStorager, hasher interfaces.PasswordHasher) (*User, error) {
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}

	if hasher == nil {
		return nil, fmt.Errorf("hasher is nil")
	}
	return &User{
		userStorager: userStorager,
		hasher:       hasher,
	}, nil
}

func (u *User) Register(
	ctx context.Context,
	createUser *entity.RegisterUser,
) (*uuid.UUID, error) {
	if createUser.Password == "" {
		return nil, entity.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	hashPassword, err := u.hasher.Hash([]byte(createUser.Password))
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:         uuid.New(),
		FirstName:  createUser.FirstName,
		SecondName: createUser.SecondName,
		BirthDate:  createUser.BirthDate,
		Biography:  createUser.Biography,
		City:       createUser.City,
		Password:   hashPassword,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	err = u.userStorager.Save(ctx, nil, user)
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (u *User) Authenticate(ctx context.Context, userID uuid.UUID, password string) (*entity.User, error) {
	if password == "" {
		return nil, entity.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	user, err := u.userStorager.Get(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	if !u.hasher.Compare(user.Password, []byte(password)) {
		return nil, entity.InvalidPasswordError{}
	}

	return user, nil
}

func (u *User) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.userStorager.Get(ctx, nil, userID)
}
