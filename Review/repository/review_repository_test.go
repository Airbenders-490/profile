package repository_test

import (
	"context"
	"github.com/airbenders/profile/Review/repository"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/airbenders/profile/utils/pgxmocks"
	"github.com/bxcodec/faker"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	_ "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAddReview(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockReviewRepo := new(mocks.ReviewRepositoryMock)
	var mockReview domain.Review
	faker.FakeData(&mockReview)

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	txMock := new(pgxmocks.TxMock)

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()
		txMock.On("Commit", mock.Anything).Return(nil).Once()
		mockReviewRepo.
			On("addTags", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		sr := repository.NewReviewRepository(mockPool)
		err := sr.AddReview(context.Background(), &domain.Review{})

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewReviewRepository(mockPool)
		err := sr.AddReview(context.Background(), &domain.Review{})

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewReviewRepository(mockPool)
		err := sr.AddReview(context.Background(), &domain.Review{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewReviewRepository(mockPool)
		err := sr.AddReview(context.Background(), &domain.Review{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})



}


func TestGetReviewsBy(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reviewColumns := []string{"id", "reviewed", "reviewer","created_at"}
	tagColumns := []string{"name"}
	 mockReview := domain.Review{
		 ID:        "asd",
		 Reviewed:  domain.Student{ID: "123"},
		 Reviewer:  domain.Student{ID: "456"},
		 CreatedAt: time.Now(),
		 Tags:      []domain.Tag{{Name: "some"}, {Name: "thing"}},
	 }

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	t.Run("success", func(t *testing.T) {
		pgxRows := pgxpoolmock.NewRows(reviewColumns).AddRow(
			mockReview.ID ,
			mockReview.Reviewer.ID,
			mockReview.Reviewed.ID,
			mockReview.CreatedAt).ToPgxRows()
		pgxTagRows := pgxpoolmock.NewRows(tagColumns).AddRow(mockReview.Tags[0].Name).AddRow(mockReview.Tags[1].Name).ToPgxRows()

		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(),  gomock.Any()).Return(pgxRows, nil)
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(),  gomock.Any()).Return(pgxTagRows, nil)
		rr := repository.NewReviewRepository(mockPool)

		reviews, err := rr.GetReviewsBy(context.Background(), mockReview.Reviewer.ID)

		assert.NoError(t, err)
		assert.EqualValues(t, []domain.Review{mockReview}, reviews)
	})




}