package infrastructure

import (
	"fmt"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/config"
	"growlink/internal/shared/crypto"

	"gorm.io/gorm"
)

// genSeed inserts development verification data when the database has no roles yet.
func genSeed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Role{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashed, err := crypto.HashPassword("password", config.GetPepper())
	if err != nil {
		return fmt.Errorf("genSeed: hash password: %w", err)
	}
	password := &hashed

	return db.Transaction(func(tx *gorm.DB) error {
		roleGeneral := &entity.Role{Name: "一般"}
		roleBoss := &entity.Role{Name: "上司"}
		roleOwner := &entity.Role{Name: "オーナー"}
		if err := tx.Create([]*entity.Role{roleGeneral, roleBoss, roleOwner}).Error; err != nil {
			return err
		}

		deptHQ := &entity.Department{Name: "本社"}
		deptDev := &entity.Department{Name: "開発部"}
		deptSales := &entity.Department{Name: "営業部"}
		if err := tx.Create([]*entity.Department{deptHQ, deptDev, deptSales}).Error; err != nil {
			return err
		}

		owner := &entity.User{
			Name:         "オーナー",
			Email:        "owner@example.com",
			Password:     password,
			RoleID:       roleOwner.ID,
			DepartmentID: deptHQ.ID,
			Authority:    valueobject.AuthorityOwner,
		}
		boss := &entity.User{
			Name:         "上司",
			Email:        "boss@example.com",
			Password:     password,
			RoleID:       roleBoss.ID,
			DepartmentID: deptDev.ID,
			Authority:    valueobject.AuthorityBoss,
		}
		mentor := &entity.User{
			Name:         "メンター",
			Email:        "mentor@example.com",
			Password:     password,
			RoleID:       roleGeneral.ID,
			DepartmentID: deptDev.ID,
			Authority:    valueobject.AuthorityGeneral,
		}
		mentee := &entity.User{
			Name:         "メンティ",
			Email:        "mentee@example.com",
			Password:     password,
			RoleID:       roleGeneral.ID,
			DepartmentID: deptDev.ID,
			Authority:    valueobject.AuthorityGeneral,
		}
		sharer := &entity.User{
			Name:         "共有者",
			Email:        "sharer@example.com",
			Password:     password,
			RoleID:       roleGeneral.ID,
			DepartmentID: deptSales.ID,
			Authority:    valueobject.AuthorityGeneral,
		}
		if err := tx.Create([]*entity.User{owner, boss, mentor, mentee, sharer}).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.DepartmentHierarchy{
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptDev.ID},
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptSales.ID},
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.DepartmentRole{
			Name:         "部長",
			UserID:       boss.ID,
			DepartmentID: deptDev.ID,
		}).Error; err != nil {
			return err
		}

		team := &entity.Team{Name: "開発チームA", DepartmentID: deptDev.ID}
		if err := tx.Create(team).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.TeamMember{
			{TeamID: team.ID, MemberID: mentor.ID, IsReader: true},
			{TeamID: team.ID, MemberID: mentee.ID, IsReader: false},
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.UserRelationship{
			RelationType: valueobject.RelationshipTypeBoss,
			SourceUserID: boss.ID,
			TargetUserID: mentor.ID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.UserReferenceUser{
			UserID:          boss.ID,
			ReferableUserID: mentee.ID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.UserReferenceTeam{
			UserID:          boss.ID,
			ReferableTeamID: team.ID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.UserReferenceDepartment{
			UserID:                owner.ID,
			ReferableDepartmentID: deptDev.ID,
		}).Error; err != nil {
			return err
		}

		template := &entity.Template{
			Title:     "1on1テンプレート",
			DisplayNo: 1,
			IsSample:  true,
			UserID:    owner.ID,
		}
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		topicRequired := &entity.Topic{
			DisplayNo:  1,
			Question:   "今月の成果",
			IsRequire:  true,
			TemplateID: template.ID,
		}
		topicOptional := &entity.Topic{
			DisplayNo:  2,
			Question:   "困っていること",
			IsRequire:  false,
			TemplateID: template.ID,
		}
		if err := tx.Create([]*entity.Topic{topicRequired, topicOptional}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.SystemSettings{
			MeetingWeek:       1,
			MeetingWeekday:    1,
			DefaultTemplateID: template.ID,
			LogOutputPath:     "",
			LogFilePrefix:     "",
		}).Error; err != nil {
			return err
		}

		recurWeek := 1
		recurWeekday := 1
		pairing := &entity.Pairing{
			MentorID:     mentor.ID,
			MenteeID:     mentee.ID,
			RecurWeek:    &recurWeek,
			RecurWeekday: &recurWeekday,
		}
		if err := tx.Create(pairing).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.PairingSharer{
			PairingID: pairing.ID,
			UserID:    sharer.ID,
		}).Error; err != nil {
			return err
		}

		implementedAt := time.Now()
		meeting := &entity.Meeting{
			ImplementedAt: &implementedAt,
			Status:        valueobject.MeetingStatusDone,
			TemplateID:    template.ID,
			PairingID:     pairing.ID,
		}
		if err := tx.Create(meeting).Error; err != nil {
			return err
		}

		answerRequired := &entity.Answer{
			Content:   "目標を達成しました。",
			MeetingID: meeting.ID,
			TopicID:   topicRequired.ID,
			UserID:    mentee.ID,
		}
		answerOptional := &entity.Answer{
			Content:   "特にありません。",
			MeetingID: meeting.ID,
			TopicID:   topicOptional.ID,
			UserID:    mentee.ID,
		}
		if err := tx.Create([]*entity.Answer{answerRequired, answerOptional}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.MeetingComment{
			Content:   "よく頑張りました。",
			IsLimited: false,
			MeetingID: meeting.ID,
			UserID:    mentor.ID,
		}).Error; err != nil {
			return err
		}

		commentText := "具体的で良いですね。"
		if err := tx.Create(&entity.AnswerComment{
			Content:  &commentText,
			AnswerID: answerRequired.ID,
			UserID:   mentor.ID,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
