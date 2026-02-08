package services

import (
	"context"
	"fmt"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/models"
)

type F1Service interface {
	GetEventsByYear(context context.Context, year int) ([]models.Event, error)
}

type f1Service struct {
	client clients.F1DataClient
}

func NewF1Service(client clients.F1DataClient) F1Service {
	return &f1Service{
		client: client,
	}
}

func (service *f1Service) GetEventsByYear(context context.Context, year int) ([]models.Event, error) {
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
