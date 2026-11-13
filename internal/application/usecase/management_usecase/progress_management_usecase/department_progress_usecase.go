package progress_mng_usecase

import (
	"fmt"
	"slices"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

type deptProgressStats struct {
	userCount      int
	pairCount      int
	notDoneCount   int
	submittedCount int
	doneCount      int
}

func (uc *ProgressViewer) visibleUserIDs(filter ProgressFilter) (map[uint]bool, error) {
	if uc.fullAccess(filter.SystemRole) {
		return nil, nil
	}
	visible, err := uc.visibleMenteeIDs(filter)
	if err != nil {
		return nil, err
	}
	if visible == nil {
		return nil, nil
	}
	if len(visible) == 0 && filter.ViewerID > 0 {
		return map[uint]bool{filter.ViewerID: true}, nil
	}
	return visible, nil
}

func intersectUsers(users map[uint]struct{}, visible map[uint]bool) map[uint]struct{} {
	if visible == nil {
		return users
	}
	out := map[uint]struct{}{}
	for id := range users {
		if visible[id] {
			out[id] = struct{}{}
		}
	}
	return out
}

func descendantDepartmentIDsLocal(rootIDs []uint, hierarchies []entity.DepartmentHierarchy) []uint {
	children := map[uint][]uint{}
	for _, hierarchy := range hierarchies {
		children[hierarchy.ParentDepartmentID] = append(children[hierarchy.ParentDepartmentID], hierarchy.ChildDepartmentID)
	}
	seen := map[uint]bool{}
	var walk func(uint)
	walk = func(id uint) {
		if seen[id] {
			return
		}
		seen[id] = true
		for _, childID := range children[id] {
			walk(childID)
		}
	}
	for _, rootID := range rootIDs {
		walk(rootID)
	}
	ids := make([]uint, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func chosenParentByChild(hierarchies []entity.DepartmentHierarchy) map[uint]uint {
	out := map[uint]uint{}
	for _, h := range hierarchies {
		childID := h.ChildDepartmentID
		parentID := h.ParentDepartmentID
		if existing, ok := out[childID]; !ok || parentID < existing {
			out[childID] = parentID
		}
	}
	return out
}

func departmentDisplayOrder(depts []entity.Department, parentByChild map[uint]uint) ([]entity.Department, []int) {
	byID := map[uint]entity.Department{}
	for _, d := range depts {
		byID[d.ID] = d
	}
	childIDs := map[uint]bool{}
	for childID := range parentByChild {
		childIDs[childID] = true
	}
	var roots []entity.Department
	for _, d := range depts {
		if !childIDs[d.ID] {
			roots = append(roots, d)
		}
	}
	slices.SortFunc(roots, func(a, b entity.Department) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	children := map[uint][]uint{}
	for childID, parentID := range parentByChild {
		children[parentID] = append(children[parentID], childID)
	}
	for parentID := range children {
		slices.SortFunc(children[parentID], func(a, b uint) int {
			na, nb := byID[a].Name, byID[b].Name
			if na < nb {
				return -1
			}
			if na > nb {
				return 1
			}
			return 0
		})
	}

	var ordered []entity.Department
	var depths []int
	var walk func(deptID uint, depth int)
	walk = func(deptID uint, depth int) {
		d, ok := byID[deptID]
		if !ok {
			return
		}
		ordered = append(ordered, d)
		depths = append(depths, depth)
		for _, childID := range children[deptID] {
			walk(childID, depth+1)
		}
	}
	for _, root := range roots {
		walk(root.ID, 0)
	}
	return ordered, depths
}

func usersInDepartmentSubtree(
	deptID uint,
	hierarchies []entity.DepartmentHierarchy,
	membersByDept map[uint][]entity.DepartmentMember,
	teamsByDept map[uint][]entity.Team,
	teamMembersByTeam map[uint][]entity.TeamMember,
) map[uint]struct{} {
	subtree := descendantDepartmentIDsLocal([]uint{deptID}, hierarchies)
	users := map[uint]struct{}{}
	for _, id := range subtree {
		for _, m := range membersByDept[id] {
			if m.UserID > 0 {
				users[m.UserID] = struct{}{}
			}
		}
		for _, team := range teamsByDept[id] {
			for _, tm := range teamMembersByTeam[team.ID] {
				if tm.MemberID > 0 {
					users[tm.MemberID] = struct{}{}
				}
			}
		}
	}
	return users
}

func allAffiliatedUserIDs(
	members []entity.DepartmentMember,
	teamMembers []entity.TeamMember,
) map[uint]struct{} {
	users := map[uint]struct{}{}
	for _, m := range members {
		if m.UserID > 0 {
			users[m.UserID] = struct{}{}
		}
	}
	for _, tm := range teamMembers {
		if tm.MemberID > 0 {
			users[tm.MemberID] = struct{}{}
		}
	}
	return users
}

func meetingInDateRange(m entity.Meeting, filter ProgressFilter) bool {
	if filter.From != nil && m.ImplementedAt != nil && m.ImplementedAt.Before(*filter.From) {
		return false
	}
	if filter.To != nil && m.ImplementedAt != nil && m.ImplementedAt.After(*filter.To) {
		return false
	}
	return true
}

func statsForUsers(
	userIDs map[uint]struct{},
	pairings []entity.Pairing,
	meetings []entity.Meeting,
	filter ProgressFilter,
) deptProgressStats {
	menteeInSet := map[uint]bool{}
	for id := range userIDs {
		menteeInSet[id] = true
	}
	pairCount := 0
	for _, p := range pairings {
		if menteeInSet[p.MenteeID] {
			pairCount++
		}
	}
	var notDone, submitted, done int
	for _, m := range meetings {
		if !meetingInDateRange(m, filter) {
			continue
		}
		menteeID := uint(0)
		if m.Pairing != nil {
			menteeID = m.Pairing.MenteeID
		}
		if menteeID == 0 || !menteeInSet[menteeID] {
			continue
		}
		switch m.Status {
		case valueobject.MeetingStatusPending, valueobject.MeetingStatusDraft:
			notDone++
		case valueobject.MeetingStatusSubmitted:
			submitted++
		case valueobject.MeetingStatusDone:
			done++
		}
	}
	return deptProgressStats{
		userCount:      len(userIDs),
		pairCount:      pairCount,
		notDoneCount:   notDone,
		submittedCount: submitted,
		doneCount:      done,
	}
}

func progressPercents(s deptProgressStats) (subPct, donePct float64) {
	den := s.notDoneCount + s.submittedCount + s.doneCount
	if den > 0 {
		subPct = float64(s.submittedCount) / float64(den) * 100
		donePct = float64(s.doneCount) / float64(den) * 100
	}
	return subPct, donePct
}

func rowFromStats(name, uuid string, depth int, isTotal bool, s deptProgressStats, teams []outputmodel.ProgressTeamRowOutput) outputmodel.ProgressDepartmentRowOutput {
	subPct, donePct := progressPercents(s)
	if teams == nil {
		teams = []outputmodel.ProgressTeamRowOutput{}
	}
	return outputmodel.ProgressDepartmentRowOutput{
		ProgressRowOutput: outputmodel.ProgressRowOutput{
			UUID:             uuid,
			Name:             name,
			Depth:            depth,
			IsTotal:          isTotal,
			UserCount:        s.userCount,
			PairCount:        s.pairCount,
			NotDoneCount:     s.notDoneCount,
			SubmittedCount:   s.submittedCount,
			SubmittedPercent: subPct,
			DoneCount:        s.doneCount,
			DonePercent:      donePct,
		},
		Teams: teams,
	}
}

func usersInTeam(teamID uint, teamMembersByTeam map[uint][]entity.TeamMember) map[uint]struct{} {
	users := map[uint]struct{}{}
	for _, tm := range teamMembersByTeam[teamID] {
		if tm.MemberID > 0 {
			users[tm.MemberID] = struct{}{}
		}
	}
	return users
}

func teamRowFromStats(team entity.Team, departmentUUID string, s deptProgressStats) outputmodel.ProgressTeamRowOutput {
	subPct, donePct := progressPercents(s)
	deptUUIDs := []string{}
	if departmentUUID != "" {
		deptUUIDs = []string{departmentUUID}
	}
	return outputmodel.ProgressTeamRowOutput{
		ProgressRowOutput: outputmodel.ProgressRowOutput{
			UUID:             team.UUID.String(),
			Name:             team.Name,
			Depth:            0,
			IsTotal:          false,
			UserCount:        s.userCount,
			PairCount:        s.pairCount,
			NotDoneCount:     s.notDoneCount,
			SubmittedCount:   s.submittedCount,
			SubmittedPercent: subPct,
			DoneCount:        s.doneCount,
			DonePercent:      donePct,
		},
		DepartmentUUIDs: deptUUIDs,
	}
}

func buildTeamRowsForDepartment(
	dept entity.Department,
	deptTeams []entity.Team,
	visible map[uint]bool,
	teamMembersByTeam map[uint][]entity.TeamMember,
	pairings []entity.Pairing,
	meetings []entity.Meeting,
	filter ProgressFilter,
) []outputmodel.ProgressTeamRowOutput {
	sorted := append([]entity.Team(nil), deptTeams...)
	slices.SortFunc(sorted, func(a, b entity.Team) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
	out := make([]outputmodel.ProgressTeamRowOutput, 0, len(sorted))
	deptUUID := dept.UUID.String()
	for _, team := range sorted {
		users := intersectUsers(usersInTeam(team.ID, teamMembersByTeam), visible)
		if len(users) == 0 {
			continue
		}
		stats := statsForUsers(users, pairings, meetings, filter)
		out = append(out, teamRowFromStats(team, deptUUID, stats))
	}
	return out
}

func (uc *ProgressViewer) DepartmentTable(filter ProgressFilter) (*outputmodel.ProgressDepartmentTableOutput, error) {
	if uc.deptRepo == nil || uc.hierarchyRepo == nil || uc.deptMemberRepo == nil ||
		uc.teamRepo == nil || uc.teamMemberRepo == nil || uc.pairingRepo == nil {
		return nil, fmt.Errorf("部署進捗の集計に必要なリポジトリが設定されていません")
	}

	visible, err := uc.visibleUserIDs(filter)
	if err != nil {
		return nil, err
	}

	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	deptMembers, err := uc.deptMemberRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamMembers, err := uc.teamMemberRepo.GetAll()
	if err != nil {
		return nil, err
	}
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, err
	}
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}

	membersByDept := map[uint][]entity.DepartmentMember{}
	for _, m := range deptMembers {
		membersByDept[m.DepartmentID] = append(membersByDept[m.DepartmentID], m)
	}

	teamsByDept := map[uint][]entity.Team{}
	for _, t := range teams {
		teamsByDept[t.DepartmentID] = append(teamsByDept[t.DepartmentID], t)
	}

	teamMembersByTeam := map[uint][]entity.TeamMember{}
	for _, tm := range teamMembers {
		teamMembersByTeam[tm.TeamID] = append(teamMembersByTeam[tm.TeamID], tm)
	}

	affiliated := allAffiliatedUserIDs(deptMembers, teamMembers)
	totalUsers := intersectUsers(affiliated, visible)
	rows := make([]outputmodel.ProgressDepartmentRowOutput, 0, len(depts)+1)

	if len(totalUsers) > 0 || visible == nil {
		totalStats := statsForUsers(totalUsers, pairings, meetings, filter)
		rows = append(rows, rowFromStats("合計", "", 0, true, totalStats, []outputmodel.ProgressTeamRowOutput{}))
	}

	parentByChild := chosenParentByChild(hierarchies)
	ordered, depths := departmentDisplayOrder(depts, parentByChild)
	for i, dept := range ordered {
		rawUsers := usersInDepartmentSubtree(
			dept.ID,
			hierarchies,
			membersByDept,
			teamsByDept,
			teamMembersByTeam,
		)
		users := intersectUsers(rawUsers, visible)
		if len(users) == 0 {
			continue
		}
		stats := statsForUsers(users, pairings, meetings, filter)
		teamRows := buildTeamRowsForDepartment(
			dept,
			teamsByDept[dept.ID],
			visible,
			teamMembersByTeam,
			pairings,
			meetings,
			filter,
		)
		rows = append(rows, rowFromStats(dept.Name, dept.UUID.String(), depths[i], false, stats, teamRows))
	}

	return &outputmodel.ProgressDepartmentTableOutput{Rows: rows}, nil
}
