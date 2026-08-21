---
id: REQ-316
title: "[impact-rule-change] Audit the calibration-log write step for the REQ-274 stale-stamp bug class"
status: pending
created_at: 2026-08-21T15:26:47Z
status_changed_at: 2026-08-21T15:26:47Z
user_request: UR-057
addendum_to: REQ-311
domain: general
impact: impact-rule-change
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work.md
---

# Audit the Calibration-Log Write Step for the REQ-274 Stale-Stamp Bug Class

## What

REQ-311 corrected nine `do-work/calibration-log.tsv` rows that disagreed with the REQ frontmatter they were computed from — every one of the nine was resolved in the frontmatter's favor, meaning `do-work/calibration-log.tsv`'s `wall_minutes` column, not the frontmatter, was wrong nine separate times. A tenth row, REQ-274, was already found and corrected earlier the same day inside REQ-281, with its cause known exactly: "the calibration arithmetic used a hardcoded `claimed_at` string instead of reading back the one actually stamped into the file" (REQ-281 `## Lessons Learned`).

REQ-311 fixed the data but was scoped (`write_set: do-work/calibration-log.tsv` only) away from touching the write path — `actions/work.md` Step 8 substep 7.5, which appends the calibration-log line at archive time by computing `completed_at − claimed_at`. This REQ is the follow-up: read that instruction and decide whether it needs hardening against the same mistake recurring a tenth-plus time.

## Why this is worth someone's time

Ten wrong rows out of seventy-two logged (REQ-274 plus the nine REQ-311 closed) is roughly 14% of the corpus the estimator (`actions/estimate-reference.md`) fits its scoring table from — caught and hand-corrected twice now, after the fact, row by row. `do-work` is a set of instructions an agent follows, not compiled code with a literal function to patch; the fix, if one is needed, is a clearer instruction, not a code change. That is exactly the class of gap `crew-members/maintenance.md`'s deletion-first questions are built to catch: is the instruction ambiguous, or genuinely missing a constraint an agent needs to avoid reusing a value it already holds in context instead of the one just written to disk.

## Requirements

- Read `actions/work.md` Step 8 substep 7.5 (line ~605) and its surrounding steps (6–8) to see what value an agent following it would actually have in hand for `claimed_at` and `completed_at` at the moment it computes the calibration-log line — a value it read earlier in the run and is carrying forward, or a value it re-reads from the archived REQ file on disk at this exact step.
- Per `crew-members/maintenance.md`'s deletion-first questions: check first whether this is already covered by an existing instruction elsewhere in `work.md` (e.g. a general "read stamps back from disk, never carry them forward" rule) that substep 7.5 merely fails to point to — in which case the fix is a cross-reference, not new prose.
- If the instruction is genuinely silent on this, add the minimum clause to substep 7.5 (or the general Timestamp rule in `actions/work-reference.md`, if that is the more central home) requiring the agent to read `claimed_at` and `completed_at` from the just-archived file's frontmatter at the point of computing `wall_minutes`, not from values held earlier in the run.
- Do not touch `do-work/calibration-log.tsv` itself — REQ-311 already reconciled it against frontmatter, and this REQ is about the instruction, not the data.
- Decide whether any other stamp-write site in `work.md` or `work-reference.md` carries the same risk (reusing a context-held timestamp instead of reading the one on disk) and is worth the same clause, or whether the risk is specific to substep 7.5's cross-file read (computing from a REQ file's stamps into a *different* file, `calibration-log.tsv`, which is what made the staleness invisible until REQ-281's probe existed).

## Context

Discovered while completing REQ-311 (`## Discovered Tasks`). REQ-309 already queued a related but distinct process gap (the pipeline's test resolution being area-scoped); this REQ is scoped narrowly to the calibration-log write step and should not widen into a general stamp-write audit unless the last requirement above turns up other genuinely at-risk sites.
