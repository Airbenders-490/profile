package usecase

import (
	"context"
	"errors"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSearchSchoolByDomain(t *testing.T) {
	mockSchoolRepo := new(mocks.SchoolRepositoryMock)
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	var mockSchool domain.School
	faker.FakeData(&mockSchool)

	t.Run("case failure empty slice", func(t *testing.T) {
		mockSchoolRepo.
			On("SearchByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return([]domain.School{}, nil).
			Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo , nil , time.Second)

		school, err := u.SearchSchoolByDomain(context.TODO(), mockSchool.Name)

		assert.Nil(t, school)
		assert.Error(t, err)

		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case success", func(t *testing.T) {
		mockSchoolRepo.
			On("SearchByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return([]domain.School{domain.School{}}, nil).
			Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo , nil , time.Second)

		school, err := u.SearchSchoolByDomain(context.TODO(), mockSchool.Name)

		assert.NoError(t, err)
		assert.NotNil(t, school)

		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case error", func(t *testing.T) {
		mockSchoolRepo.
			On("SearchByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo , nil , time.Second)

		school, err := u.SearchSchoolByDomain(context.TODO(), mockSchool.Name)

		assert.Error(t, err)
		assert.Nil(t, school)

		mockSchoolRepo.AssertExpectations(t)
	})
}

func TestSendConfirmation(t *testing.T) {

}

func TestCreateEmailBody(t *testing.T) {

}

func TestConfirmSchoolEnrollment(t *testing.T) {

}

