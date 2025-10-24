package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	_ "github.com/lib/pq"
)

func main() {
	// Get database connection string from environment
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "postgres://police:police@localhost/police_training_testing?sslmode=disable"
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
		// 1 Admin
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

	fmt.Println("Sample users creation completed!")
}
