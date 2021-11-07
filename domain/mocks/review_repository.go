package mocks

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

type ReviewRepositoryMock struct {
	mock.Mock
}

/*
type ReviewRepository interface {
	GetReviewsFor(ctx context.Context, reviewed string) ([]Review, error)
	GetReviewsBy(ctx context.Context, reviewer string) ([]Review, error)
	GetReviewByAndFor(ctx context.Context, reviewer string, reviewed string) (*Review, error)
	AddReview(ctx context.Context, review *Review) error
	UpdateReviewTags(ctx context.Context, review *Review) error
}
 */

func (m *ReviewRepositoryMock) GetReviewsFor(ctx context.Context, reviewed string) ([]domain.Review, error) {
	args := m.Called(ctx, reviewed)

	var r0 []domain.Review
	if rf, ok := args.Get(0).(func(context.Context, string) []domain.Review); ok {
		r0 = rf(ctx, reviewed)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).([]domain.Review)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, reviewed)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (r *ReviewRepositoryMock) GetReviewsBy(ctx context.Context, reviewer string) ([]domain.Review, error) {
	panic("implement me")
}

func (r *ReviewRepositoryMock) GetReviewByAndFor(ctx context.Context, reviewer string, reviewed string) (*domain.Review, error) {
	panic("implement me")
}

func (r *ReviewRepositoryMock) AddReview(ctx context.Context, review *domain.Review) error {
	panic("implement me")
}

func (r *ReviewRepositoryMock) UpdateReviewTags(ctx context.Context, review *domain.Review) error {
	panic("implement me")
}
