---
source_type: req_lesson
req_id: REQ-274
req_path: do-work/archive/UR-056/REQ-274-retire-the-hook-exits-nonzero-framing.md
date: 2026-08-20
domain: general
module: _dev/primes
tags: [general, retire, sessionstart, hook, exits]
---

# Lessons from REQ-274: Retire the "the SessionStart hook exits nonzero" framing where it is still stated

## What the REQ was about

A timestamp-repairer failure **does not fail the SessionStart hook**. `skills/do-work/hooks/session-start.sh:59` runs the script under `|| true`, deliberately, with a comment saying that on a tripped guard the script's failure lines *are* the audit trail and must reach the banner. `report_failure` writes to stdout, the hook captures it into `REPAIR_SUMMARY` and echoes it, and the hook exits **0** — verified by running the real hook against a wedge fixture.

The real consequence is still bad: the script exits 1 on every run and prints a `FAILED to repair …` line into every session's start banner, permanently, with no self-heal. But three live maintainer docs state the *mechanism* as a nonzero hook exit, and one of them is a standing decision rationale:

## Solution summary

Swept every live statement of what a timestamp-repairer failure does to the SessionStart hook and found exactly one still on the false mechanism — REQ-255's lesson link in `_dev/primes/prime-shell-commands.md`, the one site the REQ predicted would matter because a prime is loaded and acted on. Corrected it to name the true consequence. The other seven swept sites were verified individually: three already stated the rule truly, two instances had expired, one was a different claim, and two run artifacts were deliberately left as historical record.

## What worked

Refusing to inherit the REQ's premise. This REQ exists because a claim travelled three documents deep without anyone re-deriving it, so accepting "the hook exits 0, verified against a wedge fixture" on the REQ's word would have reproduced the exact failure being fixed. Re-deriving it also produced something the REQ did not have: the refusal shape it names does *not* reproduce the wedge (voicing is opt-in), so the proof needed a stub repairer that actually fails. The REQ's conclusion was right and its named fixture would not have demonstrated it.

## What didn't work

Two attempts at reproducing a real repairer failure failed before the stub worked — the refusal path is silent by design, and an unwritable-directory trip does nothing when the session runs as root. Worth knowing for any future probe of this script: root defeats every permission-based failure path in it, so simulate the failure at the seam instead of trying to provoke it.

## Worth knowing

The rule's canonical home is the comment block at `skills/do-work/scripts/repair-req-timestamps.sh:108-117`, and it has been correct all along. Every false restatement was downstream of it in a *lesson link* or a *run artifact* — narrative surfaces, not contract surfaces. That is the pattern worth watching: the contract text stayed true while the stories told about it drifted, and the stories are what the next reader reaches first.

## Back-reference

See `do-work/archive/UR-056/REQ-274-retire-the-hook-exits-nonzero-framing.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0efefa6`.
