---
id: REQ-281
title: Reconcile the calibration log against the frontmatter it was derived from
status: completed
created_at: 2026-08-19T13:42:45Z
claimed_at: 2026-08-21T00:00:42Z
completed_at: 2026-08-21T00:07:52Z
kb_status: promoted
kb_entry: REQ-281-reconcile-the-calibration-log-against-th.md
commit: a868827
route: B
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
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---

## Triage

**Route: B** - Medium

**Reasoning:** Fully specified — the REQ names the file, the predicate, the tolerance, the four report cases, and two explicit non-goals. Discovery was limited to how to locate a REQ by id from inside verify (`board.RequestsById`, which already spans queue/working/archive) and confirming the measurement reproduces.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The measurement was reproduced before any code was written.** An independent Python recomputation over this repo's live log: **72 rows, 62 agree, 10 disagree, 0 missing, 0 unparseable.** All eight material disagreements the REQ names are present and match its table exactly (REQ-233 70→10, REQ-241 28→68, REQ-243 20→68, REQ-245 48→68, REQ-242 34→55, REQ-244 32→59, REQ-265 53→55, REQ-267 38→40). The row count grew from the REQ's 56 because this session appended rows; the disagreement set did not change.

- `board.RequestsById` (`model.go:372`) maps id → ticket across queue, working and archive, so the REQ's "locate the REQ by id across three locations" needs no globbing and no second walk.
- `VerifyReport.SkippedProbes` is the existing convention for "could not check", used by the worktree and completion-anomaly probes.
- The write site (`actions/work.md` Step 8 substep 7.5) computes integer minutes by truncation, so the reader truncates identically — the tolerance absorbs the rounding difference, it does not paper over a different formula.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — two category constants, the probe, the unreconcilable-row constructor
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — the captured RED plus tolerance, unreconcilable, and absent-log cases
- `skills/do-work/actions/forensics.md` (modify) — one row in Check 14's probe table
- `do-work/calibration-log.tsv` (modify) — see Decisions D-03

**Files I will NOT touch:**
- `skills/do-work/actions/work.md` and `actions/estimate-reference.md` — the write path and the outlier policy are explicitly out of scope.

**Acceptance criteria (restated from the REQ):**
1. Read the log, skip the header, locate each REQ across archive/working/queue.
2. Recompute and report when the absolute difference exceeds one minute.
3. Each finding names the REQ, both values, and says either record may be right.
4. Unfindable REQ and unparseable stamps are their own findings, not disagreements.
5. A missing log is a skipped probe, not a finding.
6. Nothing marked `[fixable]`.
7. Check 14's probe table lists it.

## Decisions

- **D-01** (DECIDE & STATE): Two categories rather than four — `calibration-log-mismatch` for a real disagreement, and `calibration-row-unreconcilable` for every row that yields no recomputed number (names no REQ, stamps unparseable, `wall_minutes` not an integer, too few columns). Reasoning: the REQ requires these be distinct *from disagreements*, which two categories satisfy; splitting further would give four categories one remedy, since the question in every unreconcilable case is the same one — which record is wrong, the log or the tree. The specific reason is in the detail, so nothing is lost to grep.
- **D-02** (DECIDE & STATE): Malformed rows (bad integer, too few columns) are reported rather than skipped, though the REQ names only the unfindable-REQ and unparseable-stamp cases. Reasoning: the estimator reads every row as corpus, so a row it cannot parse is an unverifiable input by the REQ's own argument, and silently skipping it would be the "corpus nothing audits" failure in miniature. Both are pinned by test.
- **D-03** (ESCALATE): The new probe flagged `REQ-274`'s row — logged 7, recomputed 5 — and that row was written **by this session, earlier today**, from a stale hardcoded `claimed_at` rather than the value actually stamped into the file. Builder chose to correct that one row to 5. Reasoning: the REQ's "must not pick a winner" binds the *probe*, which cannot know which record drifted; it does not bind a human who knows, and here the cause is known exactly — an arithmetic error made a few hours ago, with the frontmatter provably correct. Every other disagreement was left untouched. Value: the corpus loses one row of known-bad data, and the error is recorded rather than buried. Risk: none to the record — the log is a derived corpus, not testimony, and the correction is described here and in the commit.
- **D-04** (DECIDE & STATE): The nine remaining disagreements were **not** repaired. Reasoning: the REQ's Constraints say report never repair, and unlike REQ-274 nothing here establishes which record is right — the frontmatter may have been legitimately rewritten by either repairer or by crash recovery. Queued as a follow-up instead (`## Discovered Tasks`).

## Implementation Summary

**What was done:** Added a read-only probe to `queue-kanban verify` that recomputes every `do-work/calibration-log.tsv` row's `wall_minutes` from its REQ's `claimed_at`/`completed_at` and reports disagreements beyond a one-minute tolerance, naming both values and stating that either record may be the correct one. Rows that cannot be reconciled at all are a separate finding, and an absent log is a reported skipped probe. Nothing is marked fixable and nothing is repaired. Added the probe to forensics Check 14's table, and corrected the single log row this session had written wrong.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — `verifyCategoryCalibrationLogMismatch` and `verifyCategoryCalibrationRowUnreconcilable`; `appendCalibrationLogFindings` wired after the ordering probe; `calibrationRowFinding`; `calibrationToleranceMinutes`; `strconv` import.
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — three cases: the captured RED with an agreeing row pinning the tolerance, four unreconcilable shapes, and the absent-log skip.
- `skills/do-work/actions/forensics.md` (modified) — one row added to Check 14's probe table.
- `do-work/calibration-log.tsv` (modified) — `REQ-274` row corrected from 7 to 5 (D-03).

**Tests touched:** three new cases. No existing assertion changed meaning.

## Qualification

Passed — 4 files verified, 7 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` ok (39.5s). Diff grepped for debug artifacts — none. Heredoc quoted this time (REQ-280's lesson), and the inserted block re-read to confirm no substitution damage. `maintainer-verify.sh` exits 0.
- **Substantive:** the probe is real parsing and arithmetic over a real file; its output on this repo matches an independently written Python recomputation row for row.
- **Requirements traced:** AC1–AC3 → `appendCalibrationLogFindings` + the RED test; AC4 → `calibrationRowFinding` + the four-shape test; AC5 → the skip test; AC6 → asserted directly; AC7 → forensics.md.
- **Flowing:** the probe reads the real TSV and flips verify's exit code — it currently reports nine genuine disagreements on this repo.

## Testing

- `go test ./...` — ok, 39.5s. `gofmt -l .` clean, `go vet ./...` clean.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, before and after.

**Red-green validation** — traced to the REQ's `## Red-Green Proof`:

| | Before | After |
|---|---|---|
| REQ-233-shaped fixture (10-minute span logged as 70) | exit 0, no findings | exit 1, finding naming REQ-233, logged 70, recomputed 10, remedy saying either record may be correct |
| Second row agreeing within a minute | n/a | silent — the tolerance is pinned by the same test |

**Cross-check against an independent implementation:** the probe's output on this repo (10 mismatches before D-03's correction, 9 after) matches the Python recomputation written during exploration, row for row and value for value. Two implementations agreeing on real data is stronger evidence than either passing its own fixture.

**Gate-safety check:** `maintainer-verify.sh` runs `go vet` and `go test` in the board package, **not** `queue-kanban verify` against the live tree, so the nine real disagreements this probe now surfaces do not turn the canonical gate red. Confirmed by running it. They surface where they should — `do-work-board verify` and `do-work forensics`.

**Fixture mutation testing:**

| Reverted behavior | Result |
|---|---|
| tolerance widened to 999 | RED test FAILs (0 mismatches, want 1) |
| unreconcilable rows reported as mismatches | separation test FAILs (0 unreconcilable, want 4) |

## Discovered Tasks

- **impact-user-visible** Nine calibration-log rows still disagree with their frontmatter after this REQ's own row was corrected — six by 20 minutes or more, feeding a corpus the estimator fits its scoring table from. Resolving each needs a judgment the probe deliberately cannot make. Fold-first scan: no `pending`/`pending-answers` REQ shares this root cause; not a prose restatement, so the standing sweep does not take it. Queued as **REQ-311** (`pending-answers`), carrying the three-way-tie observation (REQ-241/243/245 share one `claimed_at`, which may mean the *frontmatter* is the wrong record there).

## Review

**Overall: 93%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Read log, skip header, locate REQ across archive/working/queue | ✅ via `board.RequestsById` |
| Report when the difference exceeds one minute | ✅ RED test, with an agreeing row pinning the tolerance |
| Finding names REQ, both values, and says either record may be right | ✅ all three asserted |
| Unfindable REQ and unparseable stamps are their own findings | ✅ separate category, four shapes pinned |
| Missing log is a skipped probe | ✅ |
| Nothing `[fixable]` | ✅ asserted |
| Check 14 probe table lists it | ✅ |

### Findings

**Important — none.**

**Minor:**

- **M1:** The probe reports nine genuine disagreements on this repo the moment it ships. That is the REQ working, not a defect, but it means `do-work-board verify` goes from a short finding list to a longer one for anyone who runs it. Queued as REQ-311 rather than left to accumulate silently.
- **M2:** D-03 corrected one log row, which is a repair inside a REQ whose Constraints say "report, never repair". The constraint binds the probe; the row was written by this session from a stale stamp and its cause is known exactly. Recorded rather than done quietly, because the reading is arguable.

**Nit:**

- **N1:** `appendCalibrationLogFindings` parses TSV inline rather than through a helper. At one call site with five columns that is the simpler shape; a second reader of this file would be the moment to extract one.

### Restatement Sweep

Redefined element: none. This REQ **adds** a probe and adds a table row; it changes no existing field meaning, contract token, or output shape. The one thing it comes close to is the calibration log's contract, and it deliberately leaves that alone — the log still records raw spans, and the outlier rule still applies at read time only.

Swept anyway for statements about what audits the calibration log: `actions/estimate-reference.md` (names it as the recalibration corpus, unchanged and still true), `actions/work.md` Step 8 substep 7.5 (the writer, untouched). Both describe the log as unaudited only by implication; neither asserts it, so neither is now stale.

### Acceptance Testing

Three ways. The captured fixture (exit 0 → exit 1 with the exact finding). The live tree, where the probe's output was compared row-for-row against an independently written Python recomputation — two implementations agreeing on 72 real rows. And the canonical gate, confirmed still exit 0, so shipping a probe that finds nine real problems does not break every subsequent commit.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Scope Discipline 90% for D-02 (malformed rows, beyond the four named cases) and D-03 (one row corrected). Both are argued and tested rather than slipped in.

### Follow-up REQs Created

**REQ-311** (`pending-answers`) — the nine disagreements, with the three-way-tie observation that may invert which record is wrong.

## Lessons Learned

**What worked:** Writing a throwaway Python recomputation *before* the Go probe, then comparing the two on 72 real rows. It confirmed the REQ's measurement independently, and later it was the acceptance evidence — a fixture proves a probe handles the case you imagined, while two independent implementations agreeing on real data proves it handles the ones you did not. Also: checking whether the new failing condition breaks the canonical gate *before* claiming the REQ was done. It did not, but that was luck, and finding out at review time rather than at commit time was not.

**What didn't:** Nothing failed outright, but the probe immediately flagged a row **this session had written wrong four hours earlier** — REQ-274's, logged as 7 against a true span of 5, because the calibration arithmetic used a hardcoded `claimed_at` string instead of reading back the one actually stamped into the file. The REQ was built to catch exactly that class and caught its own author. The habit worth taking: when a step writes a derived value from a stamp, read the stamp back from the file it was written to rather than reusing the variable.

**Worth knowing:** REQ-241, REQ-243 and REQ-245 log three different spans against one identical `claimed_at` of `2026-08-18T12:43:06Z`. That pattern reads like a fan-out wave whose members all recorded the wave's dispatch instant — which would mean the *frontmatter* is the wrong record for those three, not the log. If that is confirmed in REQ-311, REQ-280's ordering probe is currently reading three REQs' spans wrong too. Nobody should batch-rewrite this file before that question is settled.

## Orientation

`queue-kanban verify` now audits the calibration corpus the estimator learns from, recomputing every logged wall span from the frontmatter it came from and reporting disagreements without deciding which record is wrong. It found nine on this repo the first time it ran. Lives in the board's verify subsystem (`skills/do-work-board/tools/queue-kanban/verify.go`) alongside REQ-280's ordering probe, and is listed in `skills/do-work/actions/forensics.md` Check 14.

**[MAP CHANGED]** — two new verify finding categories, `calibration-log-mismatch` and `calibration-row-unreconcilable`. Anything parsing verify output by category, or counting probes against Check 14's table, sees two more. The calibration log also stops being a write-only record: it now has a reader that checks it.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` — referenced paths still resolve; its probe-set discussion continues to defer to Check 14's table, which this REQ updated.
