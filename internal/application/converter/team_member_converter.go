package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func TeamMemberInputToEntity(in inputmodel.TeamMemberInput) *entity.TeamMember {
	return &entity.TeamMember{
		IsReader: in.IsReader,
		TeamID:   in.TeamID,
		MemberID: in.MemberID,
	}
}

func TeamMemberEntityToOutput(e entity.TeamMember) outputmodel.TeamMemberOutput {
	out := outputmodel.TeamMemberOutput{
		IsReader: e.IsReader,
	}
	if e.Team != nil {
		out.Team = teamEntityToOutputShallow(*e.Team)
	}
	if e.Member != nil {
		out.Member = UserEntityToOutput(*e.Member)
	}
	return out
}

func TeamMemberEntitiesToOutputs(es []entity.TeamMember) []outputmodel.TeamMemberOutput {
	outputs := make([]outputmodel.TeamMemberOutput, len(es))
	for i, e := range es {
		outputs[i] = TeamMemberEntityToOutput(e)
	}
	return outputs
}
