package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

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

func (repo *SystemSettingsRepository) GetSingleton() (*entity.SystemSettings, error) {
	settings, err := gorm.G[entity.SystemSettings](repo.db).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("システム設定の取得に失敗しました: %w", err)
	}
	return &settings, nil
}

func (repo *SystemSettingsRepository) Update(settings *entity.SystemSettings) error {
	if settings == nil {
		return fmt.Errorf("更新対象のシステム設定が空です")
	}
	if _, err := gorm.G[entity.SystemSettings](repo.db).Where("id = ?", settings.ID).Omit(clause.Associations).Updates(repo.ctx, *settings); err != nil {
		return fmt.Errorf("システム設定の更新に失敗しました: %w", err)
	}
	return nil
}
