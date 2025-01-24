package services

import (
	"context"

	"github.com/patos-ufscar/quack-week/models"
)

type EventService interface {
	CreateEvent(ctx context.Context, name string, ownerId uint32, orgId string, description string) (models.Event, error)
	GetEvent(ctx context.Context, eventId string) (models.Event, error)
	SetCover(ctx context.Context, id string, url string) error
	//GetAllEvents(ctx context.Context, limit int) ([]models.Event, error)
}
