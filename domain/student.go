package domain

import (
	"context"
	"time"
)

type Student struct {
	ID string `json:"id"` //uuid string
	Name string `json:"name"`
	Email string `json:"email"` //TODO: validate:required
	GeneralInfo string `json:"general_info"`
	School string `json:"school"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StudentUseCase interface {
	Create(ctx context.Context, st *Student) error
	GetById(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, st *Student) error
	Delete(ctx context.Context, id string) error
}

type StudentRepository interface {
	Create(ctx context.Context, id string, st *Student) error
	GetById(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, st *Student) error
	Delete(ctx context.Context, id string) error
}





