---
id: REQ-311
title: "[impact-user-visible] Resolve the nine calibration-log rows that disagree with their frontmatter"
status: completed
created_at: 2026-08-21T00:07:13Z
completed_at: 2026-08-21T15:20:33Z
status_changed_at: 2026-08-21T15:20:33Z
user_request: UR-057
addendum_to: REQ-281
domain: general
impact: impact-user-visible
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- do-work/calibration-log.tsv
---

# Resolve the Nine Calibration-Log Rows That Disagree With Their Frontmatter

## What

REQ-281 built the probe that finds them; finding them was its whole scope. Nine rows in `do-work/calibration-log.tsv` still disagree with the frontmatter they were derived from, and the probe deliberately does not pick a winner. Somebody has to.

Run `do-work-board verify` for the current list. As of REQ-281 it was:

| REQ | logged | recomputed from frontmatter | gap |
| --- | --- | --- | --- |
| REQ-233 | 70 | 10 | 60 |
| REQ-241 | 28 | 68 | 40 |
| REQ-243 | 20 | 68 | 48 |
| REQ-245 | 48 | 68 | 20 |
| REQ-242 | 34 | 55 | 21 |
| REQ-244 | 32 | 59 | 27 |
| REQ-265 | 53 | 55 | 2 |
| REQ-267 | 38 | 40 | 2 |
| REQ-261 | 12 | 14 | 2 |

## Why this is worth someone's time

`actions/estimate-reference.md` fits the estimator's scoring table from this file. Six of the nine are off by 20 minutes or more, and three of those (REQ-241, REQ-243, REQ-245) share the identical `claimed_at` of `2026-08-18T12:43:06Z` while logging three different spans — which looks less like drift and more like a fan-out wave whose members all recorded the wave's claim instant. If that reading is right, the frontmatter is the wrong record for those three and the log is correct, which is the opposite of the obvious assumption.

Every estimate the pipeline prints is fit from these rows. Six materially wrong rows in seventy-two is roughly 8% of the corpus.

## Context

Discovered by REQ-281's probe on its first run against this repo. One tenth row, `REQ-274`, was corrected inside REQ-281 (its `## Decisions` D-03) because that one's cause was known exactly — this session wrote it from a stale stamp hours earlier. The nine here have no such provenance.

## Open Questions

- [x] I discovered this while working on REQ-281: nine calibration-log rows disagree with their REQ frontmatter, six of them by 20 minutes or more, and the estimator fits its scoring table from that corpus. Resolving each needs a judgment about which record is right, which the probe cannot make. Should I process this as a new task? → Confirmed: Yes, resolve now
  *(2026-08-21)* User confirmed via `do-work clarify`, and asked for it to be resolved in this same session rather than left in the queue for a later work run.

**How to decide it, if you say yes**, because "fix the rows" is not one decision but three:

1. **The three-way tie (REQ-241, REQ-243, REQ-245).** All three log different spans against one shared `claimed_at`. Check whether that instant is a fan-out wave's dispatch time. If it is, the log is the better record for all three and the *frontmatter* wants correcting — which also means the ordering probe REQ-280 added is reading three REQs' spans wrong today.
2. **REQ-233's 70-vs-10.** A 60-minute gap on a single REQ with no shared instant. Read its archived body: if the work described plainly took an hour, the log is right and its `completed_at` was re-stamped later by a repair pass.
3. **The three 2-minute rows (REQ-265, REQ-267, REQ-261).** Just past the one-minute tolerance and plausibly ordinary truncation-plus-a-re-stamp. Consider whether these are worth touching at all, or whether the tolerance should be two minutes — that is a smaller change than nine edits, and it is the delete-before-you-add answer if the 2-minute class turns out to be systematic.

Do not batch-rewrite the file. Whichever record loses in each case, the losing value is evidence about a real defect in how a stamp got written, and it is worth understanding before it is erased.

## Implementation

**All nine rows: frontmatter wins, log loses.** Investigated with a source outside both disputed records — `git log`, which cannot have been rewritten by either a repair pass or a buggy log-write step.

1. **The three-way tie (REQ-241/243/245) and the pair (REQ-242/244).** Found the shared `claimed_at` really is a fan-out wave's dispatch instant — each pair traces to one `[queue] claim REQ-... for wave N` commit (`2432f45` for 241/243/245, `2ad71eb` for 242/244), timestamped within a few minutes of the frontmatter `claimed_at` it produced. That confirms the *opposite* of the REQ's working hypothesis: the shared claim is real, accurate data, not a corruption — so recomputing wall time from it is legitimate, not an error. The large gap between each REQ's code-merge commit and its stamped `completed_at` (13–35 minutes for 241/243/245) is explained by real work, not a stamping bug: both REQs carry a `## Review` section reading `PASS-WITH-FINDINGS (first pass, pre-remediation)` — the findings-and-remediation cycle after merge and before archive accounts for the gap. REQ-242 and REQ-244's gaps (2–3 minutes) are unremarkable.
2. **REQ-233's 70-vs-10.** The archived body's own git trail settles this outright: the claim commit (`b2d5840`) lands at `2026-08-18T10:51:30Z`, the merge (`9b2578b`) at `11:05:11Z`, and the review footer inside the REQ's own body is timestamped `11:09:44Z` — all three land inside the ~10-13 minute window the frontmatter states, and none is anywhere near a 70-minute span. The log's `70` has no corroborating evidence anywhere and matches the exact bug class REQ-281 already found and fixed once this same day in REQ-274 (`## Decisions` D-03): calibration arithmetic reading a stale/wrong `claimed_at` instead of the one actually stamped into the file.
3. **The three 2-minute rows (REQ-261/265/267).** Each of these carries a `**Overall: NN%** | <timestamp>` review footer, an independent record generated mid-REQ, well before archiving. For all three, `review timestamp − claimed_at` rounds to the *frontmatter*-recomputed wall_minutes, not the logged one (REQ-261: 14m16s ≈ 14, logged 12; REQ-267: 39m52s ≈ 40, logged 38; REQ-265: 54m42s ≈ 55, logged 53) — a second independent source agreeing with frontmatter over the log in all three cases.

No case produced evidence for the log over the frontmatter. `do-work/calibration-log.tsv` was corrected to the frontmatter-recomputed `wall_minutes` for all nine rows (233→10, 241→68, 243→68, 245→68, 242→55, 244→59, 261→14, 267→40, 265→55). `queue-kanban verify` now reports zero `calibration-log-mismatch` findings and exits 0.

## Testing

**Tests run:** `queue-kanban verify --repo-root .` (before: 9 `calibration-log-mismatch` findings; after: 0, exit 0) and `bash _dev/tests/maintainer-verify.sh` (canonical gate).
**Result:** both exit 0.

## Decisions

- **D-01**: Root cause of the log's own errors was not pursued to a code fix — out of this REQ's `write_set` (`do-work/calibration-log.tsv` only) and REQ-281's constraint that the write path (`actions/work.md` Step 8 substep 7.5) stays untouched here. REQ-274's precedent (stale-`claimed_at` arithmetic) is the most likely shared cause across all nine, but nothing here confirms it for certain — recorded as a discovered task rather than assumed.

## Discovered Tasks

- [low] `do-work/calibration-log.tsv`'s write step (`actions/work.md` Step 8 substep 7.5) has now been shown wrong on 9 of 72 rows (REQ-274 was a 10th, already fixed) in this repo's own history, all from the same working day. Worth a REQ that reads the write step's source and confirms or rules out the REQ-274 stale-stamp bug as the shared cause, so it can be fixed once instead of caught row-by-row after the fact.

## Orientation

`do-work/calibration-log.tsv` now agrees with the REQ frontmatter it was derived from on every row `queue-kanban verify` checks. The estimator (`actions/estimate-reference.md`) fits its scoring table from this file, so its corpus is now clean rather than carrying six rows off by 20+ minutes. Lives alongside REQ-281's probe in the board's verify subsystem; no code changed, only data.

**[MAP CHANGED]** — none. This REQ corrects data the probe already described as a known gap; it changes no contract, no probe behavior, no schema.
