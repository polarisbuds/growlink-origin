package repository_interface

import (
	"growlink/internal/domain/entity"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type TopicRepository interface {
	Create(topic *entity.Topic) error
}
