---
id: REQ-366
title: "Repair the queued tdd REQs whose write set names no test file"
status: pending
created_at: 2026-08-24T23:17:16Z
user_request: UR-070
domain: general
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-365]
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - do-work/queue/REQ-348-group-the-timelines-rows-by-user-request.md
  - do-work/queue/REQ-350-narrow-the-durations-axis-window.md
  - do-work/queue/REQ-352-state-the-durations-headline-numbers.md
  - do-work/queue/REQ-353-hide-dead-filter-knobs-on-durations.md
  - do-work/queue/REQ-362-stop-a-multi-path-bullet-disabling-the-scope-drift-check.md
---

# Repair the Queued TDD REQs Whose Write Set Names No Test File

## What

Five queued REQs carry `tdd: true` with a `write_set` that names no test file — the same shape
REQ-346 hit and REQ-365 closes at capture. REQ-365 stops new ones being minted; it does not touch
the ones already queued. Resolve each of the five so no builder has to re-decide the question.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`tdd: true` triggers the TDD-evidence gate at `actions/work.md` Step 6.5, which will not accept a
hand-back without a failing-then-passing test. A `write_set` that names no test file states a scope
that cannot satisfy that gate. Nothing breaks — `write_set` is display-only — so the cost is paid in
builder judgment: each of the five will stop, report the conflict, and wait for the orchestrator to
extend the set, one at a time, for a contradiction that is already known and fixable in one pass.

REQ-346's builder handled it correctly and the resolution lives in a commit message. That is the
part worth not repeating five more times.

## Detailed Requirements

- For each of the five REQs in `write_set`, decide which of its two facts is wrong and fix that one:
  - the write set is incomplete → add the test file its captured GREEN implies (for the board REQs,
    the Go test beside the code they change; for REQ-362, the `_dev/tests/` case that covers the
    check); or
  - the `tdd: true` verdict does not hold → set `tdd: false`, per the runnable-RED heuristic in
    `skills/do-work/actions/capture.md` Step 1's TDD assessment. Pure layout or visibility work with
    no assertable behaviour is the case that belongs here, and `tdd: false` keeps the REQ's
    `## Red-Green Proof` exactly as captured.
- Where neither is knowable without doing the REQ's own work, **clear `write_set` entirely** rather
  than inventing a path. Absence reads as *unknown* and gets no overlaps badge, which is honest; a
  wrong set makes the board badge misleadingly (`actions/capture-reference.md` → Populating
  `write_set`).
- Leave every other field alone. This is queue bookkeeping, not a re-capture: no retitling, no
  re-judging `impact:` or `effort_estimate:`, no body edits beyond what a changed `tdd` verdict
  requires.

## Constraints

- **`do-work/working/` and `do-work/archive/` are immutable** (`actions/capture.md` → Immutability
  Rule). REQ-347 carries the same shape and is `claimed`; it is deliberately **not** in this REQ's
  scope. Its builder resolves it in-flight through the stop-and-report path at `actions/work.md`
  § "Write only inside the declared scope".
- **Do not make `write_set` a gate.** No lint, no check, no scan wired into a hook or a test suite.
  The field is display-only at any builder count (`actions/work-reference.md` → Request File Schema)
  and this REQ must not change that. The RED command below is a one-off verification, not a
  deliverable.
- This REQ is `tdd: false` on purpose: its proof is a shell scan over queue data, not a test in any
  harness, and minting it `tdd: true` would reproduce the exact shape it exists to repair.
- No new REQ files, no changes under `skills/`.

## Red-Green Proof

**RED prompt/case:** run this from the repository root — it prints five queued REQ files:

```bash
for f in do-work/queue/REQ-*.md; do
  fm=$(awk 'BEGIN{n=0} /^---$/{n++; next} n==1' "$f")
  printf '%s\n' "$fm" | grep -q '^tdd: *true' || continue
  ws=$(printf '%s\n' "$fm" | awk '/^write_set:/{inws=1; next} inws && /^  *- /{print; next} inws{exit}')
  [ -n "$ws" ] || continue
  printf '%s\n' "$ws" | grep -qE '(^|/)[^/]*(_test\.|\.test\.|test_)|/tests?/' || echo "NO-TEST: $f"
done
```

**Why RED now:** REQ-348, REQ-350, REQ-352, REQ-353 and REQ-362 each declare `tdd: true` beside a
write set of implementation files only, so each states a scope its own TDD gate cannot be satisfied
inside.

**GREEN when:** the command above prints nothing, and each of the five reads coherently on its own —
a test path in the set, or a `tdd: false` verdict that matches capture's runnable-RED heuristic, or
no `write_set` at all.

**Validation:** Inferred during capture.

---
*Source: analysis of REQ-365's scope during UR-070 — REQ-365 fixes capture, this fixes what capture already emitted.*
