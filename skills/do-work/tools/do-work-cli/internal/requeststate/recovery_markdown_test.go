package requeststate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestRecoveryPreservesHiddenHeadingsAndFollowingRequirements(t *testing.T) {
	for _, example := range []string{
		"```markdown\n## Plan\nexample\n```\nMUST keep this requirement.\n",
		"~~~~\n## Scope\n```\n## Plan\n~~~~\nMUST keep this requirement.\n",
		"<!--\n## Plan\n-->\nMUST keep this requirement.\n",
		"```markdown\r\n## Plan\r\nexample\r\n```\r\nMUST keep this requirement.\r\n",
		"```markdown\n## Plan\nUnclosed example must survive.\n",
	} {
		t.Run(example, func(t *testing.T) {
			original := "---\nid: REQ-501\n---\n# Request\n\n" + example
			got, err := stripGeneratedRecoverySections([]byte(original))
			if err != nil || string(got) != original {
				t.Fatalf("unowned bytes changed: error=%v\ngot %q\nwant %q", err, got, original)
			}
		})
	}
	original := "---\nid: REQ-501\n---\n## Plan\ngenerated\n```\n## Scope\nUnclosed example.\n"
	want := "---\nid: REQ-501\n---\n```\n## Scope\nUnclosed example.\n"
	got, err := stripGeneratedRecoverySections([]byte(original))
	if err != nil || string(got) != want {
		t.Fatalf("unclosed example under generated section changed: %q, %v", got, err)
	}
}

func TestRecoverClaimCommitsWithoutDeletingFencedScopeOrRequirements(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	writeStateRequest(t, root, "do-work/queue/REQ-501.md", "REQ-501", "pending", "write_set: [owned.go]\n")
	writeStateCheckpoint(t, root, "")
	runStateGit(t, root, "add", "do-work")
	runStateGit(t, root, "commit", "-qm", "seed")
	ctx := commandruntime.ExecutionContext{RepositoryRoot: root}
	assertStateSuccess(t, handleStateCommand(ctx, TransitionClaim, []string{"REQ-501", "--request-path", "do-work/queue/REQ-501.md", "--provenance", "explicit-req", "--writer", "foreign:/checkout", "--at", "2026-09-02T01:00:00Z"}))
	example := "\n## Requirements\n\n```markdown\n## Scope\nexample\n## Plan\nexample\n```\nMUST preserve these bytes.\n"
	working := readStateFile(t, root, "do-work/working/REQ-501.md") + example + "\n## Testing\ngenerated evidence\n"
	if err := os.WriteFile(filepath.Join(root, "do-work/working/REQ-501.md"), []byte(working), 0o644); err != nil {
		t.Fatal(err)
	}
	assertStateSuccess(t, handleStateCommand(ctx, TransitionRecover, []string{"REQ-501", "--request-path", "do-work/working/REQ-501.md", "--checkpoint-writer", "foreign:/checkout", "--assume-sole-writer", "--commit", "--at", "2026-09-02T02:00:00Z"}))
	queue := readStateFile(t, root, "do-work/queue/REQ-501.md")
	if !strings.Contains(queue, example) || !strings.Contains(queue, "write_set: [owned.go]") || strings.Contains(queue, "generated evidence") {
		t.Fatalf("recovery damaged request or retained generated evidence:\n%s", queue)
	}
	if status := runStateGit(t, root, "status", "--porcelain"); status != "" {
		t.Fatalf("recovery was not committed: %s", status)
	}
}
