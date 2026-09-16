package corehelpers

import (
	"bytes"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestFrontmatterUsesSchemaNormalizationAndDecodedValues(t *testing.T) {
	cases := []struct {
		name, contents, field string
		normalize             bool
		want                  string
		status                int
		warning               bool
	}{
		{"canonical route", "route: A", "route", true, "A\n", 0, false},
		{"lowercase route", "route: b", "route", true, "B\n", 0, false},
		{"status alias", "status: done", "status", true, "completed\n", 0, false},
		{"plain alias", "status: done", "status", false, "done\n", 0, false},
		{"assignment case", "assigned_to: Cloud-Alpha_2", "assigned_to", true, "Cloud-Alpha_2\n", 0, false},
		{"path case", "custom_path: Src/Main.go", "custom_path", true, "Src/Main.go\n", 0, false},
		{"schema fallback", "domain: unknown", "domain", true, "general\n", 0, true},
		{"unknown without default", "route: unknown", "route", true, "UNKNOWN\n", 0, true},
		{"block list", "depends_on:\n  - REQ-001\n  - REQ-002", "depends_on", false, "REQ-001\nREQ-002\n", 0, false},
		{"flow list", "depends_on: ['REQ-001', REQ-002]", "depends_on", false, "REQ-001\nREQ-002\n", 0, false},
		{"empty list", "depends_on: []", "depends_on", false, "", 0, false},
		{"empty scalar", "status:", "status", false, "", 1, true},
		{"absent", "route: A", "status", false, "", 1, true},
		{"escaped quote", "title: 'It''s a test'", "title", false, "It's a test\n", 0, false},
		{"inline comment", "status: completed # confirmed", "status", false, "completed\n", 0, false},
		{"quoted hash", "title: 'Keep # literal'", "title", false, "Keep # literal\n", 0, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeMatrixFile(t, root, "request.md", "\ufeff---\r\n"+strings.ReplaceAll(test.contents, "\n", "\r\n")+"\r\n---\r\nBody\r\n")
			args := []string{"get", "request.md", test.field}
			if test.normalize {
				args = append(args, "--normalize")
			}
			result := handleFrontmatter(testContext(root), args)
			output, err := resultmodel.RenderResult(result, resultmodel.FormatText)
			if err != nil || string(output) != test.want || resultmodel.ExitCode(result.Outcome) != test.status {
				t.Fatalf("got output=%q status=%d error=%v; want %q status=%d", output, resultmodel.ExitCode(result.Outcome), err, test.want, test.status)
			}
			if (len(result.Findings) > 0) != test.warning {
				t.Fatalf("warning mismatch: %+v", result.Findings)
			}
		})
	}
}

func TestFrontmatterMembershipUsesExitStatusAndEmptyStdout(t *testing.T) {
	for _, test := range []struct {
		value, set string
		status     int
	}{
		{"pending", "terminal-success", 1},
		{"completed", "terminal-success", 0},
		{"completed-with-issues", "terminal-success", 0},
		{"done", "terminal-success", 0},
		{"cancelled", "terminal-success", 1},
		{"abandoned", "terminal-resolved", 0},
		{"completed-with-issues", "terminal-resolved", 0},
		{"failed", "terminal-resolved", 1},
		{"pending", "terminal-resolved", 1},
	} {
		t.Run(test.value+"/"+test.set, func(t *testing.T) {
			root := t.TempDir()
			writeMatrixFile(t, root, "request.md", "---\nstatus: "+test.value+"\n---\n")
			var output bytes.Buffer
			runtime := commandruntime.NewRuntime(&output, Handlers())
			status := runtime.Run([]string{"--repo-root", root, "frontmatter", "get", "request.md", "status", "--in-set", test.set})
			if status != test.status || output.Len() != 0 {
				t.Fatalf("predicate returned status=%d stdout=%q; want status=%d empty stdout", status, output.String(), test.status)
			}
		})
	}
}

func TestFrontmatterRejectsInvalidMembershipArguments(t *testing.T) {
	root := t.TempDir()
	writeMatrixFile(t, root, "request.md", "---\nstatus: completed\nroute: A\n---\n")
	for _, args := range [][]string{
		{"get", "request.md", "route", "--in-set", "terminal-success"},
		{"get", "request.md", "status", "--in-set", ""},
		{"get", "request.md", "status", "--in-set"},
		{"get", "request.md", "status", "--in-set", "unknown"},
	} {
		result := handleFrontmatter(testContext(root), args)
		if resultmodel.ExitCode(result.Outcome) != 2 {
			t.Fatalf("invalid membership arguments accepted: %v", args)
		}
	}
}
