package repository

import (
	"context"

	"frontdev333/gym/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, name string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id, name string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type ExerciseRepository interface {
	Create(ctx context.Context, title string) (*domain.Exercise, error)
	GetByID(ctx context.Context, id string) (*domain.Exercise, error)
	GetAll(ctx context.Context) ([]domain.Exercise, error)
	Update(ctx context.Context, id, title string) (*domain.Exercise, error)
	Delete(ctx context.Context, id string) error
}

type WorkoutRepository interface {
	Create(ctx context.Context, workout domain.Workout) (*domain.Workout, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.Workout, error)
}

type StatisticsRepository interface {
	GetUserStatistics(ctx context.Context, userID string) (*domain.Statistics, error)
}
