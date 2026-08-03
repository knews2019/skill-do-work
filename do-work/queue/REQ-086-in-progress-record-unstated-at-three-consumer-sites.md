---
id: REQ-086
title: The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide
status: pending
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

## Full Context

Found by `actions/review-work.md`'s Restatement Sweep during REQ-077's review — see that REQ's
`## Review` → Findings.
