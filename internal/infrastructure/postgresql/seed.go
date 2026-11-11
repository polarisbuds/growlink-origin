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
		roleBoss := &entity.Role{Name: "部長"}
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
			Name:      "オーナー",
			Email:     "owner@example.com",
			Password:  password,
			Authority: valueobject.AuthorityOwner,
		}
		boss := &entity.User{
			Name:      "上司",
			Email:     "boss@example.com",
			Password:  password,
			Authority: valueobject.AuthorityBoss,
		}
		mentor := &entity.User{
			Name:      "メンター",
			Email:     "mentor@example.com",
			Password:  password,
			Authority: valueobject.AuthorityGeneral,
		}
		mentee := &entity.User{
			Name:      "メンティ",
			Email:     "mentee@example.com",
			Password:  password,
			Authority: valueobject.AuthorityGeneral,
		}
		sharer := &entity.User{
			Name:      "共有者",
			Email:     "sharer@example.com",
			Password:  password,
			Authority: valueobject.AuthorityGeneral,
		}
		mentee2 := &entity.User{
			Name:      "メンティ2",
			Email:     "mentee2@example.com",
			Password:  password,
			Authority: valueobject.AuthorityGeneral,
		}
		if err := tx.Create([]*entity.User{owner, boss, mentor, mentee, sharer, mentee2}).Error; err != nil {
			return err
		}
		if err := tx.Create([]*entity.DepartmentMember{
			{DepartmentID: deptHQ.ID, UserID: owner.ID, RoleID: roleOwner.ID},
			{DepartmentID: deptDev.ID, UserID: boss.ID, RoleID: roleBoss.ID},
			{DepartmentID: deptDev.ID, UserID: mentor.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentee.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptSales.ID, UserID: sharer.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentee2.ID, RoleID: roleGeneral.ID},
		}).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.DepartmentHierarchy{
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptDev.ID},
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptSales.ID},
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Department{}).Where("id = ?", deptDev.ID).Update("head_user_id", boss.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.DepartmentMember{}).
			Where("department_id = ? AND user_id = ?", deptDev.ID, boss.ID).
			Updates(map[string]any{"is_head": true, "is_manager": true}).Error; err != nil {
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

		template := &entity.Template{
			Title:       "1on1テンプレート",
			Description: "メンターとメンティーの定例1on1で、成果と課題を共有するためのテンプレートです。",
			IsSample:    true,
			IsDefault:   true,
			UserID:      owner.ID,
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

		templateCheckin := &entity.Template{
			Title:       "週次チェックイン",
			Description: "週の目標とブロッカーを短時間で確認するためのテンプレートです。",
			IsSample:    false,
			UserID:      boss.ID,
		}
		if err := tx.Create(templateCheckin).Error; err != nil {
			return err
		}
		topicCheckinGoal := &entity.Topic{
			DisplayNo:  1,
			Question:   "今週の目標",
			IsRequire:  true,
			TemplateID: templateCheckin.ID,
		}
		topicCheckinBlocker := &entity.Topic{
			DisplayNo:  2,
			Question:   "ブロッカーはありますか",
			IsRequire:  false,
			TemplateID: templateCheckin.ID,
		}
		if err := tx.Create([]*entity.Topic{topicCheckinGoal, topicCheckinBlocker}).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.SystemSettings{
			MeetingWeek:       1,
			MeetingWeekday:    1,
			DefaultTemplateID: template.ID,
			LogOutputPath:     "",
			LogFilePrefix:     "",
			RemindAutoSend:    false,
			RemindWeekday:     1,
			RemindTime:        "09:00",
		}).Error; err != nil {
			return err
		}

		pairing := &entity.Pairing{
			MentorID:        mentor.ID,
			MenteeID:        mentee.ID,
			CreatedByUserID: mentor.ID,
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

		pairingBoss := &entity.Pairing{
			MentorID:        boss.ID,
			MenteeID:        mentee2.ID,
			CreatedByUserID: boss.ID,
		}
		if err := tx.Create(pairingBoss).Error; err != nil {
			return err
		}

		now := time.Now()
		pastDoneAt := now.AddDate(0, -1, 0)
		todayDoneAt := now
		todayPendingAt := now
		draftAt := now.AddDate(0, 0, -3)
		futurePendingAt := now.AddDate(0, 0, 14)
		bossPastDoneAt := now.AddDate(0, 0, -14)
		bossFuturePendingAt := now.AddDate(0, 0, 7)

		meetingPastDone := &entity.Meeting{
			ImplementedAt:   &pastDoneAt,
			Status:          valueobject.MeetingStatusDone,
			TemplateID:      template.ID,
			PairingID:       pairing.ID,
			CreatedByUserID: mentor.ID,
		}
		meetingTodaySubmitted := &entity.Meeting{
			ImplementedAt:   &todayDoneAt,
			Status:          valueobject.MeetingStatusSubmitted,
			TemplateID:      template.ID,
			PairingID:       pairing.ID,
			CreatedByUserID: mentor.ID,
		}
		meetingTodayPending := &entity.Meeting{
			ImplementedAt:   &todayPendingAt,
			Status:          valueobject.MeetingStatusPending,
			TemplateID:      templateCheckin.ID,
			PairingID:       pairing.ID,
			CreatedByUserID: mentor.ID,
		}
		meetingDraft := &entity.Meeting{
			ImplementedAt:   &draftAt,
			Status:          valueobject.MeetingStatusDraft,
			TemplateID:      template.ID,
			PairingID:       pairing.ID,
			CreatedByUserID: mentor.ID,
		}
		meetingFuturePending := &entity.Meeting{
			ImplementedAt:   &futurePendingAt,
			Status:          valueobject.MeetingStatusPending,
			TemplateID:      template.ID,
			PairingID:       pairing.ID,
			CreatedByUserID: mentor.ID,
		}
		meetingBossPastDone := &entity.Meeting{
			ImplementedAt:   &bossPastDoneAt,
			Status:          valueobject.MeetingStatusDone,
			TemplateID:      templateCheckin.ID,
			PairingID:       pairingBoss.ID,
			CreatedByUserID: boss.ID,
		}
		meetingBossFuturePending := &entity.Meeting{
			ImplementedAt:   &bossFuturePendingAt,
			Status:          valueobject.MeetingStatusPending,
			TemplateID:      templateCheckin.ID,
			PairingID:       pairingBoss.ID,
			CreatedByUserID: boss.ID,
		}
		if err := tx.Create([]*entity.Meeting{
			meetingPastDone,
			meetingTodaySubmitted,
			meetingTodayPending,
			meetingDraft,
			meetingFuturePending,
			meetingBossPastDone,
			meetingBossFuturePending,
		}).Error; err != nil {
			return err
		}
		if err := tx.Create([]*entity.MeetingSharer{
			{MeetingID: meetingPastDone.ID, UserID: sharer.ID},
			{MeetingID: meetingTodaySubmitted.ID, UserID: sharer.ID},
			{MeetingID: meetingDraft.ID, UserID: sharer.ID},
			{MeetingID: meetingBossPastDone.ID, UserID: mentor.ID},
		}).Error; err != nil {
			return err
		}

		doneMeetingAnswers := []*entity.Answer{
			{
				Content:   "先月はリリースを完了しました。",
				MeetingID: meetingPastDone.ID,
				TopicID:   topicRequired.ID,
				UserID:    mentee.ID,
			},
			{
				Content:   "チーム内の連携がスムーズでした。",
				MeetingID: meetingPastDone.ID,
				TopicID:   topicOptional.ID,
				UserID:    mentee.ID,
			},
			{
				Content:   "目標を達成しました。",
				MeetingID: meetingTodaySubmitted.ID,
				TopicID:   topicRequired.ID,
				UserID:    mentee.ID,
			},
			{
				Content:   "特にありません。",
				MeetingID: meetingTodaySubmitted.ID,
				TopicID:   topicOptional.ID,
				UserID:    mentee.ID,
			},
			{
				Content:   "下書き中の回答です。",
				MeetingID: meetingDraft.ID,
				TopicID:   topicRequired.ID,
				UserID:    mentee.ID,
			},
			{
				Content:   "週次の目標を達成しました。",
				MeetingID: meetingBossPastDone.ID,
				TopicID:   topicCheckinGoal.ID,
				UserID:    mentee2.ID,
			},
			{
				Content:   "特になし",
				MeetingID: meetingBossPastDone.ID,
				TopicID:   topicCheckinBlocker.ID,
				UserID:    mentee2.ID,
			},
		}
		if err := tx.Create(doneMeetingAnswers).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.MeetingComment{
			{
				Content:   "着実に成長しています。",
				IsLimited: false,
				MeetingID: meetingPastDone.ID,
				UserID:    mentor.ID,
			},
			{
				Content:   "次回までにブロッカーを整理してください。",
				IsLimited: false,
				MeetingID: meetingBossPastDone.ID,
				UserID:    boss.ID,
			},
		}).Error; err != nil {
			return err
		}

		pastRequiredComment := "進捗は順調です。次は品質改善に注力しましょう。"
		pastOptionalComment := "コードレビューの頻度を上げたいです。"
		todayRequiredComment := "具体的で良いですね。"
		bossGoalComment := "来週はレビュー体制を整えましょう。"
		if err := tx.Create([]*entity.AnswerComment{
			{
				Content:  &pastRequiredComment,
				AnswerID: doneMeetingAnswers[0].ID,
				UserID:   mentor.ID,
			},
			{
				Content:  &pastOptionalComment,
				AnswerID: doneMeetingAnswers[1].ID,
				UserID:   mentor.ID,
			},
			{
				Content:  &todayRequiredComment,
				AnswerID: doneMeetingAnswers[2].ID,
				UserID:   mentor.ID,
			},
			{
				Content:  &bossGoalComment,
				AnswerID: doneMeetingAnswers[5].ID,
				UserID:   boss.ID,
			},
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
