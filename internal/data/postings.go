// FileName: internal/data/postings.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// Posting Declarations
/************************************************************************************************************/

// Posting struct to represent a posting in the system
type Posting struct {
	ID        int64     `json:"id"`
	Posting   string    `json:"posting"`
	Code      *string   `json:"code,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// PostingModel struct to interact with the postings table in the database
type PostingModel struct {
	DB *sql.DB
}
