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

	"github.com/igorracki/f1/backend/internal/models"
)

type F1DataClient interface {
	GetRaceWeekendsByYear(ctx context.Context, year int) ([]models.RaceWeekend, error)
	GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error)
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

func (client *f1DataClient) GetRaceWeekendsByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	slog.InfoContext(ctx, "Fetching race weekends", "year", year, "url", client.baseURL)
	path, err := url.JoinPath(client.baseURL, "events", fmt.Sprintf("%d", year))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var raceWeekends []models.RaceWeekend
	if err := json.Unmarshal(body, &raceWeekends); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	slog.InfoContext(ctx, "Successfully fetched race weekends", "count", len(raceWeekends))
	return raceWeekends, nil
}

func (client *f1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	slog.InfoContext(ctx, "Fetching session results", "year", year, "round", round, "sessionType", sessionType)
	path, err := url.JoinPath(client.baseURL, "results", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round), sessionType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute request", "error", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "API returned error status", "status", response.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var results models.SessionResults
	if err := json.Unmarshal(body, &results); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	slog.InfoContext(ctx, "Successfully fetched session results", "drivers", len(results.Results))
	return &results, nil
}
