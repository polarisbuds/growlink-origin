package meeting_mng_usecase

import (
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/utils"
	"sync"
	"time"
)

type MeetingCreator struct {
	meetingRepo repository_interface.MeetingRepository
	pairingRepo repository_interface.PairingRepository
}

func NewMeetingCreator(meetingRepo repository_interface.MeetingRepository, pairingRepo repository_interface.PairingRepository) *MeetingCreator {
	return &MeetingCreator{
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

func implementedAtPtr(base time.Time) *time.Time {
	at := calcImplementedAt(base)
	return &at
}

func (uc *MeetingCreator) CreateMeetingEveryMonth(templateID uint, baseDate time.Time, week int, weekday time.Weekday) ([]outputmodel.MeetingOutput, error) {
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
	implementDate := utils.GetSpecificWeekday(baseDate, week, weekday)
	for _, pairing := range pairings {
		wg.Add(1)
		pairing := pairing
		go func(pairingID uint, ec chan<- error) {
			defer wg.Done()
			meeting, err := uc.meetingRepo.Create(&entity.Meeting{
				ImplementedAt: implementedAtPtr(implementDate),
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
	return converter.MeetingEntitiesToOutputs(result), nil
}

func (uc *MeetingCreator) CreateMeeting(pairingID, templateID uint, implementDate time.Time) (*outputmodel.MeetingOutput, error) {
	meeting, err := uc.meetingRepo.Create(&entity.Meeting{
		Status:        valueobject.MeetingStatusPending,
		ImplementedAt: implementedAtPtr(implementDate),
		TemplateID:    templateID,
		PairingID:     pairingID,
	})
	if err != nil {
		return nil, fmt.Errorf("ミーティングの作成に失敗しました: %w", err)
	}
	return converter.MeetingEntityToOutputPtr(meeting), nil
}

func (uc *MeetingCreator) CreateMeetingByAll(templateID uint, implementDate time.Time) ([]outputmodel.MeetingOutput, error) {
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
				ImplementedAt: implementedAtPtr(implementDate),
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
	return converter.MeetingEntitiesToOutputs(result), nil
}
