package toolboxcommands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

// secondOpinionWorktree owns its scratch directory before Git starts, including
// partially completed setup. Its caller defers cleanup before calling prepareBundle.
type secondOpinionWorktree struct {
	repositoryPath, scratchPath, treePath, bundlePath string
}

func secondOpinionWorktreeGit(ctx context.Context, repository string, args ...string) ([]byte, error) {
	// Git setup can spawn checkout/filter children too. Use the same reaping
	// boundary as the backend before inspecting or removing the worktree.
	capture := secondOpinionCapture(ctx, secondOpinionProgress{phase: "worktree", backend: "git"}, repository, "git", append([]string{"-c", "core.hooksPath=/dev/null"}, args...)...)
	if capture.err != nil {
		return []byte(capture.stdout), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), capture.err, strings.TrimSpace(capture.stderr))
	}
	return []byte(capture.stdout), nil
}

func (worktree *secondOpinionWorktree) prepareBundle(ctx context.Context, bundle string, before map[string]secondOpinionFile) error {
	var err error
	worktree.scratchPath, err = os.MkdirTemp("", "ai-report-second-opinion-worktree-")
	if err != nil {
		return err
	}
	canonical, err := filepath.EvalSymlinks(worktree.scratchPath)
	if err != nil {
		return err
	}
	worktree.scratchPath = canonical
	worktree.treePath = filepath.Join(worktree.scratchPath, "checkout")
	if _, err = secondOpinionWorktreeGit(ctx, worktree.repositoryPath, "worktree", "add", "--detach", worktree.treePath, "HEAD"); err != nil {
		return err
	}
	relative, err := filepath.Rel(worktree.repositoryPath, bundle)
	if err != nil {
		return err
	}
	worktree.bundlePath = filepath.Join(worktree.treePath, relative)
	root, err := os.OpenRoot(worktree.treePath)
	if err != nil {
		return err
	}
	defer root.Close()
	if err := validateNoLinkedAncestors(worktree.treePath, relative, false); err != nil {
		return err
	}
	// Replace, rather than overlay: tracked files deleted in the current bundle
	// must not reappear from HEAD. Rooted operations cannot follow a checkout link.
	if err := rootedfs.RemoveAll(root, relative); err != nil {
		return err
	}
	if err := rootedfs.MkdirAll(root, relative, 0755); err != nil {
		return err
	}
	names := make([]string, 0, len(before))
	for name := range before {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		item := before[name]
		destination := filepath.Join(relative, name)
		if item.info.IsDir() {
			err = rootedfs.MkdirAll(root, destination, 0700)
		} else {
			err = rootedPublishFile(worktree.treePath, destination, item.data, item.mode, false)
		}
		if err != nil {
			return err
		}
	}
	// Apply directory modes last so read-only directories can still be populated.
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if err := os.Chmod(filepath.Join(worktree.bundlePath, name), before[name].mode); err != nil {
			return err
		}
	}
	return nil
}

func (worktree *secondOpinionWorktree) cleanupWorktree() error {
	if worktree.scratchPath == "" {
		return nil
	}
	// Cancellation belongs to generation, not to teardown after the processes
	// have been reaped. Bound Git cleanup independently, even after setup fails.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := secondOpinionWorktreeGit(ctx, worktree.repositoryPath, "worktree", "list", "--porcelain", "-z")
	if err == nil {
		for _, field := range strings.Split(string(output), "\x00") {
			if strings.TrimPrefix(field, "worktree ") == worktree.treePath {
				_, err = secondOpinionWorktreeGit(ctx, worktree.repositoryPath, "worktree", "remove", "--force", "--force", worktree.treePath)
				break
			}
		}
	}
	return errors.Join(err, os.RemoveAll(worktree.scratchPath))
}

// secondOpinionImport keeps only new regular root files admitted by the staging
// and asset-prefix rules. Refusals may still retain these files, but can never
// import an edited original, an existing asset, a link, or an unrelated path.
// The stage is first: later publication errors leave its recovery bytes intact.
func secondOpinionImport(repository, bundle, isolated, stageName, prefix string, before map[string]secondOpinionFile) (bool, error) {
	if err := secondOpinionPreserved(bundle, before, nil, ""); err != nil {
		return false, err
	}
	entries, err := os.ReadDir(isolated)
	if err != nil {
		return false, err
	}
	names := []string{stageName}
	for _, entry := range entries {
		if _, existed := before[entry.Name()]; !existed && strings.HasPrefix(entry.Name(), prefix) && entry.Type().IsRegular() {
			names = append(names, entry.Name())
		}
	}
	relative, _ := filepath.Rel(repository, bundle)
	stageCopied := false
	for _, name := range names {
		info, err := os.Lstat(filepath.Join(isolated, name))
		if os.IsNotExist(err) && name == stageName {
			continue
		}
		if err != nil {
			return stageCopied, err
		}
		if !info.Mode().IsRegular() {
			return stageCopied, fmt.Errorf("copy-back source is not a regular file: %s", name)
		}
		if name == stageName && info.Size() == 0 {
			continue
		}
		data, err := rootedReadFile(isolated, name)
		if err == nil {
			err = rootedPublishFile(repository, filepath.Join(relative, name), data, info.Mode(), false)
		}
		if err != nil {
			return stageCopied, err
		}
		if name == stageName {
			stageCopied = true
		}
	}
	return stageCopied, nil
}
