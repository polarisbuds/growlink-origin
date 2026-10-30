package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func DepartmentHierarchyInputToEntity(in inputmodel.DepartmentHierarchyInput) *entity.DepartmentHierarchy {
	return &entity.DepartmentHierarchy{
		ParentDepartmentID: in.ParentDepartmentID,
		ChildDepartmentID:  in.ChildDepartmentID,
	}
}

func DepartmentHierarchyEntityToOutput(e entity.DepartmentHierarchy) outputmodel.DepartmentHierarchyOutput {
	out := outputmodel.DepartmentHierarchyOutput{
		ID:                 e.ID,
		ParentDepartmentID: e.ParentDepartmentID,
		ChildDepartmentID:  e.ChildDepartmentID,
	}
	if e.ParentDepartment != nil {
		out.ParentDepartment = DepartmentEntityToOutput(*e.ParentDepartment)
	} else if e.ParentDepartmentID != 0 {
		out.ParentDepartment = outputmodel.DepartmentOutput{ID: e.ParentDepartmentID}
	}
	if e.ChildDepartment != nil {
		out.ChildDepartment = DepartmentEntityToOutput(*e.ChildDepartment)
	} else if e.ChildDepartmentID != 0 {
		out.ChildDepartment = outputmodel.DepartmentOutput{ID: e.ChildDepartmentID}
	}
	return out
}

func DepartmentHierarchyEntitiesToOutputs(es []entity.DepartmentHierarchy) []outputmodel.DepartmentHierarchyOutput {
	outputs := make([]outputmodel.DepartmentHierarchyOutput, len(es))
	for i, e := range es {
		outputs[i] = DepartmentHierarchyEntityToOutput(e)
	}
	return outputs
}
