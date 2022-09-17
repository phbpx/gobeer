package reviewing_test

import (
	"context"
	"testing"

	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviewing"
)

// mockRepository is a mock implementation of the Repository interface.
type mockRepository struct {
	data []beers.Beer
}

// FindBeer returns the beer with the given ID.
func (r *mockRepository) FindBeer(ctx context.Context, id string) (*beers.Beer, error) {
	for _, b := range r.data {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, beers.ErrNotFound
}

// CreateReview creates a new review.
func (r *mockRepository) CreateReview(ctx context.Context, nr reviewing.NewReview) error {
	return nil
}

func TestCreateReview(t *testing.T) {
	// Create a mock repository.
	r := &mockRepository{
		data: []beers.Beer{
			{ID: "1", Name: "Beer 1"},
			{ID: "2", Name: "Beer 2"},
		},
	}

	// Create a new service with the mock repository.
	s := reviewing.NewService(r)

	t.Logf("Given the need to test creating a new review.")
	{
		t.Logf("\tWhen creating a new review for a beer that exists.")
		{
			nr := reviewing.NewReview{
				BeerID:  "1",
				UserID:  "1",
				Score:   5,
				Comment: "A very nice beer",
			}
			if err := s.CreateReview(context.Background(), nr); err != nil {
				t.Fatalf("\t\t[ERROR] Should be able to create the review. Error: %v", err)
			}
			t.Logf("\t\t[OK] Should be able to create the review.")
		}

		t.Logf("\tWhen creating a new review for a beer that does not exist.")
		{
			nr := reviewing.NewReview{
				BeerID:  "3",
				UserID:  "1",
				Score:   5,
				Comment: "A very nice beer",
			}
			if err := s.CreateReview(context.Background(), nr); err != beers.ErrNotFound {
				t.Fatalf("\t\t[ERROR] Should not be able to create the review. Error: %v", err)
			}
			t.Logf("\t\t[OK] Should not be able to create the review.")
		}
	}
}
