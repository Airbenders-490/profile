package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/airbenders/profile/Review/delivery/http"
	"github.com/airbenders/profile/app"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	e "github.com/airbenders/profile/utils/errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const reviewType = "*domain.Review"
const postReviewPath = "%s/api/review/%s"
const applicationJSON = "application/JSON"
const failureMessage = "failed to read from message"
const putReviewPath = "/api/review/%s/update"

func TestReviewHandlerAddReview(t *testing.T) {
	mockUseCase := new(mocks.ReviewUseCase)
	h := http.NewReviewHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	server := httptest.NewServer(app.Server(nil, nil, nil, h, mw))
	defer server.Close()

	var mockReview domain.Review
	err := faker.FakeData(&mockReview)
	assert.NoError(t, err)
	var mockStudent domain.Student
	err = faker.FakeData(&mockStudent)
	assert.NoError(t, err)


	t.Run("success", func(t *testing.T) {
		mockUseCase.On("AddReview", mock.Anything, mock.AnythingOfType(reviewType), mock.AnythingOfType("string")).
			Return(&mockReview, nil).
			Once()

		postBody, err := json.Marshal(&mockReview)

		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))

		response, err := server.Client().Post(
			fmt.Sprintf(postReviewPath, server.URL, mockReview.Reviewer.ID),
			applicationJSON, reader)

		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, 201, response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)

		if err != nil {
			assert.Fail(t, failureMessage)
		}

		var receivedReview domain.Review
		err = json.Unmarshal(responseBody, &receivedReview)

		assert.NoError(t, err)
		// make their times the same. TODO: Find fix
		receivedReview.CreatedAt = mockReview.CreatedAt
		receivedReview.Reviewed.CreatedAt = mockReview.Reviewed.CreatedAt
		receivedReview.Reviewer.CreatedAt = mockReview.Reviewer.CreatedAt
		receivedReview.Reviewed.UpdatedAt = mockReview.Reviewed.UpdatedAt
		receivedReview.Reviewer.UpdatedAt = mockReview.Reviewer.UpdatedAt

		assert.EqualValues(t, mockReview, receivedReview)
		mockUseCase.AssertExpectations(t)
	})
	t.Run("invalid body", func(t *testing.T) {

		reader := strings.NewReader("Invalid body")

		response, err := server.Client().Post(
			(fmt.Sprintf(postReviewPath, server.URL, mockReview.Reviewer.ID)),
			applicationJSON, reader)

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
		assert.EqualValues(t, &restError, e.NewBadRequestError("invalid review body"))
		mockUseCase.AssertExpectations(t)

	})
	t.Run("already exist", func(t *testing.T) {
		restErr := e.NewConflictError("already exists")
		mockUseCase.On("AddReview", mock.Anything, mock.AnythingOfType(reviewType), mock.AnythingOfType("string")).
			Return(nil, restErr).
			Once()
		postBody, err := json.Marshal(&mockReview)

		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))

		response, err := server.Client().Post(
			(fmt.Sprintf(postReviewPath, server.URL, mockReview.Reviewer.ID)),
			applicationJSON, reader)
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
		mockUseCase.On("AddReview", mock.Anything, mock.AnythingOfType(reviewType), mock.AnythingOfType("string")).
			Return(nil, defaultErr).
			Once()
		postBody, err := json.Marshal(&mockReview)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))

		response, err := server.Client().Post(
			(fmt.Sprintf(postReviewPath, server.URL, mockReview.Reviewer.ID)),
			applicationJSON, reader)
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

func TestReviewHandlerGetReviewsBy(t *testing.T) {
	mockUseCase := new(mocks.ReviewUseCase)
	h := http.NewReviewHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	r := app.Server(nil, nil, nil, h, mw)

	var mockReviews []domain.Review
	err := faker.FakeData(&mockReviews)

	assert.NoError(t, err)

	var mockStudent domain.Student
	err = faker.FakeData(&mockStudent)
	assert.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		mockUseCase.On("GetReviewsBy", mock.Anything, mock.AnythingOfType("string")).
			Return(mockReviews, nil).
			Once()

		reqFound := httptest.NewRequest("GET", fmt.Sprintf("/api/reviews-by/%s", mockReviews[0].Reviewer.ID), nil)
		reqFound.Header.Set("id", mockReviews[0].Reviewer.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code) // update not post yet it return a status 201
		assert.NoError(t, err)

		responseBody, err := ioutil.ReadAll(w.Body)

		if err != nil {
			assert.Fail(t, failureMessage)
		}

		var receivedReviews []domain.Review
		err = json.Unmarshal(responseBody, &receivedReviews)

		assert.NoError(t, err)
		assert.EqualValues(t, len(mockReviews), len(receivedReviews))
		mockUseCase.AssertExpectations(t)
	})

}
func TestReviewHandlerEditReview(t *testing.T) {
	mockUseCase := new(mocks.ReviewUseCase)
	h := http.NewReviewHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	r := app.Server(nil, nil, nil, h, mw)

	var mockReview domain.Review
	err := faker.FakeData(&mockReview)
	assert.NoError(t, err)
	var mockStudent domain.Student
	err = faker.FakeData(&mockStudent)
	assert.NoError(t, err)


	t.Run("success", func(t *testing.T) {
		postBody, err := json.Marshal(mockReview)
		assert.NoError(t, err)
		mockUseCase.On("EditReview", mock.Anything, mock.AnythingOfType(reviewType), mock.AnythingOfType("string")).
			Return(&mockReview, nil).
			Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putReviewPath, mockReview.Reviewer.ID),
			reader)
		reqFound.Header.Set("id", mockReview.Reviewer.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code) // update not post yet it return a status 201
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid-data-type", func(t *testing.T) {
		reader := strings.NewReader("invalid body")
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putReviewPath, mockReview.Reviewer.ID),
			reader)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-rest-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockReview)
		restErr := e.NewConflictError("error occurred")
		assert.NoError(t, err)
		mockUseCase.On("EditReview", mock.Anything, mock.AnythingOfType(reviewType), mock.AnythingOfType("string")).
			Return(nil, restErr).
			Once()
		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putReviewPath, mockReview.Reviewer.ID),
			reader)
		reqFound.Header.Set("id", mockReview.Reviewer.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("usecase-default-error", func(t *testing.T) {
		postBody, err := json.Marshal(mockReview)
		assert.NoError(t, err)
		mockUseCase.On("EditReview", mock.Anything, mock.AnythingOfType("*domain.Review"), mock.AnythingOfType("string")).
			Return(nil, errors.New("some error")).
			Once()

		reader := strings.NewReader(string(postBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putReviewPath, mockReview.Reviewer.ID),
			reader)
		reqFound.Header.Set("id", mockReview.Reviewer.ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 500, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
