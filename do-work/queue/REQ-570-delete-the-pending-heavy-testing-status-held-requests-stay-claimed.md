---
id: REQ-570
title: '[impact-rule-change] Delete the pending-heavy-testing status; held requests stay claimed'
status: pending
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
  - skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go
  - skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - _dev/tests/contracts/core-checks.sh
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
