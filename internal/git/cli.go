package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/phantom1785322-hub/gitforge/internal/core"
)

// CLIBackend implements Backend by shelling out to the git CLI.
type CLIBackend struct {
	gitPath string
	repoDir string
}

// NewCLIBackend creates a Backend that uses the git binary found in PATH.
// repoDir may be empty; when set, every command runs inside that directory.
func NewCLIBackend(repoDir string) (*CLIBackend, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, core.ErrGitNotFound
	}
	return &CLIBackend{
		gitPath: gitPath,
		repoDir: repoDir,
	}, nil
}

// run executes a git command and returns stdout. On non-zero exit it wraps the
// error together with stderr into a *core.GitError.
func (c *CLIBackend) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c.gitPath, args...)
	if c.repoDir != "" {
		cmd.Dir = c.repoDir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		cmdName := ""
		if len(args) > 0 {
			cmdName = args[0]
		}
		return "", &core.GitError{
			Command: cmdName,
			Args:    args,
			Stderr:  strings.TrimSpace(stderr.String()),
			Err:     err,
		}
	}
	return stdout.String(), nil
}

// runWithStdin executes a git command, feeding data into stdin.
func (c *CLIBackend) runWithStdin(ctx context.Context, stdin string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c.gitPath, args...)
	if c.repoDir != "" {
		cmd.Dir = c.repoDir
	}
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		cmdName := ""
		if len(args) > 0 {
			cmdName = args[0]
		}
		return "", &core.GitError{
			Command: cmdName,
			Args:    args,
			Stderr:  strings.TrimSpace(stderr.String()),
			Err:     err,
		}
	}
	return stdout.String(), nil
}

// ValidateRef rejects ref names that start with '-' to prevent argument injection.
func ValidateRef(ref string) error {
	return validateRef(ref)
}

func validateRef(ref string) error {
	if strings.HasPrefix(ref, "-") {
		return fmt.Errorf("invalid ref name %q: must not start with '-'", ref)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Backend implementation
// ---------------------------------------------------------------------------

func (c *CLIBackend) RepoRoot(ctx context.Context) (string, error) {
	out, err := c.run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *CLIBackend) Status(ctx context.Context, opts StatusOptions) ([]core.StatusEntry, error) {
	args := []string{"status", "--porcelain=v2", "-z"}
	if opts.IncludeIgnored {
		args = append(args, "--ignored")
	}
	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return ParseStatusPorcelainV2(out)
}

func (c *CLIBackend) Log(ctx context.Context, opts LogOptions) ([]core.Commit, error) {
	// Use record separator \x01 between commits, NUL between fields.
	format := LogFormatParseable
	args := []string{"log", "--format=format:" + format + "%x01"}

	if opts.MaxCount > 0 {
		args = append(args, fmt.Sprintf("-n%d", opts.MaxCount))
	}
	if opts.Skip > 0 {
		args = append(args, fmt.Sprintf("--skip=%d", opts.Skip))
	}
	if opts.Since != "" {
		args = append(args, "--since="+opts.Since)
	}
	if opts.Until != "" {
		args = append(args, "--until="+opts.Until)
	}
	if opts.Author != "" {
		args = append(args, "--author="+opts.Author)
	}
	if opts.Grep != "" {
		args = append(args, "--grep="+opts.Grep)
	}
	if opts.All {
		args = append(args, "--all")
	}
	if opts.Ref != "" {
		if err := validateRef(opts.Ref); err != nil {
			return nil, err
		}
		args = append(args, opts.Ref)
	}
	if len(opts.Paths) > 0 {
		args = append(args, "--")
		args = append(args, opts.Paths...)
	}

	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return ParseLog(out)
}

func (c *CLIBackend) Show(ctx context.Context, ref string) (*core.CommitDetail, error) {
	if err := validateRef(ref); err != nil {
		return nil, err
	}

	// Get commit metadata.
	format := LogFormatParseable
	metaOut, err := c.run(ctx, "show", "--no-patch", "--format=format:"+format, ref)
	if err != nil {
		return nil, err
	}

	commits, err := ParseLog(metaOut)
	if err != nil || len(commits) == 0 {
		return nil, fmt.Errorf("failed to parse show output for %q", ref)
	}

	// Get diff.
	diffOut, err := c.run(ctx, "show", "--format=", "--patch", ref)
	if err != nil {
		return nil, err
	}
	diffResult, err := ParseDiff(diffOut)
	if err != nil {
		return nil, err
	}

	// Get stat.
	statOut, err := c.run(ctx, "show", "--format=", "--stat", ref)
	if err != nil {
		return nil, err
	}

	return &core.CommitDetail{
		Commit: commits[0],
		Diff:   diffResult,
		Stat:   ParseDiffStat(statOut),
	}, nil
}

func (c *CLIBackend) Diff(ctx context.Context, opts DiffOptions) (*core.DiffResult, error) {
	raw, err := c.DiffRaw(ctx, opts)
	if err != nil {
		return nil, err
	}
	return ParseDiff(raw)
}

func (c *CLIBackend) DiffRaw(ctx context.Context, opts DiffOptions) (string, error) {
	args := []string{"diff"}
	if opts.Cached {
		args = append(args, "--cached")
	}
	if opts.ContextLines > 0 {
		args = append(args, fmt.Sprintf("-U%d", opts.ContextLines))
	}
	if opts.NameOnly {
		args = append(args, "--name-only")
	}
	if opts.Stat {
		args = append(args, "--stat")
	}
	if opts.Ref != "" {
		if err := validateRef(opts.Ref); err != nil {
			return "", err
		}
		if opts.RefTo != "" {
			if err := validateRef(opts.RefTo); err != nil {
				return "", err
			}
			args = append(args, opts.Ref+".."+opts.RefTo)
		} else {
			args = append(args, opts.Ref)
		}
	}
	if len(opts.Paths) > 0 {
		args = append(args, "--")
		args = append(args, opts.Paths...)
	}
	return c.run(ctx, args...)
}

func (c *CLIBackend) Branches(ctx context.Context, opts BranchOptions) ([]core.Branch, error) {
	// Handle mutations first.
	if opts.Delete != "" {
		if err := validateRef(opts.Delete); err != nil {
			return nil, err
		}
		_, err := c.run(ctx, "branch", "-d", opts.Delete)
		return nil, err
	}
	if opts.Create != "" {
		if err := validateRef(opts.Create); err != nil {
			return nil, err
		}
		_, err := c.run(ctx, "branch", opts.Create)
		return nil, err
	}
	if opts.Rename != "" {
		if err := validateRef(opts.Rename); err != nil {
			return nil, err
		}
		_, err := c.run(ctx, "branch", "-m", opts.Rename)
		return nil, err
	}

	// List branches.
	format := "%(refname:short)%00%(objectname:short)%00%(upstream:short)%00%(upstream:track,nobracket)%00%(HEAD)"
	args := []string{"branch", "--format=" + format}
	if opts.All {
		args = append(args, "-a")
	}
	if opts.Remotes {
		args = append(args, "-r")
	}

	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return ParseBranches(out)
}

func (c *CLIBackend) CurrentBranch(ctx context.Context) (string, error) {
	out, err := c.run(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(out)
	if branch == "HEAD" {
		return "", core.ErrDetachedHead
	}
	return branch, nil
}

func (c *CLIBackend) Checkout(ctx context.Context, ref string, opts CheckoutOptions) error {
	if err := validateRef(ref); err != nil {
		return err
	}
	args := []string{"checkout"}
	if opts.Create {
		args = append(args, "-b")
	}
	args = append(args, ref)
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Add(ctx context.Context, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	args := []string{"add", "--"}
	args = append(args, paths...)
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) ApplyPatch(ctx context.Context, patch string, cached bool) error {
	args := []string{"apply"}
	if cached {
		args = append(args, "--cached")
	}
	args = append(args, "-")
	_, err := c.runWithStdin(ctx, patch, args...)
	return err
}

func (c *CLIBackend) Commit(ctx context.Context, opts CommitOptions) (string, error) {
	args := []string{"commit"}
	if opts.Message != "" {
		args = append(args, "-m", opts.Message)
	}
	if opts.Amend {
		args = append(args, "--amend")
	}
	if opts.AllowEmpty {
		args = append(args, "--allow-empty")
	}

	_, err := c.run(ctx, args...)
	if err != nil {
		return "", err
	}

	// Return the hash of the new commit.
	hashOut, err := c.run(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(hashOut), nil
}

func (c *CLIBackend) Remotes(ctx context.Context) ([]core.Remote, error) {
	out, err := c.run(ctx, "remote", "-v")
	if err != nil {
		return nil, err
	}
	return parseRemotes(out), nil
}

func parseRemotes(data string) []core.Remote {
	// Each line: <name>\t<url> (fetch|push)
	remoteMap := make(map[string]*core.Remote)
	var order []string

	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		name := parts[0]
		url := parts[1]
		kind := parts[2] // "(fetch)" or "(push)"

		r, ok := remoteMap[name]
		if !ok {
			r = &core.Remote{Name: name}
			remoteMap[name] = r
			order = append(order, name)
		}
		switch kind {
		case "(fetch)":
			r.FetchURL = url
		case "(push)":
			r.PushURL = url
		}
	}

	remotes := make([]core.Remote, 0, len(order))
	for _, name := range order {
		remotes = append(remotes, *remoteMap[name])
	}
	return remotes
}

func (c *CLIBackend) Fetch(ctx context.Context, remote string, opts FetchOptions) error {
	args := []string{"fetch"}
	if opts.All {
		args = append(args, "--all")
	} else if remote != "" {
		if err := validateRef(remote); err != nil {
			return err
		}
		args = append(args, remote)
	}
	if opts.Prune {
		args = append(args, "--prune")
	}
	if opts.Tags {
		args = append(args, "--tags")
	}
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Pull(ctx context.Context, remote string, branch string, opts PullOptions) error {
	args := []string{"pull"}
	if opts.Rebase {
		args = append(args, "--rebase")
	}
	if remote != "" {
		if err := validateRef(remote); err != nil {
			return err
		}
		args = append(args, remote)
		if branch != "" {
			if err := validateRef(branch); err != nil {
				return err
			}
			args = append(args, branch)
		}
	}
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Push(ctx context.Context, remote string, branch string, opts PushOptions) error {
	args := []string{"push"}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.SetUpstream {
		args = append(args, "--set-upstream")
	}
	if opts.Tags {
		args = append(args, "--tags")
	}
	if remote != "" {
		if err := validateRef(remote); err != nil {
			return err
		}
		args = append(args, remote)
		if branch != "" {
			if err := validateRef(branch); err != nil {
				return err
			}
			args = append(args, branch)
		}
	}
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Stash(ctx context.Context, opts StashOptions) ([]core.StashEntry, error) {
	switch opts.Action {
	case StashPush:
		args := []string{"stash", "push"}
		if opts.Message != "" {
			args = append(args, "-m", opts.Message)
		}
		_, err := c.run(ctx, args...)
		return nil, err

	case StashPop:
		ref := fmt.Sprintf("stash@{%d}", opts.Index)
		_, err := c.run(ctx, "stash", "pop", ref)
		return nil, err

	case StashApply:
		ref := fmt.Sprintf("stash@{%d}", opts.Index)
		_, err := c.run(ctx, "stash", "apply", ref)
		return nil, err

	case StashDrop:
		ref := fmt.Sprintf("stash@{%d}", opts.Index)
		_, err := c.run(ctx, "stash", "drop", ref)
		return nil, err

	case StashList:
		out, err := c.run(ctx, "stash", "list", "--format=%gd%x00%gs")
		if err != nil {
			return nil, err
		}
		return parseStashList(out), nil

	case StashShow:
		ref := fmt.Sprintf("stash@{%d}", opts.Index)
		_, err := c.run(ctx, "stash", "show", ref)
		return nil, err

	default:
		return nil, fmt.Errorf("unknown stash action: %d", opts.Action)
	}
}

func parseStashList(data string) []core.StashEntry {
	var entries []core.StashEntry
	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\x00", 2)
		if len(parts) < 2 {
			continue
		}
		// parts[0] is like "stash@{0}"
		ref := parts[0]
		idx := 0
		// Extract index from stash@{N}
		if start := strings.Index(ref, "{"); start >= 0 {
			if end := strings.Index(ref, "}"); end > start {
				idx, _ = strconv.Atoi(ref[start+1 : end])
			}
		}
		entries = append(entries, core.StashEntry{
			Index:   idx,
			Message: parts[1],
			Branch:  "", // not available from this format
		})
	}
	return entries
}

func (c *CLIBackend) Tags(ctx context.Context, opts TagOptions) ([]core.Tag, error) {
	if opts.Delete != "" {
		return nil, c.DeleteTag(ctx, opts.Delete)
	}
	if opts.Create != "" {
		return nil, c.CreateTag(ctx, opts)
	}

	format := "%(refname:short)%00%(objectname:short)%00%(*objectname:short)%00%(objecttype)%00%(creatordate:iso)"
	args := []string{"tag", "-l", "--format=" + format, "--sort=-creatordate"}

	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return parseTags(out), nil
}

func parseTags(data string) []core.Tag {
	var tags []core.Tag
	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\x00")
		if len(fields) < 5 {
			continue
		}

		isAnnotated := fields[3] == "tag"

		tag := core.Tag{
			Name:        fields[0],
			Hash:        fields[1],
			IsAnnotated: isAnnotated,
		}

		// For annotated tags, the dereferenced target hash is in field 2.
		if isAnnotated && fields[2] != "" {
			tag.Hash = fields[2]
		}

		// Parse date (iso format: "2024-01-15 10:30:00 -0500")
		if dateStr := strings.TrimSpace(fields[4]); dateStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr); err == nil {
				tag.Tagger = &core.Person{Date: t}
			}
		}

		tags = append(tags, tag)
	}
	return tags
}

func (c *CLIBackend) CreateTag(ctx context.Context, opts TagOptions) error {
	if opts.Create == "" {
		return fmt.Errorf("tag name is required")
	}
	if err := validateRef(opts.Create); err != nil {
		return err
	}

	args := []string{"tag"}
	if opts.Annotated || opts.Message != "" {
		args = append(args, "-a")
		if opts.Message != "" {
			args = append(args, "-m", opts.Message)
		}
	}
	args = append(args, opts.Create)
	if opts.Ref != "" {
		if err := validateRef(opts.Ref); err != nil {
			return err
		}
		args = append(args, opts.Ref)
	}

	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) DeleteTag(ctx context.Context, name string) error {
	if err := validateRef(name); err != nil {
		return err
	}
	_, err := c.run(ctx, "tag", "-d", name)
	return err
}

func (c *CLIBackend) RemoteAdd(ctx context.Context, name, url string) error {
	if err := validateRef(name); err != nil {
		return err
	}
	_, err := c.run(ctx, "remote", "add", name, url)
	return err
}

func (c *CLIBackend) RemoteRemove(ctx context.Context, name string) error {
	if err := validateRef(name); err != nil {
		return err
	}
	_, err := c.run(ctx, "remote", "remove", name)
	return err
}

func (c *CLIBackend) Init(ctx context.Context, path string, bare bool) error {
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	if path != "" {
		args = append(args, path)
	}
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Clone(ctx context.Context, url, path string) error {
	args := []string{"clone", url}
	if path != "" {
		args = append(args, path)
	}
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) Version(ctx context.Context) (string, error) {
	out, err := c.run(ctx, "version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *CLIBackend) ConfigGet(ctx context.Context, key string, global bool) (string, error) {
	args := []string{"config"}
	if global {
		args = append(args, "--global")
	}
	args = append(args, "--get", key)
	out, err := c.run(ctx, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *CLIBackend) ConfigSet(ctx context.Context, key, value string, global bool) error {
	args := []string{"config"}
	if global {
		args = append(args, "--global")
	}
	args = append(args, key, value)
	_, err := c.run(ctx, args...)
	return err
}

func (c *CLIBackend) ConfigList(ctx context.Context, global bool) (map[string]string, error) {
	args := []string{"config", "--list"}
	if global {
		args = append(args, "--global")
	}
	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[parts[0]] = parts[1]
		}
	}
	return config, nil
}

// Compile-time interface check.
var _ Backend = (*CLIBackend)(nil)
