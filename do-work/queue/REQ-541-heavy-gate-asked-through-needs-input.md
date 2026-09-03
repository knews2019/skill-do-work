---
id: REQ-541
title: 'The heavy gate is asked for through Needs Input, never run by the loop'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-540]
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/review-work.md
---

# The Heavy Gate Is Asked for Through Needs Input, Never Run by the Loop

## What

In `actions/work.md` Step 6.5, after the fast gate is green: when this REQ's diff (`git diff --name-only <base>..HEAD`) matches any glob from `maintainer-verify.sh --heavy-surfaces`, append `- [ ] REQ-NNN touched <surface> at <revision> (fast gate green <timestamp>)` under `## Instances` of the standing sweep REQ "Heavy gate runs requested" (`sweep_key: heavy-gate-requested`, `status: pending-answers`), creating it once through the Fold-First Rule's conversion path if absent, and continue. The archived REQ records the same line in its Testing section. The loop never runs `--heavy` and never stops for it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"the factory should not stop on the heavy tier, those are being asked by placing the REQ into the need input column" (D4).

## Context

- The Needs Input column shows `pending-answers` and `blocked` queue residents. A single standing sweep, appended to, is the one-file shape the maintainer asked for; appending to a queued sweep is an edit, not REQ creation, so REQ-531's report-only rule is untouched.
- The glob list lives in the gate script (REQ-537) so no prose enumerates surfaces.
- Review and qualification consult recorded evidence since 0.266.8; any surviving sentence that lets them run the gate again is deleted here.

## Detailed Requirements

- One rule in `work-reference.md`, cited from Step 6.5; no new status, no new flag on `do-work run`.
- The sweep REQ carries `user_request:` of the first REQ that created it and a one-line Open Question: "Run `bash _dev/tests/maintainer-verify.sh --heavy` and tick the instances it covers?" with Recommended: yes.
- Ticking is manual; a run that finds every instance ticked leaves the sweep alone. The loop never flips it to `pending`.
- Delete any sentence in `work.md`, `work-reference.md`, or `review-work.md` that permits review or qualification to relaunch the gate.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** A REQ that edits `skills/do-work-board/tools/queue-kanban/web/board-filters.js` runs through `do-work run` after REQ-537 lands.
**Why RED now:** the loop runs the same gate for every diff and has no heavy ask; nothing appears in Needs Input.
**GREEN when:** that REQ completes on the fast tier, "Heavy gate runs requested" appears in the board's Needs Input column with one new instance line naming the REQ and revision, the loop proceeds to the next REQ, and a REQ editing only `actions/*.md` adds no line.
**Validation:** User confirmed (D4)

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes action routing and a pipeline step contract.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.

