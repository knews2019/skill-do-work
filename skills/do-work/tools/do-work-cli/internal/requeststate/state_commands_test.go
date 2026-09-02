package requeststate

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestLifecycleCommandsAreRegistered(t *testing.T) {
	repositoryRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "queue"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, commandName := range []string{"claim", "recover-claim", "unblock", "complete", "fail", "cancel"} {
		command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json", commandName)
		output, _ := command.CombinedOutput()
		if bytes.Contains(output, []byte("UNKNOWN-COMMAND")) || !strings.Contains(string(output), `"command": "`+commandName+`"`) {
			t.Errorf("%s is not a registered lifecycle command:\n%s", commandName, output)
		}
	}
}

func TestLifecycleTextAndJSONRenderTheSameActionableResult(t *testing.T) {
	repositoryRoot := newStateRepository(t)
	writeStateRequest(t, repositoryRoot, "do-work/queue/REQ-320.md", "REQ-320", "pending", "")
	result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, TransitionClaim,
		[]string{"REQ-320", "--dry-run", "--writer", "host:/repo", "--at", "2026-08-31T21:00:00Z"})
	textOutput, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	jsonOutput, err := resultmodel.RenderResult(result, resultmodel.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{"claim: success", "STATE-DRY-RUN", "REQ-320", "do-work/queue/REQ-320.md", "next:", "just:", "verify:"} {
		if !bytes.Contains(textOutput, []byte(token)) {
			t.Errorf("text output missing %q:\n%s", token, textOutput)
		}
	}
	if decoded.Command != "claim" || decoded.Outcome != resultmodel.OutcomeSuccess || len(decoded.Findings) != 1 || len(decoded.Changes) == 0 {
		t.Fatalf("JSON did not preserve typed result: %#v", decoded)
	}
	if bytes.Contains(jsonOutput, []byte(": null")) {
		t.Fatalf("JSON contains null collection:\n%s", jsonOutput)
	}
}
