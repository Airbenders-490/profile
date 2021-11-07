package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/airbenders/profile/Tag/usecase"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/domain/mocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTags(t *testing.T) {
	mockTagRepo := new(mocks.TagRepositoryMock)

	var mockTags []domain.Tag
	err := faker.FakeData(&mockTags)
	assert.NoError(t, err)
	t.Run("case success", func(t *testing.T) {
		mockTagRepo.
			On("FetchAllTags", mock.Anything).
			Return(mockTags, nil).
			Once()
		u := usecase.NewTagUseCase(mockTagRepo, time.Second)

		Tag, err := u.GetAllTags(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, Tag)

		mockTagRepo.AssertExpectations(t)
	})

	t.Run("case error", func(t *testing.T) {
		mockTagRepo.
			On("FetchAllTags", mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		u := usecase.NewTagUseCase(mockTagRepo, time.Second)

		Tag, err := u.GetAllTags(context.TODO())

		assert.Error(t, err)
		assert.True(t, reflect.ValueOf(Tag).IsNil())

		mockTagRepo.AssertExpectations(t)
	})

}
