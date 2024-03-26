package usecase

import "golang.org/x/crypto/bcrypt"

type BCryptPasswordHasher struct {
}

func NewPasswordHasher() *BCryptPasswordHasher {
	return &BCryptPasswordHasher{}
}

func (h *BCryptPasswordHasher) Hash(pwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
}

func (h *BCryptPasswordHasher) Compare(hashedPwd []byte, plainPwd []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPwd, plainPwd) == nil
}
