// FileName: internal/data/ranks.go
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
// Rank Declarations
/************************************************************************************************************/

// Rank struct to represent a rank in the system
type Rank struct {
	ID                  int64  `json:"id"`
	Rank                string `json:"rank"`
	Code                string `json:"code"`
	AnnualTrainingHours int    `json:"annual_training_hours"`
}

// RankModel struct to interact with the ranks table in the database
type RankModel struct {
	DB *sql.DB
}

// ValidateRank ensures rank data is valid.
func ValidateRank(v *validator.Validator, rank *Rank) {
	v.Check(rank.Rank != "", "rank", "must be provided")
	v.Check(len(rank.Rank) <= 150, "rank", "must not exceed 150 characters")
	v.Check(rank.Code != "", "code", "must be provided")
	v.Check(len(rank.Code) <= 20, "code", "must not exceed 20 characters")
	v.Check(rank.AnnualTrainingHours >= 0, "annual_training_hours", "must be zero or greater")
}

// Insert adds a new rank record.
func (m *RankModel) Insert(rank *Rank) error {
	query := `
		INSERT INTO ranks (rank, code, annual_training_hours)
		VALUES ($1, $2, $3)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, rank.Rank, rank.Code, rank.AnnualTrainingHours).Scan(&rank.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a rank by id.
func (m *RankModel) Get(id int64) (*Rank, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, rank, code, annual_training_hours FROM ranks WHERE id = $1`

	var rank Rank

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&rank.ID, &rank.Rank, &rank.Code, &rank.AnnualTrainingHours)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rank, nil
}

// GetByName retrieves a rank by name.
func (m *RankModel) GetByName(name string) (*Rank, error) {
	query := `SELECT id, rank, code, annual_training_hours FROM ranks WHERE rank = $1`

	var rank Rank

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, name).Scan(&rank.ID, &rank.Rank, &rank.Code, &rank.AnnualTrainingHours)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rank, nil
}

// GetAll returns ranks filtered by rank or code.
func (m *RankModel) GetAll(name string, code string, filters Filters) ([]*Rank, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "rank"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, rank, code, annual_training_hours
		FROM ranks
		WHERE (to_tsvector('simple', rank) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', code) @@ plainto_tsquery('simple', $2) OR $2 = '')
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
		ranks        []*Rank
		totalRecords int
	)

	for rows.Next() {
		var rank Rank
		if err := rows.Scan(&totalRecords, &rank.ID, &rank.Rank, &rank.Code, &rank.AnnualTrainingHours); err != nil {
			return nil, MetaData{}, err
		}
		ranks = append(ranks, &rank)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return ranks, metadata, nil
}

// Update modifies an existing rank.
func (m RankModel) Update(rank *Rank) error {
	query := `
		UPDATE ranks
		SET rank = $1, code = $2, annual_training_hours = $3
		WHERE id = $4
		RETURNING rank, code, annual_training_hours`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, rank.Rank, rank.Code, rank.AnnualTrainingHours, rank.ID).Scan(&rank.Rank, &rank.Code, &rank.AnnualTrainingHours); err != nil {
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

// Delete removes a rank by id.
func (m *RankModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM ranks WHERE id = $1`

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
