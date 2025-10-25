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
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://police:police@localhost/police_training?sslmode=disable"
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

	constableRank, err := rankModel.GetByName("Sergeant")
	if err != nil {
		log.Fatal("Failed to get constable rank:", err)
	}

	corporalRank, err := rankModel.GetByName("Corporal")
	if err != nil {
		log.Fatal("Failed to get corporal rank:", err)
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
		{"david.wilson@police-training.bz", "PC004", "Constable", "Relief"},
		{"lisa.garcia@police-training.bz", "CPL002", "Corporal", "Station Manager"},
	}

	fmt.Println("Creating officer records...")

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
		if od.RankName == "Corporal" {
			rank = corporalRank
		} else {
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
