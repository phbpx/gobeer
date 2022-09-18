// Package reviews defines the review domain model.
package reviews

import "time"

// Review defines the properties of a review.
type Review struct {
	ID        string    `json:"id"`
	BeerID    string    `json:"beer_id"`
	UserID    string    `json:"user_id"`
	Score     int       `json:"score"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
