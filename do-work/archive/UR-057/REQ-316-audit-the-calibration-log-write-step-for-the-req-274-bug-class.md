---
id: REQ-316
title: "[impact-rule-change] Audit the calibration-log write step for the REQ-274 stale-stamp bug class"
status: completed
created_at: 2026-08-21T15:26:47Z
status_changed_at: 2026-08-21T19:11:08Z
claimed_at: 2026-08-21T18:55:41Z
completed_at: 2026-08-21T19:11:08Z
commit: 68d9ad9
kb_status: promoted
kb_entry: REQ-316-audit-the-calibration-log-write-step-for.md
route: B
user_request: UR-057
addendum_to: REQ-311
domain: general
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work.md
- _dev/tests/contract-regressions.sh
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-21T18:55:41Z
  basis:
    - trivial short-circuit
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

---

## Triage

**Route: B** — Medium

**Reasoning:** The likely instruction edit is small, but the request explicitly requires a
deletion-first boundary audit and a semantic population check across timestamp readers. Route B
adds independent exploration before the narrow contract change.

**Planning:** Not required for Route B.

**Estimate:** 5 active minutes (P50, high confidence; effort-mechanical short-circuit).

## Exploration

Deletion-first found no stale rule or bad example to remove. The Timestamp rule owns creation of
current instants, but it does not say where a later derived record must obtain already-persisted
stamps. Step 8 stamps `completed_at`, archives the REQ, and only then projects both frontmatter
stamps into `calibration-log.tsv`; its 7.5 instruction is silent about reading that archived source.

The semantic population contains one durable cross-file derivative writer: Step 8 substep 7.5.
Estimate-reference and forensics are readers, staging prose does not compute values, and other stamp
consumers either read current files or do not persist a derivative. The minimum earned addition is
therefore one source-selection clause in 7.5 plus a narrowly extracted semantic contract that pins
both stamps, the just-archived frontmatter source, and the ban on carried context values.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — require calibration arithmetic to read persisted stamps
  from the just-archived REQ
- `_dev/tests/contract-regressions.sh` (modify) — add the deletion-first replay contract and mutations

**Files I will NOT touch:** `do-work/calibration-log.tsv`, the general Timestamp rule,
`estimate-reference.md`, forensics/board readers, or any archived historical data.

**Acceptance criteria (restated from REQ):**
- [x] Step 8.7.5 reads both `claimed_at` and `completed_at` from the just-archived REQ frontmatter at
  calculation time and forbids reuse of context-held values.
- [x] The deletion-first audit records why no existing rule can be narrowed or cross-referenced.
- [x] The semantic population proves no other stamp-derived durable writer needs the clause.
- [x] A narrowly bound contract fails when the clause, either stamp, persisted source, frontmatter,
  or no-carried-value boundary is removed or displaced.
- [x] Calibration data and broader timestamp semantics remain unchanged.
- [x] The direct canonical repository gate exits 0.

## Pre-Flight

**Git:** ✓ Clean outside this REQ's claim bookkeeping.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh`
**Dependencies:** ✓ Bash, Python, Go 1.26.1, ShellCheck 0.11.0, and Chromium are available.

*Checked by work action*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the maintenance and action/shell primes; isolate 7.5 as the sole derivative
  writer and design a clause-local semantic mutation matrix.
- [x] **[APPLY]:** Add the contract RED-first, then the minimum source-selection clause.
- [x] **[UNIFY]:** Review exact scope, replay mutations, run focused and qualification gates, and remove
  scratch artifacts.

## Decisions

<!-- D-XX counter: next is D-03 -->

### D-01: Add at the derivative writer, not the general clock rule

**Decision:** Keep the persisted-source constraint in Step 8 substep 7.5. The Timestamp rule decides
how to obtain a current instant; it does not own which already-persisted record a later cross-file
projection reads. No stale source, bad example, broad tool, or vague job could be removed or narrowed
to supply that missing boundary.

**Reasoning:** Step 7.5 is the sole durable writer that derives one file's value from a REQ's stamps.
Moving the constraint to the general rule would broaden every timestamp write without addressing a
shared failure class.

### D-02: Test the local directive and replay its semantic losses

**Decision:** Extract only substep 7.5, locate its persisted-source directive by meaning, and mutate
the live directive in memory.

**Reasoning:** Whole-file keyword checks let unrelated timestamp prose satisfy a weakened write site.
The mutation matrix separately removes the clause, each stamp, the just-archived source, frontmatter,
and the carried-value prohibition, and also displaces the intact clause outside 7.5.

## Implementation Summary

- `skills/do-work/actions/work.md` (modified) — added the earned source-selection clause at calibration-log
  calculation time: read both stamps from the just-archived REQ frontmatter and never reuse
  context-held values.
- `_dev/tests/contract-regressions.sh` (modified) — added a clause-local semantic contract plus eight in-memory
  mutations covering deletion, stamp omissions, source weakening, frontmatter loss, ban deletion or
  inversion, and displacement.

## Qualification

- Scope remained exactly the two declared implementation files.
- `tools/checks/qualify.sh`: pass; `tools/checks/scope-drift.sh`: Implementation Summary matches
  the Scope declaration.
- Deletion-first and semantic-population findings are recorded in `## Exploration` and D-01.
- The change adds one sentence at the only cross-file derivative writer; no calibration data,
  general Timestamp rule, or other reader/writer changed.

## Testing

- RED: after adding only the semantic contract, `bash _dev/tests/contract-regressions.sh` failed with
  `incomplete persisted-stamp source contract: ['directive']` and the REQ-316 failure line.
- GREEN: after adding the 7.5 clause, the same focused command exited 0 with
  `Contract regression checks passed.`
- Mutation replay: all eight in-memory mutations were rejected — clause deletion; `claimed_at`
  omission; `completed_at` omission; generic context source; frontmatter removal; carried-value ban
  deletion; ban inversion; and moving the intact clause outside substep 7.5.
- `bash -n _dev/tests/contract-regressions.sh`: pass.
- `shellcheck _dev/tests/contract-regressions.sh`: no new findings; exit 1 is from the suite's existing
  SC2030/SC2031/SC2329 informational findings.
- `git diff --check`: pass.
- Direct unpiped canonical `bash _dev/tests/maintainer-verify.sh` with the declared Chromium headless
  shell: pass; exit 0.

## Discovered Tasks

None.

## Review

**Overall: 100%** | 2026-08-21T19:10:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None.

**Minor findings:** None.
**Acceptance:** Pass — the rule is local to the sole durable derivative writer and all eight
source/stamp mutations fail non-vacuously.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Deletion-first separated stamp creation from derived-record source selection, showing that a
  cross-reference to the general Timestamp rule would be true but insufficient.
- Extracting exactly substep 7.5 and displacing the intact clause outside it proved the contract
  guards the operative writer rather than nearby vocabulary.

**What didn't:**
- The original instruction named the correct arithmetic but left the value source implicit. That
  let an earlier context-held timestamp look interchangeable with the persisted frontmatter stamp,
  producing a valid-looking but wrong calibration row.

**Worth knowing:** A durable cross-file projection needs to read the record that actually landed on
disk at the point of derivation. Stamp-generation rules do not automatically establish that
readback boundary, and broad whole-step greps cannot prove it is attached to the writer.

## Orientation

**[MAP CHANGED]** Calibration logging now treats the archived REQ frontmatter as the authority for
both endpoints of `wall_minutes`. The archive step re-reads `claimed_at` and `completed_at` at
calculation time and explicitly rejects values carried from earlier in the run; the broader
Timestamp rule and historical calibration corpus remain unchanged.
