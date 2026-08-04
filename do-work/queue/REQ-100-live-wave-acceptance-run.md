---
id: REQ-100
title: Live auto-wave acceptance run — prove real wall-clock concurrency
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-099]
maintenance: false
related: [REQ-099, REQ-095]
batch: parallel-building
---

# Live Auto-Wave Acceptance Run

## What

Run the REQ-099 automatic wave dispatch live, with genuinely concurrent builders, and record the evidence. Real wall-clock fan-out concurrency has **never been proven** in this skill — the one recorded attempt (REQ-085) logged Partial — so this run is the proof, not ceremony.

## Detailed Requirements

- Use real (or realistic dummy) REQs so at least two builders run **simultaneously** — capture timestamps proving overlap (builder start/end times from the run directory), not just "both completed".
- Exercise the full auto path: automatic set computation (deps + claims + `assigned_to` respected), bounded wave size, worktree-per-builder, serial integration of all hand-backs.
- Record the run artifact in the same form as REQ-085's fan-out run; include at least one observed imperfection or its explicit absence (the previous run's value was finding the index-settling bug).
- If the run surfaces defects in REQ-099's prescribed commands, fix them and grep the same primitive across all actions before calling it fixed (copy-paste rule).

## Red-Green Proof

**RED prompt/case:** No recorded evidence of two builders with overlapping wall-clock execution exists anywhere in the repo (REQ-085: Partial).
**Why RED now:** Fan-out has been driven by hand, serially confirmed, and auto-wave has never existed.
**GREEN when:** A recorded run shows ≥2 builders with overlapping timestamps dispatched by the automatic wave computation and integrated serially, with any found defects fixed and cross-grepped.
**Validation:** User confirmed (approved plan, Phase 3 item 9).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 3).

---
*Source: approved plan, Phase 3*
