---
id: REQ-256
title: Disclose the session hook's queue write surface in the docs
status: pending
created_at: 2026-08-18T17:48:08Z
user_request: UR-056
addendum_to: REQ-246
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- README.md
- skills/do-work/actions/capture.md
---

# Disclose the Session Hook's Queue Write Surface in the Docs

## What

REQ-246 made the SessionStart hook a *write* surface on consumer queue files — it mechanically repairs detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` at session start. Two shipped texts still describe the hook as read-only-plus-banner; a consumer auditing "what writes to my repo at session start" is misled.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Instance 1 is REQ-246's review finding I3 (restatement sweep — the omission predates REQ-246 for the reservation cleanup, but a hook that edits files unattended is a different disclosure class than one that prints a banner). Instance 2 is the optional integration seam REQ-246's builder offered and deliberately did not write. Citations follow the literal cross-package path rule REQ-249 establishes.

## Instances

- [ ] **`README.md` (~line 188):** "SessionStart hook that injects the installed version and pending REQ count" — add that the hook also reaps stale REQ-number reservations and mechanically repairs detectably wrong queue/working timestamps (one clause each; keep it one sentence if it fits).
- [ ] **`skills/do-work/actions/capture.md`:** document `scripts/repair-req-timestamps.sh` (SessionStart hook) the way `cleanup-req-reservations.sh` is documented — one line stating that detectably wrong queue/working `*_at` stamps are mechanically corrected at session start.

## Requirements

- Both texts state the hook's write behavior; no behavior change anywhere.
- `bash _dev/tests/maintainer-verify.sh` exits 0.
