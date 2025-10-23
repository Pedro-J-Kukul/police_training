// FileName: internal/data/officers.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/************************************************************************************************************/
// Officer Declarations
/************************************************************************************************************/

// Officer struct to represent an officer in the system
type Officer struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	RegulationNumber string    `json:"regulation_number"`
	RankID           int64     `json:"rank_id"`
	PostingID        int64     `json:"posting_id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	User      *User      `json:"user,omitempty"`
	Rank      *Rank      `json:"rank,omitempty"`
	Posting   *Posting   `json:"posting,omitempty"`
	Formation *Formation `json:"formation,omitempty"`
	Region    *Region    `json:"region,omitempty"`
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
	v.Check(officer.UserID > 0, "user_id", "must be provided and greater than zero")
	v.Check(officer.RegulationNumber != "", "regulation_number", "must be provided")
	v.Check(len(officer.RegulationNumber) <= 50, "regulation_number", "must not be more than 50 characters long")
	v.Check(officer.RankID > 0, "rank_id", "must be provided and greater than zero")
	v.Check(officer.PostingID > 0, "posting_id", "must be provided and greater than zero")
	v.Check(officer.FormationID > 0, "formation_id", "must be provided and greater than zero")
	v.Check(officer.RegionID > 0, "region_id", "must be provided and greater than zero")
}

/************************************************************************************************************/
// CRUD methods
/************************************************************************************************************/

// Insert creates a new officer record.
func (m *OfficerModel) Insert(officer *Officer) error {
	query := `
		INSERT INTO officers (user_id, regulation_number, rank_id, posting_id, formation_id, region_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	args := []any{
		officer.UserID,
		officer.RegulationNumber,
		officer.RankID,
		officer.PostingID,
		officer.FormationID,
		officer.RegionID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&officer.ID, &officer.CreatedAt, &officer.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "officers_user_id_key"`:
			return ErrDuplicateValue
		case err.Error() == `pq: duplicate key value violates unique constraint "officers_regulation_number_key"`:
			return ErrDuplicateValue
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Get retrieves an officer by id.
func (m *OfficerModel) Get(id int64) (*Officer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, regulation_number, rank_id, posting_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE id = $1`

	var officer Officer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&officer.ID,
		&officer.UserID,
		&officer.RegulationNumber,
		&officer.RankID,
		&officer.PostingID,
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

// GetByUserID retrieves an officer by their user ID
func (m *OfficerModel) GetByUserID(userID int64) (*Officer, error) {
	if userID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, regulation_number, rank_id, posting_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE user_id = $1`

	var officer Officer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&officer.ID,
		&officer.UserID,
		&officer.RegulationNumber,
		&officer.RankID,
		&officer.PostingID,
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

// GetByRegulationNumber retrieves an officer by their regulation number
func (m *OfficerModel) GetByRegulationNumber(regulationNumber string) (*Officer, error) {
	query := `
		SELECT id, user_id, regulation_number, rank_id, posting_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE regulation_number = $1`

	var officer Officer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, regulationNumber).Scan(
		&officer.ID,
		&officer.UserID,
		&officer.RegulationNumber,
		&officer.RankID,
		&officer.PostingID,
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

// Update modifies an existing officer in the database
func (m *OfficerModel) Update(officer *Officer) error {
	query := `
		UPDATE officers
		SET regulation_number = $1, rank_id = $2, posting_id = $3, formation_id = $4, region_id = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	args := []any{
		officer.RegulationNumber,
		officer.RankID,
		officer.PostingID,
		officer.FormationID,
		officer.RegionID,
		officer.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&officer.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "officers_regulation_number_key"`:
			return ErrDuplicateValue
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Delete removes an officer from the database
func (m *OfficerModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM officers WHERE id = $1`

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

// GetAll retrieves all officers from the database with filtering and pagination
func (m *OfficerModel) GetAll(regulationNumber string, rankID, postingID, formationID, regionID *int64, filters Filters) ([]*Officer, MetaData, error) {
	query := `
		SELECT COUNT(*) OVER(), id, user_id, regulation_number, rank_id, posting_id, formation_id, region_id, created_at, updated_at
		FROM officers
		WHERE (to_tsvector('simple', regulation_number) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (rank_id = $2 OR $2 IS NULL)
		AND (posting_id = $3 OR $3 IS NULL)
		AND (formation_id = $4 OR $4 IS NULL)
		AND (region_id = $5 OR $5 IS NULL)
		ORDER BY id ASC
		LIMIT $6 OFFSET $7`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, regulationNumber, rankID, postingID, formationID, regionID, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	totalRecords := 0
	officers := []*Officer{}

	for rows.Next() {
		var officer Officer
		err := rows.Scan(
			&totalRecords,
			&officer.ID,
			&officer.UserID,
			&officer.RegulationNumber,
			&officer.RankID,
			&officer.PostingID,
			&officer.FormationID,
			&officer.RegionID,
			&officer.CreatedAt,
			&officer.UpdatedAt,
		)
		if err != nil {
			return nil, MetaData{}, err
		}
		officers = append(officers, &officer)
	}

	if err = rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return officers, metadata, nil
}

// GetWithDetails retrieves an officer with all related information (user, rank, posting, etc.)
func (m *OfficerModel) GetWithDetails(id int64) (*Officer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
			o.id, o.user_id, o.regulation_number, o.rank_id, o.posting_id, o.formation_id, o.region_id, o.created_at, o.updated_at,
			u.first_name, u.last_name, u.email, u.gender, u.is_activated, u.is_facilitator, u.is_officer,
			r.rank, r.code as rank_code,
			p.posting, p.code as posting_code,
			f.formation,
			rg.region
		FROM officers o
		LEFT JOIN users u ON o.user_id = u.id
		LEFT JOIN ranks r ON o.rank_id = r.id
		LEFT JOIN postings p ON o.posting_id = p.id
		LEFT JOIN formations f ON o.formation_id = f.id
		LEFT JOIN regions rg ON o.region_id = rg.id
		WHERE o.id = $1`

	var officer Officer
	var user User
	var rank Rank
	var posting Posting
	var formation Formation
	var region Region

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&officer.ID,
		&officer.UserID,
		&officer.RegulationNumber,
		&officer.RankID,
		&officer.PostingID,
		&officer.FormationID,
		&officer.RegionID,
		&officer.CreatedAt,
		&officer.UpdatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Gender,
		&user.IsActivated,
		&user.IsFacilitator,
		&user.IsOfficer,
		&rank.Rank,
		&rank.Code,
		&posting.Posting,
		&posting.Code,
		&formation.Formation,
		&region.Region,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Set the related objects
	user.ID = officer.UserID
	rank.ID = officer.RankID
	posting.ID = officer.PostingID
	formation.ID = officer.FormationID
	region.ID = officer.RegionID

	officer.User = &user
	officer.Rank = &rank
	officer.Posting = &posting
	officer.Formation = &formation
	officer.Region = &region

	return &officer, nil
}
