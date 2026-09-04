package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func MeetingInputToEntity(in inputmodel.MeetingInput) *entity.Meeting {
	meeting := &entity.Meeting{
		Status:     in.Status,
		TemplateID: in.TemplateID,
		PairingID:  in.PairingID,
	}
	if !in.ImplementedAt.IsZero() {
		implementedAt := in.ImplementedAt
		meeting.ImplementedAt = &implementedAt
	}
	return meeting
}

func MeetingEntityToOutput(e entity.Meeting) outputmodel.MeetingOutput {
	out := outputmodel.MeetingOutput{
		ID:     e.ID,
		Status: e.Status,
	}
	if e.ImplementedAt != nil {
		out.ImplementedAt = *e.ImplementedAt
	}
	if e.Template != nil {
		out.Template = TemplateEntityToOutput(*e.Template)
	} else if e.TemplateID != 0 {
		out.Template = outputmodel.TemplateOutput{ID: e.TemplateID}
	}
	if e.Pairing != nil {
		out.Pairing = PairingEntityToOutput(*e.Pairing)
	} else if e.PairingID != 0 {
		out.Pairing = outputmodel.PairingOutput{ID: e.PairingID}
	}
	return out
}

func MeetingEntityToOutputPtr(e *entity.Meeting) *outputmodel.MeetingOutput {
	if e == nil {
		return nil
	}
	out := MeetingEntityToOutput(*e)
	return &out
}

func MeetingEntitiesToOutputs(es []entity.Meeting) []outputmodel.MeetingOutput {
	outputs := make([]outputmodel.MeetingOutput, len(es))
	for i, e := range es {
		outputs[i] = MeetingEntityToOutput(e)
	}
	return outputs
}
