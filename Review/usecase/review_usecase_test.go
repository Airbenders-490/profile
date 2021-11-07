package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/airbenders/profile/Review/usecase"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*type ReviewUseCase interface {
AddReview(ctx context.Context, review *Review, reviewerID string) (*Review, error)
EditReview(ctx context.Context, review *Review, reviewerID string) (*Review, error)
GetReviewsBy(ctx context.Context, reviewer string) ([]Review, error)
*/

func TestGetReviewsBy(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)

	var mockReviews []domain.Review
	faker.FakeData(&mockReviews)

	t.Run("case success", func(t *testing.T) {
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

}

func TestAddReview(t *testing.T) {
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)

	var mockReview domain.Review
	faker.FakeData(mockReview)

	t.Run("case error when todo...", func(t *testing.T) {
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).
			Once()
		mockReviewRepo.
			On("GetReviewByAndFor", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()

		u := usecase.NewReviewUseCase(mockReviewRepo, mockStudentRepo, time.Second)

		fmt.Println("entering add review")
		review, err := u.AddReview(context.TODO(), &mockReview, mockStudent.ID)

		assert.Error(t, err)
		assert.Nil(t, review)
		mockReviewRepo.AssertExpectations(t)
	})
}
