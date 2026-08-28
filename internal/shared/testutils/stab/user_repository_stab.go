package stab

import (
	"errors"

	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

var _ repository_interface.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	GetByIDFunc    func(id int64) (*entity.User, error)
	UpdateByIDFunc func(user *entity.User) error

	UpdatedUser *entity.User
}

func (s *UserRepository) GetByID(id int64) (*entity.User, error) {
	if s.GetByIDFunc == nil {
		return nil, errors.New("unexpected GetByID")
	}
	return s.GetByIDFunc(id)
}

func (s *UserRepository) UpdateByID(user *entity.User) error {
	s.UpdatedUser = user
	if s.UpdateByIDFunc == nil {
		return nil
	}
	return s.UpdateByIDFunc(user)
}
