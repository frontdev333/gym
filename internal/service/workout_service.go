package service

import (
	"context"
	"time"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/repository"
)

type WorkoutService struct {
	workoutRepo  repository.WorkoutRepository
	userRepo     repository.UserRepository
	exerciseRepo repository.ExerciseRepository
}

func NewWorkoutService(
	workoutRepo repository.WorkoutRepository,
	userRepo repository.UserRepository,
	exerciseRepo repository.ExerciseRepository,
) *WorkoutService {
	return &WorkoutService{
		workoutRepo:  workoutRepo,
		userRepo:     userRepo,
		exerciseRepo: exerciseRepo,
	}
}

type CreateWorkoutInput struct {
	UserID      string
	ExerciseID  string
	PerformedAt *time.Time
	Amount      *int64
}

func (s *WorkoutService) Create(ctx context.Context, input CreateWorkoutInput) (*domain.Workout, error) {
	if input.ExerciseID == "" {
		return nil, domain.NewValidationError("exercise_id is required")
	}

	if input.Amount != nil && *input.Amount <= 0 {
		return nil, domain.NewValidationError("amount must be greater than 0")
	}

	if _, err := s.userRepo.GetByID(ctx, input.UserID); err != nil {
		return nil, err
	}

	if _, err := s.exerciseRepo.GetByID(ctx, input.ExerciseID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	performedAt := now
	if input.PerformedAt != nil {
		performedAt = input.PerformedAt.UTC()
		if performedAt.After(now) {
			return nil, domain.NewValidationError("performed_at cannot be in the future")
		}
	}

	workout := domain.Workout{
		UserID:      input.UserID,
		ExerciseID:  input.ExerciseID,
		PerformedAt: performedAt,
		Amount:      input.Amount,
	}

	return s.workoutRepo.Create(ctx, workout)
}

func (s *WorkoutService) GetByUserID(ctx context.Context, userID string) ([]domain.Workout, error) {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.workoutRepo.GetByUserID(ctx, userID)
}
