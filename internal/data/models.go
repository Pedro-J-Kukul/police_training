// FilenName: internal/data/models.go

package data

import "database/sql"

// Wrapper for models

type Models struct {
	User  UserModel
	Token TokenModel
}

// NewModels returns a Models struct containing the initialized models.
func NewModels(db *sql.DB) Models {
	return Models{
		User:  UserModel{DB: db},
		Token: TokenModel{DB: db},
	}
}
