package vo

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"golang.org/x/crypto/bcrypt"
)

type DecryptedPassword string

func (p DecryptedPassword) Validate() error {
	conditions := map[error]bool{
		exception.ErrEmptyInput: len(p) > 0,
	}

	for err, condition := range conditions {
		if !condition {
			return err
		}
	}

	return nil
}

func (p DecryptedPassword) Encrypt() EncryptedPassword {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return EncryptedPassword(encrypted)
}

type EncryptedPassword string

func (p EncryptedPassword) Compare(decrypted DecryptedPassword) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(decrypted))
	return err == nil
}
