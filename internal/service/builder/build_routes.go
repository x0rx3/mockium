package builder

import (
	"mockium/internal/model"
	"mockium/internal/service"
	"mockium/internal/service/matcher"
	"mockium/internal/transport"
	"mockium/internal/transport/handler"
	"mockium/internal/transport/route"
	"net/http"

	"go.uber.org/zap"
)

// Build is a function type that constructs a router from a template.
// It takes a logger for logging purposes and a template defining the routing rules,
// and returns an implementation of transport.Router.
type Build func(log *zap.Logger, procLogger service.ProcessLogger, template *model.Template) transport.Router

// BuildRoutes is the default implementation of the Build function.
// It creates a router with request matchers and response builders based on the provided template.
//
// The function performs the following steps:
// 1. Creates a two-level mapping of HTTP methods to request matchers and their corresponding response builders
// 2. Processes each handle from the template to populate the matchers map
// 3. Creates HTTP handlers for each method using the configured matchers
// 4. Returns a new router configured with the path and handlers from the template
//
// Parameters:
//   - log: Logger instance for logging operations
//   - template: Routing template containing path, handles and response configurations
//
// Returns:
//   - Configured router implementing transport.Router interface
var BuildRoutes Build = func(log *zap.Logger, procLogger service.ProcessLogger, template *model.Template) transport.Router {
	// matchersMap is a two-level map:
	// 1st level: HTTP method (e.g., GET, POST)
	// 2nd level: Map of request matchers to their response builders
	matchersMap := make(map[model.Method]map[transport.RequestMatcher]transport.ResponseBuilder)

	// handlers stores the final HTTP handlers for each method
	handlers := make(map[model.Method]http.Handler)

	// Process each handle definition from the template
	for _, handle := range template.Handle {
		// Initialize the inner map if it doesn't exist for this method
		if _, exists := matchersMap[handle.MatchRequestTemplate.MustMethod]; !exists {
			matchersMap[handle.MatchRequestTemplate.MustMethod] = make(map[transport.RequestMatcher]transport.ResponseBuilder)
		}

		// Add the matcher and response builder pair to the map
		matchersMap[handle.MatchRequestTemplate.MustMethod][matcher.NewRequestMatcher(log, &handle.MatchRequestTemplate)] =
			NewResponseBuilder(handle.SetResponseTemplate)
	}

	// Create handlers for each method using the configured matchers
	for mth, mtch := range matchersMap {
		handlers[mth] = handler.New(log, procLogger, mtch)
	}

	// Create and return a new router with the configured path and handlers
	return route.New(template.Path, handlers)
}
