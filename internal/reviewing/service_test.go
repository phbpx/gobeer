package reviewing_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviewing"
	"github.com/phbpx/gobeer/internal/reviews"
)

// mockRepository is a mock implementation of the Repository interface.
type mockRepository struct {
	data []beers.Beer
}

// GetBeer returns the beer with the given ID.
func (r *mockRepository) GetBeer(ctx context.Context, id string) (*beers.Beer, error) {
	for _, b := range r.data {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, beers.ErrNotFound
}

// CreateReview creates a new review.
func (r *mockRepository) CreateReview(ctx context.Context, nr reviews.Review) error {
	return nil
}

func TestCreateReview(t *testing.T) {
	ctx := context.Background()

	beerID := uuid.NewString()

	// Create a mock repository.
	r := &mockRepository{
		data: []beers.Beer{
			{ID: beerID, Name: "Beer 1"},
			{ID: uuid.NewString(), Name: "Beer 2"},
		},
	}

	// Create a new service with the mock repository.
	s := reviewing.NewService(r)

	t.Logf("Given the need to test creating a new review.")
	{
		t.Logf("\tWhen creating a new review for a beer that exists.")
		{
			nr := reviewing.NewReview{
				BeerID:  beerID,
				UserID:  uuid.NewString(),
				Score:   5,
				Comment: "A very nice beer",
			}
			if _, err := s.CreateReview(ctx, nr); err != nil {
				t.Fatalf("\t\t[ERROR] Should be able to create the review. Error: %v", err)
			}
			t.Logf("\t\t[OK] Should be able to create the review.")
		}

		t.Logf("\tWhen creating a new review for a beer that does not exist.")
		{
			nr := reviewing.NewReview{
				BeerID:  uuid.NewString(),
				UserID:  uuid.NewString(),
				Score:   5,
				Comment: "A very nice beer",
			}
			if _, err := s.CreateReview(ctx, nr); err != beers.ErrNotFound {
				t.Fatalf("\t\t[ERROR] Should not be able to create the review. Error: %v", err)
			}
			t.Logf("\t\t[OK] Should not be able to create the review.")
		}
	}
}
