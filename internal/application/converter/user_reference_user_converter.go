package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func UserReferenceUserInputToEntity(in inputmodel.UserReferenceUserInput) *entity.UserReferenceUser {
	return &entity.UserReferenceUser{
		UserID:          in.UserID,
		ReferableUserID: in.ReferableUserID,
	}
}

func UserReferenceUserEntityToOutput(e entity.UserReferenceUser) outputmodel.UserReferenceUserOutput {
	out := outputmodel.UserReferenceUserOutput{
		ID:              e.ID,
		UserID:          e.UserID,
		ReferableUserID: e.ReferableUserID,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	} else if e.UserID != 0 {
		out.User = outputmodel.UserOutput{ID: e.UserID}
	}
	if e.ReferableUser != nil {
		out.ReferableUser = UserEntityToOutput(*e.ReferableUser)
	} else if e.ReferableUserID != 0 {
		out.ReferableUser = outputmodel.UserOutput{ID: e.ReferableUserID}
	}
	return out
}

func UserReferenceUserEntitiesToOutputs(es []entity.UserReferenceUser) []outputmodel.UserReferenceUserOutput {
	outputs := make([]outputmodel.UserReferenceUserOutput, len(es))
	for i, e := range es {
		outputs[i] = UserReferenceUserEntityToOutput(e)
	}
	return outputs
}
