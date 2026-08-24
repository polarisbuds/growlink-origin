package infrastructure

import (
	"growlink/internal/domain/entity"
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

	if err := autoMigrateOnce(db); err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrateOnce(db *gorm.DB) error {
	return db.AutoMigrate(
		// トランザクションテーブル
		&entity.User{},
		&entity.Department{},
		&entity.Team{},
		&entity.TeamMember{},
		&entity.Peering{},
		&entity.Meeting{},
		&entity.Template{},
		&entity.Topic{},
		&entity.Answer{},
		&entity.MeetingComment{},
		&entity.AnswerComment{},
		&entity.ReadCommentable{},
		&entity.UserRelationship{},

		// マスタテーブル
		&entity.Authority{},
		&entity.Role{},
		&entity.MeetingStatus{},
		&entity.UserRelationshipType{},
	)
}
