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

	t.Run("case success", func(t *testing.T) {
		mockSchoolRepo.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo , nil , time.Second)

		school, err := u.SearchSchoolByDomain(context.TODO(), &mockSchool)

		assert.NotNil(t, school)
		assert.NoError(t, err)

		mockSchoolRepo.AssertExpectations(t)
	})
}

func TestSendConfirmation(t *testing.T) {

}

func TestCreateEmailBody(t *testing.T) {

}

func TestConfirmSchoolEnrollment(t *testing.T) {

}

