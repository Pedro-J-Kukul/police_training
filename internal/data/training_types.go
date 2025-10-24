// FileName: internal/data/training_types.go
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
// TrainingType Declarations
/************************************************************************************************************/

// TrainingType struct to represent a training type in the system
type TrainingType struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

// TrainingTypeModel struct to interact with the training_types table in the database
type TrainingTypeModel struct {
	DB *sql.DB
}

// ValidateTrainingType ensures training type data is valid.
func ValidateTrainingType(v *validator.Validator, trainingType *TrainingType) {
	v.Check(trainingType.Type != "", "type", "must be provided")
	v.Check(len(trainingType.Type) <= 150, "type", "must not exceed 150 characters")
}

// Insert adds a new training type.
func (m *TrainingTypeModel) Insert(trainingType *TrainingType) error {
	query := `
		INSERT INTO training_types (name, is_active)
		VALUES ($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, trainingType.Type, trainingType.IsActive).Scan(&trainingType.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a training type by id.
func (m *TrainingTypeModel) Get(id int64) (*TrainingType, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name FROM training_types WHERE id = $1`

	var trainingType TrainingType

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&trainingType.ID, &trainingType.Type)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &trainingType, nil
}

// GetAll returns training types filtered by name.
func (m *TrainingTypeModel) GetAll(name string, filters Filters) ([]*TrainingType, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "type"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, is_active
		FROM training_types
		WHERE (to_tsvector('simple', type) @@ plainto_tsquery('simple', $1) OR $1 = '')
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
		trainingTypes []*TrainingType
		totalRecords  int
	)

	for rows.Next() {
		var trainingType TrainingType
		if err := rows.Scan(&totalRecords, &trainingType.ID, &trainingType.Type); err != nil {
			return nil, MetaData{}, err
		}
		trainingTypes = append(trainingTypes, &trainingType)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return trainingTypes, metadata, nil
}

// Update modifies an existing training type.
func (m *TrainingTypeModel) Update(trainingType *TrainingType) error {
	query := `
		UPDATE training_types
		SET type = $1, is_active = $2
		WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, trainingType.Type, trainingType.IsActive, trainingType.ID).Scan(&trainingType.Type, &trainingType.IsActive); err != nil {
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

// Delete removes a training type by id.
func (m *TrainingTypeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM training_types WHERE id = $1`

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

// GetByType retrieves a training type by its type name.
func (m *TrainingTypeModel) GetByType(typeName string) (*TrainingType, error) {
	if typeName == "" {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name, is_active FROM training_types WHERE name = $1`

	var trainingType TrainingType

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, typeName).Scan(&trainingType.ID, &trainingType.Type, &trainingType.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &trainingType, nil
}
