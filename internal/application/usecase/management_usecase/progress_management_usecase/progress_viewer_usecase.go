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
	meetingRepo      repository_interface.MeetingRepository
	userRepo         repository_interface.UserRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
	userDeptRepo     repository_interface.UserDepartmentRepository
	deptRefRepo      repository_interface.UserReferenceDepartmentRepository
	teamRefRepo      repository_interface.UserReferenceTeamRepository
	relationshipRepo repository_interface.UserRelationshipRepository
}

func NewProgressViewer(
	meetingRepo repository_interface.MeetingRepository,
	userRepo repository_interface.UserRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	userDeptRepo repository_interface.UserDepartmentRepository,
	deptRefRepo repository_interface.UserReferenceDepartmentRepository,
	teamRefRepo repository_interface.UserReferenceTeamRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
) *ProgressViewer {
	return &ProgressViewer{
		meetingRepo:      meetingRepo,
		userRepo:         userRepo,
		teamMemberRepo:   teamMemberRepo,
		userDeptRepo:     userDeptRepo,
		deptRefRepo:      deptRefRepo,
		teamRefRepo:      teamRefRepo,
		relationshipRepo: relationshipRepo,
	}
}

func (uc *ProgressViewer) fullAccess(role valueobject.SystemRole) bool {
	return role == valueobject.SystemRoleAdmin || role == valueobject.SystemRoleAnalyst
}

func (uc *ProgressViewer) visibleMenteeIDs(filter ProgressFilter) (map[uint]bool, error) {
	if uc.fullAccess(filter.SystemRole) {
		return nil, nil
	}
	user, err := uc.userRepo.GetByID(filter.ViewerID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return map[uint]bool{}, nil
	}
	visible := map[uint]bool{}
	switch user.OperationTarget {
	case valueobject.OperationTargetDepartment:
		return uc.visibleDepartmentMentees(filter.ViewerID, visible)
	case valueobject.OperationTargetTeam:
		return uc.visibleTeamMentees(filter.ViewerID, visible)
	case valueobject.OperationTargetSubordinate:
		return uc.visibleSubordinateMentees(filter.ViewerID, visible)
	default:
		return visible, nil
	}
}

func (uc *ProgressViewer) visibleDepartmentMentees(viewerID uint, visible map[uint]bool) (map[uint]bool, error) {
	refs, err := uc.deptRefRepo.GetByUserID(viewerID)
	if err != nil {
		return nil, err
	}
	deptIDs := map[uint]bool{}
	for _, ref := range refs {
		if ref.ReferableDepartmentID > 0 {
			deptIDs[ref.ReferableDepartmentID] = true
		}
	}
	if len(deptIDs) == 0 {
		return visible, nil
	}
	allUsers, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, u := range allUsers {
		if deptIDs[u.DepartmentID] {
			visible[u.ID] = true
		}
	}
	for deptID := range deptIDs {
		ids, err := uc.userDeptRepo.GetUserIDsByDepartmentID(deptID)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			if id > 0 {
				visible[id] = true
			}
		}
	}
	return visible, nil
}

func (uc *ProgressViewer) visibleTeamMentees(viewerID uint, visible map[uint]bool) (map[uint]bool, error) {
	refs, err := uc.teamRefRepo.GetByUserID(viewerID)
	if err != nil {
		return nil, err
	}
	for _, ref := range refs {
		if ref.ReferableTeamID == 0 {
			continue
		}
		members, err := uc.teamMemberRepo.GetByTeamID(ref.ReferableTeamID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			if member.MemberID > 0 {
				visible[member.MemberID] = true
			}
		}
	}
	return visible, nil
}

func (uc *ProgressViewer) visibleSubordinateMentees(viewerID uint, visible map[uint]bool) (map[uint]bool, error) {
	rels, err := uc.relationshipRepo.GetByBossID(viewerID)
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		if rel.SourceUserID > 0 {
			visible[rel.SourceUserID] = true
		}
	}
	return visible, nil
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
		if filter.To != nil {
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
	if out.Total > 0 {
		out.DonePercent = float64(out.Done) / float64(out.Total) * 100
	}
	return out, nil
}

type progressBucket struct {
	id             uint
	name           string
	departmentName string
	pending        int
	draft          int
	done           int
	usersPending   []string
	usersDraft     []string
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
	addCount := func(id uint, name, deptName string, status valueobject.MeetingStatus, userName string) {
		b, ok := buckets[id]
		if !ok {
			b = &progressBucket{id: id, name: name, departmentName: deptName}
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
				addCount(0, unassignedProgressName, "", m.Status, uname)
				continue
			}
			if mentee.Department != nil && mentee.Department.Name != "" {
				addCount(mentee.Department.ID, mentee.Department.Name, "", m.Status, uname)
				continue
			}
			if mentee.DepartmentID > 0 {
				addCount(mentee.DepartmentID, fmt.Sprintf("部署 #%d", mentee.DepartmentID), "", m.Status, uname)
				continue
			}
			addCount(0, unassignedProgressName, "", m.Status, uname)
			continue
		}

		if mentee == nil || len(mentee.TeamMembers) == 0 {
			addCount(0, unassignedProgressName, "", m.Status, uname)
			continue
		}
		seenTeam := map[uint]bool{}
		for _, tm := range mentee.TeamMembers {
			var teamID uint
			var teamName string
			var deptName string
			if tm.Team != nil {
				teamID = tm.Team.ID
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
			addCount(teamID, display, deptName, m.Status, uname)
		}
		if len(seenTeam) == 0 {
			addCount(0, unassignedProgressName, "", m.Status, uname)
		}
	}

	out := make([]outputmodel.ProgressBreakdownRowOutput, 0, len(buckets))
	for _, b := range buckets {
		total := b.pending + b.draft + b.done
		var donePct float64
		if total > 0 {
			donePct = float64(b.done) / float64(total) * 100
		}
		out = append(out, outputmodel.ProgressBreakdownRowOutput{
			ID:             b.id,
			Name:           b.name,
			DepartmentName: b.departmentName,
			Pending:        b.pending,
			Draft:          b.draft,
			Done:           b.done,
			DonePercent:    donePct,
			UsersPending:   b.usersPending,
			UsersDraft:     b.usersDraft,
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
