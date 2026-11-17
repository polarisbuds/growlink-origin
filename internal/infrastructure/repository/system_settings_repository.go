package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SystemSettingsRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewSystemSettingsRepository(ctx context.Context, dp *gorm.DB) *SystemSettingsRepository {
	return &SystemSettingsRepository{ctx: ctx, db: dp}
}

func (repo *SystemSettingsRepository) GetSingleton() (ret0 *entity.SystemSettings, err error) {
	const op = "system_settings_repository.GetSingleton"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	settings, err := gorm.G[entity.SystemSettings](repo.db).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("システム設定の取得に失敗しました: %w", err)
	}
	return &settings, nil
}

func (repo *SystemSettingsRepository) Update(settings *entity.SystemSettings) (err error) {
	const op = "system_settings_repository.Update"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if settings == nil {
		return fmt.Errorf("更新対象のシステム設定が空です")
	}
	if _, err = gorm.G[entity.SystemSettings](repo.db).Where("id = ?", settings.ID).Omit(clause.Associations).Updates(repo.ctx, *settings); err != nil {
		return fmt.Errorf("システム設定の更新に失敗しました: %w", err)
	}
	return nil
}
