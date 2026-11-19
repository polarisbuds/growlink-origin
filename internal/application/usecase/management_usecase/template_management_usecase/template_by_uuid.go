package template_mng_usecase

import (
	"errors"
	"fmt"
	"strconv"

	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

var ErrTemplatePermissionDenied = errors.New("permission denied")

type TemplateActor struct {
	UserID     uint
	SystemRole valueobject.SystemRole
}

func (uc *TemplateViewer) resolveTemplateID(param string) (uint, error) {
	if id64, err := strconv.ParseUint(param, 10, 64); err == nil {
		return uint(id64), nil
	}
	templateUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	tpl, err := uc.templateRepo.GetByUUID(templateUUID)
	if err != nil {
		return 0, err
	}
	return tpl.ID, nil
}

func (uc *TemplateViewer) GetTemplateByParam(param string) (*outputmodel.TemplateOutput, error) {
	templateID, err := uc.resolveTemplateID(param)
	if err != nil {
		return nil, err
	}
	return uc.GetTemplate(templateID)
}

func (uc *TemplateUpdater) UpdateTemplateByParam(param string, input inputmodel.TemplateInput) error {
	id, err := uc.resolveTemplateIDForUpdater(param)
	if err != nil {
		return err
	}
	if _, err := uc.templateRepo.Get(id); err != nil {
		return fmt.Errorf("テンプレートの取得に失敗しました: %w", err)
	}
	return uc.UpdateTemplate(id, input)
}

func (uc *TemplateUpdater) resolveTemplateIDForUpdater(param string) (uint, error) {
	viewer := &TemplateViewer{templateRepo: uc.templateRepo}
	return viewer.resolveTemplateID(param)
}

func (uc *TemplateUpdater) UpdateTemplateLockByParam(param string, actor TemplateActor, isLock bool) error {
	id, err := uc.resolveTemplateIDForUpdater(param)
	if err != nil {
		return err
	}
	tpl, err := uc.templateRepo.Get(id)
	if err != nil {
		return err
	}
	if !actorCanModifyTemplate(actor, tpl.UserID) {
		return ErrTemplatePermissionDenied
	}
	return uc.UpdateTemplateLock(id, isLock)
}

func (uc *TemplateUpdater) UpdateTemplateDefaultByParam(param string, isDefault bool) error {
	id, err := uc.resolveTemplateIDForUpdater(param)
	if err != nil {
		return err
	}
	if _, err := uc.templateRepo.Get(id); err != nil {
		return err
	}
	return uc.UpdateTemplateDefault(id, isDefault)
}

func (uc *TemplateDeleter) DeleteTemplateByParam(param string, actor TemplateActor) error {
	id, err := uc.resolveTemplateIDForDeleter(param)
	if err != nil {
		return err
	}
	t, err := uc.templateRepo.Get(id)
	if err != nil {
		return err
	}
	if t.UserID != actor.UserID && !actor.SystemRole.IsAdmin() {
		return ErrTemplatePermissionDenied
	}
	return uc.DeleteTemplate(id)
}

func (uc *TemplateDeleter) resolveTemplateIDForDeleter(param string) (uint, error) {
	viewer := &TemplateViewer{templateRepo: uc.templateRepo}
	return viewer.resolveTemplateID(param)
}

func actorCanModifyTemplate(actor TemplateActor, ownerID uint) bool {
	if actor.SystemRole == valueobject.SystemRoleAdmin {
		return true
	}
	if ownerID == 0 {
		return false
	}
	return actor.UserID == ownerID
}

func (uc *TemplateCopier) CopyTemplateByParam(param string, title string, createdBy uint) error {
	viewer := &TemplateViewer{templateRepo: uc.templateRepo}
	id, err := viewer.resolveTemplateID(param)
	if err != nil {
		return err
	}
	if title == "" {
		tpl, err := viewer.GetTemplate(id)
		if err != nil {
			return err
		}
		title = tpl.Title + " のコピー"
	}
	return uc.CopyTemplate(id, title, createdBy)
}