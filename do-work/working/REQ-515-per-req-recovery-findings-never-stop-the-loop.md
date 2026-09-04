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
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/run-with-recovery.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/commit.md, _dev/tests/contract-regressions.sh, _dev/tests/contracts/recovery-set-aside.sh, skills/do-work/tools/do-work-cli/internal/finalization/]
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
