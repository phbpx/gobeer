// Package reviewing provices a use case for beer reviewing.
package reviewing

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviews"
)

// NewReview defines the input parameters for creating a new review.
type NewReview struct {
	BeerID  string `json:"beer_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
	Score   int    `json:"score" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

// Repository defines the interface for the reviewing service to interact
// with the storage.
type Repository interface {
	// GetBeer returns the beer with the given ID.
	GetBeer(ctx context.Context, id string) (*beers.Beer, error)
	// CreateReview creates a new review.
	CreateReview(ctx context.Context, review reviews.Review) error
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
func (s *Service) CreateReview(ctx context.Context, nr NewReview) (*reviews.Review, error) {
	// Find the beer.
	if _, err := s.r.GetBeer(ctx, nr.BeerID); err != nil {
		return nil, err
	}

	r := reviews.Review{
		ID:        uuid.NewString(),
		BeerID:    nr.BeerID,
		UserID:    nr.UserID,
		Score:     nr.Score,
		Comment:   nr.Comment,
		CreatedAt: time.Now(),
	}

	// Create the review.
	return &r, s.r.CreateReview(ctx, r)
}
