package services

import (
	"context"
	"fmt"
	"time"

	"github.com/igorracki/motorsports/backend/internal/cache"
	"github.com/igorracki/motorsports/backend/internal/clients"
	"github.com/igorracki/motorsports/backend/internal/formatters"
	"github.com/igorracki/motorsports/backend/internal/mappers"
	"github.com/igorracki/motorsports/backend/internal/models"
)

type F1Service interface {
	GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
	GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error)
	GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error)
	Close()
}

type f1Service struct {
	client clients.F1DataClient
	policy PredictionPolicy
}

func NewF1Service(client clients.F1DataClient, policy PredictionPolicy) F1Service {
	return &f1Service{
		client: client,
		policy: policy,
	}
}

func (service *f1Service) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	schedule, err := service.client.GetScheduleByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("fetching schedule: %w", err)
	}

	filteredSchedule := make([]models.RaceWeekend, 0)
	for i := range schedule {
		if schedule[i].Round == 0 {
			continue
		}
		mappers.MapRaceWeekend(&schedule[i])
		for j := range schedule[i].Sessions {
			session := &schedule[i].Sessions[j]
			session.IsLocked = service.policy.IsLocked(session.TimeUTCMS)
			session.IsLive = service.policy.IsLive(session.TimeUTCMS)
			session.IsCompleted = service.policy.IsSessionCompleted(session.TimeUTCMS)
		}
		filteredSchedule = append(filteredSchedule, schedule[i])
	}

	return filteredSchedule, nil
}

func (service *f1Service) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	results, err := service.client.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		return nil, fmt.Errorf("fetching results: %w", err)
	}

	if results == nil {
		return nil, nil
	}

	if len(results.Results) == 0 {
		return &models.SessionResults{
			Year:        year,
			Round:       round,
			SessionType: sessionType,
			Results:     []models.DriverResult{},
		}, nil
	}

	mappers.MapSessionResults(results)
	return results, nil
}

func (service *f1Service) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	circuit, err := service.client.GetCircuit(ctx, year, round)
	if err != nil {
		return nil, fmt.Errorf("fetching circuit: %w", err)
	}

	if circuit != nil {
		mappers.MapCircuit(circuit)
	}

	return circuit, nil
}

func (service *f1Service) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	if pastDrivers, found, err := service.getDriversFromPastSessions(ctx, year, round); err == nil && found {
		for i := range pastDrivers {
			pastDrivers[i].CountryCode = formatters.GetDriverCountryCode(pastDrivers[i].CountryCode, pastDrivers[i].ID)
		}
		return pastDrivers, nil
	}

	drivers, err := service.client.GetDrivers(ctx, year, round)
	if err != nil {
		return nil, fmt.Errorf("fetching drivers: %w", err)
	}

	if drivers == nil {
		drivers = []models.DriverInfo{}
	}

	for i := range drivers {
		drivers[i].CountryCode = formatters.GetDriverCountryCode(drivers[i].CountryCode, drivers[i].ID)
	}

	return drivers, nil
}

func (service *f1Service) getDriversFromPastSessions(ctx context.Context, year int, round int) ([]models.DriverInfo, bool, error) {
	schedule, err := service.GetScheduleByYear(ctx, year)
	if err != nil {
		return nil, false, err
	}

	relevantWeekend := findWeekendByRound(schedule, round)
	if relevantWeekend == nil {
		return nil, false, nil
	}

	now := time.Now().UnixMilli()
	sessions := sortSessionsByTimeDescending(relevantWeekend.Sessions)

	for _, session := range sessions {
		if session.TimeUTCMS >= now {
			continue
		}

		sessionCode := getSessionCodeForSearch(session)
		results, err := service.GetSessionResults(ctx, year, round, sessionCode)
		if err == nil && results != nil && len(results.Results) > 0 {
			return extractDriversFromResults(results.Results), true, nil
		}
	}

	return nil, false, nil
}

func findWeekendByRound(schedule []models.RaceWeekend, round int) *models.RaceWeekend {
	for i := range schedule {
		if schedule[i].Round == round {
			return &schedule[i]
		}
	}
	return nil
}

func sortSessionsByTimeDescending(sessions []models.Session) []models.Session {
	sorted := make([]models.Session, len(sessions))
	copy(sorted, sessions)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].TimeUTCMS < sorted[j].TimeUTCMS {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func getSessionCodeForSearch(session models.Session) string {
	if session.SessionCode != "" {
		return session.SessionCode
	}
	return formatters.GetSessionCode(session.Type)
}

func extractDriversFromResults(results []models.DriverResult) []models.DriverInfo {
	drivers := make([]models.DriverInfo, 0, len(results))
	for _, res := range results {
		drivers = append(drivers, res.Driver)
	}
	return drivers
}

func (service *f1Service) Close() {}

// --- Caching Decorator ---

const (
	activeWindowBefore = 30 * time.Minute
	activeWindowAfter  = 6 * time.Hour
	historicalTTL      = 24 * time.Hour
	futureTTL          = 1 * time.Hour
)

type f1CachingService struct {
	base          F1Service
	scheduleCache cache.Cache[[]models.RaceWeekend]
	resultsCache  cache.Cache[*models.SessionResults]
	circuitCache  cache.Cache[*models.Circuit]
	driversCache  cache.Cache[[]models.DriverInfo]
}

func NewF1CachingService(base F1Service) F1Service {
	return &f1CachingService{
		base:          base,
		scheduleCache: cache.NewMemoryCache[[]models.RaceWeekend](),
		resultsCache:  cache.NewMemoryCache[*models.SessionResults](),
		circuitCache:  cache.NewMemoryCache[*models.Circuit](),
		driversCache:  cache.NewMemoryCache[[]models.DriverInfo](),
	}
}

func (service *f1CachingService) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	cacheKey := fmt.Sprintf("schedule:%d", year)
	if schedule, found := service.scheduleCache.Get(cacheKey); found {
		return schedule, nil
	}

	schedule, err := service.base.GetScheduleByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	service.scheduleCache.Set(cacheKey, schedule, 24*time.Hour)
	return schedule, nil
}

func (service *f1CachingService) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	cacheKey := fmt.Sprintf("results:%d:%d:%s", year, round, sessionType)
	if results, found := service.resultsCache.Get(cacheKey); found {
		return results, nil
	}

	results, err := service.base.GetSessionResults(ctx, year, round, sessionType)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, nil
	}

	ttl := service.calculateSessionTTL(ctx, year, round, sessionType, len(results.Results) > 0)
	service.resultsCache.Set(cacheKey, results, ttl)

	return results, nil
}

func (service *f1CachingService) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	cacheKey := fmt.Sprintf("circuit:%d:%d", year, round)
	if circuit, found := service.circuitCache.Get(cacheKey); found {
		return circuit, nil
	}

	circuit, err := service.base.GetCircuit(ctx, year, round)
	if err != nil {
		return nil, err
	}

	if circuit != nil {
		ttl := 30 * 24 * time.Hour
		service.circuitCache.Set(cacheKey, circuit, ttl)
	}

	return circuit, nil
}

func (service *f1CachingService) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	cacheKey := fmt.Sprintf("drivers:%d:%d", year, round)
	if drivers, found := service.driversCache.Get(cacheKey); found {
		return drivers, nil
	}

	drivers, err := service.base.GetDrivers(ctx, year, round)
	if err != nil {
		return nil, err
	}

	ttl := service.calculateWeekendTTL(ctx, year, round)
	service.driversCache.Set(cacheKey, drivers, ttl)

	return drivers, nil
}

func (service *f1CachingService) Close() {
	service.scheduleCache.Close()
	service.resultsCache.Close()
	service.circuitCache.Close()
	service.driversCache.Close()
}

func (service *f1CachingService) calculateSessionTTL(ctx context.Context, year int, round int, sessionType string, hasResults bool) time.Duration {
	if !hasResults {
		return 1 * time.Minute
	}

	schedule, err := service.base.GetScheduleByYear(ctx, year)
	if err != nil {
		return 10 * time.Minute
	}

	var sessionTimeMS int64
	found := false
	for _, weekend := range schedule {
		if weekend.Round == round {
			for _, session := range weekend.Sessions {
				if session.Type == sessionType || formatters.GetSessionCode(session.Type) == sessionType {
					sessionTimeMS = session.TimeUTCMS
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		return 10 * time.Minute
	}

	now := time.Now().UnixMilli()

	if now >= sessionTimeMS-activeWindowBefore.Milliseconds() && now <= sessionTimeMS+activeWindowAfter.Milliseconds() {
		return 1 * time.Minute
	}

	if now > sessionTimeMS+activeWindowAfter.Milliseconds() {
		return historicalTTL
	}

	return futureTTL
}

func (service *f1CachingService) calculateWeekendTTL(ctx context.Context, year int, round int) time.Duration {
	schedule, err := service.base.GetScheduleByYear(ctx, year)
	if err != nil {
		return 10 * time.Minute
	}

	var startMS, endMS int64
	found := false
	for _, weekend := range schedule {
		if weekend.Round == round {
			startMS = weekend.StartDateUTCMS
			endMS = weekend.EndDateUTCMS
			found = true
			break
		}
	}

	if !found {
		return 10 * time.Minute
	}

	now := time.Now().UnixMilli()

	if now >= startMS-(24*3600*1000) && now <= endMS+(12*3600*1000) {
		return 1 * time.Hour
	}

	if now > endMS {
		return 7 * 24 * time.Hour
	}

	return 24 * time.Hour
}
