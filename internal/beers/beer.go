// Package beers defines the beer domain model.
package beers

import (
	"errors"
	"time"
)

var (
	// ErrInvalidID is returned when an invalid ID is provided.
	ErrInvalidID = errors.New("invalid beer ID")

	// ErrNotFound is used when a beer is not found.
	ErrNotFound = errors.New("beer not found")

	// ErrAlreadyExists is used when a beer already exists.
	ErrAlreadyExists = errors.New("beer already exists")
)

// Beer defines the properties of a beer.
type Beer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Brewery   string    `json:"brewery"`
	Style     string    `json:"style"`
	ABV       float32   `json:"abv"`
	ShortDesc string    `json:"short_desc"`
	Score     float32   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
