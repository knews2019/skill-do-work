---
id: REQ-099
title: Automatic wave dispatch — the work loop computes and dispatches the ready set
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-096]
maintenance: false
related: [REQ-096, REQ-100]
batch: parallel-building
write_set: [actions/work.md, actions/work-reference.md]
---

# Automatic Wave Dispatch

## What

Give the work pipeline a fan-out mode where **the loop computes the wave itself** and dispatches builders without a confirmation gate. This is a deliberate contract change: today `actions/work.md:33` says the action "does not drive a fan-out wave" and `actions/work-reference.md:320` says "a human picks which REQs run together — nothing computes the set." Both sentences get rewritten, per the user's explicit choice of fully automatic set-picking over a human-confirmed set.

## Detailed Requirements

- **Wave computation:** ready = pending REQs whose `depends_on` are satisfied, unclaimed, and not `assigned_to` another session. Wave size bounded per `crew-members/background-agents.md:53` (builders per wave sized to the harness concurrency limit). `--wave N` keeps its existing meaning (dependency-depth scoping) — document how the two compose or that auto-wave supersedes it when active.
- **Dispatch:** for each REQ in the wave, follow the existing Worktree Dispatch Mode per-REQ flow unchanged (worktree per builder mandatory, run directory + briefs before any spawn, hand-back merge sequence). Silent degradation stays: no `git worktree` support → serial, no error.
- **Integration stays serial and load-bearing:** merge → qualify → test → review → changelog → archive one REQ at a time; `actions/work-reference.md:321` ("the non-interference proof is the merge, not the pick") survives unchanged and becomes the safety argument — overlapping picks are caught at merge, which is the batch philosophy. `write_set` stays display-only; the wave computation must NOT read it as a scheduling input.
- **Mode entry:** define how auto-wave is invoked (e.g., a `do-work run` fan-out flag or harness-capability trigger) — floor-first: the default single-REQ loop remains the baseline for the simplest agent; auto-wave is the advanced-harness path, consistent with the existing "Optional, advanced harnesses only" gate at `:277`.
- Update every echo of the old "human picks / nothing computes the set" claim across shipped files (Closed Enumerations Go Stale rule).

## Constraints

- No `write_set` scheduling; no computed-set gate on anything other than `depends_on`, claim state, and `assigned_to`.
- Wall-clock saving is build-phase only — keep the honest-expectations sentence (`:322`).

## Red-Green Proof

**RED prompt/case:** Ask the work action to run the queue in parallel today: `actions/work.md:33` instructs it to process one REQ at a time and provides no wave-selection or launch-before-wait path.
**Why RED now:** Fan-out is an owner-driven manual procedure by explicit design (confirmed at `work-reference.md:316`).
**GREEN when:** The rewritten instructions let an advanced harness compute a bounded ready set and dispatch builders unattended, while the floor path (serial single-REQ) is unchanged and the merge remains the non-interference proof.
**Validation:** User confirmed (ask-tool answer: "Fully automatic set-picking").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 3, item 8).

---
*Source: approved plan, Phase 3*
