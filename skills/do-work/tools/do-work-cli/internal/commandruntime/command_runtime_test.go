package commandruntime

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestReadOnlyHandlerRunsOutsideGitWithGlobalOptions(t *testing.T) {
	repositoryRoot := t.TempDir()
	var received ExecutionContext
	var receivedArgs []string
	handlers := map[string]CommandHandler{
		"inspect": func(context ExecutionContext, args []string) resultmodel.CommandResult {
			received = context
			receivedArgs = append([]string(nil), args...)
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
		},
	}
	var stdout bytes.Buffer
	runtime := NewRuntime(&stdout, handlers)
	exitCode := runtime.Run([]string{"--repo-root", repositoryRoot, "--format", "json", "inspect", "one", "two"})
	if exitCode != 0 {
		t.Fatalf("exit code = %d, output:\n%s", exitCode, stdout.String())
	}
	wantRoot, err := filepath.Abs(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if received.RepositoryRoot != wantRoot || received.Format != resultmodel.FormatJSON {
		t.Fatalf("context = %#v, want root %q and JSON", received, wantRoot)
	}
	if !reflect.DeepEqual(receivedArgs, []string{"one", "two"}) {
		t.Fatalf("handler args = %#v", receivedArgs)
	}
	var rendered resultmodel.CommandResult
	if err := json.Unmarshal(stdout.Bytes(), &rendered); err != nil {
		t.Fatalf("decode output: %v\n%s", err, stdout.String())
	}
	if rendered.Command != "inspect" || rendered.RepositoryRoot != wantRoot {
		t.Fatalf("runtime did not supply common result fields: %#v", rendered)
	}
}

func TestRuntimeSetsAsideOwnedSelfReferentialRefusal(t *testing.T) {
	handlers := map[string]CommandHandler{
		"recover-finalization": func(ExecutionContext, []string) resultmodel.CommandResult {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{{
				Code: "FINALIZATION-LIFECYCLE-APPLY", Fixability: resultmodel.FixabilityRefused,
				AffectedIDs: []string{"REQ-456"}, AutomationStopReason: "lifecycle apply refused",
				NextArgv:         []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"},
				VerificationArgv: []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"},
			}}}
		},
	}
	var stdout bytes.Buffer
	exitCode := NewRuntime(&stdout, handlers).Run([]string{"--format", "json", "recover-finalization", "--discover"})
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want findings exit 1", exitCode)
	}
	var rendered resultmodel.CommandResult
	if err := json.Unmarshal(stdout.Bytes(), &rendered); err != nil {
		t.Fatalf("decode output: %v\n%s", err, stdout.String())
	}
	if rendered.Outcome != resultmodel.OutcomeFindings || len(rendered.Findings) != 1 {
		t.Fatalf("runtime result = %#v", rendered)
	}
	if len(rendered.Findings[0].NextArgv) != 0 || len(rendered.Findings[0].VerificationArgv) == 0 {
		t.Fatalf("set-aside command fields = %#v", rendered.Findings[0])
	}
}

func TestRuntimeReportsUsageThroughSelectedRenderer(t *testing.T) {
	var stdout bytes.Buffer
	runtime := NewRuntime(&stdout, nil)
	exitCode := runtime.Run([]string{"--format=json", "unknown"})
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
	var rendered resultmodel.CommandResult
	if err := json.Unmarshal(stdout.Bytes(), &rendered); err != nil {
		t.Fatalf("decode output: %v\n%s", err, stdout.String())
	}
	if rendered.Outcome != resultmodel.OutcomeFailure || len(rendered.Findings) != 1 {
		t.Fatalf("unexpected usage result: %#v", rendered)
	}
	if rendered.Findings[0].Code != "UNKNOWN-COMMAND" {
		t.Fatalf("finding code = %q", rendered.Findings[0].Code)
	}
}

func TestRuntimeRejectsInvalidGlobalOptions(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tests := [][]string{
		{"--format", "yaml", "inspect"},
		{"--repo-root"},
		{"--unknown", "value"},
	}
	for _, args := range tests {
		var stdout bytes.Buffer
		exitCode := NewRuntime(&stdout, nil).Run(args)
		if exitCode != 2 {
			t.Errorf("Run(%q) exit = %d, want 2", args, exitCode)
		}
		if stdout.Len() == 0 {
			t.Errorf("Run(%q) emitted no actionable result (cwd %s)", args, workingDirectory)
		}
	}
}

// Every finding the CLI emits must carry a verification command, the runtime's own usage
// findings included. Driving Run is the only seam a real caller goes through, so the
// completeness contract is asserted there rather than against the finding constructor.
func TestRuntimeFindingsCarryCompleteRemediation(t *testing.T) {
	tests := [][]string{
		{"--format=json", "unknown"},
		{"--format=json", "--repo-root"},
		{"--format=json", "--unknown", "value"},
		{"--format=json"},
	}
	for _, args := range tests {
		var stdout bytes.Buffer
		NewRuntime(&stdout, nil).Run(args)
		var rendered resultmodel.CommandResult
		if err := json.Unmarshal(stdout.Bytes(), &rendered); err != nil {
			t.Fatalf("Run(%q) did not emit JSON: %v\n%s", args, err, stdout.String())
		}
		if len(rendered.Findings) != 1 {
			t.Fatalf("Run(%q) findings = %#v", args, rendered.Findings)
		}
		finding := rendered.Findings[0]
		if finding.Code == "" || finding.Severity == "" || len(finding.Evidence) == 0 ||
			finding.Fixability == "" || finding.AutomationStopReason == "" {
			t.Errorf("Run(%q) finding is incomplete: %#v", args, finding)
		}
		if len(finding.NextArgv) == 0 && finding.NextJustRecipe == "" {
			t.Errorf("Run(%q) finding names no next step: %#v", args, finding)
		}
		if len(finding.VerificationArgv) == 0 {
			t.Errorf("Run(%q) finding names no verification command: %#v", args, finding)
		}
	}
}

// The documented 0-4 exit codes and the text/JSON parity claim are only worth anything at
// the seam a real caller goes through: global options parsed by the runtime, a real Git
// mutation underneath, and the one typed result driving both renderers.
func TestExitCodeContractThroughRealGitTransactions(t *testing.T) {
	tests := []struct {
		name            string
		wantExitCode    int
		wantFindingCode string
		setUp           func(t *testing.T, repositoryRoot string)
		options         func(repositoryRoot string) gittransaction.TransactionOptions
		mutate          func(t *testing.T, repositoryRoot string) func(*gittransaction.MutationRecorder) error
	}{
		{
			name:         "clean mutation succeeds",
			wantExitCode: 0,
			setUp: func(t *testing.T, repositoryRoot string) {
				writeFixtureFile(t, repositoryRoot, "target.txt", "initial\n")
				commitFixture(t, repositoryRoot, "initial")
			},
			options: func(repositoryRoot string) gittransaction.TransactionOptions {
				return gittransaction.TransactionOptions{RepositoryRoot: repositoryRoot, TargetPaths: []string{"target.txt"}}
			},
			mutate: func(t *testing.T, repositoryRoot string) func(*gittransaction.MutationRecorder) error {
				return func(recorder *gittransaction.MutationRecorder) error {
					if err := recorder.RecordTouched("target.txt"); err != nil {
						return err
					}
					writeFixtureFile(t, repositoryRoot, "target.txt", "tool change\n")
					return nil
				}
			},
		},
		{
			name:            "dirty target is safely refused",
			wantExitCode:    1,
			wantFindingCode: "GIT-DIRTY-TARGET",
			setUp: func(t *testing.T, repositoryRoot string) {
				writeFixtureFile(t, repositoryRoot, "target.txt", "initial\n")
				commitFixture(t, repositoryRoot, "initial")
				writeFixtureFile(t, repositoryRoot, "target.txt", "the user's own edit\n")
			},
			options: func(repositoryRoot string) gittransaction.TransactionOptions {
				return gittransaction.TransactionOptions{RepositoryRoot: repositoryRoot, TargetPaths: []string{"target.txt"}}
			},
			mutate: func(*testing.T, string) func(*gittransaction.MutationRecorder) error {
				return func(*gittransaction.MutationRecorder) error { return nil }
			},
		},
		{
			name:            "mutation failure rolls back cleanly",
			wantExitCode:    3,
			wantFindingCode: "GIT-MUTATION-FAILED",
			setUp: func(t *testing.T, repositoryRoot string) {
				writeFixtureFile(t, repositoryRoot, "target.txt", "initial\n")
				commitFixture(t, repositoryRoot, "initial")
			},
			options: func(repositoryRoot string) gittransaction.TransactionOptions {
				return gittransaction.TransactionOptions{RepositoryRoot: repositoryRoot, TargetPaths: []string{"target.txt"}}
			},
			mutate: func(t *testing.T, repositoryRoot string) func(*gittransaction.MutationRecorder) error {
				return func(recorder *gittransaction.MutationRecorder) error {
					if err := recorder.RecordTouched("target.txt"); err != nil {
						return err
					}
					writeFixtureFile(t, repositoryRoot, "target.txt", "half-written\n")
					return errors.New("forced mutation failure")
				}
			},
		},
		{
			name:            "post-commit verification failure reports committed-state risk",
			wantExitCode:    4,
			wantFindingCode: "GIT-COMMITTED-STATE-RISK",
			setUp: func(t *testing.T, repositoryRoot string) {
				writeFixtureFile(t, repositoryRoot, "target.txt", "initial\n")
				commitFixture(t, repositoryRoot, "initial")
			},
			options: func(repositoryRoot string) gittransaction.TransactionOptions {
				return gittransaction.TransactionOptions{
					RepositoryRoot: repositoryRoot,
					TargetPaths:    []string{"target.txt"},
					Commit:         true,
					CommitMessage:  "exact target change",
					PostCommitVerify: func(context.Context, string) error {
						return errors.New("forced post-commit verification failure")
					},
				}
			},
			mutate: func(t *testing.T, repositoryRoot string) func(*gittransaction.MutationRecorder) error {
				return func(recorder *gittransaction.MutationRecorder) error {
					if err := recorder.RecordTouched("target.txt"); err != nil {
						return err
					}
					writeFixtureFile(t, repositoryRoot, "target.txt", "committed\n")
					return nil
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := newFixtureRepository(t)
			test.setUp(t, repositoryRoot)
			// One transaction, rendered twice. Running the scenario a second time in a second
			// repository would give the two renderings different commit SHAs to talk about.
			transactionResult := gittransaction.BuildCommandResult("apply", gittransaction.ExecuteTransaction(
				context.Background(), test.options(repositoryRoot), test.mutate(t, repositoryRoot)))

			jsonExit, jsonOutput := runFixtureCommand(t, repositoryRoot, "json", transactionResult)
			if jsonExit != test.wantExitCode {
				t.Fatalf("JSON exit = %d, want %d\n%s", jsonExit, test.wantExitCode, jsonOutput)
			}

			var envelope map[string]any
			if err := json.Unmarshal([]byte(jsonOutput), &envelope); err != nil {
				t.Fatalf("decode JSON: %v\n%s", err, jsonOutput)
			}
			if envelope["schema_version"] != float64(resultmodel.SchemaVersion) {
				t.Fatalf("schema_version = %#v", envelope["schema_version"])
			}
			for _, field := range []string{"findings", "changes", "skipped_work", "rollback", "outcome", "repository_root"} {
				if envelope[field] == nil {
					t.Fatalf("%s is null in %s", field, jsonOutput)
				}
			}
			var decoded resultmodel.CommandResult
			if err := json.Unmarshal([]byte(jsonOutput), &decoded); err != nil {
				t.Fatalf("decode typed result: %v", err)
			}
			if resultmodel.ExitCode(decoded.Outcome) != test.wantExitCode {
				t.Fatalf("outcome %q does not map to exit %d", decoded.Outcome, test.wantExitCode)
			}

			textExit, textOutput := runFixtureCommand(t, repositoryRoot, "text", transactionResult)
			if textExit != test.wantExitCode {
				t.Fatalf("text exit = %d, want %d\n%s", textExit, test.wantExitCode, textOutput)
			}

			if test.wantFindingCode == "" {
				if len(decoded.Findings) != 0 {
					t.Fatalf("clean run reported findings: %#v", decoded.Findings)
				}
				return
			}
			if len(decoded.Findings) != 1 || decoded.Findings[0].Code != test.wantFindingCode {
				t.Fatalf("findings = %#v, want one %s", decoded.Findings, test.wantFindingCode)
			}
			// The two renderings come from one result, so the text form must name the same
			// finding code and the same next step the JSON form does.
			if !strings.Contains(textOutput, test.wantFindingCode) {
				t.Fatalf("text output does not name %s:\n%s", test.wantFindingCode, textOutput)
			}
			// Build the expected "next" line with the production renderer rather than joining
			// argv by hand, so this asserts parity instead of re-implementing the quoting.
			nextLine := renderedNextLine(t, decoded.Findings[0])
			if !strings.Contains(textOutput, nextLine) {
				t.Fatalf("text output does not name the next step %q:\n%s", nextLine, textOutput)
			}
			if len(decoded.Findings[0].VerificationArgv) == 0 {
				t.Fatalf("finding names no verification command: %#v", decoded.Findings[0])
			}
		})
	}
}

func renderedNextLine(t *testing.T, finding resultmodel.CommandFinding) string {
	t.Helper()
	rendered, err := resultmodel.RenderResult(resultmodel.CommandResult{
		Findings: []resultmodel.CommandFinding{finding},
	}, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(rendered), "\n") {
		if strings.HasPrefix(line, "  next: ") {
			return line
		}
	}
	t.Fatalf("the renderer emitted no next step for %#v", finding)
	return ""
}

func runFixtureCommand(t *testing.T, repositoryRoot, outputFormat string, result resultmodel.CommandResult) (int, string) {
	t.Helper()
	handlers := map[string]CommandHandler{
		"apply": func(ExecutionContext, []string) resultmodel.CommandResult { return result },
	}
	var stdout bytes.Buffer
	exitCode := NewRuntime(&stdout, handlers).Run([]string{"--repo-root", repositoryRoot, "--format", outputFormat, "apply"})
	return exitCode, stdout.String()
}

func newFixtureRepository(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
	repositoryRoot := t.TempDir()
	runFixtureGit(t, repositoryRoot, "init", "-q")
	runFixtureGit(t, repositoryRoot, "config", "user.name", "Do Work Test")
	runFixtureGit(t, repositoryRoot, "config", "user.email", "do-work@example.invalid")
	return repositoryRoot
}

func commitFixture(t *testing.T, repositoryRoot, message string) {
	t.Helper()
	runFixtureGit(t, repositoryRoot, "add", "-A")
	runFixtureGit(t, repositoryRoot, "commit", "-q", "-m", message)
}

func runFixtureGit(t *testing.T, repositoryRoot string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func writeFixtureFile(t *testing.T, repositoryRoot, relativePath, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repositoryRoot, relativePath), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A usage finding used to emit `do-work-cli --format text <command>`, which shows the shape
// but cannot be pasted. Every usage error now names a command that actually runs — the one
// the argv was reaching for when there is one, and a registered command otherwise.
func TestUsageFindingsNameARunnableCommand(t *testing.T) {
	registered := map[string]CommandHandler{
		"install-suite": func(ExecutionContext, []string) resultmodel.CommandResult {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
		},
		"validate-manifest": func(ExecutionContext, []string) resultmodel.CommandResult {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
		},
	}
	tests := []struct {
		name                string
		arguments           []string
		expectedCode        string
		expectedNextCommand string
	}{
		{
			name:                "an unknown global option before a known command",
			arguments:           []string{"--nonsense", "validate-manifest"},
			expectedCode:        "UNKNOWN-GLOBAL-OPTION",
			expectedNextCommand: "validate-manifest",
		},
		{
			name:                "a global option missing its value",
			arguments:           []string{"--repo-root"},
			expectedCode:        "MISSING-OPTION-VALUE",
			expectedNextCommand: "help",
		},
		{
			name:                "an invalid format before a known command",
			arguments:           []string{"--format", "yaml", "install-suite"},
			expectedCode:        "INVALID-OUTPUT-FORMAT",
			expectedNextCommand: "install-suite",
		},
		{
			name:                "no command at all",
			arguments:           []string{"--format", "text"},
			expectedCode:        "MISSING-COMMAND",
			expectedNextCommand: "help",
		},
		{
			name:                "an unknown command points to help",
			arguments:           []string{"not-a-command"},
			expectedCode:        "UNKNOWN-COMMAND",
			expectedNextCommand: "help",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := NewRuntime(&stdout, registered).Run(append([]string{"--format", "json"}, test.arguments...))
			if exitCode != resultmodel.ExitCode(resultmodel.OutcomeFailure) {
				t.Fatalf("exit = %d, want %d\n%s", exitCode, resultmodel.ExitCode(resultmodel.OutcomeFailure), stdout.String())
			}
			var decoded resultmodel.CommandResult
			if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
				// A parse error before --format json took effect renders as text; read that instead.
				if !strings.Contains(stdout.String(), test.expectedCode) {
					t.Fatalf("output does not name %s:\n%s", test.expectedCode, stdout.String())
				}
				if !strings.Contains(stdout.String(), test.expectedNextCommand) {
					t.Fatalf("output does not name the runnable command %q:\n%s", test.expectedNextCommand, stdout.String())
				}
				return
			}
			if len(decoded.Findings) != 1 || decoded.Findings[0].Code != test.expectedCode {
				t.Fatalf("findings = %#v, want one %s", decoded.Findings, test.expectedCode)
			}
			finding := decoded.Findings[0]
			for _, argv := range [][]string{finding.NextArgv, finding.VerificationArgv} {
				if len(argv) == 0 {
					t.Fatalf("finding names no argv: %#v", finding)
				}
				if argv[len(argv)-1] != test.expectedNextCommand {
					t.Errorf("argv %v does not end in the runnable command %q", argv, test.expectedNextCommand)
				}
				for _, argument := range argv {
					if strings.HasPrefix(argument, "<") && strings.HasSuffix(argument, ">") {
						t.Errorf("argv %v still carries the placeholder %q", argv, argument)
					}
				}
			}
		})
	}
}

func TestHelpSucceedsAndListsAvailableCommands(t *testing.T) {
	registered := map[string]CommandHandler{
		"install-suite": func(ExecutionContext, []string) resultmodel.CommandResult { return resultmodel.CommandResult{} },
		"next":          func(ExecutionContext, []string) resultmodel.CommandResult { return resultmodel.CommandResult{} },
	}
	for _, format := range []string{"text", "json"} {
		var stdout bytes.Buffer
		exitCode := NewRuntime(&stdout, registered).Run([]string{"--format", format, "help"})
		if exitCode != 0 {
			t.Fatalf("%s help exit = %d\n%s", format, exitCode, stdout.String())
		}
		for _, command := range []string{"help", "install-suite", "next"} {
			if !strings.Contains(stdout.String(), command) {
				t.Errorf("%s help omitted %q:\n%s", format, command, stdout.String())
			}
		}
	}
}

func TestStatusStaysUnknownAndNeverNamesItselfAsRemedy(t *testing.T) {
	var stdout bytes.Buffer
	exitCode := NewRuntime(&stdout, map[string]CommandHandler{
		"next": func(ExecutionContext, []string) resultmodel.CommandResult { return resultmodel.CommandResult{} },
	}).Run([]string{"--format", "json", "status"})
	if exitCode == 0 {
		t.Fatalf("status unexpectedly succeeded:\n%s", stdout.String())
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Findings) != 1 || decoded.Findings[0].Code != "UNKNOWN-COMMAND" {
		t.Fatalf("status finding = %#v", decoded.Findings)
	}
	for _, argv := range [][]string{decoded.Findings[0].NextArgv, decoded.Findings[0].VerificationArgv} {
		if len(argv) == 0 || argv[len(argv)-1] != "help" || strings.Contains(strings.Join(argv, " "), "status") {
			t.Errorf("status remedy = %v, want help without status", argv)
		}
	}
}

// A finding for an unavailable command has to tell the reader what IS available, otherwise
// the only next step is guessing.
func TestUnknownAndMissingCommandFindingsListTheRegisteredCommands(t *testing.T) {
	registered := map[string]CommandHandler{
		"install-suite":     func(ExecutionContext, []string) resultmodel.CommandResult { return resultmodel.CommandResult{} },
		"validate-manifest": func(ExecutionContext, []string) resultmodel.CommandResult { return resultmodel.CommandResult{} },
	}
	for _, arguments := range [][]string{{"not-a-command"}, {"--format", "text"}} {
		var stdout bytes.Buffer
		NewRuntime(&stdout, registered).Run(arguments)
		for _, expectedCommand := range []string{"install-suite", "validate-manifest"} {
			if !strings.Contains(stdout.String(), expectedCommand) {
				t.Errorf("argv %v produced a finding that does not list %q:\n%s", arguments, expectedCommand, stdout.String())
			}
		}
	}
}

func TestRuntimePreservesDeclaredCompatibilityStatus(t *testing.T) {
	for _, testCase := range []struct{ override, want int }{{130, 130}, {99, 99}, {0, 0}} {
		var output bytes.Buffer
		runtime := NewRuntime(&output, map[string]CommandHandler{"media": func(ExecutionContext, []string) resultmodel.CommandResult {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExitCodeOverride: testCase.override}
		}})
		if got := runtime.Run([]string{"media"}); got != testCase.want {
			t.Fatalf("override %d: got %d want %d", testCase.override, got, testCase.want)
		}
	}
}
