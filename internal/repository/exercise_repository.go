package repository

import (
	"context"
	"errors"

	"frontdev333/gym/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type exerciseRepository struct {
	db *pgxpool.Pool
}

func NewExerciseRepository(db *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(ctx context.Context, title string) (*domain.Exercise, error) {
	const query = `
		INSERT INTO exercises (title)
		VALUES ($1)
		RETURNING id, title, created_at
	`

	var exercise domain.Exercise
	if err := r.db.QueryRow(ctx, query, title).Scan(&exercise.ID, &exercise.Title, &exercise.CreatedAt); err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	const query = `
		SELECT id, title, created_at
		FROM exercises
		WHERE id = $1
	`

	var exercise domain.Exercise
	if err := r.db.QueryRow(ctx, query, id).Scan(&exercise.ID, &exercise.Title, &exercise.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &exercise, nil
}

func (r *exerciseRepository) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	const query = `
		SELECT id, title, created_at
		FROM exercises
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]domain.Exercise, 0)
	for rows.Next() {
		var exercise domain.Exercise
		if err := rows.Scan(&exercise.ID, &exercise.Title, &exercise.CreatedAt); err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exercises, nil
}

func (r *exerciseRepository) Update(ctx context.Context, id, title string) (*domain.Exercise, error) {
	const query = `
		UPDATE exercises
		SET title = $2
		WHERE id = $1
		RETURNING id, title, created_at
	`

	var exercise domain.Exercise
	if err := r.db.QueryRow(ctx, query, id, title).Scan(&exercise.ID, &exercise.Title, &exercise.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &exercise, nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM exercises WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.NewValidationError("cannot delete exercise: it has existing workouts")
		}

		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
