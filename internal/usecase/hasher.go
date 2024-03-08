package usecase

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (h *PasswordHasher) Hash(pwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
}

func (h *PasswordHasher) Compare(hashedPwd []byte, plainPwd []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPwd, plainPwd) == nil
}
