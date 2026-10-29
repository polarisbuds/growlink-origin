package outputmodel

import (
	"growlink/internal/domain/valueobject"
	"time"
)

type MeetingOutput struct {
	ID            uint                      `json:"id"`
	UUID          string                    `json:"uuid"`
	ImplementedAt time.Time                 `json:"implementedAt"`
	Status        valueobject.MeetingStatus `json:"status"`
	Template      TemplateOutput            `json:"template"`
	Pairing       PairingOutput             `json:"pairing"`
}
