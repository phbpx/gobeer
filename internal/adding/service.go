// adding defines the use case for adding a beer.
package adding

import (
	"time"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/beers"
)

// NewBeer represents a new beer to be added to the system.
type NewBeer struct {
	Name      string  `json:"name"`
	Brewery   string  `json:"brewery"`
	Style     string  `json:"style"`
	ABV       float32 `json:"abv"`
	ShortDesc string  `json:"short_desc"`
}

// Repository defines the interface for the adding service to interact
// with the storage.
type Repository interface {
	// CreateBeer adds a new beer to the storage.
	CreateBeer(b beers.Beer) error
	// BeerExists checks if a beer with the given name and brewery already exists.
	BeerExists(name, brewery string) (bool, error)
}

// Service provides adding operations.
type Service struct {
	r Repository
}

// NewService creates an adding service with the necessary dependencies.
func NewService(r Repository) *Service {
	return &Service{r}
}

// AddBeer adds a new beer to the system.
func (s *Service) AddBeer(b NewBeer) error {
	beer := beers.Beer{
		ID:        uuid.NewString(),
		Name:      b.Name,
		Brewery:   b.Brewery,
		Style:     b.Style,
		ABV:       b.ABV,
		ShortDesc: b.ShortDesc,
		Score:     0,
		CreatedAt: time.Now(),
	}

	// Check if the beer already exists.
	exists, err := s.r.BeerExists(beer.Name, beer.Brewery)
	if err != nil {
		return err
	}

	// If the beer already exists, return an error.
	if exists {
		return beers.ErrAlreadyExists
	}

	return s.r.CreateBeer(beer)
}
