package usecase_interface

import outputmodel "growlink/internal/application/model/output_model"

type LoginUsecase interface {
	Authenticate(id, password string) (*outputmodel.LDAPUserOutput, error)
}
