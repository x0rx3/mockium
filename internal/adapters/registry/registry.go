package registry

import (
	"fmt"
	"mockium/pkg/ports"
	"sync"
)

func New[T any](m map[string]T) ports.Registry[T] {
	return &registry[T]{
		m: m,
	}
}

type registry[T any] struct {
	m   map[string]T
	mux sync.RWMutex
}

func (inst *registry[T]) Add(id string, t T) {
	inst.mux.Lock()
	defer inst.mux.Unlock()
	inst.m[id] = t
}

func (inst *registry[T]) Delete(id string) error {
	inst.mux.Lock()
	defer inst.mux.Unlock()

	if _, ok := inst.m[id]; ok {
		delete(inst.m, id)
		return nil
	}

	return fmt.Errorf("element not found in registry: %s", id)
}

func (inst *registry[T]) Get(id string) (T, bool) {
	inst.mux.RLock()
	defer inst.mux.RUnlock()
	t, ok := inst.m[id]
	return t, ok
}

func (inst *registry[T]) GetAll() map[string]T {
	inst.mux.RLock()
	defer inst.mux.RUnlock()
	return inst.m
}

func (inst *registry[T]) GetOrNil(id string) T {
	inst.mux.RLock()
	defer inst.mux.RUnlock()
	return inst.m[id]
}
