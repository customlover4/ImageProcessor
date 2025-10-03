package errs

import "errors"

var (
	ErrDBNotFound    = errors.New("not found row")
	ErrDBNotAffected = errors.New("not affected any row")
)
