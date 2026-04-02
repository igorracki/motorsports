package services

import (
	"github.com/igorracki/motorsports/backend/internal/models"
)

type ScoringService interface {
	CalculateScore(prediction *models.Prediction, results *models.SessionResults) (int, []bool)
	GetScoringRules() []models.SessionScoringRules
}

type SessionScoringStrategy interface {
	CalculateScore(predictionEntry models.PredictionEntry, actualPosition int) (int, bool)
	GetRules() []models.PositionPoints
}

type basicScoringStrategy struct {
	points map[int]int
}

func (s *basicScoringStrategy) CalculateScore(predictionEntry models.PredictionEntry, actualPosition int) (int, bool) {
	if actualPosition == predictionEntry.Position {
		if pts, ok := s.points[actualPosition]; ok {
			return pts, true
		}
	}
	return 0, false
}

func (s *basicScoringStrategy) GetRules() []models.PositionPoints {
	maxPos := 0
	for pos := range s.points {
		if pos > maxPos {
			maxPos = pos
		}
	}

	rules := make([]models.PositionPoints, 0, len(s.points))
	for i := 1; i <= maxPos; i++ {
		if pts, ok := s.points[i]; ok {
			rules = append(rules, models.PositionPoints{
				Position: i,
				Points:   pts,
			})
		}
	}
	return rules
}

func newPracticeStrategy() SessionScoringStrategy {
	return &basicScoringStrategy{
		points: map[int]int{1: 3, 2: 2, 3: 1},
	}
}

func newQualifyingStrategy() SessionScoringStrategy {
	points := make(map[int]int)
	for i := 1; i <= 12; i++ {
		points[i] = 13 - i
	}
	return &basicScoringStrategy{points: points}
}

func newSprintStrategy() SessionScoringStrategy {
	return &basicScoringStrategy{
		points: map[int]int{
			1: 18, 2: 17, 3: 16, 4: 15, 5: 14, 6: 13, 7: 12,
			8: 10, 9: 9, 10: 8, 11: 7, 12: 6, 13: 5, 14: 4, 15: 3, 16: 2, 17: 1,
		},
	}
}

func newRaceStrategy() SessionScoringStrategy {
	return &basicScoringStrategy{
		points: map[int]int{
			1: 25, 2: 24, 3: 23, 4: 20, 5: 18, 6: 17, 7: 16, 8: 15, 9: 14, 10: 13,
			11: 12, 12: 11, 13: 10, 14: 9, 15: 8, 16: 7, 17: 6, 18: 5, 19: 4, 20: 3, 21: 2, 22: 1,
		},
	}
}

type scoringService struct {
	strategies map[string]SessionScoringStrategy
}

func NewScoringService() ScoringService {
	s := &scoringService{
		strategies: make(map[string]SessionScoringStrategy),
	}
	s.initializeStrategies()
	return s
}

func (s *scoringService) initializeStrategies() {
	practice := newPracticeStrategy()
	s.strategies[models.SessionTypePractice1] = practice
	s.strategies[models.SessionTypePractice1Short] = practice
	s.strategies[models.SessionTypePractice2] = practice
	s.strategies[models.SessionTypePractice2Short] = practice
	s.strategies[models.SessionTypePractice3] = practice
	s.strategies[models.SessionTypePractice3Short] = practice

	quali := newQualifyingStrategy()
	s.strategies[models.SessionTypeQualifying] = quali
	s.strategies[models.SessionTypeQualifyingShort] = quali
	s.strategies[models.SessionTypeSprintQualifying] = quali
	s.strategies[models.SessionTypeSprintQualifyingShort] = quali

	sprint := newSprintStrategy()
	s.strategies[models.SessionTypeSprint] = sprint
	s.strategies[models.SessionTypeSprintShort] = sprint

	race := newRaceStrategy()
	s.strategies[models.SessionTypeRace] = race
	s.strategies[models.SessionTypeRaceShort] = race
}

func (s *scoringService) CalculateScore(prediction *models.Prediction, results *models.SessionResults) (int, []bool) {
	if prediction == nil || results == nil {
		return 0, nil
	}

	strategy, ok := s.strategies[results.SessionType]
	if !ok {
		return 0, nil
	}

	actualPositions := make(map[string]int)
	for _, res := range results.Results {
		actualPositions[res.Driver.ID] = res.Position
	}

	totalScore := 0
	correctness := make([]bool, len(prediction.Entries))

	for i, entry := range prediction.Entries {
		if actualPos, exists := actualPositions[entry.DriverID]; exists {
			points, correct := strategy.CalculateScore(entry, actualPos)
			if correct {
				totalScore += points
				correctness[i] = true
			}
		}
	}

	return totalScore, correctness
}

func (s *scoringService) GetScoringRules() []models.SessionScoringRules {
	types := []string{
		models.SessionTypePractice1,
		models.SessionTypeQualifying,
		models.SessionTypeSprint,
		models.SessionTypeRace,
	}

	result := make([]models.SessionScoringRules, 0, len(types))
	for _, t := range types {
		if strategy, ok := s.strategies[t]; ok {
			result = append(result, models.SessionScoringRules{
				SessionType: t,
				Rules:       strategy.GetRules(),
			})
		}
	}
	return result
}
