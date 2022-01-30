package repository_test

import (
	"context"
	"errors"
	"github.com/airbenders/profile/Student/repository"
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

func TestGetByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	columns := []string{"id", "first_name", "last_name", "email", "general_info", "school", "current_classes", "classes_taken", "created_at", "updated_at"}

	t.Run("success-with-nil-school", func(t *testing.T) {
		expectedStudent := &domain.Student{
			ID:          "a",
			FirstName:   "b",
			LastName:    "c",
			Email:       "d",
			GeneralInfo: "e",
			School:      nil,
			CurrentClasses: nil,
			ClassesTaken: nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Reviews:     nil,
		}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(
			expectedStudent.ID,
			expectedStudent.FirstName,
			expectedStudent.LastName,
			expectedStudent.Email,
			expectedStudent.GeneralInfo,
			nil,
			expectedStudent.CurrentClasses,
			expectedStudent.ClassesTaken,
			expectedStudent.CreatedAt,
			expectedStudent.UpdatedAt,
			).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf("string")).Return(pgxRows, nil)
		sr := repository.NewStudentRepository(mockPool)
		student, err := sr.GetByID(context.Background(), "a")

		assert.NoError(t, err)
		assert.EqualValues(t, expectedStudent, student)
	})

	t.Run("success-with-some-school", func(t *testing.T) {
		expectedStudent := domain.Student{
			ID:          "a",
			FirstName:   "b",
			LastName:    "c",
			Email:       "d",
			GeneralInfo: "e",
			School:      &domain.School{ID: "something"},
			CurrentClasses: nil,
			ClassesTaken: nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Reviews:     nil,
		}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(
			expectedStudent.ID,
			expectedStudent.FirstName,
			expectedStudent.LastName,
			expectedStudent.Email,
			expectedStudent.GeneralInfo,
			&expectedStudent.School.ID,
			expectedStudent.CurrentClasses,
			expectedStudent.ClassesTaken,
			expectedStudent.CreatedAt,
			expectedStudent.UpdatedAt).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, nil)
		sr := repository.NewStudentRepository(mockPool)
		student, err := sr.GetByID(context.Background(), "a")

		assert.NoError(t, err)
		assert.EqualValues(t, expectedStudent, *student)
	})

	t.Run("query-return-err", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf("string")).
			Return(nil, errors.New("err"))
		sr := repository.NewStudentRepository(mockPool)
		student, err := sr.GetByID(context.Background(), "a")

		assert.Error(t, err)
		assert.Nil(t, student)
	})
}

func TestUpdate(t *testing.T) {
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

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Update(context.Background(), &domain.Student{})

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Update(context.Background(), &domain.Student{})

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Update(context.Background(), &domain.Student{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Update(context.Background(), &domain.Student{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
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

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Delete(context.Background(), "A")

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Delete(context.Background(), "A")

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Delete(context.Background(), "A")

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Delete(context.Background(), "A")

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
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

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Create(context.Background(), "a", &domain.Student{})

		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't begin transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("err"))

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Create(context.Background(), "a", &domain.Student{})

		assert.Error(t, err)
	})

	t.Run("can't exec transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Create(context.Background(), "a", &domain.Student{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("can't commit transaction", func(t *testing.T) {
		mockPool.EXPECT().Begin(gomock.Any()).Return(txMock, nil)
		txMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()
		txMock.On("Commit", mock.Anything).Return(errors.New("err")).Once()
		txMock.On("Rollback", mock.Anything).Return(nil).Once()

		sr := repository.NewStudentRepository(mockPool)
		err := sr.Create(context.Background(), "a", &domain.Student{})

		assert.Error(t, err)
		txMock.AssertExpectations(t)
	})
}
