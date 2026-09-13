package toolboxcommands

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestReportImageBatchExactTextRetainsFailureDiagnostics(t *testing.T) {
	for _, partial := range []bool{false, true} {
		name := "all failed"
		if partial {
			name = "partial failure"
		}
		t.Run(name, func(t *testing.T) {
			if partial && runtime.GOOS == "windows" {
				t.Skip("partial-success fixture requires a Unix shell and owned process groups")
			}
			repository := toolboxTestRepository(t)
			report := filepath.Join(repository, "report")
			if err := os.Mkdir(report, 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND", "0")
			originalLookup := reportImageLookPath
			reportImageLookPath = func(string) (string, error) { return "", exec.ErrNotFound }
			t.Cleanup(func() { reportImageLookPath = originalLookup })
			if partial {
				backend := filepath.Join(repository, "imagegen")
				script := "#!/bin/sh\nout=\nprompt=\nwhile [ $# -gt 0 ]; do\ncase \"$1\" in\n--output) out=$2; shift 2;;\n--prompt) prompt=$2; shift 2;;\n*) shift;;\nesac\ndone\ncase \"$prompt\" in\n*success*) printf fresh > \"$out\";;\n*) exit 1;;\nesac\n"
				if err := os.WriteFile(backend, []byte(script), 0o755); err != nil {
					t.Fatal(err)
				}
				reportImageLookPath = func(string) (string, error) { return backend, nil }
			}
			stderr, err := os.CreateTemp(t.TempDir(), "stderr")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { stderr.Close() })
			originalStderr := os.Stderr
			os.Stderr = stderr
			t.Cleanup(func() { os.Stderr = originalStderr })
			var stdout bytes.Buffer
			command := commandruntime.NewRuntime(&stdout, map[string]commandruntime.CommandHandler{CommandReportImageBatch: handleReportImageBatch})
			status := command.Run([]string{"--repo-root", repository, "--format", "text", CommandReportImageBatch, report, "style", "one.png:success", "two.png:failure"})
			if status != 0 {
				t.Fatalf("batch status = %d, stdout = %q", status, stdout.String())
			}
			wantOutput := ""
			if partial {
				wantOutput = filepath.Join(report, "generated") + "\n"
			}
			if stdout.String() != wantOutput {
				t.Fatalf("stdout = %q, want exact %q", stdout.String(), wantOutput)
			}
			diagnostics, err := os.ReadFile(stderr.Name())
			if err != nil {
				t.Fatal(err)
			}
			for _, expected := range []string{"REPORT-IMAGE-MISSING", "two.png", "backend failed or produced no nonempty image"} {
				if !strings.Contains(string(diagnostics), expected) {
					t.Errorf("stderr %q omits %q", diagnostics, expected)
				}
			}
			if !partial && !strings.Contains(string(diagnostics), "one.png") {
				t.Errorf("all-failed stderr omits first image: %q", diagnostics)
			}
		})
	}
}
