---
id: REQ-570
title: '[impact-rule-change] Delete the pending-heavy-testing status; held requests stay claimed'
status: claimed
created_at: 2026-09-04T22:52:00Z
user_request: UR-114
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-507]
related: [REQ-571]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/cleanup.md
  - skills/do-work/actions/roadmap.md
  - skills/do-work/actions/restart-with-parallel-handoff.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go
  - skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
claimed_at: 2026-09-04T23:00:17Z
route: C
estimate:
  p50_active_minutes: 75
  confidence: low
  calculated_at: 2026-09-04T23:03:33Z
  basis:
    - Route C
    - 24-file write set
    - 7 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
exploration_at: 2026-09-04T23:07:12Z
planning_at: 2026-09-04T23:11:07Z
preflight_at: 2026-09-04T23:14:45Z
dispatch_at: 2026-09-04T23:15:08Z
builder_handback_at: 2026-09-04T23:35:48Z
---

# Delete the pending-heavy-testing Status; Held Requests Stay Claimed

## What

Make the heavy-test hold a phase of a `claimed` request instead of a status. A request is held only after fast tests, qualification and review have passed; it stays `claimed` in `do-work/working/` with its `## Heavy Verification Plan` section and its `commit:` on `main`. At queue exhaustion the same session runs the lanes and, in the same turn, finalizes each green request through Steps 8 and 9 or enters ordinary remediation for each red one. Delete the `pending-heavy-testing` value and every reader of it in the core skill.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A green heavy answer flips the request from `pending-heavy-testing` to plain `pending`. `DependencySourceReady` (`schema_normalization.go` line 240) counts a dependency as landed only when it is terminal, or `pending-heavy-testing` with a commit. So a request is more available to its dependents while it waits for slow tests than after it has passed them. Observed on 2026-09-04: REQ-504, REQ-505 and REQ-506 all merged (commits `f412a841`, `716187b8`, `06367337`, each an ancestor of `main`) and all green at 20:57 to 21:02 UTC, yet the selector offered only REQ-504 with `resume_phase: review` and excluded the other three as `DEPENDENCIES-UNMET`. Three independent reviews were forced into a chain, each costing a re-claim commit and a selection pass, and the board showed `pending` for merged and verified work. `actions/work.md` line 93 already states that intermediate phases are tracked by which sections exist in the request file, not by status changes. The hold broke that rule; the return trip through `pending` is the damage.

## Context

- Status flow today (`actions/work.md` line 91): `pending → claimed → pending-heavy-testing → pending → claimed → completed`. Target: `pending → claimed → completed`, no detour.
- The hold is written under **Qualification and Testing Judgment** in `actions/work.md` (line 344 onward); the drain at exhaustion is the paragraph starting **Heavy-lane drain**; the manual path is `actions/clarify.md` Step 2.5.
- The green and red answers are the `heavy-testing` mode in `internal/publication/answer.go` (line 254 onward) with refusal codes `ANSWER-HEAVY-STATUS`, `ANSWER-HEAVY-EVIDENCE-*`, `ANSWER-HEAVY-ARCHIVE-FORBIDDEN`; `HeavyTestingEvidence` and `HeavyLaneResult` live in `publication_types.go`.
- The review resume is `matchingHeavyReviewPhase` and `ResumePhase` in `internal/nextselection/next_selection.go` (line 322 onward).
- Other readers of the value in the core skill: `advance_commands.go` line 110 (skip case), `checkpoint_commands.go` line 131 (queue_state counts), `result_model.go`, `schema_normalization.go` (canonical enum), `actions/cleanup.md` (active-status lists), `actions/roadmap.md` (Needs clarification bucket), `actions/restart-with-parallel-handoff.md` (clarify hint), `actions/work-reference.md` (schema), `prime-do-work-cli.md`. Find them by search, not by this list.
- The board (`skills/do-work-board`) also reads the value; that is REQ-571, which depends on this one.
- Finding and proposal: `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html` section 03, lessons L24 and L25.

## Detailed Requirements

- Reorder `actions/work.md`: Step 7 (review) and Step 7.5 (lessons) run before the heavy hold. A request with selected heavy lanes is held only after fast tests, qualification and review have all passed. The hold appends `## Heavy Verification Plan` and keeps the request `claimed` in `do-work/working/`; it writes no status change and no `## Open Questions` machine line.
- The heavy-lane drain at queue exhaustion runs the unioned lanes once at HEAD as today. For each green request the same session runs Step 8 (prepare finalization, with the release version assigned at that moment) and Step 9 (the single finalization transaction) in the same turn. For each red request the session enters ordinary remediation now, then re-holds after a fresh review. No `answer` transaction is involved in either outcome.
- A skipped lane leaves the request claimed with a typed `HEAVY-RUN-LANE-SKIPPED` finding; the next drain retries; the composed exit summary names any request still held when the run ends. Delete `actions/clarify.md` Step 2.5 and its `heavy-testing` references; clarify handles `pending-answers` and `blocked` only.
- `DependencySourceReady` accepts `claimed` with a nonblank `commit:` in place of `pending-heavy-testing` with one. The dependency graph node keeps carrying the commit; no new field.
- The recover step recognizes a `claimed` request that has a `## Heavy Verification Plan` section, a `commit:` that is an ancestor of HEAD, and no live writer, and routes it to this session's exhaustion drain as held work rather than as interrupted build work. It reads the section, never the clock. Recovery findings say so with a typed reason code.
- Delete the `heavy-testing` answer mode, its refusal codes, and the `HeavyTestingEvidence` manifest types; delete `ResumePhase` and `matchingHeavyReviewPhase` from the selector and its typed result; delete the value from the canonical status enum, the `advance` skip case, the checkpoint `queue_state` counts, the result model, and every action or reference sentence that names it. Delete the matching predicates in `_dev/tests/contracts/core-checks.sh` in the same commit.
- Keep `heavy_verified_at` and `heavy_verified_revision` as durable evidence written by the drain onto the request record before finalization, so the archived record still proves which revision the lanes checked.
- Finalization accepts the held record from `do-work/working/`, where it already is; nothing moves through `do-work/queue/`.

## Constraints

- Land after REQ-504 to REQ-507 have closed through the existing path; do not add a compatibility read for the old status. If a queued or working request still carries the old value when this REQ is claimed, stop and report it as a typed finding instead of migrating it by hand.
- Delete before you add: no new status value, no new checkpoint count, no manual lane path. Held requests count as in progress; that is accurate and the plan section explains it.
- One step per REQ, never a rewrite of `work.md`. Judgment stays prose; the CLI emits typed findings, never paragraphs.
- The floor agent (reads files, runs shell) must still complete a run using only `advance` output plus the remaining prose.
- The board's reader is out of scope here; it is REQ-571.

## Dependencies

- Depends on REQ-507 (hand the archive and commit tails to `finalize`): the drain finalizes green requests through the canonical finalization that REQ-507 puts behind `advance`. REQ-507 transitively covers REQ-506, which produces the heavy plan from `advance`.
- REQ-571 (board reader cleanup) depends on this REQ.

## Builder Guidance

- High certainty on the design: the user chose "leave them claimed" over a widened readiness rule (P2) and over a new `pending-review` status (P3). Do not reintroduce either.
- Two decisions are already made and recorded above: skipped lanes stay claimed and retry on the next drain; held requests count as in progress with no separate count.
- Expect the change to be mostly deletions. If a deletion leaves a reader with no case for a held request, the answer is "it is `claimed`", not a new branch.
- Search for the value across the core skill before declaring the sweep complete; the file list in Context is illustrative.

## Red-Green Proof
**RED prompt/case:** Fixture repository with REQ-A `claimed` in `do-work/working/`, carrying `commit:` on HEAD's ancestry and a `## Heavy Verification Plan` section, and REQ-B `pending` with `depends_on: [REQ-A]`. Run `do-work-cli --format json next`. Separately, submit an `answer` manifest with `mode: heavy-testing`.
**Why RED now:** `next` excludes REQ-B as `DEPENDENCIES-UNMET` because `DependencySourceReady` rejects `claimed`; the `heavy-testing` answer mode is accepted; `recover` reports the held REQ-A as an interrupted claim.
**GREEN when:** `next` selects REQ-B with REQ-A's commit as its landed source; the `heavy-testing` mode is refused as `ANSWER-MODE-INVALID`; `recover` returns a typed held-for-heavy-lanes record for REQ-A that routes it to the drain; the string `pending-heavy-testing` no longer appears in the core skill's shipped files or contract predicates; the full Go module, contract regressions and the direct maintainer gate pass.
**Validation:** Inferred during capture from the live selector output and the user's approval of the design.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4362 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal; matched on "status contracts, downstream readers".
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6798 tokens, over budget and `slugged: partial`; matched on "condition-based classifiers, semantic recovery completeness, lifecycle-section-evidence".
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`; matched because `_dev/tests/contracts/core-checks.sh` changes.

## Full Context
See `do-work/user-requests/UR-114/input.md` for complete verbatim input.

---

## Triage

**Route: C** - Complex

**Reasoning:** The change deletes a status across six Go packages, five action files, the schema layer and the contract predicates, reorders the review and hold steps in the work action, and adds a recover branch; it touches several subsystems with a design already fixed, so it needs a plan and an exploration before one builder runs.

**Planning:** Required

## Exploration

Precondition holds: no queued or working REQ carries `status: pending-heavy-testing`; the value survives only inside the generated `queue_state:` string in `do-work/CHECKPOINT.md`.

**Go readers and writers of the value (core skill):** enum entry and `DependencySourceReady` disjunct in `schemanormalization/schema_normalization.go` (lines 30, 240-244); the `summarizeQueue` count in `nextselection/next_selection.go` (407-408) and the `PendingHeavyTesting` field in `resultmodel/result_model.go` (272, 1179-1182); the `heavy-testing` answer mode, its only writer of the status, and every `ANSWER-HEAVY-*` code in `publication/answer.go` (254-292, 610-698) with `HeavyTestingEvidence`/`HeavyLaneResult` in `publication_types.go` (81-96); the queue skip case in `lifecycleadvance/advance_commands.go` (110-112); the checkpoint count in `checkpoint_commands.go` (129-131). `matchingHeavyReviewPhase` and `ResumePhase` live in `next_selection.go` (315, 322-338) and `result_model.go` (228, 1193-1194). `heavy_verified_at`/`heavy_verified_revision` are parsed by `requestmodel/request_model.go` (57-58, 315) and today written only by `answer.go`.

**Advance phase table:** `advance_commands.go` classifies a working request by which `##` sections exist (Triage, estimate, Plan, Exploration, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation, then `finalize`). `## Heavy Verification Plan` has no Go writer or parser and is not read by `advance`; it is written and compared by prose only (`work.md` 347, 351; `clarify.md` 39). `## Heavy Verification Result` is written only by `answer.go` (684-698).

**Recover seam:** `recovery_commands.go` → `recoverableWorkingRequests` (193-209) qualifies a working request when status is `claimed`, or `blocked` with `blocked_by`; each becomes a `RecoveryClaimResult` (`result_model.go` 529-535: RequestID, RequestPath, CheckpointEvidence, Decision, Recovered) with finding codes `RECOVERY-TAKEOVER-AVAILABLE`, `RECOVERY-CLAIM-RECOVERED`, `RECOVERY-NONE`. No clock or process probe exists; the only writer evidence is the checkpoint `writer:` label. The held-for-heavy-lanes classification branches here.

**Heavy planner and runner:** `heavyverification/heavy_commands.go` registers `plan-heavy-verification`, `plan-heavy-revalidation`, `run-heavy-verification`; the runner emits `HEAVY-RUN-LANE-SKIPPED` (`heavy_run.go` 363) and red findings; default manifest `_dev/tests/heavy-lanes.json`.

**Tests pinning the old behavior:** `schema_normalization_test.go` (enum sentence at 41; `DependencySourceReady` at 162, 165); `dependency_graph_test.go` (`TestPendingHeavyDependencyIsSourceReadyUntilItReturnsToPending`, `TestPendingHeavyDependencyRequiresImplementationCommit`, helper branch at 24); `next_selection_test.go` (`TestPendingHeavyTestingIsCountedAndNeverSelected`, `TestMatchingHeavyEvidenceResumesAtReviewAndStaleEvidenceDoesNot`, `TestPendingHeavyTestingSourceAllowsDependentSelection`); `answer_test.go` (four `TestHeavyTestingAnswer*` tests and `newHeavyAnswerFixture`). Fixture helpers to reuse: `buildFixtureGraph`, `writeCommandRequest`, `writeFixture`/`canonicalREQFixture`/`initializedGitRepository`, `writeAdvanceRequest`, `legacyCheckpointRepository`.

**Prose passages:** `work.md` 16, 91, 97, 344-372, 487, 501; `work-reference.md` 14, 229, 233-244, 282, 319-323, 671-681; `clarify.md` 11, 17, 31, 35, 37-52, 60; `cleanup.md` 47, 59, 293; `roadmap.md` 54, 68; `restart-with-parallel-handoff.md` 64; `prime-do-work-cli.md` 27 (one clause). Contract shell tests pin none of these sentences and never name the status, so the write set's `_dev/tests/contracts/core-checks.sh` entry has nothing to delete.

*Generated by Explore agent; full findings consumed from the run directory.*

## Plan

Design is fixed by UR-114. The full reader inventory with line numbers is `REQ-570-exploration.md` in this run directory; this file does not repeat it.

## Findings the orchestrator must see before dispatch

- **F1 — dependency REQ-507 is `pending`, not terminal.** `do-work/queue/REQ-507-*.md` has `status: pending`, `commit: ad8bceb7`, and a `## Heavy Verification Result` (red lane). Its code and its Step 8/9 prose are on `main` (`advance` already returns the `finalize` phase in `advance_commands.go`, and `finalization_gate.go` composes `FinalizeBound`). REQ-570 was claimed by explicit naming, which bypasses `depends_on`. Nothing in this plan needs REQ-507's remaining remediation. Record this as a typed finding in the REQ, do not stop.
- **F2 — no queued or working request carries the old status.** The only hits under `do-work/` are REQ-571's filename and title, REQ-570 itself, and the generated `queue_state:` string in `CHECKPOINT.md`, which `advance --checkpoint` rewrites. The "stop and report" constraint is not triggered.
- **F3 — obstacle: a stale `commit:` widens readiness under the new rule.** Neither the `claim` transition nor the `recover-claim` transition strips `commit:` (`requeststate/state_apply.go` lines 516-524 and 526-549). Today `pending`+commit is not ready, so a stale commit is harmless. With the new rule (`claimed` + nonblank `commit:` is ready), re-claiming REQ-506 or REQ-507 makes their dependents source-ready at claim time, against the previous attempt's commit, before any new work lands. The same hole opens on a red drain: the request stays `claimed` with its red `commit:`. **Smallest resolution:** the `claim` transition deletes `commit`, `heavy_verified_at`, and `heavy_verified_revision` (three field names in the existing delete loop shape, one test), and the red branch of the drain prose deletes `commit:` and `## Heavy Verification Plan` from the record before remediation re-dispatch, using the same writer that wrote them at the hold. This adds `internal/requeststate/state_apply.go` and its test to the write set. It is an addition earned by a concrete case (two such records are in the queue now). If the orchestrator declines, record the decision and the risk in `## Decisions`.
- **F4 — no contract predicate names the status.** `grep -n 'heavy\|resume_phase\|Step 2.5' _dev/tests/contracts/core-checks.sh` returns nothing; `contract-regressions.sh` uses `heavy` only as its own tier name. The REQ's "delete the matching predicates in `core-checks.sh`" is satisfied by absence; remove `_dev/tests/contracts/core-checks.sh` from the REQ's `write_set` if it stays untouched and say so in the Implementation Summary.
- **F5 — who writes `heavy_verified_*` and `## Heavy Verification Result` after the `answer` mode is gone.** The drain writes them onto the claimed record the same way Step 7 stamps `review_at` and appends `## Review`: an action-owned edit of the record in `do-work/working/`, before the Step 8 manifest binds the request preimage. No CLI change. The runner's JSON (`run-heavy-verification`) supplies the values.

## Architectural decisions

- **D1 — held is a phase, detected by section plus ancestry.** A claimed request in `do-work/working/` with a `## Heavy Verification Plan` section and a `commit:` that is an ancestor of HEAD is held. No field, no status, no clock.
- **D2 — `DependencySourceReady(status, commit)` keeps its signature.** Only the second disjunct changes to `normalizedStatus == "claimed" && commit != ""`. `dependencygraph` and `duplicateStatusesSatisfied` already pass the node's commit; they need no code change.
- **D3 — the selector loses `ResumePhase` entirely.** A held request is in `working/`, which `Select` never considers (only `TreeSection == "queue"` records are candidates). Nothing replaces `matchingHeavyReviewPhase`; review already ran before the hold.
- **D4 — recovery preserves a held claim instead of releasing it.** The check sits after the authority check in `handleRecover` (`recovery_commands.go` line 92 onward), so an unauthorized session still gets `RECOVERY-TAKEOVER-AVAILABLE`. When authorized and held: no `requeststate` plan runs, `Decision` is `held for heavy lanes; claim preserved`, `Recovered` stays false, a new bool `HeldForHeavyLanes` (`json:"held_for_heavy_lanes"`) on `RecoveryClaimResult` is set, and an info finding `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` carries the section name and the commit as evidence with `NextArgv` `do-work-cli --format json advance REQ-NNN`. The runner argv cannot be derived by the CLI because the lane manifest path is action-owned. Held with a commit that is **not** an ancestor of HEAD falls through to ordinary recovery: the plan section alone is not proof of landed work.
- **D5 — `advance` needs no new phase.** A held record already has Review, Lessons Learned, and Orientation, so `classifyWorkingAdvance` returns `finalize`, which is exactly the green drain's next step. `advanceSections` ignores the unknown `Heavy Verification Plan` heading. Reuse `advanceSections` + `hasSection` from the same package in D4.
- **D6 — the hold and drain prose move to one named subsection between Step 7.5 and Step 8** (`### Step 7.7: Heavy-Test Hold and Drain`). Qualification and Testing Judgment keeps one sentence pointing there. This is a move of two existing paragraphs, not a rewrite; the floor agent needs the position to read in step order.
- **D7 — the ancestor check in recovery is a local three-line helper** (`git merge-base --is-ancestor <commit> HEAD` via `exec.Command`, the shape `finalization_journal.go` line 273 uses). `lifecycleadvance` already shells out to git in `checkpoint_commands.go`; do not import `nextselection` or `publication` for their unexported helpers.

## Files to modify, in compile-safe order

Run `go build ./... && go vet ./...` after every group.

1. **`internal/nextselection/next_selection.go`** — delete `ResumePhase: matchingHeavyReviewPhase(...)` from the `SelectionRecord` literal (line 315), delete `matchingHeavyReviewPhase` (322-337), delete the `case "pending-heavy-testing"` in `summarizeQueue` (407-408). Delete `selectionRevision` / `selectionRevisionIsAncestor` and the `time` import only if nothing else in the package uses them (grep first). **`next_selection_test.go`** — see Testing.
2. **`internal/resultmodel/result_model.go`** — delete `ResumePhase` (228) and `PendingHeavyTesting` (272); delete the two `resume phase` render lines (1193-1195); drop the `PendingHeavyTesting` term from the `queue:` condition and format string (1179-1182); fix the comment at 97-100 that cites publication's `HeavyLaneResult`.
3. **`internal/schemanormalization/schema_normalization.go`** — remove the enum value (30); rewrite `DependencySourceReady` per D2 and its doc comment (236-244).
4. **`internal/dependencygraph/`** — no production change. **`dependency_graph_test.go`** — see Testing.
5. **`internal/publication/answer.go`** — delete the `case "heavy-testing"` block (254-293), `validateHeavyTestingEvidence` (610-666), `appendHeavyVerificationResult` (684-697), and `resolveAnswerCommit` / `answerRevisionIsAncestor` if unused afterwards; change the `ANSWER-MODE-INVALID` message to `answer mode must be clarify, stakeholder, or verify-repair` (352). **`publication_types.go`** — delete `HeavyLaneResult`, `HeavyTestingEvidence` (81-95) and `AnswerManifest.HeavyTesting` (105). **`answer_test.go`** — see Testing.
6. **`internal/lifecycleadvance/advance_commands.go`** — remove the value from the skip case (110). **`checkpoint_commands.go`** — remove the count from the format string and lookup (129-131). **`recovery_commands.go`** — add the held branch per D4 and D7 with a `heldForHeavyLanes(repositoryRoot, request) bool` helper. **`resultmodel/result_model.go`** — add `HeldForHeavyLanes bool` to `RecoveryClaimResult`. **`recovery_commands_test.go`** — see Testing.
7. **(F3, if accepted) `internal/requeststate/state_apply.go`** — in `TransitionClaim` (516-524) delete `commit`, `heavy_verified_at`, `heavy_verified_revision`. Test in the requeststate test file that covers claim bytes.
8. **Prose** — `skills/do-work/actions/work.md`, `work-reference.md`, `clarify.md`, `cleanup.md`, `roadmap.md`, `restart-with-parallel-handoff.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (details below), plus lessons satellites and `do-work/lessons-index.md` rows.
9. **Release** — shipped files change, so finalization carries a version bump and a changelog entry per `_dev/primes/prime-releases.md` (title says what shipped; mirror byte-identical). No CHANGELOG history is rewritten.

## Testing approach (TDD: write RED first, run, then edit production)

Old tests that pin the deleted behavior and must go or be rewritten (cross-REQ test-break rule: intentional; name them in `## Testing`):

- `schema_normalization_test.go:41` warning string lists the value; `TestTerminalPredicatesKeepFailureAndCancellationDistinct` lines 162-167 pin the old disjunct.
- `dependency_graph_test.go`: `TestPendingHeavyDependencyIsSourceReadyUntilItReturnsToPending`, `TestPendingHeavyDependencyRequiresImplementationCommit`, and the `buildFixtureGraph` special case at line 24.
- `next_selection_test.go`: `TestPendingHeavyTestingIsCountedAndNeverSelected`, `TestMatchingHeavyEvidenceResumesAtReviewAndStaleEvidenceDoesNot`, `TestPendingHeavyTestingSourceAllowsDependentSelection`; `runNextSelectionGit` is used 13 times, keep it.
- `answer_test.go`: `TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure`, `TestHeavyTestingAnswerRejectsNegativeWallSeconds`, `TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane`, `TestHeavyTestingAnswerRejectsMismatchedEvidence`, helper `newHeavyAnswerFixture`; `runGitFixtureOutput` has other callers, keep it.

RED tests to write first, one per Red-Green Proof clause:

- **(a)** `schema_normalization_test.go`: `TestDependencySourceReadyAcceptsClaimedWithCommitOnly` — `("claimed","abc123")` true; `("claimed","")` false; `("pending","abc123")` false; `("pending-heavy-testing","abc123")` false; and `NormalizeField("status","pending-heavy-testing").IsRecognized` false. Update the line-41 expected warning string to the shorter list (RED until the enum changes).
- **(b)** `dependency_graph_test.go`: `TestClaimedDependencyWithCommitIsSourceReady` — extend `graphFixture` with a `commit` string and a tree section so a `claimed` fixture is written under `do-work/working/`; assert the dependent of `claimed`+commit `IsReady`, the dependent of `claimed` without commit is not, the dependent of `pending`+commit is not.
- **(c)** `next_selection_test.go`: `TestHeldClaimedSourceAllowsDependentSelection` — REQ-A `claimed` in `do-work/working/` with `commit:` and a `## Heavy Verification Plan` body, REQ-B `pending` with `depends_on: [REQ-A]`; `Select` returns exactly `[REQ-B]`. RED today: REQ-B is `DEPENDENCIES-UNMET`. `ResumePhase` removal is compile-checked; no decorative JSON-absence test.
- **(d)** `answer_test.go`: `TestAnswerRefusesHeavyTestingMode` — a `pending-answers` fixture with `Mode: "heavy-testing"` returns `Refusal.Code == "ANSWER-MODE-INVALID"`. RED today: the code is `ANSWER-HEAVY-STATUS`.
- **(e)** `recovery_commands_test.go`: `TestRecoverHoldsAClaimedRequestWithAHeavyVerificationPlanForTheDrain` — in-process `handleRecover(..., "--assume-sole-authority")` on a git fixture (shape of `recovery_set_aside_test.go`): REQ-740 `claimed` in `working/` with `commit:` = HEAD and a `## Heavy Verification Plan` section; REQ-741 `claimed` without the section. Assert REQ-740 stays in `working/`, its claim has `HeldForHeavyLanes` true and `Recovered` false, a finding with code `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` names it, and REQ-741 is recovered to `queue/`. Second case in the same test: REQ-742 with the section but a `commit:` from an abandoned branch is recovered to `queue/` (D4 fall-through).
- **(F3, if accepted)** requeststate: `TestClaimStripsPriorAttemptCommitAndHeavyEvidence` — a `pending` fixture carrying `commit:`, `heavy_verified_at`, `heavy_verified_revision` is claimed; the claimed bytes carry none of the three.

`advance_commands_test.go` and `checkpoint_commands_test.go` pin neither the skip case's status list nor the `queue_state` string; the final `grep` is their proof.

## Prose edits (one or two sentences each; judgment stays prose, mechanics point at the CLI)

**`actions/work.md`**
- Line 16 (Do NOT use): drop the held clause; keep `pending-answers` → `do-work clarify`.
- Line 91 (Status flow): reduce to `pending → claimed → completed / completed-with-issues / failed`; add that the heavy hold is a phase of `claimed`, marked by `## Heavy Verification Plan` and `commit:`, per the next paragraph's section-tracking rule.
- Lines 97 (Special statuses): delete the `pending-heavy-testing` bullet.
- Line 343 (Qualification and Testing Judgment, last sentence): "Plan affected heavy lanes with the repository's typed planner; selected lanes are recorded in the Testing section and held at Step 7.7 after review."
- Lines 344-362 (Heavy-test hold, drain, answer-manifest paragraph): move to `### Step 7.7: Heavy-Test Hold and Drain` before Step 8. Hold: after Steps 7 and 7.5 pass, land the implementation commit, record its hash in `commit:`, append `## Heavy Verification Plan`; the request stays `claimed` in `do-work/working/`; no status change, no `## Open Questions` line, no move to `queue/`, no checkpoint edit; recompute selection and continue. Drain: trigger is "no claimable pending REQ remains and a held claimed request exists" (also one routed here by recover's `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES`); keep the drift refusal, the union run, the runner paragraph, and the reuse paragraph as they are; replace the answer-manifest paragraph with: green → write `heavy_verified_at`, `heavy_verified_revision`, and `## Heavy Verification Result` (target revision, execution revision, one line per lane) onto the record, then run Steps 8 and 9 in this turn; red → delete `commit:` and the plan section (F3), enter Step 7's remediation path, re-hold after a fresh review; skipped lane → the request stays claimed, `HEAVY-RUN-LANE-SKIPPED` names the lane, the next drain retries, the exit summary names it. Delete the `clarify` route for drift and historical plans: they are typed findings for a human, never a hand edit.
- Line 487 (checklist Testing judgment): drop "park selected REQs as `pending-heavy-testing`"; add `□ Step 7.7: Heavy hold after review (held requests stay claimed); drain at exhaustion, finalize green in the same turn`.
- Line 501 (Error Handling row): condition becomes "Held claimed REQs remain after the queue is empty"; action points at Step 7.7.

**`actions/work-reference.md`**
- Line 14 (diagram): "skip pending-answers".
- Line 229: parenthetical becomes "also recorded at the heavy hold so dependents can build against the landed source".
- Lines 233-244: delete the `status: pending-heavy-testing` / `status_changed_at` lines and comment; keep `heavy_verified_at` / `heavy_verified_revision` with a two-line comment: written by the drain onto the claimed record before finalization; proves which execution revision the lanes checked.
- Line 282 (status row): remove the value; remove "Qualification and Testing Judgment's heavy-test hold" from the read sites.
- Line 321 (Dependency-source-ready set): "terminally successful, or normalized `claimed` with a nonblank `commit:` — a request held for heavy lanes after review; `claimed` without a commit fails closed, and a red drain withdraws `commit:` so the dependency is unmet again." Keep the "scheduling authority only" sentence.
- Line 671 (exit-summary item 8): condition becomes "any claimed REQ in `do-work/working/` still carries `## Heavy Verification Plan` without `heavy_verified_revision` after the drain ran"; in the render block replace the `clarify` sentence with "Plan drift or a stored historical-revalidation plan is a typed finding for a human; the next `do-work run` drain retries once the cause is fixed."
- Line 823 (Testing Section Template): unchanged.

**`actions/clarify.md`** — delete the Step 2.5 subsection (39-52) and the "Use when" bullet at 11; remove the value from lines 17, 31, 35 (report text becomes "No pending questions or blocked REQs — queue is clear"; "skip Steps 3–5"); delete the Step 2.5 sentence at 60.

**`actions/cleanup.md`** — remove the token from the three status lists (47, 59, 293).

**`actions/roadmap.md`** — line 54 drop the value; line 68 the bucket keys on `pending-answers` only and the `clarify` heavy sentence goes.

**`actions/restart-with-parallel-handoff.md`** — line 64: "to answer questions — `do-work clarify`, included only if some REQ is at `pending-answers`".

**`tools/do-work-cli/prime-do-work-cli.md`** — heavyverification bullet's last sentence: "Per-lane result evidence is written onto the claimed record by the work action's drain before finalization." Advance phase table, Tested row: "review → lessons/orientation → heavy hold when lanes were selected → `finalize` through advance". Add the D4 recover behavior to the `lifecycleadvance` bullet in one clause.

**Lessons** — add one `[family: status-as-phase-marker]` bullet to `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and one to `_dev/primes/lessons-action-files.md` (a hold status that a green result flips back to a weaker status makes work less available after verification than during it; track a hold by its section and commit), and refresh both rows in `do-work/lessons-index.md` in the same commit.

## Verification

From `skills/do-work/tools/do-work-cli/`: `go build ./...`, `go vet ./...`, `go test -count=1 ./...`. From the repository root: `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/shipped-package-reference-contract.sh` (changelog mirror), and the sweep, which must print nothing:

```bash
grep -rn 'pending-heavy-testing\|ResumePhase\|resume_phase\|HeavyTestingEvidence\|HeavyLaneResult\|matchingHeavyReviewPhase\|ANSWER-HEAVY' skills/do-work _dev --include='*.go' --include='*.md' --include='*.sh' | grep -v CHANGELOG
```

The REQ's GREEN also names the direct maintainer gate; run it only where the orchestrator's gate step calls for it.

## Tasks

- **T1 — RED.** Write the tests in Testing (a) to (e), delete or rewrite the listed old tests, run `go test ./internal/schemanormalization ./internal/dependencygraph ./internal/nextselection ./internal/publication ./internal/lifecycleadvance` and record the failures.
- **T2 — Go deletions, groups 1-5.** nextselection, resultmodel, schemanormalization, publication. Build and vet after each group; (a) to (d) turn GREEN.
- **T3 — Go recovery and lifecycle, group 6 (plus group 7 if F3 is accepted).** advance skip case, checkpoint count, recovery held branch and result field, claim strip. (e) turns GREEN; full module green.
- **T4 — Prose, group 8.** All action files, the prime, lessons satellites and index rows; run the sweep grep and both contract scripts.
- **T5 — Record.** Write F1 to F5 and the F3 decision into the REQ (`## Decisions`, `## Discovered Tasks` for the REQ-506/507 remediation note), Implementation Summary, Testing with the cross-REQ test list, then finalization with the release entry.

**Plan validation (work action):**
- Requirement coverage: every Detailed Requirement maps to a task (reorder and hold/drain → T4 prose D6; readiness rule → T2 group 3; recover branch → T3 D4/D7; deletions → T2/T3; heavy_verified fields kept → F5; finalization from working/ → D5). No uncovered requirement.
- No orphan tasks. F3 (claim strips a prior attempt's `commit:` and heavy evidence) is accepted as **D-01** below because REQ-506 and REQ-507 sit in the queue today as `pending` with a stale `commit:`, and the new readiness rule would otherwise make their dependents ready at re-claim.
- Scope sanity: 5 tasks. T1–T3 are one Go pass split by compile order and T5 is orchestrator bookkeeping, so the builder's own work is three tasks; not split into more REQs.
- Consumer field contract: recovery adds `HeldForHeavyLanes` and finding `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` with exact next argv; the drain consumes `run-heavy-verification` JSON fields already defined. Satisfied.
- Lessons satellites and `do-work/lessons-index.md` named in T4 are orchestrator-owned deferred writes (Step 7.5/8); they are removed from the builder's scope below.

*Generated by Plan agent; validated by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — status flow, Step 7.7 hold and drain, checklist, error table
- `skills/do-work/actions/work-reference.md` (modify) — schema, dependency-source-ready set, exit summary item 8
- `skills/do-work/actions/clarify.md` (modify) — delete Step 2.5 and the held-status clauses
- `skills/do-work/actions/cleanup.md` (modify) — three status lists
- `skills/do-work/actions/roadmap.md` (modify) — source list and Needs-clarification bucket
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modify) — clarify hint
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — heavyverification and lifecycleadvance bullets, phase table row
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modify) — enum value, DependencySourceReady
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modify) — RED test (a), enum sentence
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (modify) — RED test (b), old tests removed
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modify) — ResumePhase, matchingHeavyReviewPhase, count case
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modify) — RED test (c), old tests removed
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — ResumePhase, PendingHeavyTesting, HeldForHeavyLanes, renderers
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — heavy-testing mode and helpers deleted
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — RED test (d), old tests removed
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modify) — HeavyLaneResult, HeavyTestingEvidence, AnswerManifest.HeavyTesting
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modify) — queue skip case
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go` (modify) — queue_state count
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modify) — held-for-heavy-lanes branch
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modify) — RED test (e)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modify) — claim strips commit and heavy evidence (D-01)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modify) — RED test for D-01

**Files I will NOT touch:** `skills/do-work-board/**` (REQ-571), `_dev/tests/contracts/core-checks.sh` (no predicate names the status; removed from the write set), lessons satellites and `do-work/lessons-index.md` (orchestrator writes them at Step 8), `VERSION`/changelog mirrors (finalization release), any `do-work/` path.

**Acceptance criteria (restated from REQ):**
- [ ] `next` selects a dependent whose dependency is `claimed` with a nonblank `commit:` (REQ-A claimed in working/ with commit and `## Heavy Verification Plan`; REQ-B depends on it and is selected).
- [ ] An `answer` manifest with `mode: heavy-testing` is refused as `ANSWER-MODE-INVALID`.
- [ ] `recover` returns a typed `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` record for a claimed request with a `## Heavy Verification Plan` section and a `commit:` that is an ancestor of HEAD, and leaves it in working/; without the section, or with a non-ancestor commit, it recovers as today.
- [ ] `claim` strips a prior attempt's `commit:`, `heavy_verified_at` and `heavy_verified_revision` (D-01).
- [ ] The string `pending-heavy-testing` (and `ResumePhase`, `resume_phase`, `HeavyTestingEvidence`, `HeavyLaneResult`, `matchingHeavyReviewPhase`, `ANSWER-HEAVY`) no longer appears in the core skill's shipped files or `_dev` tests, CHANGELOG excepted.
- [ ] `work.md` reorders review and lessons ahead of the hold under a named Step 7.7; the held request stays `claimed` in working/; green drains run Steps 8 and 9 in the same turn; red deletes `commit:` and the plan section and enters remediation.
- [ ] Full Go module, `go vet`, contract regressions and the direct maintainer gate pass.

## Pre-Flight

**Git:** ✓ clean outside `do-work/` at the preflight run; the only `do-work/` dirt was this run's own directory and another session's run scratch, both left alone.
**Tests baseline:** ✓ `go test -count=1` over schemanormalization, dependencygraph, nextselection, publication, lifecycleadvance, requeststate exited 0 (`do-work/working/baseline.json`, launched: true).
**Repository gate:** ✓ direct unpiped `bash _dev/tests/maintainer-verify.sh` exited 0 in 98s (board 381 tests, 18s; CLI 731 tests, 56s; slowest file 22.86s); green-gate record written by `advance` for the exact argv at the current revision.
**Dependencies:** ✓ Go toolchain present; module builds.

*Checked by work action*
