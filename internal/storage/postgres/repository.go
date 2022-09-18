package postgres

import (
	"context"
	"database/sql"

	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/reviews"
)

// Storage provides an implementation if the Repositories interface.
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new Storage instance.
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// CreateBeer creates a new beer on the database.
func (s *Storage) CreateBeer(ctx context.Context, b beers.Beer) error {
	query := `
        INSERT INTO beers (
                id, 
                name, 
                brewery, 
                style, 
                abv, 
                short_desc, 
                created_at
        ) VALUES (
                $1, $2, $3, $4, $5, $6, $7
        )`

	_, err := s.db.ExecContext(ctx, query,
		b.ID,
		b.Name,
		b.Brewery,
		b.Style,
		b.ABV,
		b.ShortDesc,
		b.CreatedAt)

	return err
}

// BeerExists checks if a beer exists on the database.
func (s *Storage) BeerExists(ctx context.Context, name, brewery string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM beers WHERE name = $1 AND brewery = $2)`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, name, brewery).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetBeer returns a beer from the database.
func (s *Storage) GetBeer(ctx context.Context, id string) (*beers.Beer, error) {
	query := `
        SELECT 
                b.id,
                b.name,
                b.brewery,
                b.style,
                b.abv,
                b.short_desc,
                COALESCE(AVG(r.score), 0) AS score,
                b.created_at
        FROM 
                beers AS b
        LEFT JOIN 
                reviews AS r ON r.beer_id = b.id
        WHERE 
                b.id = $1
        GROUP BY
                b.id`

	var b beers.Beer
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.Name,
		&b.Brewery,
		&b.Style,
		&b.ABV,
		&b.ShortDesc,
		&b.Score,
		&b.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, beers.ErrNotFound
		}
		return nil, err
	}

	return &b, nil
}

// ListBeers returns a list of beers from the database.
func (s *Storage) ListBeers(ctx context.Context) ([]beers.Beer, error) {
	query := `
        SELECT 
                b.id,
                b.name,
                b.brewery,
                b.style,
                b.abv,
                b.short_desc,
                COALESCE(AVG(r.score), 0) AS score,
                b.created_at
        FROM 
                beers AS b
        LEFT JOIN 
                reviews AS r ON r.beer_id = b.id
        GROUP BY
                b.id`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []beers.Beer
	for rows.Next() {
		var b beers.Beer

		err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.Brewery,
			&b.Style,
			&b.ABV,
			&b.ShortDesc,
			&b.Score,
			&b.CreatedAt)

		if err != nil {
			return nil, err
		}

		list = append(list, b)
	}

	return list, nil
}

func (s *Storage) CreateReview(ctx context.Context, r reviews.Review) error {
	query := `
        INSERT INTO reviews (
                id,
                beer_id,
                user_id,
                score,
                comment,
                created_at
        ) VALUES (
                $1, $2, $3, $4, $5, $6
        )`

	_, err := s.db.ExecContext(ctx, query,
		r.ID,
		r.BeerID,
		r.UserID,
		r.Score,
		r.Comment,
		r.CreatedAt)

	return err

}

// ListReviews returns a list of reviews from the database.
func (s *Storage) ListReviews(ctx context.Context, id string) ([]reviews.Review, error) {
	query := `
        SELECT 
                r.id,
                r.beer_id,
                r.user_id,
                r.score,
                r.comment,
                r.created_at
        FROM 
                reviews AS r
        WHERE 
                r.beer_id = $1
        ORDER BY 
                r.created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []reviews.Review
	for rows.Next() {
		var r reviews.Review

		err := rows.Scan(
			&r.ID,
			&r.BeerID,
			&r.UserID,
			&r.Score,
			&r.Comment,
			&r.CreatedAt)

		if err != nil {
			return nil, err
		}

		list = append(list, r)
	}

	return list, nil
}
