package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mockium/internal/model"
	"mockium/internal/service"
	"mockium/internal/transport"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Handler is an HTTP handler that routes incoming requests
// based on a set of request matchers and generates responses
// using associated response builders.
type Handler struct {
	log           *zap.Logger
	matchers      map[transport.RequestMatcher]transport.ResponseBuilder
	processLogger service.ProcessLogger
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
func New(log *zap.Logger, proceLogger service.ProcessLogger, mathcers map[transport.RequestMatcher]transport.ResponseBuilder) *Handler {
	return &Handler{
		log:           log,
		matchers:      mathcers,
		processLogger: proceLogger,
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
	logReq := inst.buildLogRequest(r)

	resProvider := inst.findMatches(r)
	if resProvider == nil {
		logReq.Response.SetStatus = http.StatusNotFound
		inst.processLogger.Log(logReq)
		inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.String("Response", "StatusNotFound"))

		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	response, err := resProvider.Build(r)
	if err != nil {
		logReq.Response.SetStatus = http.StatusInternalServerError
		inst.processLogger.Log(logReq)
		inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.String("Response", "StatusInternalServerError"))

		http.Error(w, "failed prepare response", http.StatusInternalServerError)
		return
	}

	if response == nil {
		logReq.Response.SetStatus = http.StatusInternalServerError
		inst.processLogger.Log(logReq)
		inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.String("Response", "StatusInternalServerError"))

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
	} else {
		response.SetStatus = status
	}

	switch {
	case response.SetFile != nil:
		logReq.Response = *response
		inst.processLogger.Log(logReq)
		inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.Any("Response", response))

		w.Header().Set("Content-Disposition", "attachment; filename="+response.SetFile.Name())
		w.WriteHeader(status)
		http.ServeFile(w, r, response.SetFile.Name())
		return
	case response.SetBody != nil:
		logReq.Response = *response
		inst.processLogger.Log(logReq)
		inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.Any("Response", response))

		bodyByte, err := json.Marshal(response.SetBody)
		if err != nil {
			inst.log.Info("Serve HTTP",
				zap.Any("Request", logReq),
				zap.String("Response", "StatusInternalServerError"),
			)
			http.Error(w, "failed prepare response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bodyByte)
		return
	}

	logReq.Response = *response
	inst.processLogger.Log(logReq)
	inst.log.Info("Serve HTTP", zap.Any("Request", logReq), zap.Any("Response", response))

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

func (inst *Handler) buildLogRequest(r *http.Request) *model.ProcessLoggingFileds {
	logReq := &model.LogginRequest{
		Headers: make(map[string]any),
	}

	logReq.Url = r.URL.String()
	logReq.RemoteAddr = r.RemoteAddr
	logReq.Method = r.Method

	for name, values := range r.Header {
		logReq.Headers[name] = values
	}

	if r.Body != nil && r.Body != http.NoBody {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logReq.Body = string(bodyBytes)
	}

	inst.log.Info("", zap.Any("Received Request", logReq))

	return &model.ProcessLoggingFileds{
		Time:    time.Now(),
		Request: logReq,
	}
}
