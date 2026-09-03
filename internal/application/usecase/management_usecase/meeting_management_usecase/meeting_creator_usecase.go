package meeting_mng_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/utils"
	"sync"
	"time"
)

type MeetingCreatorUsecase struct {
	meetingRepo repository_interface.MeetingRepository
	pairingRepo repository_interface.PairingRepository
}

func NewMeetingCreatorUsecase(meetingRepo repository_interface.MeetingRepository, pairingRepo repository_interface.PairingRepository) *MeetingCreatorUsecase {
	return &MeetingCreatorUsecase{
		meetingRepo: meetingRepo,
		pairingRepo: pairingRepo,
	}
}

// 次の実施日を計算します
func calcImplementedAt(base time.Time) time.Time {
	nextImplementDay := base
	for utils.IsHoliday(nextImplementDay) {
		nextImplementDay = nextImplementDay.AddDate(0, 0, 1)
	}
	return nextImplementDay
}

func (uc *MeetingCreatorUsecase) CreateMeetingEveryMonth(templateID uint, baseDate time.Time, week int, weekday time.Weekday) ([]entity.Meeting, error) {
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("定期ミーティング作成処理時に全てのペアリング情報の取得に失敗しました: %w", err)
	}

	var (
		result = []entity.Meeting{}
		mux    = sync.Mutex{}
		wg     = sync.WaitGroup{}
	)
	errCh := make(chan error, len(pairings))

	implementDate := utils.GetSpecificWeekday(baseDate, 2, time.Wednesday)
	for _, pairing := range pairings {
		wg.Add(1)
		pairing := pairing

		go func(pairingID uint, ec chan<- error) {
			defer wg.Done()
			meeting, err := uc.meetingRepo.Create(&entity.Meeting{
				ImplementedAt: calcImplementedAt(implementDate),
				Status:        valueobject.MeetingStatusPending,
				TemplateID:    templateID,
				PairingID:     pairingID,
			})
			if err != nil {
				ec <- fmt.Errorf("定期ミーティングのペアリングのミーティング作成に失敗しました: %w", err)
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

	return result, nil
}

func (uc *MeetingCreatorUsecase) CreateMeeting(pairingID, templateID uint, implementDate time.Time) (*entity.Meeting, error) {
	meeting, err := uc.meetingRepo.Create(&entity.Meeting{
		Status:        valueobject.MeetingStatusPending,
		ImplementedAt: calcImplementedAt(implementDate),
		TemplateID:    templateID,
		PairingID:     pairingID,
	})
	if err != nil {
		return nil, fmt.Errorf("ミーティングの作成に失敗しました: %w", err)
	}
	return meeting, nil
}

func (uc *MeetingCreatorUsecase) CreateMeetingByAll(templateID uint, implementDate time.Time) ([]entity.Meeting, error) {
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ミーティング作成時に全てのペアリング情報の取得に失敗しました: %w", err)
	}

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
				ImplementedAt: calcImplementedAt(implementDate),
				Status:        valueobject.MeetingStatusPending,
				TemplateID:    templateID,
				PairingID:     pairingID,
			})
			if err != nil {
				ec <- fmt.Errorf("全ペアリングのミーティング作成に失敗しました: %w", err)
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

	return result, nil
}
