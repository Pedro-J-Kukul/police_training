// FileName: internal/data/postings.go
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
// Posting Declarations
/************************************************************************************************************/

// Posting struct to represent a posting in the system
type Posting struct {
	ID      int64  `json:"id"`
	Posting string `json:"posting"`
	Code    string `json:"code,omitempty"`
}

// *PostingModel struct to interact with the postings table in the database
type PostingModel struct {
	DB *sql.DB
}

// ValidatePosting ensures posting data is valid.
func ValidatePosting(v *validator.Validator, posting *Posting) {
	v.Check(posting.Posting != "", "posting", "must be provided")
	v.Check(len(posting.Posting) <= 150, "posting", "must not exceed 150 characters")
	v.Check(len(posting.Code) <= 20, "code", "must not exceed 20 characters")
}

// Insert creates a new posting.
func (m *PostingModel) Insert(posting *Posting) error {
	query := `
		INSERT INTO postings (posting, code)
		VALUES ($1, $2)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, posting.Posting, posting.Code).Scan(&posting.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a posting by id.
func (m *PostingModel) Get(id int64) (*Posting, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, posting, code FROM postings WHERE id = $1`

	var posting Posting

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&posting.ID, &posting.Posting, &posting.Code)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &posting, nil
}

// GetAll returns postings filtered by name or code.
func (m *PostingModel) GetAll(name string, code string, filters Filters) ([]*Posting, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "posting"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, posting, code
		FROM postings
		WHERE ($1 = '' OR posting ILIKE $1)
		AND ($2 = '' OR code ILIKE $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, code, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		postings     []*Posting
		totalRecords int
	)

	for rows.Next() {
		var posting Posting
		if err := rows.Scan(&totalRecords, &posting.ID, &posting.Posting, &posting.Code); err != nil {
			return nil, MetaData{}, err
		}
		postings = append(postings, &posting)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return postings, metadata, nil
}

// Update modifies an existing posting.
func (m *PostingModel) Update(posting *Posting) error {
	query := `
		UPDATE postings
		SET posting = $1, code = $2
		WHERE id = $3
		RETURNING posting, code`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, posting.Posting, posting.Code, posting.ID).Scan(&posting.Posting, &posting.Code); err != nil {
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
