package di

import (
	"fmt"
	"mockium/internal/runtime/registry"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func NewContainer() *Container {
	return &Container{
		builders:   registry.New(make(map[string]core.ServerBuilder, 0)),
		servers:    registry.New(make(map[string]core.Server, 0)),
		validators: registry.New(make(map[string]core.Validator, 0)),
		loggers:    registry.New(make(map[string]core.Logger, 0)),
	}
}

// Container is main struct for managing server builders, servers and validators
type Container struct {
	builders   registry.Registrator[core.ServerBuilder]
	servers    registry.Registrator[core.Server]
	validators registry.Registrator[core.Validator]
	loggers    registry.Registrator[core.Logger]
}

func (inst *Container) BuildServer(templ model.Template[[]model.Handle]) error {
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

func (inst *Container) BuildServers(templs *model.SortedTemplates) error {
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

func (inst *Container) DeleteBuilder(id string) error {
	return inst.builders.Delete(id)
}

func (inst *Container) AddBuilder(id string, builder core.ServerBuilder) {
	inst.builders.Add(id, builder)
}

func (inst *Container) StartAllServers() error {
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

func (inst *Container) StopAllServers() error { return nil }

func (inst *Container) StartServer(id string) error { return nil }

func (inst *Container) DeleteServer(id string) error { return nil }

func (inst *Container) GetServer(id string) (core.Server, error) { return nil, nil }

func (inst *Container) ValidateTemplate(template model.Template[[]model.HandleTemplate]) error {
	return nil
}

func (inst *Container) ValidateTemplates(template []model.Template[[]model.HandleTemplate]) error {
	return nil
}

func (inst *Container) AddValidator(id string, validator core.Validator) {}

func (inst *Container) DeleteValidator(id string) error { return nil }

func (inst *Container) GetValidator(id string) (core.Validator, error) { return nil, nil }

func (inst *Container) AddLogger(id string, logger core.Logger) {}

func (inst *Container) GetLogger(id string) core.Logger { return nil }

func (inst *Container) DeleteLogger(id string) error { return nil }
