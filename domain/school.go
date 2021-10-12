package domain

import "context"

type School struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Country string   `json:"country"`
	Domains []string `json:"domains"`
}

type SchoolUseCase interface {
	SearchSchoolByDomain(ctx context.Context, domainName string) ([]School, error)
}

type SchoolRepository interface {
	SearchByDomain(ctx context.Context, name string) ([]School, error)
}
