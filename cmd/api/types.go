package main

import "time"

// CreateRegionRequest_T represents the request payload for creating a region
type CreateRegionRequest_T struct {
	Region string `json:"region"`
}

// UpdateRegionRequest_T represents the request payload for updating a region
type UpdateRegionRequest_T struct {
	Region *string `json:"region"`
}

// CreateFormationRequest_T represents the request payload for creating a formation
type CreateFormationRequest_T struct {
	Formation string `json:"formation"`
	RegionID  int64  `json:"region_id"`
}

// UpdateFormationRequest_T represents the request payload for updating a formation
type UpdateFormationRequest_T struct {
	Formation *string `json:"formation"`
	RegionID  *int64  `json:"region_id"`
}

// CreatePostingRequest_T represents the request payload for creating a posting
type CreatePostingRequest_T struct {
	Posting string  `json:"posting"`
	Code    *string `json:"code"`
}

// UpdatePostingRequest_T represents the request payload for updating a posting
type UpdatePostingRequest_T struct {
	Posting *string `json:"posting"`
	Code    *string `json:"code"`
}

// CreateRankRequest_T represents the request payload for creating a rank
type CreateRankRequest_T struct {
	Rank                        string `json:"rank"`
	Code                        string `json:"code"`
	AnnualTrainingHoursRequired int    `json:"annual_training_hours_required"`
}

// UpdateRankRequest_T represents the request payload for updating a rank
type UpdateRankRequest_T struct {
	Rank                        *string `json:"rank"`
	Code                        *string `json:"code"`
	AnnualTrainingHoursRequired *int    `json:"annual_training_hours_required"`
}

// CreateTrainingTypeRequest_T represents the request payload for creating a training type
type CreateTrainingTypeRequest_T struct {
	Type string `json:"type"`
}

// UpdateTrainingTypeRequest_T represents the request payload for updating a training type
type UpdateTrainingTypeRequest_T struct {
	Type *string `json:"type"`
}

// CreateTrainingCategoryRequest_T represents the request payload for creating a training category
type CreateTrainingCategoryRequest_T struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

// UpdateTrainingCategoryRequest_T represents the request payload for updating a training category
type UpdateTrainingCategoryRequest_T struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

// CreateWorkshopRequest_T represents the request payload for creating a workshop
type CreateWorkshopRequest_T struct {
	WorkshopName   string  `json:"workshop_name"`
	CategoryID     int64   `json:"category_id"`
	TrainingTypeID int64   `json:"training_type_id"`
	CreditHours    int64   `json:"credit_hours"`
	Description    *string `json:"description"`
	Objectives     *string `json:"objectives"`
	IsActive       *bool   `json:"is_active"`
}

// UpdateWorkshopRequest_T represents the request payload for updating a workshop
type UpdateWorkshopRequest_T struct {
	WorkshopName   *string    `json:"workshop_name"`
	CategoryID     *int64     `json:"category_id"`
	TrainingTypeID *int64     `json:"training_type_id"`
	CreditHours    *int64     `json:"credit_hours"`
	Description    **string   `json:"description"`
	Objectives     **string   `json:"objectives"`
	IsActive       *bool      `json:"is_active"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// CreateTrainingStatusRequest_T represents the request payload for creating a training status
type CreateTrainingStatusRequest_T struct {
	Status string `json:"status"`
}

// UpdateTrainingStatusRequest_T represents the request payload for updating a training status
type UpdateTrainingStatusRequest_T struct {
	Status *string `json:"status"`
}

// CreateEnrollmentStatusRequest_T represents the request payload for creating an enrollment status
type CreateEnrollmentStatusRequest_T struct {
	Status string `json:"status"`
}

// UpdateEnrollmentStatusRequest_T represents the request payload for updating an enrollment status
type UpdateEnrollmentStatusRequest_T struct {
	Status *string `json:"status"`
}

// CreateOfficerRequest_T represents the request payload for creating an officer
type CreateOfficerRequest_T struct {
	UserID           int64  `json:"user_id"`
	RegulationNumber string `json:"regulation_number"`
	PostingID        int64  `json:"posting_id"`
	RankID           int64  `json:"rank_id"`
	FormationID      int64  `json:"formation_id"`
	RegionID         int64  `json:"region_id"`
}

// UpdateOfficerRequest_T represents the request payload for updating an officer
type UpdateOfficerRequest_T struct {
	RegulationNumber *string    `json:"regulation_number"`
	PostingID        *int64     `json:"posting_id"`
	RankID           *int64     `json:"rank_id"`
	FormationID      *int64     `json:"formation_id"`
	RegionID         *int64     `json:"region_id"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

// registerUserRequest represents the request payload for user registration
type registerUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Password    string `json:"password"`
	Facilitator bool   `json:"facilitator"`
}

// activateUserRequest represents the request payload for user activation
type activateUserRequest struct {
	Token string `json:"token"`
}

// CreatePasswordResetTokenRequest_T represents the request payload for creating a password reset token
type CreatePasswordResetTokenRequest_T struct {
	Email string `json:"email"`
}

// ResetPasswordRequest_T represents the request payload for resetting a user's password
type ResetPasswordRequest_T struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type CreateAuthenticationTokenRequest_T struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateTrainingSessionRequest_T represents the request payload for creating a training session
type CreateTrainingSessionRequest_T struct {
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

// UpdateTrainingSessionRequest_T represents the request payload for updating a training session
type UpdateTrainingSessionRequest_T struct {
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

type CreateTrainingEnrollmentRequest_T struct {
	OfficerID          int64      `json:"officer_id"`
	SessionID          int64      `json:"session_id"`
	EnrollmentStatusID int64      `json:"enrollment_status_id"`
	AttendanceStatusID *int64     `json:"attendance_status_id,omitempty"`
	ProgressStatusID   int64      `json:"progress_status_id"`
	CompletionDate     *time.Time `json:"completion_date,omitempty"`
	CertificateIssued  bool       `json:"certificate_issued"`
	CertificateNumber  *string    `json:"certificate_number,omitempty"`
}

type UpdateTrainingEnrollmentRequest_T struct {
	OfficerID          *int64      `json:"officer_id"`
	SessionID          *int64      `json:"session_id"`
	EnrollmentStatusID *int64      `json:"enrollment_status_id"`
	AttendanceStatusID **int64     `json:"attendance_status_id,omitempty"`
	ProgressStatusID   *int64      `json:"progress_status_id"`
	CompletionDate     **time.Time `json:"completion_date,omitempty"`
	CertificateIssued  *bool       `json:"certificate_issued"`
	CertificateNumber  **string    `json:"certificate_number,omitempty"`
}
