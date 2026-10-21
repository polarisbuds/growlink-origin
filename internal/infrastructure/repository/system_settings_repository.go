package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type SystemSettingsRepository struct {
	ctx            context.Context
	baseRepository BaseRepository[entity.SystemSettings]
}

func NewSystemSettingsRepository(ctx context.Context, dp *gorm.DB) *SystemSettingsRepository {
	return &SystemSettingsRepository{
		ctx: ctx,
		baseRepository: BaseRepository[entity.SystemSettings]{
			dp:   dp,
			bctx: ctx,
		},
	}
}

func (repo *SystemSettingsRepository) GetSingleton() (*entity.SystemSettings, error) {
	settings, err := repo.baseRepository.First("id ASC", "1=1")
	if err != nil {
		return nil, fmt.Errorf("システム設定の取得に失敗しました: %w", err)
	}
	if settings == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return settings, nil
}

func (repo *SystemSettingsRepository) Update(settings *entity.SystemSettings) error {
	if settings == nil {
		return fmt.Errorf("更新対象のシステム設定が空です")
	}
	if _, err := repo.baseRepository.Updates(settings, "id = ?", settings.ID); err != nil {
		return fmt.Errorf("システム設定の更新に失敗しました: %w", err)
	}
	return nil
}
