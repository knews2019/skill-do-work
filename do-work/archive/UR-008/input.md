---
id: UR-008
title: Fix confirmed findings from the deep review of REQ-035–040 and adjudicated external feedback
created_at: 2026-07-29T09:30:45Z
requests: [REQ-043, REQ-044, REQ-045, REQ-046, REQ-047, REQ-048, REQ-049, REQ-050]
word_count: 2
---

# Fix confirmed findings from the deep review of REQ-035–040

## Summary

After the REQ-035–040 batch shipped (commits fd56267..acc4722, versions 0.144.0–0.146.3), the user asked for an independent review of the work. Three parallel reviewers re-examined the six commits (with empirical verification: executing qualify.sh against bad ranges, scratch-repo merge-range checks, Go 1.26 path.Match probes, repo-wide token greps), and the user separately submitted three rounds of external concurrency-review feedback that the agent adjudicated item by item. The confirmed findings were grouped into eight fix REQs, which the user approved capturing with "capture them". One additional confirmed finding (standalone `git show <commit>` consumers beyond review-work.md) was appended as an addendum to REQ-055 rather than duplicated here — the capture proposal targeted REQ-042, but a concurrent `do-work clarify` session resolved REQ-041/042 into REQ-051–056 mid-capture (briefly colliding with this batch's ids 043–048 before that session renumbered its files to 051–056), and REQ-055 is the successor to the targeted item.

## Extracted Requests

| REQ | Title | Origin |
|---|---|---|
| REQ-043 | Worktree evidence-range hardening (REQ-037 follow-up) | Deep review of 1348d11 + external feedback convergence on the seam-range defect |
| REQ-044 | Lock claim coherence fixes (REQ-035 follow-up) | Deep review of fd56267 |
| REQ-045 | Close the Route A gap in dispatch re-validation (REQ-036 follow-up) | Deep review of 4296e11 |
| REQ-046 | Record the operative worktree name after a collision variant (REQ-038 follow-up) | Deep review of efb6300 |
| REQ-047 | Mutex pre-write ownership check | External feedback (live-owner eviction at the 60s reclaim bound), verified against current code |
| REQ-048 | UR-closure keying consistency between cleanup Pass 1 and work Step 8 | CHECKPOINT.md note from the prior session, promoted to a REQ |
| REQ-049 | Add a restatement-sweep step to the adversarial review instructions | Calibration lesson: every top finding was an un-updated cross-file restatement the 86–98% review passes missed |
| REQ-050 | Doc-accuracy fixes from the review (3 small corrections) | Reviewer nits + adjudicated external feedback on the session-start hook comment |

## Batch Constraints

- These are fixes to the skill's own operating instructions and shipped tooling — REQ-043/044/045 carry `maintenance: true` because correcting/narrowing drifting instruction text is a candidate fix; the rest are additive.
- REQ-043/044/045/046/049 all write `actions/work.md` and/or `actions/work-reference.md` — overlapping write sets, expected to serialize.
- Adjudicated and deliberately NOT captured: the capture-dedup TOCTOU (benign duplicate block since d839cf5; lock fix deliberately declined), the 45-minute synchronous-wait takeover window (documented accepted residual risk; a "quiet lease" design exists if ever wanted), and the qualify.sh false-FAIL-on-own-tokens + blocked-flip-main-tree items (already tracked in REQ-042).
- Also done under this UR: an Addendum appended to `do-work/queue/REQ-055-standalone-review-merge-commit-diff.md` (successor of REQ-042's first open question) enumerating the remaining `git show <commit>` consumers (present-work.md ×2, pipeline.md, ai-report.md — not just review-work.md).

## Full Verbatim Input

capture them

(Said in response to the agent's proposed capture set: seven grouped fix-REQs plus an edit to REQ-042, covering the confirmed findings from the three-reviewer deep pass over commits fd56267, 4296e11, 1348d11, efb6300, 9c8ef45, acc4722 and the adjudicated external feedback rounds. The full finding statements are preserved in each REQ's Detailed Requirements. The doc-accuracy nits, proposed as riders, were captured as standalone REQ-050 because no proposed REQ was a natural host.)

---
*Captured: 2026-07-29T09:30:45Z*
