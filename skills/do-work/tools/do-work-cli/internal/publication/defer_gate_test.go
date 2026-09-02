package publication

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDeferGateCreatePublishesOneAtomicDependencyLifecycle(t *testing.T) {
	root := newDeferGateRepository(t, true)
	manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal != nil {
		t.Fatalf("plan refusal = %#v", plan.Refusal)
	}
	result := ApplyPlan(t.Context(), plan, false, false)
	if result.Outcome != resultmodel.OutcomeSuccess || result.GateDeferral == nil || result.GateDeferral.RepairOutcome != "created" {
		t.Fatalf("result = %#v", result)
	}
	parentPath := filepath.Join(root, "do-work/queue/REQ-101-parent.md")
	parentBytes, err := os.ReadFile(parentPath)
	if err != nil {
		t.Fatal(err)
	}
	parentDocument, _ := requestmodel.ParseDocument(parentBytes)
	parentRecord := parentDocument.TypedRecord()
	if parentRecord.RequestStatus != "pending" || parentRecord.GateDeferredValue != "true" || !reflect.DeepEqual(parentRecord.DependsOn, []string{"REQ-901"}) || parentRecord.ClaimedAt != "" {
		t.Fatalf("parent record = %#v", parentRecord)
	}
	if !bytes.Contains(parentBytes, []byte("## Repository Gate Deferral")) || !bytes.Contains(parentBytes, []byte("- **Direct exit status:** 17")) {
		t.Fatalf("parent history missing:\n%s", parentBytes)
	}
	if _, err := os.Lstat(filepath.Join(root, "do-work/working/REQ-101-parent.md")); !os.IsNotExist(err) {
		t.Fatalf("working parent survived: %v", err)
	}
	repairBytes, err := os.ReadFile(filepath.Join(root, "do-work/queue/REQ-901-repair-repository-gate.md"))
	if err != nil {
		t.Fatal(err)
	}
	repairDocument, _ := requestmodel.ParseDocument(repairBytes)
	repairRecord := repairDocument.TypedRecord()
	if repairRecord.RequestStatus != "pending" || repairRecord.RepositoryGateRepairValue != "true" || repairRecord.FieldEvidenceByName["sweep"].ScalarValue != "true" || repairRecord.UserRequestID != "UR-095" || repairRecord.AddendumTo != "" || !reflect.DeepEqual(repairRecord.RelatedIDs, []string{"REQ-101"}) {
		t.Fatalf("repair record = %#v", repairRecord)
	}
	if !bytes.Contains(repairBytes, []byte("## Instances\n\n- [ ] repository gate: sha256:gate-red affecting REQ-101 (found by REQ-101 / UR-095)")) {
		t.Fatalf("repair is not a canonical sweep projection:\n%s", repairBytes)
	}
	checkpointBytes, err := os.ReadFile(filepath.Join(root, "do-work/CHECKPOINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(checkpointBytes, []byte("REQ-101")) || bytes.Contains(checkpointBytes, []byte("parent detail")) || !bytes.Contains(checkpointBytes, []byte("REQ-999")) || !bytes.Contains(checkpointBytes, []byte("foreign detail")) {
		t.Fatalf("checkpoint ownership rewrite damaged entries:\n%s", checkpointBytes)
	}
}

func TestDeferGateFoldsSharedFingerprintWithoutOverwritingPriorParent(t *testing.T) {
	root := newDeferGateRepository(t, true)
	first := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
	firstResult := ApplyPlan(t.Context(), BuildDeferGatePlan(root, first), false, false)
	if firstResult.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("first deferral = %#v", firstResult)
	}
	claimDeferParent(t, root, "REQ-102", "do-work/working/REQ-102-second.md", "Second parent")
	repairPath := "do-work/queue/REQ-901-repair-repository-gate.md"
	second := deferGateManifest(root, "REQ-102", "REQ-901", "do-work/working/REQ-102-second.md", &PayloadFile{SourcePath: repairPath})
	plan := BuildDeferGatePlan(root, second)
	if plan.Refusal != nil {
		t.Fatalf("fold refusal = %#v", plan.Refusal)
	}
	result := ApplyPlan(t.Context(), plan, false, false)
	if result.Outcome != resultmodel.OutcomeSuccess || result.GateDeferral == nil || result.GateDeferral.RepairOutcome != "folded" {
		t.Fatalf("fold result = %#v", result)
	}
	repairBytes, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(repairPath)))
	if err != nil {
		t.Fatal(err)
	}
	document, _ := requestmodel.ParseDocument(repairBytes)
	if !reflect.DeepEqual(document.TypedRecord().RelatedIDs, []string{"REQ-101", "REQ-102"}) || bytes.Count(repairBytes, []byte("**Parent:** REQ-101")) != 1 || bytes.Count(repairBytes, []byte("**Parent:** REQ-102")) != 1 || bytes.Count(repairBytes, []byte("- [ ] repository gate:")) != 2 {
		t.Fatalf("folded repair lost parent evidence:\n%s", repairBytes)
	}
}

func TestDeferGateClassifiesTrackedDirtyTrackedCleanAndUntrackedPreimagesIndependently(t *testing.T) {
	for name, prepare := range map[string]func(*testing.T, string){
		"dirty parent dirty checkpoint": func(*testing.T, string) {},
		"dirty parent clean checkpoint": func(t *testing.T, root string) {
			runGitFixture(t, root, "add", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "checkpoint claim", "--", "do-work/CHECKPOINT.md")
		},
		"clean parent clean checkpoint": func(t *testing.T, root string) {
			runGitFixture(t, root, "add", "do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "claimed inputs")
		},
		"untracked parent untracked checkpoint": func(t *testing.T, root string) {
			runGitFixture(t, root, "rm", "--cached", "do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "consumer keeps lifecycle files untracked")
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := newDeferGateRepository(t, false)
			prepare(t, root)
			plan := BuildDeferGatePlan(root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil))
			if plan.Refusal != nil {
				t.Fatalf("plan refusal = %#v", plan.Refusal)
			}
			dirty := strings.Join(plan.ExistingDirtyTargetPaths, " ")
			untracked := strings.Join(plan.ExistingUntrackedTargetPaths, " ")
			switch name {
			case "dirty parent dirty checkpoint":
				if !strings.Contains(dirty, "REQ-101-parent.md") || !strings.Contains(dirty, "CHECKPOINT.md") || untracked != "" {
					t.Fatalf("classification dirty=%q untracked=%q", dirty, untracked)
				}
			case "dirty parent clean checkpoint":
				if !strings.Contains(dirty, "REQ-101-parent.md") || strings.Contains(dirty, "CHECKPOINT.md") || untracked != "" {
					t.Fatalf("classification dirty=%q untracked=%q", dirty, untracked)
				}
			case "clean parent clean checkpoint":
				if dirty != "" || untracked != "" {
					t.Fatalf("classification dirty=%q untracked=%q", dirty, untracked)
				}
			case "untracked parent untracked checkpoint":
				if dirty != "" || !strings.Contains(untracked, "REQ-101-parent.md") || !strings.Contains(untracked, "CHECKPOINT.md") {
					t.Fatalf("classification dirty=%q untracked=%q", dirty, untracked)
				}
			}
			result := ApplyPlan(t.Context(), plan, false, false)
			if result.Outcome != resultmodel.OutcomeSuccess {
				t.Fatalf("apply = %#v", result)
			}
		})
	}
}

func TestDeferGateRollsBackEveryMutationPositionToExactDirtyPreimages(t *testing.T) {
	for failureIndex := 0; failureIndex < 4; failureIndex++ {
		t.Run(string(rune('0'+failureIndex)), func(t *testing.T) {
			root := newDeferGateRepository(t, false)
			parentPath := filepath.Join(root, "do-work/working/REQ-101-parent.md")
			checkpointPath := filepath.Join(root, "do-work/CHECKPOINT.md")
			parentBefore, _ := os.ReadFile(parentPath)
			checkpointBefore, _ := os.ReadFile(checkpointPath)
			plan := BuildDeferGatePlan(root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil))
			previous := afterPublicationMutation
			afterPublicationMutation = func(index int, _ PlannedMutation) error {
				if index == failureIndex {
					return errors.New("injected publication failure")
				}
				return nil
			}
			result := ApplyPlan(t.Context(), plan, false, false)
			afterPublicationMutation = previous
			if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
				t.Fatalf("index %d result = %#v", failureIndex, result)
			}
			parentAfter, parentError := os.ReadFile(parentPath)
			checkpointAfter, checkpointError := os.ReadFile(checkpointPath)
			if parentError != nil || checkpointError != nil || !bytes.Equal(parentAfter, parentBefore) || !bytes.Equal(checkpointAfter, checkpointBefore) {
				t.Fatalf("index %d preimages changed", failureIndex)
			}
			for _, absent := range []string{"do-work/queue/REQ-101-parent.md", "do-work/queue/REQ-901-repair-repository-gate.md", "do-work/.req-reservations/REQ-901"} {
				if _, statError := os.Lstat(filepath.Join(root, filepath.FromSlash(absent))); !os.IsNotExist(statError) {
					t.Fatalf("index %d left %s: %v", failureIndex, absent, statError)
				}
			}
		})
	}
}

func TestDeferGateRollsBackUntrackedCreateAndFoldTopologies(t *testing.T) {
	t.Run("dirty parent clean checkpoint", func(t *testing.T) {
		for failureIndex := 0; failureIndex < 4; failureIndex++ {
			root := newDeferGateRepository(t, false)
			runGitFixture(t, root, "add", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "clean checkpoint claim", "--", "do-work/CHECKPOINT.md")
			assertDeferGateRollback(t, root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil), failureIndex,
				[]string{"do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md"},
				[]string{"do-work/queue/REQ-101-parent.md", "do-work/queue/REQ-901-repair-repository-gate.md", "do-work/.req-reservations/REQ-901"})
		}
	})
	t.Run("tracked clean create", func(t *testing.T) {
		for failureIndex := 0; failureIndex < 4; failureIndex++ {
			root := newDeferGateRepository(t, false)
			runGitFixture(t, root, "add", "do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "tracked clean lifecycle inputs")
			assertDeferGateRollback(t, root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil), failureIndex,
				[]string{"do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md"},
				[]string{"do-work/queue/REQ-101-parent.md", "do-work/queue/REQ-901-repair-repository-gate.md", "do-work/.req-reservations/REQ-901"})
		}
	})
	t.Run("untracked create", func(t *testing.T) {
		for failureIndex := 0; failureIndex < 4; failureIndex++ {
			root := newDeferGateRepository(t, false)
			runGitFixture(t, root, "rm", "--cached", "do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "untracked lifecycle inputs")
			assertDeferGateRollback(t, root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil), failureIndex,
				[]string{"do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md"},
				[]string{"do-work/queue/REQ-101-parent.md", "do-work/queue/REQ-901-repair-repository-gate.md", "do-work/.req-reservations/REQ-901"})
		}
	})
	t.Run("untracked repair fold", func(t *testing.T) {
		for failureIndex := 0; failureIndex < 3; failureIndex++ {
			root := newDeferGateRepository(t, true)
			first := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
			if result := ApplyPlan(t.Context(), BuildDeferGatePlan(root, first), false, false); result.Outcome != resultmodel.OutcomeSuccess {
				t.Fatalf("seed deferral = %#v", result)
			}
			claimDeferParent(t, root, "REQ-102", "do-work/working/REQ-102-second.md", "Second parent")
			repairPath := "do-work/queue/REQ-901-repair-repository-gate.md"
			fold := deferGateManifest(root, "REQ-102", "REQ-901", "do-work/working/REQ-102-second.md", &PayloadFile{SourcePath: repairPath})
			assertDeferGateRollback(t, root, fold, failureIndex,
				[]string{"do-work/working/REQ-102-second.md", "do-work/CHECKPOINT.md", repairPath},
				[]string{"do-work/queue/REQ-102-second.md"})
		}
	})
	t.Run("tracked clean fold", func(t *testing.T) {
		for failureIndex := 0; failureIndex < 3; failureIndex++ {
			root := newDeferGateRepository(t, true)
			first := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
			if result := ApplyPlan(t.Context(), BuildDeferGatePlan(root, first), false, false); result.Outcome != resultmodel.OutcomeSuccess {
				t.Fatalf("seed deferral = %#v", result)
			}
			runGitFixture(t, root, "add", ".")
			runGitFixture(t, root, "commit", "-qm", "tracked repair deferral")
			claimDeferParent(t, root, "REQ-102", "do-work/working/REQ-102-second.md", "Second parent")
			runGitFixture(t, root, "add", "do-work/working/REQ-102-second.md", "do-work/CHECKPOINT.md")
			runGitFixture(t, root, "commit", "-qm", "tracked clean fold inputs")
			repairPath := "do-work/queue/REQ-901-repair-repository-gate.md"
			fold := deferGateManifest(root, "REQ-102", "REQ-901", "do-work/working/REQ-102-second.md", &PayloadFile{SourcePath: repairPath})
			assertDeferGateRollback(t, root, fold, failureIndex,
				[]string{"do-work/working/REQ-102-second.md", "do-work/CHECKPOINT.md", repairPath},
				[]string{"do-work/queue/REQ-102-second.md"})
		}
	})
}

func assertDeferGateRollback(t *testing.T, root string, manifest Manifest, failureIndex int, preimagePaths, absentPaths []string) {
	t.Helper()
	preimages := map[string][]byte{}
	for _, path := range preimagePaths {
		contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		preimages[path] = contents
	}
	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal != nil {
		t.Fatalf("plan refusal = %#v", plan.Refusal)
	}
	previous := afterPublicationMutation
	afterPublicationMutation = func(index int, _ PlannedMutation) error {
		if index == failureIndex {
			return errors.New("injected publication failure")
		}
		return nil
	}
	t.Cleanup(func() { afterPublicationMutation = previous })
	result := ApplyPlan(t.Context(), plan, false, false)
	afterPublicationMutation = previous
	if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("index %d result = %#v", failureIndex, result)
	}
	for path, before := range preimages {
		after, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || !bytes.Equal(after, before) {
			t.Fatalf("index %d preimage %s changed: %v", failureIndex, path, err)
		}
	}
	for _, path := range absentPaths {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Fatalf("index %d left %s: %v", failureIndex, path, err)
		}
	}
}

func TestDeferGateRefusesUnsafeStaleCollidingAndStagedInputs(t *testing.T) {
	for name, mutate := range map[string]func(*testing.T, string, *Manifest){
		"unsafe repair path": func(_ *testing.T, _ string, manifest *Manifest) { manifest.DeferGate.RepairPath = "../outside.md" },
		"zero exit":          func(_ *testing.T, _ string, manifest *Manifest) { manifest.DeferGate.GateExitStatus = 0 },
		"unpaired merge": func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.DeferGate.DeferredImplementationBase = "HEAD"
		},
		"collision": func(t *testing.T, root string, _ *Manifest) {
			writeFixture(t, root, "do-work/queue/REQ-901-repair-repository-gate.md", []byte("collision\n"), 0o644)
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := newDeferGateRepository(t, false)
			manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
			mutate(t, root, &manifest)
			if plan := BuildDeferGatePlan(root, manifest); plan.Refusal == nil {
				t.Fatalf("%s accepted", name)
			}
		})
	}

	t.Run("stale parent", func(t *testing.T) {
		root := newDeferGateRepository(t, false)
		writeFixture(t, root, "payload/expected-parent.md", []byte("stale\n"), 0o644)
		manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
		manifest.DeferGate.ExpectedParent.SourcePath = "payload/expected-parent.md"
		if plan := BuildDeferGatePlan(root, manifest); plan.Refusal == nil || plan.Refusal.Code != "DEFER-GATE-PARENT-STALE" {
			t.Fatalf("stale plan = %#v", plan.Refusal)
		}
	})

	t.Run("staged parent", func(t *testing.T) {
		root := newDeferGateRepository(t, false)
		manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
		plan := BuildDeferGatePlan(root, manifest)
		runGitFixture(t, root, "add", "do-work/working/REQ-101-parent.md")
		result := ApplyPlan(t.Context(), plan, false, false)
		if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) == 0 || !strings.Contains(result.Findings[0].Evidence[0], "must not be staged") {
			t.Fatalf("staged result = %#v", result)
		}
	})
}

func TestDeferGatePersistsCanonicalPairedMergeEvidence(t *testing.T) {
	root := newDeferGateRepository(t, false)
	manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
	manifest.DeferGate.DeferredImplementationBase = "HEAD~1"
	manifest.DeferGate.DeferredImplementationMerge = "HEAD"
	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal != nil || plan.GateDeferral == nil || len(plan.GateDeferral.DeferredImplementationBase) != 40 || len(plan.GateDeferral.DeferredImplementationMerge) != 40 {
		t.Fatalf("paired merge plan = %#v, refusal=%#v", plan.GateDeferral, plan.Refusal)
	}
	result := ApplyPlan(t.Context(), plan, false, false)
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("paired merge apply = %#v", result)
	}
	parentBytes, err := os.ReadFile(filepath.Join(root, "do-work/queue/REQ-101-parent.md"))
	if err != nil {
		t.Fatal(err)
	}
	document, _ := requestmodel.ParseDocument(parentBytes)
	record := document.TypedRecord()
	if record.DeferredImplementationBaseValue != plan.GateDeferral.DeferredImplementationBase || record.DeferredImplementationMergeValue != plan.GateDeferral.DeferredImplementationMerge {
		t.Fatalf("persisted merge evidence = %#v", record)
	}
}

func TestDeferGateRejectsDeferredMergeFromSideBranch(t *testing.T) {
	root := newDeferGateRepository(t, false)
	runGitFixture(t, root, "add", "do-work/working/REQ-101-parent.md", "do-work/CHECKPOINT.md")
	runGitFixture(t, root, "commit", "-qm", "claimed lifecycle state")
	mainBranch := runGitOutputFixture(t, root, "branch", "--show-current")
	base := runGitOutputFixture(t, root, "rev-parse", "HEAD")
	runGitFixture(t, root, "checkout", "-qb", "side-implementation")
	writeFixture(t, root, "side-only.txt", []byte("side implementation\n"), 0o644)
	runGitFixture(t, root, "add", "side-only.txt")
	runGitFixture(t, root, "commit", "-qm", "side implementation")
	merge := runGitOutputFixture(t, root, "rev-parse", "HEAD")
	runGitFixture(t, root, "checkout", "-q", mainBranch)

	manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
	manifest.DeferGate.DeferredImplementationBase = base
	manifest.DeferGate.DeferredImplementationMerge = merge
	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal == nil || plan.Refusal.Code != "DEFER-GATE-MERGE-EVIDENCE-INVALID" || !strings.Contains(plan.Refusal.Reason, "current HEAD") {
		t.Fatalf("side-branch merge plan = %#v", plan.Refusal)
	}
}

func TestDeferGateUsesExactFingerprintAndCheckpointClaimIdentities(t *testing.T) {
	t.Run("checkpoint prefixes survive", func(t *testing.T) {
		root := newDeferGateRepository(t, false)
		checkpointPath := filepath.Join(root, "do-work/CHECKPOINT.md")
		checkpoint, err := os.ReadFile(checkpointPath)
		if err != nil {
			t.Fatal(err)
		}
		checkpoint = append(checkpoint,
			[]byte("- REQ-1010: Prefix id — claimed now — writer: host:/repo\n  prefix id detail\n- REQ-101: Prefix writer — claimed now — writer: host:/repo-shadow\n  prefix writer detail\n")...)
		if err := os.WriteFile(checkpointPath, checkpoint, 0o644); err != nil {
			t.Fatal(err)
		}
		plan := BuildDeferGatePlan(root, deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil))
		if plan.Refusal != nil {
			t.Fatalf("plan refusal = %#v", plan.Refusal)
		}
		if result := ApplyPlan(t.Context(), plan, false, false); result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("apply = %#v", result)
		}
		after, err := os.ReadFile(checkpointPath)
		if err != nil || !bytes.Contains(after, []byte("REQ-1010")) || !bytes.Contains(after, []byte("host:/repo-shadow")) || bytes.Contains(after, []byte("parent detail")) {
			t.Fatalf("exact claim removal failed: %v\n%s", err, after)
		}
	})

	t.Run("fingerprint prefix is not a fold", func(t *testing.T) {
		root := newDeferGateRepository(t, false)
		manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
		candidate := *manifest.DeferGate
		candidate.RepairID = "REQ-902"
		candidate.RepairPath = "do-work/queue/REQ-902-repair-repository-gate.md"
		candidate.ReservationPath = "do-work/.req-reservations/REQ-902"
		candidate.DiagnosticFingerprint = manifest.DeferGate.DiagnosticFingerprint + "-extra"
		repairBytes, err := authoredRepairBytes(&candidate, "UR-095", "", "")
		if err != nil {
			t.Fatal(err)
		}
		writeFixture(t, root, candidate.RepairPath, repairBytes, 0o644)
		plan := BuildDeferGatePlan(root, manifest)
		if plan.Refusal != nil || plan.GateDeferral == nil || plan.GateDeferral.RepairOutcome != "created" {
			t.Fatalf("prefix fingerprint matched: plan=%#v refusal=%#v", plan.GateDeferral, plan.Refusal)
		}
	})
}

func TestDeferGateRefusesAmbiguousMatchingRepairs(t *testing.T) {
	root := newDeferGateRepository(t, false)
	manifest := deferGateManifest(root, "REQ-101", "REQ-901", "do-work/working/REQ-101-parent.md", nil)
	for _, repairID := range []string{"REQ-901", "REQ-902"} {
		candidate := *manifest.DeferGate
		candidate.RepairID = repairID
		repairBytes, err := authoredRepairBytes(&candidate, "UR-095", "", "")
		if err != nil {
			t.Fatal(err)
		}
		writeFixture(t, root, "do-work/queue/"+repairID+"-repair-repository-gate.md", repairBytes, 0o644)
	}
	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal == nil || plan.Refusal.Code != "DEFER-GATE-REPAIR-AMBIGUOUS" || len(plan.Refusal.Paths) != 2 {
		t.Fatalf("ambiguous plan = %#v", plan.Refusal)
	}
}

func TestDeferGateRefusesClaimedMatchingRepairInsteadOfCreatingDuplicate(t *testing.T) {
	root := newDeferGateRepository(t, false)
	manifest := deferGateManifest(root, "REQ-101", "REQ-902", "do-work/working/REQ-101-parent.md", nil)
	candidate := *manifest.DeferGate
	candidate.RepairID = "REQ-901"
	repairBytes, err := authoredRepairBytes(&candidate, "UR-095", "", "")
	if err != nil {
		t.Fatal(err)
	}
	repairBytes = bytes.Replace(repairBytes, []byte("status: pending"), []byte("status: claimed\nclaimed_at: 2026-09-02T01:00:00Z"), 1)
	writeFixture(t, root, "do-work/working/REQ-901-repair-repository-gate.md", repairBytes, 0o644)

	plan := BuildDeferGatePlan(root, manifest)
	if plan.Refusal == nil || plan.Refusal.Code != "DEFER-GATE-REPAIR-MISMATCH" {
		t.Fatalf("claimed matching repair plan = %#v", plan.Refusal)
	}
}

func newDeferGateRepository(t *testing.T, secondParent bool) string {
	t.Helper()
	root := t.TempDir()
	runGitFixture(t, root, "init", "-q")
	runGitFixture(t, root, "config", "user.name", "Test")
	runGitFixture(t, root, "config", "user.email", "test@example.com")
	writeFixture(t, root, "do-work/working/REQ-101-parent.md", pendingParentBytes("REQ-101", "Parent"), 0o644)
	if secondParent {
		writeFixture(t, root, "do-work/working/REQ-102-second.md", pendingParentBytes("REQ-102", "Second parent"), 0o644)
	}
	writeFixture(t, root, "do-work/CHECKPOINT.md", []byte("# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-999: Foreign — claimed earlier — writer: other:/repo\n  foreign detail\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "baseline")
	writeFixture(t, root, "gate-merge-evidence.txt", []byte("implementation merge\n"), 0o644)
	runGitFixture(t, root, "add", "gate-merge-evidence.txt")
	runGitFixture(t, root, "commit", "-qm", "implementation merge evidence")
	claimDeferParent(t, root, "REQ-101", "do-work/working/REQ-101-parent.md", "Parent")
	return root
}

func pendingParentBytes(identifier, title string) []byte {
	return []byte("---\nid: " + identifier + "\ntitle: '" + title + "'\nstatus: pending\nuser_request: UR-095\ndepends_on: []\n---\n\n# " + title + "\n")
}

func claimDeferParent(t *testing.T, root, identifier, relativePath, title string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	contents = bytes.Replace(contents, []byte("status: pending"), []byte("status: claimed\nclaimed_at: 2026-09-02T01:00:00Z"), 1)
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		t.Fatal(err)
	}
	checkpointPath := filepath.Join(root, "do-work/CHECKPOINT.md")
	checkpoint, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatal(err)
	}
	entry := []byte("- " + identifier + ": " + title + " — claimed 2026-09-02T01:00:00Z — writer: host:/repo\n  parent detail\n")
	checkpoint = append(checkpoint, entry...)
	if err := os.WriteFile(checkpointPath, checkpoint, 0o644); err != nil {
		t.Fatal(err)
	}
}

func deferGateManifest(root, parentID, repairID, parentPath string, expectedRepair *PayloadFile) Manifest {
	_ = root
	return Manifest{Operation: OperationDeferGate, DeferGate: &DeferGateManifest{
		ParentID: parentID, ParentPath: parentPath, ExpectedParent: PayloadFile{SourcePath: parentPath}, ExpectedStatus: "claimed",
		CheckpointPath: "do-work/CHECKPOINT.md", ExpectedCheckpoint: PayloadFile{SourcePath: "do-work/CHECKPOINT.md"}, WriterLabel: "host:/repo",
		GateCommand: []string{"bash", "_dev/tests/maintainer-verify.sh"}, GateExitStatus: 17,
		DiagnosticFingerprint: "sha256:gate-red", DiagnosticEvidence: []string{"go test failed"}, SweepKey: "repository-gate-sha256-gate-red",
		RepairID: repairID, RepairPath: "do-work/queue/" + repairID + "-repair-repository-gate.md", RepairTitle: "Repair repository gate failure",
		RepairCreatedAt: "2026-09-02T01:02:03Z", ReservationPath: "do-work/.req-reservations/" + repairID, ExpectedRepair: expectedRepair,
	}}
}

func runGitOutputFixture(t *testing.T, root string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return strings.TrimSpace(string(output))
}
