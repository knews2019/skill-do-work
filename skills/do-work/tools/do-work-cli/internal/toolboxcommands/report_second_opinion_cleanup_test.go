package toolboxcommands

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestSecondOpinionWorktreeReadOnlyCleanup(t *testing.T) {
	for _, scenario := range []string{"success", "backend failure", "cancellation", "outside link"} {
		t.Run(scenario, func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			scratch, external := t.TempDir(), t.TempDir()
			t.Setenv("TMPDIR", scratch)
			folder := filepath.Join(bundle, "readonly-assets")
			if err := os.Mkdir(folder, 0755); err != nil {
				t.Fatal(err)
			}
			secondOpinionWriteTestFile(t, filepath.Join(folder, "chart.svg"), "original asset")
			secondOpinionWriteTestFile(t, filepath.Join(external, "keep.txt"), "external asset")
			for _, path := range []string{folder, external} {
				if err := os.Chmod(path, 0555); err != nil {
					t.Fatal(err)
				}
			}
			t.Cleanup(func() {
				// Keep a failing regression from leaking its fixture's read-only copy.
				_ = os.Chmod(folder, 0755)
				_ = os.Chmod(external, 0755)
				_ = filepath.WalkDir(scratch, func(path string, entry os.DirEntry, err error) error {
					if err == nil && entry.IsDir() {
						_ = os.Chmod(path, 0755)
					}
					return nil
				})
				_ = os.RemoveAll(scratch)
				toolboxTestGit(t, repo, "worktree", "prune")
			})
			previous := secondOpinionJudge
			secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) { return "pass", nil }
			t.Cleanup(func() { secondOpinionJudge = previous })
			extra := ""
			switch scenario {
			case "backend failure":
				extra = "exit 7"
			case "outside link":
				t.Setenv("SECOND_OPINION_EXTERNAL", external)
				extra = `/bin/ln -s "$SECOND_OPINION_EXTERNAL" ../../outside-link`
			case "cancellation":
				t.Setenv("SECOND_OPINION_READY", filepath.Join(t.TempDir(), "ready"))
				extra = `printf ready > "$SECOND_OPINION_READY"; /bin/sleep 30`
			}
			secondOpinionFakeCodex(t, extra)
			var row resultmodel.AIReportSecondOpinionRun
			if scenario == "cancellation" {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				template, err := os.ReadFile(os.Getenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT"))
				if err != nil {
					t.Fatal(err)
				}
				row = resultmodel.AIReportSecondOpinionRun{BackendName: "codex", RunState: "failed"}
				done := make(chan struct{})
				go func() {
					defer close(done)
					secondOpinionRun(ctx, repo, bundle, filepath.Join(os.Getenv("PATH"), "codex"), string(template), secondOpinionDirection{}, &row, &resultmodel.CommandResult{})
				}()
				secondOpinionAwaitFile(t, os.Getenv("SECOND_OPINION_READY"))
				cancel()
				select {
				case <-done:
				case <-time.After(10 * time.Second):
					t.Fatal("cancelled run did not finish cleanup")
				}
			} else {
				row = secondOpinionInvoke(t, repo, "--backend", "codex").AIReportSecondOpinion.BackendRuns[0]
			}
			if strings.Contains(row.FailureReason, "worktree cleanup failed") {
				t.Errorf("read-only directory prevented cleanup: %s", row.FailureReason)
			}
			if scenario == "success" {
				if row.RunState != "ran" || row.BannerState != "inserted" {
					t.Errorf("valid report refused: %+v", row)
				}
			} else if row.RunState != "failed" || !strings.Contains(row.FailureReason, "staged page kept at") {
				t.Errorf("failure lost recovery page: %+v", row)
			}
			for _, path := range []string{folder, external} {
				info, err := os.Stat(path)
				if err != nil || info.Mode().Perm() != 0555 {
					t.Errorf("original permissions changed: %s: %v %v", path, info, err)
				}
			}
			for path, expected := range map[string]string{filepath.Join(folder, "chart.svg"): "original asset", filepath.Join(external, "keep.txt"): "external asset"} {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != expected {
					t.Errorf("original bytes changed: %s: %q %v", path, data, err)
				}
			}
			secondOpinionAssertNoWorktree(t, repo, scratch)
		})
	}
}
