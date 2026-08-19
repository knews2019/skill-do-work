---
id: REQ-281
title: Reconcile the calibration log against the frontmatter it was derived from
status: pending
created_at: 2026-08-19T13:42:45Z
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-280]
maintenance: false
related: [REQ-279, REQ-280, REQ-282, REQ-283]
batch: upstream-consumer-report-2026-08-19
write_set:
- skills/do-work-board/tools/queue-kanban/verify.go
- skills/do-work-board/tools/queue-kanban/verify_test.go
- skills/do-work/actions/forensics.md
---

# Reconcile the Calibration Log Against the Frontmatter It Was Derived From

## What

`do-work/calibration-log.tsv` is an independent third record of every REQ's wall span, written once at archive time by `actions/work.md` Step 8 substep 7.5 as `completed_at − claimed_at`, and read back by `actions/estimate-reference.md:94` as the corpus the scoring table is fit from. Nothing ever compares it against the frontmatter it came from.

Add a `queue-kanban verify` probe that recomputes each row's `wall_minutes` from its REQ's `claimed_at` and `completed_at` and reports rows that disagree by more than a minute.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

**It reproduces here, in the maintainer repo.** Recomputing all 56 rows of this repo's `do-work/calibration-log.tsv` from the archived REQs: 42 agree, 14 disagree. Six are ±1, the ordinary integer-truncation gap. Eight are material:

| REQ | logged | recomputed from frontmatter |
| --- | --- | --- |
| REQ-233 | 70 | 10 |
| REQ-241 | 28 | 68 |
| REQ-243 | 20 | 68 |
| REQ-245 | 48 | 68 |
| REQ-242 | 34 | 55 |
| REQ-244 | 32 | 59 |
| REQ-265 | 53 | 55 |
| REQ-267 | 38 | 40 |

The consumer repo showed the same shape at larger scale: 54 rows, 39 agree, 15 disagree, nine materially. A corpus nothing audits is a corpus that quietly decays, and this one feeds every future estimate.

## Context

Note which record is wrong is *not* knowable from the disagreement alone. The log line is written once and never revised; the frontmatter can legitimately be rewritten afterwards — by `scripts/repair-req-timestamps.sh` at session start, by `scripts/audit-archive-timestamps.sh --fix`, or by a crash-recovery pass that cleared and re-stamped a claim. So the probe reports the disagreement and names both values. It must not pick a winner, and it must not rewrite either record.

`actions/estimate-reference.md:94` and `actions/work.md:591` both state the existing contract deliberately: the log records **raw** wall spans, and the outlier rule (over 4h or negative ⇒ excluded) is applied at read time, never at write time. This probe does not touch that; it audits fidelity between two records, not the outlier policy.

## Detailed Requirements

- Read `do-work/calibration-log.tsv`, skipping the header line. For each row, locate the REQ by id across `do-work/archive/` (including `archive/UR-NNN/`), `do-work/working/`, and `do-work/queue/`.
- Recompute `completed_at − claimed_at` in integer minutes and compare to the row's `wall_minutes`. Report a finding when the absolute difference exceeds one minute — the one-minute tolerance absorbs truncation-versus-rounding, which is noise, not damage.
- Each finding names the REQ id, the logged value, the recomputed value, and states that either record may be the correct one.
- Report a row whose REQ cannot be found, and a row whose REQ has an absent or unparseable `claimed_at`/`completed_at`, as their own distinct findings — not as disagreements.
- A missing `do-work/calibration-log.tsv` is not a finding; report it as a skipped probe with the reason, per the file's existing skipped-probe convention.
- Do not mark any of these `[fixable]` — no cleanup pass resolves them, and the resolution requires judgment about which record is right.
- Add the probe to `actions/forensics.md` Check 14's probe table.

## Constraints

- **Report, never repair.** `verify` is read-only, and both records here are load-bearing history. Even an obviously-wrong log row stays as written.
- **Do not change the write path.** `actions/work.md` Step 8 substep 7.5 keeps appending raw spans; this REQ adds no validation, no retry, and no guard at write time.
- **Do not build the fabrication heuristic** (upstream remedy 2) — declined at triage, see UR-057 Batch Constraints and REQ-280.
- Depends on REQ-280: both edit `verify.go`, and serializing them keeps the two probes from colliding in the same file.

## Red-Green Proof

**RED prompt/case:** Build a verify fixture with an archived REQ whose `claimed_at`/`completed_at` span 10 minutes, and a `do-work/calibration-log.tsv` row claiming `wall_minutes` of 70 for that REQ — the REQ-233 shape observed in this repo. Run `queue-kanban verify --repo-root <fixture>`.
**Why RED now:** It exits 0. Nothing in the suite reads `calibration-log.tsv` except the estimator at recalibration time, and that path re-fits the scoring table from the rows without ever checking them against their source.
**GREEN when:** The same fixture exits 1 with a finding naming REQ-233, logged 70, recomputed 10, and stating that either record may be correct. A second row in the same fixture that agrees within a minute produces no finding, so the tolerance is pinned too.
**Validation:** User confirmed by accepting the triage; the 14-of-56 disagreement above was measured against this repo's live `do-work/calibration-log.tsv` during that triage, and the report's independent 15-of-54 measurement on the consumer tree is the second observation.

## Full Context

See `do-work/user-requests/UR-057/input.md` for the complete verbatim upstream report.

---
*Source: upstream defect report D1, severity high, remedy 3, from `g1w-game-find-the-difference` running v0.212.25 — verbatim claim: "the log is an independent third record that nothing reconciles against frontmatter … Add a calibration-log reconciliation probe. Recompute each row's `wall_minutes` from the REQ's frontmatter and report rows that disagree by more than a minute. A corpus nothing audits is a corpus that quietly decays." Accepted by `do-work-toolbox validate-feedback` triage (2026-08-19). Evidence: `skills/do-work/actions/work.md:591` is the only writer; `skills/do-work/actions/estimate-reference.md:94` names the log as the recalibration input; measured 42 agree / 14 disagree over this repo's 56 rows, eight materially. Surface-cost: Earned — incident is 14 disagreeing rows here and 15 there, feeding the estimator uncontested; surface is one read-only probe with no write path; cheaper than the alternative of trusting an unaudited corpus indefinitely; test is the REQ-233-shaped fixture above plus an agreeing row to pin the tolerance.*
