package services

import (
	"context"
	"fmt"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/models"
)

type EventsService interface {
	GetEventsByYear(context context.Context, year int) ([]models.Event, error)
}

type eventsService struct {
	client clients.ExternalAPIClient
}

func NewEventsService(client clients.ExternalAPIClient) EventsService {
	return &eventsService{
		client: client,
	}
}

func (service *eventsService) GetEventsByYear(context context.Context, year int) ([]models.Event, error) {
	if year < 1900 || year > 2050 {
		return nil, fmt.Errorf("year outside supported Formula 1 range")
	}

	events, err := service.client.GetEventsByYear(context, year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events from external API: %w", err)
	}

	// this is where we'd transform the response if needed.

	return events, nil
}
