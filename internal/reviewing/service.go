// Package reviewing provices a use case for beer reviewing.
package reviewing

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviews"
)

// NewReview defines the input parameters for creating a new review.
type NewReview struct {
	UserID  string  `json:"user_id" binding:"required,uuid"`
	Score   float32 `json:"score" binding:"required"`
	Comment string  `json:"comment" binding:"required"`
}

// Storer defines the interface for the reviewing service to interact
// with the storage.
type Storer interface {
	// GetBeer returns the beer with the given ID.
	GetBeer(ctx context.Context, id string) (*beers.Beer, error)
	// CreateReview creates a new review.
	CreateReview(ctx context.Context, review reviews.Review) error
}

// Notifier defines the interface for the reviewing service to notify users.
type Notifier interface {
	Notify(ctx context.Context, userID string) error
}

// Service provides beer reviewing operations.
type Service struct {
	storer   Storer
	notifier Notifier
}

// NewService creates a reviewing service with the necessary dependencies.
func NewService(storer Storer, notifier Notifier) *Service {
	return &Service{
		storer:   storer,
		notifier: notifier,
	}
}

// CreateReview creates a new review.
func (s *Service) CreateReview(ctx context.Context, beerID string, nr NewReview) (reviews.Review, error) {
	if _, err := uuid.Parse(beerID); err != nil {
		return reviews.Review{}, beers.ErrInvalidID
	}

	if _, err := s.storer.GetBeer(ctx, beerID); err != nil {
		return reviews.Review{}, fmt.Errorf("get beer[id=%s]: %w", beerID, err)
	}

	r := reviews.Review{
		ID:        uuid.NewString(),
		BeerID:    beerID,
		UserID:    nr.UserID,
		Score:     nr.Score,
		Comment:   nr.Comment,
		CreatedAt: time.Now(),
	}

	if err := s.storer.CreateReview(ctx, r); err != nil {
		return reviews.Review{}, fmt.Errorf("create beer[id=%s] review: %w", beerID, err)
	}

	if err := s.notifier.Notify(ctx, r.UserID); err != nil {
		return reviews.Review{}, fmt.Errorf("notify user[id=%s]: %w", r.UserID, err)
	}

	return r, nil
}
