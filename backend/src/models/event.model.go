package models

import "time"

type Event struct {
	EventId             string     `json:"eventId"`
	EventName           string     `json:"eventName"`
	CoverUrl            *string    `json:"coverUrl"`
	OwnerUserId         uint32     `json:"ownerUserId"`
	OwnerOrganizationId string     `json:"ownerOrganizationId"`
	PaymentId           *string    `json:"paymentId"`
	CreatedAt           time.Time  `json:"createdAt"`
	Exp                 *time.Time `json:"exp"`
	Tags                []string   `json:"tags"`
	Description         string     `json:"description"`
}
