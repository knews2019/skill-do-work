package knowledgecommands

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestHandlersRegisterKnowledgeCommands(t *testing.T) {
	handlers := Handlers()
	for _, name := range []string{CommandBKBInit, CommandBKBStatus, CommandBKBLintStructure, CommandDreamScan,
		CommandInterviewList, CommandInterviewStatus, CommandInterviewExport, CommandInterviewIngest, CommandInterviewReset, CommandInterviewVersions,
		CommandMemoryRemember, CommandMemoryForget, CommandMemoryRecall, CommandMemoryStatus, CommandMemoryBootstrap, CommandMemoryAudit, CommandInstallMemoryHooks, CommandLexicalMemoryRecall} {
		if handlers[name] == nil {
			t.Fatalf("handler %q is not registered", name)
		}
	}
	if len(handlers) != 18 {
		t.Fatalf("registered %d handlers, want exactly 18", len(handlers))
	}
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		relative, _ := filepath.Rel(root, path)
		if _, err := hash.Write([]byte(relative + "\x00" + entry.Type().String() + "\x00")); err != nil {
			return err
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = hash.Write(data)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func findingsSignature(findings []resultmodel.CommandFinding) string {
	parts := make([]string, 0, len(findings))
	for _, finding := range findings {
		parts = append(parts, finding.Code+":"+strings.Join(finding.AffectedPaths, ",")+":"+strings.Join(finding.Evidence, ","))
	}
	return strings.Join(parts, "\n")
}

func TestRealRuntimeTextAndJSONProjectSameDreamFindings(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "memory"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "memory", "MEMORY.md"), []byte("# Index\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "memory", "alpha.md"), []byte("# Alpha\nToday.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var jsonOutput bytes.Buffer
	jsonStatus := commandruntime.NewRuntime(&jsonOutput, Handlers()).Run([]string{"--repo-root", root, "--format", "json", CommandDreamScan, "--path", "memory"})
	if jsonStatus != 1 {
		t.Fatalf("JSON status = %d, want findings status 1: %s", jsonStatus, jsonOutput.String())
	}
	var result resultmodel.CommandResult
	if err := json.Unmarshal(jsonOutput.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	var textOutput bytes.Buffer
	textStatus := commandruntime.NewRuntime(&textOutput, Handlers()).Run([]string{"--repo-root", root, CommandDreamScan, "--path", "memory"})
	if textStatus != jsonStatus {
		t.Fatalf("text status %d != JSON status %d", textStatus, jsonStatus)
	}
	for _, finding := range result.Findings {
		if !strings.Contains(textOutput.String(), finding.Code) {
			t.Errorf("text projection omitted %s", finding.Code)
		}
	}
}

func TestRealRuntimeTextAndJSONProjectSameInterviewAndMemoryFindings(t *testing.T) {
	root := t.TempDir()
	writeInterviewFixture(t, root, "knowledge/interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "memory/working-memory.md", workingMemoryFixture)
	for _, arguments := range [][]string{
		{CommandInterviewList, "--knowledge-root", "knowledge"},
		{CommandMemoryAudit, "--engine", "memory"},
	} {
		var jsonOutput bytes.Buffer
		jsonArguments := append([]string{"--repo-root", root, "--format", "json"}, arguments...)
		jsonStatus := commandruntime.NewRuntime(&jsonOutput, Handlers()).Run(jsonArguments)
		var result resultmodel.CommandResult
		if err := json.Unmarshal(jsonOutput.Bytes(), &result); err != nil {
			t.Fatalf("%s JSON: %v\n%s", arguments[0], err, jsonOutput.String())
		}
		var textOutput bytes.Buffer
		textArguments := append([]string{"--repo-root", root}, arguments...)
		textStatus := commandruntime.NewRuntime(&textOutput, Handlers()).Run(textArguments)
		if textStatus != jsonStatus {
			t.Fatalf("%s text=%d JSON=%d", arguments[0], textStatus, jsonStatus)
		}
		for _, finding := range result.Findings {
			if !strings.Contains(textOutput.String(), finding.Code) {
				t.Fatalf("%s text omitted %s", arguments[0], finding.Code)
			}
		}
	}
}

func TestParseOptionsAcceptsDocumentedPathsAndRejectsMixedMutationModes(t *testing.T) {
	for _, target := range []string{"../outside", filepath.Join(t.TempDir(), "absolute-kb")} {
		options, err := parseBKBOptions([]string{"--kb", target}, true)
		if err != nil || options.target != filepath.Clean(target) {
			t.Fatalf("documented target %q: options=%+v err=%v", target, options, err)
		}
	}
	if _, err := parseBKBOptions([]string{"--dry-run", "--commit"}, true); err == nil {
		t.Fatal("mixed dry-run and commit was accepted")
	}
	if _, err := parseBKBOptions([]string{"--commit"}, false); err == nil {
		t.Fatal("commit was accepted for a read-only command")
	}
}

func TestKnowledgeFindingsCarryTargetPreservingNextAndVerificationCommands(t *testing.T) {
	root := t.TempDir()
	missing := handleBKBStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "knowledge base"})
	if len(missing.Findings) != 1 || strings.Join(missing.Findings[0].NextArgv, " ") != "do-work-cli bkb-init --kb knowledge base" || missing.Findings[0].NextJustRecipe != "bkb-init 'knowledge base'" {
		t.Fatalf("missing-KB remediation = %+v", missing.Findings)
	}
	if err := os.MkdirAll(filepath.Join(root, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "kb", "wiki"), 0o755); err != nil {
		t.Fatal(err)
	}
	exists := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
	if len(exists.Findings) != 1 || !strings.Contains(strings.Join(exists.Findings[0].NextArgv, "\x00"), "--fill-gaps") || !strings.Contains(strings.Join(exists.Findings[0].VerificationArgv, "\x00"), "--dry-run") {
		t.Fatalf("existing-KB remediation = %+v", exists.Findings)
	}
}
