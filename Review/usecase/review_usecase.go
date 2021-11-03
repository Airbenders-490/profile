package usecase

import (
	"context"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

// reviewUseCase struct implements ReviewUseCase interface
type reviewUseCase struct {
	rr      domain.ReviewRepository
	sr      domain.StudentRepository
	timeout time.Duration
}

// NewReviewUseCase is the constructor
func NewReviewUseCase(rr domain.ReviewRepository, sr domain.StudentRepository, timeout time.Duration) domain.ReviewUseCase {
	return &reviewUseCase{
		rr:      rr,
		sr:      sr,
		timeout: timeout,
	}
}

// AddReview first checks if the person being reviewed exists.
func (u *reviewUseCase) AddReview(c context.Context, review *domain.Review, reviewerID string) (*domain.Review, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	student, err := u.sr.GetByID(ctx, review.Reviewed.ID)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(student, &domain.Student{}) {
		return nil, errors.NewBadRequestError("the person being reviewed doesn't exist")
	}

	anyExistingReview, err := u.rr.GetReviewByAndFor(ctx, reviewerID, review.Reviewed.ID)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(anyExistingReview, &domain.Review{}) && anyExistingReview != nil {
		return nil, errors.NewBadRequestError("the review already exists. Please update instead.")
	}

	review.Reviewer.ID = reviewerID
	review.ID = uuid.NewString()
	review.CreatedAt = time.Now()

	err = u.rr.AddReview(ctx, review)
	if err != nil {
		return nil, err
	}

	review.Reviewed = *student
	return review, nil
}

func (u *reviewUseCase) EditReview(c context.Context, review *domain.Review, reviewerID string) (*domain.Review, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	if review.ID == "" {
		return nil, errors.NewBadRequestError("review ID must be provided")
	}

	student, err := u.sr.GetByID(ctx, review.Reviewed.ID)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(student, &domain.Student{}) {
		return nil, errors.NewBadRequestError("the person being reviewed doesn't exist")
	}

	anyExistingReview, err := u.rr.GetReviewByAndFor(ctx, reviewerID, review.Reviewed.ID)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(anyExistingReview, &domain.Review{}) || anyExistingReview == nil {
		return nil, errors.NewBadRequestError("the review doesn't exists. Please create instead.")
	}

	review.Reviewer.ID = reviewerID
	review.CreatedAt = time.Now()

	err = u.rr.UpdateReviewTags(ctx, review)
	if err != nil {
		return nil, err
	}

	review.Reviewed = *student
	return review, nil
}
