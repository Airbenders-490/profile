package domain

import "time"

// Review struct
type Review struct {
	ID        string  `json:"id"`
	Reviewed  Student `json:"reviewed"`
	Reviewer  Student `json:"reviewer"`
	CreatedAt time.Time
	Tags      []Tag `json:"tags"`
}
