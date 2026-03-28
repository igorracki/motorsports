package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/igorracki/motorsports/backend/internal/models"
)

type PredictionRepository interface {
	SavePrediction(ctx context.Context, prediction *models.Prediction) error
	GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error)
	GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error)
	GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error)
	GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) (map[string]int, error)
}

type predictionRepository struct {
	database *sql.DB
}

func NewPredictionRepository(db *sql.DB) PredictionRepository {
	return &predictionRepository{database: db}
}

func (predictionRepo *predictionRepository) SavePrediction(ctx context.Context, prediction *models.Prediction) error {
	slog.InfoContext(ctx, "Entry: SavePrediction", "user_id", prediction.UserID, "year", prediction.Year, "round", prediction.Round, "session", prediction.SessionType)

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

	slog.InfoContext(ctx, "Exit: SavePrediction", "prediction_id", prediction.ID)
	return nil
}

func (predictionRepo *predictionRepository) GetPrediction(ctx context.Context, userID string, year, round int, sessionType string) (*models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetPrediction", "user_id", userID, "year", year, "round", round, "session", sessionType)

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
			slog.InfoContext(ctx, "No prediction found", "user_id", userID, "year", year, "round", round)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction entries: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetPrediction", "prediction_id", prediction.ID)
	return prediction, nil
}

func (predictionRepo *predictionRepository) GetUserPredictions(ctx context.Context, userID string) ([]models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetUserPredictions", "user_id", userID)

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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user predictions: %w", err)
	}

	if err := predictionRepo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, fmt.Errorf("fetching entries: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetUserPredictions", "user_id", userID, "count", len(predictions))
	return predictions, nil
}

func (predictionRepo *predictionRepository) GetRoundPredictions(ctx context.Context, userID string, year, round int) ([]models.Prediction, error) {
	slog.InfoContext(ctx, "Entry: GetRoundPredictions", "user_id", userID, "year", year, "round", round)

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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating round predictions: %w", err)
	}

	if err := predictionRepo.fetchEntriesForPredictions(ctx, predictions); err != nil {
		return nil, fmt.Errorf("fetching entries: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetRoundPredictions", "user_id", userID, "count", len(predictions))
	return predictions, nil
}

func (predictionRepo *predictionRepository) GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) (map[string]int, error) {
	slog.InfoContext(ctx, "Entry: GetSeasonScoresByUserIDs", "count", len(userIDs), "season", season)

	if len(userIDs) == 0 {
		return make(map[string]int), nil
	}

	uniqueUserIDs := make([]string, 0, len(userIDs))
	seen := make(map[string]struct{})
	for _, id := range userIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			uniqueUserIDs = append(uniqueUserIDs, id)
		}
	}

	query := "SELECT user_id, COALESCE(SUM(score), 0) FROM predictions WHERE year = ? AND user_id IN ("
	args := make([]interface{}, len(uniqueUserIDs)+1)
	args[0] = season
	for i, id := range uniqueUserIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i+1] = id
	}
	query += ") GROUP BY user_id"

	rows, err := predictionRepo.database.QueryContext(ctx, query, args...)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating aggregated scores: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetSeasonScoresByUserIDs", "found", len(userScores))
	return userScores, nil
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

		for rows.Next() {
			var entry models.PredictionEntry
			entry.PredictionID = predictions[i].ID
			if err := rows.Scan(&entry.Position, &entry.DriverID, &entry.Correct); err != nil {
				rows.Close()
				return fmt.Errorf("scanning entry for prediction %s: %w", predictions[i].ID, err)
			}
			predictions[i].Entries = append(predictions[i].Entries, entry)
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return fmt.Errorf("iterating entries for prediction %s: %w", predictions[i].ID, err)
		}
		rows.Close()
	}
	return nil
}
