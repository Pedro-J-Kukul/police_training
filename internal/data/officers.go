// FileName: internal/data/officers.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
	"github.com/lib/pq"
)

/************************************************************************************************************/
// Officer Declarations
/************************************************************************************************************/

// Officer struct to represent an officer in the system
type Officer struct {
	ID               int64     `json:"id"`
	RegulationNumber string    `json:"regulation_number"`
	PostingID        int64     `json:"posting_id"`
	RankID           int64     `json:"rank_id"`
	UserID           int64     `json:"user_id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	UserId           int64     `json:"user,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// OfficerModel struct to interact with the officers table in the database
type OfficerModel struct {
	DB *sql.DB
}

/************************************************************************************************************/
// Validation helpers
/************************************************************************************************************/

// ValidateOfficer checks whether the officer struct contains valid data.
func ValidateOfficer(v *validator.Validator, officer *Officer) {
	v.Check(officer.ID > 0, "user_id", "must reference an existing user")
	v.Check(officer.RegulationNumber != "", "regulation_number", "must be provided")
	v.Check(len(officer.RegulationNumber) <= 50, "regulation_number", "must not exceed 50 characters")
	v.Check(officer.PostingID > 0, "posting_id", "must be provided")
	v.Check(officer.RankID > 0, "rank_id", "must be provided")
	v.Check(officer.FormationID > 0, "formation_id", "must be provided")
	v.Check(officer.RegionID > 0, "region_id", "must be provided")
	v.Check(officer.UserId > 0, "user_id", "must reference an existing user")

}

/************************************************************************************************************/
// CRUD methods
/************************************************************************************************************/

// Insert creates a new officer record.
func (m OfficerModel) Insert(officer *Officer) error {
	query := `
		INSERT INTO officers (regulation_number, posting_id, rank_id, formation_id, region_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`

	now := time.Now()
	args := []any{
		officer.RegulationNumber,
		officer.PostingID,
		officer.RankID,
		officer.FormationID,
		officer.RegionID,
		officer.UserId,
		now,
		now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&officer.CreatedAt, &officer.UpdatedAt); err != nil {
		switch {
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves an officer by id.
func (m OfficerModel) Get(id int64) (*Officer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, regulation_number, posting_id, rank_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE id = $1`

	var officer Officer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&officer.ID,
		&officer.RegulationNumber,
		&officer.PostingID,
		&officer.RankID,
		&officer.FormationID,
		&officer.RegionID,
		&officer.CreatedAt,
		&officer.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &officer, nil
}

// GetAll returns a filtered list of officers with pagination metadata.
func (m OfficerModel) GetAll(regulation string, postingID, rankID, formationID, regionID *int64, filters Filters) ([]*Officer, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "regulation_number"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, regulation_number, posting_id, rank_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE ($1 = '' OR regulation_number ILIKE $1)
		AND ($2 = 0 OR posting_id = $2)
		AND ($3 = 0 OR rank_id = $3)
		AND ($4 = 0 OR formation_id = $4)
		AND ($5 = 0 OR region_id = $5)
		ORDER BY %s %s, id ASC
		LIMIT $6 OFFSET $7`, filters.sortColumn(), filters.sortDirection())

	postingArg := int64(0)
	if postingID != nil {
		postingArg = *postingID
	}
	rankArg := int64(0)
	if rankID != nil {
		rankArg = *rankID
	}
	formationArg := int64(0)
	if formationID != nil {
		formationArg = *formationID
	}
	regionArg := int64(0)
	if regionID != nil {
		regionArg = *regionID
	}

	args := []any{
		regulation,
		postingArg,
		rankArg,
		formationArg,
		regionArg,
		filters.limit(),
		filters.offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		officers     []*Officer
		totalRecords int
	)

	for rows.Next() {
		var officer Officer
		if err := rows.Scan(
			&totalRecords,
			&officer.ID,
			&officer.RegulationNumber,
			&officer.PostingID,
			&officer.RankID,
			&officer.FormationID,
			&officer.RegionID,
			&officer.CreatedAt,
			&officer.UpdatedAt,
		); err != nil {
			return nil, MetaData{}, err
		}
		officers = append(officers, &officer)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return officers, metadata, nil
}

// Update modifies an officer record using optimistic locking on updated_at.
func (m OfficerModel) Update(officer *Officer, originalUpdatedAt time.Time) error {
	query := `
		UPDATE officers
		SET regulation_number = $1,
			posting_id = $2,
			rank_id = $3,
			formation_id = $4,
			region_id = $5,
			updated_at = now()
		WHERE id = $6 AND updated_at = $7
		RETURNING updated_at`

	args := []any{
		officer.RegulationNumber,
		officer.PostingID,
		officer.RankID,
		officer.FormationID,
		officer.RegionID,
		officer.ID,
		originalUpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&officer.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

func (m OfficerModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM officers
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// GetByUserID retrieves an officer by user_id (more convenient than officer_id)
func (m OfficerModel) GetByUserID(userID int64) (*Officer, error) {
	if userID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, regulation_number, posting_id, rank_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE user_id = $1`

	var officer Officer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&officer.ID,
		&officer.UserID,
		&officer.RegulationNumber,
		&officer.PostingID,
		&officer.RankID,
		&officer.FormationID,
		&officer.RegionID,
		&officer.CreatedAt,
		&officer.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &officer, nil
}

// GetAllByUserIDs - Enhanced GetAll that works with user_ids
func (m OfficerModel) GetAllByUserIDs(userIDs []int64, regulationNumber string, postingID, rankID, formationID, regionID *int64, filters Filters) ([]*Officer, MetaData, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), o.id, o.user_id, o.regulation_number, o.posting_id, o.rank_id, o.formation_id, o.region_id, o.created_at, o.updated_at
		FROM officers o
		WHERE ($1 = '' OR o.regulation_number ILIKE $1)
		AND ($2 = 0 OR o.posting_id = $2)
		AND ($3 = 0 OR o.rank_id = $3)
		AND ($4 = 0 OR o.formation_id = $4)
		AND ($5 = 0 OR o.region_id = $5)
		AND (CARDINALITY($6::bigint[]) = 0 OR o.user_id = ANY($6::bigint[]))
		ORDER BY %s %s, o.id ASC
		LIMIT $7 OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	// Handle optional parameters
	postingArg := int64(0)
	if postingID != nil {
		postingArg = *postingID
	}
	rankArg := int64(0)
	if rankID != nil {
		rankArg = *rankID
	}
	formationArg := int64(0)
	if formationID != nil {
		formationArg = *formationID
	}
	regionArg := int64(0)
	if regionID != nil {
		regionArg = *regionID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query,
		regulationNumber, postingArg, rankArg, formationArg, regionArg,
		pq.Array(userIDs), filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		officers     []*Officer
		totalRecords int
	)

	for rows.Next() {
		var officer Officer
		if err := rows.Scan(
			&totalRecords,
			&officer.ID,
			&officer.UserID,
			&officer.RegulationNumber,
			&officer.PostingID,
			&officer.RankID,
			&officer.FormationID,
			&officer.RegionID,
			&officer.CreatedAt,
			&officer.UpdatedAt,
		); err != nil {
			return nil, MetaData{}, err
		}
		officers = append(officers, &officer)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return officers, metadata, nil
}
