// FileName: internal/data/formations.go
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
// Formation Declarations
/************************************************************************************************************/

// Formation struct to represent a formation in the system
type Formation struct {
	ID        int64     `json:"id"`
	Formation string    `json:"formation"`
	RegionID  int64     `json:"region_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FormationModel struct to interact with the formations table in the database
type FormationModel struct {
	DB *sql.DB
}

// ValidateFormation ensures formation data is valid.
func ValidateFormation(v *validator.Validator, formation *Formation) {
	v.Check(formation.Formation != "", "formation", "must be provided")
	v.Check(len(formation.Formation) <= 150, "formation", "must not exceed 150 characters")
	v.Check(formation.RegionID > 0, "region_id", "must be provided")
}

// Insert creates a new formation.
func (m FormationModel) Insert(formation *Formation) error {
	query := `
		INSERT INTO formations (formation, region_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, formation.Formation, formation.RegionID, now).Scan(&formation.ID, &formation.CreatedAt); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a formation by id.
func (m FormationModel) Get(id int64) (*Formation, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, formation, region_id, created_at FROM formations WHERE id = $1`

	var formation Formation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&formation.ID, &formation.Formation, &formation.RegionID, &formation.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &formation, nil
}

// GetAll returns formations filtered by name and region.
func (m FormationModel) GetAll(name string, regionID *int64, filters Filters) ([]*Formation, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "formation"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, formation, region_id, created_at
		FROM formations
		WHERE ($1 = '' OR formation ILIKE $1)
		AND ($2 = 0 OR region_id = $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	regionArg := int64(0)
	if regionID != nil {
		regionArg = *regionID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, regionArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		formations   []*Formation
		totalRecords int
	)

	for rows.Next() {
		var formation Formation
		if err := rows.Scan(&totalRecords, &formation.ID, &formation.Formation, &formation.RegionID, &formation.CreatedAt); err != nil {
			return nil, MetaData{}, err
		}
		formations = append(formations, &formation)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return formations, metadata, nil
}

// Update modifies an existing formation.
func (m FormationModel) Update(formation *Formation) error {
	query := `
		UPDATE formations
		SET formation = $1, region_id = $2
		WHERE id = $3
		RETURNING formation, region_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, formation.Formation, formation.RegionID, formation.ID).Scan(&formation.Formation, &formation.RegionID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}
