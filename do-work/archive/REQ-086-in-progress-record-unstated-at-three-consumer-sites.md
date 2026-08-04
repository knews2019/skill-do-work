---
id: REQ-086
title: The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide
status: completed
claimed_at: 2026-08-04T00:04:10Z
completed_at: 2026-08-04T00:09:48Z
kb_status: pending
created_at: 2026-08-03T21:29:26Z
user_request: UR-015
domain: general
prime_files: []
tdd: false
depends_on: [REQ-077]
maintenance: true
addendum_to: REQ-077
review_generated: true
write_set: [actions/cleanup.md, actions/forensics.md, docs/work-guide.md]
---

# The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide

## What

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

## Why

A leftover entry is not corrupting — the REQ still resolves in `archive/` or `queue/`, so
`queue-kanban verify`'s checkpoint-ghost check stays quiet. The cost is that the **next run reports a
contradiction** ("named in the in-progress record but not in `working/`") on a REQ that completed
perfectly normally. A warning that fires on the happy path is the fastest way to teach a reader to
ignore the warning — and that warning is the only signal distinguishing a real stranded claim from
bookkeeping residue.

The user-guide sentence is a different failure: a reader who believes the checkpoint appears only at
session end will read a post-crash checkpoint as evidence of a clean shutdown.

## Detailed Requirements

1. **`actions/cleanup.md` Pass 0 step 5** — state that moving a terminal REQ out of `working/` also
   drops its `## In Progress (interrupted)` entry, pointing at
   `actions/work-reference.md` → **In-Progress Record (Step 2)** rather than restating the rule.
2. **`actions/forensics.md` Check 1** — same clause in the suggested manual-reset remediation.
3. **`docs/work-guide.md`** — correct both sentences so the guide says the checkpoint is written at
   claim time and refreshed at session end, and say plainly what that buys the user (a crashed run
   picks its own work back up). Keep it at guide altitude — no procedure, one or two sentences.
4. **Do not re-state the removal rule at each site.** The canonical home states the trigger
   condition; these three sites cite it. Copying the rule is how the enumeration goes stale again.

## Constraints

- REQ-077 already named these movers in the canonical procedure's illustrative list, so the contract
  is stated — this REQ only makes each site's own reader see it. If that turns out to be enough on
  its own for one of the three, say so and skip it rather than adding prose to hit a count.
- `docs/` is user-facing: `crew-members/anti-slop.md` applies.

## Dependencies

`depends_on: [REQ-077]` — the rule these sites cite ships in REQ-077.

## Builder Guidance

**Certainty: Firm.** All three sites and their exact lines are named above and were verified against
the tree during REQ-077's review. The open latitude is wording only.

## Triage

**Route: A** - Simple

**Reasoning:** All three sites and their exact lines are named in the REQ and were verified against
the tree during REQ-077's review. Firm certainty, wording-only latitude, three leaf prose edits.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `actions/cleanup.md` (modify) — Pass 0 step 5
- `actions/forensics.md` (modify) — Check 1's suggested remediation
- `docs/work-guide.md` (modify) — the two checkpoint sentences (lines 66 and 119)

**Files I will NOT touch:**
- `actions/work-reference.md` — the canonical home; restating its rule is what requirement 4 forbids

**Acceptance criteria (restated from REQ):**
- [ ] `actions/cleanup.md` Pass 0 step 5 states the entry is dropped, citing the canonical home (req 1)
- [ ] `actions/forensics.md` Check 1's manual reset carries the same clause (req 2)
- [ ] `docs/work-guide.md` says the checkpoint is written at claim time and refreshed at session end,
      and says what that buys the user, at guide altitude (req 3)
- [ ] No site restates the removal rule (req 4)

## Implementation Summary

**Files changed:**
- `actions/cleanup.md` (modified)
- `actions/forensics.md` (modified)
- `docs/work-guide.md` (modified)

**What was done:** Built in a git worktree as Builder A of REQ-085's live fan-out acceptance test —
branch `worktree-agent-REQ-086-in-progress-record-unstated`, commit `0e04b4d`, integrated by the
`--no-ff` merge `a17e6af` (range `306c1f4..a17e6af`). Pass 0 step 5 and Check 1's manual-reset
remediation each gained one clause saying the move out of `working/` also drops the REQ's
`## In Progress (interrupted)` entry, each citing `actions/work-reference.md` → **In-Progress Record
(Step 2)** rather than restating the trigger. `docs/work-guide.md`'s two sentences were rewritten to
say the checkpoint is written at claim time and refreshed at session end, and why that matters.

## Decisions

- **D-01 (DECIDE & STATE)** — *All three sites got the clause; none was skipped.* The Constraints
  allowed skipping any site where REQ-077's canonical illustrative list was already sufficient. None
  qualified: Pass 0 step 5 and Check 1 each describe the move in full procedural detail without
  mentioning the entry, so a reader following either procedure to the letter leaves it behind — which
  is the exact contradiction the next run reports.
- **D-02 (DECIDE & STATE)** — *`docs/work-guide.md` line 119 was rewritten rather than left alone.*
  The REQ named line 66 as the primary defect and line 119 as reading "the same way". Line 119's
  "writes CHECKPOINT.md before stopping" is the more actively misleading of the two for the case that
  matters — a reader hitting a context limit — so it was corrected rather than treated as an echo.

## Qualification

Passed — 3 files verified, 4 requirements traced.

- `tools/checks/qualify.sh` run with `DO_WORK_DIFF_RANGE="306c1f4..a17e6af"` (worktree dispatch mode —
  the working tree is clean post-merge, so the mechanical checks read the merge range).
- **Requirements traced:** 1 → `actions/cleanup.md` Pass 0 step 5, one added sentence. 2 →
  `actions/forensics.md` Check 1, clause folded into the existing move sentence. 3 →
  `docs/work-guide.md` lines 66 and 119, both at guide altitude, no procedure. 4 → verified by
  reading: each site says *that* the removal is owed and points at the canonical home; neither
  restates when or why, which is the copy that would go stale.
- **Post-merge verification:** the acceptance checks were re-run against the merged main tree, not the
  builder branch — `grep` confirms all three edits present at their merged paths, and the contract
  suite's pre-existing failure count is unchanged (7 update-script probes, unrelated).

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ No new failures — the same 7 pre-existing update-script probe failures as on the base
branch (documented under REQ-083's Discovered Tasks), and no prose-contract regression from these
edits.

Non-behavioral change (documentation and action-file prose), so red-green validation does not apply;
regression evidence is the contract suite plus direct reading of the merged tree.

*Verified by work action*

## Lessons Learned

**What worked:** Citing the canonical home instead of restating the rule kept all three edits to one
sentence each — requirement 4 turned out to be a size constraint as much as a correctness one.

**Worth knowing:** The user-guide sentence was the most valuable of the three despite looking like the
smallest. `actions/*` sites are read by agents that can open the cited rule; `docs/work-guide.md` is
read by a human who cannot, so a half-true sentence there has no correction path.

## Orientation

**Now the two out-of-pipeline ways a REQ can leave `do-work/working/` — cleanup's terminal sweep and
forensics' manual reset — tell their reader to drop the checkpoint's in-progress entry too**, so a
normally-completed REQ stops producing a bogus "named in the in-progress record but not in working/"
warning on the next run. The user guide now describes checkpoints honestly: written at claim time,
refreshed at session end.

Leaf change — no new module, contract, or data flow. `prime_files` is empty, so no prime staleness
check applies.

## Full Context

Found by `actions/review-work.md`'s Restatement Sweep during REQ-077's review — see that REQ's
`## Review` → Findings.
