package requeststate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

func TestResolveTargetRequiresExactUniqueContainedIdentity(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeStateRequest(t, repositoryRoot, "do-work/queue/REQ-001-one.md", "REQ-001", "pending", "")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	target, refusal := ResolveTarget(snapshot, "REQ-001", "do-work/queue/REQ-001-one.md")
	if refusal != nil || target == nil {
		t.Fatalf("exact target refused: %#v", refusal)
	}
	if _, refusal := ResolveTarget(snapshot, "REQ-001", "../REQ-001.md"); refusal == nil || refusal.Code != "REQUEST-PATH-UNSAFE" {
		t.Fatalf("unsafe path refusal = %#v", refusal)
	}
	if _, refusal := ResolveTarget(snapshot, "REQ-001", "do-work/queue/REQ-001-other.md"); refusal == nil || refusal.Code != "REQUEST-SNAPSHOT-STALE" {
		t.Fatalf("stale path refusal = %#v", refusal)
	}

	writeStateRequest(t, repositoryRoot, "do-work/archive/REQ-001-duplicate.md", "REQ-001", "completed", "")
	snapshot, _ = repositorymodel.DiscoverRepository(repositoryRoot)
	if _, refusal := ResolveTarget(snapshot, "REQ-001", ""); refusal == nil || refusal.Code != "REQUEST-AMBIGUOUS" {
		t.Fatalf("duplicate refusal = %#v", refusal)
	}
}

func writeStateRequest(t *testing.T, repositoryRoot, relativePath, requestID, status, extra string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	contents := "---\nid: " + requestID + "\ntitle: Fixture " + requestID + "\nstatus: " + status + "\n" + extra + "---\nBody\n"
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
