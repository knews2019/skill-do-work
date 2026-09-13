package commandruntime

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestExactTextWarningsAndErrorsRemainVisibleOnStderr(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			stderr, err := os.CreateTemp(t.TempDir(), "stderr")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { stderr.Close() })
			originalStderr := os.Stderr
			os.Stderr = stderr
			t.Cleanup(func() { os.Stderr = originalStderr })
			var stdout bytes.Buffer
			exact := "compatibility output\n"
			runtime := NewRuntime(&stdout, map[string]CommandHandler{"exact": func(ExecutionContext, []string) resultmodel.CommandResult {
				return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &exact, Findings: []resultmodel.CommandFinding{
					{Code: "EXACT-WARNING", Severity: resultmodel.SeverityWarning, Evidence: []string{"warning explanation"}},
					{Code: "EXACT-ERROR", Severity: resultmodel.SeverityError, Evidence: []string{"error explanation"}},
					{Code: "EXACT-INFO", Severity: resultmodel.SeverityInfo, Evidence: []string{"informational evidence"}},
				}}
			}})
			if status := runtime.Run([]string{"--format", format, "exact"}); status != 0 {
				t.Fatalf("status = %d, want declared success", status)
			}
			diagnostics, err := os.ReadFile(stderr.Name())
			if err != nil {
				t.Fatal(err)
			}
			if format == "json" {
				if len(diagnostics) != 0 || !strings.Contains(stdout.String(), "EXACT-ERROR") {
					t.Fatalf("JSON findings must remain in stdout: stdout=%q stderr=%q", stdout.String(), diagnostics)
				}
				return
			}
			if stdout.String() != exact {
				t.Fatalf("stdout = %q, want %q", stdout.String(), exact)
			}
			for _, expected := range []string{"EXACT-WARNING", "warning explanation", "EXACT-ERROR", "error explanation"} {
				if !strings.Contains(string(diagnostics), expected) {
					t.Errorf("stderr %q omits %q", diagnostics, expected)
				}
			}
			if strings.Contains(string(diagnostics), "EXACT-INFO") {
				t.Fatalf("informational finding leaked onto diagnostic stderr: %q", diagnostics)
			}
		})
	}
}
