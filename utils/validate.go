package utils

import (
	"errors"
	"net/mail"

	"github.com/davidyannick/repository-pattern/domain"
)

const (
	ErrInvalidUserEmail       = "invalid user email"
	ErrInvalidUserName        = "invalid user name"
	ErrInvalidUserEmailLength = "invalid user email length"
	ErrInvalidUserNameLength  = "invalid user name length"
)

func ValidateUser(user domain.User) error {
	if user.Name == "" {
		return errors.New(ErrInvalidUserName)
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New(ErrInvalidUserEmail)
	}
	return nil
}
