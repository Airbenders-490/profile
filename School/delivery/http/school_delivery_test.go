package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/airbenders/profile/School/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	e "github.com/airbenders/profile/utils/errors"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

const failureMessage = "failed to read from message"
const postSchoolEmailConfirmationPath = "%s/api/school/confirm?email=%s"
const testEmail = "adam.yafout@gmail.com"
const applicationJSON = "application/JSON"

func TestSchoolHandlerSearchStudentSchool(t *testing.T){
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	middleware := new(mocks.MiddlewareMock)
	server := httptest.NewServer(app.Server(nil, h, nil, nil, middleware))
	defer server.Close()

	var mockSchool domain.School
	var arrMockSchool []domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("case success", func(t *testing.T) {
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(arrMockSchool, nil).
			Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/api/school/?domain=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()
		assert.Equal(t, 200, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read response")
		}
		var receivedSchools []domain.School
		err = json.Unmarshal(responseBody, &receivedSchools)
		assert.NoError(t, err)
		assert.EqualValues(t, arrMockSchool, receivedSchools)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("case no domain", func(t *testing.T) {
		response, err := server.Client().Get(fmt.Sprintf("%s/api/school/?domain=", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()
		assert.Equal(t, 400, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		var restError e.RestError
		err = json.Unmarshal(responseBody, &restError)
		assert.NoError(t, err)
		assert.EqualValues(t, restError, e.RestError{Code: 400, Message: "please provide a query"})
		mockUseCase.AssertExpectations(t)
	})

	t.Run("case internal error", func(t *testing.T) {
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/api/school/?domain=test", server.URL))
		defaultErr := errors.New("error")
		assert.NoError(t, err)
		defer response.Body.Close()
		assert.Equal(t, 500, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)
		var receivedResponse httputils.ValidResponse
		err = json.Unmarshal(responseBody, &receivedResponse)
		assert.NoError(t, err)
		assert.EqualValues(t, httputils.ValidResponse{Message: defaultErr.Error()}, receivedResponse)
		mockUseCase.AssertExpectations(t)
	})
}

func TestSchoolHandlerConfirmSchoolRegistration(t *testing.T) {
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	middleware := new(mocks.MiddlewareMock)
	server := httptest.NewServer(app.Server(nil, h, nil, nil, middleware))
	defer server.Close()

	var mockSchool domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("case success", func(t *testing.T) {
		mockUseCase.
			On("ConfirmSchoolEnrollment", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 200, response.StatusCode)
	})
	t.Run("case error: no token", func(t *testing.T) {
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 400, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		var restError e.RestError
		err = json.Unmarshal(responseBody, &restError)
		assert.NoError(t, err)
		assert.EqualValues(t, restError, e.RestError{Code: 400, Message: "must provide a valid email"})
	})
	t.Run("case error internal error", func(t *testing.T) {
		mockUseCase.
			On("ConfirmSchoolEnrollment", mock.Anything, mock.AnythingOfType("string")).
			Return(errors.New("error")).Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=test", server.URL))
		assert.NoError(t, err)
		defaultErr := errors.New("error")
		defer response.Body.Close()

		assert.Equal(t, 500, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)
		var receivedResponse e.RestError
		err = json.Unmarshal(responseBody, &receivedResponse)
		assert.NoError(t, err)
		assert.EqualValues(t, e.NewInternalServerError(defaultErr.Error()), &receivedResponse)
	})
}

func TestSchoolHandlerSendConfirmationMail(t *testing.T){
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	middleware := new(mocks.MiddlewareMock)
	server := httptest.NewServer(app.Server(nil, h, nil, nil, middleware))
	defer server.Close()
	var mockSchool *domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)
	var arrMockSchool []domain.School
	arrMockSchool = append(arrMockSchool, *mockSchool)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(arrMockSchool, nil).
			Once()
		mockUseCase.On("SendConfirmation", mock.Anything,
			mock.AnythingOfType("*domain.Student"), mock.AnythingOfType("string"),
			mock.AnythingOfType("*domain.School")).Return(nil).Once()
		response, err := server.Client().Get(fmt.Sprintf(postSchoolEmailConfirmationPath, server.URL, testEmail))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 200)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-email", func(t *testing.T) {
		invalidEmail := "sth@sth"
		response, err := server.Client().Get(fmt.Sprintf(postSchoolEmailConfirmationPath, server.URL, invalidEmail))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 400, response.StatusCode)
		var responseBody []byte
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		var restError e.RestError
		err = json.Unmarshal(responseBody, &restError)
		assert.NoError(t, err)
		assert.EqualValues(t, &restError, e.NewBadRequestError("please provide a valid email"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("confirmation-failure", func(t *testing.T) {
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(arrMockSchool, nil).
			Once()
		mockUseCase.On("SendConfirmation", mock.Anything,
			mock.AnythingOfType("*domain.Student"), mock.AnythingOfType("string"),
			mock.AnythingOfType("*domain.School")).Return(errors.New("some error")).Once()
		assert.NoError(t, err)
		response, err := server.Client().Get(fmt.Sprintf(postSchoolEmailConfirmationPath, server.URL, testEmail))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 500)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		mockUseCase.AssertExpectations(t)
	})
}
