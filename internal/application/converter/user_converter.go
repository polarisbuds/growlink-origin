package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func UserInputToEntity(in inputmodel.UserInput) *entity.User {
	return &entity.User{
		Name:         in.Name,
		Email:        in.Email,
		ImageURI:     in.ImageURI,
		RoleID:       in.RoleID,
		DepartmentID: in.DepartmentID,
		Authority:    in.Authority,
	}
}

func UserEntityToOutput(e entity.User) outputmodel.UserOutput {
	return outputmodel.UserOutput{
		ID:       e.ID,
		Name:     e.Name,
		Email:    e.Email,
		ImageURI: e.ImageURI,
	}
}

func UserEntitiesToOutputs(es []entity.User) []outputmodel.UserOutput {
	outputs := make([]outputmodel.UserOutput, len(es))
	for i, e := range es {
		outputs[i] = UserEntityToOutput(e)
	}
	return outputs
}
