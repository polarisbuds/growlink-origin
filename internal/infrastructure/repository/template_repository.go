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
		Order("created_at ASC, id ASC").
		Preload("Topics", nil).
		Preload("User", nil)
}

func (repo *TemplateRepository) GetAll() (ret0 []entity.Template, err error) {
	templates, err := templateChain(repo.db).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("テンプレートの全取得処理に失敗しました: %w", err)
	}
	return templates, nil
}

func (repo *TemplateRepository) GetByUUID(templateUUID uuid.UUID) (ret0 *entity.Template, err error) {
	template, err := templateChain(repo.db).Where("uuid = ?", templateUUID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("UUID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Get(templateID uint) (ret0 *entity.Template, err error) {
	template, err := templateChain(repo.db).Where("id = ?", templateID).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定でのテンプレート取得処理に失敗しました: %w", err)
	}
	return &template, nil
}

func (repo *TemplateRepository) Create(input *entity.Template) (err error) {
	if err = gorm.G[entity.Template](repo.db).Create(repo.ctx, input); err != nil {
		return fmt.Errorf("テンプレートの新規作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) Update(templateID uint, newTemplate *entity.Template) (err error) {
	return repo.db.WithContext(repo.ctx).Transaction(func(tx *gorm.DB) error {
		if _, err = gorm.G[entity.Template](tx).
			Where("id = ?", templateID).
			Omit(clause.Associations).
			Updates(repo.ctx, *newTemplate); err != nil {
			return fmt.Errorf("テンプレートの更新に失敗しました: %w", err)
		}
		if err = syncTemplateTopics(repo.ctx, tx, templateID, newTemplate.Topics); err != nil {
			return err
		}
		return nil
	})
}

func syncTemplateTopics(ctx context.Context, tx *gorm.DB, templateID uint, incoming []entity.Topic) error {
	existing, err := gorm.G[entity.Topic](tx).Where("template_id = ?", templateID).Find(ctx)
	if err != nil {
		return fmt.Errorf("トピックの取得に失敗しました: %w", err)
	}
	existingByUUID := make(map[uuid.UUID]entity.Topic, len(existing))
	for _, topic := range existing {
		existingByUUID[topic.UUID] = topic
	}

	kept := make(map[uuid.UUID]struct{}, len(incoming))
	for _, topic := range incoming {
		topic.TemplateID = templateID
		if topic.UUID == uuid.Nil {
			if err := gorm.G[entity.Topic](tx).Create(ctx, &topic); err != nil {
				return fmt.Errorf("トピックの作成に失敗しました: %w", err)
			}
			continue
		}
		current, ok := existingByUUID[topic.UUID]
		if !ok {
			return fmt.Errorf("指定されたトピックがこのテンプレートに存在しません")
		}
		if _, seen := kept[topic.UUID]; seen {
			return fmt.Errorf("トピックのUUIDが重複しています")
		}
		kept[topic.UUID] = struct{}{}
		if _, err := gorm.G[entity.Topic](tx).
			Where("id = ?", current.ID).
			Select("DisplayNo", "Question", "Description", "IsRequire").
			Updates(ctx, topic); err != nil {
			return fmt.Errorf("トピックの更新に失敗しました: %w", err)
		}
	}

	for _, topic := range existing {
		if _, ok := kept[topic.UUID]; ok {
			continue
		}
		if _, err := gorm.G[entity.Topic](tx).Where("id = ?", topic.ID).Delete(ctx); err != nil {
			return fmt.Errorf("トピックの削除に失敗しました: %w", err)
		}
	}
	return nil
}

func (repo *TemplateRepository) UpdateIsLock(templateID uint, isLock bool) (err error) {
	if _, err = gorm.G[entity.Template](repo.db).Where("id = ?", templateID).Update(repo.ctx, "is_lock", isLock); err != nil {
		return fmt.Errorf("テンプレートのロック状態更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TemplateRepository) UpdateIsDefault(templateID uint, isDefault bool) (err error) {
	return repo.db.WithContext(repo.ctx).Transaction(func(tx *gorm.DB) error {
		if isDefault {
			if err = tx.Model(&entity.Template{}).
				Where("is_default = ? AND deleted_at IS NULL", true).
				Update("is_default", false).Error; err != nil {
				return fmt.Errorf("既存の初期テンプレートフラグの解除に失敗しました: %w", err)
			}
		}
		res := tx.Model(&entity.Template{}).
			Where("id = ? AND deleted_at IS NULL", templateID).
			Update("is_default", isDefault)
		if res.Error != nil {
			return fmt.Errorf("テンプレートの初期フラグ更新に失敗しました: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (repo *TemplateRepository) Delete(templateID uint) (err error) {
	if _, err = gorm.G[entity.Template](repo.db).Where("id = ?", templateID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("テンプレートの削除に失敗しました: %w", err)
	}
	return nil
}
