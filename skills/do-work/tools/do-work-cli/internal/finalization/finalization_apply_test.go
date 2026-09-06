package finalization

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMatchingHeadCommitRejectsMergeWithDisallowedPaths(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "feature.txt", "initial\n")
	writeFinalizationFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "initial")
	preparedHead := currentHead(repositoryRoot)
	baseBranch := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "branch", "--show-current"))

	// Branch 2: branched from preparedHead, updates unrelated.txt
	runFinalizationGit(t, repositoryRoot, "checkout", "-qb", "unrelated-branch", preparedHead)
	writeFinalizationFile(t, repositoryRoot, "unrelated.txt", "unrelated change\n")
	runFinalizationGit(t, repositoryRoot, "add", "unrelated.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "unrelated commit")

	// Back on base branch (preparedHead), create merge commit that merges unrelated-branch
	// and also carries the feature.txt change in the merge
	runFinalizationGit(t, repositoryRoot, "checkout", "-q", baseBranch)
	runFinalizationGit(t, repositoryRoot, "merge", "--no-ff", "--no-commit", "unrelated-branch")
	writeFinalizationFile(t, repositoryRoot, "feature.txt", "feature change\n")
	runFinalizationGit(t, repositoryRoot, "add", "feature.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "merge with feature change")
	mergeCommit := currentHead(repositoryRoot)

	// Calculate binary diff of feature.txt between preparedHead and mergeCommit
	diffOutput, err := exec.Command("git", "-C", repositoryRoot, "diff", "--binary", preparedHead, mergeCommit, "--", "feature.txt").Output()
	if err != nil {
		t.Fatal(err)
	}
	preparedDiffSHA := digestBytes(diffOutput)

	journal := &Journal{
		PreparedHead:         preparedHead,
		EffectiveCommitPaths: []string{"feature.txt"},
		PreparedDiffSHA256:   preparedDiffSHA,
	}

	// Because mergeCommit also touched unrelated.txt (brought in by merging unrelated-branch),
	// it must NOT be accepted as an exact match for EffectiveCommitPaths.
	matchedCommit, matched := matchingHeadCommit(repositoryRoot, journal)
	if matched {
		t.Fatalf("matchingHeadCommit accepted merge commit %s that touched disallowed paths", matchedCommit)
	}
}

func TestMatchingHeadCommitAcceptsCleanMergeMatchingEffectivePaths(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "part1.txt", "initial\n")
	writeFinalizationFile(t, repositoryRoot, "part2.txt", "initial\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "initial")
	preparedHead := currentHead(repositoryRoot)
	baseBranch := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "branch", "--show-current"))

	// Branch 1: update part1.txt
	writeFinalizationFile(t, repositoryRoot, "part1.txt", "part1 change\n")
	runFinalizationGit(t, repositoryRoot, "add", "part1.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "part1 commit")

	// Branch 2: branched from preparedHead, updates part2.txt
	runFinalizationGit(t, repositoryRoot, "checkout", "-qb", "feature-branch", preparedHead)
	writeFinalizationFile(t, repositoryRoot, "part2.txt", "part2 change\n")
	runFinalizationGit(t, repositoryRoot, "add", "part2.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "part2 commit")

	// Back on base branch, merge feature-branch with no-ff
	runFinalizationGit(t, repositoryRoot, "checkout", "-q", baseBranch)
	runFinalizationGit(t, repositoryRoot, "merge", "--no-ff", "-m", "clean merge", "feature-branch")
	mergeCommit := currentHead(repositoryRoot)

	diffOutput, err := exec.Command("git", "-C", repositoryRoot, "diff", "--binary", preparedHead, mergeCommit, "--", "part1.txt", "part2.txt").Output()
	if err != nil {
		t.Fatal(err)
	}
	preparedDiffSHA := digestBytes(diffOutput)

	journal := &Journal{
		PreparedHead:         preparedHead,
		EffectiveCommitPaths: []string{"part1.txt", "part2.txt"},
		PreparedDiffSHA256:   preparedDiffSHA,
	}

	matchedCommit, matched := matchingHeadCommit(repositoryRoot, journal)
	if !matched || matchedCommit != mergeCommit {
		t.Fatalf("matchingHeadCommit = (%q, %v), want (%q, true)", matchedCommit, matched, mergeCommit)
	}
}
