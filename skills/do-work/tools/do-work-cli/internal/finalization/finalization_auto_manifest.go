package finalization

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// autoManifestOptions are the inputs of finalize --auto-manifest. The judged
// fields (transition, terminal status, message, provenance, release manifest,
// failure fields) always come from the action; the command never invents them.
// Writer, extra paths and the emit path are mechanical overrides.
type autoManifestOptions struct {
	RequestID          string
	Transition         string
	TerminalStatus     string
	MessageFile        string
	Provenance         string
	ImplementationHash string
	ReleaseManifest    string
	FailureType        string
	FailureErrorFile   string
	Writer             string
	EmitPath           string
	ExtraPaths         []string
}

func parseAutoManifestArguments(arguments []string) (autoManifestOptions, error) {
	options := autoManifestOptions{}
	singleValueFlags := map[string]*string{
		"--auto-manifest": &options.RequestID, "--transition": &options.Transition, "--terminal-status": &options.TerminalStatus,
		"--message-file": &options.MessageFile, "--provenance": &options.Provenance, "--implementation-hash": &options.ImplementationHash,
		"--release-manifest": &options.ReleaseManifest, "--failure-type": &options.FailureType, "--failure-error-file": &options.FailureErrorFile,
		"--writer": &options.Writer, "--emit": &options.EmitPath,
	}
	for index := 0; index < len(arguments); index++ {
		flag := arguments[index]
		target, known := singleValueFlags[flag]
		if flag == "--manifest" {
			return autoManifestOptions{}, fmt.Errorf("--auto-manifest and --manifest are mutually exclusive")
		}
		if !known && flag != "--extra-path" {
			return autoManifestOptions{}, fmt.Errorf("unknown finalize option %q", flag)
		}
		index++
		if index >= len(arguments) {
			return autoManifestOptions{}, fmt.Errorf("%s requires one value", flag)
		}
		if flag == "--extra-path" {
			options.ExtraPaths = append(options.ExtraPaths, arguments[index])
			continue
		}
		if *target != "" {
			return autoManifestOptions{}, fmt.Errorf("%s may be supplied only once", flag)
		}
		*target = arguments[index]
	}
	return options, nil
}

// judgedInputProblem names the first missing or contradictory judged input, so
// the refusal tells the action which flag to supply. Field values the manifest
// validator already checks (terminal status words, failure types, hash shape)
// are left to validateManifest.
func judgedInputProblem(options autoManifestOptions) string {
	switch {
	case !requestIDPattern.MatchString(options.RequestID):
		return "--auto-manifest requires an exact REQ-NNN id"
	case options.Transition != "complete" && options.Transition != "fail":
		return "--transition complete or --transition fail is required"
	case options.Transition == "complete" && options.TerminalStatus == "":
		return "--terminal-status is required for --transition complete"
	case options.Transition == "complete" && (options.FailureType != "" || options.FailureErrorFile != ""):
		return "--failure-type and --failure-error-file apply only to --transition fail"
	case options.Transition == "fail" && options.TerminalStatus != "":
		return "--terminal-status applies only to --transition complete"
	case options.Transition == "fail" && (options.FailureType == "" || options.FailureErrorFile == ""):
		return "--failure-type and --failure-error-file are required for --transition fail"
	case options.MessageFile == "":
		return "--message-file is required: the commit message is judged by the action"
	case options.Provenance == "":
		return "--provenance primary_commit or --provenance supplied_commit is required"
	case options.Provenance == ProvenanceSuppliedCommit && options.ImplementationHash == "":
		return "--implementation-hash is required with --provenance supplied_commit"
	}
	return ""
}

// handleAutoManifest builds the mechanical manifest fields from the working REQ,
// the live files and the finalization planner, preflights the tree before any
// write, and then either emits the manifest (--emit) or finalizes through the
// same journal path as finalize --manifest.
func handleAutoManifest(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	repositoryRoot := executionContext.RepositoryRoot
	options, err := parseAutoManifestArguments(arguments)
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-USAGE", err.Error())
	}
	requestID := options.RequestID
	if problem := judgedInputProblem(options); problem != "" {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-AUTO-INPUT", problem)
	}
	commitMessage, err := readJudgedText(repositoryRoot, options.MessageFile)
	if err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-AUTO-INPUT", "--message-file: "+err.Error())
	}
	failureError := ""
	if options.FailureErrorFile != "" {
		if failureError, err = readJudgedText(repositoryRoot, options.FailureErrorFile); err != nil {
			return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-AUTO-INPUT", "--failure-error-file: "+err.Error())
		}
	}

	// Preflight (a): a fresh manifest carries a new completed_at, so it can never
	// match the digest an unfinished journal recorded; only recovery resumes it.
	journalPath, _, err := journalLocations(repositoryRoot, requestID)
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
	}
	if _, statError := os.Lstat(journalPath); statError == nil {
		refusal := autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-JOURNAL-UNFINISHED", "an unfinished finalization journal exists for "+requestID+"; resume it with recover-finalization", journalPath)
		refusal.Findings[0].NextArgv = []string{"do-work-cli", CommandRecoverFinalization}
		return refusal
	}
	// Preflight (b): the finalizer commits from an empty index, so name what is staged.
	staged, err := exec.Command("git", "-C", repositoryRoot, "diff", "--cached", "--name-only").Output()
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", "list staged paths: "+err.Error())
	}
	if stagedPaths := strings.FieldsFunc(string(staged), func(character rune) bool { return character == '\n' }); len(stagedPaths) > 0 {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-INDEX-NOT-EMPTY", "finalization requires an empty index; staged: "+strings.Join(stagedPaths, ", "), stagedPaths...)
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
	}
	workingPaths := []string{}
	for _, requestFile := range snapshot.RequestsByID[requestID] {
		if requestFile.TreeSection == "working" {
			workingPaths = append(workingPaths, "do-work/"+requestFile.RelativePath)
		}
	}
	if len(workingPaths) != 1 {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-REQUEST-NOT-WORKING", fmt.Sprintf("expected exactly one %s file under do-work/working/, found %d", requestID, len(workingPaths)), workingPaths...)
	}
	requestBytes, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(workingPaths[0])))
	if err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-REQUEST-NOT-WORKING", err.Error(), workingPaths[0])
	}
	checkpointBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md"))
	if err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-PREPARE-REFUSED", "read do-work/CHECKPOINT.md: "+err.Error(), "do-work/CHECKPOINT.md")
	}

	// The same instant source as do-work-cli now.
	completedAt := time.Now().UTC().Truncate(time.Second).Format("2006-01-02T15:04:05Z")
	writerLabel := options.Writer
	if writerLabel == "" {
		writerLabel = requeststate.DefaultWriterLabel(snapshot.RepositoryRoot)
	}
	manifest := Manifest{
		RequestID: requestID, RequestPath: workingPaths[0], WriterLabel: writerLabel, Transition: options.Transition,
		TerminalStatus: options.TerminalStatus, FailureError: failureError, FailureType: options.FailureType,
		CompletedAt: completedAt, ExpectedRequestSHA256: digestBytes(requestBytes), ExpectedCheckpointSHA256: digestBytes(checkpointBytes),
		CommitPaths: options.ExtraPaths, CommitMessage: commitMessage,
		ProvenanceMode: options.Provenance, ImplementationHash: options.ImplementationHash, ReleaseManifestPath: options.ReleaseManifest,
	}
	if options.ReleaseManifest != "" {
		manifest.ReleaseAt = completedAt
	}

	// Preflight (c): the draft runs the shared planner core dry to learn the
	// required commit paths. Nil manifest bytes never match a journal digest.
	_, _, requiredPaths, err := prepareManifestJournal(repositoryRoot, manifest, nil, true)
	if requiredPaths == nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-PREPARE-REFUSED", err.Error())
	}
	if manifest.CommitPaths, err = normalizeRepositoryPaths(append(requiredPaths, options.ExtraPaths...)); err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-AUTO-INPUT", "--extra-path: "+err.Error())
	}
	if err := validateManifest(repositoryRoot, manifest); err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-AUTO-INPUT", err.Error())
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
	}

	if options.EmitPath == "" {
		journal, resumed, _, err := prepareManifestJournal(repositoryRoot, manifest, manifestBytes, false)
		if err != nil {
			return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
		}
		return advanceJournal(context.Background(), repositoryRoot, journal, resumed)
	}
	if _, _, _, err := prepareManifestJournal(repositoryRoot, manifest, manifestBytes, true); err != nil {
		return autoManifestRefusal(repositoryRoot, requestID, "FINALIZATION-PREPARE-REFUSED", err.Error())
	}
	emitPath, err := containedOrAbsolute(repositoryRoot, options.EmitPath)
	if err == nil {
		err = os.WriteFile(emitPath, manifestBytes, 0o600)
	}
	if err != nil {
		return commandFailure(repositoryRoot, CommandFinalize, "FINALIZATION-EMIT", err.Error())
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: "FINALIZATION-MANIFEST-EMITTED", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{requestID}, AffectedPaths: []string{emitPath},
		Evidence:   []string{"emitted " + emitPath + " after every finalize check passed without a journal", "commit_paths: " + strings.Join(manifest.CommitPaths, ", ")},
		Fixability: resultmodel.FixabilityAutomatic,
	}}}
}

// readJudgedText reads an action-authored text input. Trailing newlines are
// dropped so a message file written by an editor or heredoc yields the same
// manifest value a hand-built manifest carries.
func readJudgedText(repositoryRoot, path string) (string, error) {
	absolutePath, err := containedOrAbsolute(repositoryRoot, path)
	if err != nil {
		return "", err
	}
	contents, err := os.ReadFile(absolutePath)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(contents), "\r\n"), nil
}

func autoManifestRefusal(repositoryRoot, requestID, code, reason string, paths ...string) resultmodel.CommandResult {
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, AffectedPaths: paths, Evidence: []string{reason},
		Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the finalization manifest was not built; nothing was written",
	}
	if requestID != "" {
		finding.AffectedIDs = []string{requestID}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{finding}}
}
