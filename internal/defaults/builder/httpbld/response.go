package httpbld

import (
	"encoding/json"
	"fmt"
	"io"

	"mockium/internal/defs"
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
	"net/http"

	"github.com/gorilla/mux"
)

func NewResponseBuilder(template *model.SetResponse) core.ResponseBuilder[httpgw.HTTPRequest] {
	return &responseBuilder{
		template: template,
	}
}

type responseBuilder struct {
	template *model.SetResponse
}

func (inst *responseBuilder) Build(req httpgw.HTTPRequest) (*model.SetResponse, error) {
	response := &model.SetResponse{}
	if inst.template.SetBody != nil {
		resp, err := inst.buildResponseBody(inst.template.SetBody, req)
		if err != nil {
			return nil, err
		}
		response.SetBody = resp
	} else if inst.template.SetFile != "" {
		if placeholderChain, ok := inst.isPlaceholder(inst.template.SetFile); ok {
			placeholderValue, err := inst.processingPlaceholder(placeholderChain, req.Raw())
			if err != nil {
				return nil, err
			}

			if _, str := placeholderValue.(string); !str {
				return nil, fmt.Errorf("invalid value from placehoder, file can be only string")
			}
			response.SetFile = placeholderValue.(string)
		} else {
			response.SetFile = inst.template.SetFile
		}
	}

	response.SetHeaders = inst.template.SetHeaders
	response.SetStatus = inst.template.SetStatus

	return response, nil
}

func (inst *responseBuilder) buildResponseBody(template map[string]any, req httpgw.HTTPRequest) (map[string]any, error) {
	if len(template) == 0 {
		return nil, nil
	}

	response := make(map[string]any, len(template))
	for field, value := range template {
		switch valueT := value.(type) {
		case string:
			if placeholdersChain, ok := inst.isPlaceholder(valueT); ok {
				placeholderValue, err := inst.processingPlaceholder(placeholdersChain, req.Raw())
				if err != nil {
					return nil, err
				}

				response[field] = placeholderValue
				continue
			}

			response[field] = value
		case map[string]any:
			buildetMap, err := inst.buildResponseBody(valueT, req)
			if err != nil {
				return nil, err
			}
			response[field] = buildetMap
		default:
			response[field] = valueT
		}
	}

	return response, nil
}

func (inst *responseBuilder) isPlaceholder(val string) ([]string, bool) {
	if defs.RegexpResponseValuePlaceholder.MatchString(val) {
		return defs.RegexpResponseValuePlaceholder.FindStringSubmatch(val), true

	}
	return nil, false
}

func (inst *responseBuilder) processingPlaceholder(placeholders []string, req *http.Request) (any, error) {
	if len(placeholders) < 4 {
		return nil, fmt.Errorf("invalid placeholders")
	}

	switch placeholders[2] {
	case string(defs.Headers):
		return req.Header.Get(placeholders[3]), nil
	case string(defs.Query):
		return req.URL.Query().Get(placeholders[3]), nil
	case string(defs.Path):
		vars := mux.Vars(req)
		return vars[placeholders[3]], nil
	case string(defs.Form):
		return req.FormValue(placeholders[3]), nil
	case string(defs.Body):
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		mBody := make(map[string]any)
		if err := json.Unmarshal(body, &mBody); err != nil {
			return nil, err
		}

		return mBody[placeholders[3]], nil
	}
	return nil, fmt.Errorf("unexpected placeholder: %s", placeholders[2])
}
