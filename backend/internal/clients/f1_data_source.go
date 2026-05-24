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
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*f1DataClient)

func NewF1DataClient(baseURL string, options ...ClientOption) F1DataClient {
	client := &f1DataClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func WithTimeout(timeout int) ClientOption {
	return func(client *f1DataClient) {
		client.httpClient.Timeout = time.Duration(timeout) * time.Second
	}
}

func (client *f1DataClient) GetScheduleByYear(ctx context.Context, year int) ([]models.RaceWeekend, error) {
	return get[[]models.RaceWeekend](client, ctx, "events", fmt.Sprintf("%d", year))
}

func (client *f1DataClient) GetSessionResults(ctx context.Context, year int, round int, sessionType string) (*models.SessionResults, error) {
	results, err := get[models.SessionResults](client, ctx, "results", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round), sessionType)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (client *f1DataClient) GetCircuit(ctx context.Context, year int, round int) (*models.Circuit, error) {
	circuit, err := get[models.Circuit](client, ctx, "circuits", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round))
	if err != nil {
		return nil, err
	}
	return &circuit, nil
}

func (client *f1DataClient) GetDrivers(ctx context.Context, year int, round int) ([]models.DriverInfo, error) {
	return get[[]models.DriverInfo](client, ctx, "drivers", fmt.Sprintf("%d", year), fmt.Sprintf("%d", round))
}

func get[T any](client *f1DataClient, ctx context.Context, paths ...string) (T, error) {
	var zero T
	fullURL, err := url.JoinPath(client.baseURL, paths...)
	if err != nil {
		return zero, fmt.Errorf("constructing URL: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return zero, fmt.Errorf("creating request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return zero, fmt.Errorf("executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return zero, fmt.Errorf("%w: remote resource not found at %s", models.ErrNotFound, fullURL)
		}
		body, _ := io.ReadAll(response.Body)
		slog.ErrorContext(ctx, "External API error", "status", response.StatusCode, "url", fullURL, "body", string(body))
		return zero, fmt.Errorf("API returned status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return zero, fmt.Errorf("reading response body: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return zero, fmt.Errorf("unmarshaling response: %w", err)
	}

	return result, nil
}
