package http_test

import (
	"errors"
	"fmt"
	"github.com/airbenders/profile/School/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker"
	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestSchoolHandler_SearchStudentSchool(t *testing.T){
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	server := httptest.NewServer(app.Server(nil, h))
	defer server.Close()

	var mockSchool domain.School
	var arrMockSchool []domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("case success", func(t *testing.T){
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(arrMockSchool, nil).
			Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/?domain=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 200, response.StatusCode)
	})

	t.Run("case no domain", func(t *testing.T){
		response, err := server.Client().Get(fmt.Sprintf("%s/school/?domain=", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()
		assert.Equal(t, 400, response.StatusCode)
	})

	t.Run("case no schools", func(t *testing.T){
		mockUseCase.
			On("SearchSchoolByDomain", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).
			Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/?domain=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()
		assert.Equal(t, 500, response.StatusCode)
	})
}

func TestSchoolHandler_ConfirmSchoolRegistration(t *testing.T) {
	mockUseCase := new(mocks.SchoolUseCase)
	h := http.NewSchoolHandler(mockUseCase)
	server := httptest.NewServer(app.Server(nil, h))
	defer server.Close()

	var mockSchool domain.School
	err := faker.FakeData(&mockSchool)
	assert.NoError(t, err)

	t.Run("case success", func(t *testing.T){
		mockUseCase.
			On("ConfirmSchoolEnrollment", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 200, response.StatusCode)
	})
	t.Run("case error: no token", func(t *testing.T){
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 400, response.StatusCode)
	})
	t.Run("case error token does not exist", func(t *testing.T){
		mockUseCase.
			On("ConfirmSchoolEnrollment", mock.Anything, mock.AnythingOfType("string")).
			Return(errors.New("error")).Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/school/confirmation/?token=test", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 500, response.StatusCode)
	})
}
