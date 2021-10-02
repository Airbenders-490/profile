package usecase_test

import (
	"context"
	"errors"
	"github.com/airbenders/profile/Student/usecase"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

/*
type StudentUseCase interface {
	Create(ctx context.Context, st *Student) error
	GetById(ctx context.Context, id string) (*Student, error)
	Update(ctx context.Context, st *Student) error
	Delete(ctx context.Context, id string) error
}
 */

func TestCreate(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockStudent := &domain.Student{
		ID:          "asd",
		Name:        "Sunny",
		Email:       "none@gmail.com",
		GeneralInfo: "I like plants",
		School:      "KGS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now().Add(72*time.Hour),
	}

	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.
			On("Create", mock.Anything, mock.AnythingOfType("string") ,mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		err := u.Create(context.TODO(), mockStudent)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error-in-repo-for-create", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.
			On("Create", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Student")).
			Return(errors.New("error")).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		err := u.Create(context.TODO(), mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error-already-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(mockStudent, nil).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		err := u.Create(context.TODO(), mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockStudent := &domain.Student{
		ID:          "asd",
		Name:        "Sunny",
		Email:       "none@gmail.com",
		GeneralInfo: "I like plants",
		School:      "KGS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now().Add(72*time.Hour),
	}

	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(mockStudent, nil).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		student, err := u.GetById(context.TODO(), mockStudent.ID)

		assert.NoError(t, err)
		assert.NotNil(t, student)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		student, err := u.GetById(context.TODO(), mockStudent.ID)

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(student).IsNil())

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()
		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)

		student, err := u.GetById(context.TODO(), mockStudent.ID)

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(student).IsNil())

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockStudent := &domain.Student{
		ID:          "asd",
		Name:        "Sunny",
		Email:       "none@gmail.com",
		GeneralInfo: "I like plants",
		School:      "KGS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now().Add(72*time.Hour),
	}

	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(mockStudent, nil).
			Once()
		mockStudentRepo.
			On("Update", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Update(context.TODO(), mockStudent)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Update(context.TODO(), mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Update(context.TODO(), mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockStudent := &domain.Student{
		ID:          "asd",
		Name:        "Sunny",
		Email:       "none@gmail.com",
		GeneralInfo: "I like plants",
		School:      "KGS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now().Add(72*time.Hour),
	}

	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(mockStudent, nil).
			Once()
		mockStudentRepo.
			On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Delete(context.TODO(), mockStudent.ID)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Delete(context.TODO(), mockStudent.ID)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetById", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(mockStudentRepo, time.Second)
		err := u.Delete(context.TODO(), mockStudent.ID)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}