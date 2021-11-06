package http_test

import (
	"encoding/json"
	"fmt"
	"github.com/airbenders/profile/School/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	e "github.com/airbenders/profile/utils/errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestSearchStudentSchool(t *testing.T){
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	server := httptest.NewServer(app.Server(nil, h))
	defer server.Close()

	var mockSchool domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("case success", func(t *testing.T){
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(mockSchool, nil).
			Once()
		var url = server.URL + "/school"
		response, err := server.Client().Get(url)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 200, response.StatusCode)
	})
}


func TestSendConfirmationMail(t *testing.T){
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)

	server := httptest.NewServer(app.Server(nil, h))
	defer server.Close()
	var mockSchool *domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("SendConfirmation", mock.Anything,
			mock.AnythingOfType("*domain.Student"), mock.AnythingOfType("string"),
		mock.AnythingOfType("*domain.School")).Return(nil).Once()
		postBody, err := json.Marshal(mockSchool)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		response, err := server.Client().Post(fmt.Sprintf("%s/school/confirm?email=%s", server.URL, "adam.yafout@gmail.com"),
			"application/JSON", reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 200)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-body", func(t *testing.T) {
		reader := strings.NewReader("Invalid body")
		response, err := server.Client().Post(fmt.Sprintf("%s/school/confirm?email=%s", server.URL, "adam.yafout@gmail.com"),
			"application/JSON", reader)
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
		assert.EqualValues(t, &restError, e.NewBadRequestError("must provide valid data"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-email", func(t *testing.T) {
		mockUseCase := new(mocks.SchoolUseCase)
		h := http.NewSchoolHandler(mockUseCase)
		r := app.Server(nil, h)
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf("/school/confirm?email=%s", ""), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 404, w.Code)
	})


}