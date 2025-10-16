// FileName: internal/data/training_categories.go
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
// TrainingCategory Declarations
/************************************************************************************************************/

// TrainingCategory struct to represent a training category in the system
type TrainingCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TrainingCategoryModel struct to interact with the training_categories table in the database
type TrainingCategoryModel struct {
	DB *sql.DB
}

// ValidateTrainingCategory ensures category data is valid.
func ValidateTrainingCategory(v *validator.Validator, category *TrainingCategory) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 150, "name", "must not exceed 150 characters")
}

// Insert creates a new training category.
func (m TrainingCategoryModel) Insert(category *TrainingCategory) error {
	query := `
		INSERT INTO training_categories (name, is_active, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, category.Name, category.IsActive, now).Scan(&category.ID, &category.CreatedAt); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a training category by id.
func (m TrainingCategoryModel) Get(id int64) (*TrainingCategory, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name, is_active, created_at FROM training_categories WHERE id = $1`

	var category TrainingCategory

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.IsActive, &category.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

// GetAll returns categories filtered by name or activation state.
func (m TrainingCategoryModel) GetAll(name string, isActive *bool, filters Filters) ([]*TrainingCategory, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "name"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, is_active, created_at
		FROM training_categories
		WHERE ($1 = '' OR name ILIKE $1)
		AND ($2 IS NULL OR is_active = $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	var activeArg any = nil
	if isActive != nil {
		activeArg = *isActive
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, activeArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		categories   []*TrainingCategory
		totalRecords int
	)

	for rows.Next() {
		var category TrainingCategory
		if err := rows.Scan(&totalRecords, &category.ID, &category.Name, &category.IsActive, &category.CreatedAt); err != nil {
			return nil, MetaData{}, err
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return categories, metadata, nil
}

// Update modifies an existing training category.
func (m TrainingCategoryModel) Update(category *TrainingCategory) error {
	query := `
		UPDATE training_categories
		SET name = $1, is_active = $2
		WHERE id = $3
		RETURNING name, is_active`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, category.Name, category.IsActive, category.ID).Scan(&category.Name, &category.IsActive); err != nil {
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
