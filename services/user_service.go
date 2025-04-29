package service

import (
	"context"
	"fmt"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/davidyannick/repository-pattern/repository"
)

// UserService provides user-related business logic and interacts with the UserRepository.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService with the given UserRepository.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// AddUser adds a new user to the repository.
func (s *UserService) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	u, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}
	return u, nil
}

// GetAllUsers retrieves all users from the repository.
func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}
