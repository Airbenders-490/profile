package usecase

import (
	"context"
	"errors"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
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
	mockSchoolRepo := new(mocks.SchoolRepositoryMock)
	mockStudentRepo := new(mocks.StudentRepositoryMock)
	var mockSchool domain.School
	var mockStudent domain.Student
	var mockMailer mocks.SimpleMail
	os.Setenv("DOMAIN", "localhost")
	env := os.Getenv("DOMAIN")
	t.Cleanup(func(){os.Setenv("DOMAIN", env)})
	faker.FakeData(&mockStudent)

	t.Run("case-success", func(t *testing.T){
		//mockMailer := mocks.SimpleMail{}
		faker.FakeData(&mockMailer)
		mockStudent.School = nil
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).Once()
		mockSchoolRepo.
			On("SaveConfirmationToken", mock.Anything, mock.AnythingOfType("*domain.Confirmation")).
			Return(nil).
			Once()
		mockMailer.On("SendSimpleMail", mock.AnythingOfType("string"), mock.Anything).
			Return(nil).Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo, mockMailer, time.Second)

		err := u.SendConfirmation(context.TODO(), &mockStudent, mockStudent.Email, &mockSchool)


		assert.NoError(t, err)
		mockStudentRepo.AssertExpectations(t)
		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case error: School-already-confirmed", func(t *testing.T){

		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo, mockMailer, time.Second)

		err := u.SendConfirmation(context.TODO(), &mockStudent, mockStudent.Email, &mockSchool)

		//assert.True(t, reflect.ValueOf())
		assert.Error(t, err)
	})

	t.Run("case error-empty-student", func(t *testing.T){
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Student{}, nil).Once()

		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo, mockMailer,time.Second)

		err := u.SendConfirmation(context.TODO(), &mockStudent, mockStudent.Email, &mockSchool)
		assert.Error(t, err)
		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("case error-save-confirmation", func(t *testing.T){
		mockStudent.School = nil
		mockStudentRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockStudent, nil).Once()
		mockSchoolRepo.
			On("SaveConfirmationToken", mock.Anything, mock.AnythingOfType("*domain.Confirmation")).
			Return(errors.New("error")).
			Once()
		u := NewSchoolUseCase(mockSchoolRepo, mockStudentRepo, mockMailer, time.Second)

		err := u.SendConfirmation(context.TODO(), &mockStudent, mockStudent.Email, &mockSchool)


		assert.Error(t, err)
		mockStudentRepo.AssertExpectations(t)
		mockSchoolRepo.AssertExpectations(t)
	})
}


func TestConfirmSchoolEnrollment(t *testing.T) {
	mockSchoolRepo := new(mocks.SchoolRepositoryMock)
	//mockStudentRepo := new(mocks.StudentRepositoryMock)
	var mockConfirmation domain.Confirmation
	faker.FakeData(&mockConfirmation)
	//var mockStudent domain.Student
	//faker.FakeData(&mockStudent)

	t.Run("case success", func(t *testing.T){
		mockSchoolRepo.
			On("GetConfirmationByToken", mock.Anything, mock.AnythingOfType("string")).
			Return(&mockConfirmation, nil).Once()
		mockSchoolRepo.
			On("AddSchoolForStudent", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil).Once()
		u := NewSchoolUseCase(mockSchoolRepo, nil, nil, time.Second)
		err := u.ConfirmSchoolEnrollment(context.TODO(), mockConfirmation.Token)

		assert.Nil(t, err)
		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case error-can't-find-token", func(t *testing.T){
		mockSchoolRepo.
			On("GetConfirmationByToken", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error no token")).Once()
		u := NewSchoolUseCase(mockSchoolRepo, nil, nil, time.Second)
		err := u.ConfirmSchoolEnrollment(context.TODO(), mockConfirmation.Token)

		assert.Error(t, err)
		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case error-invalid-token", func(t *testing.T){

		mockSchoolRepo.
			On("GetConfirmationByToken", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Confirmation{}, nil).Once()
		u := NewSchoolUseCase(mockSchoolRepo, nil, nil, time.Second)
		err := u.ConfirmSchoolEnrollment(context.TODO(), mockConfirmation.Token)

		assert.Error(t, err)
		mockSchoolRepo.AssertExpectations(t)
	})

	t.Run("case error-expired-token", func(t *testing.T){
		now := time.Now()
		then := now.Add(-25*time.Hour)
		mockSchoolRepo.
			On("GetConfirmationByToken", mock.Anything, mock.AnythingOfType("string")).
			Return(&domain.Confirmation{CreatedAt: then}, nil).Once()
		u := NewSchoolUseCase(mockSchoolRepo, nil, nil, time.Second)
		err := u.ConfirmSchoolEnrollment(context.TODO(), mockConfirmation.Token)

		assert.Error(t, err)
		mockSchoolRepo.AssertExpectations(t)
	})

}

