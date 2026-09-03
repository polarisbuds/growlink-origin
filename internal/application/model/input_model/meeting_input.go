package inputmodel

import (
	"growlink/internal/domain/valueobject"
	"time"
)

type MeetingInput struct {
	ImplementedAt time.Time                 `json:"implementedAt"`
	Status        valueobject.MeetingStatus `json:"status"`
	TemplateID    uint                      `json:"templateId"`
	PairingID     uint                      `json:"pairingId"`
}
