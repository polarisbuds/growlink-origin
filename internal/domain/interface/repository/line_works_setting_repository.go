package repository_interface

import (
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type LineWorksSettingRepository interface {
	ListByTitles(titles []string) ([]entity.LineWorksSetting, error)
	GetByUUID(id uuid.UUID) (*entity.LineWorksSetting, error)
	Create(setting *entity.LineWorksSetting) error
	Update(setting *entity.LineWorksSetting) error
}