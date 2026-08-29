package commandruntime

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

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
