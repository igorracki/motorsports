package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/models"
)

type ScoreRepository struct {
	database *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{database: db}
}

func (scoreRepo *ScoreRepository) UpdateScore(score *models.UserScore) error {
	seasonValue := "NULL"
	if score.Season != nil {
		seasonValue = fmt.Sprintf("%d", *score.Season)
	}
	log.Printf("INFO: Attempting to update score [user_id: %s, type: %s, season: %s]",
		score.UserID, score.ScoreType, seasonValue)

	_, err := scoreRepo.database.Exec(`
		INSERT INTO user_scores (user_id, score_type, season, value, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, score_type, season) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, score.UserID, score.ScoreType, score.Season, score.Value, score.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upserting user score for user %s: %w", score.UserID, err)
	}

	log.Printf("INFO: Successfully updated score [user_id: %s, type: %s, new_value: %d]",
		score.UserID, score.ScoreType, score.Value)
	return nil
}

func (scoreRepo *ScoreRepository) GetUserScores(userID string) ([]models.UserScore, error) {
	log.Printf("INFO: Fetching all scores for user [id: %s]", userID)

	rows, err := scoreRepo.database.Query(`
		SELECT user_id, score_type, season, value, updated_at 
		FROM user_scores 
		WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying user scores for %s: %w", userID, err)
	}
	defer rows.Close()

	var scores []models.UserScore
	for rows.Next() {
		var score models.UserScore
		if err := rows.Scan(&score.UserID, &score.ScoreType, &score.Season, &score.Value, &score.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning user score: %w", err)
		}
		scores = append(scores, score)
	}

	log.Printf("INFO: Successfully fetched %d scores for user %s", len(scores), userID)
	return scores, nil
}
