package domain

import "context"

// Tag struct
type Tag struct {
	Name     string `json:"name"`
	Positive bool   `json:"positive"`
}

// TagUseCase only returns all available tags for reviewing
type TagUseCase interface {
	GetAllTags(ctx context.Context) ([]Tag, error)
}

// TagRepository only returns all available tags for reviewing from the db
type TagRepository interface {
	FetchAllTags(ctx context.Context) ([]Tag, error)
}
