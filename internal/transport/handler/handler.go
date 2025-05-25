package handler

import (
	"encoding/json"
	"gomock/internal/transport"
	"net/http"

	"go.uber.org/zap"
)

// New creates a new Handler instance with the provided logger and request matchers.
// It initializes the handler with the given matchers and sets up the logger.
// The matchers are used to determine which response provider to use for a given request.
// The logger is used for logging errors and other information during request handling.
// The handler implements the http.Handler interface, allowing it to be used as an HTTP
func New(log *zap.Logger, mathcers map[transport.RequestMatcher]transport.ResponsePreparer) *Handler {
	return &Handler{
		log:      log,
		matchers: mathcers,
	}
}

// Handler is a struct that implements the http.Handler interface.
// It is responsible for handling HTTP requests and providing responses based on the request matchers and response providers.
type Handler struct {
	log      *zap.Logger
	matchers map[transport.RequestMatcher]transport.ResponsePreparer
}

// ServeHTTP is the main entry point for handling HTTP requests.
// It implements the http.Handler interface and is called when an HTTP request is received.
// It uses the request matchers to find the appropriate response provider for the request.
// If a matching response provider is found, it prepares the response and writes it to the http.ResponseWriter.
// If no matching response provider is found, it returns a 404 Not Found status.
// If an error occurs during response preparation, it returns a 500 Internal Server Error status.
// The response can include headers, a status code, a file to be served, or a JSON body.
// The response is written to the http.ResponseWriter, and the appropriate status code is set.
// The method is designed to be used in conjunction with the transport package to handle HTTP requests and responses.
func (inst *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resProvider := inst.findMatches(r)
	if resProvider == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := resProvider.Prepare(r)
	if err != nil {
		http.Error(w, "failed prepare response", http.StatusInternalServerError)
		return
	}

	if response == nil {
		http.Error(w, "nil response after prepare", http.StatusInternalServerError)
	}

	if response.SetHeaders != nil {
		for k, v := range response.SetHeaders {
			w.Header().Set(k, v)
		}
	}

	status := http.StatusOK
	if response.SetStatus != 0 {
		status = response.SetStatus
	}

	switch {
	case response.SetFile != nil:
		w.Header().Set("Content-Disposition", "attachment; filename="+response.SetFile.Name())
		w.WriteHeader(status)
		http.ServeFile(w, r, response.SetFile.Name())
		return
	case response.SetBody != nil:
		bodyByte, err := json.Marshal(response.SetBody)
		if err != nil {
			inst.log.Error("error marshal body", zap.Error(err))
			http.Error(w, "failed prepare response", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bodyByte)
		return
	}

	w.WriteHeader(status)
}

// findMatches iterates over the request matchers and checks if any of them match the given request.
// If a match is found, it returns the corresponding response provider.
// If no match is found, it returns nil.
func (inst *Handler) findMatches(req *http.Request) transport.ResponsePreparer {
	for reqMatcher, resProvider := range inst.matchers {
		if reqMatcher.Match(req) {
			return resProvider
		}
	}
	return nil
}
