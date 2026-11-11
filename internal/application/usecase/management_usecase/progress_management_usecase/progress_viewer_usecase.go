package progress_mng_usecase

import (
	"fmt"
	"slices"
	"time"

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

type ProgressFilter struct {
	From       *time.Time
	To         *time.Time
	ViewerID   uint
	SystemRole valueobject.SystemRole
}

type ProgressViewer struct {
	meetingRepo repository_interface.MeetingRepository
	referable   ReferableUsers
}

func NewProgressViewer(
	meetingRepo repository_interface.MeetingRepository,
	referable ReferableUsers,
) *ProgressViewer {
	return &ProgressViewer{
		meetingRepo: meetingRepo,
		referable:   referable,
	}
}

func (uc *ProgressViewer) fullAccess(role valueobject.SystemRole) bool {
	return role == valueobject.SystemRoleAdmin || role == valueobject.SystemRoleAnalyst
}

func (uc *ProgressViewer) visibleMenteeIDs(filter ProgressFilter) (map[uint]bool, error) {
	if uc.fullAccess(filter.SystemRole) {
		return nil, nil
	}
	if uc.referable == nil {
		return map[uint]bool{}, nil
	}
	return uc.referable.ReferableUserIDs(filter.ViewerID, filter.SystemRole)
}

func (uc *ProgressViewer) filterMeetings(meetings []entity.Meeting, filter ProgressFilter) ([]entity.Meeting, error) {
	visibleMentees, err := uc.visibleMenteeIDs(filter)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Meeting, 0, len(meetings))
	for _, m := range meetings {
		if filter.From != nil && m.ImplementedAt != nil && m.ImplementedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && m.ImplementedAt != nil {
			end := filter.To
			if m.ImplementedAt.After(*end) {
				continue
			}
		}
		if visibleMentees != nil {
			menteeID := uint(0)
			if m.Pairing != nil {
				menteeID = m.Pairing.MenteeID
			}
			if menteeID == 0 || !visibleMentees[menteeID] {
				continue
			}
		}
		out = append(out, m)
	}
	return out, nil
}

func (uc *ProgressViewer) Summary(filter ProgressFilter) (*outputmodel.ProgressSummaryOutput, error) {
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}
	meetings, err = uc.filterMeetings(meetings, filter)
	if err != nil {
		return nil, err
	}
	// TODO; GORMを用いて取得するものへと変更する
	out := &outputmodel.ProgressSummaryOutput{}
	for _, m := range meetings {
		out.Total++
		switch m.Status {
		case valueobject.MeetingStatusPending:
			out.Pending++
		case valueobject.MeetingStatusDraft:
			out.Draft++
		case valueobject.MeetingStatusSubmitted:
			out.Submitted++
		case valueobject.MeetingStatusDone:
			out.Done++
		}
	}
	if out.Total > 0 {
		out.DonePercent = float64(out.Done) / float64(out.Total) * 100
	}
	return out, nil
}

type progressBucket struct {
	uuid           string
	name           string
	departmentName string
	pending        int
	draft          int
	submitted      int
	done           int
	usersPending   []string
	usersDraft     []string
	usersSubmitted []string
	usersDone      []string
}

func menteeName(m *entity.Meeting) string {
	if m.Pairing == nil || m.Pairing.Mentee == nil {
		return ""
	}
	return m.Pairing.Mentee.Name
}

func (uc *ProgressViewer) Breakdown(group string, filter ProgressFilter) ([]outputmodel.ProgressBreakdownRowOutput, error) {
	if group != progressGroupDepartment && group != progressGroupTeam {
		return nil, fmt.Errorf("不正な group です: %s", group)
	}
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}
	meetings, err = uc.filterMeetings(meetings, filter)
	if err != nil {
		return nil, err
	}

	buckets := map[uint]*progressBucket{}
	addCount := func(id uint, rowUUID, name, deptName string, status valueobject.MeetingStatus, userName string) {
		b, ok := buckets[id]
		if !ok {
			b = &progressBucket{uuid: rowUUID, name: name, departmentName: deptName}
			buckets[id] = b
		}
		switch status {
		case valueobject.MeetingStatusPending:
			b.pending++
			if userName != "" {
				b.usersPending = append(b.usersPending, userName)
			}
		case valueobject.MeetingStatusDraft:
			b.draft++
			if userName != "" {
				b.usersDraft = append(b.usersDraft, userName)
			}
		case valueobject.MeetingStatusSubmitted:
			b.submitted++
			if userName != "" {
				b.usersSubmitted = append(b.usersSubmitted, userName)
			}
		case valueobject.MeetingStatusDone:
			b.done++
			if userName != "" {
				b.usersDone = append(b.usersDone, userName)
			}
		}
	}

	for _, m := range meetings {
		mentee := (*entity.User)(nil)
		if m.Pairing != nil {
			mentee = m.Pairing.Mentee
		}
		uname := menteeName(&m)
		if group == progressGroupDepartment {
			if mentee == nil {
				addCount(0, "", unassignedProgressName, "", m.Status, uname)
				continue
			}
			if dept := primaryDepartmentMembership(mentee); dept != nil {
				name := fmt.Sprintf("部署 #%d", dept.DepartmentID)
				deptUUID := ""
				if dept.Department != nil {
					name = dept.Department.Name
					deptUUID = dept.Department.UUID.String()
				}
				addCount(dept.DepartmentID, deptUUID, name, "", m.Status, uname)
				continue
			}
			addCount(0, "", unassignedProgressName, "", m.Status, uname)
			continue
		}

		if mentee == nil || len(mentee.TeamMembers) == 0 {
			addCount(0, "", unassignedProgressName, "", m.Status, uname)
			continue
		}
		seenTeam := map[uint]bool{}
		for _, tm := range mentee.TeamMembers {
			var teamID uint
			var teamUUID string
			var teamName string
			var deptName string
			if tm.Team != nil {
				teamID = tm.Team.ID
				teamUUID = tm.Team.UUID.String()
				teamName = tm.Team.Name
				if tm.Team.Department != nil {
					deptName = tm.Team.Department.Name
				}
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
			display := teamName
			if deptName != "" {
				display = deptName + " / " + teamName
			}
			addCount(teamID, teamUUID, display, deptName, m.Status, uname)
		}
		if len(seenTeam) == 0 {
			addCount(0, "", unassignedProgressName, "", m.Status, uname)
		}
	}

	out := make([]outputmodel.ProgressBreakdownRowOutput, 0, len(buckets))
	for _, b := range buckets {
		total := b.pending + b.draft + b.submitted + b.done
		var donePct float64
		if total > 0 {
			donePct = float64(b.done) / float64(total) * 100
		}
		out = append(out, outputmodel.ProgressBreakdownRowOutput{
			UUID:           b.uuid,
			Name:           b.name,
			DepartmentName: b.departmentName,
			Pending:        b.pending,
			Draft:          b.draft,
			Submitted:      b.submitted,
			Done:           b.done,
			DonePercent:    donePct,
			UsersPending:   b.usersPending,
			UsersDraft:     b.usersDraft,
			UsersSubmitted: b.usersSubmitted,
			UsersDone:      b.usersDone,
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

func primaryDepartmentMembership(user *entity.User) *entity.DepartmentMember {
	if user == nil || len(user.DepartmentMembers) == 0 {
		return nil
	}
	primary := user.DepartmentMembers[0]
	for _, dm := range user.DepartmentMembers[1:] {
		if dm.ID < primary.ID {
			primary = dm
		}
	}
	return &primary
}
