package infrastructure

import (
	"growlink/internal/domain/entity"
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

	if err := AutoMigrate(db); err != nil {
		return nil, err
	}
	if err := migrateAndDropPairingTargets(db); err != nil {
		return nil, err
	}
	if err := dropLegacyOperationTargetTables(db); err != nil {
		return nil, err
	}
	if err := dropTemplateDisplayNo(db); err != nil {
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
		&entity.DepartmentRole{},
		&entity.UserDepartment{},

		// マスタテーブル
		&entity.Role{},
	)
}

func dropTemplateDisplayNo(db *gorm.DB) error {
	return db.Exec("ALTER TABLE templates DROP COLUMN IF EXISTS display_no").Error
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
