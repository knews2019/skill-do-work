package requeststate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// The finalizer regenerates its lifecycle plan from a journal whose recorded
// preimages already equal the pipeline's own uncommitted edits. Without this
// acceptance, a `complete` transition refused the working REQ it had itself
// been editing all run, and startup recovery could never get past it.
func TestCompleteAcceptsDirtyTargetOnlyWhenBytesHashToRecordedPreimage(t *testing.T) {
	requestPath := "do-work/working/REQ-516.md"
	seed := func(t *testing.T) (string, string) {
		root := newStateRepository(t)
		configureStateGit(t, root)
		writeStateRequest(t, root, requestPath, "REQ-516", "claimed", "claimed_at: 2026-09-02T20:00:00Z\n")
		runStateGit(t, root, "add", requestPath)
		runStateGit(t, root, "commit", "-q", "-m", "claim")
		dirty := readStateFile(t, root, requestPath) + "\n## Review\n\nApproved after the claim commit.\n"
		if err := os.WriteFile(filepath.Join(root, requestPath), []byte(dirty), 0o644); err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256([]byte(dirty))
		return root, hex.EncodeToString(digest[:])
	}
	buildComplete := func(t *testing.T, root string, digests map[string]string) StatePlan {
		snapshot, err := discoverRepository(root)
		if err != nil {
			t.Fatal(err)
		}
		return BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionComplete, RequestID: "REQ-516", RequestPath: requestPath,
			TerminalStatus: "completed", WriterLabel: "host:/repo", Now: time.Date(2026, 9, 2, 21, 0, 0, 0, time.UTC), AcceptedPreimageDigests: digests})
	}

	t.Run("matching recorded digest is accepted", func(t *testing.T) {
		root, digest := seed(t)
		plan := buildComplete(t, root, map[string]string{requestPath: digest})
		if len(plan.ExistingDirtyTargetPaths) != 1 || plan.ExistingDirtyTargetPaths[0] != requestPath {
			t.Fatalf("accepted dirty targets = %v", plan.ExistingDirtyTargetPaths)
		}
		result := ApplyPlan(context.Background(), plan)
		assertStateSuccess(t, result)
		if _, err := os.Stat(filepath.Join(root, "do-work/archive/REQ-516.md")); err != nil {
			t.Fatalf("complete did not archive the dirty working REQ: %v", err)
		}
	})
	t.Run("mismatched recorded digest still refuses", func(t *testing.T) {
		root, _ := seed(t)
		plan := buildComplete(t, root, map[string]string{requestPath: "0000000000000000000000000000000000000000000000000000000000000000"})
		if len(plan.ExistingDirtyTargetPaths) != 0 {
			t.Fatalf("mismatched digest accepted dirty targets = %v", plan.ExistingDirtyTargetPaths)
		}
		result := ApplyPlan(context.Background(), plan)
		if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) == 0 || result.Findings[0].Code != "GIT-DIRTY-TARGET" {
			t.Fatalf("mismatched digest = %#v", result)
		}
	})
}
