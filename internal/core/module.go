package core

import (
	"context"
	"fmt"
	"sync"
)

// Module represents a feature module (internal or plugin).
type Module interface {
	// Name returns the unique name of the module.
	Name() string
	// Init initializes the module with the engine context.
	Init(ctx context.Context) error
	// Close cleans up module resources.
	Close() error
}

// ModuleRegistry manages registered modules.
type ModuleRegistry struct {
	mu      sync.RWMutex
	modules map[string]Module
}

// NewModuleRegistry creates a new ModuleRegistry.
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry.
func (r *ModuleRegistry) Register(m Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.modules[m.Name()]; exists {
		return fmt.Errorf("module %q already registered", m.Name())
	}
	r.modules[m.Name()] = m
	return nil
}

// Get returns a module by name.
func (r *ModuleRegistry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.modules[name]
	return m, ok
}

// All returns all registered modules.
func (r *ModuleRegistry) All() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Module, 0, len(r.modules))
	for _, m := range r.modules {
		result = append(result, m)
	}
	return result
}

// InitAll initializes all registered modules.
func (r *ModuleRegistry) InitAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name, m := range r.modules {
		if err := m.Init(ctx); err != nil {
			return fmt.Errorf("initializing module %q: %w", name, err)
		}
	}
	return nil
}

// CloseAll closes all registered modules.
func (r *ModuleRegistry) CloseAll() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var firstErr error
	for _, m := range r.modules {
		if err := m.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
