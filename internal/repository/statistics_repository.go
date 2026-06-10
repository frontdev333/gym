package repository

import (
	"context"
	"errors"

	"frontdev333/gym/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type statisticsRepository struct {
	db *pgxpool.Pool
}

func NewStatisticsRepository(db *pgxpool.Pool) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) GetUserStatistics(ctx context.Context, userID string) (*domain.Statistics, error) {
	if err := r.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	total, err := r.getTotalCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	today, err := r.getTodayCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	last7Days, err := r.getLast7Days(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Statistics{
		Total:     total,
		Today:     today,
		Last7Days: last7Days,
	}, nil
}

func (r *statisticsRepository) ensureUserExists(ctx context.Context, userID string) error {
	const query = `SELECT 1 FROM users WHERE id = $1`

	var exists int
	if err := r.db.QueryRow(ctx, query, userID).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *statisticsRepository) getTotalCount(ctx context.Context, userID string) (int64, error) {
	const query = `SELECT COUNT(*) FROM workouts WHERE user_id = $1`

	var total int64
	if err := r.db.QueryRow(ctx, query, userID).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *statisticsRepository) getTodayCount(ctx context.Context, userID string) (int64, error) {
	const query = `
		SELECT COUNT(*)
		FROM workouts
		WHERE user_id = $1
		  AND performed_at >= CURRENT_DATE
		  AND performed_at < CURRENT_DATE + INTERVAL '1 day'
	`

	var today int64
	if err := r.db.QueryRow(ctx, query, userID).Scan(&today); err != nil {
		return 0, err
	}

	return today, nil
}

func (r *statisticsRepository) getLast7Days(ctx context.Context, userID string) ([]domain.DailyCount, error) {
	const query = `
		WITH days AS (
			SELECT generate_series(
				(CURRENT_DATE - INTERVAL '6 days')::date,
				CURRENT_DATE::date,
				'1 day'::interval
			)::date AS day
		)
		SELECT d.day::text, COUNT(w.id)::bigint AS count
		FROM days d
		LEFT JOIN workouts w
			ON w.user_id = $1 AND w.performed_at::date = d.day
		GROUP BY d.day
		ORDER BY d.day
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dailyCounts := make([]domain.DailyCount, 0, 7)
	for rows.Next() {
		var daily domain.DailyCount
		if err := rows.Scan(&daily.Date, &daily.Count); err != nil {
			return nil, err
		}

		dailyCounts = append(dailyCounts, daily)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dailyCounts, nil
}
