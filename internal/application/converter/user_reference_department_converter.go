package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func UserReferenceDepartmentInputToEntity(in inputmodel.UserReferenceDepartmentInput) *entity.UserReferenceDepartment {
	return &entity.UserReferenceDepartment{
		UserID:                in.UserID,
		ReferableDepartmentID: in.ReferableDepartmentID,
	}
}

func UserReferenceDepartmentEntityToOutput(e entity.UserReferenceDepartment) outputmodel.UserReferenceDepartmentOutput {
	out := outputmodel.UserReferenceDepartmentOutput{
		ID:                    e.ID,
		UserID:                e.UserID,
		ReferableDepartmentID: e.ReferableDepartmentID,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	} else if e.UserID != 0 {
		out.User = outputmodel.UserOutput{ID: e.UserID}
	}
	if e.ReferableDepartment != nil {
		out.ReferableDepartment = DepartmentEntityToOutput(*e.ReferableDepartment)
	} else if e.ReferableDepartmentID != 0 {
		out.ReferableDepartment = outputmodel.DepartmentOutput{ID: e.ReferableDepartmentID}
	}
	return out
}

func UserReferenceDepartmentEntitiesToOutputs(es []entity.UserReferenceDepartment) []outputmodel.UserReferenceDepartmentOutput {
	outputs := make([]outputmodel.UserReferenceDepartmentOutput, len(es))
	for i, e := range es {
		outputs[i] = UserReferenceDepartmentEntityToOutput(e)
	}
	return outputs
}
