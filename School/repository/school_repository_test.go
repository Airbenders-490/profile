package repository_test

import (
	"context"
	"errors"
	"github.com/airbenders/profile/School/repository"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/pgxmocks"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetConfirmationByToken(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	columns := []string{"token", "sc_id", "st_id", "created_at"}
	expectedToken := domain.Confirmation{Token: "123", School: domain.School{ID: "abc"}, Student: domain.Student{ID: "def"}, CreatedAt: time.Now()}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(
		expectedToken.Token,
		expectedToken.School.ID,
		expectedToken.Student.ID,
		expectedToken.CreatedAt).ToPgxRows()

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf("string")).Return(pgxRows, nil)
		sr := repository.NewSchoolRepository(mockPool)
		token, err := sr.GetConfirmationByToken(context.Background(), "a")

		assert.NoError(t, err)
		assert.EqualValues(t, expectedToken, *token)
	})

	t.Run("query-return-err", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf("string")).
			Return(nil, errors.New("err"))
		sr := repository.NewSchoolRepository(mockPool)
		token, err := sr.GetConfirmationByToken(context.Background(), "a")

		assert.Error(t, err)
		assert.Nil(t, token)
	})
}

func TestSearchByDomain(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	columns := []string{"s.id", "s.name", "s.country"}
	schools := []domain.School{{"a", "b", "c", nil}, {"d", "e", "f", nil}}
	pgRows := pgxpoolmock.NewRows(columns).
		AddRow(schools[0].ID, schools[0].Name, schools[0].Country).
		AddRow(schools[1].ID, schools[1].Name, schools[1].Country).
		ToPgxRows()

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgRows, nil)
		sr := repository.NewSchoolRepository(mockPool)
		returnedSchools, err := sr.SearchByDomain(context.Background(), "sth")

		assert.NoError(t, err)
		assert.EqualValues(t, schools, returnedSchools)
	})

	t.Run("failure", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
		sr := repository.NewSchoolRepository(mockPool)
		returnedSchools, err := sr.SearchByDomain(context.Background(), "sth")

		assert.Error(t, err)
		assert.Nil(t, returnedSchools)
	})
}

func TestSaveConfirmationToken(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	txMock := new(pgxmocks.TxMock)

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()
		txMock.On("Commit", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.SaveConfirmationToken(context.Background(), &domain.Confirmation{})

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.SaveConfirmationToken(context.Background(), &domain.Confirmation{})

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.SaveConfirmationToken(context.Background(), &domain.Confirmation{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.SaveConfirmationToken(context.Background(), &domain.Confirmation{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})
}

func TestAddSchoolForStudent(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	txMock := new(pgxmocks.TxMock)

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()
		txMock.On("Commit", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.AddSchoolForStudent(context.Background(), "a", "b")

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.AddSchoolForStudent(context.Background(), "a", "b")

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.AddSchoolForStudent(context.Background(), "a", "b")

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewSchoolRepository(mockPool)
		err := sr.AddSchoolForStudent(context.Background(), "a", "b")

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})
}
