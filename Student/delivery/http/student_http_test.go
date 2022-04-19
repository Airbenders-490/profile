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

const failureMessage = "failed to read from message"
const studentType = "*domain.Student"
const putStudentPath = "/api/student/%s"
const postStudentPath = "%s/api/student/"
const getStudentPath = "%s/api/student/%s"

func TestStudentHandlerGetByID(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := http.NewStudentHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	server := httptest.NewServer(r)
	defer server.Close()

	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&mockStudent, nil).Once()

		response, err := server.Client().Get(fmt.Sprintf(getStudentPath, server.URL, mockStudent.ID))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 200, response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
		}
		var receivedStudent domain.Student
		err = json.Unmarshal(responseBody, &receivedStudent)
		assert.NoError(t, err)
		// make their times the same. Find fix
		receivedStudent.UpdatedAt = mockStudent.UpdatedAt
		receivedStudent.CreatedAt = mockStudent.CreatedAt
		assert.EqualValues(t, mockStudent, receivedStudent)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("not-found", func(t *testing.T) {
		mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, e.NewNotFoundError("student not found")).Once()

		response, err := server.Client().Get(fmt.Sprintf(getStudentPath, server.URL, mockStudent.ID))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 404, response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
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

		response, err := server.Client().Get(fmt.Sprintf(getStudentPath, server.URL, "asd"))
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

func TestStudentHandlerCreate(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	server := httptest.NewServer(r)
	defer server.Close()
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	const applicationJSON = "application/JSON"
	t.Run("success", func(t *testing.T) {
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType(studentType)).
			Return(nil).Once()
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		response, err := server.Client().Post(fmt.Sprintf(postStudentPath, server.URL), applicationJSON, reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 201, response.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-body", func(t *testing.T) {
		reader := strings.NewReader("Invalid body")
		response, err := server.Client().Post(fmt.Sprintf(postStudentPath, server.URL), applicationJSON, reader)
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
		assert.EqualValues(t, &restError, e.NewBadRequestError("invalid data"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-create-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("already exists")
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType(studentType)).
			Return(restErr).Once()
		response, err := server.Client().Post(fmt.Sprintf(postStudentPath, server.URL), applicationJSON, reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, restErr.Code, response.StatusCode)
		var responseBody []byte
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, failureMessage)
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
		mockUseCase.On("Create", mock.Anything, mock.AnythingOfType(studentType)).
			Return(defaultErr).Once()

		response, err := server.Client().Post(fmt.Sprintf(postStudentPath, server.URL), applicationJSON, reader)
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

func TestStudentHandlerUpdate(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(&mockStudent, nil).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putStudentPath, mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putStudentPath, mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(nil, restErr).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putStudentPath, mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Student")).
			Return(nil, errors.New("some error")).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putStudentPath, mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandlerDelete(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putStudentPath, "asd"), nil)
		reqFound.Header.Set("id", "asd")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("unauthorized-user", func(t *testing.T) {
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putStudentPath, "a"), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		restErr := e.NewConflictError("error occurred")
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(restErr).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putStudentPath, "asdasd"), nil)
		reqFound.Header.Set("id", "asdasd")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(errors.New("some error")).Once()
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putStudentPath, "asd"), nil)
		reqFound.Header.Set("id", "asd")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandlerAddClasses(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("AddClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(nil).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/addClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/addClasses/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("AddClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(restErr).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/addClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("AddClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(errors.New("error")).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/addClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandlerRemoveClassesTaken(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("RemoveClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(nil).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/removeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/removeClasses/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("RemoveClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(restErr).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/removeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("RemoveClasses", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(errors.New("error")).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/removeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandlerCompleteClass(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := &http.StudentHandler{UseCase: mockUseCase}
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	var mockStudent domain.Student
	err := faker.FakeData(&mockStudent)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("CompleteClass", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(nil).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/completeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/completeClasses/%s", mockStudent.ID), reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("CompleteClass", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(restErr).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/completeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockStudent)
		assert.NoError(t, err)
		mockUseCase.On("CompleteClass", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType(studentType)).
			Return(errors.New("error")).Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/api/completeClasses/%s", mockStudent.ID), reader)
		reqFound.Header.Set("id", mockStudent.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestStudentHandlerSearchStudents(t *testing.T) {
	mockUseCase := new(mocks.StudentUseCase)
	h := http.NewStudentHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	r := app.Server(h, nil, nil, nil, mw, mw, parser)
	var mockRetrievedStudents []domain.Student
	err := faker.FakeData(&mockRetrievedStudents)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("SearchStudents", mock.Anything, mock.Anything).
			Return(mockRetrievedStudents, nil).Once()
		reqFound := httptest.NewRequest("GET", "/api/search/?firstName=Test&lastname=Smith&classes=class", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)

	})

	t.Run("rest error", func(t *testing.T) {
		restErr := e.NewConflictError("error occurred")
		mockUseCase.On("SearchStudents", mock.Anything, mock.Anything).
			Return(nil, restErr).Once()

		reqFound := httptest.NewRequest("GET", "/api/search/?firstName=Test&lastname=Smith&classes=class", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		mockUseCase.On("SearchStudents", mock.Anything, mock.Anything).
			Return(nil, errors.New("internal error")).Once()

		reqFound := httptest.NewRequest("GET", "/api/search/?firstName=Test&lastname=Smith&classes=class", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})

}
