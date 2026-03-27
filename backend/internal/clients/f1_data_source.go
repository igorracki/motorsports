package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/igorracki/motorsports/backend/internal/models"
)

type F1DataClient interface {
	GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
	GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error)
	GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error)
}

type f1DataClient struct {
	baseURL string
	client  *http.Client
}

func NewF1DataClient(baseURL string) F1DataClient {
	return &f1DataClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (client *f1DataClient) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	slog.InfoContext(ctx, "Entry: GetScheduleByYear", "year", year, "url", client.baseURL)

	if year < 1950 || year > 2100 {
		return nil, fmt.Errorf("invalid year parameter: %d", year)
	}

	path, err := url.JoinPath(client.baseURL, "events", fmt.Sprintf("%d", year))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL for year %d: %w", year, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request for year %d: %w", year, err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request for year %d: %w", year, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d for year %d: %s", response.StatusCode, year, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body for year %d: %w", year, err)
	}

	schedule := []models.RaceWeekend{}
	if err := json.Unmarshal(body, &schedule); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response for year %d: %w", year, err)
	}

	slog.InfoContext(ctx, "Exit: GetScheduleByYear", "year", year, "count", len(schedule))
	return schedule, nil
}

func (client *f1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	slog.InfoContext(ctx, "Entry: GetSessionResults", "year", year, "round", round, "sessionType", sessionType)

	if year < 1950 || year > 2100 || round < 1 || round > 50 {
		return nil, fmt.Errorf("invalid year or round parameter")
	}

	path, err := url.JoinPath(client.baseURL, "results", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round), sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL for %d round %d (%s): %w", year, round, sessionType, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request for %d round %d (%s): %w", year, round, sessionType, err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request for %d round %d (%s): %w", year, round, sessionType, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d for %d round %d (%s): %s", response.StatusCode, year, round, sessionType, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body for %d round %d (%s): %w", year, round, sessionType, err)
	}

	results := models.SessionResults{
		Results: []models.DriverResult{},
	}
	if err := json.Unmarshal(body, &results); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response for %d round %d (%s): %w", year, round, sessionType, err)
	}

	slog.InfoContext(ctx, "Exit: GetSessionResults", "year", year, "round", round, "sessionType", sessionType, "drivers", len(results.Results))
	return &results, nil
}

func (client *f1DataClient) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	slog.InfoContext(ctx, "Entry: GetCircuit", "year", year, "round", round)

	if year < 1950 || year > 2100 || round < 1 || round > 50 {
		return nil, fmt.Errorf("invalid year or round parameter")
	}

	path, err := url.JoinPath(client.baseURL, "circuits", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL for circuit %d (%d): %w", round, year, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request for circuit %d (%d): %w", round, year, err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request for circuit %d (%d): %w", round, year, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d for circuit %d (%d): %s", response.StatusCode, round, year, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body for circuit %d (%d): %w", round, year, err)
	}

	circuit := models.Circuit{
		Layout: []models.CircuitLayoutPoint{},
	}
	if err := json.Unmarshal(body, &circuit); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response for circuit %d (%d): %w", round, year, err)
	}

	slog.InfoContext(ctx, "Exit: GetCircuit", "year", year, "round", round, "name", circuit.CircuitName)
	return &circuit, nil
}

func (client *f1DataClient) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	slog.InfoContext(ctx, "Entry: GetDrivers", "year", year, "round", round)

	if year < 1950 || year > 2100 || round < 1 || round > 50 {
		return nil, fmt.Errorf("invalid year or round parameter")
	}

	path, err := url.JoinPath(client.baseURL, "drivers", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL for drivers %d (%d): %w", round, year, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request for drivers %d (%d): %w", round, year, err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request for drivers %d (%d): %w", round, year, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d for drivers %d (%d): %s", response.StatusCode, round, year, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body for drivers %d (%d): %w", round, year, err)
	}

	drivers := []models.DriverInfo{}
	if err := json.Unmarshal(body, &drivers); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response for drivers %d (%d): %w", round, year, err)
	}

	slog.InfoContext(ctx, "Exit: GetDrivers", "year", year, "round", round, "drivers", len(drivers))
	return drivers, nil
}
