package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type SystemSettingsRepository interface {
	GetSingleton() (*entity.SystemSettings, error)
	Update(settings *entity.SystemSettings) error
}
