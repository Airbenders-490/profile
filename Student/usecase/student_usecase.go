package usecase

import (
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type studentUseCase struct {
	studentRepository domain.StudentRepository
	//TODO: add review and possibly tag repositories to fetch their data too
	contextTimeout time.Duration
}

func NewStudentUseCase(sr domain.StudentRepository, timeout time.Duration) domain.StudentUseCase {
	return &studentUseCase{
		studentRepository: sr,
		contextTimeout: timeout,
	}
}

func (s *studentUseCase) Create(c context.Context, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if st.ID != "" {
		existingStudent, _ := s.studentRepository.GetById(ctx, st.ID)
		if reflect.DeepEqual(*existingStudent, domain.Student{}) {
			return errors.NewConflictError(fmt.Sprintf("Student with ID %s already exists", st.ID))
		}
	}

	st.ID = uuid.New().String()
	st.CreatedAt = time.Now()
	st.UpdatedAt = time.Now()
	err := s.studentRepository.Create(ctx, st.ID, st)
	if err != nil {
		return err
	}
	return nil
}

func (s *studentUseCase) GetById(c context.Context, id string) (*domain.Student, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	student, err := s.studentRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentUseCase) Update(c context.Context, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existingStudent, _ := s.studentRepository.GetById(ctx, st.ID)
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf("No such student with ID %s exists", st.ID))
	}
	updateStudent(existingStudent, st)
	return s.studentRepository.Update(ctx, existingStudent)
}

func updateStudent(existing *domain.Student, toUpdate *domain.Student) {
	if toUpdate.Name != "" {
		existing.Name = toUpdate.Name
	}
	if toUpdate.Email != "" {
		existing.Email = toUpdate.Email
	}
	if toUpdate.GeneralInfo != "" {
		existing.GeneralInfo = toUpdate.GeneralInfo
	}
	if toUpdate.School != "" {
		existing.School = toUpdate.School
	}
	existing.UpdatedAt = time.Now()
}

func (s *studentUseCase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existingStudent, _ := s.studentRepository.GetById(ctx, id)
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf("No such student with ID %s exists", id))
	}

	return s.studentRepository.Delete(ctx, id)
}