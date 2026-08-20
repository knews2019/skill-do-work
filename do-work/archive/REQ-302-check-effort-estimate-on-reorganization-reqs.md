---
id: REQ-302
title: "[impact-negligible] Check whether capture's effort judgment under-sizes reorganization REQs"
status: cancelled
created_at: 2026-08-20T08:37:00Z
status_changed_at: 2026-08-20T08:37:00Z
completed_at: 2026-08-20T11:38:27Z
user_request: UR-056
addendum_to: REQ-258
domain: general
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: []
---

# Check Whether Capture's Effort Judgment Under-Sizes Reorganization REQs

## What

REQ-258 carried `effort_estimate: trivial`. That normalizes to `effort-mechanical`, which fires Step 3.6's mechanical-effort short-circuit and produced a **5-minute P50** for a wholesale restructure of 1882 lines into 19 files.

The field means *size*, and the per-case change genuinely was trivial — nothing was rewritten. But the work was not. Reorganization REQs may be a systematic blind spot for the judgment capture makes: "each individual edit is mechanical" and "this REQ is mechanical" come apart exactly when the file count is the work.

**This is a question before it is a task.** One data point is not a bias. The cheap version: read `do-work/calibration-log.tsv` and the archived REQs alongside it for other reorganize/split/extract REQs and see whether their estimates ran short the same way. If two or three did, the fix is a sentence in `actions/capture-reference.md`'s effort guidance; if not, close this.

Estimation never gates anything, which is why this is `impact-negligible` — the figure was wrong and cost nothing.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-258: `effort_estimate: trivial` produced a 5-minute estimate for a 19-file restructure, suggesting capture may judge reorganization REQs by per-edit size rather than total size. Should I process this as a new task? → Discarded
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — one data point, and estimation gates nothing.

## Cancelled

- **When:** 2026-08-20T11:38:27Z
- **Why:** explicitly `impact-negligible`; one data point, and estimation never gates work. Discard it.
- **Decided by:** user, via `do-work clarify`

## Context

Discovered during REQ-258. Its `## Discovered Tasks` section carries the original note.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
