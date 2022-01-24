package mocks

import (
	"context"

	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/mock"
)

//ReviewUseCase Mock struct
type ReviewUseCase struct {
	mock.Mock
}

// AddReview - ReviewUseCase
func (m *ReviewUseCase) AddReview(ctx context.Context, review *domain.Review, reviewerID string) (*domain.Review, error) {
	args := m.Called(ctx, review, reviewerID)

	var r0 *domain.Review
	if rf, ok := args.Get(0).(func(context.Context, *domain.Review, string) *domain.Review); ok {
		r0 = rf(ctx, review, reviewerID)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*domain.Review)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, *domain.Review, string) error); ok {
		r1 = rf(ctx, review, reviewerID)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1

}

// EditReview - ReviewUseCase
func (m *ReviewUseCase) EditReview(ctx context.Context, review *domain.Review, reviewerID string) (*domain.Review, error) {
	args := m.Called(ctx, review, reviewerID)

	var r0 *domain.Review
	if rf, ok := args.Get(0).(func(context.Context, *domain.Review, string) *domain.Review); ok {
		r0 = rf(ctx, review, reviewerID)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*domain.Review)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, *domain.Review, string) error); ok {
		r1 = rf(ctx, review, reviewerID)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1

}

// GetReviewsBy - ReviewUseCase
func (m *ReviewUseCase) GetReviewsBy(ctx context.Context, reviewer string) ([]domain.Review, error) {
	args := m.Called(ctx, reviewer)

	var r0 []domain.Review
	if rf, ok := args.Get(0).(func(context.Context, string) []domain.Review); ok {
		r0 = rf(ctx, reviewer)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).([]domain.Review)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, reviewer)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1

}
