package template_mng_usecase

import (
	inputmodel "growlink/internal/application/model/input_model"
	"strings"
)

func normalizeTopicInputs(topics []inputmodel.TopicInput) []inputmodel.TopicInput {
	out := make([]inputmodel.TopicInput, len(topics))
	for i, t := range topics {
		t.Question = strings.ReplaceAll(strings.ReplaceAll(t.Question, "\r\n", " "), "\n", " ")
		out[i] = t
	}
	return out
}
