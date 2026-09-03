package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func PairingInputToEntity(in inputmodel.PairingInput) *entity.Pairing {
	return &entity.Pairing{
		MentorID: in.MentorID,
		MenteeID: in.MenteeID,
	}
}

func PairingEntityToOutput(e entity.Pairing) outputmodel.PairingOutput {
	out := outputmodel.PairingOutput{ID: e.ID}
	if e.Mentor != nil {
		out.Mentor = UserEntityToOutput(*e.Mentor)
	}
	if e.Mentee != nil {
		out.Mentee = UserEntityToOutput(*e.Mentee)
	}
	return out
}

func PairingEntitiesToOutputs(es []entity.Pairing) []outputmodel.PairingOutput {
	outputs := make([]outputmodel.PairingOutput, len(es))
	for i, e := range es {
		outputs[i] = PairingEntityToOutput(e)
	}
	return outputs
}
