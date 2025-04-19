package errorsrepo

import "errors"

var (
	ErrNotFound = errors.New("record not found")
)

type RepositoryError struct {
	err  error
	Msg  string
	code int
}

func (e RepositoryError) Error() string {
	return e.Msg
}
