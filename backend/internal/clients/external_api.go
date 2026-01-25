package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/igorracki/f1/backend/internal/models"
)

type ExternalAPIClient interface {
	GetEventsByYear(ctx context.Context, year int) ([]models.Event, error)
}

type externalAPIClient struct {
	baseURL string
	client  *http.Client
}

func NewExternalAPIClient(baseURL string) ExternalAPIClient {
	return &externalAPIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (client *externalAPIClient) GetEventsByYear(context context.Context, year int) ([]models.Event, error) {
	url := fmt.Sprintf("%s/wrapper/events/%d", client.baseURL, year)

	request, err := http.NewRequestWithContext(context, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := client.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("API returned status %d: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var events []models.Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return events, nil
}
