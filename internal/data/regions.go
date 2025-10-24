// FileName: internal/data/regions.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/************************************************************************************************************/
// Region Declarations
/************************************************************************************************************/

// Region struct to represent a region in the system
type Region struct {
	ID     int64  `json:"id"`
	Region string `json:"region"`
}

// RegionModel struct to interact with the regions table in the database
type RegionModel struct {
	DB *sql.DB
}

// ValidateRegion ensures region data is valid.
func ValidateRegion(v *validator.Validator, region *Region) {
	v.Check(region.Region != "", "region", "must be provided")
	v.Check(len(region.Region) <= 100, "region", "must not exceed 100 characters")
}

// Insert adds a new region.
func (m *RegionModel) Insert(region *Region) error {
	query := `
		INSERT INTO regions (region)
		VALUES ($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, region.Region).Scan(&region.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a region by id.
func (m *RegionModel) Get(id int64) (*Region, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, region FROM regions WHERE id = $1`

	var region Region

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&region.ID, &region.Region)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &region, nil
}

// GetByName retrieves a region by name.
func (m *RegionModel) GetByName(name string) (*Region, error) {
	query := `SELECT id, region FROM regions WHERE region = $1`

	var region Region

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, name).Scan(&region.ID, &region.Region)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &region, nil
}

// GetAll returns regions filtered by name with pagination.
func (m *RegionModel) GetAll(name string, filters Filters) ([]*Region, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "region"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, region
		FROM regions
		WHERE (to_tsvector('simple', region) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		regions      []*Region
		totalRecords int
	)

	for rows.Next() {
		var region Region
		if err := rows.Scan(&totalRecords, &region.ID, &region.Region); err != nil {
			return nil, MetaData{}, err
		}
		regions = append(regions, &region)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return regions, metadata, nil
}

// Update modifies an existing region.
func (m *RegionModel) Update(region *Region) error {
	query := `
		UPDATE regions
		SET region = $1
		WHERE id = $2
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, region.Region, region.ID).Scan(&region.ID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Delete removes a region by id.
func (m *RegionModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM regions WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
