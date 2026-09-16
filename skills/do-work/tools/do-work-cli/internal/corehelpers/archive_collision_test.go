package corehelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestArchiveCollisionFindsNestedRequestsWithoutFollowingDirectoryLinks(t *testing.T) {
	root := t.TempDir()
	writeMatrixFile(t, root, "do-work/archive/UR-222/deeper/REQ-869-old.md", "archived")
	writeMatrixFile(t, root, "do-work/archive/REQ-869.md", "archived")
	writeMatrixFile(t, root, "do-work/archive/UR-222/REQ-8690-other.md", "different ID")
	writeMatrixFile(t, root, "outside/REQ-870.md", "outside archive")
	if err := os.Symlink(filepath.Join(root, "outside"), filepath.Join(root, "do-work/archive/linked")); err != nil {
		t.Fatal(err)
	}
	result := handleArchiveCollision(testContext(root), []string{"REQ-869"})
	want := "do-work/archive/REQ-869.md\ndo-work/archive/UR-222/deeper/REQ-869-old.md\n"
	if result.Outcome != resultmodel.OutcomeFindings || result.ExactTextOutput == nil || *result.ExactTextOutput != want {
		t.Fatalf("nested collision not reported: %+v", result)
	}
	if err := os.Remove(filepath.Join(root, "do-work/archive/REQ-869.md")); err != nil {
		t.Fatal(err)
	}
	result = handleArchiveCollision(testContext(root), []string{"REQ-869"})
	if resultmodel.ExitCode(result.Outcome) != 1 || result.ExactTextOutput == nil || *result.ExactTextOutput != "do-work/archive/UR-222/deeper/REQ-869-old.md\n" {
		t.Fatalf("nested-only collision authorized reuse: %+v", result)
	}
	for _, requestID := range []string{"REQ-86", "REQ-870", "REQ-999"} {
		result := handleArchiveCollision(testContext(root), []string{requestID})
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || *result.ExactTextOutput != "" {
			t.Fatalf("false collision for %s: %+v", requestID, result)
		}
	}
}

func TestArchiveCollisionDistinguishesMissingArchiveFromUnreadableArchive(t *testing.T) {
	root := t.TempDir()
	if result := handleArchiveCollision(testContext(root), []string{"REQ-869"}); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("missing archive failed: %+v", result)
	}
	writeMatrixFile(t, root, "do-work/archive", "not a directory")
	if result := handleArchiveCollision(testContext(root), []string{"REQ-869"}); result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("unreadable archive authorized reuse: %+v", result)
	}
}
