package planning_usecase

import (
	"errors"
	"fmt"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

var (
	ErrPairingSharerForbidden = errors.New("forbidden")
	ErrPairingSharerInvalid   = errors.New("invalid sharer uuid")
)

func (uc *PlanningAccessUsecase) ListPairingSharersByUUID(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	pairingUUID string,
) ([]outputmodel.UserOutput, error) {
	parsed, err := uuid.Parse(pairingUUID)
	if err != nil {
		return nil, err
	}
	pair, err := uc.pairingRepo.GetByUUID(parsed)
	if err != nil {
		return nil, err
	}
	ok, err := uc.CanManagePairingMentee(userID, authority, systemRole, pair.MenteeID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrPairingSharerForbidden
	}
	if uc.pairingUpdater == nil {
		return nil, fmt.Errorf("pairing updater is not configured")
	}
	return uc.pairingUpdater.ListSharers(pair.ID)
}

func (uc *PlanningAccessUsecase) ReplacePairingSharersByUUID(
	userID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
	pairingUUID string,
	sharerUUIDs []string,
) ([]outputmodel.UserOutput, error) {
	parsed, err := uuid.Parse(pairingUUID)
	if err != nil {
		return nil, err
	}
	pair, err := uc.pairingRepo.GetByUUID(parsed)
	if err != nil {
		return nil, err
	}
	ok, err := uc.CanManagePairingMentee(userID, authority, systemRole, pair.MenteeID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrPairingSharerForbidden
	}
	sharerIDs, err := uc.resolveSharerUserIDs(sharerUUIDs)
	if err != nil {
		return nil, err
	}
	if uc.pairingUpdater == nil {
		return nil, fmt.Errorf("pairing updater is not configured")
	}
	if err = uc.pairingUpdater.ReplaceSharers(pair.ID, sharerIDs); err != nil {
		return nil, err
	}
	return uc.pairingUpdater.ListSharers(pair.ID)
}

func (uc *PlanningAccessUsecase) resolveSharerUserIDs(uuids []string) ([]uint, error) {
	if len(uuids) == 0 {
		return nil, nil
	}
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(uuids))
	for _, u := range uuids {
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		unique = append(unique, u)
	}
	if len(unique) == 0 {
		return nil, nil
	}
	if uc.userRepo == nil {
		return nil, fmt.Errorf("user repository is not configured")
	}
	users, err := uc.userRepo.GetByUUIDs(unique)
	if err != nil {
		return nil, err
	}
	if len(users) != len(unique) {
		return nil, ErrPairingSharerInvalid
	}
	ids := make([]uint, len(users))
	for idx, u := range users {
		ids[idx] = u.ID
	}
	return ids, nil
}