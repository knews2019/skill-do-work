---
id: REQ-515
title: '[impact-rule-change] Per-REQ recovery findings never stop the loop'
status: claimed
priority: now
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-514]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-513, REQ-514, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/run-with-recovery.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/commit.md, _dev/tests/contract-regressions.sh, _dev/tests/contracts/recovery-set-aside.sh, skills/do-work/docs/work-guide.md, skills/do-work/tools/do-work-cli/internal/finalization/, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/]
claimed_at: 2026-09-04T18:15:54Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-04T18:17:56Z
  basis:
    - Route B
    - 5-file write set
    - 3 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
dispatch_at: 2026-09-04T18:23:51Z
builder_handback_at: 2026-09-04T19:09:41Z
integration_at: 2026-09-04T19:09:41Z
review_at: 2026-09-04T19:24:12Z
remediation_at: 2026-09-04T19:36:54Z
re_review_at: 2026-09-04T19:36:54Z
---

# Per-REQ recovery findings never stop the loop

## What

Run Step 1 recovery per REQ. Each refused finalization or claim-recovery record becomes an exclusion with its reason code in the selector output, and selection continues with what remains. The only global stop left is a finding that owns no REQ, which is what shared-target dirt looks like.

The fold-first scan found REQ-469 (Replace the unrelated canonical-gate hold with a blocked set-aside) and REQ-504 (Collapse Step 10 and Crash Recovery prose into recovery) as neighbors: REQ-469 sets aside a gate failure inside a build, REQ-504 shortens the recovery prose once commands own it. Neither changes recovery's stop-versus-continue behavior, so this is a new REQ.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both primes, REQ-514's archived set-aside contract, and the crew rules; located the whole-run gate in `handleRecoverFinalization` and its three prose restatements.
- [x] **[APPLY]:** Per-record folding in `finalization_commands.go` with two Go behavior tests; Step 1 and Step 0.1 prose rewritten per record; a new exit-summary section; a new owner contract with ten predicates. Two files outside the declared write set, both flagged before writing (D-03, seam).
- [x] **[UNIFY]:** `git diff --stat` reviewed across 8 files; gofmt clean, `go vet` clean, the do-work-cli module and the contract suite green; a restatement sweep for six contract phrases across every shipped Markdown file, not only the edited ones.

## Why

Both `run` and `rwr` put a global gate in front of the loop: Step 0.1 and Step 1 recovery are "recover everything or stop". REQ-456's stuck commit tail therefore parked 31 pending REQs. The maintainer's principle is that a failed REQ is set aside with a typed finding and the loop continues; only shared-target dirt may stop it.

## Context

`recover-finalization --discover` already returns an ordered `finalizations` list with one record per REQ. `actions/run-with-recovery.md` Step 0.1 says continue only when every record is terminal, and `actions/work.md` Step 1 has the same shape. The change is mostly on the action side, plus whatever the CLI needs to report per-record refusals as exclusions the selector understands.

## Detailed Requirements

- Step 1 in `actions/work.md` and Step 0.1 in `actions/run-with-recovery.md` iterate recovery records; a refused record excludes that REQ from this run's selection with its reason code and the loop continues.
- The composed exit summary lists set-aside REQs with their reason codes and resolving verbs.
- A finding with no owning REQ, such as dirt on a shared target that no REQ wrote, still stops the run, and it names a resolving verb per REQ-514.
- Contract predicates that pin "continue only if every record is clean" are replaced by predicates on the per-record wording, and the CLI carries a behavior test for a mixed result: one refused record, one clean record, selection proceeds.
- Serial and fan-out modes behave the same.

## Constraints

- Never widen what recovery accepts; this REQ changes what happens after a refusal, not whether it refuses.
- Keep the floor agent able to follow the loop with the command output plus the remaining prose.
- Coordinate wording with REQ-504 if both are in flight; the write sets overlap on `work.md` and `run-with-recovery.md`.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

Depends on REQ-514 for the set-aside finding shape. Related to REQ-469, REQ-472, and REQ-504.

## Builder Guidance

Certainty level: Firm on the behavior, latitude on how the exclusion is projected into the selector. Read `_dev/primes/prime-action-files.md` before touching an action file.

## Red-Green Proof

**RED prompt/case:** With REQ-456's journal at `prepared` and its checkpoint dirty, run `do-work run` on a queue with other claimable REQs.
**Why RED now:** Step 1 stops at the first refused finalization record and no other REQ is selected.
**GREEN when:** The same state reports REQ-456 as set aside with its reason code, selects the next claimable REQ, and the exit summary lists the set-aside with a resolving verb.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing action routing and status contracts.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A3 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Recovery record iteration spans two action files, the contract regression suite, and the CLI's finalization projection; the required behavior is firm, the exact projection sites need discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ⚠ `bash _dev/tests/maintainer-verify.sh` red BEFORE any change — one pre-existing failure, `_dev/tests/session-start-hook-behavior.sh took 44s; each test file must finish under 30s`. A wall-clock budget miss on a slow container, no assertion failed. Recorded in `do-work/working/baseline-failures.txt` so Step 6.5 separates it from new regressions; not attributable to this REQ and not deferred to a repair REQ.
**Dependencies:** ✓ Go 1.26.1 and ShellCheck 0.11.0 provisioned for this session (container shipped Go 1.24.7 / no ShellCheck)

*Checked by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modify) — fold recovery results one record at a time
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modify) — the mixed-result behavior test and the never-widen lock-in
- `skills/do-work/actions/work.md` (modify) — Step 1 reads recovery per REQ
- `skills/do-work/actions/run-with-recovery.md` (modify) — Step 0.1 the same, plus one rationalization row and one checklist item
- `skills/do-work/actions/work-reference.md` (modify) — the finalization paragraph under the Commit & Metadata-Commit Procedure, and a new Composed Exit Summary section
- `_dev/tests/contract-regressions.sh` (modify) — register the new owner contract
- `_dev/tests/contracts/recovery-set-aside.sh` (new) — the per-record predicates
- `skills/do-work/actions/commit.md` (modify) — **added to scope during the run**, see D-06

**Files I will NOT touch:** `internal/lifecycleadvance/` (its own recovery gate is a separate surface), `CHANGELOG.md` and `VERSION` (Step 9 finalization).

**Acceptance criteria (restated from REQ):**
- [x] A refused recovery record excludes only its own REQ; the loop continues with the rest
- [x] The composed exit summary lists set-aside REQs with reason codes and resolving verbs
- [x] A finding owning no REQ still stops the run and names a resolving verb
- [x] Whole-run contract predicates replaced by per-record predicates
- [x] The CLI carries a behavior test for a mixed result: one refused record, one clean, selection proceeds
- [x] Serial and fan-out modes behave the same
- [x] Nothing widens what recovery accepts

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/run-with-recovery.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/contracts/recovery-set-aside.sh` (new)

**What was done:** `handleRecoverFinalization` folds journal results one at a time. A refusal owned by exactly one REQ becomes that REQ's exclusion — `FINALIZATION-SET-ASIDE` appended to the record's `reason_codes` and `next_argv` cleared so recovery never names itself as the fix — and the loop continues. Anything else still stops the run: a finding naming no REQ or more than one, a command-level failure, or an incomplete rollback. Step 1 and Step 0.1 now read recovery per record, a new Composed Exit Summary section reports set-asides, and a new owner contract holds ten predicates, the strongest of which reads the `setAsideReasonCode` constant out of the Go source and requires all three action files to name that exact token, so a rename in the CLI reddens instead of leaving prose citing a code nothing emits.

## Decisions

**D-01 — DECIDE & STATE. A set-aside carries no `next_argv`; the resolving verb comes from the exit summary.** The only CLI verbs recovery could name are the ones that just ran, so naming either would be the self-remedy REQ-514 forbids. The command emits the shape REQ-514 defined and the exit summary names the verb by the fork that already exists in Stuck Runs Hand Off to Judgment: `do-work run-with-recovery` when this checkout is the only writer and releaser, `do-work cleanup` when the archive needs repair.

**D-02 — DECIDE & STATE. The REQ-scoped test is ownership plus no residue, not a path classifier.** A record is set aside only when every finding on it names that one REQ, the record exists, the outcome is not a command-level failure, and the rollback status is not incomplete. A fifth condition inspecting blocked paths for shared targets was considered and dropped: `sharedFinalizationPath` treats every `do-work/…` path as shared, including the REQ's own request file, so it cannot separate the two, and a hand-maintained path list would go stale. The residue check is the real safety half — an incomplete rollback means bytes the next claim would write through, and that still stops the run.

**D-03 — DECIDE & STATE, escalated before writing. The new predicates went into a new owner contract file outside the declared write set.** `contract-regressions.sh` is a 76-line runner under a hard 77-line ratchet, and every predicate in the suite lives in a single-subject owner contract under `_dev/tests/contracts/`. Adding ~50 lines inline would have breached the ceiling; raising the ceiling to hold predicates that do not belong there would have been worse. The new file was registered with one added line, taking the runner to exactly the ceiling without raising it. The builder flagged this rather than writing it silently; the orchestrator accepted it and extended the scope list and `write_set` above.

**D-04 — DECIDE & STATE. Recovery still returns typed `success` when it sets a REQ aside.** `OutcomeFindings` exits 1, and `lifecycleadvance/recovery_commands.go` stops `recover` on any non-success finalization result, so the loop would still have been parked and every floor agent reading "continue only on typed success" would have stopped too. Success is also the honest report: the command settled everything it could and typed the rest per REQ.

**D-05 — DECIDE & STATE. The new exit-summary section was appended as 9, not inserted.** Inserting earlier would have renumbered a section a concurrent builder was working beside. Ordering in that list carries no semantics.

**D-06 — DECIDE & STATE, orchestrator. The handed-back `commit.md` seam was applied, and `commit.md` added to the declared scope.** Line 51 restated the whole-run gate this REQ retires. After the change a set-aside record carries a reason code while the command returns success, so `do-work commit` would have stopped on exactly the state `do-work run` is being taught to walk past — the same command's output contract read two different ways in two actions, which is the `alternate-writer-contract-drift` family the prime names. Fail-closed is not a defence when the drift is introduced by this REQ, so this REQ closes it rather than leaving a follow-up.

## Discovered Tasks

- **impact-negligible, report only** — `recovery.FinalizationPassed` in `internal/lifecycleadvance/recovery_commands.go` is set to `true` even when a REQ was set aside. Defensible, since the finalization pass did complete and the per-record evidence says what was excluded, but a future REQ may want a distinct field so a consumer can see "passed with exclusions" without walking the records.
- **impact-rule-change, report only** — REQ-514's re-review left two open findings on this same surface. This REQ clears one of them for set-aside records only, by clearing their `next_argv`, and touches neither otherwise.

## Qualification

Passed — 8 files verified in the merge range `18666d7..28c1460`, 7 acceptance criteria traced, P-A-U confirmed. `qualify.sh` returned `OK: mechanical qualification passed`. Judgment checks: `finalization_commands.go` gained real per-record control flow with six named helpers, not a flag; the new contract file holds ten executable predicates rather than sentence pins; every requirement maps to a named file.

## Testing

**Tests run:** `go test ./internal/finalization/`, `bash _dev/tests/contracts/recovery-set-aside.sh`, and `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Finalization package green. Contract probes passed. Repository gate green on the rerun.

**Repository gate retry:** first run exited 1, rerun exited 0.

The first run failed on `_dev/tests/session-start-hook-behavior.sh took 30s; each test file must finish under 30s` — the budget is `< 30s`, so it missed by one second. The rerun measured the same file at 27s and the whole gate exited 0 with no failures.

**This is REQ-559's rule catching its first real case, in the same run that shipped it.** Under the previous rule the first non-zero exit was final: REQ-515 would have been deferred, a `repository_gate_repair` REQ minted and claimed, and its already-green no-op completion would have been the only thing letting REQ-515 resume — the REQ-548 sequence, reproduced exactly, over a one-second timing miss on a shared test file this REQ never touched. Instead the same argv was rerun once and the run continued. Cost: one gate run. Old cost: two gate runs, one claim, one checkpoint entry, one finalization, one archived REQ that changed nothing.

**Red-green validation:**
- `TestRecoverFinalizationSetsAsideRefusedRecordAndFinishesTheRest`: ✗ before implementation → ✓ after. The RED case from the REQ's Red-Green Proof — one refused record and one clean record, selection proceeds.
- `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest`: the never-widen lock-in. A refusal owning no REQ still stops the run.

**New tests added:**
- The two Go behavior tests above.
- `_dev/tests/contracts/recovery-set-aside.sh` — ten predicates. The strongest reads the `setAsideReasonCode` constant out of the Go source and requires all three action files to name that exact token, so a rename in the CLI reddens the suite instead of leaving prose citing a code nothing emits.

**Gate history across this run:** the same test file measured 44s, 39s, 32s, 30s and 27s on five separate runs with nothing touching it. It is a flake sitting near its budget on a slow container, which is why it produced a red gate twice and a green one three times.

## Review

**Verdict: Fail** — independent review, acceptance failed on a named requirement, overall capped at 50%. Full record: `do-work/runs/work-2026-09-04-182017/REQ-515-review.md`. Returned to Step 6 for remediation.

**F1 — impact-critical. A set-aside REQ is released back into the queue and re-selected.** Because the command now returns typed `success`, `recover --assume-sole-authority` goes on to claim recovery in `internal/lifecycleadvance/recovery_commands.go`. The reviewer reproduced it on this REQ's own `seedTwoPlannedFinalizations` fixture: after the set-aside the REQ's claim is recovered, it reappears in `do-work/queue/`, and `next` selects the very REQ the run just excluded — while its journal still sits at phase `prepared`. A builder then redoes the work and the finalize tail refuses with "an unfinished journal exists with a different manifest".

This is worse than the behaviour it replaced. The old code parked the queue; this re-dispatches work onto a REQ that cannot finalize. The prose this REQ added asserts the opposite in three places — `run-with-recovery.md` ("one REQ this run will not select"), `work.md` ("excludes that one REQ from this run's selection"), and `work-reference.md` ("the REQ was not claimed and nothing was written to it", when the claim reset is written and committed). Shipped prose describing shipped behaviour incorrectly.

**F2 — impact-rule-change. The shared-cause stop is effectively unreachable.** `requestScopedRefusal` returns false only for conditions `advanceJournal` never produces, plus an incomplete rollback. So `FINALIZATION-DIRTY-INDEX` and `FINALIZATION-AMBIGUOUS-SHARED-STATE` — the two codes that mean shared-target dirt — are stamped with the journal's REQ and become per-REQ exclusions. The maintainer's principle is that only shared-target dirt may stop the loop; setting shared-target dirt aside inverts it. Requirement 3 passes today only because discovery-level refusals return before the folding loop, which this diff did not touch.

**F3 — impact-rule-change. The never-widen lock-in does not exercise the code this REQ added.** `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest` refuses inside `discoverFinalizationJournals` and returns before `consumeRecoveryRecord` runs. It would pass unchanged if `requestScopedRefusal` returned true for everything. No test covers any false branch, including `RollbackIncomplete`, which is the only live stop.

**Recorded as report-only, not remediated:** F4 (the exclusion never reaches the selector — `advance` and `next` read no finalization records, so the exclusion lives in the recovery result and in prose; under plain `run` the two agree only because the REQ stays in `working/`), M2 (under `run-with-recovery` the prescribed resolving verb is the verb that just ran — circular advice, though not a REQ-514 violation since that rule governs `next_argv`), M4 (seven of ten contract predicates are exact-phrase pins; the Batch Constraint is still met because the Go RED test is the behaviour half), M6 (P-A-U boxes were unchecked; the hand-back shows the equivalent work).

**What the review affirmed:** the per-record folding, the REQ-514-conformant set-aside shape with empty `next_argv` and preserved reason codes, the honest RED→GREEN, and the `setAsideReasonCode` cross-artifact predicate. The plain `do-work run` path in the Red-Green Proof does what it promises.

Remediation dispatched for F1, F2, F3, M1, M3 and M5, with `internal/lifecycleadvance/recovery_commands.go` authorized into scope because the REQ's completion proof requires it.

## Remediation

Dispatched after the failed review; merged at `fe2de1e`, keeping `<pre>` at `18666d7` so the range stays cumulative.

**F1 fixed — a set-aside REQ keeps its claim.** `handleRecover` now collects the REQ ids whose finalization records carry `FINALIZATION-SET-ASIDE` and skips claim recovery for each, recording `finalization set aside; claim preserved`. Stopping `recover` outright was the alternative and was rejected: that is the whole-run gate this REQ exists to retire. The exclusion is read from the record's own `reason_codes`, so there is no new field and no second contract.

The orchestrator verified this independently rather than accepting the report: with `recovery_commands.go` restored to the pre-remediation revision, `TestRecoverPreservesTheClaimOfARequestFinalizationSetAside` fails with `the set-aside REQ lost its claim: stat .../do-work/working/REQ-730.md: no such file or directory` — the reviewer's symptom exactly — and passes on the merged tree.

**F2 fixed — shared-cause refusals lose REQ ownership at the producer.** `FINALIZATION-DIRTY-INDEX` and `FINALIZATION-AMBIGUOUS-SHARED-STATE` are exactly the refusals whose cause is state the journal never declared: the repository's single index, and shared lifecycle, release or protected paths outside the recovery group. They now emit through `sharedStateRefusal` with no `AffectedIDs`, so the existing ownership gate returns false and the run stops. Keying it on the condition — no REQ owns the cause — rather than on a code list that would go stale on the third such code is what makes the "a finding that owns no REQ stops the run" branch reachable from the journal path for the first time. It also matches `discoveryRefusal`, which has always produced unowned refusals for the same reason. The unowned refusal's resolving verb is `uncommitted-inventory`, copied from that established shape, so REQ-514's rule that a refusal never names itself still holds.

**F3 and M1 fixed — both false branches now pinned, red-proven.** `TestRecoverFinalizationStopsOnSharedDirtInsteadOfSettingOneRequestAside` fails without its fix with outcome `success` and *both* REQs carrying `FINALIZATION-DIRTY-INDEX, FINALIZATION-SET-ASIDE` — the old code set the whole repository's dirty index aside once per REQ. `TestRecoverFinalizationStopsWhenRollbackLeavesResidue` fails without its guard with REQ-721 dragged into a broken state. `namesRequestID` became the exclusive `namesOnlyRequestID`, closing M1.

**M3 and M5 fixed.** The exit-summary line dropped `[title]`, which is not a field of `FinalizationResult`, rather than sending a floor agent to open the REQ file for one line. The contract's missing-action-file branch is guarded on the file existing instead of on an unrelated variable.

## Decisions (remediation)

**D-07 — DECIDE & STATE. Claim recovery skips set-aside REQs; it does not refuse for them.** Stopping `recover` on a set-aside would reinstate the whole-run gate this REQ removes.

**D-08 — DECIDE & STATE. Shared-cause refusals lose ownership at the producer, not at the consumer.** Stated as a condition rather than a code list, per CLAUDE.md.

**D-09 — DECIDE & STATE. The unowned refusal's resolving verb is `uncommitted-inventory`**, the established shape for a global stop, and a genuinely different verb from the one that just ran.

**D-10 — DECIDE & STATE, orchestrator-authorized. `internal/lifecycleadvance/` was written despite the first pass declaring it out of scope.** F1 cannot be fixed anywhere else — the release happens in `handleRecover`, not in the finalization package — so the REQ's completion proof requires the file class. Added to the scope list and `write_set`.

**D-11 — DECIDE & STATE. `SetAsideReasonCode` is exported** so `lifecycleadvance` uses the same token the finalization package emits, rather than a second literal in a second package. The contract's extraction follows the rename in the same commit, so it still reddens on a future rename.

**D-12 — DECIDE & STATE. The exit-summary line dropped `[title]` rather than sourcing it.**

**D-13 — DECIDE & STATE, orchestrator-authorized. Two restatement sites outside the declared write set were corrected.** `work-reference.md` Crash Recovery and `docs/work-guide.md` both said `--assume-sole-authority` resets every working claim. That is false as of this commit and false *because of* it, so the same rule that pulled `commit.md` into scope pulls these in. One clause each.

**D-14 — DECIDE & STATE. `recover --take-over REQ-NNN` on a set-aside REQ preserves the claim instead of taking it over.** The same command's finalization pass just refused that REQ's tail, so handing the claim to the run would dispatch a builder onto a REQ whose stale journal refuses again at finalize. Both the set-aside finding and the claim decision appear in the result, so the user sees why.
