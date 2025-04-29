package utils

import (
	"errors"
	"net/mail"

	"github.com/davidyannick/repository-pattern/domain"
)

const (
	// ErrInvalidUserEmail represents an invalid user email error.
	ErrInvalidUserEmail = "invalid user email"
	// ErrInvalidUserName represents an invalid user name error.
	ErrInvalidUserName = "invalid user name"
	// ErrInvalidUserEmailLength represents an invalid user email length error.
	ErrInvalidUserEmailLength = "invalid user email length"
	// ErrInvalidUserNameLength represents an invalid user name length error.
	ErrInvalidUserNameLength = "invalid user name length"
)

// ValidateUser checks if the user name and email are valid.
func ValidateUser(user domain.User) error {
	if user.Name == "" {
		return errors.New(ErrInvalidUserName)
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New(ErrInvalidUserEmail)
	}
	return nil
}
