package repository

import (
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"sync"
)

type studentRepository struct {
	repo map[string]domain.Student
	m sync.Mutex
}

func NewStudentRepository() domain.StudentRepository {
	return &studentRepository{
		repo: make(map[string]domain.Student),
	}
}

func (r *studentRepository) Create(ctx context.Context, id string, st *domain.Student) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.repo[id] = *st
	return nil
}

func (r *studentRepository) GetById(ctx context.Context, id string) (*domain.Student, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if val, ok := r.repo[id]; ok {
		return &val, nil
	}
	return nil, errors.NewNotFoundError(fmt.Sprintf("Student with id %s not found.", id))
}

func (r *studentRepository) Update(ctx context.Context, st *domain.Student) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.repo[st.ID]; ok {
		r.repo[st.ID] = *st
		return nil
	}

	return errors.NewNotFoundError(fmt.Sprintf("Student with id %s not found.", st.ID))
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.repo[id]
	if ok {
		delete(r.repo, id)
		return nil
	}
	return errors.NewNotFoundError(fmt.Sprintf("Student with id %s not found.", id))
}
