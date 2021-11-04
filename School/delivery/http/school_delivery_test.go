package http_test

import (
	"github.com/airbenders/profile/School/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
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