package domain

import (
	"time"
)

type Student struct {
	ID string `json:"id"` //uuid string
	Name string `json:"name"`
	Email string `json:"email"` //TODO: validate:required
	GeneralInfo string `json:"general_info"`
	School string `json:"school"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Reviews []Review
}




