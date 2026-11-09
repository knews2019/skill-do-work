---
title: "Lessons from REQ-086: The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-086-the-in-progress-record-s-rule-is-unstate.md]
related:
  - page: REQ-077-crash-recovery-s-own-crash-branch-is-unr
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-086: The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

REQ-077 made `do-work/CHECKPOINT.md`'s `## In Progress (interrupted)` list a claim-time record and
stated the removal rule as a **trigger condition** — *whenever a REQ leaves `do-work/working/`, its
entry goes with it* — in the canonical home (`actions/work-reference.md` → **In-Progress Record
(Step 2)**). REQ-077's own scope covered the pipeline's three movers (Step 8's archive on success and
on failure, and the mid-run blocked flip). Three sites outside that scope were found by REQ-077's
Restatement Sweep and routed here:

1. **`actions/cleanup.md` Pass 0 step 5** moves a terminal-status REQ out of `working/` to
   `archive/`. It says nothing about the in-progress entry, so the entry survives the move.
2. **`actions/forensics.md` Check 1**'s suggested manual reset returns a stuck `claimed` REQ to
   `do-work/queue/`. Same gap. This one matters more than it looks: it is the documented remedy for
   a stranded claim, so it is precisely the procedure a user runs *because* recovery is involved.
3. **`docs/work-guide.md:66`** tells the user "At session end, a `do-work/CHECKPOINT.md` is
   written…" — which was the whole of the truth before REQ-077 and is now the half of it that made
   the 0.164.0 regression possible. Line 119 ("The system writes `do-work/CHECKPOINT.md` before
   stopping") reads the same way.

## Solution summary

Built in a git worktree as Builder A of REQ-085's live fan-out acceptance test — branch `worktree-agent-REQ-086-in-progress-record-unstated`, commit `0e04b4d`, integrated by the `--no-ff` merge `a17e6af` (range `306c1f4..a17e6af`). Pass 0 step 5 and Check 1's manual-reset remediation each gained one clause saying the move out of `working/` also drops the REQ's `## In Progress (interrupted)` entry, each citing `actions/work-reference.md` → **In-Progress Record (Step 2)** rather than restating the trigger. `docs/work-guide.md`'s two sentences were rewritten to say the checkpoint is written at claim time and refreshed at session end, and why that matters.

## What worked

- Citing the canonical home instead of restating the rule kept all three edits to one sentence each — requirement 4 turned out to be a size constraint as much as a correctness one.

## Worth knowing

- The user-guide sentence was the most valuable of the three despite looking like the smallest. `actions/*` sites are read by agents that can open the cited rule; `docs/work-guide.md` is read by a human who cannot, so a half-true sentence there has no correction path.

## Back-reference

See `do-work/archive/UR-015/REQ-086-in-progress-record-unstated-at-three-consumer-sites.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a17e6af`.
