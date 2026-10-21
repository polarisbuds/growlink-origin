package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Template]
	unScopedBaseRepository BaseUnScopedRepository[entity.Template]
}

func NewTemplateRepository(ctx context.Context, dp *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		baseRepository: BaseRepository[entity.Template]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Template]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *TemplateRepository) GetAll() ([]entity.Template, error) {
	templateChain := repo.baseRepository.Chain("display_no ASC, id ASC", "Topics", "User")
	templates, err := templateChain.Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("テンプレートの全取得処理に失敗しました: %w", err)
	}
	return templates, nil
}

func (repo *TemplateRepository) GetByUUID(templateUUID uuid.UUID) (*entity.Template, error) {
	templateChain := repo.baseRepository.Chain("display_no ASC, id ASC", "Topics", "User")
	template, err := templateChain.Where("uuid = ?", templateUUID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("UUID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Get(templateID uint) (*entity.Template, error) {
	templateChain := repo.baseRepository.Chain("display_no ASC, id ASC", "Topics", "User")
	template, err := templateChain.Where("id = ?", templateID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Create(input *entity.Template) error {
	if err := repo.baseRepository.Insert(input); err != nil {
		return fmt.Errorf("テンプレートの新規作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) Update(templateID uint, newTemplate *entity.Template) error {
	if _, err := repo.baseRepository.Updates(newTemplate, "id = ?", templateID); err != nil {
		return fmt.Errorf("テンプレートの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) UpdateDisplayNo(templateID uint, displayNo int) error {
	if _, err := repo.baseRepository.Update("display_no", displayNo, "id = ?", templateID); err != nil {
		return fmt.Errorf("テンプレートの並び順更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) Delete(templateID uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", templateID); err != nil {
		return fmt.Errorf("テンプレートの削除に失敗しました: %w", err)
	}
	return nil
}
