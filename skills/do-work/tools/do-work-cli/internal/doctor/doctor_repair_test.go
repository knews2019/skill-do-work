package doctor

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

func TestRepairTimestampsPreservesUnrelatedBytesCommentsCRLFAndMode(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	relativePath := "do-work/archive/REQ-030-time.md"
	original := []byte("\xef\xbb\xbf---\r\nid: REQ-030\r\ncreated_at: 2099-01-01 00:00:00 # preserve\r\nunknown: 'keep me'\r\n---\r\nBody\r\n")
	writeDoctorFixture(t, repositoryRoot, relativePath, string(original))
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.Chmod(absolutePath, 0o640); err != nil {
		t.Fatal(err)
	}
	commitDoctorFixture(t, repositoryRoot, "fixture")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plans, findings := BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	if len(plans) != 1 || len(findings) != 1 || findings[0].Code != "TIMESTAMP-FUTURE" {
		t.Fatalf("plan=%#v findings=%#v", plans, findings)
	}
	result := ApplyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{})
	if result.Outcome != "success" {
		t.Fatalf("repair result = %#v", result)
	}
	updated, err := os.ReadFile(absolutePath)
	if err != nil {
		t.Fatal(err)
	}
	want := bytes.Replace(original, []byte("2099-01-01 00:00:00"), []byte("2026-08-30T19:00:00Z"), 1)
	if !bytes.Equal(updated, want) {
		t.Fatalf("updated bytes = %q, want %q", updated, want)
	}
	info, err := os.Stat(absolutePath)
	if err != nil || info.Mode().Perm() != 0o640 {
		t.Fatalf("mode = %v, err=%v", info.Mode().Perm(), err)
	}
}

func TestUncommittedTimestampApplyAcceptsDirtyPreimageAndRefusesRace(t *testing.T) {
	root := t.TempDir()
	relative := "do-work/queue/REQ-099.md"
	absolute := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	before := []byte("before\n")
	after := []byte("after\n")
	if err := os.WriteFile(absolute, before, 0o640); err != nil {
		t.Fatal(err)
	}
	snapshot := &repositorymodel.RepositorySnapshot{RepositoryRoot: root}
	plan := TimestampRepairPlan{RelativePath: relative, ExpectedBytes: before, UpdatedBytes: after, Changes: []TimestampFieldChange{{FieldName: "created_at", LineNumber: 2, OldValue: "old", NewValue: "new", Source: "file mtime"}}}
	result := ApplyUncommittedTimestampPlans(snapshot, []TimestampRepairPlan{plan})
	if result.Outcome != "success" || string(readDoctorFixture(t, absolute)) != "after\n" {
		t.Fatalf("result=%#v", result)
	}
	if err := os.WriteFile(absolute, []byte("raced\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	result = ApplyUncommittedTimestampPlans(snapshot, []TimestampRepairPlan{plan})
	if result.Outcome != "findings" || string(readDoctorFixture(t, absolute)) != "raced\n" {
		t.Fatalf("race result=%#v", result)
	}
}

func readDoctorFixture(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func TestRepairTimestampsRefusesDirtyTargetAndCommitWithDirtyIndex(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-031-time.md", doctorRequest("REQ-031", "pending", "created_at: 2099-01-01T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "unrelated.txt", "first\n")
	commitDoctorFixture(t, repositoryRoot, "fixture")
	requestPath := filepath.Join(repositoryRoot, "do-work/queue/REQ-031-time.md")
	requestBytes, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(requestPath, append(requestBytes, []byte("dirty\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plans, _ := BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	result := ApplyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{})
	if result.Outcome != "refused" && result.Outcome != "findings" {
		t.Fatalf("dirty target result = %#v", result)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("dirty target findings = %#v", result.Findings)
	}
	dirtyTargetFinding := result.Findings[0]
	if dirtyTargetFinding.Code != "GIT-DIRTY-TARGET" ||
		!reflect.DeepEqual(dirtyTargetFinding.NextArgv, []string{"git", "status", "--short", "--", "do-work/queue/REQ-031-time.md"}) ||
		!reflect.DeepEqual(dirtyTargetFinding.VerificationArgv, []string{"git", "diff", "--quiet", "--exit-code", "--", "do-work/queue/REQ-031-time.md"}) {
		t.Fatalf("dirty target remediation = %#v", dirtyTargetFinding)
	}
	command := exec.Command("git", "-C", repositoryRoot, "restore", "--", "do-work/queue/REQ-031-time.md")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("restore: %v: %s", err, output)
	}
	if err := os.WriteFile(filepath.Join(repositoryRoot, "unrelated.txt"), []byte("staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	command = exec.Command("git", "-C", repositoryRoot, "add", "unrelated.txt")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("stage: %v: %s", err, output)
	}
	snapshot, _ = repositorymodel.DiscoverRepository(repositoryRoot)
	plans, _ = BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	result = ApplyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{Commit: true})
	if result.Outcome != "refused" {
		t.Fatalf("dirty index result = %#v", result)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("dirty index findings = %#v", result.Findings)
	}
	finding := result.Findings[0]
	if finding.Code != "GIT-DIRTY-INDEX" ||
		!reflect.DeepEqual(finding.NextArgv, []string{"git", "diff", "--cached", "--name-only"}) ||
		!reflect.DeepEqual(finding.VerificationArgv, []string{"git", "diff", "--cached", "--quiet", "--exit-code"}) {
		t.Fatalf("dirty index remediation = %#v", finding)
	}
}

func TestRepairTimestampsDryRunAndCommitTouchOnlyTheExactTarget(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	relativePath := "do-work/queue/REQ-032-time.md"
	writeDoctorFixture(t, repositoryRoot, relativePath, doctorRequest("REQ-032", "pending", "created_at: 2099-01-01T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "unrelated.txt", "keep\n")
	commitDoctorFixture(t, repositoryRoot, "fixture")
	original, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plans, _ := BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	dryRunResult := ApplyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{DryRun: true})
	afterDryRun, _ := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if dryRunResult.Outcome != "success" || !bytes.Equal(original, afterDryRun) || len(dryRunResult.Changes) != 1 {
		t.Fatalf("dry run result=%#v changed=%t", dryRunResult, !bytes.Equal(original, afterDryRun))
	}
	commitResult := ApplyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{Commit: true})
	if commitResult.Outcome != "success" {
		t.Fatalf("commit result = %#v", commitResult)
	}
	command := exec.Command("git", "-C", repositoryRoot, "diff-tree", "--root", "--no-commit-id", "--name-only", "-r", "HEAD")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != relativePath+"\n" {
		t.Fatalf("committed paths = %q", output)
	}
}
