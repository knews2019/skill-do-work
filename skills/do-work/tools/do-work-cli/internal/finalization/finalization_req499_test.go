package finalization

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/publication"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRecoverFinalizationAssumeSoleReleaserRequiresDiscover(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	for _, arguments := range [][]string{{"--assume-sole-releaser"}, {"--discover", "--discover", "--assume-sole-releaser"}} {
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, arguments)
		if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "FINALIZATION-USAGE" {
			t.Fatalf("arguments %v result = %#v", arguments, result)
		}
	}
}

func TestRecoverFinalizationAssumeSoleReleaserAttributesOnlySharedMetadata(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSoleReleaserLegacyTail(t, repositoryRoot)
	checkpointBefore := readFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md")

	strict := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if strict.Outcome != resultmodel.OutcomeRefused || len(strict.Finalizations) != 1 || !reflect.DeepEqual(strict.Finalizations[0].ReasonCodes, []string{"FINALIZATION-DISCOVERY-AMBIGUOUS"}) || !containsFinalizationPath(strict.Finalizations[0].BlockedPaths, "do-work/CHECKPOINT.md") {
		t.Fatalf("strict shared-metadata refusal = %#v", strict)
	}
	if got := readFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md"); got != checkpointBefore {
		t.Fatalf("strict refusal changed checkpoint:\n%s", got)
	}

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover", "--assume-sole-releaser"})
	if recovered.Outcome != resultmodel.OutcomeSuccess || len(recovered.Finalizations) != 1 {
		t.Fatalf("sole-releaser recovery = %#v", recovered)
	}
	record := recovered.Finalizations[0]
	if record.Phase != string(PhaseCleanupComplete) || record.PrimaryCommit == "" || record.MetadataCommit == "" {
		t.Fatalf("sole-releaser terminal record = %#v", record)
	}
	attributedPaths := []string{}
	for _, finding := range recovered.Findings {
		if finding.Code == "FINALIZATION-SOLE-RELEASER-ATTRIBUTED" {
			attributedPaths = finding.AffectedPaths
		}
	}
	wantAttributed := []string{"do-work/CHECKPOINT.md", "do-work/calibration-log.tsv"}
	if !reflect.DeepEqual(attributedPaths, wantAttributed) {
		t.Fatalf("attributed paths = %v, want %v", attributedPaths, wantAttributed)
	}
	if got := readFinalizationFile(t, repositoryRoot, "notes.txt"); got != "unrelated\n" {
		t.Fatalf("unrelated project bytes changed: %q", got)
	}
	archived := readFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-740-fixture.md")
	if !strings.Contains(archived, "commit: "+record.PrimaryCommit) {
		t.Fatalf("recovered request lacks provenance:\n%s", archived)
	}
}

func containsFinalizationPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestRecoverFinalizationAssumeSoleReleaserRefusesMultipleTails(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedTwoSimpleDiscoveredTails(t, repositoryRoot)
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover", "--assume-sole-releaser"})
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Finalizations) != 1 || !reflect.DeepEqual(result.Finalizations[0].ReasonCodes, []string{"FINALIZATION-MULTIPLE-TAILS"}) {
		t.Fatalf("multiple-tail assertion result = %#v", result)
	}
}

func TestRecoverFinalizationAssumeSoleReleaserStillRefusesStagedProtectedPath(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSoleReleaserLegacyTail(t, repositoryRoot)
	writeFinalizationFile(t, repositoryRoot, ".env.local", "secret\n")
	runFinalizationGit(t, repositoryRoot, "add", ".env.local")
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover", "--assume-sole-releaser"})
	assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-PROTECTED-STAGED", ".env.local")
}

func TestRecoverFinalizationRefusesPartialConfiguredReleaseMirrors(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)
	writeFinalizationFile(t, repositoryRoot, "package-lock.json", "{\"name\":\"fixture\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"fixture\",\"version\":\"1.0.0\"}}}\n")
	before := readFinalizationFile(t, repositoryRoot, "package-lock.json")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", "package-lock.json")
	if got := readFinalizationFile(t, repositoryRoot, "package-lock.json"); got != before {
		t.Fatalf("partial mirror refusal changed lockfile: %q", got)
	}
}

func TestRecoverFinalizationRefusesForeignEditInTrackedFollowup(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	archivePath := "do-work/archive/REQ-750-fixture.md"
	followupPath := "do-work/queue/REQ-751-follow-up.md"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: REQ-750\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `owned.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, "---\nid: REQ-751\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-750\n---\n\nOriginal body.\n")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, "---\nid: REQ-751\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-750\n---\n\nOriginal body plus foreign edit.\n")
	before := readFinalizationFile(t, repositoryRoot, followupPath)

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", followupPath)
	if got := readFinalizationFile(t, repositoryRoot, followupPath); got != before {
		t.Fatalf("tracked follow-up refusal changed bytes: %q", got)
	}
}

func TestRecoverFinalizationAcceptsExactTrackedFollowupFold(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	archivePath := "do-work/archive/REQ-752-fixture.md"
	followupPath := "do-work/queue/REQ-753-follow-up.md"
	followupBefore := "---\nid: REQ-753\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-752\n---\n\nOriginal body.\n"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: REQ-752\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `owned.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, followupBefore)
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, followupBefore+"\n## Review Fold — REQ-752\n\nExact originating finding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-752 -->\n")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || !containsFinalizationPath(result.Finalization.CommitPaths, followupPath) {
		t.Fatalf("exact tracked follow-up fold = %#v", result)
	}
}

func TestFollowupPathProvesTheEntireSingleNamedFold(t *testing.T) {
	tests := []struct {
		name   string
		append string
		want   bool
	}{
		{name: "review fold", append: "\n## Review Fold — REQ-754\n\nExact originating finding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-754 -->\n", want: true},
		{name: "recovery fold", append: "\n## Recovery Fold — REQ-754\n\nExact originating recovery.\n\n<!-- do-work:finalization-followup-fold-end kind=recovery request=REQ-754 -->\n", want: true},
		{name: "foreign bytes before", append: "\nforeign preface\n\n## Review Fold — REQ-754\n\nExact originating finding.\n"},
		{name: "foreign section after", append: "\n## Review Fold — REQ-754\n\nExact originating finding.\n\n## Foreign Notes\n\nUnowned tail.\n"},
		{name: "second allowed section after", append: "\n## Review Fold — REQ-754\n\nExact originating finding.\n\n## Recovery Fold — REQ-754\n\nUnowned second fold.\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := newFinalizationRepository(t)
			path := "do-work/queue/REQ-755-follow-up.md"
			before := "---\nid: REQ-755\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-754\n---\n\nOriginal body.\n"
			writeFinalizationFile(t, repositoryRoot, path, before)
			runFinalizationGit(t, repositoryRoot, "add", ".")
			runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
			writeFinalizationFile(t, repositoryRoot, path, before+test.append)
			if got := followupPathProves(repositoryRoot, path, "REQ-754"); got != test.want {
				t.Fatalf("followupPathProves() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors(t *testing.T) {
	tests := []struct {
		name         string
		manifestPath string
		lockPath     string
		manifestOld  string
		manifestNew  string
		lockOld      string
		lockNew      string
		extraFiles   map[string]string
	}{
		{
			name: "npm", manifestPath: "consumer/packages/widget/package.json", lockPath: "consumer/package-lock.json",
			manifestOld: "{\"name\":\"widget\",\"version\":\"1.0.0\"}\n", manifestNew: "{\"name\":\"widget\",\"version\":\"1.0.1\"}\n",
			lockOld: "{\"name\":\"consumer\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\"},\"packages/widget\":{\"name\":\"widget\",\"version\":\"1.0.0\"}}}\n",
			lockNew: "{\"name\":\"consumer\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\"},\"packages/widget\":{\"name\":\"widget\",\"version\":\"1.0.1\"}}}\n",
			extraFiles: map[string]string{
				"package.json":          "{\"name\":\"root\",\"version\":\"1.0.0\",\"private\":true,\"workspaces\":[\"consumer\"]}\n",
				"package-lock.json":     "{\"name\":\"root\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"root\",\"version\":\"1.0.0\"}}}\n",
				"consumer/package.json": "{\"name\":\"consumer\",\"private\":true,\"workspaces\":[\"packages/*\"]}\n",
			},
		},
		{
			name: "cargo", manifestPath: "consumer-rust/crates/widget/Cargo.toml", lockPath: "consumer-rust/Cargo.lock",
			manifestOld: "[package]\nname = \"widget\"\nversion = \"1.0.0\"\n", manifestNew: "[package]\nname = \"widget\"\nversion = \"1.0.1\"\n",
			lockOld: "[[package]]\nname = \"widget\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"helper\"\nversion = \"2.0.0\"\n",
			lockNew: "[[package]]\nname = \"widget\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"helper\"\nversion = \"2.0.0\"\n",
			extraFiles: map[string]string{
				"Cargo.toml":               "[workspace]\nmembers = [\"consumer-rust\"]\n",
				"consumer-rust/Cargo.toml": "[workspace]\nmembers = [\"crates/*\"]\n",
			},
		},
		{
			name: "uv", manifestPath: "consumer-python/packages/widget/pyproject.toml", lockPath: "consumer-python/uv.lock",
			manifestOld: "[project]\nname = \"widget\"\nversion = \"1.0.0\"\n", manifestNew: "[project]\nname = \"widget\"\nversion = \"1.0.1\"\n",
			lockOld: "[[package]]\nname = \"widget\"\nversion = \"1.0.0\"\nsource = { editable = \"packages/widget\" }\n\n[[package]]\nname = \"helper\"\nversion = \"2.0.0\"\nsource = { editable = \"packages/helper\" }\n",
			lockNew: "[[package]]\nname = \"widget\"\nversion = \"1.0.1\"\nsource = { editable = \"packages/widget\" }\n\n[[package]]\nname = \"helper\"\nversion = \"2.0.0\"\nsource = { editable = \"packages/helper\" }\n",
			extraFiles: map[string]string{
				"pyproject.toml":                 "[tool.uv.workspace]\nmembers = [\"consumer-python\"]\n",
				"consumer-python/pyproject.toml": "[tool.uv.workspace]\nmembers = [\"packages/*\"]\n",
			},
		},
	}
	for _, test := range tests {
		for _, updateLock := range []bool{false, true} {
			name := "stale"
			if updateLock {
				name = "updated"
			}
			t.Run(test.name+"/"+name, func(t *testing.T) {
				repositoryRoot := newFinalizationRepository(t)
				seedSemanticLegacyTail(t, repositoryRoot)
				writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.manifestOld)
				writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockOld)
				tracked := []string{test.manifestPath, test.lockPath}
				for path, contents := range test.extraFiles {
					writeFinalizationFile(t, repositoryRoot, path, contents)
					tracked = append(tracked, path)
				}
				runFinalizationGit(t, repositoryRoot, append([]string{"add", "--"}, tracked...)...)
				runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed workspace")
				writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.manifestNew)
				if updateLock {
					writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockNew)
				}

				result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
				if !updateLock {
					assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", test.lockPath)
					return
				}
				if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) || !containsFinalizationPath(result.Finalization.CommitPaths, test.lockPath) {
					t.Fatalf("updated workspace recovery = %#v", result)
				}
			})
		}
	}
}

func TestRecoverFinalizationReleaseEnumerationFailureIsTypedAndFailClosed(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)
	previous := enumerateTrackedReleasePaths
	enumerateTrackedReleasePaths = func(string) ([]string, error) { return nil, os.ErrPermission }
	t.Cleanup(func() { enumerateTrackedReleasePaths = previous })

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || !reflect.DeepEqual(result.Finalization.ReasonCodes, []string{"FINALIZATION-DISCOVERY-RELEASE-ENUMERATION"}) {
		t.Fatalf("release enumeration refusal = %#v", result)
	}
	if status := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "status", "--short")); status == "" {
		t.Fatal("enumeration refusal unexpectedly committed or removed the legacy tail")
	}
}

func TestRecoverFinalizationRefusesReleaseMetadataWithoutProjectOwnership(t *testing.T) {
	for _, path := range []string{"third_party/do-work/VERSION", "dist/skills/do-work/VERSION", "cache-tree-9f3a/rebuilt/VERSION"} {
		t.Run(path, func(t *testing.T) {
			repositoryRoot := newFinalizationRepository(t)
			seedSemanticLegacyTail(t, repositoryRoot)
			writeFinalizationFile(t, repositoryRoot, path, "1.0.0\n")
			runFinalizationGit(t, repositoryRoot, "add", "--", path)
			runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed unowned release metadata")
			writeFinalizationFile(t, repositoryRoot, path, "1.0.1\n")

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || !reflect.DeepEqual(result.Finalization.ReasonCodes, []string{"FINALIZATION-DISCOVERY-RELEASE-OWNERSHIP"}) || !containsFinalizationPath(result.Finalization.BlockedPaths, path) {
				t.Fatalf("unowned release metadata refusal = %#v", result)
			}
			if got := readFinalizationFile(t, repositoryRoot, path); got != "1.0.1\n" {
				t.Fatalf("ownership refusal changed %s: %q", path, got)
			}
		})
	}
}

func TestRecoverFinalizationRefusesCompleteInstalledSuiteWithoutMaintainerTopology(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)
	runFinalizationGit(t, repositoryRoot, "rm", "-q", "--", "suite/modules.tsv")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "remove suite maintainer topology")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	wantPaths := []string{"skills/do-work/CHANGELOG.md", "skills/do-work/VERSION", "skills/do-work/actions/version.md"}
	if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || !reflect.DeepEqual(result.Finalization.ReasonCodes, []string{"FINALIZATION-DISCOVERY-RELEASE-OWNERSHIP"}) || !reflect.DeepEqual(result.Finalization.BlockedPaths, wantPaths) {
		t.Fatalf("complete installed-suite refusal = %#v", result)
	}
}

func TestRecoverFinalizationAcceptsSuiteMirrorsWithMaintainerTopology(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("declared suite topology recovery = %#v", result)
	}
}

func TestAffirmativeReleaseOwnershipRequiresRootedWorkspaceChain(t *testing.T) {
	tests := []struct {
		name             string
		files            map[string]string
		ownedPaths       []string
		unownedPaths     []string
		ownedManifests   []string
		unownedManifests []string
	}{
		{
			name: "nested npm workspace cannot own itself",
			files: map[string]string{
				"third_party/dependency/package.json": "{\"name\":\"dependency\",\"workspaces\":[\"packages/*\"]}\n",
				"third_party/dependency/VERSION":      "1.0.0\n",
			},
			unownedPaths: []string{"third_party/dependency/VERSION"}, unownedManifests: []string{"third_party/dependency/package.json"},
		},
		{
			name: "nested cargo workspace cannot own itself",
			files: map[string]string{
				"third_party/dependency/Cargo.toml": "[workspace]\nmembers = [\"crates/*\"]\n",
				"third_party/dependency/VERSION":    "1.0.0\n",
			},
			unownedPaths: []string{"third_party/dependency/VERSION"}, unownedManifests: []string{"third_party/dependency/Cargo.toml"},
		},
		{
			name: "nested uv workspace cannot own itself",
			files: map[string]string{
				"third_party/dependency/pyproject.toml": "[tool.uv.workspace]\nmembers = [\"packages/*\"]\n",
				"third_party/dependency/VERSION":        "1.0.0\n",
			},
			unownedPaths: []string{"third_party/dependency/VERSION"}, unownedManifests: []string{"third_party/dependency/pyproject.toml"},
		},
		{
			name: "repository-root chain owns nested workspace member",
			files: map[string]string{
				"package.json":                                "{\"name\":\"root\",\"workspaces\":[\"packages/*\"]}\n",
				"packages/owned/package.json":                 "{\"name\":\"owned\",\"workspaces\":[\"children/*\"]}\n",
				"packages/owned/children/widget/package.json": "{\"name\":\"widget\",\"version\":\"1.0.0\"}\n",
				"packages/owned/children/widget/VERSION":      "1.0.0\n",
			},
			ownedPaths: []string{"packages/owned/children/widget/VERSION"}, ownedManifests: []string{"package.json", "packages/owned/package.json", "packages/owned/children/widget/package.json"},
		},
		{
			name: "unrooted chain remains unowned",
			files: map[string]string{
				"package.json":                                        "{\"name\":\"root\",\"workspaces\":[\"packages/*\"]}\n",
				"third_party/dependency/package.json":                 "{\"name\":\"dependency\",\"workspaces\":[\"children/*\"]}\n",
				"third_party/dependency/children/widget/package.json": "{\"name\":\"widget\",\"version\":\"1.0.0\"}\n",
				"third_party/dependency/children/widget/VERSION":      "1.0.0\n",
			},
			unownedPaths: []string{"third_party/dependency/children/widget/VERSION"}, unownedManifests: []string{"third_party/dependency/package.json", "third_party/dependency/children/widget/package.json"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := newFinalizationRepository(t)
			for path, contents := range test.files {
				writeFinalizationFile(t, repositoryRoot, path, contents)
			}
			runFinalizationGit(t, repositoryRoot, "add", ".")
			runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed ownership graph")
			tracked, err := enumerateTrackedReleasePaths(repositoryRoot)
			if err != nil {
				t.Fatal(err)
			}
			trackedSet := map[string]bool{}
			for _, path := range tracked {
				trackedSet[filepath.ToSlash(filepath.Clean(path))] = true
			}
			ownedPaths, ownedManifests, err := affirmativeReleaseOwnership(repositoryRoot, tracked, trackedSet)
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range test.ownedPaths {
				if !ownedPaths[path] {
					t.Errorf("owned path %s was not classified", path)
				}
			}
			for _, path := range test.unownedPaths {
				if ownedPaths[path] {
					t.Errorf("unowned path %s was classified", path)
				}
			}
			for _, path := range test.ownedManifests {
				if !ownedManifests[path] {
					t.Errorf("owned manifest %s was not classified", path)
				}
			}
			for _, path := range test.unownedManifests {
				if ownedManifests[path] {
					t.Errorf("unowned manifest %s was classified", path)
				}
			}
		})
	}
}

func TestPublicRecoverFinalizationMovesURThenAllowsRealClaim(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)
	nextPath := "do-work/queue/REQ-799-next.md"
	writeFinalizationFile(t, repositoryRoot, nextPath, "---\nid: REQ-799\ntitle: Next request\nstatus: pending\ncreated_at: 2026-09-02T09:10:00Z\n---\n")
	runFinalizationGit(t, repositoryRoot, "add", "--", nextPath)
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed next request")

	recovered := runPublicFinalizationCommand(t, repositoryRoot, "recover-finalization", "--discover", "--assume-sole-releaser")
	if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || recovered.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("public recovery = %#v", recovered)
	}
	for _, path := range []string{"do-work/archive/UR-700/input.md", "do-work/user-requests/UR-700/input.md"} {
		if !containsFinalizationPath(recovered.Finalization.CommitPaths, path) {
			t.Fatalf("public recovery commit paths %v omit %s", recovered.Finalization.CommitPaths, path)
		}
	}
	claimed := runPublicFinalizationCommand(t, repositoryRoot, "claim", "REQ-799", "--commit", "--writer", "fixture:/repo", "--at", "2026-09-02T09:15:00Z")
	if claimed.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("public claim after recovery = %#v", claimed)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "REQ-799-next.md")); err != nil {
		t.Fatalf("claimed request was not moved to working: %v", err)
	}
}

func runPublicFinalizationCommand(t *testing.T, repositoryRoot string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	commandArguments := append([]string{"run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json"}, arguments...)
	command := exec.Command("go", commandArguments...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("public command %v: %v\n%s", arguments, err, output)
	}
	var result resultmodel.CommandResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode public command %v: %v\n%s", arguments, err, output)
	}
	return result
}

func TestRecoverFinalizationResumesAfterRealPreCommitHookFailure(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	hookPath := filepath.Join(repositoryRoot, ".git", "hooks", "pre-commit")
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\nexit 1\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	first := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if first.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("failing commit hook unexpectedly succeeded: %#v", first)
	}
	if err := os.Remove(hookPath); err != nil {
		t.Fatal(err)
	}
	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || recovered.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("hook recovery = %#v", recovered)
	}
}

func TestRecoverFinalizationReleaseReplacementPreservesModeZeroSentinel(t *testing.T) {
	for _, interruptedPhase := range []Phase{PhaseLifecycleApplied, PhaseReleaseApplied, PhasePrimaryCommitted, PhaseMetadataCommitted, PhaseVerified} {
		t.Run(string(interruptedPhase), func(t *testing.T) {
			repositoryRoot, manifestPath := seedPlannedReleaseFinalization(t)
			previousHook := afterFinalizationPhase
			afterFinalizationPhase = func(phase Phase) error {
				if phase == interruptedPhase {
					return context.Canceled
				}
				return nil
			}
			first := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
			afterFinalizationPhase = previousHook
			t.Cleanup(func() { afterFinalizationPhase = previousHook })
			if first.Outcome == resultmodel.OutcomeSuccess {
				t.Fatalf("release interruption unexpectedly succeeded: %#v", first)
			}
			recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
			if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || recovered.Finalization.Phase != string(PhaseCleanupComplete) {
				t.Fatalf("release recovery = %#v", recovered)
			}
			info, err := os.Stat(filepath.Join(repositoryRoot, "VERSION"))
			if err != nil {
				t.Fatal(err)
			}
			if got := info.Mode().Perm(); got != 0o600 {
				t.Fatalf("VERSION mode = %#o, want preserved 0600", got)
			}
		})
	}
}

func TestRecoverFinalizationAlreadyGreenNoReleaseManifest(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-765.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-765\ntitle: Already green\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-765: fixture — writer: host:/repo\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	manifest := Manifest{RequestID: "REQ-765", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete", TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath), CommitPaths: []string{requestPath, "do-work/archive/REQ-765.md", checkpointPath}, CommitMessage: "[REQ-765] finalize already-green fixture", ProvenanceMode: ProvenancePrimaryCommit}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	manifestBytes, _ := json.Marshal(manifest)
	if err := os.WriteFile(manifestPath, manifestBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	journal, resumed, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
	if err != nil || resumed {
		t.Fatalf("already-green prepare: resumed=%t err=%v", resumed, err)
	}
	if len(journal.ReleasePostimages) != 0 {
		t.Fatalf("already-green manifest unexpectedly planned a release: %#v", journal.ReleasePostimages)
	}
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("already-green recovery = %#v", result)
	}
}

func TestPublicationPostimagesKeepExplicitReplacementMode(t *testing.T) {
	plan := publication.PublicationPlan{Mutations: []publication.PlannedMutation{{Kind: publication.MutationReplace, Path: "VERSION", Contents: []byte("2.0.0\n"), Mode: 0o640}}}
	images := publicationPostimages(plan)
	if len(images) != 1 || images[0].Mode != 0o640 {
		t.Fatalf("explicit replacement mode = %#v", images)
	}
}

func TestReleaseVersionRecognizesOwnedCargoAndUVLockEntries(t *testing.T) {
	tests := []struct {
		path   string
		before string
		after  string
	}{
		{path: "Cargo.lock", before: "[[package]]\nname = \"fixture\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.0\"\nsource = \"registry+x\"\n", after: "[[package]]\nname = \"fixture\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.0\"\nsource = \"registry+x\"\n"},
		{path: "uv.lock", before: "[[package]]\nname = \"fixture\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.0\"\nsource = { registry = \"https://example.invalid\" }\n", after: "[[package]]\nname = \"fixture\"\nversion = \"1.0.1\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.0\"\nsource = { registry = \"https://example.invalid\" }\n"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			if version, ok := releaseVersion(test.path, []byte(test.before)); !ok || version != "1.0.0" {
				t.Fatalf("owned lock version = %q, %t", version, ok)
			}
			if !semanticVersionReplacement(test.path, []byte(test.before), []byte(test.after), "1.0.0", "1.0.1") {
				t.Fatal("owned lock replacement was not exact")
			}
		})
	}
}

func TestReleaseVersionReadsOnlyProjectTOMLSections(t *testing.T) {
	for path, contents := range map[string]string{
		"Cargo.toml":     "[workspace.package]\nversion = \"9.9.9\"\n\n[package]\nname = \"fixture\"\nversion = \"1.0.0\"\n",
		"pyproject.toml": "[tool.example]\nversion = \"9.9.9\"\n\n[project]\nname = \"fixture\"\nversion = \"1.0.0\"\n",
	} {
		if version, ok := releaseVersion(path, []byte(contents)); !ok || version != "1.0.0" {
			t.Fatalf("%s project version = %q, %t", path, version, ok)
		}
	}
}

func seedSoleReleaserLegacyTail(t *testing.T, repositoryRoot string) {
	t.Helper()
	archivePath := "do-work/archive/REQ-740-fixture.md"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: REQ-740\ntitle: Fixture\nstatus: completed\nroute: C\nclaimed_at: 2026-09-02T08:00:00Z\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\nestimate:\n  p50_active_minutes: 60\n---\n\n## Implementation Summary\n- `owned.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-740: fixture\n\n## Session Notes\nold\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/calibration-log.tsv", "req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\n")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "unrelated\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Checkpoint\n\n## In Progress (interrupted)\n\n## Session Notes\nnew unrelated note\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/calibration-log.tsv", "req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\nREQ-740\tC\t60\t61\t2026-09-02T09:00:00Z\n")
}

func seedPlannedReleaseFinalization(t *testing.T) (string, string) {
	t.Helper()
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-760.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-760\ntitle: Planned release\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-760: fixture — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "VERSION", "1.0.0\n")
	writeFinalizationFile(t, repositoryRoot, "CHANGELOG.md", "# Changelog\n\n## 1.0.0 — Seed\n")
	if err := os.Chmod(filepath.Join(repositoryRoot, "VERSION"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	payloadRoot := t.TempDir()
	writeFinalizationFile(t, payloadRoot, "version-old", "1.0.0\n")
	writeFinalizationFile(t, payloadRoot, "version-new", "1.0.1\n")
	writeFinalizationFile(t, payloadRoot, "changelog-old", "# Changelog\n\n## 1.0.0 — Seed\n")
	writeFinalizationFile(t, payloadRoot, "changelog-new", "# Changelog\n\n## 1.0.1 — REQ-760 Planned Release\n\nDelivered.\n\n## 1.0.0 — Seed\n")
	releaseManifest := publication.Manifest{Operation: publication.OperationRelease, Release: &publication.ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1",
		ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:             []publication.ReleaseTarget{{Path: "VERSION", ExpectedPayload: publication.PayloadFile{SourcePath: filepath.Join(payloadRoot, "version-old")}, NewPayload: publication.PayloadFile{SourcePath: filepath.Join(payloadRoot, "version-new")}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs:          []publication.ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: publication.PayloadFile{SourcePath: filepath.Join(payloadRoot, "changelog-old")}, NewPayload: publication.PayloadFile{SourcePath: filepath.Join(payloadRoot, "changelog-new")}, InsertionAnchor: "## 1.0.0", EntryKey: "1.0.1", EntryTitle: "Planned Release"}},
	}}
	releasePath := filepath.Join(t.TempDir(), "release.json")
	releaseBytes, _ := json.Marshal(releaseManifest)
	if err := os.WriteFile(releasePath, releaseBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := Manifest{RequestID: "REQ-760", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete", TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath), CommitPaths: []string{requestPath, "do-work/archive/REQ-760.md", checkpointPath, "VERSION", "CHANGELOG.md"}, CommitMessage: "[REQ-760] finalize planned release", ProvenanceMode: ProvenancePrimaryCommit, ReleaseManifestPath: releasePath, ReleaseAt: "2026-09-02T09:05:00Z"}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	manifestBytes, _ := json.Marshal(manifest)
	if err := os.WriteFile(manifestPath, manifestBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot, manifestPath
}
