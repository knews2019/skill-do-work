# Builder brief — REQ-575

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-575-append-only-stamps`
- **Your branch (already checked out there):** `worktree-agent-REQ-575-append-only-stamps`
- **Route:** A
- **Base commit:** 09a13839 (main)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it to the orchestrator in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)

Also read every path in the REQ's `prime_files`, and the `lessons-<name>.md` satellite beside each prime whose Read-first or Traps entries your change touches.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes), verify no debug artifacts (`console.log`, `debugger`, stray `TODO`) in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use:
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./internal/requeststate/... ./internal/publication/...`
- Then the whole module once: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./...`

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-120117/REQ-575-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, and for a `tdd: true` REQ the RED observation (test name + failure) and the GREEN observation.
- `## Lesson evidence` — each `required_lessons` entry you read (whole-satellite or family-targeted) and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN`, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings. Do not fix them inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.

`tdd: true`. The `## Red-Green Proof` in the REQ is the captured RED/GREEN pair and is not yours to rewrite: write the failing test first, observe it fail, then make it pass, and report the test name plus both observations.

---

# The request

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
  - skills/do-work/actions/work-reference.md
claimed_at: 2026-09-05T12:00:55Z
---

# Keep Every Lifecycle Stamp: No Transition Deletes an Existing `*_at` Field

## What

Make REQ frontmatter timestamps append-only. Once a `*_at` field exists on a request, no lifecycle transition deletes it or writes a different value over it. A transition that re-enters a phase writes its stamp only when the field is absent. Delete the two code paths that still remove stamps: the recover transition's field-strip loop in `state_apply.go` and the gate-deferral parent edit in `defer_gate.go`. State the rule once in the Request File Schema so the next writer inherits it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

