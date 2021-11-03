package domain

import (
	"context"
	"time"
)

// Review struct
type Review struct {
	ID        string  `json:"id"`
	Reviewed  Student `json:"reviewed"`
	Reviewer  Student `json:"reviewer"`
	CreatedAt time.Time
	Tags      []Tag `json:"tags"`
}

// ReviewUseCase is the contract every use case must employ
type ReviewUseCase interface {
	AddReview(ctx context.Context, review *Review, reviewerID string) (*Review, error)
	EditReview(ctx context.Context, review *Review, reviewerID string) (*Review, error)
}

// ReviewRepository is the contract every review repository must employ
type ReviewRepository interface {
	GetReviewByAndFor(ctx context.Context, reviewer string, reviewed string) (*Review, error)
	AddReview(ctx context.Context, review *Review) error
	UpdateReviewTags(ctx context.Context, review *Review) error
}
