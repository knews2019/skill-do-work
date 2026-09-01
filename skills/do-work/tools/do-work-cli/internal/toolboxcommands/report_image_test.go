package toolboxcommands

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestReportImagePublishesStatusBackedNonemptyOutput(t *testing.T) {
	repository := toolboxTestRepository(t)
	backend := filepath.Join(repository, "imagegen")
	script := "#!/bin/sh\nout=\nwhile [ $# -gt 0 ]; do if [ \"$1\" = --output ]; then out=$2; shift 2; else shift; fi; done\nprintf fresh > \"$out\"\n"
	if err := os.WriteFile(backend, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	originalLookup := reportImageLookPath
	reportImageLookPath = func(name string) (string, error) {
		if name == "imagegen" {
			return backend, nil
		}
		return exec.LookPath(name)
	}
	t.Cleanup(func() { reportImageLookPath = originalLookup })
	output := filepath.Join(repository, "image.png")
	result := handleReportImage(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{output, "style", "literal; $(touch nope)"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result=%+v", result)
	}
	contents, err := os.ReadFile(output)
	if err != nil || string(contents) != "fresh" {
		t.Fatalf("output=%q err=%v", contents, err)
	}
	if _, err := os.Stat(filepath.Join(repository, "nope")); !os.IsNotExist(err) {
		t.Fatal("prompt metacharacters were evaluated")
	}
}

func TestReportImageBatchAllFailedIsSuccessfulAndLeavesNoStage(t *testing.T) {
	repository := toolboxTestRepository(t)
	report := filepath.Join(repository, "report")
	if err := os.Mkdir(report, 0o755); err != nil {
		t.Fatal(err)
	}
	originalLookup := reportImageLookPath
	reportImageLookPath = func(string) (string, error) { return "", exec.ErrNotFound }
	t.Cleanup(func() { reportImageLookPath = originalLookup })
	result := handleReportImageBatch(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{report, "style", "one.png:first", "two.png:second"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result=%+v", result)
	}
	if result.ExactTextOutput == nil || *result.ExactTextOutput != "" || len(result.Findings) != 2 {
		t.Fatalf("per-item fallback evidence missing: %+v", result)
	}
	if _, err := os.Stat(filepath.Join(report, "generated")); !os.IsNotExist(err) {
		t.Fatal("all-failed batch published generated/")
	}
	matches, _ := filepath.Glob(filepath.Join(report, ".generated.staging.*"))
	if len(matches) != 0 {
		t.Fatalf("staging leaked: %v", matches)
	}
}

func toolboxTestRepository(t *testing.T) string {
	t.Helper()
	repository := t.TempDir()
	for _, arguments := range [][]string{{"init"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Test"}} {
		command := exec.Command("git", append([]string{"-C", repository}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", arguments, err, output)
		}
	}
	return repository
}

func toolboxTestGit(t *testing.T, repository string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repository}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
	return string(output)
}
