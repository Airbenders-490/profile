package repository_test

import (
	"context"
	"github.com/airbenders/profile/Tag/repository"
	"github.com/airbenders/profile/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
)

func TestFetchAllTags(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	columns := []string{"name", "positive"}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow("1", true).AddRow("2", false).ToPgxRows()

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any()).Return(pgxRows, nil)
		tr := repository.NewTagRepository(mockPool)
		tags, err := tr.FetchAllTags(context.Background())
		assert.NoError(t, err, "error found when not expected")
		assert.EqualValues(t, []domain.Tag{{"1", true}, {"2", false}}, tags)
	})

	t.Run("success", func(t *testing.T) {
		mockPool.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
		tr := repository.NewTagRepository(mockPool)
		tags, err := tr.FetchAllTags(context.Background())
		assert.Error(t, err, "error not found when expected")
		assert.Nil(t, tags)
	})
}
