package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/driftprogramming/pgxpoolmock"
)

// tagRepository struct
type tagRepository struct {
	db pgxpoolmock.PgxPool
}

// NewTagRepository is a constructor for TagRepository
func NewTagRepository(db pgxpoolmock.PgxPool) domain.TagRepository {
	return &tagRepository{db: db}
}

const (
	fetchAll = "SELECT name, positive FROM tag"
)

// FetchAllTags returns all the tags
func (u *tagRepository) FetchAllTags(ctx context.Context) ([]domain.Tag, error) {
	rows, err := u.db.Query(ctx, fetchAll)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err = rows.Scan(&tag.Name, &tag.Positive)
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
