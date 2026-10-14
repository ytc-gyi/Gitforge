package core

import "context"

// Version information, injected at build time.
var (
	BuildVersion = "dev"
	BuildCommit  = "unknown"
	BuildDate    = "unknown"
)

// Engine is the central facade for GitForge.
// All front ends (CLI, TUI, Web) interact with the system through this interface.
type Engine struct {
	Events   *EventBus
	Modules  *ModuleRegistry
	repoPath string
}

// NewEngine creates a new Engine.
func NewEngine() *Engine {
	return &Engine{
		Events:  NewEventBus(),
		Modules: NewModuleRegistry(),
	}
}

// SetRepoPath sets the repository path for the engine.
func (e *Engine) SetRepoPath(path string) {
	e.repoPath = path
}

// RepoPath returns the current repository path.
func (e *Engine) RepoPath() string {
	return e.repoPath
}

// Init initializes the engine and all registered modules.
func (e *Engine) Init(ctx context.Context) error {
	return e.Modules.InitAll(ctx)
}

// Close shuts down the engine and all modules.
func (e *Engine) Close() error {
	return e.Modules.CloseAll()
}
