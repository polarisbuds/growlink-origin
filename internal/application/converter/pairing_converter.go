package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func PairingInputToEntity(in inputmodel.PairingInput) *entity.Pairing {
	menteeID := in.MenteeID
	if len(in.MenteeIDs) > 0 {
		menteeID = in.MenteeIDs[0]
	}
	return &entity.Pairing{
		MentorID: in.MentorID,
		MenteeID: menteeID,
	}
}

func PairingEntityToOutput(e entity.Pairing) outputmodel.PairingOutput {
	out := outputmodel.PairingOutput{
		ID:              e.ID,
		UUID:            e.UUID.String(),
		RecurWeek:       e.RecurWeek,
		RecurWeekday:    e.RecurWeekday,
		CreatedByUserID: e.CreatedByUserID,
	}
	if e.Mentor != nil {
		out.Mentor = UserEntityToOutput(*e.Mentor)
	}
	if e.Mentee != nil {
		out.Mentee = UserEntityToOutput(*e.Mentee)
	}
	for _, m := range e.Mentees {
		if m.User != nil {
			out.Mentees = append(out.Mentees, UserEntityToOutput(*m.User))
		}
	}
	if len(out.Mentees) == 0 && e.Mentee != nil {
		out.Mentees = []outputmodel.UserOutput{out.Mentee}
	}
	for _, s := range e.Sharers {
		out.SharerUserIDs = append(out.SharerUserIDs, s.UserID)
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
