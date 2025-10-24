// FileName: internal/data/workshops.go
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
// Workshop Declarations
/************************************************************************************************************/

// Workshop struct to represent a workshop in the system
type Workshop struct {
	ID           int64     `json:"id"`
	WorkshopName string    `json:"workshop_name"`
	CategoryID   int64     `json:"category_id"`
	TypeID       int64     `json:"type_id"`
	CreditHours  int       `json:"credit_hours"`
	Description  *string   `json:"description,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WorkshopModel struct to interact with the workshops table in the database
type WorkshopModel struct {
	DB *sql.DB
}

// ValidateWorkshop ensures workshop data is valid.
func ValidateWorkshop(v *validator.Validator, workshop *Workshop) {
	v.Check(workshop.WorkshopName != "", "workshop_name", "must be provided")
	v.Check(len(workshop.WorkshopName) <= 200, "workshop_name", "must not exceed 200 characters")
	v.Check(workshop.CategoryID > 0, "category_id", "must be provided")
	v.Check(workshop.TypeID > 0, "type_id", "must be provided")
	v.Check(workshop.CreditHours >= 0, "credit_hours", "must be zero or greater")
}

// Insert creates a new workshop.
func (m *WorkshopModel) Insert(workshop *Workshop) error {
	query := `
		INSERT INTO workshops (workshop_name, category_id, type_id, credit_hours, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		workshop.WorkshopName,
		workshop.CategoryID,
		workshop.TypeID,
		workshop.CreditHours,
		workshop.Description,
		workshop.IsActive,
	).Scan(&workshop.ID, &workshop.CreatedAt, &workshop.UpdatedAt); err != nil {
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

// Get retrieves a workshop by id.
func (m *WorkshopModel) Get(id int64) (*Workshop, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, workshop_name, category_id, type_id, credit_hours, description, is_active, created_at, updated_at
		FROM workshops
		WHERE id = $1`

	var workshop Workshop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&workshop.ID,
		&workshop.WorkshopName,
		&workshop.CategoryID,
		&workshop.TypeID,
		&workshop.CreditHours,
		&workshop.Description,
		&workshop.IsActive,
		&workshop.CreatedAt,
		&workshop.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &workshop, nil
}

// GetByName retrieves a workshop by its name
func (m *WorkshopModel) GetByName(workshopName string) (*Workshop, error) {
	query := `
		SELECT id, workshop_name, category_id, type_id, credit_hours, description, is_active, created_at, updated_at
		FROM workshops
		WHERE workshop_name = $1`

	var workshop Workshop
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, workshopName).Scan(
		&workshop.ID,
		&workshop.WorkshopName,
		&workshop.CategoryID,
		&workshop.TypeID,
		&workshop.CreditHours,
		&workshop.Description,
		&workshop.IsActive,
		&workshop.CreatedAt,
		&workshop.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &workshop, nil
}

// GetAll returns workshops filtered by name, category, type, or active state.
func (m *WorkshopModel) GetAll(name string, categoryID, typeID *int64, isActive *bool, filters Filters) ([]*Workshop, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "workshop_name"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, workshop_name, category_id, type_id, credit_hours, description, is_active, created_at, updated_at
		FROM workshops
		WHERE ($1 = '' OR workshop_name ILIKE $1)
		AND ($2 = 0 OR category_id = $2)
		AND ($3 = 0 OR type_id = $3)
		AND ($4::boolean IS NULL OR is_active = $4::boolean)
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	categoryArg := int64(0)
	if categoryID != nil {
		categoryArg = *categoryID
	}

	typeArg := int64(0)
	if typeID != nil {
		typeArg = *typeID
	}

	var activeArg interface{} = nil
	if isActive != nil {
		activeArg = *isActive
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, categoryArg, typeArg, activeArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		workshops    []*Workshop
		totalRecords int
	)

	for rows.Next() {
		var workshop Workshop
		if err := rows.Scan(
			&totalRecords,
			&workshop.ID,
			&workshop.WorkshopName,
			&workshop.CategoryID,
			&workshop.TypeID,
			&workshop.CreditHours,
			&workshop.Description,
			&workshop.IsActive,
			&workshop.CreatedAt,
			&workshop.UpdatedAt,
		); err != nil {
			return nil, MetaData{}, err
		}
		workshops = append(workshops, &workshop)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return workshops, metadata, nil
}

// Update modifies an existing workshop.
func (m *WorkshopModel) Update(workshop *Workshop) error {
	query := `
		UPDATE workshops
		SET workshop_name = $1, category_id = $2, type_id = $3, credit_hours = $4, description = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		workshop.WorkshopName,
		workshop.CategoryID,
		workshop.TypeID,
		workshop.CreditHours,
		workshop.Description,
		workshop.IsActive,
		workshop.ID,
	).Scan(&workshop.UpdatedAt); err != nil {
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

// Delete removes a workshop from the database
func (m *WorkshopModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM workshops WHERE id = $1`

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
