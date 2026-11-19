package meeting_mng_usecase

import (
	"context"
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
	"growlink/internal/shared/utils"
	"sync"
	"time"
)

type MeetingCreator struct {
	meetingRepo  repository_interface.MeetingRepository
	pairingRepo  repository_interface.PairingRepository
	userRepo     repository_interface.UserRepository
	templateRepo repository_interface.TemplateRepository
}

func NewMeetingCreator(
	meetingRepo repository_interface.MeetingRepository,
	pairingRepo repository_interface.PairingRepository,
	userRepo repository_interface.UserRepository,
	templateRepo repository_interface.TemplateRepository,
) *MeetingCreator {
	return &MeetingCreator{
		meetingRepo:  meetingRepo,
		pairingRepo:  pairingRepo,
		userRepo:     userRepo,
		templateRepo: templateRepo,
	}
}

func calcImplementedAt(base time.Time) time.Time {
	nextImplementDay := base
	for utils.IsHoliday(nextImplementDay) {
		nextImplementDay = nextImplementDay.AddDate(0, 0, 1)
	}
	return nextImplementDay
}

func implementedAtPtr(base time.Time) *time.Time {
	at := calcImplementedAt(base)
	return &at
}

func (uc *MeetingCreator) saveSharers(meetingID uint, sharerUserIDs []uint) error {
	if len(sharerUserIDs) == 0 {
		return nil
	}
	return uc.meetingRepo.ReplaceSharers(meetingID, sharerUserIDs)
}

func filterPairingsByIDs(pairings []entity.Pairing, pairingIDs []uint) []entity.Pairing {
	pairings = excludeTemporaryPairings(pairings)
	if len(pairingIDs) == 0 {
		return pairings
	}
	allowed := make(map[uint]bool, len(pairingIDs))
	for _, id := range pairingIDs {
		allowed[id] = true
	}
	filtered := make([]entity.Pairing, 0, len(pairingIDs))
	for _, p := range pairings {
		if allowed[p.ID] {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func excludeTemporaryPairings(pairings []entity.Pairing) []entity.Pairing {
	out := make([]entity.Pairing, len(pairings))
	for idx, p := range pairings {
		if !p.IsTemporary {
			out[idx] = p
		}
	}
	return out
}

func (uc *MeetingCreator) CreatePlannedMeetings(mentorID, templateID uint, menteeIDs []uint, implementDate time.Time, createdByUserID uint, sharerUserIDs []uint) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreatePlannedMeetings"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if mentorID == 0 {
		return nil, fmt.Errorf("メンターが指定されていません")
	}
	if len(menteeIDs) == 0 {
		return nil, fmt.Errorf("メンティが指定されていません")
	}
	result := make([]entity.Meeting, len(menteeIDs))
	for idx, menteeID := range menteeIDs {
		if menteeID == 0 {
			return nil, fmt.Errorf("無効なメンティが含まれています")
		}
		if menteeID == mentorID {
			return nil, fmt.Errorf("メンターと同じユーザをメンティに指定できません")
		}
		pair, err := uc.pairingRepo.CreateTemporaryPairing(menteeID, mentorID, createdByUserID)
		if err != nil {
			return nil, err
		}
		meeting, err := uc.meetingRepo.Create(&entity.Meeting{
			ImplementedAt:   implementedAtPtr(implementDate),
			Status:          valueobject.MeetingStatusPending,
			TemplateID:      templateID,
			PairingID:       pair.ID,
			CreatedByUserID: createdByUserID,
		})
		if err != nil {
			return nil, fmt.Errorf("ミーティングの作成に失敗しました: %w", err)
		}
		if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
			return nil, err
		}
		result[idx] = *meeting
	}
	return converter.MeetingEntitiesToOutputs(result), nil
}

func (uc *MeetingCreator) CreateMeetingEveryMonth(mentorID, templateID uint, menteeIDs []uint, from, to time.Time, createdByUserID uint, sharerUserIDs []uint) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreateMeetingEveryMonth"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if mentorID == 0 {
		return nil, fmt.Errorf("メンターが指定されていません")
	}
	if len(menteeIDs) == 0 {
		return nil, fmt.Errorf("メンティが指定されていません")
	}
	implementDates := utils.ListSpecificMonth(from, to)
	if len(implementDates) == 0 {
		return converter.MeetingEntitiesToOutputs(nil), nil
	}
	pairingIDs := make([]uint, len(menteeIDs))
	for idx, menteeID := range menteeIDs {
		if menteeID == 0 {
			return nil, fmt.Errorf("無効なメンティが含まれています")
		}
		if menteeID == mentorID {
			return nil, fmt.Errorf("メンターと同じユーザをメンティに指定できません")
		}
		pair, err := uc.pairingRepo.CreateTemporaryPairing(menteeID, mentorID, createdByUserID)
		if err != nil {
			return nil, err
		}
		pairingIDs[idx] = pair.ID
	}
	var (
		result = []entity.Meeting{}
		mux    = sync.Mutex{}
		wg     = sync.WaitGroup{}
	)
	totalCreates := len(implementDates) * len(pairingIDs)
	errCh := make(chan error, totalCreates)
	for _, implementDate := range implementDates {
		implementDate := implementDate
		for _, pairingID := range pairingIDs {
			wg.Add(1)
			pairingID := pairingID
			go func(pairingID uint, ec chan<- error) {
				defer wg.Done()
				meeting, err := uc.meetingRepo.Create(&entity.Meeting{
					ImplementedAt:   implementedAtPtr(implementDate),
					Status:          valueobject.MeetingStatusPending,
					TemplateID:      templateID,
					PairingID:       pairingID,
					CreatedByUserID: createdByUserID,
				})
				if err != nil {
					ec <- fmt.Errorf("定期ミーティングのペアリングのミーティング作成に失敗しました: %w", err)
					return
				}
				if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
					ec <- err
					return
				}
				mux.Lock()
				result = append(result, *meeting)
				mux.Unlock()
			}(pairingID, errCh)
		}
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return converter.MeetingEntitiesToOutputs(result), nil
}

func (uc *MeetingCreator) CreateMeetingEveryMonthWeekday(mentorID, templateID uint, menteeIDs []uint, from, to time.Time, week int, weekday time.Weekday, createdByUserID uint, sharerUserIDs []uint) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreateMeetingEveryMonthWeekday"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if mentorID == 0 {
		return nil, fmt.Errorf("メンターが指定されていません")
	}
	if len(menteeIDs) == 0 {
		return nil, fmt.Errorf("メンティが指定されていません")
	}
	implementDates := utils.ListSpecificWeekdays(from, to, week, weekday)
	if len(implementDates) == 0 {
		return converter.MeetingEntitiesToOutputs(nil), nil
	}
	pairingIDs := make([]uint, len(menteeIDs))
	for idx, menteeID := range menteeIDs {
		if menteeID == 0 {
			return nil, fmt.Errorf("無効なメンティが含まれています")
		}
		if menteeID == mentorID {
			return nil, fmt.Errorf("メンターと同じユーザをメンティに指定できません")
		}
		pair, err := uc.pairingRepo.CreateTemporaryPairing(menteeID, mentorID, createdByUserID)
		if err != nil {
			return nil, err
		}
		pairingIDs[idx] = pair.ID
	}
	var (
		result = []entity.Meeting{}
		mux    = sync.Mutex{}
		wg     = sync.WaitGroup{}
	)
	totalCreates := len(implementDates) * len(pairingIDs)
	errCh := make(chan error, totalCreates)
	for _, implementDate := range implementDates {
		implementDate := implementDate
		for _, pairingID := range pairingIDs {
			wg.Add(1)
			pairingID := pairingID
			go func(pairingID uint, ec chan<- error) {
				defer wg.Done()
				meeting, err := uc.meetingRepo.Create(&entity.Meeting{
					ImplementedAt:   implementedAtPtr(implementDate),
					Status:          valueobject.MeetingStatusPending,
					TemplateID:      templateID,
					PairingID:       pairingID,
					CreatedByUserID: createdByUserID,
				})
				if err != nil {
					ec <- fmt.Errorf("定期ミーティングのペアリングのミーティング作成に失敗しました: %w", err)
					return
				}
				if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
					ec <- err
					return
				}
				mux.Lock()
				result = append(result, *meeting)
				mux.Unlock()
			}(pairingID, errCh)
		}
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return converter.MeetingEntitiesToOutputs(result), nil
}

func (uc *MeetingCreator) CreateMeetingEveryMonthByPairings(templateID uint, from, to time.Time, week int, weekday time.Weekday, createdByUserID uint, sharerUserIDs []uint, pairingIDs []uint) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreateMeetingEveryMonthByPairings"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if len(pairingIDs) == 0 {
		return nil, fmt.Errorf("ペアが指定されていません")
	}
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ミーティング作成時に全てのペアリング情報の取得に失敗しました: %w", err)
	}
	pairings = filterPairingsByIDs(pairings, pairingIDs)
	if len(pairings) == 0 {
		return nil, fmt.Errorf("ペアがありません")
	}
	implementDates := utils.ListSpecificWeekdays(from, to, week, weekday)
	if len(implementDates) == 0 {
		return converter.MeetingEntitiesToOutputs(nil), nil
	}
	var (
		result = []entity.Meeting{}
		mux    = sync.Mutex{}
		wg     = sync.WaitGroup{}
	)
	totalCreates := len(implementDates) * len(pairings)
	errCh := make(chan error, totalCreates)
	for _, implementDate := range implementDates {
		implementDate := implementDate
		for _, pairing := range pairings {
			wg.Add(1)
			pairingID := pairing.ID
			go func(pairingID uint, ec chan<- error) {
				defer wg.Done()
				meeting, err := uc.meetingRepo.Create(&entity.Meeting{
					ImplementedAt:   implementedAtPtr(implementDate),
					Status:          valueobject.MeetingStatusPending,
					TemplateID:      templateID,
					PairingID:       pairingID,
					CreatedByUserID: createdByUserID,
				})
				if err != nil {
					ec <- fmt.Errorf("定期ミーティングのペアリングのミーティング作成に失敗しました: %w", err)
					return
				}
				if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
					ec <- err
					return
				}
				mux.Lock()
				result = append(result, *meeting)
				mux.Unlock()
			}(pairingID, errCh)
		}
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return converter.MeetingEntitiesToOutputs(result), nil
}

func (uc *MeetingCreator) CreateMeeting(pairingID, templateID uint, implementDate time.Time, createdByUserID uint, sharerUserIDs []uint) (ret0 *outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreateMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	meeting, err := uc.meetingRepo.Create(&entity.Meeting{
		Status:          valueobject.MeetingStatusPending,
		ImplementedAt:   implementedAtPtr(implementDate),
		TemplateID:      templateID,
		PairingID:       pairingID,
		CreatedByUserID: createdByUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("ミーティングの作成に失敗しました: %w", err)
	}
	if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
		return nil, err
	}
	return converter.MeetingEntityToOutputPtr(meeting), nil
}

func (uc *MeetingCreator) CreateMeetingByAll(templateID uint, implementDate time.Time, createdByUserID uint, sharerUserIDs []uint, pairingIDs []uint) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "meeting_mng_usecase.CreateMeetingByAll"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ミーティング作成時に全てのペアリング情報の取得に失敗しました: %w", err)
	}
	pairings = filterPairingsByIDs(pairings, pairingIDs)
	var (
		result = []entity.Meeting{}
		mux    = sync.Mutex{}
		wg     = sync.WaitGroup{}
	)
	errCh := make(chan error, len(pairings))
	for _, pairing := range pairings {
		wg.Add(1)
		pairing := pairing
		go func(pairingID uint, ec chan<- error) {
			defer wg.Done()
			meeting, err := uc.meetingRepo.Create(&entity.Meeting{
				ImplementedAt:   implementedAtPtr(implementDate),
				Status:          valueobject.MeetingStatusPending,
				TemplateID:      templateID,
				PairingID:       pairingID,
				CreatedByUserID: createdByUserID,
			})
			if err != nil {
				ec <- fmt.Errorf("全ペアリングのミーティング作成に失敗しました: %w", err)
				return
			}
			if err = uc.saveSharers(meeting.ID, sharerUserIDs); err != nil {
				ec <- err
				return
			}
			mux.Lock()
			result = append(result, *meeting)
			mux.Unlock()
		}(pairing.ID, errCh)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return converter.MeetingEntitiesToOutputs(result), nil
}
