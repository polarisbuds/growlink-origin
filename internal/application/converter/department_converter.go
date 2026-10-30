package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func DepartmentInputToEntity(in inputmodel.DepartmentInput) *entity.Department {
	return &entity.Department{Name: in.Name}
}

func DepartmentEntityToOutput(e entity.Department) outputmodel.DepartmentOutput {
	return outputmodel.DepartmentOutput{ID: e.ID, Name: e.Name}
}

func DepartmentEntitiesToOutputs(es []entity.Department) []outputmodel.DepartmentOutput {
	outputs := make([]outputmodel.DepartmentOutput, len(es))
	for i, e := range es {
		outputs[i] = DepartmentEntityToOutput(e)
	}
	return outputs
}
