package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	FirstName  *string
	SecondName *string
	BirthDate  *time.Time
	Biography  *string
	City       *string
	Password   []byte
}

func (u *User) Validate() error {
	if len(u.Password) == 0 {
		return ValidationError{Err: fmt.Errorf("password is empty")}
	}
	return nil
}

type RegisterUser struct {
	FirstName  *string
	SecondName *string
	BirthDate  *time.Time
	Biography  *string
	City       *string
	Password   string
}
