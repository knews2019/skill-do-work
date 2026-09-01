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
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestAuditPercentilesAndBands(t *testing.T) {
	s := auditSummary([]int{1, 2, 3, 4, 5})
	if s.median != 3 || s.p90 != 5 || s.p95 != 5 || s.max != 5 {
		t.Fatalf("summary=%+v", s)
	}
	if auditBand(10, 10, 20) != "" || auditBand(11, 10, 20) != "WATCH" || auditBand(21, 10, 20) != "FLAG" {
		t.Fatal("strict band edge regressed")
	}
}

func TestRemediationAuditStandaloneDifferentialAndTypedOrder(t *testing.T) {
	if strings.HasPrefix(runtime.Version(), "go1.25") {
		t.Skip("the retained standalone oracle declares Go 1.26; the canonical CLI itself remains covered by the Go 1.25 lane")
	}
	repository := toolboxTestRepository(t)
	if err := os.MkdirAll(filepath.Join(repository, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "src", "large.go"), []byte("one two\nthree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "README"), []byte("unterminated text"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "binary.bin"), []byte{'a', 0, 'b'}, 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", ".")
	toolboxTestGit(t, repository, "commit", "-m", "initial")
	if err := os.Rename(filepath.Join(repository, "src", "large.go"), filepath.Join(repository, "src", "renamed.go")); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "-A")
	toolboxTestGit(t, repository, "commit", "-m", "rename")
	oracleSource, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "..", "do-work-toolbox", "tools", "audit-metrics"))
	if err != nil {
		t.Fatal(err)
	}
	oracle := filepath.Join(t.TempDir(), "audit-metrics")
	command := exec.Command("go", "build", "-o", oracle, ".")
	command.Dir = oracleSource
	if output, buildErr := command.CombinedOutput(); buildErr != nil {
		t.Fatalf("build standalone oracle: %v: %s", buildErr, output)
	}
	context := commandruntime.ExecutionContext{RepositoryRoot: repository}
	for _, testCase := range []struct {
		mode  string
		flags []string
	}{{"inventory", []string{"--top-count", "2", "--watch-lines", "1"}}, {"folders", []string{"--top-count", "2", "--watch-files", "1"}}, {"churn", []string{"--since-window", "20 years"}}, {"hotspots", []string{"--since-window", "20 years"}}} {
		arguments := append([]string{testCase.mode, "--repo-root", repository}, testCase.flags...)
		oracleCommand := exec.Command(oracle, arguments...)
		oracleOutput, oracleErr := oracleCommand.Output()
		if oracleErr != nil {
			t.Fatalf("oracle %s: %v", testCase.mode, oracleErr)
		}
		result := handleAuditMetrics(context, append([]string{testCase.mode, "--repo-root", repository}, testCase.flags...))
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || !bytes.Equal(oracleOutput, []byte(*result.ExactTextOutput)) {
			t.Fatalf("%s differential mismatch\nold=%s\nnew=%+v", testCase.mode, oracleOutput, result)
		}
		if result.AuditMetrics == nil || result.AuditMetrics.Kind != testCase.mode {
			t.Fatalf("%s missing typed projection", testCase.mode)
		}
		if testCase.mode == "churn" && (len(result.AuditMetrics.Churn) == 0 || result.AuditMetrics.Churn[0].Path != "src/renamed.go") {
			t.Fatalf("rename-normalized typed order=%+v", result.AuditMetrics.Churn)
		}
	}
}

func TestRemediationAuditComparatorRejectsStatusTextAndOrderMutations(t *testing.T) {
	compare := func(oldStatus int, oldText string, newStatus int, newText string) bool {
		return oldStatus == newStatus && oldText == newText
	}
	for _, mutation := range []struct {
		status int
		text   string
	}{{1, "a\nb\n"}, {0, "changed\n"}, {0, "b\na\n"}} {
		if compare(0, "a\nb\n", mutation.status, mutation.text) {
			t.Fatalf("comparator missed mutation %+v", mutation)
		}
	}
}

func TestParseAuditOptions(t *testing.T) {
	o, err := parseAuditOptions(".", []string{"inventory", "--exclude-path", "CHANGELOG.md", "--top-count=3"})
	if err != nil || o.top != 3 || len(o.excludes) != 1 {
		t.Fatalf("options=%+v err=%v", o, err)
	}
	if _, err := parseAuditOptions(".", []string{"inventory", "leftover"}); err == nil {
		t.Fatal("leftover accepted")
	}
}

func TestRemediationAuditExplicitDefaultsRemainWrongSubcommandFlags(t *testing.T) {
	for _, arguments := range [][]string{
		{"inventory", "--since-window", "12 months"},
		{"churn", "--watch-lines", "-1"},
		{"inventory", "--watch-files", "-1"},
		{"folders", "--watch-lines", "-1"},
	} {
		if _, err := parseAuditOptions(".", arguments); err == nil {
			t.Errorf("accepted malformed explicit-default argv %v", arguments)
		}
	}
}
