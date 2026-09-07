package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func DepartmentRoleInputToEntity(in inputmodel.DepartmentRoleInput) *entity.DepartmentRole {
	return &entity.DepartmentRole{
		Name:         in.Name,
		UserID:       in.UserID,
		DepartmentID: in.DepartmentID,
	}
}

func DepartmentRoleEntityToOutput(e entity.DepartmentRole) outputmodel.DepartmentRoleOutput {
	out := outputmodel.DepartmentRoleOutput{
		ID:           e.ID,
		Name:         e.Name,
		UserID:       e.UserID,
		DepartmentID: e.DepartmentID,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	} else if e.UserID != 0 {
		out.User = outputmodel.UserOutput{ID: e.UserID}
	}
	if e.Department != nil {
		out.Department = DepartmentEntityToOutput(*e.Department)
	} else if e.DepartmentID != 0 {
		out.Department = outputmodel.DepartmentOutput{ID: e.DepartmentID}
	}
	return out
}

func DepartmentRoleEntitiesToOutputs(es []entity.DepartmentRole) []outputmodel.DepartmentRoleOutput {
	outputs := make([]outputmodel.DepartmentRoleOutput, len(es))
	for i, e := range es {
		outputs[i] = DepartmentRoleEntityToOutput(e)
	}
	return outputs
}
