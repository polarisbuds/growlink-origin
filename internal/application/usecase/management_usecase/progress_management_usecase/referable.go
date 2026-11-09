package progress_mng_usecase

import "growlink/internal/domain/valueobject"

// ReferableUsers resolves which user IDs may appear in progress aggregates.
type ReferableUsers interface {
	ReferableUserIDs(userID uint, systemRole valueobject.SystemRole) (map[uint]bool, error)
}
