package toolboxcommands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const last30DaysSource = "https://github.com/mvanhorn/last30days-skill"

var last30DaysLookPath = exec.LookPath
var last30DaysPythonQualifies = func(path string) bool {
	return exec.Command(path, "-c", "import sys; raise SystemExit(0 if sys.version_info >= (3, 12) else 1)").Run() == nil
}
var last30DaysCopyTree = copyLast30DaysTree
var last30DaysRename = os.Rename
var last30DaysMkdir = os.Mkdir

func handleLast30Days(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(args)
	if err != nil {
		return usageResult(CommandLast30Days, err.Error())
	}
	if commit {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{toolboxFinding(CommandLast30Days, "LAST30DAYS-COMMIT-PRIVATE", resultmodel.SeverityWarning, nil, "last30days is deliberately project-local and Git-ignored", resultmodel.FixabilityRefused, "there is no committable target")}}
	}
	if len(rest) < 2 || len(rest) > 3 {
		return usageResult(CommandLast30Days, "Usage: install-last30days [--dry-run|--commit] (check|install) <project-root> [source-repository]")
	}
	mode, project := rest[0], absoluteFromRoot(ctx.RepositoryRoot, rest[1])
	if mode != "check" && mode != "install" {
		return usageResult(CommandLast30Days, "unknown last30days mode: "+mode)
	}
	source := last30DaysSource
	if len(rest) == 3 {
		source = rest[2]
	}
	target := filepath.Join(project, ".claude", "skills", "last30days")
	changes := []resultmodel.RecordedChange{}
	if mode == "install" && last30DaysComplete(target) && !dryRun {
		_, excludeErr := prepareLast30DaysExclude(project)
		if excludeErr != nil {
			return last30DaysFailure(excludeErr)
		}
	}
	if mode == "install" && !last30DaysComplete(target) {
		if eligibilityErr := preflightLast30DaysTarget(project, target); eligibilityErr != nil {
			return last30DaysFailure(eligibilityErr)
		}
		undoExclude := func() error { return nil }
		var excludeErr error
		if !dryRun {
			undoExclude, excludeErr = prepareLast30DaysExclude(project)
			if excludeErr != nil {
				return last30DaysFailure(excludeErr)
			}
		}
		keepExclude := false
		defer func() {
			if !keepExclude {
				_ = undoExclude()
			}
		}()
		invocationContext, stopSignals, interrupted := mutationSignalContext()
		defer stopSignals()
		cloneRoot, cloneErr := os.MkdirTemp("", "do-work-last30days-*")
		if cloneErr != nil {
			return last30DaysFailure(cloneErr)
		}
		defer os.RemoveAll(cloneRoot)
		command := exec.CommandContext(invocationContext, "git", "clone", "--depth", "1", source, cloneRoot)
		if output, runErr := command.CombinedOutput(); runErr != nil {
			result := last30DaysFailure(fmt.Errorf("clone/source validation FAILED: %v: %s", runErr, output))
			result.ExitCodeOverride = interrupted()
			return result
		}
		sourceTree := filepath.Join(cloneRoot, "skills", "last30days")
		if !last30DaysComplete(sourceTree) {
			return last30DaysFailure(fmt.Errorf("clone/source validation FAILED"))
		}
		if dryRun {
			changes = append(changes, resultmodel.RecordedChange{Path: filepath.ToSlash(target), Kind: "planned", Detail: "would install complete last30days payload"})
		} else if publishErr := publishLast30Days(invocationContext, sourceTree, target); publishErr != nil {
			result := last30DaysFailure(publishErr)
			result.ExitCodeOverride = interrupted()
			return result
		} else {
			changes = append(changes, resultmodel.RecordedChange{Path: filepath.ToSlash(target), Kind: "published", Detail: "installed complete last30days payload"})
			keepExclude = true
		}
	}
	lines := []string{}
	failed := false
	if last30DaysComplete(target) || dryRun && len(changes) > 0 {
		lines = append(lines, "runtime payload: OK")
	} else {
		lines = append(lines, "runtime payload: FAILED")
		failed = true
	}
	if gitRepository(project) {
		if ignored := exec.Command("git", "-C", project, "check-ignore", "-q", ".claude/skills/last30days/SKILL.md").Run() == nil; ignored || dryRun && mode == "install" {
			lines = append(lines, "ignore rule: OK")
		} else {
			lines = append(lines, "ignore rule: FAILED")
			failed = true
		}
	} else {
		lines = append(lines, "ignore rule: n/a (not a git repo)")
	}
	python := findTargetPython()
	if python != "" {
		lines = append(lines, "python 3.12+: OK ("+python+")")
	} else {
		lines = append(lines, "python 3.12+: FAILED")
		failed = true
	}
	output := strings.Join(lines, "\n") + "\n"
	outcome := resultmodel.OutcomeSuccess
	if failed {
		outcome = resultmodel.OutcomeFindings
	}
	findings := []resultmodel.CommandFinding{}
	if failed {
		findings = append(findings, toolboxFinding(CommandLast30Days, "LAST30DAYS-CHECK-FAILED", resultmodel.SeverityWarning, []string{filepath.ToSlash(target)}, output, resultmodel.FixabilityManual, "the installed target is not runnable"))
	}
	return resultmodel.CommandResult{Outcome: outcome, Changes: changes, Findings: findings, ExactTextOutput: &output}
}

func last30DaysComplete(root string) bool {
	for _, relative := range []string{"SKILL.md", "scripts/last30days.py"} {
		info, err := os.Stat(filepath.Join(root, relative))
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			return false
		}
	}
	return true
}
func last30DaysFailure(err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{toolboxFinding(CommandLast30Days, "LAST30DAYS-FAILED", resultmodel.SeverityError, nil, err.Error(), resultmodel.FixabilityManual, "installation or verification did not complete")}}
}

func publishLast30Days(ctx context.Context, source, target string) error {
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(parent, ".last30days.staging.*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	if err := last30DaysCopyTree(source, stage); err != nil {
		return fmt.Errorf("clone/copy FAILED: %w", err)
	}
	if !last30DaysComplete(stage) {
		return fmt.Errorf("clone/copy FAILED")
	}
	backup := ""
	keepBackup := false
	if _, err := os.Lstat(target); err == nil {
		backup, err = os.MkdirTemp(parent, ".last30days.backup.*")
		if err != nil {
			return err
		}
		defer func() {
			if !keepBackup {
				_ = os.RemoveAll(backup)
			}
		}()
		if err := last30DaysRename(target, filepath.Join(backup, "previous")); err != nil {
			return fmt.Errorf("existing-tree backup FAILED: %w", err)
		}
	}
	restoreBackup := func() error {
		if backup != "" {
			return last30DaysRename(filepath.Join(backup, "previous"), target)
		}
		return nil
	}
	if err := ctx.Err(); err != nil {
		if restoreErr := restoreBackup(); restoreErr != nil {
			keepBackup = true
			return fmt.Errorf("publication interrupted; rollback FAILED; prior tree remains at %s: %w", filepath.Join(backup, "previous"), restoreErr)
		}
		return err
	}
	if err := last30DaysMkdir(target, 0o700); err != nil {
		if _, collisionErr := os.Lstat(target); collisionErr == nil {
			if backup != "" {
				keepBackup = true
				return fmt.Errorf("publication FAILED; %s reappeared — prior tree remains at %s", target, filepath.Join(backup, "previous"))
			}
			return fmt.Errorf("publication FAILED; %s reappeared", target)
		}
		if restoreErr := restoreBackup(); restoreErr != nil {
			keepBackup = true
			return fmt.Errorf("publication claim FAILED: %v; rollback FAILED; prior tree remains at %s: %w", err, filepath.Join(backup, "previous"), restoreErr)
		}
		return fmt.Errorf("publication claim FAILED: %w", err)
	}
	publishedInfo, statErr := os.Lstat(target)
	if statErr != nil {
		return statErr
	}
	restore := func() error {
		currentInfo, currentErr := os.Lstat(target)
		if currentErr != nil || !os.SameFile(publishedInfo, currentInfo) {
			return fmt.Errorf("published target changed; replacement preserved")
		}
		if removeErr := os.RemoveAll(target); removeErr != nil {
			return removeErr
		}
		if backup != "" {
			return last30DaysRename(filepath.Join(backup, "previous"), target)
		}
		return nil
	}
	if err := last30DaysCopyTree(stage, target); err != nil {
		if restoreErr := restore(); restoreErr != nil {
			keepBackup = true
			return fmt.Errorf("publication copy FAILED: %v; rollback FAILED; prior tree remains at %s: %w", err, filepath.Join(backup, "previous"), restoreErr)
		}
		return fmt.Errorf("publication copy FAILED: %w", err)
	}
	if err := ctx.Err(); err != nil {
		if restoreErr := restore(); restoreErr != nil {
			keepBackup = true
			return fmt.Errorf("publication interrupted; rollback FAILED; prior tree remains at %s: %w", filepath.Join(backup, "previous"), restoreErr)
		}
		return err
	}
	if !last30DaysComplete(target) {
		if restoreErr := restore(); restoreErr != nil {
			keepBackup = true
			return fmt.Errorf("publication verification FAILED; rollback FAILED; prior tree remains at %s: %w", filepath.Join(backup, "previous"), restoreErr)
		}
		return fmt.Errorf("publication verification FAILED; previous tree restored")
	}
	return nil
}

func copyLast30DaysTree(source, target string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		destination := filepath.Join(target, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(destination, info.Mode().Perm())
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported source entry %s", relative)
		}
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
		if err != nil {
			input.Close()
			return err
		}
		_, copyErr := io.Copy(output, input)
		closeErr := output.Close()
		inputCloseErr := input.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		return inputCloseErr
	})
}
func gitRepository(root string) bool {
	return exec.Command("git", "-C", root, "rev-parse", "--git-dir").Run() == nil
}

func preflightLast30DaysTarget(project, target string) error {
	relative, err := filepath.Rel(project, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("last30days target escapes project root")
	}
	if err := validateNoLinkedAncestors(project, filepath.ToSlash(relative), false); err != nil {
		return err
	}
	if !gitRepository(project) {
		return nil
	}
	tracked, err := exec.Command("git", "-C", project, "ls-files", "-z", "--", filepath.ToSlash(relative)).Output()
	if err != nil {
		return err
	}
	if len(tracked) > 0 {
		return fmt.Errorf("tracked vendored path must be untracked before installation; existing bytes were preserved")
	}
	return nil
}

func prepareLast30DaysExclude(root string) (func() error, error) {
	noUndo := func() error { return nil }
	if !gitRepository(root) {
		return noUndo, nil
	}
	pathBytes, err := exec.Command("git", "-C", root, "rev-parse", "--git-path", "info/exclude").Output()
	if err != nil {
		return noUndo, err
	}
	exclude := strings.TrimSpace(string(pathBytes))
	if !filepath.IsAbs(exclude) {
		exclude = filepath.Join(root, exclude)
	}
	original, readErr := os.ReadFile(exclude)
	existed := readErr == nil
	if readErr != nil && !os.IsNotExist(readErr) {
		return noUndo, readErr
	}
	undo := func() error {
		if existed {
			return os.WriteFile(exclude, original, 0o600)
		}
		return os.Remove(exclude)
	}
	if exec.Command("git", "-C", root, "check-ignore", "-q", ".claude/skills/last30days/SKILL.md").Run() != nil {
		handle, openErr := os.OpenFile(exclude, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if openErr != nil {
			return noUndo, openErr
		}
		_, writeErr := handle.WriteString("\n**/.claude/skills/last30days/\n")
		closeErr := handle.Close()
		if writeErr != nil {
			_ = undo()
			return noUndo, writeErr
		}
		if closeErr != nil {
			_ = undo()
			return noUndo, closeErr
		}
	}
	publishedInfo, publishedStatErr := os.Lstat(exclude)
	publishedBytes, publishedReadErr := os.ReadFile(exclude)
	if publishedStatErr != nil || publishedReadErr != nil || !publishedInfo.Mode().IsRegular() {
		return noUndo, fmt.Errorf("local exclude publication identity unavailable")
	}
	undo = func() error {
		currentInfo, statErr := os.Lstat(exclude)
		currentBytes, currentReadErr := os.ReadFile(exclude)
		if statErr != nil || currentReadErr != nil || !os.SameFile(publishedInfo, currentInfo) || !bytes.Equal(publishedBytes, currentBytes) {
			return fmt.Errorf("local exclude changed after publication; preserved replacement")
		}
		if existed {
			return os.WriteFile(exclude, original, 0o600)
		}
		return os.Remove(exclude)
	}
	if exec.Command("git", "-C", root, "check-ignore", "-q", ".claude/skills/last30days/SKILL.md").Run() != nil {
		_ = undo()
		return noUndo, fmt.Errorf("local exclude verification FAILED")
	}
	return undo, nil
}

func ensureLast30DaysExclude(root string) error {
	if !gitRepository(root) {
		return nil
	}
	if exec.Command("git", "-C", root, "ls-files", "--error-unmatch", "--", ".claude/skills/last30days").Run() == nil {
		return fmt.Errorf("tracked vendored path must be untracked before local ignore can protect it")
	}
	if exec.Command("git", "-C", root, "check-ignore", "-q", ".claude/skills/last30days/SKILL.md").Run() == nil {
		return nil
	}
	pathBytes, err := exec.Command("git", "-C", root, "rev-parse", "--git-path", "info/exclude").Output()
	if err != nil {
		return err
	}
	exclude := strings.TrimSpace(string(pathBytes))
	if !filepath.IsAbs(exclude) {
		exclude = filepath.Join(root, exclude)
	}
	handle, err := os.OpenFile(exclude, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer handle.Close()
	_, err = handle.WriteString("\n**/.claude/skills/last30days/\n")
	return err
}
func findTargetPython() string {
	for _, candidate := range []string{"python3.13", "python3.12", "python3", "python"} {
		path, err := last30DaysLookPath(candidate)
		if err != nil {
			continue
		}
		if last30DaysPythonQualifies(path) {
			return candidate
		}
	}
	return ""
}
