# Re-Review: REQ-515 (per-REQ recovery findings never stop the loop)

**Approve with follow-ups** — every finding sent back for fixing (F1, F2, F3, M1, M3, M5) is genuinely closed and independently red-green verified. One new Important finding limits how long the loop keeps running after a set-aside, but it is a fail-safe stop, it pre-dates the remediation, and it is strictly better than the behaviour this REQ replaced.

Route B | cumulative range `18666d7..fe2de1e` (original `28c1460`, remediation `fe2de1e`)
Independent re-reviewer; wrote none of this code. Every claim below was reproduced, not read.

## Verification method

Work was done on a copy of the Go module at `/tmp/.../scratchpad/cli`; the real working tree was never modified (`git status --porcelain` is empty, `gofmt -l` clean on both changed packages, `git diff --check 18666d7..fe2de1e` clean).

- `go test -count=1 ./...` across the whole `do-work-cli` module — all packages pass.
- `go vet ./internal/finalization ./internal/lifecycleadvance` — clean.
- `bash _dev/tests/contracts/recovery-set-aside.sh` → exit 0; `shellcheck` on it → clean. `bash _dev/tests/contract-regressions.sh` → exit 0.
- Three fix-removal experiments (below), each confirming the matching test turns red and passes again on restore.
- Four constructed probes on the copy for paths no test covers.

## Per-finding verdicts

### F1 — a set-aside REQ was released back into the queue — **CLOSED**

The fix is in `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go:61` (collect the ids) and `:84-88` (skip claim recovery). It sits **before** the `authorized` test, so it covers every authority mode, not only `--assume-sole-authority`.

Fix-removal proof: with the skip block deleted and `setAsideRequestIDs` reduced to a no-op,

```
--- FAIL: TestRecoverPreservesTheClaimOfARequestFinalizationSetAside
    the set-aside REQ lost its claim: stat .../do-work/working/REQ-730.md: no such file or directory
```

which is the exact symptom the first review reproduced. Restored → pass.

Path coverage, each constructed and run on the copy:

| Path | Result |
|---|---|
| `recover --assume-sole-authority` | REQ-730 keeps its claim, `do-work/queue/REQ-730*` empty, unrelated REQ-731 still recovered |
| `recover --take-over REQ-730` | claim preserved, decision `finalization set aside; claim preserved`, outcome `success`, `FINALIZATION-SET-ASIDE` finding names REQ-730 (this is D-14; correct, see M9) |
| plain `recover` (observe) | claim preserved; **no** `RECOVERY-TAKEOVER-AVAILABLE` finding for the set-aside REQ (see M8) |
| second `recover` in the same run | claim still preserved and still not queued — but the run now stops, see F5 |
| REQ set aside by discovery rather than the journal path | not constructible: discovery refusals return from `finalization_commands.go:68-74` before the folding loop, so no `FINALIZATION-SET-ASIDE` record exists and `handleRecover` returns before claim recovery |

### F2 — shared-cause refusals were being set aside per REQ — **CLOSED**

`finalization_apply.go:127` now routes `commitSafety` refusals through the new `sharedStateRefusal` (`:645-656`), which clears `AffectedIDs`, so `requestScopedRefusal` returns false and the run stops.

Fix-removal proof: with `:127` reverted to `finalizationFailure`,

```
--- FAIL: TestRecoverFinalizationStopsOnSharedDirtInsteadOfSettingOneRequestAside
    shared dirt no REQ owns must stop the run: outcome "success",
    both REQ-720 and REQ-721 carrying [FINALIZATION-DIRTY-INDEX FINALIZATION-SET-ASIDE]
```

The whole repository's one dirty index was being set aside once per REQ — exactly what F2 described. Restored → pass.

**Consumer sweep on the removed `affected_ids` (this is the part most likely to break something, and it did not):**

- Go: full module test suite green. The only code that reads `AffectedIDs` off finalization-shaped findings is `internal/hookcommands/session_start.go:90` (`firstFinalizationTailRequestID`), and it is fed by `doctor.FinalizationTailFindings`, a different producer that this diff never touches.
- The `FinalizationResult` record still carries `RequestID` — only the `CommandFinding` lost the id — so every selector, exit summary, and JSON consumer that identifies a record by `request_id` is unaffected.
- Prose/shell: no action file, hook, or test script parses `affected_ids` for these codes. `skills/do-work/actions/forensics.md:54` maps `affected_ids` generically and renders an empty list fine. `skills/do-work/next-steps.md:32` names `FINALIZATION-AMBIGUOUS-SHARED-STATE` as a whole-run stop resolved by `do-work run-with-recovery` — still true, and truer than before.
- Exit code unchanged: `result_model.go:515` maps `OutcomeFindings` and `OutcomeRefused` to the same status.
- One observable shape change, harmless: for the standalone `recover-finalization` verb, the REQ-514 normalizer at `result_model.go:783-796` used to rewrite a dirty-index refusal to outcome `findings` (because its `next_argv` named the command that just ran). It now stays `refused`, with `next_argv` = `uncommitted-inventory`. Same exit code, and every action reads only "typed `success`". Recorded as N3.

### F3 — the never-widen lock-in never reached the new code — **CLOSED**

Two new tests reach `consumeRecoveryRecord` and drive `requestScopedRefusal` to false on both live branches, and both were proven red by removing the thing they protect:

- `TestRecoverFinalizationStopsOnSharedDirtInsteadOfSettingOneRequestAside` — the unowned branch (proof above).
- `TestRecoverFinalizationStopsWhenRollbackLeavesResidue` — the residue branch. With `|| result.Rollback.Status == resultmodel.RollbackIncomplete` deleted from `finalization_commands.go:151`, it fails with outcome `success` and both REQs set aside despite `rollback: incomplete`. Restored → pass.

The original `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest` still only exercises the discovery path, but it is no longer the only never-widen pin, so F3's actual complaint is answered.

### M1 — membership instead of exclusivity — **CLOSED**

`namesRequestID` became `namesOnlyRequestID` (`finalization_commands.go:170-172`) with `len(affectedIDs) == 1 && affectedIDs[0] == requestID`. A two-id finding can no longer become one REQ's private exclusion.

### M3 — `[title]` was not a field of `FinalizationResult` — **CLOSED**

The exit-summary template at `skills/do-work/actions/work-reference.md:678-684` now renders `REQ-NNN (set aside: <reason codes, comma-separated>)` with no title, so "render each line from the record's own fields" is true as written.

### M5 — the contract's missing-file branch was guarded on the wrong variable — **CLOSED**

`_dev/tests/contracts/recovery-set-aside.sh:28-41`: the missing-file branch moved inside `require_action_phrase` and is guarded on `[ ! -f "$core_actions/$action_file" ]`. A missing action file now prints "is missing", not a misleading "must state" plus grep stderr.

## New findings

**Important:**

- **F5 — the set-aside survives only until the shared checkpoint moves; the next queue boundary then stops the whole run.** `recover` runs at every queue boundary (`actions/work.md` Step 1). Constructed on the copy: set REQ-730 aside, then make an ordinary later commit to `do-work/CHECKPOINT.md` — the write every subsequent claim and finalize performs — and re-run `recover --assume-sole-authority`. Result: outcome `refused`, record `REQ-730 codes=[FINALIZATION-LIFECYCLE-CONFLICT FINALIZATION-ROLLBACK-INCOMPLETE]`, whole run parked. Mechanism: at `PhasePrepared` the journal refuses with `FINALIZATION-LIFECYCLE-CONFLICT` before mutating anything (`finalization_apply.go:54`), but the deferred rollback at `:32-43` still runs `convergeImages(postimages → preimages)` on `do-work/CHECKPOINT.md`, which now matches neither image, so the status becomes `RollbackIncomplete` and `requestScopedRefusal` refuses to set the record aside. There is no residue to protect against — nothing was applied in that pass. So `run-with-recovery.md:33` ("every other REQ still runs") and `work.md:127` ("the run continues with what remains") hold for the boundary where the set-aside happened and stop holding at the next one. This is not a remediation regression: the `RollbackIncomplete` stop shipped in the first pass, and the remediation only added the test that pins it. It is still better than the behaviour being replaced, which parked the queue at the first refusal. — `impact-user-visible` → report only

**Minor:**

- **M7 — "the whole-run stop is dirt no REQ owns" is now the documented rule, and one live stop contradicts it.** `work-reference.md:1081` states "its finding names no REQ" and `work.md:127` glosses the stop as "dirt no REQ owns". The `RollbackIncomplete` stop names exactly one REQ and still stops the run, and `_dev/tests/contracts/recovery-set-aside.sh:66` pins the phrase `its finding names no REQ`, so the contract now protects a statement that is true of one of the two live stop conditions. — `impact-rule-change` → report only
- **M8 — plain `recover` no longer offers takeover for a set-aside REQ, and the prose still promises it does.** Verified: with REQ-730 set aside and REQ-731 clean, plain `recover` emits `RECOVERY-TAKEOVER-AVAILABLE` for REQ-731 only. `work-reference.md:362` still opens with "plain recovery preserves claims and returns exact takeover argv"; the set-aside exception the remediation added attaches only to the explicit-authority half of that sentence. The behaviour is right (it is D-14); the sentence is half-updated. — `impact-negligible` → report only
- **M9 — D-14 has no test.** `recover --take-over <set-aside REQ>` preserving the claim is behaviour a future change can silently reverse: `recovery_set_aside_test.go` only drives `--assume-sole-authority`. I verified `--take-over` by hand and it is correct, including that the `FINALIZATION-SET-ASIDE` finding still names the REQ so the user sees why nothing was handed over. One extra case in the existing test would pin it. — `impact-negligible` → report only

**Nit:**

- **N2 — `sharedStateRefusal` also swallows two causes the prose does not describe.** `commitSafety` returns `FINALIZATION-INVENTORY-FAILED` and an `FINALIZATION-AMBIGUOUS-SHARED-STATE` variant for a protected path that is *inside* the journal's allowlist (`finalization_apply.go:487-489`), while `work-reference.md:1081` describes the latter as paths "outside the recovery group". Both are genuinely repository-wide (an inventory read failure, and the global quarantine at `corehelpers/inventory.go:122-131` that reclassifies every `A` row to `X` when any secret path is dirty), so stopping the run is the right call — the sentence just does not cover them. — `impact-negligible` → report only
- **N3 — standalone `recover-finalization` on shared dirt reports `refused` where it used to report `findings`.** Same exit code, no consumer reads the difference. Recorded so the shape change is on the record. — `impact-negligible` → report only

## Findings carried forward unchanged

F4 (the exclusion never reaches the selector — `advance` and `next` read no finalization records), M2 (the resolving verb under `run-with-recovery` is the verb that just ran), M4 (seven exact-phrase contract pins), M6 (P-A-U boxes — now checked `[x]`) were accepted as report-only by the first review and were not in the remediation's scope. Nothing in the merged diff changes them, except that F4's practical exposure shrank: with F1 fixed, a set-aside REQ stays in `do-work/working/` under both `run` and `run-with-recovery`, so queue selection cannot see it in either mode.

## Never-widen constraint

I could not construct a refusal that used to stop the run and is now set aside in a way that lets shared or foreign bytes through. The full classification of refusals reachable from `advanceJournal`:

- Unowned, stop the run: `FINALIZATION-DIRTY-INDEX`, `FINALIZATION-AMBIGUOUS-SHARED-STATE`, `FINALIZATION-INVENTORY-FAILED` (all via `sharedStateRefusal`), plus every discovery-level refusal, which returns before the folding loop.
- Owned, stop the run anyway: any record whose rollback is incomplete, and any command-level `failure`.
- Owned, set aside: lifecycle conflict/plan/apply/recovery, journal write, commit and verify refusals. These are the per-REQ exclusions the REQ asks for.

The remediation moved refusals from "set aside" to "stop" only. It moved none the other way. Downstream is still fail-closed: `gittransaction` refuses every commit on a non-empty index.

## Restatement sweep

Run on what the remediation redefines: the `SetAsideReasonCode` constant became exported, and `commitSafety` refusals changed finding shape.

- `setAsideReasonCode` (unexported) has zero remaining references anywhere in the repo; the contract's extraction at `_dev/tests/contracts/recovery-set-aside.sh:12-13` reads the exported form and passes.
- Every sentence the REQ added or edited was re-read against the merged behaviour: `work.md:127`, `run-with-recovery.md:33` and its two new checklist/rationalization rows, `work-reference.md:362`, `:678-687`, `:1081`, `commit.md:51`, `docs/work-guide.md:130`. All are accurate for the entry points they name, with the three exceptions recorded above (F5's duration claim, M7's "names no REQ", M8's takeover-argv half-sentence). The specific sentences the first review found false — "the REQ was not claimed and nothing was written to it", and the `--assume-sole-authority` "resets every working claim" statements in Crash Recovery and `work-guide.md` — are all corrected and now match the code.

## Requirements Checklist

- [x] Step 1 in `work.md` and Step 0.1 in `run-with-recovery.md` iterate recovery records; a refused record excludes that REQ and the loop continues — delivered at the boundary where recovery ran, in both entry points, verified by test. F5 limits how many boundaries later it still holds.
- [x] The composed exit summary lists set-aside REQs with reason codes and resolving verbs — delivered, three contract predicates.
- [x] A finding with no owning REQ still stops the run and names a resolving verb per REQ-514 — now delivered through the journal path, not only the untouched discovery path. The verb is `uncommitted-inventory`, a genuinely different command from the one that refused.
- [x] Per-record contract predicates plus a CLI behavior test for a mixed result — delivered.
- [x] Serial and fan-out modes behave the same — delivered; the mode split happens after selection.
- [x] Never widen what recovery accepts — held; see the classification above.
- [x] Floor agent can follow the loop from the output plus prose — held.
- [x] Behavior test or contract predicate, never a sentence pin alone — held; four Go behavior tests now, three of them red-proven here.

## Acceptance Testing

**Result: Partial**

- Whole `do-work-cli` module: `go test -count=1 ./...` green. `go vet` clean, `gofmt -l` empty.
- `_dev/tests/contracts/recovery-set-aside.sh` exit 0, ShellCheck clean; `_dev/tests/contract-regressions.sh` exit 0.
- Three fix-removal experiments: all three new tests go red for the right reason and green on restore. No red-green claim in the `## Remediation` section was taken on trust; all were re-run.
- Five constructed probes on `handleRecover`: sole-authority, take-over, plain, second invocation, and post-checkpoint-mutation.
- Partial rather than Pass because of F5: the set-aside holds for the boundary that produced it and for a repeat `recover` against unchanged state, but the run stops at the next boundary once `do-work/CHECKPOINT.md` has moved.
- Not tested: a live end-to-end `do-work run` against a real queue in this checkout.

## Suggested Additional Testing

1. A Go test that sets a REQ aside, mutates and commits `do-work/CHECKPOINT.md`, and asserts the next `recover` still returns `success` — F5's regression pin, and the shape the fix would take (do not score a rollback as incomplete when the journal is at `PhasePrepared` and applied nothing).
2. Extend `TestRecoverPreservesTheClaimOfARequestFinalizationSetAside` with a `--take-over` case, pinning D-14.
3. A fixture where a finalization refusal names two REQs, pinning `namesOnlyRequestID`'s exclusivity rather than leaving it structural-only.
4. Manual: run `do-work run` against a queue with one stuck finalization tail and three pending REQs, and count how many get built before the run stops.

## Scores (on the record — not the headline)

**Overall: 79%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 85% | All eight acceptance criteria delivered; the "loop continues" promise is durable for one boundary, not the whole run (F5) |
| Code Quality | 90% | Small helpers with why-comments; ownership stated as a condition, not a code list; the `PhasePrepared` rollback that manufactures a false residue is the one real blemish |
| Test Adequacy | 85% | Three genuine red-green tests verified by fix removal; D-14 and the multi-boundary path untested |
| Scope | 95% | Four files outside the first write set, each escalated and recorded (D-10, D-13); all four are restatements of what this REQ changed |
| Risk | Low | Every remaining gap fails closed — the run stops rather than dispatching onto a broken journal |
| Acceptance | Partial | Named path works in both entry points; the run stops one boundary later than the prose implies |

Formula: average of 85/90/85/95 = 88.75%; Acceptance = Partial applies a 10-point penalty → 78.75%, rounded to 79%. No cap applies (Risk is not Critical, Acceptance is not Fail).

## Follow-ups created

None written by this reviewer — the orchestrator owns writes under `do-work/`. All six new findings are noncritical and report only.

## Self-Validation

Three things I re-checked before writing this up.

1. **That the F1 fix is not just passing its own test.** I removed it and watched the test fail with the first review's exact error string, then constructed three more invocation paths by hand. All four preserve the claim.
2. **That removing `affected_ids` did not break a consumer I had not thought of.** I grepped every `.go`, `.md`, and `.sh` in the repo for `AffectedIDs`/`affected_ids`, traced the one finalization-shaped consumer (`session_start.go`) to a different producer, and confirmed the exit-code mapping is identical for `findings` and `refused`. The full module suite is green.
3. **That F5 is real and not an artifact of my fixture.** The first observation came from a second `recover` where the checkpoint had moved because the *first* recover released an unrelated claim. I reproduced it a second time with the checkpoint moved by an ordinary commit instead, with no claim recovery involved, and got the same `LIFECYCLE-CONFLICT` + `ROLLBACK-INCOMPLETE` stop.

---

**Verdict: Pass**

Acceptance reasoning: the remediation had six jobs and did all six. F1 is fixed at the only place it could be fixed, and the fix sits ahead of the authority check so it covers `--assume-sole-authority`, `--take-over`, plain observe mode, and a repeat invocation — I removed the fix and watched the first review's exact failure come back, then restored it and constructed each path by hand. F2 is fixed at the producer rather than with a code list, which is what makes the "a finding that owns no REQ stops the run" branch reachable from the journal path for the first time; removing affected ids from those two findings breaks no Go, shell, or prose consumer, and the exit code is unchanged. F3 now has two red-proven tests on the two live false branches instead of one test that never reached the code. M1, M3 and M5 are small and correct. The three prose sentences that asserted the opposite of the shipped behaviour are corrected, and the sweep that found them was repeated across every edited sentence.

What keeps this at Partial acceptance rather than a clean pass is F5: a set-aside REQ stops the whole run at the next queue boundary once `do-work/CHECKPOINT.md` moves, because a lifecycle conflict detected before anything is applied still drags a failed no-op rollback behind it and scores as residue. That undercuts "every other REQ still runs" beyond the first boundary. It is not a remediation regression — the stop shipped in the first pass and the remediation only pinned it with a test — and its direction is fail-safe: the run parks with typed evidence instead of dispatching a builder onto a journal that will refuse. Against the behaviour this REQ replaces, where the first refused record parked the entire queue, this is a real improvement and not a reason to send the work back a second time. It belongs in the queue as a follow-up, not in another remediation round.

*Reviewed by review-work action (independent re-reviewer)*
