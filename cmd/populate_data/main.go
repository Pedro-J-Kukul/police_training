package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	_ "github.com/lib/pq"
)

func main() {
	// Define command line flags
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "Database connection string (required)")
	flag.Parse()

	// Check if DSN is provided
	if dsn == "" {
		log.Fatal("Database DSN is required. Use -dsn flag to specify the connection string.")
	}

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Database connection successful!")
	fmt.Println("Starting data population process...")

	// Step 1: Populate Users
	fmt.Println("\n=== Step 1: Populating Users ===")
	populateUsers(db)

	// Step 2: Populate Officers
	fmt.Println("\n=== Step 2: Populating Officers ===")
	populateOfficers(db)

	// Step 3: Populate Training Sessions
	fmt.Println("\n=== Step 3: Populating Training Sessions ===")
	populateSessions(db)

	// Step 4: Populate Training Enrollments
	fmt.Println("\n=== Step 4: Populating Training Enrollments ===")
	populateEnrollments(db)

	fmt.Println("\n=== Data population completed successfully! ===")
}

func populateUsers(db *sql.DB) {
	// Initialize models
	userModel := data.UserModel{DB: db}
	roleModel := data.RoleModel{DB: db}

	// Create users
	users := []struct {
		FirstName     string
		LastName      string
		Email         string
		Gender        string
		IsFacilitator bool
		IsOfficer     bool
		Role          string
	}{
		// 2 Admin
		{"Pedro", "Kukul", "admin1@police-training.bz", "m", true, false, "Admin"},
		{"Immanuel", "Garcia", "admin2.garcia@police-training.bz", "m", true, false, "Admin"},

		// 3 Content Creators/Facilitators
		{"Maria", "Rodriguez", "maria.rodriguez@police-training.bz", "f", true, false, "Content-Contributor"},
		{"Carlos", "Martinez", "carlos.martinez@police-training.bz", "m", true, false, "Content-Contributor"},
		{"Ana", "Lopez", "ana.lopez@police-training.bz", "f", true, false, "Content-Contributor"},

		// 6 Officers
		{"John", "Smith", "john.smith@police-training.bz", "m", false, true, "Officer"},
		{"Sarah", "Johnson", "sarah.johnson@police-training.bz", "f", false, true, "Officer"},
		{"Michael", "Brown", "michael.brown@police-training.bz", "m", false, true, "Officer"},
		{"Jessica", "Davis", "jessica.davis@police-training.bz", "f", false, true, "Officer"},
		{"David", "Wilson", "david.wilson@police-training.bz", "m", false, true, "Officer"},
		{"Lisa", "Garcia", "lisa.garcia@police-training.bz", "f", false, true, "Officer"},
	}

	fmt.Println("Creating sample users...")

	for _, userData := range users {
		// Check if user already exists
		existingUser, err := userModel.GetByEmail(userData.Email)
		if err == nil {
			fmt.Printf("User %s already exists (ID: %d), skipping...\n", userData.Email, existingUser.ID)
			continue
		}

		// Create new user
		user := &data.User{
			FirstName:     userData.FirstName,
			LastName:      userData.LastName,
			Email:         userData.Email,
			Gender:        userData.Gender,
			IsActivated:   true,
			IsFacilitator: userData.IsFacilitator,
			IsOfficer:     userData.IsOfficer,
		}

		// Set password - this automatically hashes it
		err = user.Password.Set("TrainingPass123!")
		if err != nil {
			log.Printf("Failed to set password for %s: %v", userData.Email, err)
			continue
		}

		// Insert user
		err = userModel.Insert(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", userData.Email, err)
			continue
		}

		fmt.Printf("Created user: %s %s (ID: %d)\n", user.FirstName, user.LastName, user.ID)

		// Assign role
		err = roleModel.AssignToUser(user.ID, userData.Role)
		if err != nil {
			log.Printf("Failed to assign role %s to user %s: %v", userData.Role, userData.Email, err)
		} else {
			fmt.Printf("Assigned role '%s' to %s\n", userData.Role, userData.Email)
		}
	}

	fmt.Println("Users creation completed!")
}

func populateOfficers(db *sql.DB) {
	// Initialize models
	userModel := data.UserModel{DB: db}
	officerModel := data.OfficerModel{DB: db}
	regionModel := data.RegionModel{DB: db}
	formationModel := data.FormationModel{DB: db}
	postingModel := data.PostingModel{DB: db}
	rankModel := data.RankModel{DB: db}

	// Get reference data (using your seed data)
	region, err := regionModel.GetByName("Northern Region")
	if err != nil {
		log.Fatal("Failed to get region:", err)
	}

	formation, err := formationModel.Get(1) // Corozal Police Formation
	if err != nil {
		log.Fatal("Failed to get formation:", err)
	}

	posting, err := postingModel.GetByName("Relief")
	if err != nil {
		log.Fatal("Failed to get posting:", err)
	}

	constableRank, err := rankModel.GetByName("Constable")
	if err != nil {
		log.Fatal("Failed to get constable rank:", err)
	}

	corporalRank, err := rankModel.GetByName("Corporal")
	if err != nil {
		log.Fatal("Failed to get corporal rank:", err)
	}

	sergeantRank, err := rankModel.GetByName("Sergeant")
	if err != nil {
		log.Fatal("Failed to get sergeant rank:", err)
	}

	// Define officer data for our users
	officerData := []struct {
		Email            string
		RegulationNumber string
		RankName         string
		PostingName      string
	}{
		{"john.smith@police-training.bz", "PC001", "Constable", "Relief"},
		{"sarah.johnson@police-training.bz", "PC002", "Constable", "Relief"},
		{"michael.brown@police-training.bz", "CPL001", "Corporal", "Station Manager"},
		{"jessica.davis@police-training.bz", "PC003", "Constable", "Relief"},
		{"david.wilson@police-training.bz", "SGT001", "Sergeant", "Station Manager"},
		{"lisa.garcia@police-training.bz", "CPL002", "Corporal", "Relief"},
	}

	fmt.Println("Creating officer renilcords...")

	for _, od := range officerData {
		// Get the user
		user, err := userModel.GetByEmail(od.Email)
		if err != nil {
			log.Printf("Failed to find user %s: %v", od.Email, err)
			continue
		}

		// Check if officer record already exists
		existingOfficer, err := officerModel.GetByUserID(user.ID)
		if err == nil {
			fmt.Printf("Officer record for %s already exists (ID: %d), skipping...\n", od.Email, existingOfficer.ID)
			continue
		}

		// Get rank for this officer
		var rank *data.Rank
		switch od.RankName {
		case "Corporal":
			rank = corporalRank
		case "Sergeant":
			rank = sergeantRank
		default:
			rank = constableRank
		}

		// Get posting for this officer
		var officerPosting *data.Posting
		if od.PostingName == "Station Manager" {
			officerPosting, err = postingModel.GetByName("Station Manager")
			if err != nil {
				log.Printf("Failed to get Station Manager posting: %v", err)
				officerPosting = posting // fallback to Relief
			}
		} else {
			officerPosting = posting
		}

		// Create officer record
		officer := &data.Officer{
			UserID:           user.ID,
			RegulationNumber: od.RegulationNumber,
			RankID:           rank.ID,
			PostingID:        officerPosting.ID,
			FormationID:      formation.ID,
			RegionID:         region.ID,
		}

		err = officerModel.Insert(officer)
		if err != nil {
			log.Printf("Failed to create officer for %s: %v", od.Email, err)
			continue
		}

		fmt.Printf("Created officer: %s %s (Reg: %s, Rank: %s, Officer ID: %d)\n",
			user.FirstName, user.LastName, od.RegulationNumber, od.RankName, officer.ID)
	}

	fmt.Println("Officer records creation completed!")
}

func populateSessions(db *sql.DB) {
	// Initialize models
	userModel := data.UserModel{DB: db}
	sessionModel := data.TrainingSessionModel{DB: db}
	formationModel := data.FormationModel{DB: db}
	regionModel := data.RegionModel{DB: db}
	trainingStatusModel := data.TrainingStatusModel{DB: db}

	// Get facilitators
	facilitators := []string{
		"maria.rodriguez@police-training.bz",
		"carlos.martinez@police-training.bz",
		"ana.lopez@police-training.bz",
	}

	// Get workshops (assuming workshop IDs 1-4 exist from seed data)
	workshops := []int64{1, 2, 3, 4}

	// Get formation and region
	formation, err := formationModel.Get(1)
	if err != nil {
		log.Fatal("Failed to get formation:", err)
	}

	region, err := regionModel.GetByName("Northern Region")
	if err != nil {
		log.Fatal("Failed to get region:", err)
	}

	// Get training status (assuming "Scheduled" status exists)
	trainingStatus, err := trainingStatusModel.GetByName("completed")
	if err != nil {
		log.Fatal("Failed to get training status:", err)
	}

	fmt.Println("Creating training sessions...")

	// Create sessions for next 30 days
	baseDate := time.Now()
	sessionCount := 0

	for i := 0; i < 12; i++ { // Create 12 sessions over next few weeks
		sessionDate := baseDate.AddDate(0, 0, i*2+1) // Every other day

		facilitatorEmail := facilitators[i%len(facilitators)]
		workshopID := workshops[i%len(workshops)]

		// Get facilitator user
		facilitator, err := userModel.GetByEmail(facilitatorEmail)
		if err != nil {
			log.Printf("Failed to find facilitator %s: %v", facilitatorEmail, err)
			continue
		}

		startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 9, 0, 0, 0, sessionDate.Location())
		endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 17, 0, 0, 0, sessionDate.Location())

		location := fmt.Sprintf("Training Room %d", (i%3)+1)
		maxCapacity := 20
		notes := fmt.Sprintf("Training session %d conducted by %s %s", i+1, facilitator.FirstName, facilitator.LastName)

		// Create session
		session := &data.TrainingSession{
			FacilitatorID:    facilitator.ID,
			WorkshopID:       workshopID,
			FormationID:      formation.ID,
			RegionID:         region.ID,
			TrainingStatusID: trainingStatus.ID,
			SessionDate:      sessionDate,
			StartTime:        startTime,
			EndTime:          endTime,
			Location:         &location,
			MaxCapacity:      &maxCapacity,
			Notes:            &notes,
		}

		err = sessionModel.Insert(session)
		if err != nil {
			log.Printf("Failed to create session: %v", err)
			continue
		}

		fmt.Printf("Created session: Workshop %d on %s by %s %s (Session ID: %d)\n",
			workshopID, sessionDate.Format("2006-01-02"), facilitator.FirstName, facilitator.LastName, session.ID)
		sessionCount++
	}

	fmt.Printf("Training sessions creation completed! Created %d sessions.\n", sessionCount)
}

func populateEnrollments(db *sql.DB) {
	// Initialize models
	officerModel := data.OfficerModel{DB: db}
	sessionModel := data.TrainingSessionModel{DB: db}
	enrollmentModel := data.TrainingEnrollmentModel{DB: db}
	enrollmentStatusModel := data.EnrollmentStatusModel{DB: db}
	attendanceStatusModel := data.AttendanceStatusModel{DB: db}
	progressStatusModel := data.ProgressStatusModel{DB: db}

	// Get all officers// CORRECT:
	officers, _, err := officerModel.GetAll("", nil, nil, nil, nil, data.Filters{
		Page:     1,
		PageSize: 1000, // Get all records
	})
	if err != nil {
		log.Fatal("Failed to get officers:", err)
	}

	sessions, _, err := sessionModel.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
		Page:         1,
		PageSize:     1000,
		Sort:         "id",
		SortSafelist: []string{"id", "-id"},
	})

	// Get status references
	enrolledStatus, err := enrollmentStatusModel.GetByName("Enrolled")
	if err != nil {
		log.Fatal("Failed to get enrolled status:", err)
	}

	presentStatus, err := attendanceStatusModel.GetByName("Present")
	if err != nil {
		log.Fatal("Failed to get present status:", err)
	}

	inProgressStatus, err := progressStatusModel.GetByName("In Progress")
	if err != nil {
		log.Fatal("Failed to get in progress status:", err)
	}

	completedStatus, err := progressStatusModel.GetByName("Completed")
	if err != nil {
		log.Fatal("Failed to get completed status:", err)
	}

	fmt.Println("Creating training enrollments...")

	enrollmentCount := 0

	// Enroll officers in sessions (not all officers in all sessions to be realistic)
	for i, session := range sessions {
		// Enroll 3-5 officers per session
		numEnrollments := 3 + (i % 3) // 3, 4, or 5 officers per session

		for j := 0; j < numEnrollments && j < len(officers); j++ {
			officerIndex := (i*2 + j) % len(officers) // Distribute officers across sessions
			officer := officers[officerIndex]

			// Determine if this is a completed or in-progress enrollment
			var progressStatus *data.ProgressStatus
			var completionDate *time.Time
			var certificateIssued bool
			var certificateNumber string

			// If session is in the past, mark as completed
			if session.SessionDate.Before(time.Now()) {
				progressStatus = completedStatus
				completionTime := session.SessionDate.AddDate(0, 0, 1)
				completionDate = &completionTime
				certificateIssued = true
				certificateNumber = fmt.Sprintf("CERT-%d-%d-%d", session.ID, officer.ID, session.SessionDate.Year())
			} else {
				progressStatus = inProgressStatus
			}

			// Create enrollment
			enrollment := &data.TrainingEnrollment{
				OfficerID:          officer.ID,
				SessionID:          session.ID,
				EnrollmentStatusID: enrolledStatus.ID,
				AttendanceStatusID: &presentStatus.ID,
				ProgressStatusID:   progressStatus.ID,
				CompletionDate:     completionDate,
				CertificateIssued:  certificateIssued,
				CertificateNumber:  &certificateNumber,
			}

			err = enrollmentModel.Insert(enrollment)
			if err != nil {
				log.Printf("Failed to create enrollment for officer %d in session %d: %v", officer.ID, session.ID, err)
				continue
			}

			status := "In Progress"
			if certificateIssued {
				status = "Completed"
			}

			fmt.Printf("Enrolled officer %d in session %d (Status: %s, Enrollment ID: %d)\n",
				officer.ID, session.ID, status, enrollment.ID)
			enrollmentCount++
		}
	}

	fmt.Printf("Training enrollments creation completed! Created %d enrollments.\n", enrollmentCount)
}
