package repository_interface

import (
	"growlink/internal/domain/entity"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type TemplateRepository interface {
	GetAll() ([]entity.Template, error)
	Get(templateID uint) (*entity.Template, error)
	Create(input *entity.Template) error
	Update(templateID uint, newTemplate *entity.Template) error
	Delete(templateID uint) error
}
