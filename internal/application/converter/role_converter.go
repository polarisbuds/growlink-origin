package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func RoleInputToEntity(in inputmodel.RoleInput) *entity.Role {
	return &entity.Role{Name: in.Name}
}

func RoleEntityToOutput(e entity.Role) outputmodel.RoleOutput {
	return outputmodel.RoleOutput{ID: e.ID, Name: e.Name}
}

func RoleEntitiesToOutputs(es []entity.Role) []outputmodel.RoleOutput {
	outputs := make([]outputmodel.RoleOutput, len(es))
	for i, e := range es {
		outputs[i] = RoleEntityToOutput(e)
	}
	return outputs
}
