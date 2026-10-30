package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TemplateRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTemplateRepository(ctx context.Context, dp *gorm.DB) *TemplateRepository {
	return &TemplateRepository{ctx: ctx, db: dp}
}

func templateChain(db *gorm.DB) gorm.ChainInterface[entity.Template] {
	return gorm.G[entity.Template](db).
		Order("display_no ASC, id ASC").
		Preload("Topics", nil).
		Preload("User", nil)
}

func (repo *TemplateRepository) GetAll() ([]entity.Template, error) {
	templates, err := templateChain(repo.db).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("テンプレートの全取得処理に失敗しました: %w", err)
	}
	return templates, nil
}

func (repo *TemplateRepository) GetByUUID(templateUUID uuid.UUID) (*entity.Template, error) {
	template, err := templateChain(repo.db).Where("uuid = ?", templateUUID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("UUID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Get(templateID uint) (*entity.Template, error) {
	template, err := templateChain(repo.db).Where("id = ?", templateID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Create(input *entity.Template) error {
	if err := gorm.G[entity.Template](repo.db).Create(repo.ctx, input); err != nil {
		return fmt.Errorf("テンプレートの新規作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) Update(templateID uint, newTemplate *entity.Template) error {
	if _, err := gorm.G[entity.Template](repo.db).Where("id = ?", templateID).Omit(clause.Associations).Updates(repo.ctx, *newTemplate); err != nil {
		return fmt.Errorf("テンプレートの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) UpdateDisplayNo(templateID uint, displayNo int) error {
	if _, err := gorm.G[entity.Template](repo.db).Where("id = ?", templateID).Update(repo.ctx, "display_no", displayNo); err != nil {
		return fmt.Errorf("テンプレートの並び順更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) Delete(templateID uint) error {
	if _, err := gorm.G[entity.Template](repo.db).Where("id = ?", templateID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("テンプレートの削除に失敗しました: %w", err)
	}
	return nil
}
