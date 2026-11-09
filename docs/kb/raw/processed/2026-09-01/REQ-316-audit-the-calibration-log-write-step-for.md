---
source_type: req_lesson
req_id: REQ-316
req_path: do-work/archive/UR-057/REQ-316-audit-the-calibration-log-write-step-for-the-req-274-bug-class.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, audit, calibration, write, step]
---

# Lessons from REQ-316: Audit the calibration-log write step for the REQ-274 stale-stamp bug class

## What the REQ was about

REQ-311 corrected nine `do-work/calibration-log.tsv` rows that disagreed with the REQ frontmatter they were computed from — every one of the nine was resolved in the frontmatter's favor, meaning `do-work/calibration-log.tsv`'s `wall_minutes` column, not the frontmatter, was wrong nine separate times. A tenth row, REQ-274, was already found and corrected earlier the same day inside REQ-281, with its cause known exactly: "the calibration arithmetic used a hardcoded `claimed_at` string instead of reading back the one actually stamped into the file" (REQ-281 `## Lessons Learned`).

REQ-311 fixed the data but was scoped (`write_set: do-work/calibration-log.tsv` only) away from touching the write path — `actions/work.md` Step 8 substep 7.5, which appends the calibration-log line at archive time by computing `completed_at − claimed_at`. This REQ is the follow-up: read that instruction and decide whether it needs hardening against the same mistake recurring a tenth-plus time.

## Solution summary

**[MAP CHANGED]** Calibration logging now treats the archived REQ frontmatter as the authority for
both endpoints of `wall_minutes`. The archive step re-reads `claimed_at` and `completed_at` at
calculation time and explicitly rejects values carried from earlier in the run; the broader
Timestamp rule and historical calibration corpus remain unchanged.

## What worked

- Deletion-first separated stamp creation from derived-record source selection, showing that a
  cross-reference to the general Timestamp rule would be true but insufficient.
- Extracting exactly substep 7.5 and displacing the intact clause outside it proved the contract
  guards the operative writer rather than nearby vocabulary.

## What didn't work

- The original instruction named the correct arithmetic but left the value source implicit. That
  let an earlier context-held timestamp look interchangeable with the persisted frontmatter stamp,
  producing a valid-looking but wrong calibration row.

## Worth knowing

A durable cross-file projection needs to read the record that actually landed on
disk at the point of derivation. Stamp-generation rules do not automatically establish that
readback boundary, and broad whole-step greps cannot prove it is attached to the writer.

## Back-reference

See `do-work/archive/UR-057/REQ-316-audit-the-calibration-log-write-step-for-the-req-274-bug-class.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `68d9ad9`.
