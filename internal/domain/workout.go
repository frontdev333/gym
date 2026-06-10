package domain

import "time"

type Workout struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	ExerciseID  string    `json:"exercise_id" db:"exercise_id"`
	PerformedAt time.Time `json:"performed_at" db:"performed_at"`
	Amount      *int64    `json:"amount,omitempty" db:"amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
