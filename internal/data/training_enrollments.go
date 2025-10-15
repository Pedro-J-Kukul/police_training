// FileName: internal/data/training_enrollments.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// TrainingEnrollment Declarations
/************************************************************************************************************/

// TrainingEnrollment struct to represent a training enrollment in the system
type TrainingEnrollment struct {
	ID                 int64      `json:"id"`
	OfficerID          int64      `json:"officer_id"`
	SessionID          int64      `json:"session_id"`
	EnrollmentStatusID int64      `json:"enrollment_status_id"`
	AttendanceStatusID *int64     `json:"attendance_status_id,omitempty"`
	ProgressStatusID   int64      `json:"progress_status_id"`
	CompletionDate     *time.Time `json:"completion_date,omitempty"`
	CertificateIssued  bool       `json:"certificate_issued"`
	CertificateNumber  *string    `json:"certificate_number,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// TrainingEnrollmentModel struct to interact with the training_enrollments table in the database
type TrainingEnrollmentModel struct {
	DB *sql.DB
}
