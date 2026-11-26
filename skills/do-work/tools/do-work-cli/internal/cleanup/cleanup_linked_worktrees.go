package cleanup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// Include every secondary checkout, plus unattached branches in builder namespaces.
// Keep the primary checkout, main, and the checkout running this command.
func cleanupCandidates(ctx context.Context, root, baseName string, trees []worktreeRecord) ([]worktreeRecord, error) {
	caller, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	here, _ := cleanupGit(ctx, ".", "rev-parse", "--show-toplevel")
	current, _ := os.Stat(strings.TrimSuffix(here, "\n"))
	attached := map[string]int{}
	for _, tree := range trees {
		// Rebasing or bisecting can detach HEAD while the original branch is
		// still in use. Keep both the checkout and that branch.
		dir, err := tree.gitDir(ctx, root)
		if err != nil {
			return nil, err
		}
		for _, marker := range []string{"rebase-merge/head-name", "rebase-apply/head-name", "BISECT_START"} {
			data, err := os.ReadFile(filepath.Join(dir, marker))
			if err == nil {
				name := strings.TrimPrefix(strings.TrimSpace(string(data)), "refs/heads/")
				attached[name]++
			} else if !os.IsNotExist(err) {
				return nil, err
			}
		}

		if !tree.Detached {
			attached[tree.Name]++
		}
	}
	var candidates []worktreeRecord
	for i, tree := range trees {
		info, _ := os.Stat(tree.Path)
		if i == 0 || os.SameFile(info, caller) || os.SameFile(info, current) || tree.Name == baseName || attached[tree.Name] > 1 {
			continue
		}
		candidates = append(candidates, tree)
	}
	output, err := cleanupGit(ctx, root, "for-each-ref", "--format=%(objectname) %(refname)", "refs/heads/t3/", "refs/heads/worktree-agent-*")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		head, name, ok := strings.Cut(line, " refs/heads/")
		if ok && attached[name] == 0 && name != baseName {
			candidates = append(candidates, worktreeRecord{Name: name, Head: head})
		}
	}
	return candidates, nil
}

// A squash or cherry-pick can deliver a branch's files without its history.
func cleanupMergeProof(ctx context.Context, root, base, head string) (bool, string) {
	if head == "" {
		return false, "missing commit"
	}
	if gitExitSuccess(ctx, root, "merge-base", "--is-ancestor", head, base) {
		return true, "already merged"
	}
	tree, err := cleanupGit(ctx, root, "rev-parse", "--verify", base+"^{tree}")
	if err != nil {
		return false, "could not read integration tree"
	}
	merged, err := cleanupGit(ctx, root, "merge-tree", "--write-tree", base, head)
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return false, "merge conflicts"
		}
		return false, "merge preview unavailable"
	}
	if strings.TrimSpace(merged) != strings.TrimSpace(tree) {
		return false, "merge would change files"
	}
	return true, "merge would leave files unchanged"
}

func removeCleanupCandidate(ctx context.Context, root string, tree worktreeRecord, links []string, baseRef, baseHead string, force bool) error {
	if tree.Path != "" {
		listing, err := cleanupGitBytes(ctx, root, "worktree", "list", "--porcelain", "-z")
		if err != nil {
			return err
		}
		found := false
		for _, live := range allWorktrees(listing) {
			if live == tree {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("worktree changed since discovery: %s", tree.Path)
		}
		baseNow, err := cleanupGit(ctx, root, "rev-parse", "--verify", baseRef)
		if err != nil || strings.TrimSpace(baseNow) != baseHead {
			return fmt.Errorf("integration branch changed before removing %s", tree.Path)
		}
		args := []string{"worktree", "remove"}
		if force {
			args = append(args, "--force")
		} else {
			if !tree.safeCheckout(ctx, root, links) {
				return fmt.Errorf("worktree gained local changes: %s", tree.Path)
			}
			for _, link := range links {
				if err := os.Remove(filepath.Join(tree.Path, filepath.FromSlash(link))); err != nil {
					return err
				}
			}
		}
		if _, err := cleanupGit(ctx, root, append(args, "--", tree.Path)...); err != nil {
			return err
		}
	}
	if tree.Detached {
		return nil
	}
	ref := "refs/heads/" + tree.Name
	occupied, err := cleanupGit(ctx, root, "for-each-ref", "--format=%(refname) %(worktreepath)", ref)
	if err != nil {
		return err
	}
	if strings.TrimSuffix(occupied, "\n") != ref+" " {
		return fmt.Errorf("branch missing or checked out: %s", tree.Name)
	}
	// Compare and delete. A branch advance or a change of base since the proof
	// must leave the branch intact, including for a content-equivalent squash.
	cmd := exec.CommandContext(ctx, "git", "-C", root, "update-ref", "--no-deref", "--stdin")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("verify %s %s\ndelete %s %s\n", baseRef, baseHead, ref, tree.Head))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("delete %s: %w: %s", ref, err, out)
	}
	// Match porcelain branch deletion by clearing this branch's local tracking
	// configuration, if any, so a later branch of the same name starts clean.
	section := "branch." + tree.Name
	probe := exec.CommandContext(ctx, "git", "-C", root, "config", "--local", "--get-regexp", "^"+regexp.QuoteMeta(section)+`\.`)
	if err := probe.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return nil
		}
		return err
	}
	_, err = cleanupGit(ctx, root, "config", "--local", "--remove-section", section)
	return err
}

// A missing linked directory can still have staged files or a rebase in its
// private Git directory. Find the exact registration before pruning it.
func (tree worktreeRecord) gitDir(ctx context.Context, root string) (string, error) {
	if tree.Prunable {
		common, err := cleanupGit(ctx, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
		if err != nil {
			return "", err
		}
		parent := filepath.Join(strings.TrimSuffix(common, "\n"), "worktrees")
		entries, err := os.ReadDir(parent)
		if err != nil {
			return "", err
		}
		for _, entry := range entries {
			dir := filepath.Join(parent, entry.Name())
			link, err := os.ReadFile(filepath.Join(dir, "gitdir"))
			if err == nil && strings.TrimSuffix(string(link), "\n") == filepath.Join(tree.Path, ".git") {
				return dir, nil
			}
		}
		return "", fmt.Errorf("could not locate Git metadata for %s", tree.Path)
	}
	dir, err := cleanupGit(ctx, tree.Path, "rev-parse", "--absolute-git-dir")
	return strings.TrimSuffix(dir, "\n"), err
}

func (tree worktreeRecord) safeCheckout(ctx context.Context, root string, ownedLinks []string) bool {
	if tree.Path == "" {
		return true
	}
	dir, err := tree.gitDir(ctx, root)
	if err != nil {
		return false
	}
	for _, marker := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "BISECT_START", "rebase-apply", "rebase-merge", "sequencer", "index.lock"} {
		if _, err := os.Lstat(filepath.Join(dir, marker)); !os.IsNotExist(err) {
			return false
		}
	}
	if tree.Prunable {
		_, err := os.Lstat(tree.Path)
		return os.IsNotExist(err) && gitExitSuccess(ctx, root, "--git-dir="+dir, "diff-index", "--quiet", "--cached", tree.Head, "--")
	}
	dirty, err := worktreeDirtyPaths(ctx, tree.Path, ownedLinks)
	if err != nil || len(dirty) > 0 {
		return false
	}
	ignored, err := cleanupGit(ctx, tree.Path, "ls-files", "-z", "--others", "--ignored", "--exclude-standard", "--directory")
	if err != nil {
		return false
	}
	for _, file := range splitZero(ignored) {
		if !slices.Contains(ownedLinks, file) {
			return false
		}
	}
	return true
}
