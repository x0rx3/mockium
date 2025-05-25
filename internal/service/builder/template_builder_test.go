package builder

import (
	"gomock/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTemplateBuilder_SuccessBuild(t *testing.T) {
	_, err := NewTemplateBuilder(zap.NewNop()).Build("testdata_template_builder/success")
	assert.NoError(t, err)
}

func TestTemplateBuilder_ErrorUnmarshal(t *testing.T) {
	_, err := NewTemplateBuilder(zap.NewNop()).Build("testdata_template_builder/error_unmarshal")
	assert.Error(t, err)
}

func TestTemplateBuilder_ErrorNotFoundDir(t *testing.T) {
	_, err := NewTemplateBuilder(zap.NewNop()).Build("error_path")
	assert.Error(t, err)
}

func TestTemplateBuilder_ErrorValidate(t *testing.T) {
	builder := NewTemplateBuilder(zap.NewNop())

	template := model.Template{
		Path: "/user",
		Handle: []model.HandleTemplate{
			model.HandleTemplate{
				SetResponseTemplate: model.SetResponseTemplate{
					SetFile: "test",
					SetBody: map[string]any{"filed": "value"},
				},
			},
		},
	}

	err := builder.validate([]model.Template{template})
	assert.Error(t, err)
}
