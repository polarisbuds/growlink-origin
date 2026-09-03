package inputmodel

import "growlink/internal/domain/valueobject"

type HistoryInput struct {
	UserID    uint                  `json:"userId"`
	EventType valueobject.EventType `json:"eventType"`
}
