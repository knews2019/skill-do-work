package cleanup

// `do-work-cli worktree new|status|merge|cleanup` runs the builder worktree lifecycle that
// actions/fan-out-reference.md specifies by hand: Naming and its collision rule, Where
// worktrees live, the hand-back merge steps 1 to 4, and Cleanup — happy path. Every git
// call here is non-forcing by construction: a refusal from git is reported and the command
// stops. Every precondition is checked before the first side effect of a subcommand.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// worktreeLinksConfig is the optional per-repository list of paths, one repository-relative
// path per line, that `worktree new` symlinks from the main tree into each builder worktree.
const worktreeLinksConfig = "do-work/worktree-links"

var lifecycleRequestPattern = regexp.MustCompile(`^REQ-[0-9]+$`)

// lifecycleFilenamePrefix is the REQ-NNN- id prefix the naming suffix drops from a filename.
var lifecycleFilenamePrefix = regexp.MustCompile(`^REQ-[0-9]+-?`)

// lifecycleRun carries one subcommand's observations into its typed result and findings.
type lifecycleRun struct {
	ctx      context.Context
	mainRoot string
	report   *resultmodel.WorktreeLifecycleResult
	findings []resultmodel.CommandFinding
}

func handleWorktree(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	run := &lifecycleRun{ctx: context.Background(), report: &resultmodel.WorktreeLifecycleResult{}}
	if len(arguments) == 0 {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-USAGE", nil, nil, "worktree requires a subcommand: new, status, merge or cleanup")
	}
	run.report.Subcommand = arguments[0]
	requestID, optionValue, usageError := parseLifecycleArguments(arguments[0], arguments[1:])
	if usageError != nil {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-USAGE", nil, nil, usageError.Error())
	}
	run.report.RequestID = requestID
	listing, listError := cleanupGitBytes(run.ctx, executionContext.RepositoryRoot, "worktree", "list", "--porcelain", "-z")
	if listError != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, listError.Error())
	}
	// The first porcelain record is the main working tree, so a run from a linked tree
	// still derives names, links and the integration branch from the main tree.
	run.mainRoot = strings.TrimPrefix(string(bytes.SplitN(listing, []byte{0}, 2)[0]), "worktree ")
	builderWorktrees := parseWorktrees(listing)
	integrationBranch, branchFinding := run.integrationBranch(arguments[0] == "new" && optionValue != "")
	if branchFinding != nil {
		return *branchFinding
	}
	switch arguments[0] {
	case "new":
		if optionValue != "" {
			integrationBranch = optionValue
		}
		run.report.IntegrationBranch = integrationBranch
		return run.createWorktree(requestID, builderWorktrees)
	case "status":
		run.report.IntegrationBranch = integrationBranch
		return run.reportStatus(builderWorktrees)
	case "merge":
		run.report.IntegrationBranch = integrationBranch
		return run.mergeBuilderBranch(requestID, optionValue)
	default:
		run.report.IntegrationBranch = integrationBranch
		return run.removeBuilderWorktree(requestID, optionValue, builderWorktrees)
	}
}

// parseLifecycleArguments accepts `new REQ-N [--from <branch>]`, `status`, and
// `merge|cleanup REQ-N [--name <operative_name>]`.
func parseLifecycleArguments(subcommand string, arguments []string) (string, string, error) {
	optionName := map[string]string{"new": "--from", "merge": "--name", "cleanup": "--name", "status": ""}
	allowedOption, known := optionName[subcommand]
	if !known {
		return "", "", fmt.Errorf("unknown worktree subcommand %q: use new, status, merge or cleanup", subcommand)
	}
	if subcommand == "status" {
		if len(arguments) > 0 {
			return "", "", fmt.Errorf("worktree status takes no arguments")
		}
		return "", "", nil
	}
	requestID, optionValue := "", ""
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == allowedOption && index+1 < len(arguments) && arguments[index+1] != "":
			index++
			optionValue = arguments[index]
		case requestID == "" && lifecycleRequestPattern.MatchString(arguments[index]):
			requestID = arguments[index]
		default:
			return "", "", fmt.Errorf("worktree %s: unexpected argument %q; usage: worktree %s REQ-N [%s <value>]", subcommand, arguments[index], subcommand, allowedOption)
		}
	}
	if requestID == "" {
		return "", "", fmt.Errorf("worktree %s requires a request id such as REQ-42", subcommand)
	}
	return requestID, optionValue, nil
}

// integrationBranch is the branch checked out in the main tree. A detached main tree
// refuses unless `new` names its base explicitly; merge and cleanup also refuse a builder
// branch checked out there, because `branch -d` from it answers merged-ness wrongly.
func (run *lifecycleRun) integrationBranch(baseNamed bool) (string, *resultmodel.CommandResult) {
	output, err := cleanupGit(run.ctx, run.mainRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
	branch := strings.TrimSpace(output)
	if baseNamed {
		return branch, nil
	}
	if err != nil || branch == "" {
		result := run.stop(resultmodel.OutcomeRefused, "WORKTREE-DETACHED-HEAD", []string{run.mainRoot}, []string{"git", "-C", run.mainRoot, "status", "--short", "--branch"},
			"the main tree has no checked-out branch to integrate into")
		return "", &result
	}
	if strings.HasPrefix(branch, "worktree-agent-") && (run.report.Subcommand == "merge" || run.report.Subcommand == "cleanup") {
		result := run.stop(resultmodel.OutcomeRefused, "WORKTREE-BUILDER-BRANCH-CHECKED-OUT", []string{run.mainRoot}, []string{"git", "-C", run.mainRoot, "status", "--short", "--branch"},
			"the main tree has builder branch "+branch+" checked out; check out the integration branch first")
		return "", &result
	}
	return branch, nil
}

func (run *lifecycleRun) createWorktree(requestID string, builderWorktrees []worktreeRecord) resultmodel.CommandResult {
	snapshot, discoveryError := repositorymodel.DiscoverRepository(run.mainRoot)
	if discoveryError != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-DISCOVERY-FAILED", nil, nil, discoveryError.Error())
	}
	requestFiles := snapshot.RequestsByID[requestID]
	if len(requestFiles) != 1 {
		paths := []string{}
		for _, requestFile := range requestFiles {
			paths = append(paths, requestFile.RelativePath)
		}
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-REQUEST-NOT-UNIQUE", paths, nil,
			fmt.Sprintf("%s resolves to %d request files; worktree new needs exactly one", requestID, len(requestFiles)))
	}
	if _, err := cleanupGit(run.ctx, run.mainRoot, "rev-parse", "--verify", "--quiet", "refs/heads/"+run.report.IntegrationBranch); err != nil {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-BASE-NOT-A-BRANCH", nil, []string{"git", "-C", run.mainRoot, "branch", "--list"},
			"base "+run.report.IntegrationBranch+" is not a local branch")
	}
	linkPaths, linkFinding := run.readLinkConfig()
	if linkFinding != nil {
		return *linkFinding
	}
	branchNames, branchError := run.builderBranches("")
	if branchError != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, branchError.Error())
	}

	// Naming: the filename slug after REQ-NNN-, reduced to [a-z0-9-] as a text operation.
	slug := lifecycleFilenamePrefix.ReplaceAllString(strings.TrimSuffix(filepath.Base(requestFiles[0].RelativePath), ".md"), "")
	suffix := strings.Trim(strings.Map(func(character rune) rune {
		switch {
		case character >= 'A' && character <= 'Z':
			return character + ('a' - 'A')
		case character >= 'a' && character <= 'z', character >= '0' && character <= '9', character == '-':
			return character
		}
		return -1
	}, slug), "-")
	derivedName := "worktree-agent-" + requestID
	if suffix != "" {
		derivedName += "-" + suffix
	}
	worktreesParent := filepath.Join(filepath.Dir(run.mainRoot), filepath.Base(run.mainRoot)+"-worktrees")
	operativeName := derivedName
	for attempt := 2; lifecycleNameTaken(operativeName, filepath.Join(worktreesParent, operativeName), branchNames, builderWorktrees); attempt++ {
		operativeName = derivedName + "-" + strconv.Itoa(attempt)
	}
	worktreePath := filepath.Join(worktreesParent, operativeName)
	if operativeName != derivedName {
		run.finding(resultmodel.SeverityWarning, "WORKTREE-NAME-COLLISION", []string{derivedName}, nil,
			derivedName+" already exists as a branch, worktree or path; created "+operativeName+" beside it and left the original to its owners")
	}
	if _, err := cleanupGit(run.ctx, run.mainRoot, "worktree", "add", "-b", operativeName, worktreePath, run.report.IntegrationBranch); err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", []string{worktreePath}, nil, err.Error())
	}
	run.report.OperativeName, run.report.WorktreePath = operativeName, worktreePath

	for _, linkPath := range linkPaths {
		sourcePath := filepath.Join(run.mainRoot, filepath.FromSlash(linkPath))
		destinationPath := filepath.Join(worktreePath, filepath.FromSlash(linkPath))
		skipReason := ""
		if _, err := os.Lstat(sourcePath); err != nil {
			skipReason = "the main tree has no " + linkPath + " to link"
		} else if _, err := os.Lstat(destinationPath); err == nil {
			skipReason = "the worktree already has " + linkPath + "; it is never overwritten"
		}
		if skipReason != "" {
			run.report.LinksSkipped = append(run.report.LinksSkipped, linkPath)
			run.finding(resultmodel.SeverityWarning, "WORKTREE-LINK-SKIPPED", []string{linkPath}, nil, skipReason)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-LINK-FAILED", []string{linkPath}, nil, err.Error())
		}
		if err := os.Symlink(sourcePath, destinationPath); err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-LINK-FAILED", []string{linkPath}, nil, err.Error())
		}
		run.report.LinksCreated = append(run.report.LinksCreated, linkPath)
	}
	return run.finish(resultmodel.OutcomeSuccess)
}

func lifecycleNameTaken(name, path string, branchNames []string, builderWorktrees []worktreeRecord) bool {
	for _, branchName := range branchNames {
		if branchName == name {
			return true
		}
	}
	for _, record := range builderWorktrees {
		if record.Name == name || filepath.Base(record.Path) == name {
			return true
		}
	}
	_, err := os.Lstat(path)
	return err == nil
}

func (run *lifecycleRun) reportStatus(builderWorktrees []worktreeRecord) resultmodel.CommandResult {
	linkPaths, linkFinding := run.readLinkConfig()
	if linkFinding != nil {
		return *linkFinding
	}
	run.report.StatusRows = []resultmodel.WorktreeStatusRow{}
	for _, record := range builderWorktrees {
		// HEAD, not the name: a detached builder worktree has no branch to count from.
		counts, countError := cleanupGit(run.ctx, run.mainRoot, "rev-list", "--left-right", "--count", run.report.IntegrationBranch+"..."+record.Head)
		commitTime, timeError := cleanupGit(run.ctx, run.mainRoot, "log", "-1", "--format=%ct", record.Head)
		dirtyPaths, dirtyError := run.dirtyPaths(record.Path, linkPaths)
		if err := errors.Join(countError, timeError, dirtyError); err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", []string{record.Path}, nil, err.Error())
		}
		row := resultmodel.WorktreeStatusRow{OperativeName: record.Name, WorktreePath: record.Path, Dirty: len(dirtyPaths) > 0}
		if countFields := strings.Fields(counts); len(countFields) == 2 {
			row.Behind, _ = strconv.Atoi(countFields[0])
			row.Ahead, _ = strconv.Atoi(countFields[1])
		}
		row.LastCommitUnix, _ = strconv.ParseInt(strings.TrimSpace(commitTime), 10, 64)
		run.report.StatusRows = append(run.report.StatusRows, row)
	}
	return run.finish(resultmodel.OutcomeSuccess)
}

// mergeBuilderBranch is fan-out-reference.md "When to merge" steps 1 to 4. Step 0 (staging
// run artifacts) and integration seams stay with the integrator's hand path.
func (run *lifecycleRun) mergeBuilderBranch(requestID, requestedName string) resultmodel.CommandResult {
	operativeName, nameFinding := run.resolveOperativeName(requestID, requestedName)
	if nameFinding != nil {
		return *nameFinding
	}
	if gitExitSuccess(run.ctx, run.mainRoot, "rev-parse", "--quiet", "--verify", "MERGE_HEAD") {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-MERGE-IN-PROGRESS", nil, []string{"git", "-C", run.mainRoot, "status", "--short"},
			"a merge is already in progress in the main tree")
	}
	stagedPaths, err := cleanupGit(run.ctx, run.mainRoot, "diff", "--cached", "--name-only")
	if err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
	}
	if staged := strings.Fields(stagedPaths); len(staged) > 0 {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-INDEX-NOT-EMPTY", staged, []string{"git", "-C", run.mainRoot, "diff", "--cached", "--name-only"},
			"the index holds staged paths, and the merge commit would take them")
	}
	preOutput, preError := cleanupGit(run.ctx, run.mainRoot, "rev-parse", "--short", "HEAD")
	aheadOutput, aheadError := cleanupGit(run.ctx, run.mainRoot, "rev-list", "--count", "HEAD.."+operativeName)
	if err := errors.Join(preError, aheadError); err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
	}
	run.report.Pre = strings.TrimSpace(preOutput)
	// `git merge` prints "Already up to date." and exits 0 here, so count instead.
	if strings.TrimSpace(aheadOutput) == "0" {
		run.report.EmptyHandBack = true
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-EMPTY-HAND-BACK", nil, []string{"git", "-C", run.mainRoot, "log", "--oneline", "-1", operativeName},
			operativeName+" has no commits beyond the integration tip; nothing was merged or committed")
	}
	queuePaths, err := cleanupGit(run.ctx, run.mainRoot, "diff", "--name-only", run.report.Pre+"..."+operativeName, "--", "do-work/")
	if err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
	}
	if queueWrites := strings.Fields(queuePaths); len(queueWrites) > 0 {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-QUEUE-GUARD", queueWrites, []string{"git", "-C", run.mainRoot, "diff", "--name-only", run.report.Pre + "..." + operativeName, "--", "do-work/"},
			operativeName+" commits queue state under do-work/; drop or revert those commits on the branch before integrating")
	}
	if _, mergeError := cleanupGit(run.ctx, run.mainRoot, "merge", "--no-ff", "--no-commit", operativeName); mergeError != nil {
		conflicts, _ := cleanupGit(run.ctx, run.mainRoot, "diff", "--name-only", "--diff-filter=U")
		run.report.ConflictedPaths = strings.Fields(conflicts)
		if len(run.report.ConflictedPaths) == 0 {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, mergeError.Error())
		}
		return run.stop(resultmodel.OutcomeFindings, "WORKTREE-MERGE-CONFLICT", run.report.ConflictedPaths, []string{"git", "-C", run.mainRoot, "status", "--short"},
			"the merge stopped on conflicts and is left in progress; resolve and commit by hand with the [REQ-N] merge subject")
	}
	if _, err := cleanupGit(run.ctx, run.mainRoot, "commit", "-m", "["+requestID+"] merge builder branch "+operativeName); err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
	}
	mergeHash, err := cleanupGit(run.ctx, run.mainRoot, "rev-parse", "--short", "HEAD")
	if err != nil {
		return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
	}
	run.report.MergeHash = strings.TrimSpace(mergeHash)
	return run.finish(resultmodel.OutcomeSuccess)
}

// removeBuilderWorktree is fan-out-reference.md "Cleanup — happy path", preceded by a
// preflight so a refusal leaves the link, the worktree and the branch exactly as they were.
func (run *lifecycleRun) removeBuilderWorktree(requestID, requestedName string, builderWorktrees []worktreeRecord) resultmodel.CommandResult {
	operativeName, nameFinding := run.resolveOperativeName(requestID, requestedName)
	if nameFinding != nil {
		return *nameFinding
	}
	if !gitExitSuccess(run.ctx, run.mainRoot, "merge-base", "--is-ancestor", operativeName, "HEAD") {
		return run.stop(resultmodel.OutcomeRefused, "WORKTREE-NOT-MERGED", nil, []string{"git", "-C", run.mainRoot, "log", "--oneline", "HEAD.." + operativeName},
			operativeName+" is not merged into "+run.report.IntegrationBranch)
	}
	linkPaths, linkFinding := run.readLinkConfig()
	if linkFinding != nil {
		return *linkFinding
	}
	worktreePath := ""
	for _, record := range builderWorktrees {
		if record.Name == operativeName {
			worktreePath = record.Path
		}
	}
	run.report.WorktreePath = worktreePath
	ownedLinks := []string{}
	if worktreePath != "" {
		dirtyPaths, err := run.dirtyPaths(worktreePath, linkPaths)
		if err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", []string{worktreePath}, nil, err.Error())
		}
		if len(dirtyPaths) > 0 {
			return run.stop(resultmodel.OutcomeRefused, "WORKTREE-DIRTY", dirtyPaths, []string{"git", "-C", worktreePath, "status", "--short"},
				operativeName+" has uncommitted builder work; commit it on the branch or discard it through cleanup Pass 5")
		}
		ownedLinks = run.ownedLinks(worktreePath, linkPaths)
	}

	for _, linkPath := range ownedLinks {
		if err := os.Remove(filepath.Join(worktreePath, filepath.FromSlash(linkPath))); err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-LINK-FAILED", []string{linkPath}, nil, err.Error())
		}
	}
	steps := [][]string{{"branch", "-d", operativeName}, {"worktree", "prune"}}
	if worktreePath != "" {
		steps = append([][]string{{"worktree", "remove", worktreePath}}, steps...)
	}
	for _, step := range steps {
		if _, err := cleanupGit(run.ctx, run.mainRoot, step...); err != nil {
			return run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", []string{operativeName}, nil, err.Error())
		}
	}
	return run.finish(resultmodel.OutcomeSuccess)
}

// resolveOperativeName enumerates worktree-agent-REQ-N-* branches; requestIDFromWorktree
// keeps REQ-66 from matching REQ-660. Zero or several matches refuse unless --name picks one.
func (run *lifecycleRun) resolveOperativeName(requestID, requestedName string) (string, *resultmodel.CommandResult) {
	candidates, err := run.builderBranches(requestID)
	if err != nil {
		result := run.stop(resultmodel.OutcomeFailure, "WORKTREE-GIT-FAILED", nil, nil, err.Error())
		return "", &result
	}
	for _, candidate := range candidates {
		if candidate == requestedName || (requestedName == "" && len(candidates) == 1) {
			run.report.OperativeName = candidate
			return candidate, nil
		}
	}
	evidence := fmt.Sprintf("%d worktree-agent-%s-* branches match; pass --name with the exact operative name", len(candidates), requestID)
	if requestedName != "" {
		evidence = fmt.Sprintf("--name %s is not one of the %d worktree-agent-%s-* branches", requestedName, len(candidates), requestID)
	}
	result := run.stop(resultmodel.OutcomeRefused, "WORKTREE-OPERATIVE-NAME-UNRESOLVED", candidates, []string{"git", "-C", run.mainRoot, "branch", "--list", "worktree-agent-" + requestID + "*"}, evidence)
	return "", &result
}

// builderBranches lists local worktree-agent-* branch names, only those for requestID when
// it is set.
func (run *lifecycleRun) builderBranches(requestID string) ([]string, error) {
	output, err := cleanupGit(run.ctx, run.mainRoot, "for-each-ref", "--format=%(refname)", "refs/heads/worktree-agent-*")
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, reference := range strings.Fields(output) {
		name := strings.TrimPrefix(reference, "refs/heads/")
		if requestID == "" || requestIDFromWorktree(name) == requestID {
			names = append(names, name)
		}
	}
	return names, nil
}

// readLinkConfig reads worktreeLinksConfig from the main tree. An absent file means no links;
// an absolute path or a .. segment refuses the whole command before anything is created.
func (run *lifecycleRun) readLinkConfig() ([]string, *resultmodel.CommandResult) {
	contents, err := os.ReadFile(filepath.Join(run.mainRoot, filepath.FromSlash(worktreeLinksConfig)))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		result := run.stop(resultmodel.OutcomeFailure, "WORKTREE-LINK-CONFIG-UNREADABLE", []string{worktreeLinksConfig}, nil, err.Error())
		return nil, &result
	}
	linkPaths := []string{}
	for _, line := range strings.Split(string(contents), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cleaned := filepath.ToSlash(filepath.Clean(line))
		if filepath.IsAbs(line) || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.Contains("/"+line+"/", "/../") {
			result := run.stop(resultmodel.OutcomeRefused, "WORKTREE-LINK-CONFIG-INVALID", []string{worktreeLinksConfig}, nil,
				fmt.Sprintf("%s line %q must be a repository-relative path without .. segments", worktreeLinksConfig, line))
			return nil, &result
		}
		linkPaths = append(linkPaths, cleaned)
	}
	return linkPaths, nil
}

// ownedLinks are the configured paths that are symlinks to their own main-tree path, the
// only links this command created and the only ones cleanup may remove.
func (run *lifecycleRun) ownedLinks(worktreePath string, linkPaths []string) []string {
	owned := []string{}
	for _, linkPath := range linkPaths {
		target, err := os.Readlink(filepath.Join(worktreePath, filepath.FromSlash(linkPath)))
		if err == nil && target == filepath.Join(run.mainRoot, filepath.FromSlash(linkPath)) {
			owned = append(owned, linkPath)
		}
	}
	return owned
}

// dirtyPaths is the worktree's porcelain status minus the owned links: a symlink is not
// matched by a `node_modules/` ignore line, so without this filter every linked worktree
// would read as dirty.
func (run *lifecycleRun) dirtyPaths(worktreePath string, linkPaths []string) ([]string, error) {
	output, err := cleanupGit(run.ctx, worktreePath, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	owned := map[string]bool{}
	for _, linkPath := range run.ownedLinks(worktreePath, linkPaths) {
		owned[linkPath] = true
	}
	dirty := []string{}
	entries := splitZero(output)
	for index := 0; index < len(entries); index++ {
		if len(entries[index]) < 4 {
			continue
		}
		statusCode, entryPath := entries[index][:2], entries[index][3:]
		if statusCode[0] == 'R' || statusCode[0] == 'C' {
			index++ // the rename or copy source follows as its own NUL field
		}
		if !owned[entryPath] {
			dirty = append(dirty, entryPath)
		}
	}
	return dirty, nil
}

// finding records one observation; nextArgv is always a read-only command, never the
// destructive step it would take to resolve the finding.
func (run *lifecycleRun) finding(severity resultmodel.FindingSeverity, code string, affectedPaths, nextArgv []string, evidence string) {
	fixability := resultmodel.FixabilityManual
	stopReason := ""
	if severity == resultmodel.SeverityError {
		fixability = resultmodel.FixabilityRefused
		stopReason = "worktree " + run.report.Subcommand + " stopped before any further step"
	}
	affectedIDs := []string{}
	if run.report.RequestID != "" {
		affectedIDs = []string{run.report.RequestID}
	}
	run.findings = append(run.findings, resultmodel.CommandFinding{Code: code, Severity: severity, AffectedIDs: affectedIDs, AffectedPaths: affectedPaths,
		Evidence: []string{evidence}, Fixability: fixability, AutomationStopReason: stopReason, NextArgv: nextArgv,
		VerificationArgv: []string{"do-work-cli", "--format", "json", "worktree", "status"}})
}

func (run *lifecycleRun) stop(outcome resultmodel.CommandOutcome, code string, affectedPaths, nextArgv []string, evidence string) resultmodel.CommandResult {
	run.finding(resultmodel.SeverityError, code, affectedPaths, nextArgv, evidence)
	return run.finish(outcome)
}

func (run *lifecycleRun) finish(outcome resultmodel.CommandOutcome) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: outcome, Findings: run.findings, Worktree: run.report}
}
