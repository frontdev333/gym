package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/repository"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) Create(ctx context.Context, title string) (*domain.Exercise, error) {
	if err := validateExerciseTitle(title); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, strings.TrimSpace(title))
}

func (s *ExerciseService) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExerciseService) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExerciseService) Update(ctx context.Context, id, title string) (*domain.Exercise, error) {
	if err := validateExerciseTitle(title); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, strings.TrimSpace(title))
}

func (s *ExerciseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func validateExerciseTitle(title string) error {
	trimmed := strings.TrimSpace(title)
	length := utf8.RuneCountInString(trimmed)
	if length < minExerciseTitleLength || length > maxExerciseTitleLength {
		return domain.NewValidationError("title must be between 1 and 200 characters")
	}

	return nil
}
