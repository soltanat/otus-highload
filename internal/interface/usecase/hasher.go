package usecase

type PasswordHasher interface {
	Hash(pwd []byte) ([]byte, error)
	Compare(hashedPwd []byte, plainPwd []byte) bool
}
