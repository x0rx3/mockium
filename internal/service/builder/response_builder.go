package builder

import (
	"encoding/json"
	"fmt"
	"gomock/internal/model"
	"gomock/internal/service"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// NewResponseBuilder creates a new ResponsePreparator instance with the provided template response.
// It initializes the preparator with the given template and sets up the response preparation logic.
func NewResponseBuilder(templResp model.SetResponseTemplate) *ResponseBuilder {
	return &ResponseBuilder{
		templResp: templResp,
	}
}

// ResponseBuilder is a struct that implements the ResponseProvider interface.
// It is responsible for preparing the response based on the provided template response.
// The struct contains the template response and provides methods to prepare the response.
// It checks for the presence of required headers, status code, and response body.
type ResponseBuilder struct {
	log       *zap.Logger
	templResp model.SetResponseTemplate
}

// Build prepares the response based on the provided request and template response.
// It builds the response based on the template and the incoming request.
// The method can handle different types of response formats, including JSON bodies and files.
// It checks for the presence of required headers, status code, and response body.
// The method returns the prepared response and any error that occurred during preparation.
// If the response body is a JSON object, it is built using the build method.
// If the response body is a file, it checks if the file exists and sets the file path in the response.
// The method also sets the response headers and status code based on the template response.
// If an error occurs during preparation, it returns the error.
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

// build builds the response based on the provided template response and request.
// It processes the template response and replaces any placeholders with actual values from the request.
// The method supports different types of response formats, including JSON bodies and files.
// It returns the built response and any error that occurred during the process.
func (inst *ResponseBuilder) build(templResp map[string]any, req *http.Request) (map[string]any, error) {
	if len(templResp) == 0 {
		return nil, nil
	}

	response := make(map[string]any, len(templResp))
	for filedName, fieldValue := range templResp {
		switch fieldValT := fieldValue.(type) {
		case string:
			if service.RegexpResponseValuePlaceholder.MatchString(fieldValT) {
				placeholders := service.RegexpResponseValuePlaceholder.FindStringSubmatch(fieldValT)
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

// valueByPlacehoders compares the placeholders in the template response with the request parameters.
func (inst *ResponseBuilder) valueByPlacehoders(placeholders []string, req *http.Request) (any, error) {
	if len(placeholders) < 4 {
		return nil, fmt.Errorf("invalid placeholders")
	}

	switch placeholders[2] {
	case string(service.Headers):
		return req.Header.Get(placeholders[3]), nil
	case string(service.Query):
		return req.URL.Query().Get(placeholders[3]), nil
	case string(service.Path):
		vars := mux.Vars(req)
		return vars[placeholders[3]], nil
	case string(service.Form):
		return req.FormValue(placeholders[3]), nil
	case string(service.Body):
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
	return nil, fmt.Errorf("unxpected placeholder: %s", placeholders[2])
}
