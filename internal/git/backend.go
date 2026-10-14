package git

import (
	"context"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// StatusOptions configures the status command.
type StatusOptions struct {
	// IncludeIgnored includes ignored files in the result.
	IncludeIgnored bool
}

// LogOptions configures the log command.
type LogOptions struct {
	MaxCount int
	Skip     int
	Since    string
	Until    string
	Author   string
	Grep     string
	All      bool
	Graph    bool
	Oneline  bool
	Ref      string   // specific ref to start from
	Paths    []string // file paths to filter
}

// DiffOptions configures the diff command.
type DiffOptions struct {
	Cached       bool   // --cached (staged changes)
	Ref          string // diff against a ref
	RefTo        string // end of range (Ref..RefTo)
	Paths        []string
	ContextLines int // -U<n>
	NameOnly     bool
	Stat         bool
}

// BranchOptions configures branch operations.
type BranchOptions struct {
	All     bool
	Remotes bool
	Delete  string
	Create  string
	Rename  string
}

// CommitOptions configures the commit command.
type CommitOptions struct {
	Message    string
	Amend      bool
	AllowEmpty bool
}

// CheckoutOptions configures the checkout command.
type CheckoutOptions struct {
	Create bool // -b flag
}

// PushOptions configures the push command.
type PushOptions struct {
	Force       bool
	SetUpstream bool
	Tags        bool
}

// PullOptions configures the pull command.
type PullOptions struct {
	Rebase bool
}

// FetchOptions configures the fetch command.
type FetchOptions struct {
	All   bool
	Prune bool
	Tags  bool
}

// StashAction represents a stash operation.
type StashAction int

const (
	StashPush StashAction = iota
	StashPop
	StashApply
	StashDrop
	StashList
	StashShow
)

// StashOptions configures stash operations.
type StashOptions struct {
	Action  StashAction
	Message string
	Index   int // for pop/apply/drop/show
}

// TagOptions configures tag operations.
type TagOptions struct {
	List      bool
	Create    string
	Delete    string
	Message   string // for annotated tags
	Annotated bool
	Ref       string // what to tag (defaults to HEAD)
}

// Backend defines the interface for git operations.
// All methods accept context.Context for cancellation.
type Backend interface {
	// RepoRoot returns the root directory of the current repository.
	RepoRoot(ctx context.Context) (string, error)

	// Status returns the current working tree status.
	Status(ctx context.Context, opts StatusOptions) ([]core.StatusEntry, error)

	// Log returns commit history.
	Log(ctx context.Context, opts LogOptions) ([]core.Commit, error)

	// Show returns detailed information about a commit.
	Show(ctx context.Context, ref string) (*core.CommitDetail, error)

	// Diff returns parsed diff output.
	Diff(ctx context.Context, opts DiffOptions) (*core.DiffResult, error)

	// DiffRaw returns raw diff text (for patch application).
	DiffRaw(ctx context.Context, opts DiffOptions) (string, error)

	// Branches returns branch information.
	Branches(ctx context.Context, opts BranchOptions) ([]core.Branch, error)

	// CurrentBranch returns the current branch name.
	CurrentBranch(ctx context.Context) (string, error)

	// Checkout switches branches or restores working tree files.
	Checkout(ctx context.Context, ref string, opts CheckoutOptions) error

	// Add stages files.
	Add(ctx context.Context, paths []string) error

	// ApplyPatch applies a patch to the index or working tree.
	ApplyPatch(ctx context.Context, patch string, cached bool) error

	// Commit creates a new commit.
	Commit(ctx context.Context, opts CommitOptions) (string, error)

	// Remotes returns configured remotes.
	Remotes(ctx context.Context) ([]core.Remote, error)

	// Fetch fetches from a remote.
	Fetch(ctx context.Context, remote string, opts FetchOptions) error

	// Pull pulls from a remote.
	Pull(ctx context.Context, remote string, branch string, opts PullOptions) error

	// Push pushes to a remote.
	Push(ctx context.Context, remote string, branch string, opts PushOptions) error

	// Stash performs stash operations.
	Stash(ctx context.Context, opts StashOptions) ([]core.StashEntry, error)

	// Tags returns tag information.
	Tags(ctx context.Context, opts TagOptions) ([]core.Tag, error)

	// CreateTag creates a new tag.
	CreateTag(ctx context.Context, opts TagOptions) error

	// DeleteTag deletes a tag.
	DeleteTag(ctx context.Context, name string) error

	// RemoteAdd adds a remote.
	RemoteAdd(ctx context.Context, name, url string) error

	// RemoteRemove removes a remote.
	RemoteRemove(ctx context.Context, name string) error

	// Init initializes a new repository.
	Init(ctx context.Context, path string, bare bool) error

	// Clone clones a repository.
	Clone(ctx context.Context, url, path string) error

	// Version returns the git version string.
	Version(ctx context.Context) (string, error)

	// ConfigGet returns a git config value.
	ConfigGet(ctx context.Context, key string, global bool) (string, error)

	// ConfigSet sets a git config value.
	ConfigSet(ctx context.Context, key, value string, global bool) error

	// ConfigList returns all git config entries.
	ConfigList(ctx context.Context, global bool) (map[string]string, error)
}
