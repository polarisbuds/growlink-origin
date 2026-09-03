package outputmodel

import (
	"growlink/internal/domain/valueobject"
	"time"
)

type MeetingOutput struct {
	ID            uint                      `json:"id"`
	ImplementedAt time.Time                 `json:"implementedAt"`
	Status        valueobject.MeetingStatus `json:"status"`
	Template      TemplateOutput            `json:"template"`
	Pairing       PairingOutput             `json:"pairing"`
}
