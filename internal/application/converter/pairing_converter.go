package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/entity"
)

func PairingTargetInputsToOptions(in []inputmodel.PairingTargetInput) []repository_interface.PairingTargetOption {
	out := make([]repository_interface.PairingTargetOption, len(in))
	for i, t := range in {
		out[i] = repository_interface.PairingTargetOption{Kind: t.Kind, TargetID: t.TargetID}
	}
	return out
}

func PairingInputToEntity(in inputmodel.PairingInput) *entity.Pairing {
	menteeID := in.MentorID
	for _, t := range in.Targets {
		if t.Kind == entity.PairingTargetKindUser {
			menteeID = t.TargetID
			break
		}
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
		CreatedByUserID: e.CreatedByUserID,
	}
	if e.Mentor != nil {
		out.Mentor = UserEntityToOutput(*e.Mentor)
	}
	if e.Mentee != nil {
		out.Mentee = UserEntityToOutput(*e.Mentee)
	}
	for _, t := range e.Targets {
		out.Targets = append(out.Targets, outputmodel.PairingTargetOutput{
			Kind:     t.Kind,
			TargetID: t.TargetID,
		})
	}
	if len(out.Targets) == 0 {
		for _, m := range e.Mentees {
			if m.User != nil {
				out.Targets = append(out.Targets, outputmodel.PairingTargetOutput{
					Kind:     entity.PairingTargetKindUser,
					TargetID: m.UserID,
					Label:    m.User.Name,
				})
			}
		}
		if len(out.Targets) == 0 && e.Mentee != nil {
			out.Targets = append(out.Targets, outputmodel.PairingTargetOutput{
				Kind:     entity.PairingTargetKindUser,
				TargetID: e.MenteeID,
				Label:    e.Mentee.Name,
			})
		}
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
