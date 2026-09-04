# Exploration — REQ-566 (Explore agent, verified anchors and patterns to copy)

## 1. Command handler pattern — `heavy_commands.go`

`skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go:16-33`

```go
func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandPlanHeavyVerification: handlePlanHeavyVerification,
		CommandPlanHeavyRevalidation: handlePlanHeavyRevalidation,
	}
}

func handlePlanHeavyVerification(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, baseRevision, targetRevision, forceAll, err := parsePlanArguments(arguments)
	if err != nil {
		return planFailure(CommandPlanHeavyVerification, "HEAVY-PLAN-USAGE", err)
	}
	plan, err := Plan(executionContext.RepositoryRoot, manifestPath, baseRevision, targetRevision, forceAll)
	if err != nil {
		return planFailure(CommandPlanHeavyVerification, "HEAVY-PLAN-UNVERIFIABLE", err)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, HeavyVerification: &plan}
}
```

Refusal shape, `heavy_commands.go:153-161`:

```go
func planFailure(command, code string, err error) resultmodel.CommandResult {
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "heavy verification cannot be planned safely",
		NextArgv:         []string{"git", "status", "--short"},
		VerificationArgv: []string{"do-work-cli", "--format", "json", command},
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{finding}}
}
```

Argument parser conventions, `heavy_commands.go:47-96`: default manifest `_dev/tests/heavy-lanes.json` at line 48; `strings.Cut(argument, "=")` at line 63 accepts both `--opt=value` and `--opt value`; a `seen` map rejects duplicates; unknown options return `unknown <command> option %q`. Registration point: `cmd/do-work-cli/main.go:52`, `for name, handler := range heavyverification.Handlers()`. Constants at `heavy_commands.go:11-14`. `commandruntime.ExecutionContext` is `{RepositoryRoot string; Format resultmodel.OutputFormat}` at `internal/commandruntime/command_runtime.go:14-19`.

## 2. Owned-process-group pattern

`internal/toolboxcommands/report_image_process.go:23-66` — the fail-closed copy target:

```go
func configureOwnedProcess(command *exec.Cmd) error {
	if !ownedprocess.ConfigureGroup(command) {
		return errors.New("report image generation is unavailable: descendant process ownership cannot be proved on this platform")
	}
	return nil
}

func runOwnedProcess(ctx context.Context, directory string, argv ...string) ownedProcessResult {
	command := exec.Command(argv[0], argv[1:]...)
	command.Dir = directory
	if err := configureOwnedProcess(command); err != nil { return ownedProcessResult{Status: 1, Err: err} }
	if err := command.Start(); err != nil { return ownedProcessResult{Status: 1, Err: err} }
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case err := <-done:
		return completedProcessResult(err)
	case <-ctx.Done():
		_ = ownedprocess.TerminateGroup(command.Process.Pid, reportImageGracePeriod)
		<-done
		return ownedProcessResult{Status: 1, Interrupted: true, Err: ctx.Err()}
	}
}
```

Exit-status extraction, same file lines 68-76: `*exec.ExitError` gives `exitError.ExitCode()`, anything else is status 1.

Signatures, `internal/ownedprocess/owned_process_group.go`:
- `func ConfigureGroup(command *exec.Cmd) bool` (line 26) — true means group ownership established; a caller whose safety depends on reaching descendants must fail closed on false.
- `func TerminateGroup(leaderPID int, grace time.Duration) error` (line 35) — blocks until the group is gone; returns `os.ErrProcessDone` when nothing was there.
- `const DefaultGracePeriod = 250 * time.Millisecond` (line 19).

### Section 2 (rest) — build tags and the 124 timeout

Build tags in `internal/ownedprocess/`:
- `owned_process_group_unix.go:1` is `//go:build unix`. Its `configureOwnedGroup` sets `command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}` and returns true.
- `owned_process_group_unsupported.go:1` is `//go:build !unix`. `configureOwnedGroup(*exec.Cmd) bool { return false }`, `terminateOwnedGroup(int, time.Duration) error { return os.ErrProcessDone }`. Named for the condition, not `_windows.go`, on purpose.

Timeout and 124: `internal/nextselection/blocked_probe.go:10-13`:

```go
const (
	BlockedProbeTimeoutStatus = 124
	BlockedProbeLaunchStatus  = 125
)
```

`blocked_probe_unix.go:33-53` is the select the runner should copy:

```go
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case waitError := <-done:
		status := probeExitStatus(waitError)
		cleanupReapedProcessGroup(processGroup)
		return status, nil
	case <-timer.C:
		terminateOwnedProcessGroup(processGroup, syscall.SIGTERM, done)
		return BlockedProbeTimeoutStatus, nil
	case received := <-signalChannel:
		terminateOwnedProcessGroup(processGroup, forwarded, done)
		return 128 + int(forwarded), BlockedProbeInterruption{ExitStatus: status}
	}
```

`probeExitStatus` (lines 92-106) maps a signalled child to `128 + signal` and an ordinary exit to its status, falling back to 125. That package predates `ownedprocess` and inlines `Setpgid` plus a `syscall.Getpgid` isolation check (lines 20-32). The new runner calls `ownedprocess.ConfigureGroup` / `TerminateGroup` instead.

Test shape, `blocked_probe_test.go:1` is `//go:build unix`; lines 24-45:

```go
func TestBlockedProbeTimeoutKillsDescendantGroup(t *testing.T) {
	pidPath := filepath.Join(t.TempDir(), "child.pid")
	probe := fmt.Sprintf("(trap '' TERM; sleep 30) & echo $! > %q; wait", pidPath)
	status, err := RunBlockedProbe([]byte(probe), 1)
	if err != nil || status != BlockedProbeTimeoutStatus { t.Fatalf("status=%d err=%v", status, err) }
	// read pid, then poll syscall.Kill(pid, 0) for 2s; return on error, else t.Fatalf
}
```

## 3. Manifest read at a revision

`internal/heavyverification/heavy_verification.go:201-221`:

```go
func readManifestAtRevision(repositoryRoot, manifestRelativePath, targetRevision string) (laneManifest, error) {
	treeEntry, err := runGit(repositoryRoot, "ls-tree", "-z", targetRevision, "--", manifestRelativePath)
	entryPrefix, entryPath, found := bytes.Cut(bytes.TrimSuffix(treeEntry, []byte{0}), []byte{'\t'})
	entryFields := bytes.Fields(entryPrefix)
	if !found || len(entryFields) != 3 || string(entryPath) != manifestRelativePath { /* missing or ambiguous */ }
	// entryFields[0] mode must be 100644 or 100755; entryFields[1] must be "blob"
	contents, err := runGit(repositoryRoot, "cat-file", "blob", string(entryFields[2]))
	return decodeManifest(contents)
}
```

`decodeManifest` (line 223) sets `DisallowUnknownFields`, rejects a trailing JSON value, and requires `schema_version == 1` (`heavyManifestSchemaVersion`, line 16).

Struct fields, lines 18-35, all unexported (the runner lives in the same package, so it reads them directly):

```go
type laneManifest struct {
	SchemaVersion    int            `json:"schema_version"`
	Lanes            []manifestLane `json:"lanes"`
	NonHeavyCoverage []coverageRule `json:"non_heavy_coverage"`
}
type manifestLane struct {
	ID       string         `json:"id"`
	Argv     []string       `json:"argv"`
	Coverage []coverageRule `json:"coverage"`
}
```

Every lane argv in `_dev/tests/heavy-lanes.json` is `["bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "<lane-id>"]`.

## 4-10. Remaining sites (read them directly; anchors from the plan)
- resultmodel: `internal/resultmodel/result_model.go` — `CommandResult` ~360-380 (`HeavyVerification` at 373), nil-slice normalization ~420-440, text renderer ~720-740; `CommandFinding` fields (Code, Severity, Evidence, Fixability, AutomationStopReason, NextArgv, VerificationArgv); `SeverityWarning` for red/skipped lanes, `SeverityError` for refusals.
- Test helpers: `internal/heavyverification/heavy_verification_test.go:391-430` — `newHeavyTestRepository`, `writeHeavyTestFile`, `commitHeavyTestChanges`.
- Answer test: `internal/publication/answer_test.go:262` (`TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure`), assertion near 307 (`"browser: exit 0"`).
- Shell: `_dev/tests/maintainer-verify.sh` — Node discovery ~554-567, browser discovery ~584-605, `run_heavy_lane` ~662-690, `--self-test` `heavy-no-node` fixture ~380-400 (SKIP line asserted at 388).
- Board: `skills/do-work-board/tools/queue-kanban/timeline.go:525` label `waiting for permission to run heavy tests`; grep confirms no test pins it (re-verify).
- Version drift: `VERSION` and `skills/do-work/VERSION` = 0.274.0; `skills/do-work/actions/version.md` and `CHANGELOG.md` head = 0.274.1. Not yours to fix (release tail).
- CLI stdout is reserved for the JSON result (`internal/commandruntime`): tee lane output to stderr only.
