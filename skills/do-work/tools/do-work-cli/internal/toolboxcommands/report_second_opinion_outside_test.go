package toolboxcommands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestSecondOpinionOutsideSnapshotBaseline(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionWriteTestFile(t, filepath.Join(repo, ".gitignore"), "ignored.txt\n")
	secondOpinionWriteTestFile(t, filepath.Join(repo, "ignored.txt"), "baseline")
	external := t.TempDir()
	secondOpinionWriteTestFile(t, filepath.Join(external, "asset"), "external")
	if err := os.Symlink(external, filepath.Join(repo, "linked")); err != nil {
		t.Fatal(err)
	}
	before, err := secondOpinionOutsideSnapshot(context.Background(), repo, bundle)
	if err != nil {
		t.Fatal(err)
	}
	// A real ignored, untracked baseline is allowed. Links remain opaque, even
	// when their targets change, and bundle publication is checked elsewhere.
	secondOpinionWriteTestFile(t, filepath.Join(external, "asset"), "changed externally")
	secondOpinionWriteTestFile(t, filepath.Join(bundle, "new.html"), "allowed by another check")
	row := resultmodel.AIReportSecondOpinionRun{}
	if err := secondOpinionOutsideCheck(context.Background(), repo, bundle, before, &row); err != nil {
		t.Fatal(err)
	}
	secondOpinionWriteTestFile(t, filepath.Join(repo, "ignored.txt"), "changed")
	if err := secondOpinionOutsideCheck(context.Background(), repo, bundle, before, &row); err == nil || !slices.Equal(row.OutsidePaths, []string{"ignored.txt"}) {
		t.Fatalf("ignored baseline edit missed: %v %v", row.OutsidePaths, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := secondOpinionOutsideSnapshot(ctx, repo, bundle); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled snapshot continued: %v", err)
	}
}

func TestSecondOpinionWorktreeOutsideSnapshot(t *testing.T) {
	for _, phase := range []string{"probe", "generation"} {
		for _, scenario := range []struct{ name, script, path string }{
			{"committed edit", "printf changed > ../../README.md\ngit -C ../.. add README.md\ngit -C ../.. commit -qm changed", "README.md"},
			{"ignored new file", "printf changed > ../../scratch-output.txt", "scratch-output.txt"},
			{"staged edit", "printf changed > ../../README.md\ngit -C ../.. add README.md", "README.md"},
			{"unstaged edit", "printf changed > ../../README.md", "README.md"},
			{"committed deletion", "git -C ../.. rm -q README.md\ngit -C ../.. commit -qm deleted", "README.md"},
			{"changed ignore rules", "printf hidden.txt >> ../../.gitignore\nprintf hidden > ../../hidden.txt\ngit -C ../.. add .gitignore\ngit -C ../.. commit -qm ignored", "hidden.txt"},
			{"unchanged ignored input", "test \"$(/bin/cat ../../scratch-input.txt)\" = baseline", ""},
		} {
			t.Run(phase+"/"+scenario.name, func(t *testing.T) {
				repo, bundle := secondOpinionFixture(t)
				secondOpinionWriteTestFile(t, filepath.Join(repo, "README.md"), "original")
				secondOpinionWriteTestFile(t, filepath.Join(repo, ".gitignore"), "scratch-*.txt\n")
				secondOpinionWriteTestFile(t, filepath.Join(repo, "scratch-input.txt"), "baseline")
				toolboxTestGit(t, repo, "add", "README.md", ".gitignore")
				toolboxTestGit(t, repo, "add", "-f", "scratch-input.txt")
				toolboxTestGit(t, repo, "commit", "-qm", "outside baseline")
				// Real-checkout dirt is not the isolated checkout's baseline.
				secondOpinionWriteTestFile(t, filepath.Join(repo, "README.md"), "user edit")
				scratch := t.TempDir()
				t.Setenv("TMPDIR", scratch)
				extra := scenario.script
				if phase == "probe" {
					extra = "exit 93" // A refused probe must never reach generation.
					if scenario.path == "" {
						extra = ""
					}
				}
				secondOpinionFakeCodex(t, extra)
				if phase == "probe" {
					path := filepath.Join(os.Getenv("PATH"), "codex")
					script, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					// The fake backend handles its probe before writing the page.
					script = []byte(strings.Replace(string(script), "for prompt do", scenario.script+"\nfor prompt do", 1))
					if err := os.WriteFile(path, script, 0755); err != nil {
						t.Fatal(err)
					}
				}
				previous := secondOpinionJudge
				secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) { return "pass", nil }
				t.Cleanup(func() { secondOpinionJudge = previous })
				row := secondOpinionInvoke(t, repo, "--backend", "codex").AIReportSecondOpinion.BackendRuns[0]
				if scenario.path == "" {
					if row.RunState != "ran" || row.BannerState != "inserted" || len(row.OutsidePaths) != 0 {
						t.Fatalf("unchanged input refused: %+v", row)
					}
				} else {
					if row.RunState != "failed" || !slices.Contains(row.OutsidePaths, scenario.path) || !strings.Contains(row.FailureReason, "outside-bundle changes") {
						t.Errorf("outside change not detected: %+v", row)
					}
					kept, _ := filepath.Glob(filepath.Join(bundle, ".*-staging-*.html"))
					if phase == "generation" && (len(kept) != 1 || !strings.Contains(row.FailureReason, "staged page kept at")) {
						t.Errorf("finished refusal page lost: %+v", row)
					}
					if phase == "probe" && len(kept) != 0 {
						t.Errorf("generation ran after refused probe: %v", kept)
					}
				}
				data, err := os.ReadFile(filepath.Join(repo, "README.md"))
				if err != nil || string(data) != "user edit" {
					t.Errorf("real outside file changed: %q %v", data, err)
				}
				if _, err := os.Stat(filepath.Join(repo, "scratch-output.txt")); !os.IsNotExist(err) {
					t.Errorf("ignored output leaked to real checkout: %v", err)
				}
				secondOpinionAssertNoWorktree(t, repo, scratch)
			})
		}
	}
}
