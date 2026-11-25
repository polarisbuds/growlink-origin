package converter

import (
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func PairingEntityToOutput(e entity.Pairing) outputmodel.PairingOutput {
	out := outputmodel.PairingOutput{
		ID:      e.ID,
		UUID:    e.UUID.String(),
		Deleted: e.DeletedAt.Valid,
	}
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

func PairingSharersToUserOutputs(sharers []entity.PairingSharer) []outputmodel.UserOutput {
	if len(sharers) == 0 {
		return nil
	}
	out := make([]outputmodel.UserOutput, len(sharers))
	for idx, s := range sharers {
		if s.User != nil {
			out[idx] = UserEntityToOutput(*s.User)
		} else if s.UserID != 0 {
			out[idx] = outputmodel.UserOutput{ID: s.UserID}
		}
	}
	return out
}
