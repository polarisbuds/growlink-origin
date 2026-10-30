package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func UserReferenceTeamInputToEntity(in inputmodel.UserReferenceTeamInput) *entity.UserReferenceTeam {
	return &entity.UserReferenceTeam{
		UserID:          in.UserID,
		ReferableTeamID: in.ReferableTeamID,
	}
}

func UserReferenceTeamEntityToOutput(e entity.UserReferenceTeam) outputmodel.UserReferenceTeamOutput {
	out := outputmodel.UserReferenceTeamOutput{
		ID:              e.ID,
		UserID:          e.UserID,
		ReferableTeamID: e.ReferableTeamID,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	} else if e.UserID != 0 {
		out.User = outputmodel.UserOutput{ID: e.UserID}
	}
	if e.ReferableTeam != nil {
		out.ReferableTeam = teamEntityToOutputShallow(*e.ReferableTeam)
	} else if e.ReferableTeamID != 0 {
		out.ReferableTeam = outputmodel.TeamOutput{ID: e.ReferableTeamID}
	}
	return out
}

func UserReferenceTeamEntitiesToOutputs(es []entity.UserReferenceTeam) []outputmodel.UserReferenceTeamOutput {
	outputs := make([]outputmodel.UserReferenceTeamOutput, len(es))
	for i, e := range es {
		outputs[i] = UserReferenceTeamEntityToOutput(e)
	}
	return outputs
}
