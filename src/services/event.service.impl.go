package services

import (
	"context"
	"database/sql"

	"github.com/patos-ufscar/quack-week/models"
)

type EventServicePgImpl struct {
	db *sql.DB
}

func NewEventServicePgImpl(db *sql.DB) EventService {
	return &EventServicePgImpl{
		db: db,
	}
}

func (s *EventServicePgImpl) CreateEvent(ctx context.Context, name string, ownerId uint32, orgId string) (models.Event, error) {
	e := models.Event{}
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO events (event_name, owner_user_id, owner_organization_id)
		VALUES
			($1, $2, $3)
		RETURNING
			event_id,
			event_name,
			cover_url,
			owner_user_id,
			owner_organization_id,
			payment_id,
			created_at,
			exp;
		`,
		name,
		ownerId,
		orgId,
	).Scan(
		&e.EventId,
		&e.EventName,
		&e.CoverUrl,
		&e.OwnerUserId,
		&e.OwnerOrganizationId,
		&e.PaymentId,
		&e.CreatedAt,
		&e.Exp,
	)

	return e, err
}

func (s *EventServicePgImpl) SetCover(ctx context.Context, id string, url string) error {
	_, err := s.db.ExecContext(ctx, `
			UPDATE events
			SET 
				cover_url = $1
			WHERE event_id = $2;
		`,
		url,
		id,
	)
	return err
}
