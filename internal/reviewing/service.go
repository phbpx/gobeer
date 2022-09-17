// Package reviewing provices a use case for beer reviewing.
package reviewing

import (
	"context"

	"github.com/phbpx/gobeer/internal/beers"
)

// NewReview defines the input parameters for creating a new review.
type NewReview struct {
	BeerID  string `json:"beer_id"`
	UserID  string `json:"user_id"`
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

type Repository interface {
	// FindBeer returns the beer with the given ID.
	FindBeer(ctx context.Context, id string) (*beers.Beer, error)
	// CreateReview creates a new review.
	CreateReview(ctx context.Context, review NewReview) error
}

// Service provides beer reviewing operations.
type Service struct {
	r Repository
}

// NewService creates a reviewing service with the necessary dependencies.
func NewService(r Repository) *Service {
	return &Service{r}
}

// CreateReview creates a new review.
func (s *Service) CreateReview(ctx context.Context, review NewReview) error {
	// Find the beer.
	if _, err := s.r.FindBeer(ctx, review.BeerID); err != nil {
		return err
	}

	// Create the review.
	return s.r.CreateReview(ctx, review)
}
