package service

import (
	"context"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/davidyannick/repository-pattern/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	return s.repo.AddUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}
