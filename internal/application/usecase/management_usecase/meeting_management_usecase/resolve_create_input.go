package meeting_mng_usecase

import (
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"

	"github.com/google/uuid"
)

func (uc *MeetingCreator) resolveUserIDFromUUID(userUUID string) (uint, error) {
	parsed, err := uuid.Parse(userUUID)
	if err != nil {
		return 0, err
	}
	user, err := uc.userRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (uc *MeetingCreator) resolveUserIDsFromUUIDs(uuids []string) ([]uint, error) {
	if len(uuids) == 0 {
		return nil, nil
	}
	users, err := uc.userRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(users))
	for idx, user := range users {
		ids[idx] = user.ID
	}
	return ids, nil
}

func (uc *MeetingCreator) resolveTemplateIDFromUUIDString(templateUUID string) (uint, error) {
	parsed, err := uuid.Parse(templateUUID)
	if err != nil {
		return 0, err
	}
	tpl, err := uc.templateRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return tpl.ID, nil
}

func (uc *MeetingCreator) resolvePairingIDs(pairingUUIDs []string) ([]uint, error) {
	if len(pairingUUIDs) == 0 {
		return nil, nil
	}
	ids := make([]uint, len(pairingUUIDs))
	pairs, err := uc.pairingRepo.GetByUUIDs(pairingUUIDs)
	if err != nil {
		return nil, err
	}
	for idx, pair := range pairs {
		ids[idx] = pair.ID
	}
	return ids, nil
}

func (uc *MeetingCreator) ResolveCreateMeetingInput(req *inputmodel.CreateMeetingInput) error {
	if req.MentorUuid == "" || req.TemplateUuid == "" {
		return fmt.Errorf("mentorUuid and templateUuid are required")
	}
	if len(req.MenteeUuids) == 0 {
		return fmt.Errorf("menteeUuids are required")
	}
	mentorID, err := uc.resolveUserIDFromUUID(req.MentorUuid)
	if err != nil {
		return err
	}
	templateID, err := uc.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if err != nil {
		return err
	}
	menteeIDs, err := uc.resolveUserIDsFromUUIDs(req.MenteeUuids)
	if err != nil {
		return err
	}
	sharerIDs, err := uc.resolveUserIDsFromUUIDs(req.SharerUserUuids)
	if err != nil {
		return err
	}
	req.MentorID = mentorID
	req.TemplateID = templateID
	req.MenteeIDs = menteeIDs
	req.SharerUserIDs = sharerIDs
	return nil
}

func (uc *MeetingCreator) ResolveCreateMeetingByAllInput(req *inputmodel.CreateMeetingByAllInput) error {
	if req.TemplateUuid == "" {
		return fmt.Errorf("templateUuid is required")
	}
	templateID, err := uc.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if err != nil {
		return err
	}
	sharerIDs, err := uc.resolveUserIDsFromUUIDs(req.SharerUserUuids)
	if err != nil {
		return err
	}
	pairingIDs, err := uc.resolvePairingIDs(req.PairingUuids)
	if err != nil {
		return err
	}
	req.TemplateID = templateID
	req.SharerUserIDs = sharerIDs
	req.PairingIDs = pairingIDs
	return nil
}

func (uc *MeetingCreator) ResolveCreateMeetingMonthlyInput(req *inputmodel.CreateMeetingMonthlyInput) error {
	if req.MentorUuid == "" || req.TemplateUuid == "" {
		return fmt.Errorf("mentorUuid and templateUuid are required")
	}
	if len(req.MenteeUuids) == 0 {
		return fmt.Errorf("menteeUuids are required")
	}
	mentorID, err := uc.resolveUserIDFromUUID(req.MentorUuid)
	if err != nil {
		return err
	}
	templateID, err := uc.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if err != nil {
		return err
	}
	menteeIDs, err := uc.resolveUserIDsFromUUIDs(req.MenteeUuids)
	if err != nil {
		return err
	}
	sharerIDs, err := uc.resolveUserIDsFromUUIDs(req.SharerUserUuids)
	if err != nil {
		return err
	}
	req.MentorID = mentorID
	req.TemplateID = templateID
	req.MenteeIDs = menteeIDs
	req.SharerUserIDs = sharerIDs
	return nil
}

func (uc *MeetingCreator) ResolveCreateMeetingMonthlyByPairingsInput(req *inputmodel.CreateMeetingMonthlyByPairingsInput) error {
	if req.TemplateUuid == "" {
		return fmt.Errorf("templateUuid is required")
	}
	if len(req.PairingUuids) == 0 {
		return fmt.Errorf("pairingUuids are required")
	}
	templateID, err := uc.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if err != nil {
		return err
	}
	sharerIDs, err := uc.resolveUserIDsFromUUIDs(req.SharerUserUuids)
	if err != nil {
		return err
	}
	pairingIDs, err := uc.resolvePairingIDs(req.PairingUuids)
	if err != nil {
		return err
	}
	req.TemplateID = templateID
	req.SharerUserIDs = sharerIDs
	req.PairingIDs = pairingIDs
	return nil
}