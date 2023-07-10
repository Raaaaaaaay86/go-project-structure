package exception

import "github.com/pkg/errors"

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrEmptyInput            = errors.New("empty input")
	ErrLoginFailed           = errors.New("login failed: username or password is incorrect")
	ErrUserNotFound          = errors.New("user not found")
	ErrWrongPassword         = errors.New("wrong password")
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrEmptyFile             = errors.New("file is empty")
)
