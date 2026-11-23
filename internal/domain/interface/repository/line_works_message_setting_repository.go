package repository_interface

import (
	"context"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type LineWorksMessageSettingRepository interface {
	ListByTitles(titles []string) ([]entity.LineWorksMessageSetting, error)
	GetByUUID(id uuid.UUID) (*entity.LineWorksMessageSetting, error)
	GetLineWorksMessageByKind(ctx context.Context, kind valueobject.LineWorksMessageKind) (*entity.LineWorksMessageSetting, error)
	Create(setting *entity.LineWorksMessageSetting) error
	Update(setting *entity.LineWorksMessageSetting) error
}