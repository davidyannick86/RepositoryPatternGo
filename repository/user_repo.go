package repository

import (
	"context"

	"github.com/davidyannick/repository-pattern/domain"
)

// UserRepository defines the methods for user data persistence.
type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}
