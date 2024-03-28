package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/soltanat/otus-highload/internal/entity"
)

const defaultLimit = 10

type UserStorager interface {
	Save(ctx context.Context, tx entity.Tx, user *entity.User) error
	Get(ctx context.Context, tx entity.Tx, userID uuid.UUID) (*entity.User, error)
	Find(ctx context.Context, tx entity.Tx, filter *entity.UserFilter) ([]*entity.User, error)
}

type PasswordHasher interface {
	Hash(pwd []byte) ([]byte, error)
	Compare(hashedPwd []byte, plainPwd []byte) bool
}

type User struct {
	writeUserStorager UserStorager
	readUserStorager  UserStorager
	hasher            PasswordHasher
}

func NewUser(writeUserStorager UserStorager, readUserStorager UserStorager, hasher PasswordHasher) (*User, error) {
	if writeUserStorager == nil {
		return nil, fmt.Errorf("writeUserStorager is nil")
	}
	if readUserStorager == nil {
		return nil, fmt.Errorf("readUserStorager is nil")
	}
	if hasher == nil {
		return nil, fmt.Errorf("hasher is nil")
	}
	return &User{
		writeUserStorager: writeUserStorager,
		readUserStorager:  readUserStorager,
		hasher:            hasher,
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

	err = u.writeUserStorager.Save(ctx, nil, user)
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (u *User) Authenticate(ctx context.Context, userID uuid.UUID, password string) (*entity.User, error) {
	if password == "" {
		return nil, entity.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	user, err := u.writeUserStorager.Get(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	if !u.hasher.Compare(user.Password, []byte(password)) {
		return nil, entity.InvalidPasswordError{}
	}

	return user, nil
}

func (u *User) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.readUserStorager.Get(ctx, nil, userID)
}

func (u *User) Search(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.FirstName.Like == "" && filter.SecondName.Like == "" {
		return nil, fmt.Errorf("filter.FirstName and filter.SecondName are nil")
	}

	if filter.Limit == nil {
		filter.Limit = intPtr(10)
	}

	return u.readUserStorager.Find(ctx, nil, filter)
}

func intPtr(i int) *int {
	return &i
}
