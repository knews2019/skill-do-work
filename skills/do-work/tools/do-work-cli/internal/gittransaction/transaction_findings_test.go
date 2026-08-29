package gittransaction

import (
	"os"
	"regexp"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// declaredFailureKinds reads the kinds back out of their single declaration site instead of
// restating them here. A hand-copied list in this file would go stale the moment a kind is
// added, which is the exact shape this test exists to catch.
func declaredFailureKinds(t *testing.T) []FailureKind {
	t.Helper()
	source, err := os.ReadFile("git_transaction.go")
	if err != nil {
		t.Fatal(err)
	}
	matches := regexp.MustCompile(`FailureKind = "([a-z_]+)"`).FindAllSubmatch(source, -1)
	if len(matches) == 0 {
		t.Fatal("no FailureKind constants found; the declaration shape changed")
	}
	kinds := make([]FailureKind, 0, len(matches))
	for _, match := range matches {
		kinds = append(kinds, FailureKind(match[1]))
	}
	return kinds
}

func TestFindingCodeIsDerivedFromTheFailureKind(t *testing.T) {
	tests := map[FailureKind]string{
		FailureDirtyTarget:   "GIT-DIRTY-TARGET",
		FailureCommittedRisk: "GIT-COMMITTED-STATE-RISK",
		FailureNotGit:        "GIT-NOT-GIT-REPOSITORY",
	}
	for kind, want := range tests {
		if got := FindingCode(kind); got != want {
			t.Errorf("FindingCode(%q) = %q, want %q", kind, got, want)
		}
	}
}

func TestEveryDeclaredFailureKindProducesACompleteFinding(t *testing.T) {
	for _, kind := range declaredFailureKinds(t) {
		result := TransactionResult{
			Outcome:        resultmodel.OutcomeRefused,
			RepositoryRoot: "/tmp/example",
			CommitSHA:      "0123456789abcdef0123456789abcdef01234567",
			RevertArgv:     []string{"git", "revert", "0123456789abcdef0123456789abcdef01234567"},
			Failure:        &TransactionFailure{Kind: kind, Reason: "forced " + string(kind), Paths: []string{"target.txt"}},
		}
		commandResult := BuildCommandResult(result)
		if len(commandResult.Findings) != 1 {
			t.Fatalf("kind %q produced %d findings", kind, len(commandResult.Findings))
		}
		finding := commandResult.Findings[0]
		if finding.Code != FindingCode(kind) {
			t.Errorf("kind %q code = %q, want %q", kind, finding.Code, FindingCode(kind))
		}
		if finding.Severity == "" || finding.Fixability == "" || finding.AutomationStopReason == "" {
			t.Errorf("kind %q finding is incomplete: %#v", kind, finding)
		}
		if len(finding.Evidence) == 0 || len(finding.AffectedPaths) == 0 {
			t.Errorf("kind %q finding names no evidence or paths: %#v", kind, finding)
		}
		if len(finding.NextArgv) == 0 && finding.NextJustRecipe == "" {
			t.Errorf("kind %q finding names no next step: %#v", kind, finding)
		}
		if len(finding.VerificationArgv) == 0 {
			t.Errorf("kind %q finding names no verification command: %#v", kind, finding)
		}
	}
}

// A kind with no remediation template must fail loudly rather than emit a finding that
// looks complete and says nothing.
func TestUnmappedFailureKindFailsLoudly(t *testing.T) {
	result := TransactionResult{
		Outcome: resultmodel.OutcomeFailure,
		Failure: &TransactionFailure{Kind: FailureKind("invented_after_this_test"), Reason: "synthetic"},
	}
	findings := BuildCommandResult(result).Findings
	if len(findings) != 1 || findings[0].Code != "GIT-UNMAPPED-FAILURE" {
		t.Fatalf("unmapped kind findings = %#v", findings)
	}
	if findings[0].Severity != resultmodel.SeverityError || len(findings[0].VerificationArgv) == 0 {
		t.Fatalf("unmapped fallback is not itself a complete error finding: %#v", findings[0])
	}
}

func TestSuccessfulTransactionCarriesTruthfulChangeKinds(t *testing.T) {
	result := TransactionResult{
		Outcome:      resultmodel.OutcomeSuccess,
		ChangedPaths: []string{"created.txt", "tracked.txt"},
		CreatedPaths: []string{"created.txt"},
	}
	changes := BuildCommandResult(result).Changes
	if len(changes) != 2 || changes[0].Kind != "created" || changes[1].Kind != "modified" {
		t.Fatalf("changes = %#v", changes)
	}
	if changes[0].Detail == "" || changes[1].Detail == "" {
		t.Fatalf("changes carry no detail: %#v", changes)
	}
}

// A completed rollback undid every recorded change, so reporting them would be a lie.
func TestCompletedRollbackReportsNoSurvivingChanges(t *testing.T) {
	result := TransactionResult{
		Outcome:      resultmodel.OutcomeRolledBack,
		ChangedPaths: []string{"tracked.txt"},
		Rollback:     resultmodel.RollbackResult{Status: resultmodel.RollbackSucceeded},
		Failure:      &TransactionFailure{Kind: FailureMutation, Reason: "forced", Paths: []string{"tracked.txt"}},
	}
	if changes := BuildCommandResult(result).Changes; len(changes) != 0 {
		t.Fatalf("rolled-back result reported surviving changes: %#v", changes)
	}
}
