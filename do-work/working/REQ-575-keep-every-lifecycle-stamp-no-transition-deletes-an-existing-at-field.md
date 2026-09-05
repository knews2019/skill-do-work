---
id: REQ-575
title: '[impact-rule-change] Keep every lifecycle stamp: no transition deletes an existing *_at field'
status: claimed
created_at: 2026-09-04T23:52:00Z
user_request: UR-116
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T12:16:57Z
  basis:
    - trivial short-circuit
depends_on: []
related: [REQ-570, REQ-562, REQ-576]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
write_set:
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/actions/work-reference.md
claimed_at: 2026-09-05T12:00:55Z
dispatch_at: 2026-09-05T12:06:43Z
route: A
builder_handback_at: 2026-09-05T12:29:52Z
integration_at: 2026-09-05T12:29:52Z
---

# Keep Every Lifecycle Stamp: No Transition Deletes an Existing `*_at` Field

## What

Make REQ frontmatter timestamps append-only. Once a `*_at` field exists on a request, no lifecycle transition deletes it or writes a different value over it. A transition that re-enters a phase writes its stamp only when the field is absent. Delete the two code paths that still remove stamps: the recover transition's field-strip loop in `state_apply.go` and the gate-deferral parent edit in `defer_gate.go`. State the rule once in the Request File Schema so the next writer inherits it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Both deletion sites were named in the request, so no discovery was needed. The builder settled on satisfying the suffix condition structurally — no stamp list left to add to by mistake — with one helper carrying the rule, the schema stating it once, and a table-driven test reading the stamp set out of the fixture. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-575-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`afc30a9f`). Seven files, six of them the declared write set plus the selector pair the request's Builder Guidance authorized (D-03, accepted as D-06).
- [x] **[UNIFY]:** `git diff --stat` reviewed (7 files, +186/-14); `gofmt -l .` empty; `go vet ./...` clean; debug-artifact scan over added lines empty; each file reviewed for what it kept as well as what it changed, with the per-file checks listed in the hand-back.

## Why

REQ-505 (moving selection and claim behind `advance`) was claimed at 2026-09-04 16:39 UTC and worked until 17:26. At 21:02 the green heavy-verification result moved it from `pending-heavy-testing` back to `pending` and deleted `claimed_at` (commit d4e4a985). The re-claim at 23:00 and completion at 23:01 left the board card saying "wall time 1m 23s" for six hours of work, and the detail drawer showing "-6h 10m wall since Claimed" on the Planning row. The stamps are the only durable timing record the suite has, so deleting one destroys evidence that nothing else holds. REQ-570 (deleting the pending-heavy-testing status) removes that specific transition; this REQ is the general rule so no future hold, requeue, or deferral path can repeat it.

## Context

Current deletion sites, both in the CLI:

- `internal/requeststate/state_apply.go` `TransitionRecover` deletes `claimed_at`, `route`, `planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, `release_at` when it returns an interrupted request to `pending`.
- `internal/publication/defer_gate.go` deletes `claimed_at` on the parent when a gate failure is deferred to a repair REQ.

The board's phase breakdown (`queue-kanban/durations.go` `buildPhaseBreakdown`) already renders stamps in declared pipeline order and shows reversed bookkeeping instead of hiding it, so keeping an old stamp through a recover does not need new display work. `route` and `write_set` are not stamps and stay under their current rules.

## Detailed Requirements

- The condition is the rule: any frontmatter field whose name ends in `_at` is written at most once by the lifecycle and never deleted by it. Do not encode this as a list of ten field names; key it on the suffix so a stamp added later is covered.
- `TransitionRecover` keeps every existing `*_at` stamp. It may still reset `status`, `status_changed_at` (which is itself a stamp the transition legitimately advances; see the next point), `route`, `write_set`, and the generated recovery sections as it does today.
- `status_changed_at` is the one stamp that records "the most recent status change" by definition, so it keeps being overwritten. Name this exception once next to the rule.
- Claim on a request that already carries `claimed_at` leaves the existing value; the drawer's phase order then shows the true first claim and the re-entry is visible from `status_changed_at` and git history.
- Gate deferral leaves `claimed_at` on the parent.
- Record the rule in `actions/work-reference.md` Request File Schema in one or two sentences beside the Timestamp rule. Delete any sentence there that says recovery strips phase observations.

## Constraints

- No new timing file, stream, or writer. REQ-562 (recording lightweight per-REQ lifecycle timings) owns command-level attribution.
- Do not touch the doctor's session-start stamp repair or `scripts/audit-archive-timestamps.sh`; correcting a detectably wrong stamp from git history is a different operation from a lifecycle transition and stays allowed.
- Keep `route` and `write_set` behavior in recover unchanged.

## Builder Guidance

Certainty: firm on "never delete a stamp"; the user confirmed at verify that `status_changed_at` is the one stamp a transition may overwrite; firm that a re-claim keeps the first `claimed_at` because the user's goal is true wall time from first claim to completion. The board's calibration span (estimate versus actual) also reads `claimed_at`, so after a recover-and-re-claim it will measure from the first claim; that is the intended reading. If the builder finds a reader that breaks when a recovered request keeps its old stamps, fix the reader, not the writer.

## Red-Green Proof

**RED prompt/case:** Plan and apply `TransitionRecover` on a claimed request fixture carrying `claimed_at: 2026-09-04T16:39:30Z` and `planning_at: 2026-09-04T16:49:45Z`. Today the applied document has neither field. Plan and apply a gate deferral on a claimed parent fixture; today the parent loses `claimed_at`.
**Why RED now:** `state_apply.go` deletes the ten named stamp fields on recover, and `defer_gate.go` deletes `claimed_at` on the parent.
**GREEN when:** Both applied documents still carry the original `claimed_at` and `planning_at` values byte for byte, `status` and `status_changed_at` changed as before, and a table-driven test over the recover fixture asserts that no `*_at` field present before the transition is absent or different after it, except `status_changed_at`.
**Validation:** User confirmed (verify-requests, 2026-09-05)

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 7300 tokens, over the 2000-token budget and `slugged: partial`; matched because the change edits lifecycle state writers and the `lifecycle-section-evidence` family.
- `_dev/primes/lessons-action-files.md` — 4362 tokens, over budget and `slugged: partial`; matched because the change edits a status contract and its schema prose.

## Full Context

See `do-work/user-requests/UR-116/input.md` for the verbatim input and the REQ-505 trace.

---
*Source: "capture a req for append-only stamps and the board wall time change"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both deletion sites by file and function (`state_apply.go` `TransitionRecover`, `defer_gate.go` parent edit), declares the five-file write set, and carries a captured RED/GREEN pair. Nothing about the location or the pattern needs discovery, so exploration and planning are skipped.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*

## Decisions

- **D-06 — the selector's two files were added to this REQ's write set, by the orchestrator, before integration. DECIDE & STATE.** The builder reported D-03 as a scope expansion rather than writing it silently: keeping `claimed_at` through a recover made `internal/nextselection` veto every recovered request as `ALREADY-CLAIMED`, so the writer change alone would have made recovery a one-way door. The REQ's own Builder Guidance settles the direction — "If the builder finds a reader that breaks when a recovered request keeps its old stamps, fix the reader, not the writer" — so the expansion is what the request asked for, not drift. `write_set` now carries `internal/nextselection/next_selection.go` and its test.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)

**What was done:** A new helper states the append-only rule once and keys it on the field-name suffix; the claim transition writes its stamp through it, so a re-claim keeps the first value. The recover transition's ten-name strip loop is gone, replaced by a single field deletion for the route, which is what makes the suffix condition structural rather than a longer list. Gate deferral no longer deletes the parent's claim stamp. The request file schema states the rule beside the Timestamp rule and marks each of the four fields that carry current state rather than history. The selector was corrected under the request's own guidance so a recovered request carrying its old claim stamp stays selectable. Merge range `cd686ed7..a2c6f4cf`; builder branch head `afc30a9f`. Builder-authored `## Decisions` (D-01 to D-05) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-05-120117/REQ-575-handback.md`; the orchestrator's D-06 is in this file.

## Qualification

**Passed.** Read from the merge range `cd686ed7..a2c6f4cf`.

- The rule is a condition, not a list. `state_apply.go` gains one helper that writes a stamp only when the field is absent, and recover's ten-name deletion loop is replaced by a single route deletion. There is no enumeration left for a future stamp to be forgotten from, which is exactly what the request asked for.
- Recover still does everything else it did: status, the status-change stamp, route, write set, and the generated recovery sections. The claim's withdrawal of a prior attempt's commit and heavy-verification fields is untouched.
- `defer_gate.go` lost exactly one deletion line and kept the parent's move to the queue with its deferral fields and history entry.
- The selector change is one condition: a frontmatter claim stamp counts as live-claim evidence only while the status says claimed. Checkpoint-writer evidence is untouched and still vetoes on its own. Without it, keeping the stamp would have made every recovered request permanently unselectable — caught by an existing test in another package, not by the new ones.
- The schema paragraph names its exceptions instead of asserting an absolute rule, and each exception is also marked where its field is defined. A schema sentence claiming transitions never remove a stamp, in a package that removes three, would mislead the next writer.

Requirements traced: the suffix condition rather than a name list, recover keeping every stamp, the status-change stamp still advancing, a re-claim keeping the first claim stamp, gate deferral leaving the parent's stamp, and the rule recorded in the schema. The three deletions the builder kept are judged below.

**Judgment on the builder's two escalations:**
- **D-01 accepted.** The two remaining deletions withdraw a field together with the state it describes — the heavy-verification stamp with the commit it verified, the blocked stamp with the condition that blocked it — rather than discarding timing history, which is what the request set out to protect. Deleting the heavy-verification stamp is an active guard that stops dependents building against work a remediation withdrew, shipped one day earlier. The request's Context lists only two deletion sites, which suggests it was written before that guard existed. This goes to the user as a follow-up rather than being settled here.
- **D-03 accepted as D-06.** The request's Builder Guidance names this exact case and says to fix the reader. The builder reported it instead of writing it silently, and the write set was extended before integration.

*Checked by work action*

## Testing

**Focused tests (builder, in its worktree at `afc30a9f`):**
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./internal/requeststate/... ./internal/publication/...` — pass.
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./...` — pass, 760 tests, wall 59s, slowest file `internal/finalization/finalization_recovery_test.go` at 22.78s against the 30s per-file budget.
- `go test -count=1 ./...` in the module — exit 0, 30 packages ok.
- `bash _dev/tests/contracts/core-checks.sh`, `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/select-simple-reqs-behavior.sh` — all exit 0, run because the change edits shipped action prose and the selector.

**Load note worth keeping:** the builder's first whole-module run reported three files over the 30-second per-file budget with no failing test — 42.63s, 40.25s and 32.75s — and the same three files came in at 22.78s and below on an immediate re-run. All three are untouched by this change. That is machine load from the parallel builders, not this diff, and it is the same shape the 5 September run record already documented.

**Red-green validation** (traced to `## Red-Green Proof`): RED — `TestRecoverAndReclaimPreserveEveryLifecycleStamp` failed with one subtest per stamp, nine in all, each naming the field recover deleted and its original value, and the dumped post-recover document carried only the status fields; `TestDeferGateCreatePublishesOneAtomicDependencyLifecycle` failed at `defer_gate_test.go:39` with an empty claim stamp where the fixture had one. GREEN — both pass, every one of the nine stamps survives byte for byte, the status-change stamp is the recover instant, route and write set are gone, and a re-claim two hours later leaves the original claim stamp in place.

