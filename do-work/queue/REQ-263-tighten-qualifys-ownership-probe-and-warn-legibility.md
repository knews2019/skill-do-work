---
id: REQ-263
title: Tighten qualify's ownership probe and make its WARN legible
status: pending
created_at: 2026-08-18T19:52:15Z
status_changed_at: 2026-08-18T20:55:14Z
user_request: UR-055
addendum_to: REQ-254
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Tighten Qualify's Ownership Probe and Make Its WARN Legible

## What

REQ-254's ownership condition ("printed output belongs to whoever owns the process exit") is implemented as a whole-file grep for exit-idiom text, which is weaker than the condition as stated. Reproduced by its review, three ways: adding `sys.exit(0)` in the same diff as a debug print flips FAIL to WARN; a pre-existing `__main__`-guarded exit makes every debug print in a dual-use module WARN; and a library file whose **docstring merely says "exit 1 on failure"** WARNs. Also: the WARN branch omits the matched lines the FAIL branch prints, so "confirm from the diff" costs a manual dig.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-254 review, Important finding 1 (gate: trivial — the WARN's confirm-from-diff instruction plus Step 6.3's judgment contract mitigate; threat model is forgetfulness, not adversarial builders) plus its folded Minor (WARN legibility). Created `pending-answers` per the generation-≥2 depth stop. The categorical fix (exit-added-in-same-diff ⇒ FAIL) is known-wrong: a legitimately new checker adds its prints and its exit in one diff.

## Requirements

- The ownership probe moves toward code-shaped exit occurrences (not docstring/comment text) and/or base-revision ownership for pre-existing files — direction is builder latitude; the boundary that ships is pinned by a lock-in either way (REQ-250's lesson: pin the documented limitation with a fixture that can fail, e.g. the docstring-"exit 1" case).
- The WARN branch prints the matched lines exactly as the FAIL branch does.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [ ] REQ-254's review found the ownership probe satisfiable by non-semantic bytes (same-diff exit, `__main__` guard, docstring prose) and the WARN branch less legible than FAIL. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the WARN channel plus orchestrator judgment is mitigation enough.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.
