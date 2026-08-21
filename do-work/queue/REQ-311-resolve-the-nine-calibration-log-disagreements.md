---
id: REQ-311
title: "[impact-user-visible] Resolve the nine calibration-log rows that disagree with their frontmatter"
status: pending-answers
created_at: 2026-08-21T00:07:13Z
status_changed_at: 2026-08-21T00:07:13Z
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

- [ ] I discovered this while working on REQ-281: nine calibration-log rows disagree with their REQ frontmatter, six of them by 20 minutes or more, and the estimator fits its scoring table from that corpus. Resolving each needs a judgment about which record is right, which the probe cannot make. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the outlier rule already excludes spans over four hours at read time, and an 8% corpus error may be inside the noise the P50 estimate already carries.

**How to decide it, if you say yes**, because "fix the rows" is not one decision but three:

1. **The three-way tie (REQ-241, REQ-243, REQ-245).** All three log different spans against one shared `claimed_at`. Check whether that instant is a fan-out wave's dispatch time. If it is, the log is the better record for all three and the *frontmatter* wants correcting — which also means the ordering probe REQ-280 added is reading three REQs' spans wrong today.
2. **REQ-233's 70-vs-10.** A 60-minute gap on a single REQ with no shared instant. Read its archived body: if the work described plainly took an hour, the log is right and its `completed_at` was re-stamped later by a repair pass.
3. **The three 2-minute rows (REQ-265, REQ-267, REQ-261).** Just past the one-minute tolerance and plausibly ordinary truncation-plus-a-re-stamp. Consider whether these are worth touching at all, or whether the tolerance should be two minutes — that is a smaller change than nine edits, and it is the delete-before-you-add answer if the 2-minute class turns out to be systematic.

Do not batch-rewrite the file. Whichever record loses in each case, the losing value is evidence about a real defect in how a stamp got written, and it is worth understanding before it is erased.
