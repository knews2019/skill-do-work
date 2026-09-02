---
id: REQ-518
title: '[impact-rule-change] Run the full gate once per REQ'
status: pending
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-519, REQ-520, REQ-521, REQ-522, REQ-523]
batch: cheap-maintainer-gate
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/tools/do-work-cli/, _dev/tests/contract-regressions.sh]
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-02T21:34:07Z
  basis:
    - Route B
    - 4-file write set
    - 1 new files
    - 2 subsystems involved
    - 7 acceptance criteria
    - persistence changes
    - cross-route regression gates
---

# Run the Full Gate Once per REQ

## What

Stop running the identical canonical repository gate twice per REQ. Record the revision hash after every green gate run; before dispatch, when `HEAD` equals that recorded revision, take the baseline as green without running the gate. The Step 6.5 run after implementation stays mandatory.

The fold-first scan found no pending or pending-answers REQ that owns this: REQ-469 (Replace the unrelated canonical-gate hold with a blocked set-aside) and REQ-470/REQ-471 change what happens when the gate is red, not how often it runs.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`actions/work.md` runs the same gate argv at the pre-dispatch baseline (line 402) and again at Step 6.5 (line 531). The gate takes 6.5 minutes, and 167 REQ commits landed in 14 days: about 36 hours of gate wall-clock, half of it a re-run of a revision that was already proven green. Capture-time answer Q3: skip the baseline when `HEAD` is the last green revision.

## Context

The baseline exists for attribution: a red Step 6.5 result is compared to the pre-dispatch fingerprint to decide whether the failure is the REQ's own (`actions/work-reference.md` → Repository Gate Deferral and Resumption). A recorded green revision carries the same information: if the baseline revision was green and Step 6.5 is red, the diff is the cause. The recorded revision must be written only by a green run of the exact gate argv, and read only when it equals `HEAD` exactly.

## Detailed Requirements

- Mechanics in the Go CLI, not in prose: one command records the green revision for a gate argv after a zero exit; one command answers whether `HEAD` matches the recorded green revision for that argv. Where the record lives is the builder's choice; it must survive a session restart and must not be a hand-edited pipeline field.
- `actions/work.md` Step 5.75 and `actions/work-reference.md` → Session state and baseline: when the check reports a match, save it as the green baseline and dispatch; otherwise run the gate as today. A launch failure of the check stops safely, as a gate launch failure does today.
- Every green gate run in the pipeline (baseline or Step 6.5) records its revision, so the next REQ in the same run skips its baseline.
- Attribution at Step 6.5 and the deferral lifecycle work unchanged from a recorded baseline; the fingerprint procedure is untouched.
- A fingerprint or record that names a different argv, a different repository, or a revision not in `HEAD` ancestry never counts as a match.
- The record points at the revision after the gate's own log commit (REQ-523, Log and commit every maintainer gate run): the gate records its green revision after it has committed its log, and a commit whose only paths are under `_dev/gate-runs/` never invalidates a recorded green revision. Without this the log commit moves `HEAD` off the record after every run and the skip never fires. (verify-requests, 2026-09-03)
- The self-referential refusal invariant from REQ-514 applies to any new finding.

## Constraints

- Never waive the gate: a REQ still needs a green Step 6.5 run to complete.
- No new sentence predicates in `_dev/tests/contract-regressions.sh`; adjust the existing baseline predicates it already pins, and cover the behavior in a Go test.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

Roots this batch: REQ-519 (Path-scoped fast lane for the maintainer gate) depends on it because both change the same gate lanes in `work.md` and the maintainer asked for this first. Related to REQ-469 through REQ-471 by write-set overlap only.

## Builder Guidance

Certainty level: Firm on the rule (one full gate run per REQ when the revision is already green), latitude on the record's location and command names. Read the CLI prime first.

## Red-Green Proof

**RED prompt/case:** Run `do-work run REQ-NNN` on a revision whose gate was just green, and count gate invocations before the builder is dispatched.
**Why RED now:** The pipeline runs the full gate again at the baseline although `HEAD` has not moved since the green run; there is no record to consult.
**GREEN when:** A Go test records a green revision, sets `HEAD` to it, and asserts the check reports a match; a second test moves `HEAD` and asserts no match; a third adds one commit touching only `_dev/gate-runs/` on top of the record and asserts the match still holds. The work action reads that check and a run on an already-green revision performs exactly one full gate run, at Step 6.5.
**Validation:** User confirmed (capture-time answer Q3, 2026-09-03).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing pipeline fields and the run action's gate lanes.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on structured evidence projection and git-transaction boundaries in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A1 of the analysis report's improvements, captured by UR-100.*
