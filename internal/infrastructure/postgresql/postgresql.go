package infrastructure

import (
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresqlDB(dsn string) (*gorm.DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// GORMがレコード作成・更新時に使う時刻をJSTに固定
		NowFunc: func() time.Time {
			return time.Now().In(jst)
		},
	})

	if err != nil {
		return nil, err
	}

	if err := renameLineWorksSettingsTable(db); err != nil {
		return nil, err
	}

	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	if err := backfillLineWorksMessageKind(db); err != nil {
		return nil, err
	}

	if err := ensureLineWorksDefaultMessages(db); err != nil {
		return nil, err
	}

	if config.IsDevelop() {
		if err := genSeed(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// AutoMigrate applies the GORM schema for all domain entities.
// Production and test environments share this migration list.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// トランザクションテーブル
		&entity.User{},
		&entity.Department{},
		&entity.Team{},
		&entity.TeamMember{},
		&entity.DepartmentMember{},
		&entity.Pairing{},
		&entity.PairingSharer{},
		&entity.Meeting{},
		&entity.MeetingSharer{},
		&entity.Template{},
		&entity.Topic{},
		&entity.Answer{},
		&entity.MeetingComment{},
		&entity.AnswerComment{},
		&entity.SystemSettings{},
		&entity.DepartmentHierarchy{},
		&entity.LineWorksMessageSetting{},

		// マスタテーブル
		&entity.Role{},
	)
}

func renameLineWorksSettingsTable(db *gorm.DB) error {
	if db.Migrator().HasTable("line_works_settings") && !db.Migrator().HasTable("line_works_message_settings") {
		return db.Exec("ALTER TABLE line_works_settings RENAME TO line_works_message_settings").Error
	}
	return nil
}

var lineWorksDefaultKinds = []valueobject.LineWorksMessageKind{
	valueobject.LineWorksMessageKindDraft,
	valueobject.LineWorksMessageKindSubmit,
	valueobject.LineWorksMessageKindNotConducted,
}

func backfillLineWorksMessageKind(db *gorm.DB) error {
	if !db.Migrator().HasTable("line_works_message_settings") {
		return nil
	}
	for title, kind := range map[string]valueobject.LineWorksMessageKind{
		"未提出通知": valueobject.LineWorksMessageKindDraft,
		"提出通知":  valueobject.LineWorksMessageKindSubmit,
		"未実施通知": valueobject.LineWorksMessageKindNotConducted,
	} {
		if err := db.Model(&entity.LineWorksMessageSetting{}).
			Where("title = ? AND (kind = 0 OR kind IS NULL)", title).
			Update("kind", kind).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureLineWorksDefaultMessages(db *gorm.DB) error {
	for _, kind := range lineWorksDefaultKinds {
		var count int64
		if err := db.Model(&entity.LineWorksMessageSetting{}).Where("kind = ?", kind).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		content, ok := kind.DefaultContent()
		if !ok {
			continue
		}
		if err := db.Create(&entity.LineWorksMessageSetting{
			Kind:    kind,
			Title:   kind.DefaultTitle(),
			Content: content,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func dropTemplateDisplayNo(db *gorm.DB) error {
	return db.Exec("ALTER TABLE templates DROP COLUMN IF EXISTS display_no").Error
}

func dropRoleRank(db *gorm.DB) error {
	return db.Exec("ALTER TABLE roles DROP COLUMN IF EXISTS rank").Error
}

func dropDepartmentRolesTable(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS department_roles").Error
}

func ensureTemplateDefaultIndex(db *gorm.DB) error {
	return db.Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS idx_templates_one_default
ON templates ((true))
WHERE is_default AND deleted_at IS NULL
`).Error
}

// backfillTemplateIsDefault sets is_default on system_settings.default_template_id when no template is marked default yet.
func backfillTemplateIsDefault(db *gorm.DB) error {
	var defaultCount int64
	if err := db.Model(&entity.Template{}).Where("is_default = ? AND deleted_at IS NULL", true).Count(&defaultCount).Error; err != nil {
		return err
	}
	if defaultCount > 0 {
		return nil
	}
	var settings entity.SystemSettings
	if err := db.First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if settings.DefaultTemplateID == 0 {
		return nil
	}
	return db.Model(&entity.Template{}).
		Where("id = ? AND deleted_at IS NULL", settings.DefaultTemplateID).
		Update("is_default", true).Error
}

func dropLegacyOperationTargetTables(db *gorm.DB) error {
	tables := []string{
		"user_reference_users",
		"user_reference_teams",
		"user_reference_departments",
		"user_relationships",
	}
	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			return err
		}
	}
	return nil
}
