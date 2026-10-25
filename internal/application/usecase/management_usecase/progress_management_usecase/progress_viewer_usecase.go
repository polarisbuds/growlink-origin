package progress_mng_usecase

import (
	"fmt"
	"slices"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

const (
	progressGroupDepartment = "department"
	progressGroupTeam       = "team"
	unassignedProgressName  = "未所属"
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

type progressBucket struct {
	id      uint
	name    string
	pending int
	draft   int
	done    int
}

func (uc *ProgressViewer) Breakdown(group string) ([]outputmodel.ProgressBreakdownRowOutput, error) {
	if group != progressGroupDepartment && group != progressGroupTeam {
		return nil, fmt.Errorf("不正な group です: %s", group)
	}
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}

	buckets := map[uint]*progressBucket{}
	addCount := func(id uint, name string, status valueobject.MeetingStatus) {
		b, ok := buckets[id]
		if !ok {
			b = &progressBucket{id: id, name: name}
			buckets[id] = b
		}
		switch status {
		case valueobject.MeetingStatusPending:
			b.pending++
		case valueobject.MeetingStatusDraft:
			b.draft++
		case valueobject.MeetingStatusDone:
			b.done++
		}
	}

	for _, m := range meetings {
		mentee := (*entity.User)(nil)
		if m.Pairing != nil {
			mentee = m.Pairing.Mentee
		}
		if group == progressGroupDepartment {
			if mentee == nil {
				addCount(0, unassignedProgressName, m.Status)
				continue
			}
			if mentee.Department != nil && mentee.Department.Name != "" {
				addCount(mentee.Department.ID, mentee.Department.Name, m.Status)
				continue
			}
			if mentee.DepartmentID > 0 {
				addCount(mentee.DepartmentID, fmt.Sprintf("部署 #%d", mentee.DepartmentID), m.Status)
				continue
			}
			addCount(0, unassignedProgressName, m.Status)
			continue
		}

		if mentee == nil || len(mentee.TeamMembers) == 0 {
			addCount(0, unassignedProgressName, m.Status)
			continue
		}
		seenTeam := map[uint]bool{}
		for _, tm := range mentee.TeamMembers {
			var teamID uint
			var teamName string
			if tm.Team != nil && tm.Team.Name != "" {
				teamID = tm.Team.ID
				teamName = tm.Team.Name
			} else if tm.TeamID > 0 {
				teamID = tm.TeamID
				teamName = fmt.Sprintf("チーム #%d", tm.TeamID)
			} else {
				continue
			}
			if seenTeam[teamID] {
				continue
			}
			seenTeam[teamID] = true
			addCount(teamID, teamName, m.Status)
		}
		if len(seenTeam) == 0 {
			addCount(0, unassignedProgressName, m.Status)
		}
	}

	out := make([]outputmodel.ProgressBreakdownRowOutput, 0, len(buckets))
	for _, b := range buckets {
		out = append(out, outputmodel.ProgressBreakdownRowOutput{
			ID:      b.id,
			Name:    b.name,
			Pending: b.pending,
			Draft:   b.draft,
			Done:    b.done,
		})
	}
	slices.SortFunc(out, func(a, b outputmodel.ProgressBreakdownRowOutput) int {
		if a.Name == unassignedProgressName && b.Name != unassignedProgressName {
			return 1
		}
		if b.Name == unassignedProgressName && a.Name != unassignedProgressName {
			return -1
		}
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
	return out, nil
}
