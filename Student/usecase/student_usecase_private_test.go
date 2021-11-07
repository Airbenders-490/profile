package usecase

import (
	"github.com/airbenders/profile/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateStudent(t *testing.T) {
	now := time.Now()
	existing := &domain.Student{
		ID:          "asd",
		FirstName:        "Sunny",
		LastName:        "Moony",
		Email:       "none@gmail.com",
		GeneralInfo: "I like plants",
		School:      "KGS",
		CreatedAt:   now,
		UpdatedAt:   now.Add(72 * time.Hour),
	}
	toUpdate := &domain.Student{
		ID:          "",
		Email:       "something@gmail.com",
		GeneralInfo: "",
		School:      "Concordia University",
		CreatedAt:   now,
		UpdatedAt:   now.Add(72 * time.Hour),
	}
	expected := &domain.Student{
		ID:          "asd",
		FirstName:        "Sunny",
		LastName:        "Moony",
		Email:       "something@gmail.com",
		GeneralInfo: "I like plants",
		School:      "Concordia University",
		CreatedAt:   now,
		UpdatedAt:   time.Now(),
	}

	updateStudent(existing, toUpdate)

	existing.UpdatedAt = expected.UpdatedAt
	assert.EqualValues(t, expected, existing)
}
