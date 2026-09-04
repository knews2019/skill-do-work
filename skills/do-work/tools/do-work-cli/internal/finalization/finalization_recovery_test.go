package finalization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRecoverFinalizationDiscoversLegacyNoJournalTail(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	workingPath := "do-work/working/REQ-700-legacy.md"
	archivePath := "do-work/archive/REQ-700-legacy.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	requestBefore := "---\nid: REQ-700\ntitle: Legacy fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	requestAfter := "---\nid: REQ-700\ntitle: Legacy fixture\nstatus: completed\nclaimed_at: 2026-09-02T08:00:00Z\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	writeFinalizationFile(t, repositoryRoot, workingPath, requestBefore)
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-700: Legacy fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/queue/REQ-701-next.md", "---\nid: REQ-701\ntitle: Next fixture\nstatus: pending\n---\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")

	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repositoryRoot, filepath.FromSlash(workingPath)), filepath.Join(repositoryRoot, filepath.FromSlash(archivePath))); err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, archivePath, requestAfter)
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "unrelated\n")

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || len(recovered.Finalizations) != 1 {
		t.Fatalf("discovery result = %#v", recovered)
	}
	record := recovered.Finalizations[0]
	if record.RequestID != "REQ-700" || !record.Discovered || !record.Resumed {
		t.Fatalf("discovered record = %#v", record)
	}
	if record.PrimaryCommit == "" || record.MetadataCommit == "" || record.PrimaryCommit == record.MetadataCommit {
		t.Fatalf("commit identities = %#v", record)
	}
	archived := readFinalizationFile(t, repositoryRoot, archivePath)
	if !strings.Contains(archived, "commit: "+record.PrimaryCommit) {
		t.Fatalf("canonical provenance missing:\n%s", archived)
	}
	if got := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "show", record.PrimaryCommit+":implementation.txt")); got != "after" {
		t.Fatalf("primary implementation bytes = %q", got)
	}
	if got := readFinalizationFile(t, repositoryRoot, "notes.txt"); got != "unrelated\n" {
		t.Fatalf("unrelated unstaged file changed: %q", got)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-700")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(journalPath); !os.IsNotExist(err) {
		t.Fatalf("journal remains after discovery: %v", err)
	}

	second := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if second.Outcome != resultmodel.OutcomeSuccess || second.Finalization != nil || len(second.Finalizations) != 0 {
		t.Fatalf("second discovery is not idempotent: %#v", second)
	}
	selected := nextselection.Handlers()[nextselection.CommandNext](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if selected.Outcome != resultmodel.OutcomeSuccess || len(selected.Selected) != 1 || selected.Selected[0].RequestID != "REQ-701" {
		t.Fatalf("selection after recovery = %#v", selected)
	}
	claimed := requeststate.Handlers()["claim"](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"REQ-701", "--request-path", selected.Selected[0].RequestPath, "--provenance", "default", "--writer", "host:/repo", "--at", "2026-09-02T10:00:00Z",
	})
	if claimed.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("claim after recovery = %#v", claimed)
	}
}

func TestRecoverFinalizationAcceptsCoherentClaimOnlyTopology(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "do-work/queue/REQ-699.md", "---\nid: REQ-699\ntitle: Claim-only fixture\nstatus: pending\n---\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "claim-only fixture")
	claim := requeststate.Handlers()["claim"](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"REQ-699", "--request-path", "do-work/queue/REQ-699.md", "--provenance", "explicit-req", "--writer", "host:/repo", "--at", "2026-09-02T01:00:00Z",
	})
	if claim.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("claim fixture = %#v", claim)
	}

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 0 {
		t.Fatalf("claim-only topology was mistaken for finalization: %#v", result)
	}
}

func TestValidateManifestRequiresExplicitProvenanceMode(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "seed.txt", "seed\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	manifest := Manifest{
		RequestID: "REQ-701", RequestPath: "do-work/working/REQ-701.md", WriterLabel: "writer", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: strings.Repeat("a", 64),
		ExpectedCheckpointSHA256: strings.Repeat("b", 64), CommitPaths: []string{"seed.txt"}, CommitMessage: "fixture",
	}
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "provenance_mode") {
		t.Fatalf("missing provenance mode error = %v", err)
	}
	manifest.ProvenanceMode = ProvenancePrimaryCommit
	manifest.ImplementationHash = "abcdef0"
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "forbids") {
		t.Fatalf("primary hash error = %v", err)
	}
	manifest.ProvenanceMode = ProvenanceSuppliedCommit
	manifest.ImplementationHash = strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "--short", "HEAD"))
	if err := validateManifest(repositoryRoot, manifest); err != nil {
		t.Fatalf("valid supplied commit refused: %v", err)
	}
	manifest.ImplementationHash = "ABCDEF0"
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "lowercase") {
		t.Fatalf("uppercase supplied hash error = %v", err)
	}
}

func TestRecoverFinalizationDiscoveryRefusesExistingIndex(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "staged\n")
	runFinalizationGit(t, repositoryRoot, "add", "ordinary.txt")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) != 1 || result.Findings[0].Code != "FINALIZATION-DISCOVERY-FOREIGN-STAGED" {
		t.Fatalf("staged discovery result = %#v", result)
	}
	if got := readFinalizationFile(t, repositoryRoot, "ordinary.txt"); got != "staged\n" {
		t.Fatalf("staged path changed: %q", got)
	}
}

func TestMatchingHeadCommitUsesPreparedDiffIdentityAcrossLaterCommit(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "later.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	preparedHead, preparedDiff, err := preparedCommitIdentity(repositoryRoot, []string{"owned.txt"})
	if err != nil {
		t.Fatal(err)
	}
	runFinalizationGit(t, repositoryRoot, "add", "owned.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "primary")
	primaryCommit := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
	writeFinalizationFile(t, repositoryRoot, "later.txt", "after\n")
	runFinalizationGit(t, repositoryRoot, "add", "later.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "later")

	journal := &Journal{PreparedHead: preparedHead, PreparedDiffSHA256: preparedDiff, EffectiveCommitPaths: []string{"owned.txt"}}
	if matched, ok := matchingHeadCommit(repositoryRoot, journal); !ok || matched != primaryCommit {
		t.Fatalf("matched commit = %q, %t; want %q", matched, ok, primaryCommit)
	}
}

func TestRecoverFinalizationDiscoversCompleteSemanticLegacyTail(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	fixture := seedSemanticLegacyTail(t, repositoryRoot)

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if recovered.Outcome != resultmodel.OutcomeSuccess || len(recovered.Finalizations) != 1 {
		t.Fatalf("semantic discovery result = %#v", recovered)
	}
	record := recovered.Finalizations[0]
	if record.Phase != string(PhaseCleanupComplete) || record.PrimaryCommit == "" || record.MetadataCommit == "" || record.CreatedPrimaryCommit == "" || record.CreatedMetadataCommit == "" {
		t.Fatalf("terminal commit evidence = %#v", record)
	}
	if record.PrimaryCommit != record.CreatedPrimaryCommit || record.MetadataCommit != record.CreatedMetadataCommit {
		t.Fatalf("created and settled hashes diverged on first invocation: %#v", record)
	}
	if !reflect.DeepEqual(record.CommitPaths, fixture.ownedPaths) {
		t.Fatalf("commit paths = %v, want %v", record.CommitPaths, fixture.ownedPaths)
	}
	encoded, err := json.Marshal(resultmodel.NormalizeResult(recovered))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"\"primary_commit\"", "\"metadata_commit\"", "\"created_primary_commit\"", "\"created_metadata_commit\"", "\"blocked_paths\":[]", "\"reason_codes\":[]", "\"collection_argv\"", "\"next_argv\"", "\"verification_argv\""} {
		if !bytes.Contains(encoded, []byte(key)) {
			t.Fatalf("typed success JSON omits %s: %s", key, encoded)
		}
	}
	for _, path := range fixture.ownedPaths {
		if output := runFinalizationGit(t, repositoryRoot, "diff", "--name-only", record.PrimaryCommit+"^", record.PrimaryCommit, "--", path); strings.TrimSpace(output) == "" {
			t.Fatalf("primary commit omitted owned path %s", path)
		}
	}
	if got := readFinalizationFile(t, repositoryRoot, "notes.txt"); got != "unrelated\n" {
		t.Fatalf("unrelated ordinary dirt changed: %q", got)
	}
	if got := strings.Count(readFinalizationFile(t, repositoryRoot, "CHANGELOG.md"), "REQ-700"); got != 1 {
		t.Fatalf("release entry count = %d", got)
	}
	if got := strings.Count(readFinalizationFile(t, repositoryRoot, "do-work/calibration-log.tsv"), "REQ-700\t"); got != 1 {
		t.Fatalf("calibration row count = %d", got)
	}
	if second := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"}); second.Outcome != resultmodel.OutcomeSuccess || len(second.Finalizations) != 0 {
		t.Fatalf("semantic retry duplicated effects: %#v", second)
	}
}

// Pins the real failure: a Finder .DS_Store under do-work/ made discovery refuse
// FINALIZATION-DISCOVERY-AMBIGUOUS with no terminal request in sight, and again
// blocked commit safety when a real legacy tail was there to recover.
func TestRecoverFinalizationIgnoresFinderMetadataUnderDoWork(t *testing.T) {
	t.Run("no terminal request", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n")
		runFinalizationGit(t, repositoryRoot, "add", ".")
		runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
		writeFinalizationFile(t, repositoryRoot, "do-work/.DS_Store", "finder\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 0 {
			t.Fatalf("Finder metadata stopped discovery: %#v", result)
		}
	})
	t.Run("legacy tail recovers around it", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		seedSemanticLegacyTail(t, repositoryRoot)
		writeFinalizationFile(t, repositoryRoot, "do-work/.DS_Store", "finder\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 1 || result.Finalizations[0].Phase != string(PhaseCleanupComplete) {
			t.Fatalf("Finder metadata blocked legacy recovery: %#v", result)
		}
		if status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain", "--untracked-files=all", "--", "do-work/.DS_Store"); strings.TrimSpace(status) != "?? do-work/.DS_Store" {
			t.Fatalf("recovery touched or committed Finder metadata: %q", status)
		}
	})
}

func TestRecoverFinalizationRefusesForeignSharedHunkByteIdentically(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSemanticLegacyTail(t, repositoryRoot)
	writeFinalizationFile(t, repositoryRoot, "CHANGELOG.md", "# Changelog\n\n## 1.0.1 — REQ-700 semantic recovery\n\nOwned entry.\n\n## 1.0.0 — Seed plus FOREIGN-HUNK\n")
	before := readFinalizationFile(t, repositoryRoot, "CHANGELOG.md")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Finalizations) != 1 || result.Finalizations[0].ReasonCodes[0] != "FINALIZATION-DISCOVERY-AMBIGUOUS" {
		t.Fatalf("foreign shared hunk result = %#v", result)
	}
	if got := readFinalizationFile(t, repositoryRoot, "CHANGELOG.md"); got != before {
		t.Fatalf("refusal changed shared bytes:\n%s", got)
	}
}

func TestRecoverFinalizationRefusesMultiHunkProjectPathByteIdentically(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	archivePath := "do-work/archive/REQ-705-fixture.md"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: REQ-705\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `owned.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "one\ntwo\nthree\nfour\nfive\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "ONE\ntwo\nthree\nfour\nFIVE\n")
	before := readFinalizationFile(t, repositoryRoot, "owned.txt")
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", "owned.txt")
	if got := readFinalizationFile(t, repositoryRoot, "owned.txt"); got != before {
		t.Fatalf("multi-hunk refusal changed project bytes: %q", got)
	}
}

func TestRecoverFinalizationPreservesUnstagedProtectedAndDistinguishesStagedRefusals(t *testing.T) {
	t.Run("unstaged protected is preserved", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		seedSimpleDiscoveredTail(t, repositoryRoot, "REQ-710", "owned.txt")
		writeFinalizationFile(t, repositoryRoot, ".env.local", "secret bytes\n")
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("unstaged protected result = %#v", result)
		}
		if got := readFinalizationFile(t, repositoryRoot, ".env.local"); got != "secret bytes\n" {
			t.Fatalf("unstaged protected bytes changed: %q", got)
		}
	})

	t.Run("staged protected refuses distinctly", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		writeFinalizationFile(t, repositoryRoot, ".env.local", "seed\n")
		runFinalizationGit(t, repositoryRoot, "add", ".env.local")
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-PROTECTED-STAGED", ".env.local")
	})

	t.Run("staged ordinary refuses distinctly", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "seed\n")
		runFinalizationGit(t, repositoryRoot, "add", "ordinary.txt")
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-FOREIGN-STAGED", "ordinary.txt")
	})
}

func TestMatchingHeadCommitSearchesPastEarlierIndependentGroup(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "group-a.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "group-b.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	preparedHead := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
	writeFinalizationFile(t, repositoryRoot, "group-b.txt", "after\n")
	_, preparedDiff, err := preparedCommitIdentity(repositoryRoot, []string{"group-b.txt"})
	if err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, "group-a.txt", "after\n")
	runFinalizationGit(t, repositoryRoot, "add", "group-a.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "group a")
	runFinalizationGit(t, repositoryRoot, "add", "group-b.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "group b")
	targetCommit := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))

	journal := &Journal{PreparedHead: preparedHead, PreparedDiffSHA256: preparedDiff, EffectiveCommitPaths: []string{"group-b.txt"}}
	if matched, ok := matchingHeadCommit(repositoryRoot, journal); !ok || matched != targetCommit {
		t.Fatalf("matched commit = %q, %t; want later target %q", matched, ok, targetCommit)
	}
}

func TestDiscoveryRefusalRendersCompleteOrderedTypedEvidence(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "staged\n")
	runFinalizationGit(t, repositoryRoot, "add", "ordinary.txt")
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	normalized := resultmodel.NormalizeResult(result)
	encoded, err := json.Marshal(normalized)
	if err != nil {
		t.Fatal(err)
	}
	if normalized.Outcome != resultmodel.OutcomeRefused || len(normalized.Finalizations) != 1 {
		t.Fatalf("normalized refusal = %#v", normalized)
	}
	record := normalized.Finalizations[0]
	if record.Phase != string(PhaseDiscoveryRefused) || len(record.BlockedPaths) != 1 || len(record.ReasonCodes) != 1 || record.NextArgv == nil || record.VerificationArgv == nil || record.CollectionArgv == nil {
		t.Fatalf("typed refusal record = %#v; json=%s", record, encoded)
	}
}

func TestDiscoveryRefusalNamesInventoryAsTheResolvingVerb(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	result := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-AMBIGUOUS", "shared state has no exact owner", []string{"do-work/CHECKPOINT.md"})
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) != 1 {
		t.Fatalf("discovery refusal = %#v", result)
	}
	want := []string{"do-work-cli", "--format", "json", "uncommitted-inventory"}
	if !reflect.DeepEqual(result.Findings[0].NextArgv, want) {
		t.Fatalf("next argv = %#v, want inventory resolver %#v", result.Findings[0].NextArgv, want)
	}
	if result.Finalization == nil || !reflect.DeepEqual(result.Finalization.NextArgv, want) || len(result.Finalizations) != 1 || !reflect.DeepEqual(result.Finalizations[0].NextArgv, want) {
		t.Fatalf("typed finalization remedies = singular %#v ordered %#v, want %#v", result.Finalization, result.Finalizations, want)
	}
	verification := []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"}
	if !reflect.DeepEqual(result.Findings[0].VerificationArgv, verification) {
		t.Fatalf("verification argv = %#v, want %#v", result.Findings[0].VerificationArgv, verification)
	}
}

func TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce(t *testing.T) {
	phases := []Phase{PhasePrepared, PhaseLifecycleApplied, PhaseReleaseApplied, PhasePrimaryCommitted, PhaseMetadataCommitted, PhaseVerified, PhaseCleanupComplete}
	for _, interruptedPhase := range phases {
		t.Run(string(interruptedPhase), func(t *testing.T) {
			repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
			if interruptedPhase == PhasePrepared {
				if _, resumed, err := prepareJournal(context.Background(), repositoryRoot, manifestPath); err != nil || resumed {
					t.Fatalf("prepare interruption setup: resumed=%t err=%v", resumed, err)
				}
			} else {
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
					t.Fatalf("interruption at %s unexpectedly succeeded", interruptedPhase)
				}
			}
			recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
			if recovered.Outcome != resultmodel.OutcomeSuccess || len(recovered.Finalizations) != 1 || recovered.Finalizations[0].Phase != string(PhaseCleanupComplete) {
				t.Fatalf("recovery from %s = %#v", interruptedPhase, recovered)
			}
			if interruptedPhase == PhasePrimaryCommitted && recovered.Finalizations[0].CreatedPrimaryCommit != "" {
				t.Fatalf("settled pre-existing primary was reported as created this invocation: %#v", recovered.Finalizations[0])
			}
			if got := strings.Count(readFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-720.md"), "commit:"); got != 1 {
				t.Fatalf("phase %s provenance field count = %d", interruptedPhase, got)
			}
			if second := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil); second.Outcome != resultmodel.OutcomeSuccess || len(second.Finalizations) != 0 {
				t.Fatalf("phase %s second recovery = %#v", interruptedPhase, second)
			}
		})
	}
}

func TestFinalizeSupportsSuppliedWorktreeProvenanceWithoutMetadataCommit(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenanceSuppliedCommit)
	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 1 {
		t.Fatalf("supplied provenance result = %#v", result)
	}
	record := result.Finalizations[0]
	if record.Phase != string(PhaseCleanupComplete) || record.PrimaryCommit == "" || record.CreatedPrimaryCommit == "" || record.MetadataCommit != "" || record.CreatedMetadataCommit != "" {
		t.Fatalf("supplied provenance record = %#v", record)
	}
	archived := readFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-720.md")
	seedCommit := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", record.PrimaryCommit+"^"))
	if !strings.Contains(archived, "commit: "+seedCommit[:12]) {
		t.Fatalf("supplied implementation hash missing:\n%s", archived)
	}
}

func TestFinalizeFailureReportsTerminalCleanupWithoutProvenanceMetadata(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-725.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-725\ntitle: Failure fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\n---\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-725: Failure fixture — writer: host:/repo\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	manifest := Manifest{RequestID: "REQ-725", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "fail", FailureError: "fixture failure", FailureType: "code", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath), CommitPaths: []string{requestPath, "do-work/archive/REQ-725.md", checkpointPath}, CommitMessage: "[REQ-725] finalize failure fixture", ProvenanceMode: ProvenancePrimaryCommit}
	manifestPath := filepath.Join(t.TempDir(), "failure.json")
	contents, _ := json.Marshal(manifest)
	if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.TerminalStatus != "failed" || result.Finalization.Phase != string(PhaseCleanupComplete) || result.Finalization.MetadataCommit != "" {
		t.Fatalf("failed finalization result = %#v", result)
	}
}

func TestRecoverFinalizationProcessesTwoSafeGroupsInStableOrder(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedTwoSimpleDiscoveredTails(t, repositoryRoot)
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 2 {
		t.Fatalf("two-group recovery = %#v", result)
	}
	if got := []string{result.Finalizations[0].RequestID, result.Finalizations[1].RequestID}; !reflect.DeepEqual(got, []string{"REQ-730", "REQ-731"}) {
		t.Fatalf("group order = %v", got)
	}
	for _, record := range result.Finalizations {
		if record.Phase != string(PhaseCleanupComplete) || record.PrimaryCommit == "" || record.MetadataCommit == "" {
			t.Fatalf("incomplete group record = %#v", record)
		}
	}
}

func TestRecoverFinalizationRefusesCorruptJournalImage(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	if _, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath); err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, "do-work/working/REQ-720.md", "corrupt outside journal images\n")
	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || len(result.Finalization.ReasonCodes) < 1 || result.Finalization.ReasonCodes[0] != "FINALIZATION-LIFECYCLE-CONFLICT" {
		t.Fatalf("corrupt image result = %#v", result)
	}
}

func TestReadJournalRejectsCorruptStoredImage(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	journal.LifecyclePostimages[0].Bytes = append(journal.LifecyclePostimages[0].Bytes, []byte("corrupt")...)
	contents, _ := json.Marshal(journal)
	if err := os.WriteFile(journal.JournalPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readJournal(repositoryRoot, journal.JournalPath); err == nil || !strings.Contains(err.Error(), "integrity") {
		t.Fatalf("corrupt stored image error = %v", err)
	}
}

func TestVerifyFinalStateAllowsProvenanceAfterReleaseStamp(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	archivePath := "do-work/archive/REQ-726.md"
	beforeProvenance := "---\nid: REQ-726\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\nrelease_at: 2026-09-02T09:05:00Z\ncommit:\n---\n"
	afterProvenance := strings.Replace(beforeProvenance, "commit:\n", "commit: abcdef0\n", 1)
	writeFinalizationFile(t, repositoryRoot, archivePath, afterProvenance)
	writeFinalizationFile(t, repositoryRoot, "VERSION", "1.0.1\n")
	journal := &Journal{ArchivedPath: archivePath, PrimaryCommit: "abcdef0", Manifest: Manifest{RequestID: "REQ-726", Transition: "complete", TerminalStatus: "completed", ReleaseAt: "2026-09-02T09:05:00Z"}, ReleasePostimages: []FileImage{{Path: archivePath, Exists: true, Bytes: []byte(beforeProvenance), Mode: 0o644}, {Path: "VERSION", Exists: true, Bytes: []byte("1.0.1\n"), Mode: 0o644}}}
	if err := verifyFinalState(repositoryRoot, journal); err != nil {
		t.Fatalf("release verification rejected provenance-only archive delta: %v", err)
	}
}

type semanticLegacyFixture struct {
	ownedPaths []string
}

func seedSemanticLegacyTail(t *testing.T, repositoryRoot string) semanticLegacyFixture {
	t.Helper()
	workingPath := "do-work/working/REQ-700-semantic.md"
	archivePath := "do-work/archive/UR-700/REQ-700-semantic.md"
	requestBefore := "---\nid: REQ-700\ntitle: Semantic fixture\nstatus: claimed\nuser_request: UR-700\nroute: C\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\nestimate:\n  p50_active_minutes: 60\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	requestAfter := "---\nid: REQ-700\ntitle: Semantic fixture\nstatus: completed\nuser_request: UR-700\nroute: C\nclaimed_at: 2026-09-02T08:00:00Z\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\nestimate:\n  p50_active_minutes: 60\nrelease_at: 2026-09-02T09:05:00Z\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	checkpointBefore := "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-700: Semantic fixture — claimed now — writer: host:/repo\n  - branch: worktree-agent-REQ-700\n  - note: enriched continuation\n\n## Decisions\n\n- keep\n"
	checkpointAfter := "# Session Checkpoint\n\n## In Progress (interrupted)\n\n## Decisions\n\n- keep\n"
	writeFinalizationFile(t, repositoryRoot, workingPath, requestBefore)
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", checkpointBefore)
	writeFinalizationFile(t, repositoryRoot, "do-work/calibration-log.tsv", "req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/user-requests/UR-700/input.md", "---\nid: UR-700\nrequests: [REQ-700, REQ-699]\n---\nInput\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-699-prior.md", "---\nid: REQ-699\ntitle: Prior\nstatus: completed\nuser_request: UR-700\ncommit: abcdef0\n---\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "VERSION", "1.0.0\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/VERSION", "1.0.0\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/actions/version.md", "**Current version**: 1.0.0\n")
	writeFinalizationFile(t, repositoryRoot, "suite/modules.tsv", "source\tdestination\nskills/do-work\t.claude/skills/do-work\n")
	writeFinalizationFile(t, repositoryRoot, "package.json", "{\"name\":\"fixture\",\"version\":\"1.0.0\"}\n")
	writeFinalizationFile(t, repositoryRoot, "package-lock.json", "{\"name\":\"fixture\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"fixture\",\"version\":\"1.0.0\"}}}\n")
	writeFinalizationFile(t, repositoryRoot, "CHANGELOG.md", "# Changelog\n\n## 1.0.0 — Seed\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/CHANGELOG.md", "# Changelog\n\n## 1.0.0 — Seed\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")

	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "archive", "UR-700"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repositoryRoot, filepath.FromSlash(workingPath)), filepath.Join(repositoryRoot, filepath.FromSlash(archivePath))); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repositoryRoot, "do-work", "user-requests", "UR-700", "input.md"), filepath.Join(repositoryRoot, "do-work", "archive", "UR-700", "input.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repositoryRoot, "do-work", "archive", "REQ-699-prior.md"), filepath.Join(repositoryRoot, "do-work", "archive", "UR-700", "REQ-699-prior.md")); err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, archivePath, requestAfter)
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", checkpointAfter)
	writeFinalizationFile(t, repositoryRoot, "do-work/calibration-log.tsv", "req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\nREQ-700\tC\t60\t60\t2026-09-02T09:00:00Z\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/queue/REQ-701-follow-up.md", "---\nid: REQ-701\ntitle: Follow-up\nstatus: pending\nuser_request: UR-700\naddendum_to: REQ-700\n---\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, "VERSION", "1.0.1\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/VERSION", "1.0.1\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/actions/version.md", "**Current version**: 1.0.1\n")
	writeFinalizationFile(t, repositoryRoot, "package.json", "{\"name\":\"fixture\",\"version\":\"1.0.1\"}\n")
	writeFinalizationFile(t, repositoryRoot, "package-lock.json", "{\"name\":\"fixture\",\"version\":\"1.0.1\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"fixture\",\"version\":\"1.0.1\"}}}\n")
	changelog := "# Changelog\n\n## 1.0.1 — REQ-700 semantic recovery\n\nOwned entry.\n\n## 1.0.0 — Seed\n"
	writeFinalizationFile(t, repositoryRoot, "CHANGELOG.md", changelog)
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/CHANGELOG.md", changelog)
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "unrelated\n")
	return semanticLegacyFixture{ownedPaths: uniqueSorted([]string{
		"CHANGELOG.md", "VERSION", archivePath, "do-work/CHECKPOINT.md", "do-work/archive/REQ-699-prior.md",
		"do-work/archive/UR-700/REQ-699-prior.md", "do-work/archive/UR-700/input.md", "do-work/calibration-log.tsv",
		"do-work/queue/REQ-701-follow-up.md", "do-work/user-requests/UR-700/input.md", "implementation.txt",
		"package-lock.json", "package.json", "skills/do-work/CHANGELOG.md", "skills/do-work/VERSION", "skills/do-work/actions/version.md", workingPath,
	})}
}

func seedSimpleDiscoveredTail(t *testing.T, repositoryRoot, requestID, ownedPath string) {
	t.Helper()
	archivePath := "do-work/archive/" + requestID + "-fixture.md"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: "+requestID+"\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `"+ownedPath+"` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, ownedPath, "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, ownedPath, "after\n")
}

func seedTwoSimpleDiscoveredTails(t *testing.T, repositoryRoot string) {
	t.Helper()
	for index, requestID := range []string{"REQ-730", "REQ-731"} {
		ownedPath := strings.ToLower(requestID) + ".txt"
		completedAt := fmt.Sprintf("2026-09-02T09:0%d:00Z", index)
		writeFinalizationFile(t, repositoryRoot, "do-work/archive/"+requestID+"-fixture.md", "---\nid: "+requestID+"\ntitle: Fixture\nstatus: completed\ncompleted_at: "+completedAt+"\ncommit:\n---\n\n## Implementation Summary\n- `"+ownedPath+"` (modified)\n")
		writeFinalizationFile(t, repositoryRoot, ownedPath, "before\n")
	}
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "req-730.txt", "after 730\n")
	writeFinalizationFile(t, repositoryRoot, "req-731.txt", "after 731\n")
}

func seedPlannedFinalization(t *testing.T, provenanceMode string) (string, string) {
	t.Helper()
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-720.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-720\ntitle: Planned fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-720: Planned fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	seedCommit := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "--short=12", "HEAD"))
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	manifest := Manifest{
		RequestID: "REQ-720", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete", TerminalStatus: "completed",
		CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		CommitPaths: []string{requestPath, "do-work/archive/REQ-720.md", checkpointPath, "implementation.txt"}, CommitMessage: "[REQ-720] finalize planned fixture", ProvenanceMode: provenanceMode,
	}
	if provenanceMode == ProvenanceSuppliedCommit {
		manifest.ImplementationHash = seedCommit
	}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	contents, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot, manifestPath
}

func assertDiscoveryReason(t *testing.T, result resultmodel.CommandResult, code, path string) {
	t.Helper()
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Finalizations) != 1 {
		t.Fatalf("discovery refusal = %#v", result)
	}
	record := result.Finalizations[0]
	if !reflect.DeepEqual(record.ReasonCodes, []string{code}) || !reflect.DeepEqual(record.BlockedPaths, []string{path}) {
		t.Fatalf("refusal record = %#v", record)
	}
}

func TestRecoverFinalizationPreservesAndMatchesPrivateFileModeInLifecyclePostimages(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	requestPath := "do-work/working/REQ-720.md"
	reqAbs := filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
	if err := os.Chmod(reqAbs, 0o600); err != nil {
		t.Fatal(err)
	}

	journal, resumed, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
	if err != nil || resumed {
		t.Fatalf("prepare: resumed=%t, err=%v", resumed, err)
	}

	plan, err := lifecyclePlan(repositoryRoot, journal.Manifest, journal.LifecyclePreimages)
	if err != nil {
		t.Fatal(err)
	}
	applied := requeststate.ApplyPlan(context.Background(), plan)
	if applied.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("apply failed: %#v", applied)
	}

	archiveAbs := filepath.Join(repositoryRoot, "do-work", "archive", "REQ-720.md")
	info, err := os.Lstat(archiveAbs)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("archived file mode = %o, want 0600", info.Mode().Perm())
	}

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if recovered.Outcome != resultmodel.OutcomeSuccess || len(recovered.Finalizations) != 1 || recovered.Finalizations[0].Phase != string(PhaseCleanupComplete) {
		t.Fatalf("recovery failed: %#v", recovered)
	}
}

// REQ-515: one refused record is that REQ's exclusion, not the run's stop. The
// fixture is REQ-456's shape — a journal that cannot finish its tail while a
// second, unrelated journal is perfectly finishable. Before this REQ, recovery
// returned at the first refusal and the second REQ never ran.
func TestRecoverFinalizationSetsAsideRefusedRecordAndFinishesTheRest(t *testing.T) {
	repositoryRoot := seedTwoPlannedFinalizations(t)
	previousHook := afterFinalizationPhase
	afterFinalizationPhase = func(phase Phase) error {
		if phase != PhaseLifecycleApplied {
			return nil
		}
		if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "archive", "REQ-720.md")); err == nil {
			return context.Canceled
		}
		return nil
	}
	t.Cleanup(func() { afterFinalizationPhase = previousHook })

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	afterFinalizationPhase = previousHook
	if recovered.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("a REQ-scoped refusal stopped the whole recovery: %#v", recovered)
	}
	if len(recovered.Finalizations) != 2 {
		t.Fatalf("mixed recovery must report both records, got %#v", recovered.Finalizations)
	}
	setAside, settled := recovered.Finalizations[0], recovered.Finalizations[1]
	if setAside.RequestID != "REQ-720" || settled.RequestID != "REQ-721" {
		t.Fatalf("record order = %q, %q", setAside.RequestID, settled.RequestID)
	}
	if !containsFinalizationReason(setAside.ReasonCodes, setAsideReasonCode) {
		t.Fatalf("set-aside record must carry %s: %#v", setAsideReasonCode, setAside)
	}
	if !containsFinalizationReason(setAside.ReasonCodes, "FINALIZATION-JOURNAL-WRITE") {
		t.Fatalf("set-aside record must keep the refusal's own reason code: %#v", setAside)
	}
	if len(setAside.NextArgv) != 0 {
		t.Fatalf("a set-aside names no next verb of its own (REQ-514): %#v", setAside.NextArgv)
	}
	if settled.Phase != string(PhaseCleanupComplete) || len(settled.ReasonCodes) != 0 || len(settled.BlockedPaths) != 0 {
		t.Fatalf("the clean record did not settle: %#v", settled)
	}

	setAsideFinding := resultmodel.CommandFinding{}
	for _, finding := range recovered.Findings {
		if finding.Code == setAsideReasonCode {
			setAsideFinding = finding
		}
	}
	if len(setAsideFinding.AffectedIDs) != 1 || setAsideFinding.AffectedIDs[0] != "REQ-720" {
		t.Fatalf("the set-aside finding must name the one REQ it excludes: %#v", setAsideFinding)
	}
	if setAsideFinding.AutomationStopReason == "" || len(setAsideFinding.NextArgv) != 0 {
		t.Fatalf("set-aside finding shape = %#v", setAsideFinding)
	}

	selected := nextselection.Handlers()[nextselection.CommandNext](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if selected.Outcome != resultmodel.OutcomeSuccess || len(selected.Selected) != 1 || selected.Selected[0].RequestID != "REQ-722" {
		t.Fatalf("selection after a set-aside = %#v", selected)
	}
}

// REQ-515: dirt that no REQ owns still stops the run, and its finding names a
// verb other than the one that refused.
func TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedSimpleDiscoveredTail(t, repositoryRoot, "REQ-740", "req-740.txt")
	writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-999: foreign — writer: other:/repo\n")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("unowned shared dirt must stop the run: %#v", result)
	}
	for _, finding := range result.Findings {
		if len(finding.AffectedIDs) == 0 && len(finding.NextArgv) == 0 {
			t.Fatalf("a global stop must name a resolving verb: %#v", finding)
		}
	}
}

func containsFinalizationReason(reasonCodes []string, wanted string) bool {
	for _, code := range reasonCodes {
		if code == wanted {
			return true
		}
	}
	return false
}

// seedTwoPlannedFinalizations leaves two prepared journals plus one untouched
// pending REQ, so a run that sets the first journal aside can still be seen to
// finish the second and select the third.
func seedTwoPlannedFinalizations(t *testing.T) string {
	t.Helper()
	repositoryRoot := newFinalizationRepository(t)
	checkpointPath := "do-work/CHECKPOINT.md"
	checkpoint := "# Session Checkpoint\n\n## In Progress (interrupted)\n\n" +
		"- REQ-720: First fixture — claimed now — writer: host:/repo\n" +
		"- REQ-721: Second fixture — claimed now — writer: host:/repo\n"
	writeFinalizationFile(t, repositoryRoot, checkpointPath, checkpoint)
	writeFinalizationFile(t, repositoryRoot, "do-work/queue/REQ-722-next.md", "---\nid: REQ-722\ntitle: Next fixture\nstatus: pending\n---\n")
	for _, requestID := range []string{"REQ-720", "REQ-721"} {
		ownedPath := strings.ToLower(requestID) + ".txt"
		writeFinalizationFile(t, repositoryRoot, "do-work/working/"+requestID+".md",
			"---\nid: "+requestID+"\ntitle: Planned fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `"+ownedPath+"` (modified)\n")
		writeFinalizationFile(t, repositoryRoot, ownedPath, "before\n")
	}
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	for _, requestID := range []string{"REQ-720", "REQ-721"} {
		ownedPath := strings.ToLower(requestID) + ".txt"
		writeFinalizationFile(t, repositoryRoot, ownedPath, "after\n")
	}
	for _, requestID := range []string{"REQ-720", "REQ-721"} {
		requestPath := "do-work/working/" + requestID + ".md"
		manifest := Manifest{
			RequestID: requestID, RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete",
			TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z",
			ExpectedRequestSHA256:    digestFile(t, repositoryRoot, requestPath),
			ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
			CommitPaths: []string{requestPath, "do-work/archive/" + requestID + ".md", checkpointPath,
				strings.ToLower(requestID) + ".txt"},
			CommitMessage: "[" + requestID + "] finalize planned fixture", ProvenanceMode: ProvenancePrimaryCommit,
		}
		manifestPath := filepath.Join(t.TempDir(), requestID+"-manifest.json")
		contents, err := json.Marshal(manifest)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath); err != nil {
			t.Fatal(err)
		}
	}
	return repositoryRoot
}
