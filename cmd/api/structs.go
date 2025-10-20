package main

import "time"

// CreateRegionRequest represents the request payload for creating a region
var CreateRegionRequest struct {
	Region string `json:"region"`
}

// UpdateRegionRequest represents the request payload for updating a region
var UpdateRegionRequest struct {
	Region *string `json:"region"`
}

// CreateFormationRequest represents the request payload for creating a formation
var CreateFormationRequest struct {
	Formation string `json:"formation"`
	RegionID  int64  `json:"region_id"`
}

// UpdateFormationRequest represents the request payload for updating a formation
var UpdateFormationRequest struct {
	Formation *string `json:"formation"`
	RegionID  *int64  `json:"region_id"`
}

// CreatePostingRequest represents the request payload for creating a posting
var CreatePostingRequest struct {
	Posting string  `json:"posting"`
	Code    *string `json:"code"`
}

// UpdatePostingRequest represents the request payload for updating a posting
var UpdatePostingRequest struct {
	Posting *string `json:"posting"`
	Code    *string `json:"code"`
}

// CreateRankRequest represents the request payload for creating a rank
var CreateRankRequest struct {
	Rank                        string `json:"rank"`
	Code                        string `json:"code"`
	AnnualTrainingHoursRequired int    `json:"annual_training_hours_required"`
}

// UpdateRankRequest represents the request payload for updating a rank
var UpdateRankRequest struct {
	Rank                        *string `json:"rank"`
	Code                        *string `json:"code"`
	AnnualTrainingHoursRequired *int    `json:"annual_training_hours_required"`
}

// CreateTrainingTypeRequest represents the request payload for creating a training type
var CreateTrainingTypeRequest struct {
	Type string `json:"type"`
}

// UpdateTrainingTypeRequest represents the request payload for updating a training type
var UpdateTrainingTypeRequest struct {
	Type *string `json:"type"`
}

// CreateTrainingCategoryRequest represents the request payload for creating a training category
var CreateTrainingCategoryRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

// UpdateTrainingCategoryRequest represents the request payload for updating a training category
var UpdateTrainingCategoryRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

// CreateWorkshopRequest represents the request payload for creating a workshop
var CreateWorkshopRequest struct {
	WorkshopName   string  `json:"workshop_name"`
	CategoryID     int64   `json:"category_id"`
	TrainingTypeID int64   `json:"training_type_id"`
	CreditHours    int64   `json:"credit_hours"`
	Description    *string `json:"description"`
	Objectives     *string `json:"objectives"`
	IsActive       *bool   `json:"is_active"`
}

// UpdateWorkshopRequest represents the request payload for updating a workshop
var UpdateWorkshopRequest struct {
	WorkshopName   *string    `json:"workshop_name"`
	CategoryID     *int64     `json:"category_id"`
	TrainingTypeID *int64     `json:"training_type_id"`
	CreditHours    *int64     `json:"credit_hours"`
	Description    *string    `json:"description"`
	Objectives     *string    `json:"objectives"`
	IsActive       *bool      `json:"is_active"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// CreateTrainingStatusRequest represents the request payload for creating a training status
var CreateTrainingStatusRequest struct {
	Status string `json:"status"`
}

// UpdateTrainingStatusRequest represents the request payload for updating a training status
var UpdateTrainingStatusRequest struct {
	Status *string `json:"status"`
}

// CreateEnrollmentStatusRequest represents the request payload for creating an enrollment status
var CreateEnrollmentStatusRequest struct {
	Status string `json:"status"`
}

// UpdateEnrollmentStatusRequest represents the request payload for updating an enrollment status
var UpdateEnrollmentStatusRequest struct {
	Status *string `json:"status"`
}

// CreateOfficerRequest represents the request payload for creating an officer
var CreateOfficerRequest struct {
	UserID           int64  `json:"user_id"`
	RegulationNumber string `json:"regulation_number"`
	PostingID        int64  `json:"posting_id"`
	RankID           int64  `json:"rank_id"`
	FormationID      int64  `json:"formation_id"`
	RegionID         int64  `json:"region_id"`
}

// UpdateOfficerRequest represents the request payload for updating an officer
var UpdateOfficerRequest struct {
	RegulationNumber *string    `json:"regulation_number"`
	PostingID        *int64     `json:"posting_id"`
	RankID           *int64     `json:"rank_id"`
	FormationID      *int64     `json:"formation_id"`
	RegionID         *int64     `json:"region_id"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

// CreateTrainingSessionRequest represents the request payload for creating a training session
var CreateTrainingSessionRequest struct {
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	FacilitatorID    int64     `json:"facilitator_id"`
	WorkshopID       int64     `json:"workshop_id"`
	SessionDate      time.Time `json:"session_date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Location         *string   `json:"location"`
	MaxCapacity      *int64    `json:"max_capacity"`
	TrainingStatusID int64     `json:"training_status_id"`
	Notes            *string   `json:"notes"`
}

// UpdateTrainingSessionRequest represents the request payload for updating a training session
var UpdateTrainingSessionRequest struct {
	FormationID      *int64     `json:"formation_id"`
	RegionID         *int64     `json:"region_id"`
	FacilitatorID    *int64     `json:"facilitator_id"`
	WorkshopID       *int64     `json:"workshop_id"`
	SessionDate      *time.Time `json:"session_date"`
	StartTime        *time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time"`
	Location         *string    `json:"location"`
	MaxCapacity      *int64     `json:"max_capacity"`
	TrainingStatusID *int64     `json:"training_status_id"`
	Notes            *string    `json:"notes"`
}

// CreateTrainingEnrollmentRequest represents the request payload for creating a training enrollment
var CreateTrainingEnrollmentRequest struct {
	OfficerID          int64   `json:"officer_id"`
	SessionID          int64   `json:"session_id"`
	EnrollmentStatusID int64   `json:"enrollment_status_id"`
	AttendanceStatusID *int64  `json:"attendance_status_id,omitempty"`
	ProgressStatusID   int64   `json:"progress_status_id"`
	CompletionDate     *string `json:"completion_date,omitempty"`
	CertificateIssued  bool    `json:"certificate_issued"`
	CertificateNumber  *string `json:"certificate_number,omitempty"`
}

var UpdateTrainingEnrollmentRequest struct {
	OfficerID          *int64     `json:"officer_id,omitempty"`
	SessionID          *int64     `json:"session_id,omitempty"`
	EnrollmentStatusID *int64     `json:"enrollment_status_id,omitempty"`
	AttendanceStatusID *int64     `json:"attendance_status_id,omitempty"` // Fixed: was **int64
	ProgressStatusID   *int64     `json:"progress_status_id,omitempty"`
	CompletionDate     *time.Time `json:"completion_date,omitempty"`
	CertificateIssued  *bool      `json:"certificate_issued,omitempty"`
	CertificateNumber  *string    `json:"certificate_number,omitempty"`
}

// Attendance Status Request Variables
var CreateAttendanceStatusRequest struct {
	Status          string `json:"status"`
	CountsAsPresent bool   `json:"counts_as_present"`
}

var UpdateAttendanceStatusRequest struct {
	Status          *string `json:"status,omitempty"`
	CountsAsPresent *bool   `json:"counts_as_present,omitempty"`
}

// Progress Status Request Variables
var CreateProgressStatusRequest struct {
	Status string `json:"status"`
}

var UpdateProgressStatusRequest struct {
	Status *string `json:"status,omitempty"`
}
