package vo

import (
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"net/mail"
)

type Email string

func (e Email) Validate() error {
	_, err := mail.ParseAddress(string(e))
	if err != nil {
		return exception.ErrInvalidEmail
	}

	return err
}
