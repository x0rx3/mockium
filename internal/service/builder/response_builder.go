package builder

import (
	"encoding/json"
	"fmt"
	"io"
	"mockium/internal/model"
	"mockium/internal/service/constants"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// ResponseBuilder is responsible for constructing a response object
// based on a provided response template. It supports dynamic value substitution
// using placeholders from the incoming HTTP request.
type ResponseBuilder struct {
	log       *zap.Logger               // Logger for error or debug output (optional, not used in current logic).
	templResp model.SetResponseTemplate // Template used to build the response.
}

// NewResponseBuilder creates a new instance of ResponseBuilder with the given response template.
//
// Parameters:
//   - templResp: the template that defines how the response should be structured.
//
// Returns a pointer to a ResponseBuilder.
func NewResponseBuilder(templResp model.SetResponseTemplate) *ResponseBuilder {
	return &ResponseBuilder{
		templResp: templResp,
	}
}

// Build constructs a model.SetResponse object from the template and the provided HTTP request.
// It evaluates dynamic placeholders in the template using request values.
//
// Parameters:
//   - req: the incoming HTTP request used for extracting dynamic values.
//
// Returns a constructed SetResponse object or an error if placeholder resolution fails.
func (inst *ResponseBuilder) Build(req *http.Request) (*model.SetResponse, error) {
	response := &model.SetResponse{}
	if inst.templResp.SetBody != nil {
		resp, err := inst.build(inst.templResp.SetBody, req)
		if err != nil {
			return nil, err
		}
		response.SetBody = resp
	} else if inst.templResp.SetFile != "" {
		f, err := os.Open(inst.templResp.SetFile)
		if err != nil {
			return nil, err
		}

		response.SetFile = f
	}

	response.SetHeaders = inst.templResp.SetHeaders
	response.SetStatus = inst.templResp.SetStatus

	return response, nil
}

// build recursively constructs the response body map, resolving any dynamic
// placeholders using values from the request.
//
// Parameters:
//   - templResp: a nested map representing the body structure with possible placeholders.
//   - req: the HTTP request from which values can be extracted.
//
// Returns a fully resolved map or an error if a placeholder fails to resolve.
func (inst *ResponseBuilder) build(templResp map[string]any, req *http.Request) (map[string]any, error) {
	if len(templResp) == 0 {
		return nil, nil
	}

	response := make(map[string]any, len(templResp))
	for filedName, fieldValue := range templResp {
		switch fieldValT := fieldValue.(type) {
		case string:
			if constants.RegexpResponseValuePlaceholder.MatchString(fieldValT) {
				placeholders := constants.RegexpResponseValuePlaceholder.FindStringSubmatch(fieldValT)
				if placeholderValue, err := inst.valueByPlacehoders(placeholders, req); err != nil {
					return nil, err
				} else {
					response[filedName] = placeholderValue
				}
				continue
			}
			response[filedName] = fieldValT
		case map[string]any:
			buildetMap, err := inst.build(fieldValT, req)
			if err != nil {
				return nil, err
			}
			response[filedName] = buildetMap
		default:
			response[filedName] = fieldValT
		}
	}

	return response, nil
}

// valueByPlacehoders resolves a value from the HTTP request based on the parsed
// placeholder format.
//
// Expected format for placeholders: {{<type>:<key>}}
// Supported types: headers, query, path, form, body
//
// Parameters:
//   - placeholders: array of matched strings from the placeholder regex.
//   - req: the HTTP request used to extract the actual value.
//
// Returns the resolved value or an error if the placeholder is invalid or cannot be fulfilled.
func (inst *ResponseBuilder) valueByPlacehoders(placeholders []string, req *http.Request) (any, error) {
	if len(placeholders) < 4 {
		return nil, fmt.Errorf("invalid placeholders")
	}

	switch placeholders[2] {
	case string(constants.Headers):
		return req.Header.Get(placeholders[3]), nil
	case string(constants.Query):
		return req.URL.Query().Get(placeholders[3]), nil
	case string(constants.Path):
		vars := mux.Vars(req)
		return vars[placeholders[3]], nil
	case string(constants.Form):
		return req.FormValue(placeholders[3]), nil
	case string(constants.Body):
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
