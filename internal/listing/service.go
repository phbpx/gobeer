// Package listing provides a use case for listing beers and reviews.
package listing

import (
	"context"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviews"
)

// Repository defines the interface for the listing service to interact
// with the storage.
type Repository interface {
	// ListBeers returns a list of beers.
	ListBeers(ctx context.Context) ([]beers.Beer, error)
	// ListReviews returns a list of reviews.
	ListReviews(ctx context.Context, id string) ([]reviews.Review, error)
}

// Service provides beer listing operations.
type Service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies.
func NewService(r Repository) *Service {
	return &Service{r}
}

// ListBeers lists all the beers.
func (s *Service) ListBeers(ctx context.Context) ([]beers.Beer, error) {
	return s.r.ListBeers(ctx)
}

// ListReviews lists all the reviews for a given beer.
func (s *Service) ListReviews(ctx context.Context, id string) ([]reviews.Review, error) {
	// Validate the beer ID.
	if _, err := uuid.Parse(id); err != nil {
		return nil, beers.ErrInvalidID
	}

	return s.r.ListReviews(ctx, id)
}
