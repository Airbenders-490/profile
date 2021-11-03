package domain

import (
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

type ReviewUseCase interface {
	AddReview(review *Review, reviewerID int) (*Review, error)
	EditReview(review *Review, reviewerID int) (*Review, error)
	DeleteReview(review *Review, reviewerID int) error
}
