// FilenName: internal/data/models.go

package data

import "database/sql"

// Wrapper for models

type Models struct {
	User                UserModel
	Token               TokenModel
	Permission          PermissionModel
	Role                RoleModel
	RolePermission      RolePermissionModel
	RoleUser            RoleUserModel
	Region              RegionModel
	Formation           FormationModel
	Posting             PostingModel
	Rank                RankModel
	Officer             OfficerModel
	TrainingType        TrainingTypeModel
	TrainingCategory    TrainingCategoryModel
	Workshop            WorkshopModel
	TrainingStatus      TrainingStatusModel
	TrainingSession     TrainingSessionModel
	EnrollmentStatus    EnrollmentStatusModel
	AttendanceStatus    AttendanceStatusModel
	ProgressStatus      ProgressStatusModel
	TrainingEnrollment  TrainingEnrollmentModel
}

// NewModels returns a Models struct containing the initialized models.
func NewModels(db *sql.DB) Models {
	return Models{
		User:               UserModel{DB: db},
		Token:              TokenModel{DB: db},
		Permission:         PermissionModel{DB: db},
		Role:               RoleModel{DB: db},
		RolePermission:     RolePermissionModel{DB: db},
		RoleUser:           RoleUserModel{DB: db},
		Region:             RegionModel{DB: db},
		Formation:          FormationModel{DB: db},
		Posting:            PostingModel{DB: db},
		Rank:               RankModel{DB: db},
		Officer:            OfficerModel{DB: db},
		TrainingType:       TrainingTypeModel{DB: db},
		TrainingCategory:   TrainingCategoryModel{DB: db},
		Workshop:           WorkshopModel{DB: db},
		TrainingStatus:     TrainingStatusModel{DB: db},
		TrainingSession:    TrainingSessionModel{DB: db},
		EnrollmentStatus:   EnrollmentStatusModel{DB: db},
		AttendanceStatus:   AttendanceStatusModel{DB: db},
		ProgressStatus:     ProgressStatusModel{DB: db},
		TrainingEnrollment: TrainingEnrollmentModel{DB: db},
	}
}
