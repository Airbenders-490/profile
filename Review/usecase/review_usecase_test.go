package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/airbenders/profile/Review/usecase"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const casesuccess = "case success"
const stardomainreview = "*domain.Review"
// TestEditReviewsBy function
func TestEditReviewsBy(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockReview domain.Review
	faker.FakeData(&mockReview)


	t.Run(casesuccess, func(t *testing.T) {

		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockReview, nil).
			Once()
		mockReviewRepo.
			On("UpdateReviewTags", mock.Anything, mock.AnythingOfType(stardomainreview)).
			Return(nil).
			Once()
		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)
		review, err := u.EditReview(context.TODO(), &mockReview, mockReview.Reviewer.ID)

		assert.NoError(t, err)
		assert.NotNil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("case student does not exist", func(t *testing.T) {

		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)
		review, err := u.EditReview(context.TODO(), &mockReview, mockReview.Reviewer.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("case review does not exist", func(t *testing.T) {

		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)
		review, err := u.EditReview(context.TODO(), &mockReview, mockReview.Reviewer.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("case failed update", func(t *testing.T) {

		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockReview, nil).
			Once()
		mockReviewRepo.
			On("UpdateReviewTags", mock.Anything, mock.AnythingOfType(stardomainreview)).
			Return(errors.New("error")).
			Once()
		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)
		review, err := u.EditReview(context.TODO(), &mockReview, mockReview.Reviewer.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})

}

// TestGetReviewsBy mock function
func TestGetReviewsBy(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)

	var mockReviews []domain.Review
	faker.FakeData(&mockReviews)

	t.Run(casesuccess, func(t *testing.T) {
		mockReviewRepo.
			On("GetReviewsBy", mock.Anything, mock.AnythingOfType("string")).
			Return(mockReviews, nil).
			Once()
		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		Reviews, err := u.GetReviewsBy(context.TODO(), mockReviews[0].Reviewer.ID)

		assert.NoError(t, err)
		assert.NotNil(t, Reviews)
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("case error", func(t *testing.T) {
		mockReviewRepo.
			On("GetReviewsBy", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		Reviews, err := u.GetReviewsBy(context.TODO(), mockReviews[0].Reviewer.ID)

		assert.Error(t, err)
		assert.Nil(t, Reviews)
		mockReviewRepo.AssertExpectations(t)
	})
}

func TestAddReview(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)

	var mockReview domain.Review
	faker.FakeData(mockReview)

	t.Run(casesuccess, func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil, nil).
			Once()
		mockReviewRepo.
			On("AddReview", mock.Anything, mock.AnythingOfType(stardomainreview)).
			Return(nil).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		review, err := u.AddReview(context.TODO(), &mockReview, mockStudent.ID)

		assert.NoError(t, err)
		assert.NotNil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("case review exists", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockReview, nil).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		review, err := u.AddReview(context.TODO(), &mockReview, mockStudent.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("case student does not exist", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		review, err := u.AddReview(context.TODO(), &mockReview, mockStudent.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
}
