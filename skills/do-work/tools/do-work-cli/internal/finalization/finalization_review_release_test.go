package finalization

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrimaryReleaseRequiresAnActualShippedAllowlistChange(t *testing.T) {
	for _, change := range []string{"unchanged", "addition", "edit", "deletion", "outside allowlist"} {
		t.Run(change, func(t *testing.T) {
			repositoryRoot := newFinalizationRepository(t)
			seedMaintainerSuite(t, repositoryRoot)
			shippedPath := "skills/do-work/actions/work.md"
			commitTouching(t, repositoryRoot, shippedPath)
			writeFinalizationFile(t, repositoryRoot, "_dev/tests/probe.sh", "maintainer-only change\n")
			writeFinalizationFile(t, repositoryRoot, "skills/do-work/VERSION", "1.0.1\n")
			paths := []string{shippedPath, "_dev/tests/probe.sh", "skills/do-work/VERSION"}
			switch change {
			case "addition":
				shippedPath = "skills/do-work/actions/new.md"
				paths = append(paths, shippedPath)
				writeFinalizationFile(t, repositoryRoot, shippedPath, "new shipped content\n")
			case "edit":
				writeFinalizationFile(t, repositoryRoot, shippedPath, "updated shipped content\n")
			case "deletion":
				if err := os.Remove(filepath.Join(repositoryRoot, shippedPath)); err != nil {
					t.Fatal(err)
				}
			case "outside allowlist":
				writeFinalizationFile(t, repositoryRoot, "skills/do-work/actions/unrelated.md", "unrelated shipped change\n")
			}
			err := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenancePrimaryCommit, CommitPaths: paths})
			if change == "unchanged" || change == "outside allowlist" {
				if err == nil || !strings.Contains(err.Error(), "RELEASE-WITHOUT-SHIPPED-CHANGE") {
					t.Fatalf("unchanged shipped allowlist entry authorized maintainer-only release: %v", err)
				}
			} else if err != nil {
				t.Fatalf("actual shipped %s was refused: %v", change, err)
			}
		})
	}
}
