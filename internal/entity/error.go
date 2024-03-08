package entity

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type InvalidPasswordError struct {
	Err error
}

func (e InvalidPasswordError) Error() string {
	return e.Err.Error()
}

type StorageError struct {
	Err error
}

func (e StorageError) Error() string {
	return e.Err.Error()
}

type ExistUserError struct {
	Err error
}

func (e ExistUserError) Error() string {
	return e.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}
