package preparer

import (
	"encoding/json"
	"gomock/internal/model"
	"gomock/internal/service"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// NewResponsePreparer creates a new ResponsePreparator instance with the provided template response.
// It initializes the preparator with the given template and sets up the response preparation logic.
func NewResponsePreparer(templResp model.SetResponseTemplate) *ResponsePreparer {
	return &ResponsePreparer{
		templResp: templResp,
	}
}

// ResponsePreparer is a struct that implements the ResponseProvider interface.
// It is responsible for preparing the response based on the provided template response.
// The struct contains the template response and provides methods to prepare the response.
// It checks for the presence of required headers, status code, and response body.
type ResponsePreparer struct {
	templResp model.SetResponseTemplate
}

// Prepare prepares the response based on the provided request and template response.
// It builds the response based on the template and the incoming request.
// The method can handle different types of response formats, including JSON bodies and files.
// It checks for the presence of required headers, status code, and response body.
// The method returns the prepared response and any error that occurred during preparation.
// If the response body is a JSON object, it is built using the build method.
// If the response body is a file, it checks if the file exists and sets the file path in the response.
// The method also sets the response headers and status code based on the template response.
// If an error occurs during preparation, it returns the error.
func (inst *ResponsePreparer) Prepare(req *http.Request) (*model.SetResponse, error) {
	var response = &model.SetResponse{}
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
func (inst *ResponsePreparer) build(templResp map[string]any, req *http.Request) (map[string]any, error) {
	var resp = make(map[string]any)

	for filedName, fieldValue := range templResp {
		switch fieldValT := fieldValue.(type) {
		case string:
			if service.RegexpResponseValuePlaceholder.MatchString(fieldValT) {
				placeholders := service.RegexpResponseValuePlaceholder.FindStringSubmatch(fieldValT)
				if err := inst.compareRegexpField(filedName, resp, placeholders, req); err != nil {
					return nil, err
				}
				continue
			}
			resp[filedName] = fieldValT
		case map[string]any:
			buildetMap, err := inst.build(fieldValT, req)
			if err != nil {
				return nil, err
			}
			resp[filedName] = buildetMap
		default:
			resp[filedName] = fieldValT
		}
	}

	return resp, nil
}

// compareRegexpField compares the placeholders in the template response with the request parameters.
func (inst *ResponsePreparer) compareRegexpField(key string, res map[string]any, placeholders []string, req *http.Request) error {
	switch placeholders[2] {
	case string(service.Headers):
		res[key] = req.Header.Get(placeholders[3])
	case string(service.Query):
		res[key] = req.URL.Query().Get(placeholders[3])
	case string(service.Path):
		vars := mux.Vars(req)
		res[key] = vars[placeholders[3]]
	case string(service.Form):
		res[key] = req.FormValue(placeholders[3])
	case string(service.Body):
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		mBody := make(map[string]any)

		if err := json.Unmarshal(body, &mBody); err != nil {
			return err
		}

		res[key] = mBody[placeholders[2]]
	}
	return nil
}
