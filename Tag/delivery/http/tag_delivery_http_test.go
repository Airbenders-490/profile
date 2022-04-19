package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/airbenders/profile/Tag/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/airbenders/profile/utils/httputils"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTagHandlerGetAllTags(t *testing.T) {
	mockUseCase := new(mocks.TagUseCase)
	h := http.NewTagHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	parser := new(mocks.ClaimsParserMock)
	server := httptest.NewServer(app.Server(nil, nil, h, nil, mw, mw, parser))
	defer server.Close()

	var mockTag []domain.Tag
	err := faker.FakeData(&mockTag)
	assert.NoError(t, err)
	fmt.Print("sdajklhalksjdhakjlsdh", len(mockTag))
	t.Run("success", func(t *testing.T) {
		mockUseCase.On("GetAllTags", mock.Anything).Return(mockTag, nil).Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/api/all-tags", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, 200)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, "failed to read from message")
		}
		var receivedTags []domain.Tag
		err = json.Unmarshal(responseBody, &receivedTags)
		assert.NoError(t, err)

		assert.EqualValues(t, mockTag, receivedTags)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("some-internal-error", func(t *testing.T) {
		defaultErr := errors.New("some error occurred")
		mockUseCase.On("GetAllTags", mock.Anything).Return(nil, defaultErr).Once()
		response, err := server.Client().Get(fmt.Sprintf("%s/api/all-tags", server.URL))
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
}
