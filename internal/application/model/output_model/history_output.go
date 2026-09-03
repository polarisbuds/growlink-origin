package outputmodel

import "growlink/internal/domain/valueobject"

type HistoryOutput struct {
	ID        uint                  `json:"id"`
	UserID    uint                  `json:"userId"`
	EventType valueobject.EventType `json:"eventType"`
	User      UserOutput            `json:"user"`
}
