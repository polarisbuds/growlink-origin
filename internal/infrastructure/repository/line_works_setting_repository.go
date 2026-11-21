package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LineWorksSettingRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewLineWorksSettingRepository(ctx context.Context, dp *gorm.DB) *LineWorksSettingRepository {
	return &LineWorksSettingRepository{ctx: ctx, db: dp}
}

func (repo *LineWorksSettingRepository) ListByTitles(titles []string) (ret0 []entity.LineWorksSetting, err error) {
	if len(titles) == 0 {
		return nil, nil
	}
	rows, err := gorm.G[entity.LineWorksSetting](repo.db).Where("title IN ?", titles).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("LINE WORKS 通知文の取得に失敗しました: %w", err)
	}
	byTitle := make(map[string]entity.LineWorksSetting, len(rows))
	for _, row := range rows {
		byTitle[row.Title] = row
	}
	out := make([]entity.LineWorksSetting, 0, len(titles))
	for _, t := range titles {
		if row, ok := byTitle[t]; ok {
			out = append(out, row)
		}
	}
	return out, nil
}

func (repo *LineWorksSettingRepository) GetByUUID(id uuid.UUID) (ret0 *entity.LineWorksSetting, err error) {
	row, err := gorm.G[entity.LineWorksSetting](repo.db).Where("uuid = ?", id).First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("LINE WORKS 通知文の取得に失敗しました: %w", err)
	}
	return &row, nil
}

func (repo *LineWorksSettingRepository) Create(setting *entity.LineWorksSetting) (err error) {
	if setting == nil {
		return fmt.Errorf("作成対象の LINE WORKS 通知文が空です")
	}
	if err = gorm.G[entity.LineWorksSetting](repo.db).Create(repo.ctx, setting); err != nil {
		return fmt.Errorf("LINE WORKS 通知文の作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *LineWorksSettingRepository) Update(setting *entity.LineWorksSetting) (err error) {
	if setting == nil {
		return fmt.Errorf("更新対象の LINE WORKS 通知文が空です")
	}
	if _, err = gorm.G[entity.LineWorksSetting](repo.db).Where("id = ?", setting.ID).Omit(clause.Associations).Updates(repo.ctx, *setting); err != nil {
		return fmt.Errorf("LINE WORKS 通知文の更新に失敗しました: %w", err)
	}
	return nil
}