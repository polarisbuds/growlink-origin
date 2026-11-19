package planning_usecase

import (
	"errors"
	"fmt"

	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

var ErrBreakPairingForbidden = errors.New("forbidden")

func (uc *PlanningAccessUsecase) BreakPairingByUUID(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	pairingUUID string,
	deleter *pairing_mng_usecase.PairingDeleter,
) error {
	parsed, err := uuid.Parse(pairingUUID)
	if err != nil {
		return err
	}
	pair, err := uc.pairingRepo.GetByUUID(parsed)
	if err != nil {
		return err
	}
	ok, err := uc.CanBreakPairing(userID, authority, systemRole, pair.CreatedByUserID, pair.MenteeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrBreakPairingForbidden
	}
	if deleter == nil {
		return fmt.Errorf("pairing deleter is not configured")
	}
	return deleter.BreakPairing(pair.ID)
}