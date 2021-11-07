package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/airbenders/profile/Student/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	e "github.com/airbenders/profile/utils/errors"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStudentHandler_GetByID(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := http.NewStudentHandler(mockUseCase)
	server := httptest.NewServer(app.Server(h, nil))
	defer server.Close()

	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&mockStudent, nil).Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/student/%s", server.URL, mockStudent.ID))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 200)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		var receivedStudent domain.Student
		err = json.Unmarshal(responseBody, &receivedStudent)
		assert.NoError(t, err)
		// make their times the same. TODO: Find fix
		receivedStudent.UpdatedAt = mockStudent.UpdatedAt
		receivedStudent.CreatedAt = mockStudent.CreatedAt
		assert.EqualValues(t, mockStudent, receivedStudent)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("not-found", func(t *testing.T) {
		mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, e.NewNotFoundError("student not found")).Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/student/%s", server.URL, mockStudent.ID))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 404)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		var restError e.RestError
		err = json.Unmarshal(responseBody, &restError)
		assert.NoError(t, err)
		assert.EqualValues(t, restError, e.RestError{Code: 404, Message: "student not found"})
		mockUseCase.AssertExpectations(t)
	})

	t.Run("some-internal-error", func(t *testing.T) {
		defaultErr := errors.New("some error occurred")
		mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, defaultErr).Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/student/%s", server.URL, "asd"))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 500)
		responseBody, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)
		var receivedResponse httputils.ValidResponse
		err = json.Unmarshal(responseBody, &receivedResponse)
		assert.NoError(t, err)
		assert.EqualValues(t, httputils.ValidResponse{Message: defaultErr.Error()}, receivedResponse)
		mockUseCase.AssertExpectations(t)
	})

	//t.Run("id-not-provided", func(t *testing.T) {
	//	response, err := server.Client().Get(fmt.Sprintf("%s/student/%s", server.URL, ""))
	//	assert.NoError(t, err)
	//	defer response.Body.Close()
	//
	//	assert.Equal(t, response.StatusCode, 400)
	//	responseBody, err := ioutil.ReadAll(response.Body)
	//	fmt.Println(string(responseBody))
	//	assert.NoError(t, err)
	//	var receivedResponse e.RestError
	//	err = json.Unmarshal(responseBody, &receivedResponse)
	//	assert.NoError(t, err)
	//	assert.EqualValues(t, e.NewBadRequestError("id must be provided"), &receivedResponse)
	//	mockUseCase.AssertExpectations(t)
	//})
}

func TestStudentHandler_Create(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	server := httptest.NewServer(app.Server(h, nil))
	defer server.Close()
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(nil).Once()
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		response, err := server.Client().Post(fmt.Sprintf("%s/student", server.URL), "application/JSON", reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 201, response.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-body", func(t *testing.T) {
		reader := strings.NewReader("Invalid body")
		response, err := server.Client().Post(fmt.Sprintf("%s/student", server.URL), "application/JSON", reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 400, response.StatusCode)
		var responseBody []byte
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		var restError e.RestError
		err = json.Unmarshal(responseBody, &restError)
		assert.NoError(t, err)
		assert.EqualValues(t, &restError, e.NewBadRequestError("invalid data"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-create-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("already exists")
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(restErr).Once()
		response, err := server.Client().Post(fmt.Sprintf("%s/student", server.URL), "application/JSON", reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, restErr.Code, response.StatusCode)
		var responseBody []byte
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		var receivedError e.RestError
		err = json.Unmarshal(responseBody, &receivedError)
		assert.NoError(t, err)
		assert.EqualValues(t, restErr, &receivedError)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("some-internal-error", func(t *testing.T) {
		defaultErr := errors.New("some error occurred")
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Student")).
			Return(defaultErr).Once()

		response, err := server.Client().Post(fmt.Sprintf("%s/student", server.URL), "application/JSON", reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 500, response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)
		var receivedResponse e.RestError
		err = json.Unmarshal(responseBody, &receivedResponse)
		assert.NoError(t, err)
		assert.EqualValues(t, e.NewInternalServerError(defaultErr.Error()), &receivedResponse)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandler_Update(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	r := app.Server(h, nil)

	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Student")).
			Return(nil).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/student/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-id", func(t *testing.T) {
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/student/%s", ""), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/student/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Student")).
			Return(restErr).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/student/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Student")).
			Return(errors.New("some error")).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/student/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandler_Delete(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	r := app.Server(h, nil)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf("/student/%s", "asd"), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-id", func(t *testing.T) {
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf("/student/%s", ""), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		restErr := e.NewConflictError("error occurred")
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(restErr).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf("/student/%s", "asdasd"), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(errors.New("some error")).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf("/student/%s", "asd"), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
