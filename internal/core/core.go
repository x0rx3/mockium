// Package core implements the core application logic for managing
// server builders, servers, validators, and loggers.
package core

import (
	"fmt"
	"mockium/internal/adapters/registry"
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

// NewApplicationCore creates a new instance of ApplicationCore.
func NewApplicationCore() *ApplicationCore {
	return &ApplicationCore{
		builders:   registry.New(make(map[string]ports.ServerBuilder, 0)),
		servers:    registry.New(make(map[string]ports.Server, 0)),
		validators: registry.New(make(map[string]ports.Validator, 0)),
		loggers:    registry.New(make(map[string]ports.Logger, 0)),
	}
}

// ApplicationCore manages server builders, servers, validators, and loggers.
type ApplicationCore struct {
	builders   ports.Registry[ports.ServerBuilder]
	servers    ports.Registry[ports.Server]
	validators ports.Registry[ports.Validator]
	loggers    ports.Registry[ports.Logger]
}

// BuildServer constructs and registers a server based on the provided template.
func (inst *ApplicationCore) BuildServer(templ model.Template[[]model.Handle]) error {
	idx := templ.Type
	builder, ok := inst.builders.Get(idx)
	if !ok {
		return fmt.Errorf("builder not found: %s", idx)
	}

	logger, ok := inst.loggers.Get(idx)
	if !ok {
		return fmt.Errorf("logger not found: %s", idx)
	}

	server := builder.Build(logger, []model.Template[[]model.Handle]{templ})
	inst.servers.Add(idx, server)

	return nil
}

// BuildServers constructs and registers servers for all provided templates.
func (inst *ApplicationCore) BuildServers(templs *model.SortedTemplates) error {
	for t, templs := range templs.GetAll() {
		builder, ok := inst.builders.Get(t)
		if !ok {
			return fmt.Errorf("builder not found: %s", t)
		}

		logger, ok := inst.loggers.Get(t)
		if !ok {
			return fmt.Errorf("logger not found: %s", t)
		}

		server := builder.Build(logger, templs)
		inst.servers.Add(t, server)
	}

	return nil
}

// DeleteBuilder removes a server builder by its identifier.
func (inst *ApplicationCore) DeleteBuilder(id string) error {
	return inst.builders.Delete(id)
}

// AddBuilder registers a new server builder with the given identifier.
func (inst *ApplicationCore) AddBuilder(id string, builder ports.ServerBuilder) {
	inst.builders.Add(id, builder)
}

// StartAllServers initiates all registered servers concurrently.
func (inst *ApplicationCore) StartAllServers() error {
	for _, server := range inst.servers.GetAll() {
		s := server
		go func() {
			if err := s.Configure(); err != nil {
				fmt.Printf("failed to configure server: %s", err.Error())
				return
			}

			if err := s.Start(); err != nil {
				fmt.Printf("failed to start server: %s", err.Error())
			}
		}()

		if !server.IsRunning() {
			return fmt.Errorf("failed to start server")
		}

	}

	return nil
}

// StopAllServers halts all running servers.
func (inst *ApplicationCore) StopAllServers() error {
	for _, server := range inst.servers.GetAll() {
		if server.IsRunning() {
			server.Stop()
		}
	}

	return nil
}

// StartServer initiates a specific server by its identifier.
func (inst *ApplicationCore) StartServer(id string) error {
	server, ok := inst.servers.Get(id)
	if !ok {
		return fmt.Errorf("server not found: %s", id)
	}

	s := server
	go func() {
		if err := s.Configure(); err != nil {
			fmt.Printf("failed to configure server: %s", err.Error())
			return
		}

		if err := s.Start(); err != nil {
			fmt.Printf("failed to start server: %s", err.Error())
		}
	}()

	if !server.IsRunning() {
		return fmt.Errorf("failed to start server")
	}

	return nil
}

// DeleteServer removes a server by its identifier, stopping it first if it's running.
func (inst *ApplicationCore) DeleteServer(id string) error {
	if server, ok := inst.servers.Get(id); ok {
		if server.IsRunning() {
			server.Stop()
		}

		return inst.servers.Delete(id)
	}

	return fmt.Errorf("server not found: %s", id)
}

// GetServer retrieves a server by its identifier.
func (inst *ApplicationCore) GetServer(id string) (ports.Server, error) {
	server, ok := inst.servers.Get(id)
	if !ok {
		return nil, fmt.Errorf("server not found: %s", id)
	}

	return server, nil
}

// ValidateTemplate checks if a single template is valid using the appropriate validator.
func (inst *ApplicationCore) ValidateTemplate(template model.Template[[]model.HandleTemplate]) error {
	templType := template.Type
	validator, ok := inst.validators.Get(templType)
	if !ok {
		return fmt.Errorf("validator not found: %s", templType)
	}

	return validator.Validate(template)
}

// ValidateTemplates checks if multiple templates are valid using their respective validators.
func (inst *ApplicationCore) ValidateTemplates(template []model.Template[[]model.HandleTemplate]) error {
	for _, templ := range template {
		templType := templ.Type
		validator, ok := inst.validators.Get(templType)
		if !ok {
			return fmt.Errorf("validator not found: %s", templType)
		}

		if err := validator.Validate(templ); err != nil {
			return err
		}
	}

	return nil
}

// AddValidator registers a new validator with the given identifier.
func (inst *ApplicationCore) AddValidator(id string, validator ports.Validator) {
	inst.validators.Add(id, validator)
}

func (inst *ApplicationCore) DeleteValidator(id string) error {
	return inst.validators.Delete(id)
}

// GetValidator retrieves a validator by its identifier.
func (inst *ApplicationCore) GetValidator(id string) (ports.Validator, error) {
	validator, ok := inst.validators.Get(id)
	if !ok {
		return nil, fmt.Errorf("validator not found: %s", id)
	}

	return validator, nil
}

// AddLogger registers a new logger with the given identifier.
func (inst *ApplicationCore) AddLogger(id string, logger ports.Logger) {
	inst.loggers.Add(id, logger)
}

// GetLogger retrieves a logger by its identifier.
func (inst *ApplicationCore) GetLogger(id string) ports.Logger { return inst.loggers.GetOrNil(id) }

// DeleteLogger removes a logger by its identifier.|
func (inst *ApplicationCore) DeleteLogger(id string) error { return inst.loggers.Delete(id) }
