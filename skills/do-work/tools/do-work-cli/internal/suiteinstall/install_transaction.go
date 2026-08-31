// Package suiteinstall carries the install and update transactions and the five commands
// that expose them.
//
// The transaction keeps its own snapshot-and-restore recovery rather than using
// gittransaction: it replaces four whole module trees (which gittransaction refuses — it
// declares exact files) and it must restore the Git index, which gittransaction cannot do.
// It reports its outcome in the same resultmodel vocabulary so callers still see one
// contract.
//
// It also keeps every subprocess the shell installer used — cp -Rp, tar, diff, git,
// just — because those preserve byte, mode and symlink semantics for free. Only the Python
// and jq branches moved into Go.
package suiteinstall

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/archivefetch"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/managedsection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/settingshooks"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/suitemanifest"
)

// SkipCodeInstallCancelled marks a declined confirmation. Cancellation is a success with
// skipped work rather than a refusal, which keeps today's exit 0 through the public shell
// path without a second outcome-to-number table.
const SkipCodeInstallCancelled = "INSTALL-CANCELLED"

// The agent instructions live behind their own marker pair, not the Just recipe pair.
const (
	instructionsBeginMarker = "<!-- >>> do-work:communication-style >>> -->"
	instructionsEndMarker   = "<!-- <<< do-work:communication-style <<< -->"
)

// reservedRecipeNames must all survive into the installed Justfile; a candidate missing one
// is a broken managed section, not a customisation.
var reservedRecipeNames = []string{"run-kanban", "run-kanban-cli", "kanban-static", "kanban-summary", "run-do-work-update"}

// justfileCandidateNames is the discovery order. The first existing entry wins, and its real
// directory-entry spelling is what recovery restores (REQ-180).
var justfileCandidateNames = []string{"justfile", "Justfile", ".justfile"}

// InstallOptions names everything an install needs from its caller. ExtractedSourceRoot lets
// update-suite hand over an archive it already fetched, extracted and validated, so the
// update path never downloads twice.
type InstallOptions struct {
	ProjectRoot         string
	SuppliedArchivePath string
	UpstreamURL         string
	ToolDirectory       string
	ExtractedSourceRoot string
	Narration           io.Writer
	ConfirmationInput   io.Reader
}

// InstallResult reports the transaction in resultmodel's vocabulary. FailureReason carries
// the diagnostic verbatim so the command layer can put it in a finding.
type InstallResult struct {
	Outcome       resultmodel.CommandOutcome
	SuiteVersion  string
	Changes       []resultmodel.RecordedChange
	SkippedWork   []resultmodel.SkippedWork
	Rollback      resultmodel.RollbackResult
	FailurePaths  []string
	FailureReason string
	Cancelled     bool
}

// installTransaction holds the state recovery needs. Its field names mirror the shell
// installer's variables so the two can be read side by side during review.
type installTransaction struct {
	options       InstallOptions
	projectRoot   string
	gitIndexPath  string
	installTmp    string
	sourceRoot    string
	backupRoot    string
	suiteVersion  string
	modules       []moduleInstallPlan
	justTarget    string
	justExisted   bool
	justCandidate string

	settingsTarget    string
	settingsExisted   bool
	settingsCandidate string

	instructionsTarget    string
	instructionsExisted   bool
	instructionsCandidate string

	gitIndexExisted bool
	writeStarted    bool
	installVerified bool
	recoveryFailed  bool
	recoveryRan     bool

	// recoveryFinished is closed once this transaction's cleanup has run. The signal handler
	// waits on it so an interrupted install exits only after recovery has completed.
	recoveryFinished chan struct{}
}

type moduleInstallPlan struct {
	sourcePath      string
	destinationPath string
	relativePath    string
	existed         bool
	backupPath      string
}

// installFailure is a diagnostic plus the paths it concerns.
type installFailure struct {
	reason string
	paths  []string
}

func (failure *installFailure) Error() string { return failure.reason }

func failInstall(format string, arguments ...any) *installFailure {
	return &installFailure{reason: fmt.Sprintf(format, arguments...)}
}

// RunInstall performs the whole transaction: every guard runs before the single
// confirmation, every write after it, and any failure past the first write recovers every
// managed path plus the Git index.
func RunInstall(ctx context.Context, options InstallOptions) InstallResult {
	// Work runs under a cancellable context so an arriving signal STOPS the in-flight
	// subprocess rather than racing it: recovery then runs once, on this goroutine, through
	// the same path a failed write takes.
	workContext, cancelWork := context.WithCancel(ctx)
	defer cancelWork()
	transaction := &installTransaction{options: options, recoveryFinished: make(chan struct{})}

	// Deferred calls run last-in-first-out, so these three run as: cleanup (which recovers),
	// then the signal handler's release, then the signal subscription teardown. Signals stay
	// armed for the whole of cleanup, which is when recovery actually happens.
	stopSignals := transaction.armSignalRecovery(cancelWork)
	defer stopSignals()
	defer close(transaction.recoveryFinished)
	defer transaction.cleanup()

	if err := transaction.prepare(workContext); err != nil {
		return transaction.failureResult(err)
	}
	confirmed, err := transaction.reviewAndConfirm(workContext)
	if err != nil {
		return transaction.failureResult(err)
	}
	if !confirmed {
		transaction.narrate("Installation cancelled; no files were changed.")
		return InstallResult{
			Outcome:      resultmodel.OutcomeSuccess,
			SuiteVersion: transaction.suiteVersion,
			Cancelled:    true,
			SkippedWork: []resultmodel.SkippedWork{{
				Code:   SkipCodeInstallCancelled,
				Reason: "the single install confirmation was declined; no files were changed",
			}},
		}
	}
	if err := transaction.writeAndVerify(workContext); err != nil {
		return transaction.failureResult(err)
	}
	transaction.installVerified = true
	transaction.narrate("Installed do-work suite v%s with four verified modules.", transaction.suiteVersion)
	return InstallResult{
		Outcome:      resultmodel.OutcomeSuccess,
		SuiteVersion: transaction.suiteVersion,
		Changes:      transaction.recordedChanges(),
		Rollback:     resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
	}
}

// interruptedInstallExitStatus is what an install killed by HUP, INT or TERM exits with. All
// three map to 130 here because that is what the shell installer did, and unifying them with
// the fetcher's 129/130/143 would change a status the update suite asserts elsewhere. This
// is a signal status, not an outcome: resultmodel.ExitCode stays the only outcome-to-number
// authority.
const interruptedInstallExitStatus = 130

// armSignalRecovery makes an interruption take the same recovery path a failed write does. A
// signal arriving mid-write must not leave a half-installed suite behind.
//
// The handler never recovers itself. It cancels the work context — which kills the in-flight
// cp, tar or git and turns it into an ordinary write failure — and then waits for the main
// goroutine to finish that recovery before exiting. Recovering from the handler instead would
// race the writes it is trying to undo, and the recovery would sometimes report itself
// incomplete because the main goroutine was still copying into a directory it had removed.
func (transaction *installTransaction) armSignalRecovery(cancelWork context.CancelFunc) func() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	released := make(chan struct{})
	go func() {
		select {
		case <-signalChannel:
			cancelWork()
			<-transaction.recoveryFinished
			os.Exit(interruptedInstallExitStatus)
		case <-released:
		}
	}()
	return func() {
		signal.Stop(signalChannel)
		close(released)
	}
}

func (transaction *installTransaction) failureResult(err error) InstallResult {
	failure, isInstallFailure := err.(*installFailure)
	reason := err.Error()
	var paths []string
	if isInstallFailure {
		paths = failure.paths
	}
	result := InstallResult{
		SuiteVersion:  transaction.suiteVersion,
		FailureReason: reason,
		FailurePaths:  paths,
		Outcome:       resultmodel.OutcomeFailure,
		Rollback:      resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
	}
	// cleanup runs from the deferred call and decides whether recovery was needed; the
	// outcome below is finalised there through recoveryOutcome.
	transaction.runRecoveryIfNeeded()
	if transaction.recoveryRan {
		result.Rollback = resultmodel.RollbackResult{
			Status:  resultmodel.RollbackSucceeded,
			Actions: []string{"restored every managed path and the Git index to their exact pre-install state"},
		}
		result.Outcome = resultmodel.OutcomeRolledBack
		if transaction.recoveryFailed {
			result.Rollback = resultmodel.RollbackResult{
				Status: resultmodel.RollbackIncomplete,
				Errors: []string{"automatic recovery was incomplete; the managed paths and the Git index need a person"},
			}
			result.Outcome = resultmodel.OutcomeRisk
		}
	}
	return result
}

func (transaction *installTransaction) recordedChanges() []resultmodel.RecordedChange {
	changes := make([]resultmodel.RecordedChange, 0, len(transaction.modules)+3)
	for _, module := range transaction.modules {
		changes = append(changes, resultmodel.RecordedChange{
			Path: module.relativePath, Kind: changeKindFor(module.existed),
			Detail: "installed do-work suite v" + transaction.suiteVersion,
		})
	}
	changes = append(changes,
		resultmodel.RecordedChange{
			Path: transaction.relativeToProject(transaction.justTarget), Kind: changeKindFor(transaction.justExisted),
			Detail: "managed Just recipe section reconciled",
		},
		resultmodel.RecordedChange{
			Path: transaction.relativeToProject(transaction.instructionsTarget), Kind: changeKindFor(transaction.instructionsExisted),
			Detail: "managed agent-instructions section reconciled",
		},
		resultmodel.RecordedChange{
			Path: ".claude/settings.json", Kind: changeKindFor(transaction.settingsExisted),
			Detail: "core hooks composed into existing settings",
		})
	return changes
}

func changeKindFor(existed bool) string {
	if existed {
		return managedsection.ChangeKindModified
	}
	return managedsection.ChangeKindCreated
}

// prepare runs every guard and builds every candidate. Nothing under it writes a client
// path, so a failure here leaves the project exactly as it was found.
func (transaction *installTransaction) prepare(ctx context.Context) error {
	if err := transaction.resolveProjectRoot(ctx); err != nil {
		return err
	}
	if err := transaction.createWorkingDirectories(); err != nil {
		return err
	}
	if err := transaction.obtainAndValidateSource(ctx); err != nil {
		return err
	}
	if err := transaction.guardModulePlan(); err != nil {
		return err
	}
	if err := transaction.buildJustCandidate(ctx); err != nil {
		return err
	}
	if err := transaction.buildInstructionsCandidate(); err != nil {
		return err
	}
	return transaction.buildSettingsCandidate()
}

func (transaction *installTransaction) resolveProjectRoot(ctx context.Context) error {
	projectRoot := transaction.options.ProjectRoot
	if projectRoot == "" {
		return failInstall("--project-root is required")
	}
	info, err := os.Stat(projectRoot)
	if err != nil || !info.IsDir() {
		return failInstall("project root does not exist: %s", projectRoot)
	}
	physicalRoot, err := filepath.EvalSymlinks(projectRoot)
	if err != nil {
		return failInstall("project root does not exist: %s", projectRoot)
	}
	physicalRoot, err = filepath.Abs(physicalRoot)
	if err != nil {
		return failInstall("project root does not exist: %s", projectRoot)
	}
	transaction.projectRoot = physicalRoot

	gitRoot, err := runCapturing(ctx, physicalRoot, "git", "rev-parse", "--show-toplevel")
	if err != nil || strings.TrimSpace(gitRoot) == "" {
		return failInstall("--project-root must be a Git repository so recovery is deterministic")
	}
	physicalGitRoot, err := filepath.EvalSymlinks(strings.TrimSpace(gitRoot))
	if err != nil {
		return failInstall("--project-root must be a Git repository so recovery is deterministic")
	}
	physicalGitRoot, _ = filepath.Abs(physicalGitRoot)
	if physicalGitRoot != physicalRoot {
		return failInstall("--project-root must name the Git worktree root (%s)", physicalGitRoot)
	}
	indexPath, err := runCapturing(ctx, physicalRoot, "git", "rev-parse", "--git-path", "index")
	if err != nil {
		return failInstall("could not resolve the project Git index")
	}
	transaction.gitIndexPath = strings.TrimSpace(indexPath)
	if !filepath.IsAbs(transaction.gitIndexPath) {
		transaction.gitIndexPath = filepath.Join(physicalRoot, transaction.gitIndexPath)
	}

	transaction.settingsTarget = filepath.Join(physicalRoot, ".claude", "settings.json")
	transaction.instructionsTarget = filepath.Join(physicalRoot, "CLAUDE.md")
	return nil
}

func (transaction *installTransaction) createWorkingDirectories() error {
	installTmp, err := os.MkdirTemp("", "do-work-suite-install.*")
	if err != nil {
		return failInstall("could not allocate a private install workspace")
	}
	transaction.installTmp = installTmp
	transaction.sourceRoot = filepath.Join(installTmp, "source")
	transaction.backupRoot = filepath.Join(installTmp, "originals")
	for _, directory := range []string{transaction.sourceRoot, transaction.backupRoot} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			return failInstall("could not allocate a private install workspace")
		}
	}
	return nil
}

// obtainAndValidateSource fetches or accepts an archive, extracts it, and validates the
// manifest before anything downstream reads a byte of it.
func (transaction *installTransaction) obtainAndValidateSource(ctx context.Context) error {
	if transaction.options.ExtractedSourceRoot != "" {
		transaction.sourceRoot = transaction.options.ExtractedSourceRoot
	} else {
		archivePath, err := transaction.resolveArchive(ctx)
		if err != nil {
			return err
		}
		if err := runQuietly(ctx, "", "tar", "xzf", archivePath, "-C", transaction.sourceRoot, "--strip-components=1"); err != nil {
			return failInstall("suite archive extraction failed; no client files were changed")
		}
	}
	validation, err := suitemanifest.ValidateSuite(transaction.sourceRoot)
	if err != nil {
		return failInstall("suite manifest validation failed; no client files were changed: %v", err)
	}
	transaction.suiteVersion = validation.SuiteVersion
	transaction.modules = make([]moduleInstallPlan, 0, len(validation.Modules))
	for index, module := range validation.Modules {
		transaction.modules = append(transaction.modules, moduleInstallPlan{
			sourcePath:      filepath.Join(transaction.sourceRoot, module.Source),
			destinationPath: filepath.Join(transaction.projectRoot, module.Destination),
			relativePath:    module.Destination,
			backupPath:      filepath.Join(transaction.backupRoot, "modules", fmt.Sprint(index)),
		})
	}
	if len(transaction.modules) != 4 {
		return failInstall("validated suite did not produce exactly four install targets")
	}
	return nil
}

func (transaction *installTransaction) resolveArchive(ctx context.Context) (string, error) {
	if supplied := transaction.options.SuppliedArchivePath; supplied != "" {
		info, err := os.Lstat(supplied)
		if err != nil || !info.Mode().IsRegular() {
			return "", failInstall("--archive must name a regular file: %s", supplied)
		}
		absolute, err := filepath.Abs(supplied)
		if err != nil {
			return "", failInstall("--archive must name a regular file: %s", supplied)
		}
		return absolute, nil
	}
	archivePath := filepath.Join(transaction.installTmp, "upstream.tar.gz")
	upstreamURL := transaction.options.UpstreamURL
	if upstreamURL == "" {
		upstreamURL = archivefetch.UpstreamURLFromEnvironment()
	}
	if _, err := archivefetch.FetchArchive(ctx, archivefetch.Request{
		ArchiveTargetPath:    archivePath,
		UpstreamTarballURL:   upstreamURL,
		AtomicDownloadScript: archivefetch.LocateAtomicDownloadScript(transaction.options.ToolDirectory),
	}); err != nil {
		return "", failInstall("upstream archive could not be fetched by any route; no client files were changed")
	}
	return archivePath, nil
}

// guardModulePlan rejects every shape that could redirect a managed write. The manifest
// validator already constrained the destinations textually; here the nearest existing parent
// is resolved physically as well, so a project-local symlink cannot escape the project.
func (transaction *installTransaction) guardModulePlan() error {
	for index := range transaction.modules {
		module := &transaction.modules[index]
		info, err := os.Lstat(module.destinationPath)
		switch {
		case err == nil && info.Mode()&os.ModeSymlink != 0:
			return &installFailure{
				reason: fmt.Sprintf("managed destination must not be a symlink: %s", module.destinationPath),
				paths:  []string{module.relativePath},
			}
		case err == nil && !info.IsDir():
			return &installFailure{
				reason: fmt.Sprintf("managed destination must be a directory when it exists: %s", module.destinationPath),
				paths:  []string{module.relativePath},
			}
		case err == nil:
			module.existed = true
		}
		destinationParent := filepath.Dir(module.destinationPath)
		for {
			if info, statErr := os.Stat(destinationParent); statErr == nil && info.IsDir() {
				break
			}
			nextParent := filepath.Dir(destinationParent)
			if nextParent == destinationParent {
				return failInstall("cannot resolve managed destination parent: %s", module.destinationPath)
			}
			destinationParent = nextParent
		}
		physicalParent, err := filepath.EvalSymlinks(destinationParent)
		if err != nil {
			return failInstall("cannot resolve managed destination parent: %s", module.destinationPath)
		}
		if physicalParent != transaction.projectRoot &&
			!strings.HasPrefix(physicalParent, transaction.projectRoot+string(filepath.Separator)) {
			return failInstall("managed destination resolves outside the project: %s", module.destinationPath)
		}

		sourceInfo, err := os.Lstat(module.sourcePath)
		if err != nil || !sourceInfo.IsDir() {
			return failInstall("managed module source is not a real directory: %s", module.sourcePath)
		}
		if containsSymlink(module.sourcePath) {
			return failInstall("managed module source contains a symlink: %s", module.sourcePath)
		}
	}
	return nil
}

func containsSymlink(root string) bool {
	found := false
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

// buildJustCandidate reconciles the managed recipe section into a candidate held outside the
// project, validates it, and only then becomes eligible to be written.
func (transaction *installTransaction) buildJustCandidate(ctx context.Context) error {
	boardTemplate := filepath.Join(transaction.sourceRoot, "skills", "do-work-board", "justfile.template")
	if err := requireNonEmptyRegularFile(boardTemplate); err != nil {
		return failInstall("board Justfile template is missing or unsafe")
	}
	if err := requireNonEmptyRegularFile(transaction.coreHooksPath()); err != nil {
		return failInstall("core hook fragment is missing or unsafe")
	}
	if err := requireNonEmptyRegularFile(transaction.instructionsTemplatePath()); err != nil {
		return failInstall("core agent instructions template is missing or unsafe")
	}

	if err := transaction.discoverJustTarget(); err != nil {
		return err
	}
	managedSectionPath := filepath.Join(transaction.installTmp, "do-work-recipes.just")
	if err := extractManagedSection(boardTemplate, managedSectionPath); err != nil {
		return failInstall("board template does not contain one complete managed recipe section")
	}

	transaction.justCandidate = filepath.Join(transaction.installTmp, "justfile.candidate")
	if transaction.justExisted {
		if err := copyPreservingMode(ctx, transaction.justTarget, transaction.justCandidate); err != nil {
			return failInstall("Justfile ownership validation failed; no client files were changed")
		}
		if _, err := managedsection.ReplaceSection(managedsection.ReplaceRequest{
			TargetPath:             transaction.justCandidate,
			SectionFilePath:        managedSectionPath,
			RejectRecipeCollisions: true,
		}); err != nil {
			return &installFailure{
				reason: fmt.Sprintf("Justfile ownership validation failed; no client files were changed: %v", err),
				paths:  []string{transaction.relativeToProject(transaction.justTarget)},
			}
		}
	} else if _, err := managedsection.ReplaceSection(managedsection.ReplaceRequest{
		TargetPath:       transaction.justCandidate,
		SectionFilePath:  managedSectionPath,
		TemplateFilePath: boardTemplate,
	}); err != nil {
		return failInstall("complete Justfile candidate creation failed; no client files were changed: %v", err)
	}
	return transaction.validateJustCandidate(ctx, transaction.justCandidate,
		"Justfile candidate does not parse; no client files were changed")
}

func (transaction *installTransaction) validateJustCandidate(ctx context.Context, candidatePath, parseFailureMessage string) error {
	candidateData, err := os.ReadFile(candidatePath)
	if err != nil {
		return failInstall("Justfile candidate could not be read")
	}
	if countExactLines(candidateData, managedsection.DefaultBeginMarker) != 1 ||
		countExactLines(candidateData, managedsection.DefaultEndMarker) != 1 {
		return failInstall("Justfile candidate has invalid managed markers")
	}
	definedNames := managedsection.JustDefinitionNames(candidateData)
	for _, recipeName := range reservedRecipeNames {
		if _, defined := definedNames[recipeName]; !defined {
			return failInstall("Justfile candidate is missing %s", recipeName)
		}
	}
	if _, err := exec.LookPath("just"); err == nil {
		if err := runQuietly(ctx, "", "just", "--justfile", candidatePath, "--list"); err != nil {
			return failInstall("%s", parseFailureMessage)
		}
	}
	return nil
}

// discoverJustTarget records the real directory-entry spelling, so a failed install restores
// `Justfile` as `Justfile` and never creates a second `justfile` beside it (REQ-180).
func (transaction *installTransaction) discoverJustTarget() error {
	entries, err := os.ReadDir(transaction.projectRoot)
	if err != nil {
		return failInstall("could not resolve the real Justfile directory entry in %s", transaction.projectRoot)
	}
	for _, candidateName := range justfileCandidateNames {
		for _, entry := range entries {
			if entry.Name() != candidateName {
				continue
			}
			candidatePath := filepath.Join(transaction.projectRoot, entry.Name())
			info, statErr := os.Lstat(candidatePath)
			if statErr != nil || !info.Mode().IsRegular() {
				return failInstall("Justfile target must be a regular file: %s", candidatePath)
			}
			transaction.justTarget = candidatePath
			transaction.justExisted = true
			return nil
		}
	}
	transaction.justTarget = filepath.Join(transaction.projectRoot, "justfile")
	transaction.justExisted = false
	return nil
}

// buildInstructionsCandidate reconciles CLAUDE.md's managed section. The shipped template is
// itself exactly one managed section, so it doubles as the section file and the create-from
// template on both paths.
func (transaction *installTransaction) buildInstructionsCandidate() error {
	instructionsTemplate := transaction.instructionsTemplatePath()
	if info, err := os.Lstat(transaction.instructionsTarget); err == nil {
		if !info.Mode().IsRegular() {
			return failInstall("agent instructions target must be a regular file: %s", transaction.instructionsTarget)
		}
		transaction.instructionsExisted = true
	}

	transaction.instructionsCandidate = filepath.Join(transaction.installTmp, "agent-instructions.candidate")
	if transaction.instructionsExisted {
		if err := copyPreservingMode(context.Background(), transaction.instructionsTarget, transaction.instructionsCandidate); err != nil {
			return failInstall("agent instructions reconciliation failed; no client files were changed")
		}
	}
	if _, err := managedsection.ReplaceSection(managedsection.ReplaceRequest{
		TargetPath:       transaction.instructionsCandidate,
		SectionFilePath:  instructionsTemplate,
		TemplateFilePath: instructionsTemplate,
		BeginMarker:      instructionsBeginMarker,
		EndMarker:        instructionsEndMarker,
	}); err != nil {
		return &installFailure{
			reason: fmt.Sprintf("agent instructions reconciliation failed; no client files were changed: %v", err),
			paths:  []string{transaction.relativeToProject(transaction.instructionsTarget)},
		}
	}
	candidateData, err := os.ReadFile(transaction.instructionsCandidate)
	if err != nil {
		return failInstall("agent instructions candidate could not be read")
	}
	if bytes.Count(candidateData, []byte(instructionsBeginMarker)) != 1 ||
		bytes.Count(candidateData, []byte(instructionsEndMarker)) != 1 {
		return failInstall("agent instructions candidate has invalid managed markers")
	}
	if !bytes.Contains(candidateData, []byte("crew-members/communication-style.md")) {
		return failInstall("agent instructions candidate does not link the communication-style crew member")
	}
	return nil
}

// buildSettingsCandidate composes the core hooks into the consumer's settings. With Go
// always able to reconcile, "no JSON tool available" is no longer a reachable state, so the
// old jq/python3/manual three-way branch and its manual instruction are gone.
func (transaction *installTransaction) buildSettingsCandidate() error {
	settingsMode := os.FileMode(0o644)
	settingsData := []byte("{}\n")
	if info, err := os.Lstat(transaction.settingsTarget); err == nil {
		if !info.Mode().IsRegular() {
			return failInstall("Claude settings must be a regular file: %s", transaction.settingsTarget)
		}
		transaction.settingsExisted = true
		settingsMode = info.Mode().Perm() | (info.Mode() & (os.ModeSetuid | os.ModeSetgid | os.ModeSticky))
		settingsData, err = os.ReadFile(transaction.settingsTarget)
		if err != nil {
			return failInstall("Claude settings could not be read: %s", transaction.settingsTarget)
		}
	}
	fragmentData, err := os.ReadFile(transaction.coreHooksPath())
	if err != nil {
		return failInstall("core hook fragment is missing or unsafe")
	}
	composed, err := settingshooks.ComposeSettings(settingsData, fragmentData)
	if err != nil {
		return &installFailure{
			reason: fmt.Sprintf("Claude settings are invalid or cannot accept composed core hooks; no client files were changed: %v", err),
			paths:  []string{".claude/settings.json"},
		}
	}
	if !bytes.Contains(composed, []byte("do-work/hooks/session-start.sh")) {
		return failInstall("composed settings omitted the core SessionStart hook")
	}
	if bytes.Contains(composed, []byte("do-work/hooks/pipeline-guard.sh")) {
		return failInstall("composed settings retained the retired pipeline Stop hook")
	}
	transaction.settingsCandidate = filepath.Join(transaction.installTmp, "settings.candidate.json")
	if err := os.WriteFile(transaction.settingsCandidate, composed, settingsMode); err != nil {
		return failInstall("composed Claude settings could not be staged")
	}
	return os.Chmod(transaction.settingsCandidate, settingsMode)
}

// reviewAndConfirm prints the whole managed plan and diff, then takes the one confirmation.
// Everything here goes to the narration stream so stdout carries only the rendered result.
func (transaction *installTransaction) reviewAndConfirm(ctx context.Context) (bool, error) {
	transaction.narrate("Ready to install do-work suite v%s into %s:", transaction.suiteVersion, transaction.projectRoot)
	for _, module := range transaction.modules {
		transaction.narrate("  %s", module.relativePath)
	}
	transaction.narrate("  Justfile: %s", transaction.relativeToProject(transaction.justTarget))
	transaction.narrate("  agent instructions: %s", transaction.relativeToProject(transaction.instructionsTarget))

	transaction.narrate("Reviewing the complete managed install before overwrite:")
	for _, module := range transaction.modules {
		transaction.narrate("")
		transaction.narrate("--- managed destination: %s ---", module.relativePath)
		if err := transaction.narrateDiff(ctx, "could not compare managed destination "+module.relativePath,
			"diff", "-ruN", module.destinationPath, module.sourcePath); err != nil {
			return false, err
		}
	}
	for _, configuration := range []struct {
		label         string
		existingPath  string
		existed       bool
		candidatePath string
		failureReason string
	}{
		{transaction.relativeToProject(transaction.justTarget), transaction.justTarget, transaction.justExisted,
			transaction.justCandidate, "could not compare the managed Justfile candidate"},
		{transaction.relativeToProject(transaction.instructionsTarget), transaction.instructionsTarget, transaction.instructionsExisted,
			transaction.instructionsCandidate, "could not compare the agent instructions candidate"},
		{".claude/settings.json", transaction.settingsTarget, transaction.settingsExisted,
			transaction.settingsCandidate, "could not compare the Claude settings candidate"},
	} {
		transaction.narrate("")
		transaction.narrate("--- managed configuration: %s ---", configuration.label)
		leftPath := "/dev/null"
		if configuration.existed {
			leftPath = configuration.existingPath
		}
		if err := transaction.narrateDiff(ctx, configuration.failureReason,
			"diff", "-u", leftPath, configuration.candidatePath); err != nil {
			return false, err
		}
	}

	if err := transaction.warnAboutDirtyManagedPaths(ctx); err != nil {
		return false, err
	}
	transaction.narrateWithoutNewline("Install this complete four-skill suite? [y/N] ")
	return transaction.readConfirmation(), nil
}

// narrateDiff treats a diff status above 1 as a hard failure: 1 means "files differ", which
// is the normal case, and anything higher means the comparison itself did not run.
func (transaction *installTransaction) narrateDiff(ctx context.Context, failureReason string, name string, arguments ...string) error {
	command := exec.CommandContext(ctx, name, arguments...)
	output, err := command.CombinedOutput()
	transaction.writeNarration(output)
	if err == nil {
		return nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return nil
	}
	return failInstall("%s", failureReason)
}

func (transaction *installTransaction) warnAboutDirtyManagedPaths(ctx context.Context) error {
	statusArgs := []string{"status", "--porcelain", "--"}
	for _, module := range transaction.modules {
		statusArgs = append(statusArgs, module.relativePath)
	}
	statusArgs = append(statusArgs,
		transaction.relativeToProject(transaction.justTarget),
		".claude/settings.json",
		transaction.relativeToProject(transaction.instructionsTarget))
	dirtyManaged, err := runCapturing(ctx, transaction.projectRoot, "git", statusArgs...)
	if err != nil {
		return failInstall("could not inspect managed path status")
	}
	if strings.TrimSpace(dirtyManaged) == "" {
		return nil
	}
	transaction.narrate("Managed install paths have uncommitted changes. Continuing discards those changes "+
		"in managed modules and replaces only the owned configuration bytes:\n%s", strings.TrimRight(dirtyManaged, "\n"))
	return nil
}

func (transaction *installTransaction) readConfirmation() bool {
	if transaction.options.ConfirmationInput == nil {
		return false
	}
	reader := bufio.NewReader(transaction.options.ConfirmationInput)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return false
	}
	switch strings.TrimRight(line, "\r\n") {
	case "y", "Y", "yes", "YES":
		return true
	default:
		return false
	}
}

// writeAndVerify takes the backups, then writes, then proves every installed byte. Recovery
// becomes possible the moment writeStarted is set, which is BEFORE the unstage loop so an
// index mutation already made is recovered too.
func (transaction *installTransaction) writeAndVerify(ctx context.Context) error {
	if err := transaction.takeBackups(ctx); err != nil {
		return err
	}
	transaction.writeStarted = true

	// The confirmation authorised discarding dirty managed module content. Clear only module
	// paths from the index, so a previously staged customisation cannot survive beneath
	// installed bytes; Just and settings stay project configuration and are never reset.
	for _, module := range transaction.modules {
		tracked, err := runCapturing(ctx, transaction.projectRoot, "git", "ls-files", "--", module.relativePath)
		if err != nil {
			return failInstall("could not inspect the Git index for %s", module.relativePath)
		}
		if strings.TrimSpace(tracked) == "" {
			continue
		}
		if err := runQuietly(ctx, transaction.projectRoot, "git", "restore", "--staged", "--", module.relativePath); err != nil {
			return failInstall("could not unstage managed module path %s", module.relativePath)
		}
	}

	for _, module := range transaction.modules {
		if err := os.RemoveAll(module.destinationPath); err != nil {
			return failInstall("could not clear managed destination %s", module.relativePath)
		}
		if err := os.MkdirAll(module.destinationPath, 0o755); err != nil {
			return failInstall("could not create managed destination %s", module.relativePath)
		}
		if err := runQuietly(ctx, "", "cp", "-Rp", module.sourcePath+string(filepath.Separator)+".",
			module.destinationPath+string(filepath.Separator)); err != nil {
			return failInstall("could not install managed module %s", module.relativePath)
		}
	}

	// Each configuration file is replaced from the candidate this transaction already built
	// and validated, so the install never executes a helper it just downloaded.
	for _, write := range []struct {
		candidatePath string
		targetPath    string
		temporaryStem string
	}{
		{transaction.justCandidate, transaction.justTarget, ".do-work-just.install.*"},
		{transaction.settingsCandidate, transaction.settingsTarget, ".settings.json.install.*"},
		{transaction.instructionsCandidate, transaction.instructionsTarget, ".do-work-instructions.install.*"},
	} {
		if err := publishCandidate(ctx, write.candidatePath, write.targetPath, write.temporaryStem); err != nil {
			return failInstall("could not write managed configuration %s", transaction.relativeToProject(write.targetPath))
		}
	}

	return transaction.verifyInstalledBytes(ctx)
}

func (transaction *installTransaction) verifyInstalledBytes(ctx context.Context) error {
	for _, module := range transaction.modules {
		if err := runQuietly(ctx, "", "diff", "-qr", module.sourcePath, module.destinationPath); err != nil {
			return failInstall("installed bytes do not match %s", module.relativePath)
		}
	}
	if !filesAreIdentical(transaction.justCandidate, transaction.justTarget) {
		return failInstall("installed Justfile does not match its validated candidate")
	}
	if !filesAreIdentical(transaction.instructionsCandidate, transaction.instructionsTarget) {
		return failInstall("installed agent instructions do not match their validated candidate")
	}
	if _, err := exec.LookPath("just"); err == nil {
		if err := runQuietly(ctx, "", "just", "--justfile", transaction.justTarget, "--list"); err != nil {
			return failInstall("installed Justfile failed post-write validation")
		}
	}
	if !filesAreIdentical(transaction.settingsCandidate, transaction.settingsTarget) {
		return failInstall("installed Claude settings do not match the validated candidate")
	}
	installedSettings, err := os.ReadFile(transaction.settingsTarget)
	if err != nil {
		return failInstall("installed Claude settings failed post-write validation")
	}
	if _, err := settingshooks.ComposeSettings(installedSettings, mustReadFile(transaction.coreHooksPath())); err != nil {
		return failInstall("installed Claude settings failed post-write validation")
	}

	installedVersion, err := suitemanifest.ReadSuiteVersion(
		filepath.Join(transaction.projectRoot, ".claude", "skills", "do-work", "VERSION"))
	if err != nil || installedVersion != transaction.suiteVersion {
		return failInstall("installed version mismatch (expected %s, found %s)",
			transaction.suiteVersion, orUnknown(installedVersion))
	}
	actionVersion, markerCount, err := suitemanifest.ReadActionVersion(
		filepath.Join(transaction.projectRoot, ".claude", "skills", "do-work", "actions", "version.md"))
	if err != nil {
		return failInstall("installed actions/version.md must contain exactly one Current version line")
	}
	if markerCount != 1 {
		return failInstall("installed actions/version.md must contain exactly one Current version line")
	}
	if actionVersion != transaction.suiteVersion {
		return failInstall("installed action version mismatch (expected %s, found %s)",
			transaction.suiteVersion, orUnknown(actionVersion))
	}
	return nil
}

func (transaction *installTransaction) takeBackups(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Join(transaction.backupRoot, "modules"), 0o755); err != nil {
		return failInstall("could not allocate the recovery snapshot")
	}
	for index := range transaction.modules {
		module := &transaction.modules[index]
		if !module.existed {
			continue
		}
		if err := os.MkdirAll(module.backupPath, 0o755); err != nil {
			return failInstall("could not snapshot %s", module.relativePath)
		}
		if err := runQuietly(ctx, "", "cp", "-Rp", module.destinationPath+string(filepath.Separator)+".",
			module.backupPath+string(filepath.Separator)); err != nil {
			return failInstall("could not snapshot %s", module.relativePath)
		}
	}
	for _, backup := range []struct {
		existed    bool
		sourcePath string
		backupName string
	}{
		{transaction.justExisted, transaction.justTarget, "justfile"},
		{transaction.settingsExisted, transaction.settingsTarget, "settings.json"},
		{transaction.instructionsExisted, transaction.instructionsTarget, "agent-instructions"},
	} {
		if !backup.existed {
			continue
		}
		if err := copyPreservingMode(ctx, backup.sourcePath, filepath.Join(transaction.backupRoot, backup.backupName)); err != nil {
			return failInstall("could not snapshot %s", backup.sourcePath)
		}
	}
	if info, err := os.Lstat(transaction.gitIndexPath); err == nil {
		if !info.Mode().IsRegular() {
			return failInstall("Git index must be a regular file: %s", transaction.gitIndexPath)
		}
		transaction.gitIndexExisted = true
		if err := copyPreservingMode(ctx, transaction.gitIndexPath, filepath.Join(transaction.backupRoot, "git-index")); err != nil {
			return failInstall("could not snapshot the Git index")
		}
	}
	return nil
}

// runRecoveryIfNeeded restores in a fixed order — modules, Justfile, settings, agent
// instructions, then the Git index — and reports a single failure flag rather than stopping
// at the first problem, because a partial restore is still better than none.
func (transaction *installTransaction) runRecoveryIfNeeded() {
	if !transaction.writeStarted || transaction.installVerified || transaction.recoveryRan {
		return
	}
	transaction.recoveryRan = true
	transaction.narrate("do-work suite install: installation did not complete; recovering managed paths and Git index.")
	ctx := context.Background()

	for _, module := range transaction.modules {
		if err := os.RemoveAll(module.destinationPath); err != nil {
			transaction.recoveryFailed = true
		}
		if !module.existed {
			continue
		}
		if err := os.MkdirAll(module.destinationPath, 0o755); err != nil {
			transaction.recoveryFailed = true
			continue
		}
		if err := runQuietly(ctx, "", "cp", "-Rp", module.backupPath+string(filepath.Separator)+".",
			module.destinationPath+string(filepath.Separator)); err != nil {
			transaction.recoveryFailed = true
		}
	}
	for _, restore := range []struct {
		targetPath string
		existed    bool
		backupName string
	}{
		{transaction.justTarget, transaction.justExisted, "justfile"},
		{transaction.settingsTarget, transaction.settingsExisted, "settings.json"},
		{transaction.instructionsTarget, transaction.instructionsExisted, "agent-instructions"},
	} {
		if restore.targetPath == "" {
			continue
		}
		if err := os.RemoveAll(restore.targetPath); err != nil {
			transaction.recoveryFailed = true
		}
		if !restore.existed {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(restore.targetPath), 0o755); err != nil {
			transaction.recoveryFailed = true
			continue
		}
		if err := copyPreservingMode(ctx, filepath.Join(transaction.backupRoot, restore.backupName), restore.targetPath); err != nil {
			transaction.recoveryFailed = true
		}
	}
	if transaction.gitIndexExisted {
		if err := publishCandidate(ctx, filepath.Join(transaction.backupRoot, "git-index"),
			transaction.gitIndexPath, ".do-work-index.restore.*"); err != nil {
			transaction.recoveryFailed = true
		}
	} else if err := os.Remove(transaction.gitIndexPath); err != nil && !os.IsNotExist(err) {
		transaction.recoveryFailed = true
	}

	if transaction.recoveryFailed {
		transaction.narrate("do-work suite install: automatic recovery was incomplete; inspect the four skill "+
			"directories, %s, %s, %s, and Git index %s",
			transaction.justTarget, transaction.settingsTarget, transaction.instructionsTarget, transaction.gitIndexPath)
		return
	}
	transaction.narrate("do-work suite install: restored every managed path and the Git index to their exact pre-install state.")
}

func (transaction *installTransaction) cleanup() {
	transaction.runRecoveryIfNeeded()
	if transaction.installTmp != "" {
		_ = os.RemoveAll(transaction.installTmp)
		transaction.installTmp = ""
	}
}

func (transaction *installTransaction) coreHooksPath() string {
	return filepath.Join(transaction.sourceRoot, "skills", "do-work", "hooks", "hooks.json")
}

func (transaction *installTransaction) instructionsTemplatePath() string {
	return filepath.Join(transaction.sourceRoot, "skills", "do-work", "agent-instructions.template.md")
}

func (transaction *installTransaction) relativeToProject(path string) string {
	relative, err := filepath.Rel(transaction.projectRoot, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(relative)
}

func (transaction *installTransaction) narrate(format string, arguments ...any) {
	transaction.writeNarration([]byte(fmt.Sprintf(format, arguments...) + "\n"))
}

func (transaction *installTransaction) narrateWithoutNewline(text string) {
	transaction.writeNarration([]byte(text))
}

func (transaction *installTransaction) writeNarration(data []byte) {
	if transaction.options.Narration == nil || len(data) == 0 {
		return
	}
	_, _ = transaction.options.Narration.Write(data)
}

// extractManagedSection copies the board template's managed span into its own file, which is
// what the section replacer consumes. The template is entirely one section today, so this is
// a byte-identical copy — but reading the span rather than the whole file keeps that an
// observation instead of an assumption.
func extractManagedSection(templatePath, sectionPath string) error {
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}
	beginOffset := indexOfExactLine(templateData, managedsection.DefaultBeginMarker)
	endOffset := indexOfExactLine(templateData, managedsection.DefaultEndMarker)
	if beginOffset < 0 || endOffset < 0 || endOffset < beginOffset {
		return errors.New("board template does not contain one complete managed recipe section")
	}
	endOfSection := endOffset + len(managedsection.DefaultEndMarker)
	if endOfSection < len(templateData) && templateData[endOfSection] == '\n' {
		endOfSection++
	}
	return os.WriteFile(sectionPath, templateData[beginOffset:endOfSection], 0o644)
}

// indexOfExactLine finds a line equal to the marker, never a line merely containing it.
func indexOfExactLine(data []byte, marker string) int {
	offset := 0
	for _, line := range bytes.SplitAfter(data, []byte("\n")) {
		if string(bytes.TrimRight(line, "\r\n")) == marker {
			return offset
		}
		offset += len(line)
	}
	return -1
}

func countExactLines(data []byte, marker string) int {
	count := 0
	for _, line := range bytes.SplitAfter(data, []byte("\n")) {
		if string(bytes.TrimRight(line, "\r\n")) == marker {
			count++
		}
	}
	return count
}

// publishCandidate stages beside the target and renames, so a target is never observed
// half-written. `cp -p` carries the candidate's mode onto the staged file the way the shell
// installer did.
func publishCandidate(ctx context.Context, candidatePath, targetPath, temporaryStem string) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	temporaryFile, err := os.CreateTemp(filepath.Dir(targetPath), temporaryStem)
	if err != nil {
		return err
	}
	temporaryPath := temporaryFile.Name()
	_ = temporaryFile.Close()
	if err := runQuietly(ctx, "", "cp", "-p", candidatePath, temporaryPath); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	if err := os.Rename(temporaryPath, targetPath); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	return nil
}

func copyPreservingMode(ctx context.Context, sourcePath, destinationPath string) error {
	return runQuietly(ctx, "", "cp", "-p", sourcePath, destinationPath)
}

func filesAreIdentical(firstPath, secondPath string) bool {
	firstData, firstErr := os.ReadFile(firstPath)
	secondData, secondErr := os.ReadFile(secondPath)
	return firstErr == nil && secondErr == nil && bytes.Equal(firstData, secondData)
}

func requireNonEmptyRegularFile(path string) error {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("%s is missing or unsafe", path)
	}
	return nil
}

func mustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func orUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func runQuietly(ctx context.Context, workingDirectory, name string, arguments ...string) error {
	command := exec.CommandContext(ctx, name, arguments...)
	command.Dir = workingDirectory
	return command.Run()
}

func runCapturing(ctx context.Context, workingDirectory, name string, arguments ...string) (string, error) {
	command := exec.CommandContext(ctx, name, arguments...)
	command.Dir = workingDirectory
	var standardOutput bytes.Buffer
	command.Stdout = &standardOutput
	if err := command.Run(); err != nil {
		return "", err
	}
	return standardOutput.String(), nil
}
