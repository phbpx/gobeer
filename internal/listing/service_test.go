package listing_test

import (
	"context"
	"testing"

	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/listing"
	"github.com/phbpx/gobeer/internal/reviews"
)

// mockRepository is a mock implementation of the Repository interface.
type mockRepository struct {
	beers   []beers.Beer
	reviews []reviews.Review
}

// ListBeers returns a list of beers.
func (r *mockRepository) ListBeers(ctx context.Context) ([]beers.Beer, error) {
	return r.beers, nil
}

// ListReviews returns a list of reviews.
func (r *mockRepository) ListReviews(ctx context.Context, id string) ([]reviews.Review, error) {
	for _, review := range r.reviews {
		if review.BeerID == id {
			return []reviews.Review{review}, nil
		}
	}
	return []reviews.Review{}, nil
}

func TestListing(t *testing.T) {
	// Create a mock repository.
	r := &mockRepository{
		beers: []beers.Beer{
			{ID: "1", Name: "Beer 1", Brewery: "Brewery 1"},
			{ID: "2", Name: "Beer 2", Brewery: "Brewery 2"},
		},
		reviews: []reviews.Review{
			{ID: "1", BeerID: "1", UserID: "1", Score: 5, Comment: "Comment 1"},
			{ID: "2", BeerID: "2", UserID: "2", Score: 4, Comment: "Comment 2"},
		},
	}

	// Create a listing service with the mock repository.
	service := listing.NewService(r)

	t.Log("Given the need to list beers.")
	{
		t.Log("\tWhen handling the list beers request.")
		{
			// List the beers.
			_, err := service.ListBeers(context.Background())
			if err != nil {
				t.Fatalf("\t\t[ERROR] Should be able to list the beers. Error: %s", err)
			}
			t.Log("\t\t[OK] Should be able to list the beers.")
		}
	}

	t.Log("Given the need to list reviews.")
	{
		t.Log("\tWhen handling the list reviews request.")
		{
			// List the reviews.
			_, err := service.ListReviews(context.Background(), "1")
			if err != nil {
				t.Fatalf("\t\t[ERROR] Should be able to list the reviews. Error: %s", err)
			}
			t.Log("\t\t[OK] Should be able to list the reviews.")
		}
	}
}
