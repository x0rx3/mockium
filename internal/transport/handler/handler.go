package handler

import (
	"encoding/json"
	"mockium/internal/transport"
	"net/http"

	"go.uber.org/zap"
)

// Handler is an HTTP handler that routes incoming requests
// based on a set of request matchers and generates responses
// using associated response builders.
type Handler struct {
	log      *zap.Logger
	matchers map[transport.RequestMatcher]transport.ResponseBuilder
}

// New creates a new instance of Handler.
//
// Parameters:
//   - log: a zap.Logger instance for logging request/response activity.
//   - matchers: a map of RequestMatcher to corresponding ResponseBuilder.
//
// Returns:
//
//	A pointer to an initialized Handler.
func New(log *zap.Logger, mathcers map[transport.RequestMatcher]transport.ResponseBuilder) *Handler {
	return &Handler{
		log:      log,
		matchers: mathcers,
	}
}

// ServeHTTP handles incoming HTTP requests by matching them
// against configured request matchers. If a match is found,
// the corresponding response is built and sent.
//
// If no match is found, it responds with 404 Not Found.
// If an error occurs during response building, it responds with 500 Internal Server Error.
//
// Parameters:
//   - w: the HTTP response writer.
//   - r: the HTTP request.
func (inst *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inst.log.Info("request", zap.String("path", r.URL.Path))

	resProvider := inst.findMatches(r)
	if resProvider == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := resProvider.Build(r)
	if err != nil {
		http.Error(w, "failed prepare response", http.StatusInternalServerError)
		return
	}

	if response == nil {
		http.Error(w, "nil response after prepare", http.StatusInternalServerError)
		return
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
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bodyByte)
		return
	}

	w.WriteHeader(status)
}

// findMatches finds the first matching response builder for the incoming request
// by iterating over the registered request matchers.
//
// Parameters:
//   - req: the incoming HTTP request.
//
// Returns:
//
//	The first matching ResponseBuilder, or nil if no match is found.
func (inst *Handler) findMatches(req *http.Request) transport.ResponseBuilder {
	for reqMatcher, resProvider := range inst.matchers {
		if reqMatcher.Match(req) {
			return resProvider
		}
	}
	return nil
}
