package repository

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

// tagRepository struct
type tagRepository struct {
	db *pgxpool.Pool
}

// NewTagRepository is a constructor for TagRepository
func NewTagRepository(db *pgxpool.Pool) domain.TagRepository {
	return &tagRepository{db: db}
}

const (
	fetchAll = `SELECT * FROM tag`
)

// FetchAllTags returns all the tags
func (u *tagRepository) FetchAllTags(ctx context.Context) ([]domain.Tag, error) {
	rows, err := u.db.Query(ctx, fetchAll)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var tags []domain.Tag
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		var tag domain.Tag
		tag.Name = value[0].(string)
		tag.Positive = value[1].(bool)

		tags = append(tags, tag)
	}

	return tags, nil
}
