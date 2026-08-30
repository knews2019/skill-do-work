package dependencygraph

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

type graphFixture struct {
	requestID     string
	status        string
	dependencyKey string
	dependencies  []string
}

func buildFixtureGraph(t *testing.T, fixtures []graphFixture) *DependencyGraph {
	t.Helper()
	repositoryRoot := t.TempDir()
	for _, fixture := range fixtures {
		frontmatter := "---\nid: " + fixture.requestID + "\nstatus: " + fixture.status + "\n"
		if fixture.dependencyKey != "" {
			frontmatter += fixture.dependencyKey + ": ["
			for dependencyIndex, dependencyID := range fixture.dependencies {
				if dependencyIndex > 0 {
					frontmatter += ", "
				}
				frontmatter += dependencyID
			}
			frontmatter += "]\n"
		}
		frontmatter += "---\nBody\n"
		requestPath := filepath.Join(repositoryRoot, "do-work", "queue", fixture.requestID+".md")
		if err := os.MkdirAll(filepath.Dir(requestPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(requestPath, []byte(frontmatter), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	return BuildGraph(snapshot)
}

func TestBuildGraphComputesReadinessReverseEdgesAndDepth(t *testing.T) {
	graph := buildFixtureGraph(t, []graphFixture{
		{"REQ-001", "completed", "", nil},
		{"REQ-002", "pending", "dependencies", []string{"REQ-001"}},
		{"REQ-003", "pending", "depends_on", []string{"REQ-001", "REQ-002", "REQ-002"}},
		{"REQ-004", "pending", "depends_on", []string{"REQ-999"}},
		{"REQ-005", "pending", "depends_on", []string{"REQ-006"}},
		{"REQ-006", "pending", "depends_on", []string{"REQ-005"}},
	})

	if !graph.NodesByID["REQ-002"].IsReady || graph.NodesByID["REQ-002"].DependencyDepth != 1 {
		t.Fatalf("REQ-002 = %#v", graph.NodesByID["REQ-002"])
	}
	if graph.NodesByID["REQ-003"].IsReady || graph.NodesByID["REQ-003"].DependencyDepth != 2 || !reflect.DeepEqual(graph.NodesByID["REQ-003"].UnmetDependencies, []string{"REQ-002"}) {
		t.Fatalf("REQ-003 = %#v", graph.NodesByID["REQ-003"])
	}
	if !reflect.DeepEqual(graph.NodesByID["REQ-001"].DependentIDs, []string{"REQ-002", "REQ-003"}) {
		t.Fatalf("REQ-001 reverse edges = %v", graph.NodesByID["REQ-001"].DependentIDs)
	}
	if !reflect.DeepEqual(graph.NodesByID["REQ-004"].MissingTargets, []string{"REQ-999"}) || graph.NodesByID["REQ-004"].DependencyDepth != -1 {
		t.Fatalf("REQ-004 = %#v", graph.NodesByID["REQ-004"])
	}
	if len(graph.DependencyCycles) != 1 || !reflect.DeepEqual(graph.DependencyCycles[0].RequestIDs, []string{"REQ-005", "REQ-006"}) {
		t.Fatalf("cycles = %#v", graph.DependencyCycles)
	}
	if !graph.NodesByID["REQ-005"].IsCyclic || !graph.NodesByID["REQ-006"].IsCyclic || graph.NodesByID["REQ-005"].DependencyDepth != -1 {
		t.Fatalf("cycle nodes = %#v %#v", graph.NodesByID["REQ-005"], graph.NodesByID["REQ-006"])
	}
	if len(graph.DependencyEdges) != 6 {
		t.Fatalf("edges = %#v, want 6 de-duplicated edges", graph.DependencyEdges)
	}
}

func TestDependencySatisfactionUsesOnlyTerminalSuccess(t *testing.T) {
	graph := buildFixtureGraph(t, []graphFixture{
		{"REQ-001", "completed-with-issues", "", nil},
		{"REQ-002", "cancelled", "", nil},
		{"REQ-003", "failed", "", nil},
		{"REQ-010", "pending", "depends_on", []string{"REQ-001"}},
		{"REQ-011", "pending", "depends_on", []string{"REQ-002"}},
		{"REQ-012", "pending", "depends_on", []string{"REQ-003"}},
	})
	if !graph.NodesByID["REQ-010"].IsReady {
		t.Fatal("completed-with-issues must satisfy a dependency")
	}
	if graph.NodesByID["REQ-011"].IsReady || graph.NodesByID["REQ-012"].IsReady {
		t.Fatal("cancelled and failed targets must not satisfy dependencies")
	}
}

func TestFilenameFrontmatterCollisionMakesDependencyAmbiguous(t *testing.T) {
	repositoryRoot := t.TempDir()
	fixtures := map[string]string{
		"do-work/queue/REQ-020-first.md":     "---\nid: REQ-021\nstatus: completed\n---\nBody\n",
		"do-work/archive/REQ-021-second.md":  "---\nid: REQ-022\nstatus: completed\n---\nBody\n",
		"do-work/queue/REQ-030-dependent.md": "---\nid: REQ-030\nstatus: pending\ndepends_on: [REQ-021]\n---\nBody\n",
	}
	for relativePath, contents := range fixtures {
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.CollisionEntries) != 1 || snapshot.CollisionEntries[0].RequestID != "REQ-021" {
		t.Fatalf("fixture collision evidence = %#v", snapshot.CollisionEntries)
	}

	graph := BuildGraph(snapshot)
	dependent := graph.NodesByID["REQ-030"]
	if dependent == nil {
		t.Fatal("REQ-030 node is missing")
	}
	if dependent.IsReady || dependent.DependenciesSatisfied {
		t.Fatalf("REQ-030 was ready through a collided dependency: %#v", dependent)
	}
	if !reflect.DeepEqual(dependent.AmbiguousTargets, []string{"REQ-021"}) || !reflect.DeepEqual(dependent.UnmetDependencies, []string{"REQ-021"}) {
		t.Fatalf("REQ-030 ambiguity evidence = %#v", dependent)
	}
	if dependent.DependencyDepth != -1 {
		t.Fatalf("REQ-030 dependency depth = %d, want unresolved -1", dependent.DependencyDepth)
	}
	if target := graph.NodesByID["REQ-021"]; target == nil || !target.IsAmbiguous {
		t.Fatalf("REQ-021 target = %#v, want ambiguous", target)
	}
}
