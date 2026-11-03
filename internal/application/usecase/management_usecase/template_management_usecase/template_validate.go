package template_mng_usecase

import (
	"errors"
	"fmt"
	inputmodel "growlink/internal/application/model/input_model"
	"strings"
)

var ErrInvalidTemplateInput = errors.New("invalid template input")

func validateTemplateInput(input inputmodel.TemplateInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: タイトルを入力してください", ErrInvalidTemplateInput)
	}
	if strings.TrimSpace(input.Description) == "" {
		return fmt.Errorf("%w: 説明を入力してください", ErrInvalidTemplateInput)
	}
	if len(input.Topics) == 0 {
		return fmt.Errorf("%w: 質問を1件以上追加してください", ErrInvalidTemplateInput)
	}
	for i, t := range input.Topics {
		if strings.TrimSpace(t.Question) == "" {
			return fmt.Errorf("%w: Q%d の質問文を入力してください", ErrInvalidTemplateInput, i+1)
		}
	}
	return nil
}
