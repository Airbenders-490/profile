package domain

import (
	"context"
	"time"
)

// Student struct
type Student struct {
	ID          string `json:"id"` //uuid string
	Name        string `json:"name"`
	Email       string `json:"email"` //TODO: validate:required
	GeneralInfo string `json:"general_info"`
	School      string `json:"school"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StudentUseCase interface defines the functions all studentUseCases should have
type StudentUseCase interface {
	Create(ctx context.Context, st *Student) error
	GetByID(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, id string, st *Student) error
	Delete(ctx context.Context, id string) error
}

// StudentRepository interface defines the functions all studentRepositories should have
type StudentRepository interface {
	Create(ctx context.Context, id string, st *Student) error
	GetByID(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, st *Student) error
	Delete(ctx context.Context, id string) error
}
