package domain

import "time"

// Confirmation stores the confirmation token for validation. Also stores the student and school id associated with
// the confirmation. Lastly has a timestamp for expiry of records
type Confirmation struct {
	Token     string    `json:"token"`
	School    School    `json:"in_school"`
	Student   Student   `json:"for_student"`
	CreatedAt time.Time `json:"created_at"`
}
