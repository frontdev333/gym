package repository

import (
	"context"
	"errors"

	"frontdev333/gym/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, name string) (*domain.User, error) {
	const query = `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`

	var user domain.User
	if err := r.db.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
		SELECT id, name, created_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	if err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT id, name, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id, name string) (*domain.User, error) {
	const query = `
		UPDATE users
		SET name = $2
		WHERE id = $1
		RETURNING id, name, created_at
	`

	var user domain.User
	if err := r.db.QueryRow(ctx, query, id, name).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
