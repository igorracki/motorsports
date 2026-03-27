package services

import (
	"log/slog"

	"github.com/igorracki/motorsports/backend/internal/models"
)

type ScoringService interface {
	CalculateScore(prediction *models.Prediction, results *models.SessionResults) (int, []bool)
	GetScoringRules() []models.SessionScoringRules
}

type scoringService struct {
	rules map[string]map[int]int
}

func NewScoringService() ScoringService {
	s := &scoringService{
		rules: make(map[string]map[int]int),
	}
	s.initializeRules()
	return s
}

func (s *scoringService) initializeRules() {
	practiceRules := map[int]int{
		1: 3, 2: 2, 3: 1,
	}
	s.rules[models.SessionTypePractice1] = practiceRules
	s.rules[models.SessionTypePractice2] = practiceRules
	s.rules[models.SessionTypePractice3] = practiceRules
	s.rules["FP1"] = practiceRules
	s.rules["FP2"] = practiceRules
	s.rules["FP3"] = practiceRules

	qualiRules := make(map[int]int)
	for i := 1; i <= 12; i++ {
		qualiRules[i] = 13 - i
	}
	s.rules[models.SessionTypeQualifying] = qualiRules
	s.rules[models.SessionTypeQualifyingShort] = qualiRules
	s.rules[models.SessionTypeSprintQualifying] = qualiRules
	s.rules[models.SessionTypeSprintQualifyingShort] = qualiRules

	sprintRules := map[int]int{
		1: 18, 2: 17, 3: 16, 4: 15, 5: 14, 6: 13, 7: 12,
		8: 10, 9: 9, 10: 8, 11: 7, 12: 6, 13: 5, 14: 4, 15: 3, 16: 2, 17: 1,
	}
	s.rules[models.SessionTypeSprint] = sprintRules
	s.rules[models.SessionTypeSprintShort] = sprintRules

	raceRules := map[int]int{
		1: 25, 2: 24, 3: 23, 4: 20, 5: 18, 6: 17, 7: 16, 8: 15, 9: 14, 10: 13,
		11: 12, 12: 11, 13: 10, 14: 9, 15: 8, 16: 7, 17: 6, 18: 5, 19: 4, 20: 3, 21: 2, 22: 1,
	}
	s.rules[models.SessionTypeRace] = raceRules
	s.rules[models.SessionTypeRaceShort] = raceRules
}

func (s *scoringService) CalculateScore(prediction *models.Prediction, results *models.SessionResults) (int, []bool) {
	if prediction == nil || results == nil {
		return 0, nil
	}

	slog.Info("Entry: CalculateScore",
		"user_id", prediction.UserID,
		"year", prediction.Year,
		"round", prediction.Round,
		"session_type", prediction.SessionType,
		"entries_count", len(prediction.Entries),
	)

	sessionRules, ok := s.rules[results.SessionType]
	if !ok {
		slog.Warn("No scoring rules found for session type", "session_type", results.SessionType)
		return 0, nil
	}

	actualPositions := make(map[string]int)
	for _, res := range results.Results {
		actualPositions[res.Driver.ID] = res.Position
	}

	totalScore := 0
	correctness := make([]bool, len(prediction.Entries))

	for i, entry := range prediction.Entries {
		actualPos, exists := actualPositions[entry.DriverID]
		if exists && actualPos == entry.Position {
			points, hasPoints := sessionRules[entry.Position]
			if hasPoints {
				totalScore += points
				correctness[i] = true
			}
		}
	}

	slog.Info("Exit: CalculateScore",
		"user_id", prediction.UserID,
		"total_score", totalScore,
	)

	return totalScore, correctness
}

func (s *scoringService) GetScoringRules() []models.SessionScoringRules {
	slog.Info("Entry: GetScoringRules")
	sessionTypes := []string{
		models.SessionTypePractice1,
		models.SessionTypeQualifying,
		models.SessionTypeSprint,
		models.SessionTypeRace,
	}

	result := make([]models.SessionScoringRules, 0, len(sessionTypes))
	for _, st := range sessionTypes {
		rules, ok := s.rules[st]
		if !ok {
			continue
		}

		posPoints := make([]models.PositionPoints, 0, len(rules))
		maxPos := 0
		for pos := range rules {
			if pos > maxPos {
				maxPos = pos
			}
		}

		for i := 1; i <= maxPos; i++ {
			if pts, ok := rules[i]; ok {
				posPoints = append(posPoints, models.PositionPoints{
					Position: i,
					Points:   pts,
				})
			}
		}

		result = append(result, models.SessionScoringRules{
			SessionType: st,
			Rules:       posPoints,
		})
	}

	slog.Info("Exit: GetScoringRules", "count", len(result))
	return result
}
