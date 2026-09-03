package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func TeamInputToEntity(in inputmodel.TeamInput) *entity.Team {
	return &entity.Team{
		Name:         in.Name,
		DepartmentID: in.DepartmentID,
	}
}

func teamEntityToOutputShallow(e entity.Team) outputmodel.TeamOutput {
	out := outputmodel.TeamOutput{
		ID:           e.ID,
		Name:         e.Name,
		DepartmentID: e.DepartmentID,
	}
	if e.Department != nil {
		out.Department = DepartmentEntityToOutput(*e.Department)
	}
	return out
}

func TeamEntityToOutput(e entity.Team) outputmodel.TeamOutput {
	out := teamEntityToOutputShallow(e)
	if len(e.TeamMembers) > 0 {
		out.TeamMembers = TeamMemberEntitiesToOutputs(e.TeamMembers)
	}
	return out
}

func TeamEntitiesToOutputs(es []entity.Team) []outputmodel.TeamOutput {
	outputs := make([]outputmodel.TeamOutput, len(es))
	for i, e := range es {
		outputs[i] = TeamEntityToOutput(e)
	}
	return outputs
}
