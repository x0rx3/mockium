package registry

import (
	"fmt"
	"sync"
)

func New[T any](m map[string]T) *Registry[T] {
	return &Registry[T]{
		m: m,
	}
}

type Registrator[T any] interface {
	Add(string, T)
	Delete(string) error
	Get(string) (T, bool)
	GetAll() map[string]T
}

type Registry[T any] struct {
	m   map[string]T
	mux sync.RWMutex
}

func (inst *Registry[T]) Add(id string, t T) {
	inst.mux.Lock()
	defer inst.mux.Unlock()
	inst.m[id] = t
}

func (inst *Registry[T]) Delete(id string) error {
	inst.mux.Lock()
	defer inst.mux.Unlock()

	if _, ok := inst.m[id]; ok {
		delete(inst.m, id)
		return nil
	}

	return fmt.Errorf("element not found in registry: %s", id)
}

func (inst *Registry[T]) Get(id string) (T, bool) {
	inst.mux.RLock()
	defer inst.mux.RUnlock()
	t, ok := inst.m[id]
	return t, ok
}

func (inst *Registry[T]) GetAll() map[string]T {
	inst.mux.RLock()
	defer inst.mux.RUnlock()
	return inst.m
}
