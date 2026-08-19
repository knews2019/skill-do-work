---
id: REQ-274
title: Retire the "the SessionStart hook exits nonzero" framing where it is still stated
status: pending-answers
created_at: 2026-08-18T23:38:35Z
status_changed_at: 2026-08-18T23:38:35Z
user_request: UR-056
addendum_to: REQ-267
domain: general
review_generated: true
sweep: true
sweep_key: repairer-hook-failure-framing-restated
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- do-work/RESTART-PROMPT.md
- do-work/CHECKPOINT.md
- _dev/primes/prime-shell-commands.md
---

# Retire the "the SessionStart Hook Exits Nonzero" Framing Where It Is Still Stated

## What

A timestamp-repairer failure **does not fail the SessionStart hook**. `skills/do-work/hooks/session-start.sh:59` runs the script under `|| true`, deliberately, with a comment saying that on a tripped guard the script's failure lines *are* the audit trail and must reach the banner. `report_failure` writes to stdout, the hook captures it into `REPAIR_SUMMARY` and echoes it, and the hook exits **0** — verified by running the real hook against a wedge fixture.

The real consequence is still bad: the script exits 1 on every run and prints a `FAILED to repair …` line into every session's start banner, permanently, with no self-heal. But three live maintainer docs state the *mechanism* as a nonzero hook exit, and one of them is a standing decision rationale:

- `do-work/RESTART-PROMPT.md:39` — "can wedge the SessionStart hook into exiting nonzero every session"
- `do-work/CHECKPOINT.md:76` — REQ-255 D-04's rationale: "refusal would have made the SessionStart hook exit nonzero every session". **The decision is still right; the argument for it names a mechanism that does not exist.**
- `_dev/primes/prime-shell-commands.md:34` — "the fuzz found the shape that wedges the hook"

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [x] `do-work/RESTART-PROMPT.md` — **already corrected.** Both files were rewritten wholesale at the end of the session that
  created this REQ, and the orchestrator wrote the true framing rather than reproducing a claim it knew to be false. Verify
  rather than assume: re-grep before deciding this instance is closed.
- [x] `do-work/CHECKPOINT.md` — **already corrected**, same rewrite. REQ-255 D-04's rationale no longer appears there at all,
  because the checkpoint was replaced by the next session's; if that decision's reasoning is still wanted in a durable place,
  its archived REQ is the home, and archived REQs are immutable.
- [ ] `_dev/primes/prime-shell-commands.md:34` — **still live.** This is the one that matters: a prime is loaded and acted on.

## Note added after the instance list was written

Two of the three instances sat in `do-work/` state files that the session rewrites at its end. Rather than deliberately re-writing a known-false sentence to preserve a tidy sweep, the orchestrator wrote the correct framing into both. **This REQ is therefore mostly about the prime**, plus the sweep that checks nothing else restates the old mechanism. Do not read the two ticked boxes as work done by this REQ's builder — read them as instances that expired, and re-grep to confirm.

## Requirements

- The rule is **stated once**: a repairer failure prints into the session banner; it does not fail the hook. Nothing restates it on the old mechanism.
- Where the false mechanism carried a decision's argument (REQ-255 D-04), the decision keeps standing on the true consequence — a permanent banner failure line with no self-heal supports it just as well. **Do not silently reverse a decision while correcting its reason.**
- Sweep the primitive rather than fixing three lines: grep for every statement about what a repairer failure does to the hook, in any spelling, and check each against `session-start.sh`.
- Archived REQs under `do-work/archive/` are immutable record and stay untouched.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-267's independent review, Important finding 1 (gate: rule-change). Created `pending-answers` per the generation-≥2 cascade stop, since REQ-267 is itself `review_generated: true`.

This is REQ-257's lesson arriving one REQ later, and it is worth stating plainly: **a claim inherited rather than re-derived can be right in its conclusion and false in its mechanism, and the false mechanism travels.** This one travelled far enough to set REQ-267's own approval framing — the orchestrator repeated it to the maintainer twice before a builder checked the hook.

## Open Questions

- [ ] REQ-267's review found that three live maintainer docs state a repairer failure as making the SessionStart hook exit nonzero, which it does not — including a standing decision rationale that rests on it. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the affected decisions are all still correct, so accept the wrong mechanism in the record.
