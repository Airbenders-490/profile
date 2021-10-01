package domain

import "time"

type Review struct {
	ID string `json:"id"`
	reviewed Student `json:"reviewed"`
	reviewer Student `json:"reviewer"`
	CreatedAt time.Time
	tags []Tag `json:"tags"`
}
