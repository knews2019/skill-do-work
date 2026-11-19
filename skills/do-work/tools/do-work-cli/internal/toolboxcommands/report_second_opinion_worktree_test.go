package toolboxcommands

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func secondOpinionAssertNoWorktree(t *testing.T, repository, scratch string) {
	t.Helper()
	output, err := exec.Command("git", "-C", repository, "worktree", "list", "--porcelain").CombinedOutput()
	if err != nil || strings.Count(string(output), "worktree ") != 1 {
		t.Fatalf("worktree residue: %s %v", output, err)
	}
	entries, err := os.ReadDir(scratch)
	if err != nil || len(entries) != 0 {
		t.Fatalf("temporary worktree residue: %v %v", entries, err)
	}
}

func TestSecondOpinionWorktreeDirtyReadme(t *testing.T) {
	for _, outsideWrite := range []bool{false, true} {
		t.Run(map[bool]string{false: "unrelated dirt", true: "backend outside write"}[outsideWrite], func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			readme := filepath.Join(repo, "README.md")
			if err := os.WriteFile(readme, []byte("committed\n"), 0644); err != nil {
				t.Fatal(err)
			}
			toolboxTestGit(t, repo, "add", "README.md")
			toolboxTestGit(t, repo, "commit", "-qm", "tracked README")
			dirty := "uncommitted user edit\n"
			if err := os.WriteFile(readme, []byte(dirty), 0644); err != nil {
				t.Fatal(err)
			}
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			extra := "git symbolic-ref -q HEAD && exit 91"
			if outsideWrite {
				extra += "\nprintf backend > ../../README.md"
			}
			secondOpinionFakeCodex(t, extra)
			previous := secondOpinionJudge
			secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) { return "pass", nil }
			t.Cleanup(func() { secondOpinionJudge = previous })
			result := secondOpinionInvoke(t, repo, "--backend", "codex")
			row := secondOpinionRows(t, result)[0]
			if outsideWrite {
				if row["state"] != "failed" || !strings.Contains(strings.Join(result.AIReportSecondOpinion.BackendRuns[0].OutsidePaths, ","), "README.md") {
					t.Fatalf("outside write accepted: %v", row)
				}
				kept, _ := filepath.Glob(filepath.Join(bundle, ".*-staging-*.html"))
				if len(kept) != 1 || !strings.Contains(row["reason"].(string), filepath.Base(kept[0])) {
					t.Fatalf("finished refusal page not retained: %v", row)
				}
			} else if result.Outcome != resultmodel.OutcomeSuccess || row["banner_state"] != "inserted" {
				t.Fatalf("dirty README blocked publication: %v", row)
			}
			after, err := os.ReadFile(readme)
			if err != nil || string(after) != dirty {
				t.Fatalf("real README changed: %q %v", after, err)
			}
			secondOpinionAssertNoWorktree(t, repo, scratch)
		})
	}
}

func TestSecondOpinionWorktreeCurrentBundleAndRefusals(t *testing.T) {
	for _, scenario := range []string{"exact bundle", "backend failure", "original replacement", "unprefixed asset", "asset link", "probe write", "real original replacement", "real asset collision"} {
		t.Run(scenario, func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			secondOpinionWriteTestFile(t, filepath.Join(bundle, "deleted.txt"), "old")
			secondOpinionWriteTestFile(t, filepath.Join(bundle, "edited.txt"), "old")
			toolboxTestGit(t, repo, "add", ".")
			toolboxTestGit(t, repo, "commit", "-qm", "bundle baseline")
			if err := os.Remove(filepath.Join(bundle, "deleted.txt")); err != nil {
				t.Fatal(err)
			}
			secondOpinionWriteTestFile(t, filepath.Join(bundle, "edited.txt"), "current")
			secondOpinionWriteTestFile(t, filepath.Join(bundle, "added.txt"), "new")
			if err := os.Mkdir(filepath.Join(bundle, "empty"), 0750); err != nil {
				t.Fatal(err)
			}
			original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
			t.Setenv("REAL_BUNDLE", bundle)
			extra := `test ! -e deleted.txt && test "$(/bin/cat edited.txt)" = current && test "$(/bin/cat added.txt)" = new && test -d empty || exit 92
repository=${prompt#*'Repository: '}; repository=${repository%%'
'*}
bundle=${prompt#*'Bundle: '}; bundle=${bundle%%'
'*}
test "$repository" = "$(git rev-parse --show-toplevel)" && test "$bundle" = "$PWD" || exit 93
printf new-asset > codex-actual-model-1-chart.svg`
			switch scenario {
			case "backend failure":
				extra += "\nexit 7"
			case "original replacement":
				extra += "\n/bin/cp index.html original-copy; /bin/mv original-copy index.html"
			case "unprefixed asset":
				extra += "\nprintf forbidden > unrelated.txt"
			case "asset link":
				extra += "\n/bin/ln -s index.html codex-actual-model-1-link.svg"
			case "real original replacement":
				extra += "\n/bin/cp \"$REAL_BUNDLE/index.html\" \"$REAL_BUNDLE/copy\"; /bin/mv \"$REAL_BUNDLE/copy\" \"$REAL_BUNDLE/index.html\""
			case "real asset collision":
				extra += "\nprintf user > \"$REAL_BUNDLE/codex-actual-model-1-chart.svg\""
			}
			secondOpinionFakeCodex(t, extra)
			if scenario == "probe write" {
				script, _ := os.ReadFile(filepath.Join(os.Getenv("PATH"), "codex"))
				script = []byte(strings.Replace(string(script), "for prompt do", "printf changed > ../../README.md\nfor prompt do", 1))
				if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), "codex"), script, 0755); err != nil {
					t.Fatal(err)
				}
			}
			row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
			if scenario == "exact bundle" {
				if row["state"] != "ran" {
					t.Fatalf("current bundle not reproduced: %v", row)
				}
				asset, _ := os.ReadFile(filepath.Join(bundle, "codex-actual-model-1-chart.svg"))
				if string(asset) != "new-asset" {
					t.Fatalf("asset not imported: %q", asset)
				}
			} else {
				if row["state"] != "failed" {
					t.Fatalf("violation accepted: %v", row)
				}
				if scenario != "probe write" && !strings.HasPrefix(scenario, "real ") && !strings.Contains(row["reason"].(string), "staged page kept at") {
					t.Fatalf("finished page lost: %v", row)
				}
			}
			if scenario == "real asset collision" {
				data, _ := os.ReadFile(filepath.Join(bundle, "codex-actual-model-1-chart.svg"))
				if string(data) != "user" {
					t.Fatalf("collision overwritten: %q", data)
				}
			}
			after, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
			if string(after) != string(original) {
				t.Fatal("real original changed")
			}
			if _, err := os.Lstat(filepath.Join(bundle, "unrelated.txt")); !os.IsNotExist(err) {
				t.Fatal("unprefixed file imported")
			}
			if _, err := os.Lstat(filepath.Join(bundle, "codex-actual-model-1-link.svg")); !os.IsNotExist(err) {
				t.Fatal("link imported")
			}
			secondOpinionAssertNoWorktree(t, repo, scratch)
		})
	}
}

func secondOpinionWriteTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestSecondOpinionSessionIgnoresOutsideDirtThroughPublication(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionWriteTestFile(t, filepath.Join(repo, "README.md"), "dirty before prepare")
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--prepare", "--backend", "codex", "--model", "session-model"))[0]
	if row["state"] != "prepared" {
		t.Fatalf("prepare blocked: %v", row)
	}
	stage := row["staging_page"].(string)
	secondOpinionWriteTestFile(t, stage, `<html><body><p>Original report: <a href="index.html">index.html</a></p></body></html>`)
	secondOpinionWriteTestFile(t, row["staging_reply"].(string), "editorial reply")
	previous := secondOpinionJudge
	secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) {
		secondOpinionWriteTestFile(t, filepath.Join(repo, "README.md"), "changed during judge")
		return "pass", nil
	}
	t.Cleanup(func() { secondOpinionJudge = previous })
	result := secondOpinionInvoke(t, repo, "--finalize", stage)
	row = secondOpinionRows(t, result)[0]
	if result.Outcome != resultmodel.OutcomeSuccess || row["banner_state"] != "inserted" {
		t.Fatalf("publication tail blocked: %v", row)
	}
	page, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if !strings.Contains(string(page), "second-opinion-banner:") {
		t.Fatal("banner missing")
	}
}

func TestSecondOpinionWorktreeCancellationAndDeadline(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(map[bool]string{false: "cancellation", true: "deadline"}[timeout], func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			scratch, signals := t.TempDir(), t.TempDir()
			t.Setenv("TMPDIR", scratch)
			t.Setenv("SECOND_OPINION_READY", filepath.Join(signals, "ready"))
			t.Setenv("SECOND_OPINION_PID", filepath.Join(signals, "pid"))
			secondOpinionFakeCodex(t, `/bin/sleep 30 &
printf '%s' "$!" > "$SECOND_OPINION_PID"
printf ready > "$SECOND_OPINION_READY"
wait`)
			template, err := os.ReadFile(os.Getenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT"))
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			if timeout {
				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			}
			defer cancel()
			row := resultmodel.AIReportSecondOpinionRun{BackendName: "codex", RunState: "failed"}
			done := make(chan struct{})
			go func() {
				defer close(done)
				secondOpinionRun(ctx, repo, bundle, filepath.Join(os.Getenv("PATH"), "codex"), string(template), secondOpinionDirection{}, &row, &resultmodel.CommandResult{})
			}()
			secondOpinionAwaitFile(t, os.Getenv("SECOND_OPINION_READY"))
			if !timeout {
				cancel()
			}
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				t.Fatal("run did not finish cleanup")
			}
			expected := "context canceled"
			if timeout {
				expected = "context deadline exceeded"
			}
			if !strings.Contains(row.FailureReason, expected) || !strings.Contains(row.FailureReason, "staged page kept at") {
				t.Fatalf("cancellation lost failure/page: %+v", row)
			}
			secondOpinionAssertProcessGone(t, os.Getenv("SECOND_OPINION_PID"))
			secondOpinionAssertNoWorktree(t, repo, scratch)
		})
	}
}

func secondOpinionAwaitFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("backend did not signal readiness: %s", path)
}

func secondOpinionAssertProcessGone(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("/bin/kill", "-0", string(data)).Run(); err == nil {
		t.Fatalf("backend descendant %s survived cleanup", data)
	}
}

func TestSecondOpinionWorktreeSIGINT(t *testing.T) {
	if repo := os.Getenv("SECOND_OPINION_SIGNAL_REPO"); repo != "" {
		result := secondOpinionInvoke(t, repo, "--backend", "codex")
		if result.ExitCodeOverride != 130 {
			t.Fatalf("interrupt status=%d", result.ExitCodeOverride)
		}
		return
	}
	repo, _ := secondOpinionFixture(t)
	scratch, signals := t.TempDir(), t.TempDir()
	t.Setenv("TMPDIR", scratch)
	t.Setenv("SECOND_OPINION_READY", filepath.Join(signals, "ready"))
	t.Setenv("SECOND_OPINION_PID", filepath.Join(signals, "pid"))
	secondOpinionFakeCodex(t, `/bin/sleep 30 &
printf '%s' "$!" > "$SECOND_OPINION_PID"
printf ready > "$SECOND_OPINION_READY"
wait`)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, executable, "-test.run=^TestSecondOpinionWorktreeSIGINT$")
	command.Env = append(os.Environ(), "SECOND_OPINION_SIGNAL_REPO="+repo)
	var output bytes.Buffer
	command.Stdout, command.Stderr = &output, &output
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = command.Process.Kill() })
	secondOpinionAwaitFile(t, os.Getenv("SECOND_OPINION_READY"))
	if err := command.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	if err := command.Wait(); err != nil {
		t.Fatalf("SIGINT helper failed: %v\n%s", err, output.String())
	}
	secondOpinionAssertProcessGone(t, os.Getenv("SECOND_OPINION_PID"))
	secondOpinionAssertNoWorktree(t, repo, scratch)
}

func TestSecondOpinionWorktreeSetupFailureAndCleanupFailure(t *testing.T) {
	for _, scenario := range []string{"setup failure", "cleanup failure"} {
		t.Run(scenario, func(t *testing.T) {
			repo, _ := secondOpinionFixture(t)
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			gitShim := filepath.Join(os.Getenv("PATH"), "git")
			realGit, err := filepath.EvalSymlinks(gitShim)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(gitShim); err != nil {
				t.Fatal(err)
			}
			t.Setenv("SECOND_OPINION_REAL_GIT", realGit)
			t.Setenv("SECOND_OPINION_GIT_FAILURE", scenario)
			script := `#!/bin/sh
case "$*" in
 *'worktree add '*)
  "$SECOND_OPINION_REAL_GIT" "$@" || exit
  if [ "$SECOND_OPINION_GIT_FAILURE" = 'setup failure' ]; then exit 71; fi
  exit 0;;
 *'worktree remove '*)
  if [ "$SECOND_OPINION_GIT_FAILURE" = 'cleanup failure' ]; then exit 72; fi;;
esac
exec "$SECOND_OPINION_REAL_GIT" "$@"
`
			if err := os.WriteFile(gitShim, []byte(script), 0755); err != nil {
				t.Fatal(err)
			}
			secondOpinionFakeCodex(t, "")
			row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
			if row["state"] != "failed" {
				t.Fatalf("failure not reported: %v", row)
			}
			if scenario == "cleanup failure" {
				if !strings.Contains(row["reason"].(string), "worktree cleanup failed") {
					t.Fatalf("cleanup failure hidden: %v", row)
				}
				// The injected Git refusal leaves an honestly reported stale registration.
				// Remove this fixture-only residue before asserting no other worktree remains.
				t.Setenv("SECOND_OPINION_GIT_FAILURE", "")
				toolboxTestGit(t, repo, "worktree", "prune")
			}
			secondOpinionAssertNoWorktree(t, repo, scratch)
		})
	}
}

func TestSecondOpinionWorktreeSetupCancellation(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	scratch, signals := t.TempDir(), t.TempDir()
	t.Setenv("TMPDIR", scratch)
	t.Setenv("SECOND_OPINION_READY", filepath.Join(signals, "ready"))
	t.Setenv("SECOND_OPINION_PID", filepath.Join(signals, "pid"))
	gitShim := filepath.Join(os.Getenv("PATH"), "git")
	realGit, err := filepath.EvalSymlinks(gitShim)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(gitShim); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SECOND_OPINION_REAL_GIT", realGit)
	script := `#!/bin/sh
case "$*" in *'worktree add '*)
 "$SECOND_OPINION_REAL_GIT" "$@" || exit
 /bin/sleep 30 &
 printf '%s' "$!" > "$SECOND_OPINION_PID"
 printf ready > "$SECOND_OPINION_READY"
 wait
 exit 0;;
esac
exec "$SECOND_OPINION_REAL_GIT" "$@"
`
	if err := os.WriteFile(gitShim, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	row := resultmodel.AIReportSecondOpinionRun{BackendName: "codex", RunState: "failed"}
	done := make(chan struct{})
	go func() {
		defer close(done)
		secondOpinionRun(ctx, repo, bundle, "unused", "unused", secondOpinionDirection{}, &row, &resultmodel.CommandResult{})
	}()
	secondOpinionAwaitFile(t, os.Getenv("SECOND_OPINION_READY"))
	cancel()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("setup cancellation did not clean up")
	}
	if !strings.Contains(row.FailureReason, "context canceled") || strings.Contains(row.FailureReason, "cleanup failed") {
		t.Fatalf("setup failure=%+v", row)
	}
	secondOpinionAssertProcessGone(t, os.Getenv("SECOND_OPINION_PID"))
	secondOpinionAssertNoWorktree(t, repo, scratch)
}
