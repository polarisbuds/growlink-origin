package repository

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"

	"gorm.io/gorm"
)

func meetingNeedsPairingFilter(f repository_interface.MeetingListFilter) bool {
	if len(f.MentorUUIDs) > 0 || len(f.MenteeUUIDs) > 0 || len(f.DepartmentUUIDs) > 0 {
		return true
	}
	return !f.IncludeDeleted
}

func (repo *MeetingRepository) departmentMemberUserIDsSubquery(deptUUIDs []string) gorm.ChainInterface[entity.DepartmentMember] {
	deptIDs := gorm.G[entity.Department](repo.db).Select("id").Where("uuid IN ?", deptUUIDs)
	return gorm.G[entity.DepartmentMember](repo.db).Select("user_id").Where("department_id IN (?)", deptIDs)
}

func (repo *MeetingRepository) pairingListSubquery(f repository_interface.MeetingListFilter) gorm.ChainInterface[entity.Pairing] {
	var chain gorm.ChainInterface[entity.Pairing]
	if f.IncludeDeleted {
		chain = gorm.G[entity.Pairing](repo.unscoped).Select("id").
			Where("is_temporary = ?", false).
			Scopes(func(stmt *gorm.Statement) { stmt.Unscoped = true })
	} else {
		chain = gorm.G[entity.Pairing](repo.db).Select("id").Where("is_temporary = ?", false)
	}
	if len(f.MentorUUIDs) > 0 {
		mentors := gorm.G[entity.User](repo.db).Select("id").Where("uuid IN ?", f.MentorUUIDs)
		chain = chain.Where("mentor_id IN (?)", mentors)
	}
	if len(f.MenteeUUIDs) > 0 {
		mentees := gorm.G[entity.User](repo.db).Select("id").Where("uuid IN ?", f.MenteeUUIDs)
		chain = chain.Where("mentee_id IN (?)", mentees)
	}
	if len(f.DepartmentUUIDs) > 0 {
		byMentor := repo.departmentMemberUserIDsSubquery(f.DepartmentUUIDs)
		byMentee := repo.departmentMemberUserIDsSubquery(f.DepartmentUUIDs)
		chain = chain.Where("(mentor_id IN (?) OR mentee_id IN (?))", byMentor, byMentee)
	}
	return chain
}

func applyMeetingListFilter(
	chain gorm.ChainInterface[entity.Meeting],
	f repository_interface.MeetingListFilter,
	subQuery gorm.ChainInterface[entity.Pairing],
) gorm.ChainInterface[entity.Meeting] {
	if f.LimitToMeetingIDs {
		chain = chain.Where("id IN ?", f.MeetingIDs)
	}
	if subQuery != nil {
		chain = chain.Where("pairing_id IN (?)", subQuery)
	}
	if f.Status != nil {
		chain = chain.Where("status = ?", *f.Status)
	}
	hasDate := f.From != nil || f.To != nil
	if hasDate {
		chain = chain.Where("implemented_at IS NOT NULL")
		if f.From != nil {
			chain = chain.Where("implemented_at >= ?", *f.From)
		}
		if f.To != nil {
			chain = chain.Where("implemented_at <= ?", *f.To)
		}
	}
	return chain
}

func (repo *MeetingRepository) GetByFilter(f repository_interface.MeetingListFilter) (ret0 []entity.Meeting, err error) {
	if f.LimitToMeetingIDs && len(f.MeetingIDs) == 0 {
		return []entity.Meeting{}, nil
	}

	var subq gorm.ChainInterface[entity.Pairing]
	if meetingNeedsPairingFilter(f) {
		subq = repo.pairingListSubquery(f)
	}

	base := repo.db
	if f.IncludeDeleted {
		base = repo.unscoped
	}
	chain := preloadBasic(gorm.G[entity.Meeting](base).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Order("id DESC"))
	chain = applyMeetingListFilter(chain, f, subq)

	meetings, err := chain.Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("フィルタ済み面談取得に失敗しました: %w", err)
	}
	return meetings, nil
}
