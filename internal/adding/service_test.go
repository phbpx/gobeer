package adding_test

import (
	"context"
	"testing"

	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/beers"
)

// mockRepository is a mock implementation of the Repository interface.
type mockRepository struct {
	data []beers.Beer
}

func (m *mockRepository) CreateBeer(ctx context.Context, b beers.Beer) error {
	m.data = append(m.data, b)
	return nil
}

func (m *mockRepository) BeerExists(ctx context.Context, name, brewery string) (bool, error) {
	for _, b := range m.data {
		if b.Name == name && b.Brewery == brewery {
			return true, nil
		}
	}
	return false, nil
}

func TestAddingBeer(t *testing.T) {
	ctx := context.Background()

	// Create a mock repository.
	repo := &mockRepository{}

	// Create a new service with the mock repository.
	s := adding.NewService(repo)

	// Create a new beer.
	b := adding.NewBeer{
		Name:      "IPA",
		Brewery:   "BrewDog",
		Style:     "IPA",
		ABV:       5.5,
		ShortDesc: "A very nice IPA",
	}

	t.Log("Given the need to add a new beer to the system")
	{
		t.Log("\tWhen adding a new beer")
		{
			err := s.AddBeer(ctx, b)
			if err != nil {
				t.Fatalf("\t\t[ERROR] Should be able to add the beer without error: %v", err)
			}
			t.Log("\t\t[OK] Should be able to add the beer without error.")
		}

		t.Log("\tWhen adding a beer that already exists")
		{
			err := s.AddBeer(ctx, b)
			if err != beers.ErrAlreadyExists {
				t.Fatalf("\t\t[ERROR] Should not be able to add the beer: %v", err)
			}
			t.Log("\t\t[OK] Should not be able to add the beer.")
		}
	}
}
