package cleanup

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func EnrichDocumentationLinks(ctx context.Context, plan *CleanupPlan) {
	if plan == nil || len(plan.Groups) == 0 {
		return
	}
	command := exec.CommandContext(ctx, "git", "-C", plan.RepositoryRoot, "ls-files", "-z", "--", "*.md", ":(exclude)do-work/**")
	outputBytes, err := command.Output()
	if err != nil {
		plan.Findings = append(plan.Findings, manualFinding("CLEANUP-LINK-SCAN-FAILED", nil, nil, err.Error()))
		return
	}
	output := string(outputBytes)
	trackedPaths := splitZero(output)
	moveMap := map[string]string{}
	attachmentGroup := -1
	for groupIndex := range plan.Groups {
		for _, operation := range plan.Groups[groupIndex].Operations {
			if operation.Kind == OperationMove {
				moveMap[operation.SourcePath] = operation.DestinationPath
				if attachmentGroup < 0 {
					attachmentGroup = groupIndex
				}
			}
		}
	}
	if len(moveMap) == 0 {
		return
	}
	for _, documentPath := range trackedPaths {
		absolutePath := filepath.Join(plan.RepositoryRoot, filepath.FromSlash(documentPath))
		contents, readError := os.ReadFile(absolutePath)
		if readError != nil {
			continue
		}
		updated, replacements := rewriteMarkdownTargets(contents, documentPath, moveMap)
		if replacements > 0 {
			plan.Groups[attachmentGroup].Operations = append(plan.Groups[attachmentGroup].Operations, CleanupOperation{Kind: OperationReplace, SourcePath: documentPath, Contents: updated})
		}
	}
}

func rewriteMarkdownTargets(contents []byte, documentPath string, moveMap map[string]string) ([]byte, int) {
	updated := append([]byte(nil), contents...)
	replacements := 0
	searchOffset := 0
	for {
		openingRelative := bytes.Index(updated[searchOffset:], []byte("]("))
		if openingRelative < 0 {
			break
		}
		targetStart := searchOffset + openingRelative + 2
		closingRelative := bytes.IndexByte(updated[targetStart:], ')')
		if closingRelative < 0 {
			break
		}
		targetEnd := targetStart + closingRelative
		wholeTarget := string(updated[targetStart:targetEnd])
		pathToken := wholeTarget
		if whitespace := strings.IndexAny(pathToken, " \t"); whitespace >= 0 {
			pathToken = pathToken[:whitespace]
		}
		fragment := ""
		if hash := strings.Index(pathToken, "#"); hash >= 0 {
			fragment = pathToken[hash:]
			pathToken = pathToken[:hash]
		}
		if pathToken != "" && strings.Contains(pathToken, "/") {
			resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(filepath.FromSlash(documentPath)), filepath.FromSlash(pathToken))))
			if destination, moved := moveMap[resolved]; moved {
				relativeDestination, err := filepath.Rel(filepath.Dir(filepath.FromSlash(documentPath)), filepath.FromSlash(destination))
				if err == nil {
					replacement := filepath.ToSlash(relativeDestination) + fragment
					updated = append(updated[:targetStart], append([]byte(replacement), updated[targetStart+len(pathToken)+len(fragment):]...)...)
					targetEnd += len(replacement) - len(pathToken) - len(fragment)
					replacements++
				}
			}
		}
		searchOffset = targetEnd + 1
	}
	return updated, replacements
}

func AddBlankedRepairs(ctx context.Context, plan *CleanupPlan, restoreTargets []string) {
	approved := map[string]bool{}
	for _, target := range restoreTargets {
		cleaned := filepath.ToSlash(filepath.Clean(target))
		approved[cleaned] = true
	}
	for _, scanRoot := range []string{"do-work/archive", "do-work/queue", "do-work/working"} {
		absoluteScanRoot := filepath.Join(plan.RepositoryRoot, filepath.FromSlash(scanRoot))
		_ = filepath.WalkDir(absoluteScanRoot, func(path string, entry os.DirEntry, entryError error) error {
			if entryError != nil {
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			lowerName := strings.ToLower(entry.Name())
			upperName := strings.ToUpper(entry.Name())
			if entry.IsDir() || !(strings.HasPrefix(upperName, "REQ-") || strings.HasPrefix(upperName, "UR-")) || !strings.HasSuffix(lowerName, ".md") {
				return nil
			}
			contents, readError := os.ReadFile(path)
			if readError != nil || !isBlankedDocument(contents) {
				return nil
			}
			relativePath := repoRelative(plan.RepositoryRoot, path)
			recovered, sourceCommit, recordedCommit, recoveryError := recoverGitContent(ctx, plan.RepositoryRoot, relativePath)
			if recoveryError != nil {
				plan.Findings = append(plan.Findings, manualFinding("BLANKED-RECORD-UNRECOVERABLE", nil, []string{relativePath}, recoveryError.Error()))
				return nil
			}
			if !approved[relativePath] {
				plan.Findings = append(plan.Findings, resultmodel.CommandFinding{Code: "BLANKED-RECORD-REQUIRES-CONSENT", Severity: resultmodel.SeverityWarning,
					AffectedPaths: []string{relativePath}, Evidence: []string{fmt.Sprintf("recoverable from %s (%d bytes); unreferenced objects can be garbage-collected", sourceCommit, len(recovered))},
					Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "restoring overwrites the current file and requires its exact --restore-blanked target",
					NextArgv: []string{"do-work-cli", "cleanup", "--restore-blanked", relativePath}, VerificationArgv: []string{"do-work-cli", "cleanup", "--dry-run"}})
				return nil
			}
			if recoveredDocument, parseError := requestmodel.ParseDocument(recovered); parseError == nil && recordedCommit != "" {
				if len(recordedCommit) > 7 {
					recordedCommit = recordedCommit[:7]
				}
				if setError := recoveredDocument.SetScalar("commit", recordedCommit); setError != nil {
					plan.Findings = append(plan.Findings, manualFinding("BLANKED-COMMIT-RECORD-FAILED", nil, []string{relativePath}, setError.Error()))
					return nil
				}
				recovered = recoveredDocument.DocumentBytes()
			}
			plan.Groups = append(plan.Groups, OperationGroup{Code: "RESTORE-" + strings.ReplaceAll(relativePath, "/", "-"), PassNumber: 6,
				Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: relativePath, Contents: recovered}}})
			delete(approved, relativePath)
			return nil
		})
	}
	for target := range approved {
		plan.Findings = append(plan.Findings, manualFinding("RESTORE-TARGET-NOT-BLANKED", nil, []string{target}, "the exact restore target is not a blanked archived Markdown record"))
	}
	sortPlan(plan)
}

func isBlankedDocument(contents []byte) bool {
	if len(bytes.TrimSpace(contents)) == 0 {
		return true
	}
	_, err := requestmodel.ParseDocument(contents)
	return err != nil
}

var recordedHashPattern = regexp.MustCompile(`record commit hash ([0-9a-f]{7,40})`)

// RecoveryEvidence separates the blob supplying recoverable bytes from the
// implementation provenance recorded by the later metadata commit.
type RecoveryEvidence struct {
	ContentBytes         []byte
	SourceCommit         string
	ImplementationCommit string
}

// RecoverGitContent resolves the newest parseable bytes from full path history.
// It is read-only; cleanup remains the only caller authorized to publish them.
func RecoverGitContent(ctx context.Context, repositoryRoot, relativePath string) (RecoveryEvidence, error) {
	logOutput, err := cleanupGit(ctx, repositoryRoot, "log", "--full-history", "--format=%H", "--", relativePath)
	if err != nil {
		return RecoveryEvidence{}, err
	}
	recordedHash := ""
	for _, commitSHA := range strings.Fields(logOutput) {
		contents, showError := cleanupGitBytes(ctx, repositoryRoot, "show", commitSHA+":"+relativePath)
		if showError == nil && !isBlankedDocument(contents) {
			return RecoveryEvidence{ContentBytes: contents, SourceCommit: commitSHA, ImplementationCommit: recordedHash}, nil
		}
		if recordedHash == "" {
			subject, subjectError := cleanupGit(ctx, repositoryRoot, "log", "-1", "--format=%s", commitSHA)
			if subjectError == nil {
				if match := recordedHashPattern.FindStringSubmatch(strings.TrimSpace(subject)); match != nil {
					recordedHash = match[1]
				}
			}
		}
	}
	return RecoveryEvidence{}, fmt.Errorf("Git history has no recoverable nonblank content")
}

func recoverGitContent(ctx context.Context, repositoryRoot, relativePath string) ([]byte, string, string, error) {
	evidence, err := RecoverGitContent(ctx, repositoryRoot, relativePath)
	return evidence.ContentBytes, evidence.SourceCommit, evidence.ImplementationCommit, err
}

type WorktreeRepairOptions struct {
	DryRun       bool
	DiscardNames []string
}

func ApplyWorktreeRepairs(ctx context.Context, repositoryRoot string, options WorktreeRepairOptions) ([]resultmodel.RecordedChange, []resultmodel.CommandFinding) {
	changes := []resultmodel.RecordedChange{}
	findings := []resultmodel.CommandFinding{}
	discard := map[string]bool{}
	for _, name := range options.DiscardNames {
		discard[name] = true
	}
	if !options.DryRun {
		if _, pruneError := cleanupGit(ctx, repositoryRoot, "worktree", "prune"); pruneError != nil {
			return changes, []resultmodel.CommandFinding{manualFinding("WORKTREE-ENUMERATION-FAILED", nil, nil, pruneError.Error())}
		}
	}
	listOutput, err := cleanupGitBytes(ctx, repositoryRoot, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return changes, []resultmodel.CommandFinding{manualFinding("WORKTREE-ENUMERATION-FAILED", nil, nil, err.Error())}
	}
	trees := parseWorktrees(listOutput)
	branchOutput, branchError := cleanupGit(ctx, repositoryRoot, "branch", "--format=%(refname:short)", "--list", "worktree-agent-*")
	if branchError != nil {
		return changes, []resultmodel.CommandFinding{manualFinding("WORKTREE-ENUMERATION-FAILED", nil, nil, branchError.Error())}
	}
	branches := map[string]bool{}
	candidates := map[string]bool{}
	for _, branch := range strings.Fields(branchOutput) {
		branches[branch] = true
		candidates[branch] = true
	}
	for _, tree := range trees {
		candidates[tree.Name] = true
	}
	orderedCandidates := make([]string, 0, len(candidates))
	for name := range candidates {
		orderedCandidates = append(orderedCandidates, name)
	}
	sort.Strings(orderedCandidates)
	for _, name := range orderedCandidates {
		discardApproved := discard[name]
		delete(discard, name)
		tree := worktreeByName(trees, name)
		path := tree.Path
		clean := path == "" || worktreeClean(ctx, path)
		ancestryTarget := name
		if !branches[name] {
			ancestryTarget = tree.Head
		}
		merged := ancestryTarget != "" && gitExitSuccess(ctx, repositoryRoot, "merge-base", "--is-ancestor", ancestryTarget, "HEAD")
		if (!clean || !merged) && !discardApproved {
			findings = append(findings, resultmodel.CommandFinding{Code: "WORKTREE-REQUIRES-CONSENT", Severity: resultmodel.SeverityWarning,
				AffectedIDs: []string{requestIDFromWorktree(name)}, AffectedPaths: nonEmpty(path), Evidence: []string{fmt.Sprintf("%s is dirty=%t merged=%t", name, !clean, merged)},
				Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "unmerged or dirty builder work can be the only copy",
				NextArgv: []string{"do-work-cli", "cleanup", "--discard-worktree", name}, VerificationArgv: []string{"git", "worktree", "list", "--porcelain"}})
			continue
		}
		if options.DryRun {
			changes = append(changes, resultmodel.RecordedChange{Path: path, Kind: "worktree", Detail: "planned removal of " + name})
			continue
		}
		if path != "" {
			args := []string{"worktree", "remove", path}
			if discardApproved {
				args = append(args, "--force")
			}
			if _, removeError := cleanupGit(ctx, repositoryRoot, args...); removeError != nil {
				findings = append(findings, manualFinding("WORKTREE-REMOVE-FAILED", nil, nonEmpty(path), removeError.Error()))
				continue
			}
		}
		deleteFlag := "-d"
		if discardApproved {
			deleteFlag = "-D"
		}
		if branches[name] {
			if _, deleteError := cleanupGit(ctx, repositoryRoot, "branch", deleteFlag, name); deleteError != nil {
				findings = append(findings, manualFinding("WORKTREE-BRANCH-REMOVE-FAILED", nil, nonEmpty(path), deleteError.Error()))
				continue
			}
		}
		changes = append(changes, resultmodel.RecordedChange{Path: path, Kind: "worktree", Detail: "removed " + name})
	}
	for name := range discard {
		findings = append(findings, manualFinding("WORKTREE-DISCARD-NOT-FOUND", nil, nil, "requested worktree discard target "+name+" was not enumerated"))
	}
	return changes, findings
}

type worktreeRecord struct{ Name, Path, Head string }

func parseWorktrees(output []byte) []worktreeRecord {
	result := []worktreeRecord{}
	record := worktreeRecord{}
	flush := func() {
		if record.Path == "" {
			return
		}
		if record.Name == "" && strings.HasPrefix(filepath.Base(record.Path), "worktree-agent-") {
			record.Name = filepath.Base(record.Path)
		}
		if strings.HasPrefix(record.Name, "worktree-agent-") {
			result = append(result, record)
		}
		record = worktreeRecord{}
	}
	for _, field := range bytes.Split(output, []byte{0}) {
		if len(field) == 0 {
			flush()
			continue
		}
		value := string(field)
		switch {
		case strings.HasPrefix(value, "worktree "):
			record.Path = strings.TrimPrefix(value, "worktree ")
		case strings.HasPrefix(value, "HEAD "):
			record.Head = strings.TrimPrefix(value, "HEAD ")
		case strings.HasPrefix(value, "branch refs/heads/"):
			record.Name = strings.TrimPrefix(value, "branch refs/heads/")
		}
	}
	flush()
	return result
}

func worktreeByName(records []worktreeRecord, name string) worktreeRecord {
	for _, record := range records {
		if record.Name == name {
			return record
		}
	}
	return worktreeRecord{Name: name}
}

func worktreeClean(ctx context.Context, path string) bool {
	command := exec.CommandContext(ctx, "git", "-C", path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	output, err := command.Output()
	return err == nil && len(output) == 0
}

func requestIDFromWorktree(name string) string {
	marker := "REQ-"
	markerIndex := strings.Index(name, marker)
	if markerIndex < 0 {
		return ""
	}
	endIndex := markerIndex + len(marker)
	for endIndex < len(name) && name[endIndex] >= '0' && name[endIndex] <= '9' {
		endIndex++
	}
	if endIndex > markerIndex+len(marker) {
		return name[markerIndex:endIndex]
	}
	return ""
}

func cleanupGit(ctx context.Context, repositoryRoot string, args ...string) (string, error) {
	output, err := cleanupGitBytes(ctx, repositoryRoot, args...)
	return string(output), err
}

func cleanupGitBytes(ctx context.Context, repositoryRoot string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", repositoryRoot, "--literal-pathspecs"}, args...)...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

func gitExitSuccess(ctx context.Context, repositoryRoot string, args ...string) bool {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", repositoryRoot}, args...)...)
	return command.Run() == nil
}

func splitZero(value string) []string {
	parts := strings.Split(value, "\x00")
	result := []string{}
	for _, part := range parts {
		if part != "" {
			result = append(result, filepath.ToSlash(part))
		}
	}
	sort.Strings(result)
	return result
}

func nonEmpty(value string) []string {
	if value == "" {
		return []string{}
	}
	return []string{value}
}
