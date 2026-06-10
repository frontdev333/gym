package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/repository"
)

const (
	minUserNameLength      = 1
	maxUserNameLength      = 100
	minExerciseTitleLength = 1
	maxExerciseTitleLength = 200
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, name string) (*domain.User, error) {
	if err := validateUserName(name); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, strings.TrimSpace(name))
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) Update(ctx context.Context, id, name string) (*domain.User, error) {
	if err := validateUserName(name); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, strings.TrimSpace(name))
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func validateUserName(name string) error {
	trimmed := strings.TrimSpace(name)
	length := utf8.RuneCountInString(trimmed)
	if length < minUserNameLength || length > maxUserNameLength {
		return domain.NewValidationError("name must be between 1 and 100 characters")
	}

	return nil
}
