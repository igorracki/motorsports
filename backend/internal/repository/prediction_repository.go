package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
)

type PredictionRepository interface {
	SavePrediction(ctx context.Context, prediction *models.Prediction) error
	GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error)
	GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error)
	GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error)
}

type predictionRepository struct {
	database *sql.DB
}

func NewPredictionRepository(db *sql.DB) PredictionRepository {
	return &predictionRepository{database: db}
}

func (predictionRepo *predictionRepository) SavePrediction(ctx context.Context, prediction *models.Prediction) error {
	log.Printf("INFO: Attempting to save prediction [user_id: %s, year: %d, round: %d, session: %s]",
		prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType)

	transaction, err := predictionRepo.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction for prediction: %w", err)
	}
	defer transaction.Rollback()

	_, err = transaction.ExecContext(ctx, `
		INSERT INTO predictions (id, user_id, year, round, session_type, score, revalidate_until, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, year, round, session_type) DO UPDATE SET
			score = excluded.score,
			revalidate_until = excluded.revalidate_until,
			updated_at = excluded.updated_at
	`, prediction.ID, prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType,
		prediction.Score, prediction.RevalidateUntil, prediction.CreatedAt, prediction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upserting prediction header: %w", err)
	}

	// Fetch the final ID and CreatedAt (handles retrieving original values if it was an update)
	err = transaction.QueryRowContext(ctx, `
		SELECT id, created_at FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`,
		prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType,
	).Scan(&prediction.ID, &prediction.CreatedAt)
	if err != nil {
		return fmt.Errorf("retrieving prediction metadata: %w", err)
	}

	_, err = transaction.ExecContext(ctx, "DELETE FROM prediction_entries WHERE prediction_id = ?", prediction.ID)
	if err != nil {
		return fmt.Errorf("clearing old prediction entries: %w", err)
	}

	for _, entry := range prediction.Entries {
		_, err = transaction.ExecContext(ctx,
			"INSERT INTO prediction_entries (prediction_id, position, driver_id, correct) VALUES (?, ?, ?, ?)",
			prediction.ID, entry.Position, entry.DriverID, entry.Correct,
		)
		if err != nil {
			return fmt.Errorf("inserting prediction entry [pos: %d, driver: %s]: %w", entry.Position, entry.DriverID, err)
		}
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing prediction [id: %s]: %w", prediction.ID, err)
	}

	log.Printf("INFO: Successfully saved prediction [id: %s, entries: %d]", prediction.ID, len(prediction.Entries))
	return nil
}

func (predictionRepo *predictionRepository) GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error) {
	log.Printf("INFO: Fetching prediction [user_id: %s, year: %d, round: %d, session: %s]",
		userID, year, round, sessionType)

	prediction := &models.Prediction{
		UserID:      userID,
		Year:        year,
		Round:       round,
		SessionType: sessionType,
		Entries:     []models.PredictionEntry{},
	}

	err := predictionRepo.database.QueryRowContext(ctx, `
		SELECT id, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`,
		userID, year, round, sessionType,
	).Scan(&prediction.ID, &prediction.Score, &prediction.RevalidateUntil, &prediction.CreatedAt, &prediction.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: No prediction found [user_id: %s, year: %d, round: %d]", userID, year, round)
			return nil, nil
		}
		return nil, fmt.Errorf("querying prediction: %w", err)
	}

	rows, err := predictionRepo.database.QueryContext(ctx, `
		SELECT position, driver_id, correct 
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
		if err := rows.Scan(&entry.Position, &entry.DriverID, &entry.Correct); err != nil {
			return nil, fmt.Errorf("scanning prediction entry: %w", err)
		}
		prediction.Entries = append(prediction.Entries, entry)
	}

	log.Printf("INFO: Successfully fetched prediction [id: %s, entries: %d]", prediction.ID, len(prediction.Entries))
	return prediction, nil
}

func (predictionRepo *predictionRepository) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	log.Printf("INFO: Fetching all predictions for user [user_id: %s]", userID)

	rows, err := predictionRepo.database.QueryContext(ctx, `
		SELECT id, year, round, session_type, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ?
		ORDER BY year DESC, round DESC, session_type ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying user predictions: %w", err)
	}
	defer rows.Close()

	predictions := []models.Prediction{}
	for rows.Next() {
		var p models.Prediction
		p.UserID = userID
		p.Entries = []models.PredictionEntry{}
		if err := rows.Scan(&p.ID, &p.Year, &p.Round, &p.SessionType, &p.Score, &p.RevalidateUntil, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning prediction: %w", err)
		}
		predictions = append(predictions, p)
	}

	if err := predictionRepo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, err
	}

	log.Printf("INFO: Successfully fetched %d predictions for user [user_id: %s]", len(predictions), userID)
	return predictions, nil
}

func (predictionRepo *predictionRepository) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	log.Printf("INFO: Fetching round predictions [user_id: %s, year: %d, round: %d]", userID, year, round)

	rows, err := predictionRepo.database.QueryContext(ctx, `
		SELECT id, session_type, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ?
		ORDER BY session_type ASC`,
		userID, year, round,
	)
	if err != nil {
		return nil, fmt.Errorf("querying round predictions: %w", err)
	}
	defer rows.Close()

	predictions := []models.Prediction{}
	for rows.Next() {
		var p models.Prediction
		p.UserID = userID
		p.Year = year
		p.Round = round
		p.Entries = []models.PredictionEntry{}
		if err := rows.Scan(&p.ID, &p.SessionType, &p.Score, &p.RevalidateUntil, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning round prediction: %w", err)
		}
		predictions = append(predictions, p)
	}

	if err := predictionRepo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, err
	}

	log.Printf("INFO: Successfully fetched %d round predictions [user_id: %s, year: %d, round: %d]",
		len(predictions), userID, year, round)
	return predictions, nil
}

func (predictionRepo *predictionRepository) fetchEntriesForPredictions(ctx context.Context, predictions []models.Prediction) error {
	if len(predictions) == 0 {
		return nil
	}

	for i := range predictions {
		rows, err := predictionRepo.database.QueryContext(ctx, `
			SELECT position, driver_id, correct 
			FROM prediction_entries 
			WHERE prediction_id = ? 
			ORDER BY position ASC`,
			predictions[i].ID,
		)
		if err != nil {
			return fmt.Errorf("querying entries for prediction %s: %w", predictions[i].ID, err)
		}
		defer rows.Close()

		for rows.Next() {
			var entry models.PredictionEntry
			entry.PredictionID = predictions[i].ID
			if err := rows.Scan(&entry.Position, &entry.DriverID, &entry.Correct); err != nil {
				return fmt.Errorf("scanning entry for prediction %s: %w", predictions[i].ID, err)
			}
			predictions[i].Entries = append(predictions[i].Entries, entry)
		}
	}
	return nil
}
