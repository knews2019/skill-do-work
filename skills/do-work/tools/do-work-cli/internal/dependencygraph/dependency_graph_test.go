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
	// commit is the fixture's implementation commit; treeSection defaults to
	// "queue" so an existing fixture keeps its original location.
	commit      string
	treeSection string
}

func buildFixtureGraph(t *testing.T, fixtures []graphFixture) *DependencyGraph {
	t.Helper()
	repositoryRoot := t.TempDir()
	for _, fixture := range fixtures {
		frontmatter := "---\nid: " + fixture.requestID + "\nstatus: " + fixture.status + "\n"
		if fixture.commit != "" {
			frontmatter += "commit: " + fixture.commit + "\n"
		}
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
		treeSection := fixture.treeSection
		if treeSection == "" {
			treeSection = "queue"
		}
		requestPath := filepath.Join(repositoryRoot, "do-work", treeSection, fixture.requestID+".md")
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
		{requestID: "REQ-001", status: "completed"},
		{requestID: "REQ-002", status: "pending", dependencyKey: "dependencies", dependencies: []string{"REQ-001"}},
		{requestID: "REQ-003", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-001", "REQ-002", "REQ-002"}},
		{requestID: "REQ-004", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-999"}},
		{requestID: "REQ-005", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-006"}},
		{requestID: "REQ-006", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-005"}},
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
		{requestID: "REQ-001", status: "completed-with-issues"},
		{requestID: "REQ-002", status: "cancelled"},
		{requestID: "REQ-003", status: "failed"},
		{requestID: "REQ-010", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-001"}},
		{requestID: "REQ-011", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-002"}},
		{requestID: "REQ-012", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-003"}},
	})
	if !graph.NodesByID["REQ-010"].IsReady {
		t.Fatal("completed-with-issues must satisfy a dependency")
	}
	if graph.NodesByID["REQ-011"].IsReady || graph.NodesByID["REQ-012"].IsReady {
		t.Fatal("cancelled and failed targets must not satisfy dependencies")
	}
}

// TestClaimedDependencyWithCommitIsSourceReady pins REQ-570: a request held for
// heavy lanes stays claimed in do-work/working/ carrying its landed commit, and
// that is what makes its dependents buildable. A claimed request without a
// commit, and a pending request carrying a prior attempt's commit, both stay
// unmet.
func TestClaimedDependencyWithCommitIsSourceReady(t *testing.T) {
	heldGraph := buildFixtureGraph(t, []graphFixture{
		{requestID: "REQ-563", status: "claimed", commit: "0123456789abcdef", treeSection: "working"},
		{requestID: "REQ-564", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-563"}},
	})
	if !heldGraph.NodesByID["REQ-564"].IsReady {
		t.Fatalf("held claimed source did not unblock dependent: %#v", heldGraph.NodesByID["REQ-564"])
	}

	commitlessGraph := buildFixtureGraph(t, []graphFixture{
		{requestID: "REQ-563", status: "claimed", treeSection: "working"},
		{requestID: "REQ-564", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-563"}},
	})
	dependent := commitlessGraph.NodesByID["REQ-564"]
	if dependent.IsReady || !reflect.DeepEqual(dependent.UnmetDependencies, []string{"REQ-563"}) {
		t.Fatalf("commit-less claimed source satisfied dependency: %#v", dependent)
	}

	pendingGraph := buildFixtureGraph(t, []graphFixture{
		{requestID: "REQ-563", status: "pending", commit: "0123456789abcdef"},
		{requestID: "REQ-564", status: "pending", dependencyKey: "depends_on", dependencies: []string{"REQ-563"}},
	})
	if pendingGraph.NodesByID["REQ-564"].IsReady {
		t.Fatalf("pending remediation source unblocked dependent: %#v", pendingGraph.NodesByID["REQ-564"])
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

func TestFilenameOnlyCollisionMakesAbsentDependencyAmbiguous(t *testing.T) {
	repositoryRoot := t.TempDir()
	fixtures := map[string]string{
		"do-work/queue/REQ-021-first.md":     "---\nid: REQ-030\nstatus: completed\n---\nBody\n",
		"do-work/archive/REQ-021-second.md":  "---\nid: REQ-031\nstatus: completed\n---\nBody\n",
		"do-work/queue/REQ-040-dependent.md": "---\nid: REQ-040\nstatus: pending\ndepends_on: [REQ-021]\n---\nBody\n",
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
	dependent := graph.NodesByID["REQ-040"]
	if dependent == nil {
		t.Fatal("REQ-040 node is missing")
	}
	if dependent.IsReady || dependent.DependenciesSatisfied {
		t.Fatalf("REQ-040 was ready through a collided dependency: %#v", dependent)
	}
	if !reflect.DeepEqual(dependent.AmbiguousTargets, []string{"REQ-021"}) ||
		!reflect.DeepEqual(dependent.UnmetDependencies, []string{"REQ-021"}) ||
		len(dependent.MissingTargets) != 0 {
		t.Fatalf("REQ-040 collision evidence = %#v", dependent)
	}
	if dependent.DependencyDepth != -1 {
		t.Fatalf("REQ-040 dependency depth = %d, want unresolved -1", dependent.DependencyDepth)
	}
	if !reflect.DeepEqual(graph.WarningMessages, []string{"REQ-040 depends on ambiguous target REQ-021"}) {
		t.Fatalf("warnings = %#v", graph.WarningMessages)
	}
}

func TestSameIdFilenameFrontmatterCollisionKeepsDependencyAmbiguous(t *testing.T) {
	repositoryRoot := t.TempDir()
	fixtures := map[string]string{
		"do-work/archive/REQ-020-first.md":   "---\nid: REQ-021\nstatus: completed\n---\nBody\n",
		"do-work/archive/REQ-021-second.md":  "---\nid: REQ-021\nstatus: completed\n---\nBody\n",
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
	if !reflect.DeepEqual(dependent.AmbiguousTargets, []string{"REQ-021"}) ||
		!reflect.DeepEqual(dependent.UnmetDependencies, []string{"REQ-021"}) {
		t.Fatalf("REQ-030 ambiguity evidence = %#v", dependent)
	}
	if dependent.DependencyDepth != -1 {
		t.Fatalf("REQ-030 dependency depth = %d, want unresolved -1", dependent.DependencyDepth)
	}
	if target := graph.NodesByID["REQ-021"]; target == nil || !target.IsAmbiguous {
		t.Fatalf("REQ-021 target = %#v, want ambiguous", target)
	}
}
