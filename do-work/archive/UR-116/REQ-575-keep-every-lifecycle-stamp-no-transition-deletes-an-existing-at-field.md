---
id: REQ-575
title: '[impact-rule-change] Keep every lifecycle stamp: no transition deletes an existing *_at field'
status: completed
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
review_at: 2026-09-05T12:40:22Z
commit: a2c6f4cf977f36217d21fe88c62810ec17d2afb4
heavy_verified_at: 2026-09-05T12:55:39Z
heavy_verified_revision: c78a0d3dfe98ce57e136ea504115f8ae41436f36
completed_at: 2026-09-05T12:56:14Z
release_at: 2026-09-05T12:56:14Z
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


## Review

**Overall: 96%** | 2026-09-05T12:37:17Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Approve with follow-ups.** Both deletion sites are gone, the rule is keyed on the `_at` suffix instead of a field list, and the two new tests fail when the old behavior is put back (verified independently by mutation, below). Three Important findings are all about the *statement* of the rule and its reach, not about the code that shipped.

### Judgment on the three surviving deletions (asked explicitly)

Faithful, not hollowing. Each of the three withdraws a stamp together with the state it describes, and none of them destroys timing history that nothing else holds — which is the harm UR-116 (the board card reading 1m 23s for six hours of work) set out to stop:

- `heavy_verified_at` on re-claim goes with the `commit` it verified. Keeping it would leave a claim advertising heavy evidence for source a remediation already withdrew, which is a dependency-safety regression on a guard one day old (REQ-570, deleting the pending-heavy-testing status).
- `blocked_at` on unblock goes with `blocked_by`. The pair is a live condition, not a phase observation; `status_changed_at` and the `## Blocked` history entry both keep the trace.
- `completed_at` re-stamped on `failed` → `cancelled` was already documented, and the failure instant survives in the `## Cancelled` section's `Previously:` line.

The request's own Context names only the two sites, and REQ-570's guard landed after the capture, so the builder is closing a gap in the request rather than narrowing it. Routing this to the user as a follow-up is the right call. **What is wrong is not the three exceptions but how they are written down** — see finding I2.

### Findings

**Important:**
- Only `claimed_at` is structurally write-once. Every other lifecycle stamp is written by an agent following prose — `actions/work.md:199` (`planning_at`), `:284` (`dispatch_at`, `builder_handback_at`), `:312` (`integration_at`), `:378` (`review_at`, `remediation_at`, `re_review_at`) — and not one of those four instructions carries the "only when the field is absent" condition. `TransitionRecover` strips `route` and the generated sections, so a recovered Route C REQ is re-triaged and re-planned, and the agent then overwrites `planning_at` with the second attempt's instant. That is exactly the case the new schema paragraph names ("a phase re-entered after a recovery … writes its stamp only when the field is **absent**"), and it is the same class of evidence loss as REQ-505. The CLI cannot fix this alone; the four prose sites need the condition. — impact-user-visible → report only
- The exception clause is a closed list, and it is already short. `work-reference.md:75` says "**Four** fields are documented exceptions" one sentence after correctly insisting "**The suffix is the condition**, never a list of today's fields". `testing_updated_at` is a fifth `*_at` field defined in the same schema (`work-reference.md:209`), overwritten on every board testing transition and removed outright by the clear action (`skills/do-work-board/tools/queue-kanban/testing.go:451,458`). It is arguably out of scope because the testing track is orthogonal to the work pipeline (`durations.go:168`) and the board is read-only toward it — but nothing in the paragraph or at the field says so, so the next writer reads a contradiction and either "fixes" the board or widens the rule. Re-key the exception clause on its condition (a stamp that carries current state rather than a phase observation is withdrawn with the state it describes; today those are …, illustrative) and mark `testing_updated_at` as outside the lifecycle. This is the project's own **Closed Enumerations Go Stale** rule applied to the exception list rather than to the rule. — impact-rule-change → report only
- Stale restatement in the board's calibration-log probe. `skills/do-work-board/tools/queue-kanban/verify.go:1374` (doc comment) and `:1445` (the `Remedy` text a human reads out of `queue-kanban verify`) both name "a crash-recovery pass that cleared and re-stamped a claim" as one of the three legitimate reasons a `calibration-log.tsv` row can disagree with frontmatter. After this change recovery never clears or re-stamps `claimed_at`, so that cause is gone and a maintainer reconciling a row is sent to look for something that can no longer happen. The other two causes (the SessionStart repairer, `audit-archive-timestamps.sh --fix`) are still real. — impact-user-visible → report only

**Minor:**
- `internal/publication/answer.go:321-322` deletes `blocked_at` together with `blocked_by` on the stakeholder-terminal completion path. The schema's `blocked_at` exception is worded "removed together with `blocked_by` when an **unblock** clears the condition"; this path is a terminal stakeholder disposition, not an unblock. The deletion is correct behavior, the exception wording does not reach it. — impact-negligible → report only
- `TransitionComplete`, `TransitionFail` and `TransitionCancel` all `SetScalar("completed_at", …)` unconditionally (`state_apply.go:623,635,650`). Only the failed → cancelled overwrite is documented as an exception; nothing guards the other two, so any future path that reaches a second terminal transition silently overwrites the first terminal instant without a reader noticing. `release_at` shows the shape a guard can take — `finalization_apply.go:447-449` refuses a different existing value. — impact-negligible → report only

**Nit:**
- The new comment in `next_selection.go:322-329` breaks mid-clause ("It is evidence of a *live* / claim only while the status still / says claimed"), a leftover of the comment-polish pass. Cosmetic only. — impact-negligible → report only

### Independent verification of the three questions the orchestrator asked

**The selector change does not weaken any real veto.** The `ALREADY-CLAIMED` veto runs at `next_selection.go:186`, *before* the `STATUS-NOT-PENDING` gate at `:240`. So for the case D-03 flags as its risk — a request carrying `claimed_at` under an unrecognized status such as a legacy `in-progress` — the request is still excluded, now by `STATUS-NOT-PENDING` instead of `ALREADY-CLAIMED`; only the exclusion code and its `ClaimEvidence` payload change, never the outcome. `status` has no default in the Schema Read Contract (`schema_normalization.go:30`), so an unrecognized value resolves to itself with a warning and can never resolve to `claimed`; the recognized non-pending holds (`blocked-archive-collision`, `blocked-dependency-cycle`) are excluded by the same gate. The checkpoint-writer branch is untouched and still vetoes on its own, which is what carries a live claim across checkouts. The one behavior that genuinely changes is the intended one: a `pending` request carrying a stamp is selectable again. The gate-deferral parent is unaffected in either direction — it never carried the stamp before (the deletion this REQ removed) and is held by `depends_on` on the repair, exactly as the new comment claims.

**The reader sweep is complete for the readers that matter, and the builder's two specific claims check out.** The board's stale-claim probe (`verify.go:625` `appendClaimFindings`) walks `board.Columns.Claimed` only, so a recovered `pending` request carrying a stamp raises no finding — confirmed at the source. The calibration span is `claimed_at` → `completed_at` at both its writer (`state_plan.go:410-416`) and its independent proof (`finalization_discovery.go:429-438`); after a recover-and-re-claim it measures from the first claim, which the Builder Guidance states as the intended reading, and the calibration's own read-time rule excludes spans over four hours (`durations.go:32`, `actions/estimate-reference.md:92,98`) so an interrupted REQ drops out of the estimator corpus rather than skewing it. Two more readers the sweep did not name, both checked and both clean: `doctor`'s `STUCK-WORK` finding computes claim age from `claimed_at` but is gated on `TreeSection == "working"` (`doctor_scan.go:274`), so a recovered request in `queue/` never reaches it — though note that a *re-claimed* request in `working/` now reports its age from the first claim, which is the intended reading and not a defect; and `addStaleQueueFinding` reads `created_at`/`blocked_at` only. The board's completion-anomaly and `created_at ≤ claimed_at ≤ completed_at` ordering checks (`model.go:1466`, `verify.go:518`) get *safer* under an earlier kept stamp, never noisier.

**Restatement Sweep — applied, one stale restatement found.** The redefined elements are (a) what a lifecycle transition may do to an `*_at` field and (b) what `claimed_at` means on a request that is not currently claimed. Swept: every `claimed_at` / `ClaimedAt` reader in Go across both modules, every `*_at` mention in shipped action prose, `docs/`, primes and `_dev/`, and every phrase pairing recovery with clearing or re-stamping. Verified as still correct: `actions/abandon.md:59`, `actions/estimate-reference.md:92,98`, `actions/work.md:537`, `work-reference.md:112` (`write_set` is not a stamp and recovery still clears it), `work-reference.md:254` and `:512`, `work-reference.md:403` (a concurrent double claim still writes two different values, because each side claims from a stamp-less pending file), `docs/work-guide.md:123,130`, `docs/board-guide.md:23,36,38`, `docs/forensics-guide.md:24,27`. Found stale: the calibration-probe pair above. The `CHANGELOG.md` entries that describe the old behavior (`:4254` and others) are historical record and correctly left alone. The requirement "delete any sentence there that says recovery strips phase observations" is vacuously satisfied — checked `git show cd686ed7:skills/do-work/actions/work-reference.md`, no such sentence existed.

### Requirements Checklist

- [x] Rule keyed on the `_at` suffix, not a name list — delivered (`setLifecycleStampWhenAbsent`, and the ten-name loop replaced by one `DeleteField("route")`, so no list survives to forget a field from)
- [x] `TransitionRecover` keeps every existing `*_at` stamp — delivered, nine stamps pinned by subtest
- [x] Recover still resets `status`, `status_changed_at`, `route`, `write_set`, generated sections — delivered, asserted in the same test
- [x] `status_changed_at` still overwritten, named as the exception — delivered
- [x] A re-claim keeps the first `claimed_at` — delivered, asserted at a re-claim two hours later
- [x] Gate deferral leaves `claimed_at` on the parent — delivered, `defer_gate_test.go:35` now requires the original value
- [x] Rule recorded in the Request File Schema beside the Timestamp rule — delivered (`work-reference.md:75`), with the exception-list caveat in finding I2
- [x] Delete any sentence saying recovery strips phase observations — N/A, no such sentence existed at `cd686ed7`
- [x] No new timing file, stream or writer — delivered
- [x] Doctor stamp repair and `audit-archive-timestamps.sh` untouched — delivered, and the schema paragraph names the repair path as a non-transition
- [x] `route` and `write_set` recover behavior unchanged — delivered (`write_set` still gated on a `## Scope` section)

### Acceptance Testing

**Result: Pass**

- `go test -count=1 ./internal/requeststate/... ./internal/publication/... ./internal/nextselection/... ./internal/lifecycleadvance/... ./internal/doctor/...` — all five packages `ok` (7.0s / 21.0s / 3.3s / 24.3s / 4.5s). `lifecycleadvance` and `doctor` were added by this review because they hold the two readers most exposed to a kept stamp.
- `TestRecoverAndReclaimPreserveEveryLifecycleStamp -v` — nine subtests, one per stamp, matching the builder's RED report field for field.
- **Independent mutation check** (run on a scratch copy of the module, never in the checkout): restoring the ten-name deletion loop in `state_apply.go` fails the test with nine "recover deleted X" lines plus "re-claim overwrote the first claim"; restoring the unconditional `claimed_at` veto in `next_selection.go` fails `TestRecoveredRequestKeepingItsClaimStampStaysSelectable` with `ALREADY-CLAIMED` on REQ-705. Both new assertions are load-bearing, not decorative.
- Cross-REQ test updates traced: all three changed assertions (`state_apply_test.go:228`, `defer_gate_test.go:35`, `next_selection_test.go:135`) name REQ-575 in a comment stating the behavior they now pin.
- P-A-U boxes: all three `[x]`, and the `[UNIFY]` claims match the diff (7 files, +191/−14 in the merge range).

### Suggested Additional Testing

- Run one real recover-and-re-run of a Route C REQ end to end and read the board drawer: this is where finding I1 becomes visible, because the agent-written phase stamps are the ones with no write-once guard.
- Open the board on a re-claimed REQ and confirm the Claimed-row stopwatch and the three-hour stale-claim finding now count from the first claim. That is intended, but it is the first time a live card will show it, and it is worth seeing once before it surprises someone mid-run.
- Run `queue-kanban verify` on this repo after the next recovery and check that no new `calibration-log-mismatch` rows appear (none are expected — the row is written at completion from the same stamp that now survives).

### Follow-ups created
None (6 findings report only)

**Important findings (each with its recorded impact token):**
- Agent-written phase stamps (`work.md:199,284,312,378`) carry no write-once condition, so a re-run after recovery overwrites `planning_at` and its siblings — impact-user-visible → report only
- The schema's exception clause is a closed list of four and already misses `testing_updated_at` (`work-reference.md:209`, `queue-kanban/testing.go:451,458`) — impact-rule-change → report only
- `queue-kanban/verify.go:1374,1445` still name a crash-recovery claim re-stamp as a calibration-log mismatch cause that can no longer happen — impact-user-visible → report only

**Minor findings:** `answer.go:321-322` deletes `blocked_at` outside an unblock, which the schema exception's wording does not reach — impact-negligible → report only; `completed_at` is written unconditionally on all three terminal transitions with no guard like `release_at`'s — impact-negligible → report only
**Acceptance:** Pass — five packages green, both new tests proven mutation-sensitive on a scratch copy
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Keying the rule on the field-name suffix and then *deleting the list* rather than shortening it. Recover no longer has a stamp enumeration at all, so a stamp added to the schema next month is covered without anyone remembering to edit anything. The test does the same thing: it reads the stamp set out of the fixture by suffix, so it covers a future stamp too.

**What didn't:** Changing the writer alone was not enough and would have shipped a one-way door. Keeping the claim stamp through recovery made the selector veto every recovered request as already-claimed, so recovery would have made a request permanently unselectable. An existing test in a third package caught it, not the new ones — which is the argument for running the whole module rather than the two packages the change names.

**Worth knowing:** The exception list in the new schema paragraph is itself a closed enumeration, and the review found a fifth `*_at` field outside it within minutes. The durable form is the condition — a stamp that carries current state is withdrawn with the state it describes, a stamp that records a phase observation is never withdrawn — with today's fields as illustration. Also: only the claim stamp is structurally write-once. The other lifecycle stamps are written by an agent following prose in the work action, and none of those instructions carries the "only when absent" condition, so a recovered request that is re-planned still overwrites its planning stamp. That gap is the same class of evidence loss this request set out to stop, and it goes to the user with the exceptions question.

## Orientation

Lifecycle transitions no longer delete timestamps, so an interrupted request keeps the record of when its work actually started instead of reporting the last few minutes of it. Lives in the do-work-cli request-state subsystem (`skills/do-work/tools/do-work-cli/prime-do-work-cli.md`) with the rule stated once in the request file schema (`_dev/primes/prime-action-files.md`). [MAP CHANGED] — a frontmatter timestamp is now append-only by contract rather than by each writer's habit, and the selector's notion of a live claim moved from "the stamp exists" to "the status says claimed", which is what it always meant. Board-visible consequence: after a recovery and re-claim, the card's wall time and the drawer's Claimed row measure from the first claim. No prime was made stale; the schema section they both point at carries the new rule.

## Heavy Verification Plan

- **Base revision:** cd686ed76961be125db6d3fff214cac5800da819
- **Target revision:** a2c6f4cf977f36217d21fe88c62810ec17d2afb4
- **Planned at:** 2026-09-05T12:40:22Z, from `_dev/tests/heavy-lanes.json`

| Lane | Argv | Why it was selected |
| --- | --- | --- |
| `do-work-cli-integrations` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` | the change edits do-work-cli lifecycle state writers and the selector |
| `staged-skills` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` | shipped action prose changed |
| `updater` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` | shipped package content changed |
| `installer` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` | shipped package content changed |

No path was left uncovered by the manifest. The request stays `claimed` with its `commit:` landed until the queue-exhaustion drain.

## Heavy Verification Result

- **Target revision:** a2c6f4cf977f36217d21fe88c62810ec17d2afb4
- **Execution revision:** c78a0d3dfe98ce57e136ea504115f8ae41436f36
- **Run at:** 2026-09-05T12:55:39Z

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `do-work-cli-integrations` | 0 | 66s | executed (fingerprint_mismatch) |
| `staged-skills` | 0 | 38s | executed (fingerprint_mismatch) |
| `updater` | 0 | 70s | executed (fingerprint_mismatch) |
| `installer` | 0 | 42s | executed (fingerprint_mismatch) |

Every lane this request selected was present in the run, exited 0, and none was skipped. No lane was reused from earlier evidence; all four executed against this tree.

