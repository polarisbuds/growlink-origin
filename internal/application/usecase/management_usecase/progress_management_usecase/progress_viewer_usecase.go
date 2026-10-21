package progress_mng_usecase

import (
	"fmt"

	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type ProgressViewer struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewProgressViewer(meetingRepo repository_interface.MeetingRepository) *ProgressViewer {
	return &ProgressViewer{meetingRepo: meetingRepo}
}

func (uc *ProgressViewer) Summary() (*outputmodel.ProgressSummaryOutput, error) {
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}
	out := &outputmodel.ProgressSummaryOutput{}
	for _, m := range meetings {
		out.Total++
		switch m.Status {
		case valueobject.MeetingStatusPending:
			out.Pending++
		case valueobject.MeetingStatusDraft:
			out.Draft++
		case valueobject.MeetingStatusDone:
			out.Done++
		}
	}
	return out, nil
}
