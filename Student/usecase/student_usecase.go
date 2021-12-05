package usecase

import (
	"context"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"log"
	"reflect"
	"time"
)


const errorMessage = "No such student with ID %s exists"

type studentUseCase struct {
	studentRepository domain.StudentRepository
	// add review and possibly tag repositories to fetch their data too
	reviewRepository domain.ReviewRepository
	contextTimeout   time.Duration
}

// NewStudentUseCase returns a configured StudentUseCase
func NewStudentUseCase(sr domain.StudentRepository, rr domain.ReviewRepository, timeout time.Duration) domain.StudentUseCase {
	return &studentUseCase{
		studentRepository: sr,
		reviewRepository:  rr,
		contextTimeout:    timeout,
	}
}

// Create stores the student in the db. ID must be provided or return an error
func (s *studentUseCase) Create(c context.Context, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if st.ID != "" {
		existingStudent, err := s.studentRepository.GetByID(ctx, st.ID)
		if err == nil {
			if !reflect.DeepEqual(*existingStudent, domain.Student{}) {
				return errors.NewConflictError(fmt.Sprintf(errorMessage, st.ID))
			}
		}
	} else {
		return errors.NewBadRequestError("The student should have an ID from auth service")
	}

	st.CreatedAt = time.Now()
	st.UpdatedAt = time.Now()
	err := s.studentRepository.Create(ctx, st.ID, st)
	if err != nil {
		return err
	}
	return nil
}

// GetByID seeks student from repo layer and returns if it exists, else return error
func (s *studentUseCase) GetByID(c context.Context, id string) (*domain.Student, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	student, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(student, &domain.Student{}) {

		return nil, errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}

	reviews, err := s.reviewRepository.GetReviewsFor(ctx, student.ID)
	if err != nil {
		log.Println("Can't get the reviews right now.")
	}
	student.Reviews = reviews

	return student, nil
}

// Update checks if the student exists and updates if so. Otherwise, returns error
func (s *studentUseCase) Update(c context.Context, id string, st *domain.Student) (*domain.Student, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	st.ID = id
	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, st.ID))
	}
	updateStudent(existingStudent, st)
	return existingStudent, s.studentRepository.Update(ctx, existingStudent)
}

func updateStudent(existing *domain.Student, toUpdate *domain.Student) {
	if toUpdate.FirstName != "" {
		existing.FirstName = toUpdate.FirstName
	}
	if toUpdate.LastName != "" {
		existing.LastName = toUpdate.LastName
	}
	if toUpdate.Email != "" {
		existing.Email = toUpdate.Email
	}
	if toUpdate.GeneralInfo != "" {
		existing.GeneralInfo = toUpdate.GeneralInfo
	}
	existing.UpdatedAt = time.Now()
}

// Delete removes the student if it exists. Otherwise, returns error
func (s *studentUseCase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}

	return s.studentRepository.Delete(ctx, id)
}
