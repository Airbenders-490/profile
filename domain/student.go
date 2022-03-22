package domain

import (
	"context"
	"time"
)

// Student struct
type Student struct {
	ID             string   `json:"id"` //uuid string
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Email          string   `json:"email"` //TODO: validate:required
	GeneralInfo    string   `json:"general_info"`
	School         *School  `json:"school"`
	CurrentClasses []string `json:"current_classes" faker:"-"`
	ClassesTaken   []string `json:"classes_taken" faker:"-"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Reviews        []Review `json:"reviews" faker:"-"`
}

// StudentUseCase interface defines the functions all studentUseCases should have
type StudentUseCase interface {
	Create(ctx context.Context, st *Student) error
	GetByID(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, id string, st *Student) (*Student, error)
	Delete(ctx context.Context, id string) error
	AddClasses(c context.Context, id string, st *Student) error
	RemoveClasses(c context.Context, id string, st *Student) error
	CompleteClass(c context.Context, id string, st *Student) error
	SearchStudents(ctx context.Context, st *Student) ([]Student, error)
	CreateStudentTopic()
	UpdateStudentTopic()
	DeleteStudentTopic()
}

// StudentRepository interface defines the functions all studentRepositories should have
type StudentRepository interface {
	Create(ctx context.Context, id string, st *Student) error
	GetByID(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, st *Student) error
	Delete(ctx context.Context, id string) error
	UpdateClasses(ctx context.Context, st *Student) error
	SearchStudents(ctx context.Context, st *Student) ([]Student, error)
}
