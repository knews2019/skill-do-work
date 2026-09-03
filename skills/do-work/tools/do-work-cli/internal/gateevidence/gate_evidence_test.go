package gateevidence

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestGreenGateEvidenceIdentityAndHistory(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	recorded, err := RecordGreenGate(repositoryRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	if recorded.State != resultmodel.GateEvidenceRecorded || recorded.RecordPath == "" {
		t.Fatalf("recorded = %#v", recorded)
	}

	exact, err := CheckGreenGate(repositoryRoot, gateCommand)
	if err != nil || !exact.Matches || exact.MatchBasis != "exact_revision" || exact.BaselineRevision != exact.HeadRevision {
		t.Fatalf("exact = %#v, err=%v", exact, err)
	}

	writeGateEvidenceTestFile(t, repositoryRoot, "_dev/gate-runs/one.log", "pass\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "_dev/gate-runs/one.log")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "gate log")
	logOnly, err := CheckGreenGate(repositoryRoot, gateCommand)
	if err != nil || !logOnly.Matches || logOnly.State != resultmodel.GateEvidenceLogDescendantMatch || logOnly.BaselineRevision != logOnly.HeadRevision {
		t.Fatalf("log only = %#v, err=%v", logOnly, err)
	}

	writeGateEvidenceTestFile(t, repositoryRoot, "project.txt", "project change\n")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qam", "project change")
	invalidated, err := CheckGreenGate(repositoryRoot, gateCommand)
	if err != nil || invalidated.Matches || invalidated.State != resultmodel.GateEvidenceInvalidated {
		t.Fatalf("invalidated = %#v, err=%v", invalidated, err)
	}

	differentArgv, err := CheckGreenGate(repositoryRoot, []string{"bash", "other.sh"})
	if err != nil || differentArgv.Matches || differentArgv.State != resultmodel.GateEvidenceMissing {
		t.Fatalf("different argv = %#v, err=%v", differentArgv, err)
	}
}

// An already-green repository-gate repair no-op has no diff, so an unrelated commit
// moving HEAD cannot make it the cause of a red gate. Its recorded evidence must stay
// verifiable at its own recorded revision while ordinary REQs keep the HEAD-bound rule.
func TestGreenGateEvidenceStaysVerifiableAtRecordedRevisionAfterUnrelatedHeadMove(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	recorded, err := RecordGreenGate(repositoryRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	recordedRevision := recorded.RecordedRevision

	writeGateEvidenceTestFile(t, repositoryRoot, "peer.txt", "someone else's commit\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "peer.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "unrelated peer commit")

	headBound, err := CheckGreenGate(repositoryRoot, gateCommand)
	if err != nil || headBound.Matches || headBound.State != resultmodel.GateEvidenceInvalidated {
		t.Fatalf("ordinary HEAD-bound check must stop matching: %#v, err=%v", headBound, err)
	}

	atRecorded, err := checkGreenGateAtRevision(repositoryRoot, gateCommand, recordedRevision)
	if err != nil || !atRecorded.Matches || atRecorded.State != resultmodel.GateEvidenceExactRevisionMatch {
		t.Fatalf("recorded-revision check = %#v, err=%v", atRecorded, err)
	}
	if atRecorded.TargetRevision != recordedRevision || atRecorded.BaselineRevision != recordedRevision {
		t.Fatalf("target/baseline must be the recorded revision, not HEAD: %#v", atRecorded)
	}
	if atRecorded.HeadRevision == recordedRevision {
		t.Fatalf("fixture did not move HEAD: %#v", atRecorded)
	}

	atMovedHead, err := checkGreenGateAtRevision(repositoryRoot, gateCommand, atRecorded.HeadRevision)
	if err != nil || atMovedHead.Matches || atMovedHead.State != resultmodel.GateEvidenceInvalidated {
		t.Fatalf("an explicit target still honors the project-change rule: %#v, err=%v", atMovedHead, err)
	}

	differentArgv, err := checkGreenGateAtRevision(repositoryRoot, []string{"bash", "other.sh"}, recordedRevision)
	if err != nil || differentArgv.Matches || differentArgv.State != resultmodel.GateEvidenceMissing {
		t.Fatalf("different argv at the recorded revision = %#v, err=%v", differentArgv, err)
	}

	unresolvable, err := checkGreenGateAtRevision(repositoryRoot, gateCommand, "no-such-revision")
	if err == nil || unresolvable.Matches || unresolvable.State != resultmodel.GateEvidenceInvalidRecord {
		t.Fatalf("an unresolvable target is unverifiable, not a miss: %#v, err=%v", unresolvable, err)
	}
}

// A gate-run log commit on top of the record must not invalidate an explicit target
// either — the same tolerance the HEAD-bound lane already grants.
func TestGreenGateEvidenceTargetToleratesGateLogCommits(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	if _, err := RecordGreenGate(repositoryRoot, gateCommand); err != nil {
		t.Fatal(err)
	}
	writeGateEvidenceTestFile(t, repositoryRoot, "_dev/gate-runs/one.log", "pass\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "_dev/gate-runs/one.log")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "gate log")
	logRevision := strings.TrimSpace(runGateEvidenceTestGit(t, repositoryRoot, "rev-parse", "HEAD"))

	writeGateEvidenceTestFile(t, repositoryRoot, "peer.txt", "unrelated\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "peer.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "unrelated peer commit")

	atLog, err := checkGreenGateAtRevision(repositoryRoot, gateCommand, logRevision)
	if err != nil || !atLog.Matches || atLog.State != resultmodel.GateEvidenceLogDescendantMatch || atLog.BaselineRevision != logRevision {
		t.Fatalf("gate-log descendant target = %#v, err=%v", atLog, err)
	}
}

func TestGreenGateEvidenceRejectsDivergentAndForeignRecords(t *testing.T) {
	gateCommand := []string{"bash", "verify.sh"}
	repositoryRoot := newGateEvidenceRepository(t)
	runGateEvidenceTestGit(t, repositoryRoot, "checkout", "-qb", "recorded-branch")
	writeGateEvidenceTestFile(t, repositoryRoot, "branch.txt", "recorded\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "branch.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "recorded branch")
	if _, err := RecordGreenGate(repositoryRoot, gateCommand); err != nil {
		t.Fatal(err)
	}
	runGateEvidenceTestGit(t, repositoryRoot, "checkout", "-q", "master")
	writeGateEvidenceTestFile(t, repositoryRoot, "other.txt", "other\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "other.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "other branch")
	divergent, err := CheckGreenGate(repositoryRoot, gateCommand)
	if err != nil || divergent.Matches || divergent.State != resultmodel.GateEvidenceRecordedRevisionNotAncestor {
		t.Fatalf("divergent = %#v, err=%v", divergent, err)
	}

	foreignRoot := newGateEvidenceRepository(t)
	foreignContext, err := resolveEvidenceContext(foreignRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Dir(foreignContext.recordPath), 0o700); err != nil {
		t.Fatal(err)
	}
	sourceContext, err := resolveEvidenceContext(repositoryRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(sourceContext.recordPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(foreignContext.recordPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	foreign, err := CheckGreenGate(foreignRoot, gateCommand)
	if err != nil || foreign.Matches || foreign.State != resultmodel.GateEvidenceDifferentRepository {
		t.Fatalf("foreign = %#v, err=%v", foreign, err)
	}
}

func TestGreenGateEvidenceFailsClosedForInvalidTargets(t *testing.T) {
	gateCommand := []string{"bash", "verify.sh"}
	t.Run("malformed record", func(t *testing.T) {
		repositoryRoot := newGateEvidenceRepository(t)
		context, err := resolveEvidenceContext(repositoryRoot, gateCommand)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Dir(context.recordPath), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(context.recordPath, []byte("not-json\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		result, err := CheckGreenGate(repositoryRoot, gateCommand)
		if err == nil || result.State != resultmodel.GateEvidenceInvalidRecord {
			t.Fatalf("result = %#v, err=%v", result, err)
		}
	})

	t.Run("missing recorded revision", func(t *testing.T) {
		repositoryRoot := newGateEvidenceRepository(t)
		result, err := RecordGreenGate(repositoryRoot, gateCommand)
		if err != nil {
			t.Fatal(err)
		}
		contents, err := os.ReadFile(result.RecordPath)
		if err != nil {
			t.Fatal(err)
		}
		var record storedGateEvidence
		if err := json.Unmarshal(contents, &record); err != nil {
			t.Fatal(err)
		}
		record.RecordedRevision = "0000000000000000000000000000000000000000"
		contents, _ = json.Marshal(record)
		if err := os.WriteFile(result.RecordPath, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		checked, err := CheckGreenGate(repositoryRoot, gateCommand)
		if err != nil || checked.Matches || checked.State != resultmodel.GateEvidenceRecordedRevisionMissing {
			t.Fatalf("checked = %#v, err=%v", checked, err)
		}
	})

	t.Run("record argv mismatch", func(t *testing.T) {
		repositoryRoot := newGateEvidenceRepository(t)
		result, err := RecordGreenGate(repositoryRoot, gateCommand)
		if err != nil {
			t.Fatal(err)
		}
		contents, err := os.ReadFile(result.RecordPath)
		if err != nil {
			t.Fatal(err)
		}
		var record storedGateEvidence
		if err := json.Unmarshal(contents, &record); err != nil {
			t.Fatal(err)
		}
		record.GateCommand = []string{"bash", "other.sh"}
		contents, _ = json.Marshal(record)
		if err := os.WriteFile(result.RecordPath, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		checked, err := CheckGreenGate(repositoryRoot, gateCommand)
		if err != nil || checked.Matches || checked.State != resultmodel.GateEvidenceDifferentArgv {
			t.Fatalf("checked = %#v, err=%v", checked, err)
		}
	})

	t.Run("non-regular record", func(t *testing.T) {
		repositoryRoot := newGateEvidenceRepository(t)
		context, err := resolveEvidenceContext(repositoryRoot, gateCommand)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(context.recordPath, 0o700); err != nil {
			t.Fatal(err)
		}
		checked, err := CheckGreenGate(repositoryRoot, gateCommand)
		if err == nil || checked.State != resultmodel.GateEvidenceInvalidRecord {
			t.Fatalf("checked = %#v, err=%v", checked, err)
		}
	})

	t.Run("non-private record", func(t *testing.T) {
		repositoryRoot := newGateEvidenceRepository(t)
		result, err := RecordGreenGate(repositoryRoot, gateCommand)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(result.RecordPath, 0o644); err != nil {
			t.Fatal(err)
		}
		checked, err := CheckGreenGate(repositoryRoot, gateCommand)
		if err == nil || checked.State != resultmodel.GateEvidenceInvalidRecord {
			t.Fatalf("checked = %#v, err=%v", checked, err)
		}
	})
}

func TestGreenGateEvidenceReplacesRecordAndSharesCommonGitDirectory(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	first, err := RecordGreenGate(repositoryRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	writeGateEvidenceTestFile(t, repositoryRoot, "project.txt", "second\n")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qam", "second")
	second, err := RecordGreenGate(repositoryRoot, gateCommand)
	if err != nil || second.RecordedRevision == first.RecordedRevision {
		t.Fatalf("second = %#v, first = %#v, err=%v", second, first, err)
	}
	info, err := os.Stat(second.RecordPath)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("record mode = %v, err=%v", info.Mode(), err)
	}
	directoryInfo, err := os.Stat(filepath.Dir(second.RecordPath))
	if err != nil || directoryInfo.Mode().Perm()&0o077 != 0 {
		t.Fatalf("record directory mode = %v, err=%v", directoryInfo.Mode(), err)
	}

	linkedRoot := filepath.Join(t.TempDir(), "linked")
	runGateEvidenceTestGit(t, repositoryRoot, "worktree", "add", "-q", "--detach", linkedRoot, "HEAD")
	t.Cleanup(func() { runGateEvidenceTestGit(t, repositoryRoot, "worktree", "remove", linkedRoot) })
	linked, err := CheckGreenGate(linkedRoot, gateCommand)
	if err != nil || !linked.Matches || linked.RecordPath != second.RecordPath || linked.RepositoryIdentity != second.RepositoryIdentity {
		t.Fatalf("linked = %#v, second = %#v, err=%v", linked, second, err)
	}
}

func newGateEvidenceRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	runGateEvidenceTestGit(t, repositoryRoot, "init", "-q", "-b", "master")
	runGateEvidenceTestGit(t, repositoryRoot, "config", "user.email", "gate-evidence@example.invalid")
	runGateEvidenceTestGit(t, repositoryRoot, "config", "user.name", "Gate Evidence")
	writeGateEvidenceTestFile(t, repositoryRoot, "project.txt", "initial\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "project.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "initial")
	return repositoryRoot
}

func runGateEvidenceTestGit(t *testing.T, repositoryRoot string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return string(output)
}

func writeGateEvidenceTestFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
