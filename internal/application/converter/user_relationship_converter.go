package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func UserRelationshipInputToEntity(in inputmodel.UserRelationshipInput) *entity.UserRelationship {
	return &entity.UserRelationship{
		RelationType: in.RelationType,
		SourceUserID: in.SourceUserID,
		TargetUserID: in.TargetUserID,
	}
}

func UserRelationshipEntityToOutput(e entity.UserRelationship) outputmodel.UserRelationshipOutput {
	out := outputmodel.UserRelationshipOutput{
		ID:           e.ID,
		RelationType: e.RelationType,
		SourceUserID: e.SourceUserID,
		TargetUserID: e.TargetUserID,
	}
	if e.SourceUser != nil {
		out.SourceUser = UserEntityToOutput(*e.SourceUser)
	}
	if e.TargetUser != nil {
		out.TargetUser = UserEntityToOutput(*e.TargetUser)
	}
	return out
}

func UserRelationshipEntitiesToOutputs(es []entity.UserRelationship) []outputmodel.UserRelationshipOutput {
	outputs := make([]outputmodel.UserRelationshipOutput, len(es))
	for i, e := range es {
		outputs[i] = UserRelationshipEntityToOutput(e)
	}
	return outputs
}
