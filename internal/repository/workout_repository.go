package repository

import (
	"context"
	"errors"

	"frontdev333/gym/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type workoutRepository struct {
	db *pgxpool.Pool
}

func NewWorkoutRepository(db *pgxpool.Pool) WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(ctx context.Context, workout domain.Workout) (*domain.Workout, error) {
	const query = `
		INSERT INTO workouts (user_id, exercise_id, performed_at, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, exercise_id, performed_at, amount, created_at
	`

	var created domain.Workout
	if err := r.db.QueryRow(
		ctx,
		query,
		workout.UserID,
		workout.ExerciseID,
		workout.PerformedAt,
		workout.Amount,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.ExerciseID,
		&created.PerformedAt,
		&created.Amount,
		&created.CreatedAt,
	); err != nil {
		if isForeignKeyViolation(err) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &created, nil
}

func (r *workoutRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Workout, error) {
	const query = `
		SELECT id, user_id, exercise_id, performed_at, amount, created_at
		FROM workouts
		WHERE user_id = $1
		ORDER BY performed_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := make([]domain.Workout, 0)
	for rows.Next() {
		var workout domain.Workout
		if err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.ExerciseID,
			&workout.PerformedAt,
			&workout.Amount,
			&workout.CreatedAt,
		); err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
