package core

import "time"

// StatusCode represents the status of a file in the working tree or index.
type StatusCode byte

const (
	StatusUnmodified  StatusCode = ' '
	StatusModified    StatusCode = 'M'
	StatusTypeChanged StatusCode = 'T'
	StatusAdded       StatusCode = 'A'
	StatusDeleted     StatusCode = 'D'
	StatusRenamed     StatusCode = 'R'
	StatusCopied      StatusCode = 'C'
	StatusUnmerged    StatusCode = 'U'
	StatusUntracked   StatusCode = '?'
	StatusIgnored     StatusCode = '!'
)

func (s StatusCode) String() string { return string(s) }

// StatusEntry represents a single file's status.
type StatusEntry struct {
	Index     StatusCode
	Worktree  StatusCode
	Path      string
	OrigPath  string // non-empty for renames/copies
	Submodule string
}

// IsStaged returns true if the file has changes in the index.
func (e StatusEntry) IsStaged() bool {
	return e.Index != StatusUnmodified && e.Index != StatusUntracked && e.Index != StatusIgnored
}

// IsUnstaged returns true if the file has changes in the worktree.
func (e StatusEntry) IsUnstaged() bool {
	return e.Worktree != StatusUnmodified && e.Worktree != StatusUntracked && e.Worktree != StatusIgnored
}

// Commit represents a parsed git commit.
type Commit struct {
	Hash         string
	AbbrevHash   string
	ParentHashes []string
	TreeHash     string
	Author       Person
	Committer    Person
	Subject      string
	Body         string
	Refs         []Ref
}

// Person represents a git author or committer.
type Person struct {
	Name  string
	Email string
	Date  time.Time
}

// Ref represents a git reference (branch, tag, etc.).
type Ref struct {
	Name   string
	Type   RefType
	Remote string // non-empty for remote refs
}

// RefType categorizes a git reference.
type RefType int

const (
	RefTypeBranch RefType = iota
	RefTypeRemoteBranch
	RefTypeTag
	RefTypeHEAD
	RefTypeStash
)

func (t RefType) String() string {
	switch t {
	case RefTypeBranch:
		return "branch"
	case RefTypeRemoteBranch:
		return "remote-branch"
	case RefTypeTag:
		return "tag"
	case RefTypeHEAD:
		return "HEAD"
	case RefTypeStash:
		return "stash"
	default:
		return "unknown"
	}
}

// Branch represents a git branch.
type Branch struct {
	Name     string
	Hash     string
	Upstream string
	Ahead    int
	Behind   int
	Current  bool
}

// Remote represents a git remote.
type Remote struct {
	Name     string
	FetchURL string
	PushURL  string
}

// StashEntry represents a single stash entry.
type StashEntry struct {
	Index   int
	Message string
	Branch  string
}

// Tag represents a git tag.
type Tag struct {
	Name        string
	Hash        string
	Message     string // non-empty for annotated tags
	Tagger      *Person
	IsAnnotated bool
}

// DiffResult contains parsed diff output.
type DiffResult struct {
	Files []DiffFile
}

// DiffFile represents changes to a single file.
type DiffFile struct {
	OldName string
	NewName string
	Status  DiffStatus
	Binary  bool
	Hunks   []DiffHunk
	OldMode string
	NewMode string
}

// DiffStatus describes how a file changed.
type DiffStatus int

const (
	DiffStatusModified DiffStatus = iota
	DiffStatusAdded
	DiffStatusDeleted
	DiffStatusRenamed
	DiffStatusCopied
	DiffStatusTypeChanged
)

func (d DiffStatus) String() string {
	switch d {
	case DiffStatusModified:
		return "modified"
	case DiffStatusAdded:
		return "added"
	case DiffStatusDeleted:
		return "deleted"
	case DiffStatusRenamed:
		return "renamed"
	case DiffStatusCopied:
		return "copied"
	case DiffStatusTypeChanged:
		return "type-changed"
	default:
		return "unknown"
	}
}

// DiffHunk represents a single hunk in a diff.
type DiffHunk struct {
	OldStart int
	OldLines int
	NewStart int
	NewLines int
	Header   string
	Lines    []DiffLine
}

// DiffLine represents a single line in a diff hunk.
type DiffLine struct {
	Type    DiffLineType
	Content string
}

// DiffLineType categorizes a diff line.
type DiffLineType int

const (
	DiffLineContext DiffLineType = iota
	DiffLineAdded
	DiffLineDeleted
)

// CommitDetail extends Commit with full show information.
type CommitDetail struct {
	Commit
	Diff *DiffResult
	Stat []FileStat
}

// FileStat represents diffstat for a single file.
type FileStat struct {
	File      string
	Additions int
	Deletions int
}
