---
id: REQ-212
title: Record estimate-vs-wall calibration log at archive time
status: pending
created_at: 2026-08-17T08:05:43Z
user_request: UR-048
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-211, REQ-213]
batch: estimator-calibration
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/estimate-reference.md]
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-17T08:05:43Z
  basis:
    - Route B
    - 2-file write set
    - (priced with the pre-calibration table)
---

# Record Estimate-vs-Wall Calibration Log at Archive Time

## What

When the work action archives a REQ that carries an `estimate:` block, append one line to `do-work/calibration-log.tsv` recording the estimate against the observed wall span — so the next table re-fit reads a log instead of mining git history.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- In work.md's archive step (success path), after the terminal frontmatter flip: if the REQ has an `estimate:` block with `p50_active_minutes` and both `claimed_at` and `completed_at` parse, append one TSV line to `do-work/calibration-log.tsv`: `req_id  route  estimated_p50_minutes  wall_minutes  completed_at` (wall = completed − claimed, integer minutes). Create the file with a header line on first write.
- **Raw wall time is recorded deliberately** — no outlier judgment at write time. Analysis applies the outlier rule (>4h ⇒ assumed pause, excluded) at read time; the reference file documents this split.
- No estimate block, or unparseable stamps → no line, no error, never blocks archiving.
- The log line is appended by the orchestrator and staged with the archive move in the commit phase.
- `estimate-reference.md` Calibration Honesty section: name the log, its columns, and the read-time outlier rule as the recalibration input.
- This is a work-action write, not a board write — queue-kanban neither reads nor writes the log in this REQ.

## Constraints

- Batch constraint: must not add a queue-kanban write surface; the three-write-surface sentence stays untouched.

## Builder Guidance

Firm on raw-recording + read-time filtering; latitude on exact column order and the header wording.

## Red-Green Proof
**RED prompt/case:** Archive a REQ carrying an `estimate:` block today — no calibration record exists anywhere; comparing estimate to reality requires git archaeology.
**Why RED now:** Nothing in the pipeline persists the estimate/actual pair at completion.
**GREEN when:** The next archived estimated REQ appends its line to `do-work/calibration-log.tsv` (header created on first write), and the reference file documents the log as the recalibration input.
**Validation:** User confirmed — accepted as part of "apply it" (the calibration-log fold-in offered alongside the table re-fit).

## Full Context
See `do-work/user-requests/UR-048/input.md` for complete verbatim input.

---
*Source: UR-048 — calibration application*
