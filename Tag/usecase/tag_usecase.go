package usecase

import (
	"context"
	"github.com/airbenders/profile/domain"
	"time"
)

// tagUseCase struct
type tagUseCase struct {
	r       domain.TagRepository
	timeout time.Duration
}

// NewTagUseCase is a constructor for tagUseCase
func NewTagUseCase(r domain.TagRepository, timeout time.Duration) domain.TagUseCase {
	return &tagUseCase{r: r, timeout: timeout}
}

// GetAllTags returns all the tags
func (u *tagUseCase) GetAllTags(c context.Context) ([]domain.Tag, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	tags, err := u.r.FetchAllTags(ctx)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
