package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/igorracki/motorsports/backend/internal/models"
)

type ScoreRepository interface {
	UpdateScore(ctx context.Context, score *models.UserScore) error
	GetUserScores(ctx context.Context, userID string) ([]models.UserScore, error)
	GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) ([]models.UserScore, error)
}

type scoreRepository struct {
	database *sql.DB
}

func NewScoreRepository(db *sql.DB) ScoreRepository {
	return &scoreRepository{database: db}
}

func (scoreRepo *scoreRepository) UpdateScore(ctx context.Context, score *models.UserScore) error {
	seasonValue := "NULL"
	if score.Season != nil {
		seasonValue = fmt.Sprintf("%d", *score.Season)
	}
	slog.InfoContext(ctx, "Entry: UpdateScore", "user_id", score.UserID, "type", score.ScoreType, "season", seasonValue)

	_, err := scoreRepo.database.ExecContext(ctx, `
		INSERT INTO user_scores (user_id, score_type, season, value, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, score_type, season) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, score.UserID, score.ScoreType, score.Season, score.Value, score.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upserting user score for user %s: %w", score.UserID, err)
	}

	slog.InfoContext(ctx, "Exit: UpdateScore", "user_id", score.UserID, "type", score.ScoreType, "new_value", score.Value)
	return nil
}

func (scoreRepo *scoreRepository) GetUserScores(ctx context.Context, userID string) ([]models.UserScore, error) {
	slog.InfoContext(ctx, "Entry: GetUserScores", "user_id", userID)

	rows, err := scoreRepo.database.QueryContext(ctx, `
		SELECT user_id, score_type, season, value, updated_at 
		FROM user_scores 
		WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying user scores for %s: %w", userID, err)
	}
	defer rows.Close()

	scores := []models.UserScore{}
	for rows.Next() {
		var score models.UserScore
		if err := rows.Scan(&score.UserID, &score.ScoreType, &score.Season, &score.Value, &score.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning user score: %w", err)
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user scores: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetUserScores", "user_id", userID, "count", len(scores))
	return scores, nil
}

func (scoreRepo *scoreRepository) GetSeasonScoresByUserIDs(ctx context.Context, userIDs []string, season int) ([]models.UserScore, error) {
	slog.InfoContext(ctx, "Entry: GetSeasonScoresByUserIDs", "count", len(userIDs), "season", season)

	if len(userIDs) == 0 {
		return []models.UserScore{}, nil
	}

	query := "SELECT user_id, score_type, season, value, updated_at FROM user_scores WHERE season = ? AND user_id IN ("
	args := make([]interface{}, len(userIDs)+1)
	args[0] = season
	for i, id := range userIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i+1] = id
	}
	query += ")"

	rows, err := scoreRepo.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying season scores: %w", err)
	}
	defer rows.Close()

	scores := []models.UserScore{}
	for rows.Next() {
		var s models.UserScore
		if err := rows.Scan(&s.UserID, &s.ScoreType, &s.Season, &s.Value, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning user score: %w", err)
		}
		scores = append(scores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating season scores: %w", err)
	}

	slog.InfoContext(ctx, "Exit: GetSeasonScoresByUserIDs", "found", len(scores))
	return scores, nil
}
