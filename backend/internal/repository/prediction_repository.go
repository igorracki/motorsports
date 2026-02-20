package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
)

type PredictionRepository struct {
	database *sql.DB
}

func NewPredictionRepository(db *sql.DB) *PredictionRepository {
	return &PredictionRepository{database: db}
}

func (predictionRepo *PredictionRepository) SavePrediction(prediction *models.Prediction) error {
	log.Printf("INFO: Attempting to save prediction [user_id: %s, year: %d, round: %d, session: %s]",
		prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType)

	transaction, err := predictionRepo.database.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction for prediction: %w", err)
	}
	defer transaction.Rollback()

	_, err = transaction.Exec(`
		INSERT INTO predictions (id, user_id, year, round, session_type, score, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, year, round, session_type) DO UPDATE SET
			score = excluded.score,
			updated_at = excluded.updated_at
	`, prediction.ID, prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType,
		prediction.Score, prediction.CreatedAt, prediction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upserting prediction header: %w", err)
	}

	var actualID string
	err = transaction.QueryRow(`
		SELECT id FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`,
		prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType,
	).Scan(&actualID)
	if err != nil {
		return fmt.Errorf("retrieving prediction ID: %w", err)
	}

	_, err = transaction.Exec("DELETE FROM prediction_entries WHERE prediction_id = ?", actualID)
	if err != nil {
		return fmt.Errorf("clearing old prediction entries: %w", err)
	}

	for _, entry := range prediction.Entries {
		_, err = transaction.Exec(
			"INSERT INTO prediction_entries (prediction_id, position, driver_id) VALUES (?, ?, ?)",
			actualID, entry.Position, entry.DriverID,
		)
		if err != nil {
			return fmt.Errorf("inserting prediction entry [pos: %d, driver: %s]: %w", entry.Position, entry.DriverID, err)
		}
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing prediction [id: %s]: %w", actualID, err)
	}

	log.Printf("INFO: Successfully saved prediction [id: %s, entries: %d]", actualID, len(prediction.Entries))
	return nil
}

func (predictionRepo *PredictionRepository) GetPrediction(userID string, year, round int, sessionType string) (*models.Prediction, error) {
	log.Printf("INFO: Fetching prediction [user_id: %s, year: %d, round: %d, session: %s]",
		userID, year, round, sessionType)

	prediction := &models.Prediction{
		UserID:      userID,
		Year:        year,
		Round:       round,
		SessionType: sessionType,
	}

	err := predictionRepo.database.QueryRow(`
		SELECT id, score, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`,
		userID, year, round, sessionType,
	).Scan(&prediction.ID, &prediction.Score, &prediction.CreatedAt, &prediction.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: No prediction found [user_id: %s, year: %d, round: %d]", userID, year, round)
			return nil, nil
		}
		return nil, fmt.Errorf("querying prediction: %w", err)
	}

	rows, err := predictionRepo.database.Query(`
		SELECT position, driver_id 
		FROM prediction_entries 
		WHERE prediction_id = ? 
		ORDER BY position ASC`,
		prediction.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying prediction entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.PredictionEntry
		entry.PredictionID = prediction.ID
		if err := rows.Scan(&entry.Position, &entry.DriverID); err != nil {
			return nil, fmt.Errorf("scanning prediction entry: %w", err)
		}
		prediction.Entries = append(prediction.Entries, entry)
	}

	log.Printf("INFO: Successfully fetched prediction [id: %s, entries: %d]", prediction.ID, len(prediction.Entries))
	return prediction, nil
}
