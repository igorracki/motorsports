package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/models"
)

const (
	upsertPredictionSQL = `
		INSERT INTO predictions (id, user_id, year, round, session_type, score, revalidate_until, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, year, round, session_type) DO UPDATE SET
			score = excluded.score,
			revalidate_until = excluded.revalidate_until,
			updated_at = excluded.updated_at`
	getPredictionMetadataSQL = `
		SELECT id, created_at FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`
	deletePredictionEntriesSQL = "DELETE FROM prediction_entries WHERE prediction_id = ?"
	insertPredictionEntrySQL   = "INSERT INTO prediction_entries (prediction_id, position, driver_id, correct) VALUES (?, ?, ?, ?)"
	getPredictionSQL           = `
		SELECT id, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ? AND session_type = ?`
	getPredictionEntriesSQL = `
		SELECT position, driver_id, correct 
		FROM prediction_entries 
		WHERE prediction_id = ? 
		ORDER BY position ASC`
	getUserPredictionsSQL = `
		SELECT id, year, round, session_type, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ?
		ORDER BY year DESC, round DESC, session_type ASC`
	getRoundPredictionsSQL = `
		SELECT id, session_type, score, revalidate_until, created_at, updated_at 
		FROM predictions 
		WHERE user_id = ? AND year = ? AND round = ?
		ORDER BY session_type ASC`
)

type PredictionRepository interface {
	SavePrediction(ctx context.Context, prediction *models.Prediction) error
	GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error)
	GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error)
	GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error)
	GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) (map[string]int, error)
}

type predictionRepository struct {
	manager *database.Manager
}

func NewPredictionRepository(manager *database.Manager) PredictionRepository {
	return &predictionRepository{manager: manager}
}

func (repo *predictionRepository) SavePrediction(ctx context.Context, prediction *models.Prediction) error {
	return repo.manager.Transaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, upsertPredictionSQL, prediction.ID, prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType,
			prediction.Score, prediction.RevalidateUntil, prediction.CreatedAt, prediction.UpdatedAt)
		if err != nil {
			return fmt.Errorf("upserting prediction header: %w", err)
		}

		err = tx.QueryRowContext(ctx, getPredictionMetadataSQL, prediction.UserID, prediction.Year, prediction.Round, prediction.SessionType).Scan(&prediction.ID, &prediction.CreatedAt)
		if err != nil {
			return fmt.Errorf("retrieving prediction metadata: %w", err)
		}

		if _, err := tx.ExecContext(ctx, deletePredictionEntriesSQL, prediction.ID); err != nil {
			return fmt.Errorf("clearing old prediction entries: %w", err)
		}

		for _, entry := range prediction.Entries {
			if _, err := tx.ExecContext(ctx, insertPredictionEntrySQL, prediction.ID, entry.Position, entry.DriverID, entry.Correct); err != nil {
				return fmt.Errorf("inserting prediction entry [pos: %d, driver: %s]: %w", entry.Position, entry.DriverID, err)
			}
		}
		return nil
	})
}

func (repo *predictionRepository) GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error) {
	prediction := &models.Prediction{
		UserID:      userID,
		Year:        year,
		Round:       round,
		SessionType: sessionType,
		Entries:     []models.PredictionEntry{},
	}

	err := repo.manager.DB().QueryRowContext(ctx, getPredictionSQL, userID, year, round, sessionType).Scan(&prediction.ID, &prediction.Score, &prediction.RevalidateUntil, &prediction.CreatedAt, &prediction.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying prediction: %w", err)
	}

	rows, err := repo.manager.DB().QueryContext(ctx, getPredictionEntriesSQL, prediction.ID)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction entries: %w", err)
	}

	return prediction, nil
}

func (repo *predictionRepository) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	rows, err := repo.manager.DB().QueryContext(ctx, getUserPredictionsSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("querying user predictions: %w", err)
	}
	defer rows.Close()

	predictions := make([]models.Prediction, 0)
	for rows.Next() {
		var p models.Prediction
		p.UserID = userID
		p.Entries = []models.PredictionEntry{}
		if err := rows.Scan(&p.ID, &p.Year, &p.Round, &p.SessionType, &p.Score, &p.RevalidateUntil, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning prediction: %w", err)
		}
		predictions = append(predictions, p)
	}

	if err := repo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, fmt.Errorf("fetching entries: %w", err)
	}

	return predictions, nil
}

func (repo *predictionRepository) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	rows, err := repo.manager.DB().QueryContext(ctx, getRoundPredictionsSQL, userID, year, round)
	if err != nil {
		return nil, fmt.Errorf("querying round predictions: %w", err)
	}
	defer rows.Close()

	predictions := make([]models.Prediction, 0)
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

	if err := repo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, fmt.Errorf("fetching entries: %w", err)
	}

	return predictions, nil
}

func (repo *predictionRepository) GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) (map[string]int, error) {
	if len(userIDs) == 0 {
		return make(map[string]int), nil
	}

	uniqueIDs := make([]string, 0, len(userIDs))
	seen := make(map[string]struct{})
	for _, id := range userIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	placeholders := database.GeneratePlaceholders(len(uniqueIDs))
	query := fmt.Sprintf("SELECT user_id, COALESCE(SUM(score), 0) FROM predictions WHERE year = ? AND user_id IN (%s) GROUP BY user_id", placeholders)

	args := make([]any, 0, len(uniqueIDs)+1)
	args = append(args, season)
	args = append(args, database.ToAnySlice(uniqueIDs)...)

	rows, err := repo.manager.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying season scores from predictions: %w", err)
	}
	defer rows.Close()

	userScores := make(map[string]int)
	for rows.Next() {
		var userID string
		var totalScore int
		if err := rows.Scan(&userID, &totalScore); err != nil {
			return nil, fmt.Errorf("scanning aggregated score: %w", err)
		}
		userScores[userID] = totalScore
	}

	return userScores, nil
}

func (repo *predictionRepository) fetchEntriesForPredictions(ctx context.Context, predictions []models.Prediction) error {
	if len(predictions) == 0 {
		return nil
	}

	for i := range predictions {
		rows, err := repo.manager.DB().QueryContext(ctx, getPredictionEntriesSQL, predictions[i].ID)
		if err != nil {
			return fmt.Errorf("querying entries for prediction %s: %w", predictions[i].ID, err)
		}

		for rows.Next() {
			var entry models.PredictionEntry
			entry.PredictionID = predictions[i].ID
			if err := rows.Scan(&entry.Position, &entry.DriverID, &entry.Correct); err != nil {
				rows.Close()
				return fmt.Errorf("scanning entry for prediction %s: %w", predictions[i].ID, err)
			}
			predictions[i].Entries = append(predictions[i].Entries, entry)
		}
		rows.Close()
	}
	return nil
}
