package usecase_test

import (
	"context"
	"errors"
	"github.com/airbenders/profile/Student/usecase"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	mocks2 "github.com/airbenders/profile/utils/channelmocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)

	const studentType = "*domain.Student"
	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.
			On("Create", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(nil).
			Once()
		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		var student domain.Student
		go func() {
			student = <-mm.Created
			assert.NotNil(t, student)
			return
		}()
		err := u.Create(context.TODO(), &mockStudent)
		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error-in-repo-for-create", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.
			On("Create", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(errors.New("error")).
			Once()
		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)

		err := u.Create(context.TODO(), &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error-already-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)

		err := u.Create(context.TODO(), &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestCreateStudentTopic(t *testing.T) {
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	t.Run("success", func(t *testing.T) {
		channelMock.
			On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"), false, false, mock.Anything).
			Return(nil).
			Twice()

		u := usecase.NewStudentUseCase(mm, nil, nil, time.Second)
		go u.CreateStudentTopic()
		mm.Created <- domain.Student{}
		mm.Created <- domain.Student{}
		// wait a bit so the goroutine runs. This will ensure the other goroutine works (hopefully)
		time.Sleep(10 * time.Millisecond)
		channelMock.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	channelMock.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.Anything)

	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewsFor", mock.Anything, mock.AnythingOfType("string")).
			Return([]domain.Review{domain.Review{}, domain.Review{}}, nil).
			Once()
		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)

		student, err := u.GetByID(context.TODO(), mockStudent.ID)

		assert.NoError(t, err)
		assert.NotNil(t, student)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)

		student, err := u.GetByID(context.TODO(), mockStudent.ID)

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(student).IsNil())

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)

		student, err := u.GetByID(context.TODO(), mockStudent.ID)

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(student).IsNil())

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	channelMock.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.Anything)
	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockStudentRepo.
			On("Update", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		var student domain.Student
		go func() {
			student = <-mm.Edited
			assert.NotNil(t, student)
			return
		}()
		updatedStudent, err := u.Update(context.TODO(), mockStudent.ID, &mockStudent)
		assert.NotNil(t, student)
		assert.NoError(t, err)
		assert.EqualValues(t, mockStudent, *updatedStudent)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		_, err := u.Update(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		_, err := u.Update(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestUpdateStudentTopic(t *testing.T) {
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	t.Run("success", func(t *testing.T) {
		channelMock.
			On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"), false, false, mock.Anything).
			Return(nil).
			Twice()

		u := usecase.NewStudentUseCase(mm, nil, nil, time.Second)
		go u.UpdateStudentTopic()
		mm.Edited <- domain.Student{}
		mm.Edited <- domain.Student{}
		time.Sleep(10 * time.Millisecond)
		channelMock.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	channelMock.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.Anything)
	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockStudentRepo.
			On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		go func() {
			id := <-mm.Deleted
			assert.Equal(t, mockStudent.ID, id)
			return
		}()
		err := u.Delete(context.TODO(), mockStudent.ID)
		assert.NoError(t, err)
		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.Delete(context.TODO(), mockStudent.ID)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(mm, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.Delete(context.TODO(), mockStudent.ID)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestDeleteStudentTopic(t *testing.T) {
	channelMock := new(mocks2.ChannelMock)
	mm := usecase.NewMessagingManager(channelMock)
	t.Run("success", func(t *testing.T) {
		channelMock.
			On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("string"), false, false, mock.Anything).
			Return(nil).
			Twice()

		u := usecase.NewStudentUseCase(mm, nil, nil, time.Second)
		go u.DeleteStudentTopic()
		mm.Deleted <- "asd"
		mm.Deleted <- "cde"
		time.Sleep(10 * time.Millisecond)
		channelMock.AssertExpectations(t)
	})
}

func TestAddClasses(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)

	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockStudentRepo.
			On("UpdateClasses", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.AddClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.AddClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.AddClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestRemoveClasses(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)

	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockStudentRepo.
			On("UpdateClasses", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.RemoveClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.RemoveClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.RemoveClasses(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestCompleteClass(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)

	t.Run("success", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockStudentRepo.
			On("UpdateClasses", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.CompleteClass(context.TODO(), mockStudent.ID, &mockStudent)

		assert.NoError(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-no-student-exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.CompleteClass(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("err-empty-student", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).
			Once()

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)
		err := u.CompleteClass(context.TODO(), mockStudent.ID, &mockStudent)

		assert.Error(t, err)

		mockStudentRepo.AssertExpectations(t)
	})
}

func TestSearchStudents(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var retrievedStudents []domain.Student
	faker.FakeData(&retrievedStudents)

	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.
			On("SearchStudents", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(retrievedStudents, nil).
			Once()
		mockReviewRepo.
			On("GetReviewsFor", mock.Anything, mock.AnythingOfType("string")).
			Return([]domain.Review{domain.Review{}, domain.Review{}}, nil)

		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)

		student, err := u.SearchStudents(context.TODO(), &mockStudent)

		assert.NoError(t, err)
		assert.NotNil(t, student)

		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		mockStudentRepo.
			On("SearchStudents", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil, errors.New("error retrieving students")).
			Once()
		u := usecase.NewStudentUseCase(nil, mockStudentRepo, mockReviewRepo, time.Second)

		student, err := u.SearchStudents(context.TODO(), &mockStudent)

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(student).IsNil())

		mockStudentRepo.AssertExpectations(t)
	})

}
