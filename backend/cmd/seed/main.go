package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/igorracki/f1/backend/internal/auth"
	"github.com/igorracki/f1/backend/internal/cache"
	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/config"
	"github.com/igorracki/f1/backend/internal/database"
	"github.com/igorracki/f1/backend/internal/models"
	"github.com/igorracki/f1/backend/internal/services"
)

func main() {
	ctx := context.Background()

	// 1. Initialize Database
	configuration := config.Load()
	manager, err := database.NewManager(configuration.DatabasePath)
	if err != nil {
		log.Fatalf("CRITICAL: Failed to connect to database: %v", err)
	}
	defer manager.Close()
	db := manager.DB()

	// 2. Initialize F1 Service to fetch schedules/drivers
	f1Client := clients.NewF1DataClient(configuration.ExternalAPIURL)
	f1Cache := cache.NewMemoryCache()
	f1Service := services.NewF1Service(f1Client, f1Cache)

	// 3. Create Users
	users := []struct {
		Email       string
		DisplayName string
		Password    string
	}{
		{"test1@test.com", "test_one", "test1234"},
		{"test2@test.com", "test_two", "test1234"},
	}

	userIDs := make(map[string]string)

	for _, u := range users {
		passwordHash, err := auth.HashPassword(u.Password)
		if err != nil {
			log.Fatalf("ERROR: Failed to hash password for %s: %v", u.Email, err)
		}

		userID := uuid.New().String()
		now := time.Now().UTC()

		// Insert User
		_, err = db.Exec("INSERT OR IGNORE INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
			userID, u.Email, passwordHash, now)
		if err != nil {
			log.Printf("WARN: Failed to insert user %s: %v", u.Email, err)
		}

		// Get the actual ID (in case it existed)
		err = db.QueryRow("SELECT id FROM users WHERE email = ?", u.Email).Scan(&userID)
		if err != nil {
			log.Fatalf("ERROR: Failed to retrieve user ID for %s: %v", u.Email, err)
		}
		userIDs[u.Email] = userID

		// Insert/Update Profile
		_, err = db.Exec("INSERT INTO profiles (user_id, display_name) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET display_name = excluded.display_name",
			userID, u.DisplayName)
		if err != nil {
			log.Printf("WARN: Failed to insert profile for %s: %v", u.Email, err)
		}

		log.Printf("INFO: User ready: %s (%s)", u.DisplayName, userID)
	}

	// 4. Seed Predictions for 2025 and 2026
	years := []int{2025, 2026}
	nowMS := time.Now().UnixMilli()

	for _, year := range years {
		log.Printf("INFO: Processing schedule for %d...", year)
		schedule, err := f1Service.GetScheduleByYear(ctx, year)
		if err != nil {
			log.Printf("ERROR: Failed to fetch schedule for %d: %v", year, err)
			continue
		}

		for _, weekend := range schedule {
			log.Printf("  Round %d: %s", weekend.Round, weekend.FullName)

			// Fetch drivers for this weekend
			drivers, err := f1Service.GetDrivers(ctx, year, weekend.Round)
			if err != nil {
				log.Printf("    ERROR: Failed to fetch drivers for round %d: %v", weekend.Round, err)
				continue
			}

			if len(drivers) == 0 {
				continue
			}

			for _, session := range weekend.Sessions {
				// Only seed for completed sessions
				if session.TimeUTCMS >= nowMS {
					continue
				}

				log.Printf("    Seeding session: %s", session.Type)

				for _, email := range []string{"test1@test.com", "test2@test.com"} {
					userID := userIDs[email]
					seedRandomPrediction(db, userID, year, weekend.Round, session.Type, drivers)
				}
			}
		}
	}

	log.Println("SUCCESS: Seeding completed!")
}

func seedRandomPrediction(db *sql.DB, userID string, year, round int, sessionType string, drivers []models.DriverInfo) {
	// 1. Create Prediction Header
	predictionID := uuid.New().String()
	now := time.Now().UTC()

	_, err := db.Exec(`
		INSERT INTO predictions (id, user_id, year, round, session_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, year, round, session_type) DO UPDATE SET updated_at = excluded.updated_at
	`, predictionID, userID, year, round, sessionType, now, now)

	if err != nil {
		log.Printf("      ERROR: Failed to upsert prediction header: %v", err)
		return
	}

	// Get the actual ID if it was an update
	db.QueryRow("SELECT id FROM predictions WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?",
		userID, year, round, sessionType).Scan(&predictionID)

	// 2. Clear existing entries
	db.Exec("DELETE FROM prediction_entries WHERE prediction_id = ?", predictionID)

	// 3. Generate Random Top 10
	shuffled := make([]models.DriverInfo, len(drivers))
	copy(shuffled, drivers)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := 10
	if len(shuffled) < 10 {
		count = len(shuffled)
	}

	for i := 0; i < count; i++ {
		_, err = db.Exec("INSERT INTO prediction_entries (prediction_id, position, driver_id) VALUES (?, ?, ?)",
			predictionID, i+1, shuffled[i].ID)
		if err != nil {
			log.Printf("      ERROR: Failed to insert entry for %s: %v", shuffled[i].ID, err)
		}
	}
}
