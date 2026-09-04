---
id: REQ-566
title: '[impact-rule-change] Run held heavy lanes at queue exhaustion without asking, and record per-lane wall time'
status: claimed
created_at: 2026-09-04T13:19:11Z
user_request: UR-111
domain: backend
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
priority: now
effort_estimate: effort-substantive
claimed_at: 2026-09-04T13:24:15Z
---

# Run Held Heavy Lanes At Queue Exhaustion Without Asking, And Record Per-Lane Wall Time

## What

Two changes, in this order.

1. **Record per-lane wall time.** The generated `## Heavy Verification Result` section records wall seconds for every lane next to its exit status, and the typed `heavy_testing` evidence (`HeavyLaneResult` in `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go`) carries that field so the result is machine-readable.
2. **Lift the human permission prompt.** At queue exhaustion, the work loop runs the batched heavy plan itself at HEAD for every `pending-heavy-testing` REQ, then records green or red through the same `answer --manifest` transaction that `do-work clarify` uses today. The `pending-heavy-testing` status, the non-blocking hold, and the source-ready rule for dependents all stay. Only the permission question goes. `do-work clarify` keeps working for a maintainer who runs the lanes by hand.

The fold-first scan found no pending or pending-answers REQ in any UR that shares this root cause. REQ-564 (reuse matching per-lane verification evidence for four hours) avoids rerunning unaffected lanes and REQ-559 (retry a red repository gate once before a repair REQ) handles a red gate; neither removes the permission prompt or records lane durations.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The permission gate was a cost control from when the whole maintainer gate cost about seven minutes per REQ. It now protects roughly three and a half minutes of CPU per batch, and it is the only thing that stops an unattended queue drain. The user's goal is to drain the queue faster; the human wait is the bottleneck, not the tests.

## Context

Measured from `do-work/test-durations.tsv` on 2026-09-04:

| run | wall time |
|---|---|
| typical full gate, 126 files | 180 to 220 s |
| worst of the last ten runs | 326 s |
| fast tier alone, 14 files | 43 to 63 s |

Batching already exists: on 2026-09-04 one heavy run at revision `c0d8ce1c` stamped five held REQs (REQ-475, REQ-483, REQ-485, REQ-502, REQ-561) with the same `heavy_verified_at`. The lane selector (REQ-563, archived) already limits each plan to the lanes whose coverage paths the diff touches.

Risks the builder must handle:

- **R1, a skipped lane must stop the run, not loop.** The browser lane skips when no Chrome binary is found, and a skip counts as not confirmed. An automatic run that meets a skip records the skip, leaves the REQ at `pending-heavy-testing`, reports a typed finding naming the lane, and ends the run. It never routes the REQ through remediation for a lane that did not execute.
- **R2, lane commands must survive the deletion of the maintainer gate script.** Every lane argv in `_dev/tests/heavy-lanes.json` still calls `_dev/tests/maintainer-verify.sh`, which is slated for deletion. Do not add a new dependency on that script; the lane runner reads argv from the manifest so the commands can move without touching this feature.

Out of scope: a duration budget guard ("run automatically under N minutes, ask above"). Add it only when recorded lane durations show a lane that needs it.

## Red-Green Proof
**RED prompt/case:** One REQ sits at `pending-heavy-testing` with a green fast tier and a recorded heavy plan, the ready set is empty, and the user runs `do-work run` unattended.
**Why RED now:** The run ends by asking for permission to run the lanes. The REQ stays held until a human answers, and the existing `## Heavy Verification Result` shape records exit status only, so nothing measures how long a lane took.
**GREEN when:** The same run executes the planned lanes at HEAD, records exit status and wall seconds per lane in `## Heavy Verification Result` and in the typed `heavy_testing` evidence, returns the REQ to `pending` with `heavy_verified_at` and `heavy_verified_revision` set, and the next selection returns it with `resume_phase: review`. A red or skipped lane leaves the REQ held with the lane named in a typed finding and the run ends cleanly.
**Validation:** User confirmed

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, `slugged: partial` so only selectable bare; matches because this REQ changes a status contract and downstream readers in `actions/work.md` and `actions/clarify.md`. Over the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens, `slugged: partial`; matches because this REQ changes typed evidence the answer command writes and reads (family `opaque-evidence-projection`). Over the 2000-token budget.

## Assets

None.

---
*Source: "stopping on pending-heavy-testing because that needs your permission." <- I asked for this, but it turns out it's not very good, because the work just stops, what I need is an inteligent way to run it (either faster by 80/20 rule), or groupped with multiple tasks (but that can get complicated). Any suggestions? I would start by monitoring the duration and lifting the human blocking, after all the goal is to do the queue faster.*
