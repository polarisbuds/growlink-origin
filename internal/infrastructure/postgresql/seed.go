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
// Login: *@example.com / password — see inline user comments for roles and pairing scenarios.
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
			Name: "オーナー", Email: "owner@example.com", Password: password,
			Authority: valueobject.AuthorityOwner,
		}
		admin := &entity.User{
			Name: "システム管理者", Email: "admin@example.com", Password: password,
			Authority: valueobject.AuthorityGeneral, SystemRole: valueobject.SystemRoleAdmin,
		}
		analyst := &entity.User{
			Name: "分析者", Email: "analyst@example.com", Password: password,
			Authority: valueobject.AuthorityGeneral, SystemRole: valueobject.SystemRoleAnalyst,
		}
		sharer := &entity.User{
			Name: "共有者", Email: "sharer@example.com", Password: password,
			Authority: valueobject.AuthorityGeneral,
		}
		hqHead := &entity.User{
			Name: "本社責任者", Email: "hq-head@example.com", Password: password,
			Authority: valueobject.AuthorityBoss,
		}
		devHead := &entity.User{
			Name: "開発責任者", Email: "dev-head@example.com", Password: password,
			Authority: valueobject.AuthorityBoss,
		}
		devManager := &entity.User{
			Name: "開発管理職", Email: "dev-manager@example.com", Password: password,
			Authority: valueobject.AuthorityBoss,
		}
		salesHead := &entity.User{
			Name: "営業責任者", Email: "sales-head@example.com", Password: password,
			Authority: valueobject.AuthorityBoss,
		}
		teamLeader := &entity.User{
			Name: "チームリーダー", Email: "leader@example.com", Password: password,
			Authority: valueobject.AuthorityBoss,
		}
		crossMentor := &entity.User{
			Name: "他部署メンター", Email: "cross@example.com", Password: password,
			Authority: valueobject.AuthorityGeneral,
		}
		mentor1 := &entity.User{Name: "メンター1", Email: "mentor1@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		mentor2 := &entity.User{Name: "メンター2", Email: "mentor2@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		mentor3 := &entity.User{Name: "メンター3", Email: "mentor3@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		mentor4 := &entity.User{Name: "メンター4", Email: "mentor4@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		mentor5 := &entity.User{Name: "メンター5", Email: "mentor5@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		salesMentee := &entity.User{
			Name: "営業メンティー", Email: "sales-mentee@example.com", Password: password,
			Authority: valueobject.AuthorityGeneral,
		}
		team1 := &entity.User{Name: "チームメンバー1", Email: "team1@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		team2 := &entity.User{Name: "チームメンバー2", Email: "team2@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		team3 := &entity.User{Name: "チームメンバー3", Email: "team3@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		team4 := &entity.User{Name: "チームメンバー4", Email: "team4@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		pure1 := &entity.User{Name: "メンティー（純）1", Email: "pure1@example.com", Password: password, Authority: valueobject.AuthorityGeneral}
		pure2 := &entity.User{Name: "メンティー（純）2", Email: "pure2@example.com", Password: password, Authority: valueobject.AuthorityGeneral}

		users := []*entity.User{
			owner, admin, analyst, sharer,
			hqHead, devHead, devManager, salesHead, teamLeader,
			crossMentor,
			mentor1, mentor2, mentor3, mentor4, mentor5,
			salesMentee, team1, team2, team3, team4, pure1, pure2,
		}
		if err := tx.Create(users).Error; err != nil {
			return err
		}

		deptMembers := []*entity.DepartmentMember{
			{DepartmentID: deptHQ.ID, UserID: owner.ID, RoleID: roleOwner.ID},
			{DepartmentID: deptHQ.ID, UserID: admin.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptHQ.ID, UserID: analyst.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptHQ.ID, UserID: hqHead.ID, RoleID: roleBoss.ID, IsHead: true, IsManager: true},
			{DepartmentID: deptDev.ID, UserID: devHead.ID, RoleID: roleBoss.ID, IsHead: true, IsManager: true},
			{DepartmentID: deptDev.ID, UserID: devManager.ID, RoleID: roleBoss.ID, IsManager: true},
			{DepartmentID: deptDev.ID, UserID: teamLeader.ID, RoleID: roleBoss.ID},
			{DepartmentID: deptDev.ID, UserID: crossMentor.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentor1.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentor2.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentor3.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentor4.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: mentor5.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: team1.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: team2.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: team3.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: team4.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: pure1.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptDev.ID, UserID: pure2.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptSales.ID, UserID: salesHead.ID, RoleID: roleBoss.ID, IsHead: true, IsManager: true},
			{DepartmentID: deptSales.ID, UserID: sharer.ID, RoleID: roleGeneral.ID},
			{DepartmentID: deptSales.ID, UserID: salesMentee.ID, RoleID: roleGeneral.ID},
		}
		if err := tx.Create(deptMembers).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.DepartmentHierarchy{
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptDev.ID},
			{ParentDepartmentID: deptHQ.ID, ChildDepartmentID: deptSales.ID},
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Department{}).Where("id = ?", deptHQ.ID).Update("head_user_id", hqHead.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Department{}).Where("id = ?", deptDev.ID).Update("head_user_id", devHead.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Department{}).Where("id = ?", deptSales.ID).Update("head_user_id", salesHead.ID).Error; err != nil {
			return err
		}

		team := &entity.Team{Name: "開発チームA", DepartmentID: deptDev.ID}
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		if err := tx.Create([]*entity.TeamMember{
			{TeamID: team.ID, MemberID: teamLeader.ID, IsReader: true},
			{TeamID: team.ID, MemberID: crossMentor.ID, IsReader: false},
			{TeamID: team.ID, MemberID: team1.ID, IsReader: false},
			{TeamID: team.ID, MemberID: team2.ID, IsReader: false},
			{TeamID: team.ID, MemberID: team3.ID, IsReader: false},
			{TeamID: team.ID, MemberID: team4.ID, IsReader: false},
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
		topicRequired := &entity.Topic{DisplayNo: 1, Question: "今月の成果", IsRequire: true, TemplateID: template.ID}
		topicOptional := &entity.Topic{DisplayNo: 2, Question: "困っていること", IsRequire: false, TemplateID: template.ID}
		if err := tx.Create([]*entity.Topic{topicRequired, topicOptional}).Error; err != nil {
			return err
		}

		templateCheckin := &entity.Template{
			Title:       "週次チェックイン",
			Description: "週の目標とブロッカーを短時間で確認するためのテンプレートです。",
			IsSample:    false,
			UserID:      devHead.ID,
		}
		if err := tx.Create(templateCheckin).Error; err != nil {
			return err
		}
		topicCheckinGoal := &entity.Topic{DisplayNo: 1, Question: "今週の目標", IsRequire: true, TemplateID: templateCheckin.ID}
		topicCheckinBlocker := &entity.Topic{DisplayNo: 2, Question: "ブロッカーはありますか", IsRequire: false, TemplateID: templateCheckin.ID}
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

		type pairDef struct {
			mentor, mentee *entity.User
		}
		pairDefs := []pairDef{
			{hqHead, devHead},
			{devHead, devManager},
			{devManager, crossMentor},
			{crossMentor, salesMentee},
			{teamLeader, team1},
			{mentor1, team2},
			{mentor2, team3},
			{mentor3, team4},
			{mentor4, pure1},
			{mentor5, pure2},
		}
		pairings := make([]*entity.Pairing, 0, len(pairDefs))
		for _, p := range pairDefs {
			pairings = append(pairings, &entity.Pairing{
				MentorID:        p.mentor.ID,
				MenteeID:        p.mentee.ID,
				CreatedByUserID: p.mentor.ID,
			})
		}
		if err := tx.Create(pairings).Error; err != nil {
			return err
		}

		now := time.Now()
		pastDoneAt := now.AddDate(0, 0, -14)
		futurePendingAt := now.AddDate(0, 0, 7)

		doneMeetings := make([]*entity.Meeting, 0, 4)
		for i := 0; i < 4; i++ {
			at := pastDoneAt.AddDate(0, 0, -i)
			doneMeetings = append(doneMeetings, &entity.Meeting{
				ImplementedAt:   &at,
				Status:          valueobject.MeetingStatusDone,
				TemplateID:      template.ID,
				PairingID:       pairings[i].ID,
				CreatedByUserID: pairDefs[i].mentor.ID,
			})
		}
		if err := tx.Create(doneMeetings).Error; err != nil {
			return err
		}

		pendingMeetings := make([]*entity.Meeting, 0, 6)
		for i := 4; i < 10; i++ {
			at := futurePendingAt.AddDate(0, 0, i-4)
			tpl := template.ID
			if i%2 == 0 {
				tpl = templateCheckin.ID
			}
			pendingMeetings = append(pendingMeetings, &entity.Meeting{
				ImplementedAt:   &at,
				Status:          valueobject.MeetingStatusPending,
				TemplateID:      tpl,
				PairingID:       pairings[i].ID,
				CreatedByUserID: pairDefs[i].mentor.ID,
			})
		}
		if err := tx.Create(pendingMeetings).Error; err != nil {
			return err
		}

		// 開発責任者の面談（本社責任者→開発責任者）を共有者へ
		if err := tx.Create(&entity.MeetingSharer{
			MeetingID: doneMeetings[0].ID,
			UserID:    sharer.ID,
		}).Error; err != nil {
			return err
		}

		answers := make([]*entity.Answer, 0, 8)
		for i, m := range doneMeetings {
			mentee := pairDefs[i].mentee
			answers = append(answers,
				&entity.Answer{
					Content:   fmt.Sprintf("（ペア%d）今月の成果の回答", i+1),
					MeetingID: m.ID,
					TopicID:   topicRequired.ID,
					UserID:    mentee.ID,
				},
				&entity.Answer{
					Content:   fmt.Sprintf("（ペア%d）困っていることの回答", i+1),
					MeetingID: m.ID,
					TopicID:   topicOptional.ID,
					UserID:    mentee.ID,
				},
			)
		}
		if err := tx.Create(answers).Error; err != nil {
			return err
		}

		if err := tx.Create([]*entity.MeetingComment{
			{
				Content:   "責任者として期待に沿っています。",
				IsLimited: false,
				MeetingID: doneMeetings[0].ID,
				UserID:    hqHead.ID,
			},
			{
				Content:   "管理職としての成長を続けてください。",
				IsLimited: false,
				MeetingID: doneMeetings[1].ID,
				UserID:    devHead.ID,
			},
		}).Error; err != nil {
			return err
		}

		c0 := "良い進捗です。"
		c1 := "課題共有ありがとうございます。"
		if err := tx.Create([]*entity.AnswerComment{
			{Content: &c0, AnswerID: answers[0].ID, UserID: hqHead.ID},
			{Content: &c1, AnswerID: answers[2].ID, UserID: devHead.ID},
		}).Error; err != nil {
			return err
		}

		return nil
	})
}