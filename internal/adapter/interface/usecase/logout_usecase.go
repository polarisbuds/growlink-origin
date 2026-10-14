package usecase_interface

import inputmodel "growlink/internal/application/model/input_model"

type LogoutUsecase interface {
	Signout(user *inputmodel.LDAPUserInput) error
}
