package domain

import "time"

type Confirmation struct {
	Token string `json:"token"`
	School School `json:"in_school"`
	Student Student `json:"for_student"`
	CreatedAt time.Time `json:"created_at"`
}