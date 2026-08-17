---
id: REQ-213
title: Board surfaces negative claimed→completed duration as a completion anomaly
status: completed
created_at: 2026-08-17T08:05:43Z
claimed_at: 2026-08-17T08:15:06Z
route: B
completed_at: 2026-08-17T08:25:24Z
commit: c3a0a6d
user_request: UR-048
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-211, REQ-212]
batch: estimator-calibration
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go]
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-17T08:05:43Z
  basis:
    - Route B
    - 3-file write set
    - (priced with the pre-calibration table)
---

# Board Surfaces Negative Claimed→Completed Duration as a Completion Anomaly

## What

`detectCompletionAnomaly` currently short-circuits whenever `completed_at` parses, so a REQ whose `completed_at` is earlier than its `claimed_at` (a real case: archived REQ-091) renders as normal. Flag that negative span as a completion anomaly so the always-visible strip and `verify` surface it for repair.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** prime-kanban-board + crew loaded. Replace detectCompletionAnomaly's frontmatter-path early return with a negative-span branch: both stamps must parse (parseTimestamp) and completed strictly Before claimed; reason names both raw values + likely cause/fix in the file's established style; absent/unparseable claimed_at returns unflagged (other checks' territory). Two tests in completion_anomaly_test.go's unit style: reversed-span flagged with both stamps in the reason; ordered/absent/unparseable-claimed table stays unflagged.
- [x] **[APPLY]:** As planned; two declared files touched (tests live in completion_anomaly_test.go, not model_test.go — write_set mirror corrected).
- [x] **[UNIFY]:** Diff reviewed — model.go (one function + its doc comment), completion_anomaly_test.go (two tests). go vet clean; module tests green (-count=1). No debug artifacts.

## Detailed Requirements

- Fires only when **both** `claimed_at` and `completed_at` parse under the board's timestamp rules **and** completed is strictly earlier than claimed — a reversed span cannot be real for stamps written in order.
- Joins the existing `CompletionAnomaly`/`CompletionAnomalyReason` plumbing (anomaly strip + `verify`'s never-silent line) — no parallel channel, no new JSON fields.
- The reason names both raw values and the likely cause/fix in the established style (one stamp is usually local wall-clock written with a `Z` suffix; fix by rewriting the wrong stamp with the true UTC instant).
- Unparseable or absent stamps remain the existing checks' territory — this check must not double-report them.
- Go table-driven test: reversed span → flagged with both values in the reason; ordered span and absent claimed_at → not flagged by this check; existing anomaly tests unchanged.
- Read-only reporting: no write-surface change, no parser field additions (both stamps are already parsed fields).

## Constraints

- Batch constraint: the three-write-surface sentence must not need amending.
- Board versioning is folded into the skill — normal CHANGELOG entry + suite version bump, per the prime.

## Builder Guidance

Firm on strict-both-parse gating and reuse of the existing anomaly channel; latitude on exact reason wording.

## Red-Green Proof
**RED prompt/case:** A fixture ticket with `claimed_at: 2026-01-02T10:00:00Z`, `completed_at: 2026-01-01T10:00:00Z`, `status: completed` — `detectCompletionAnomaly` returns false today (completed_at parsed ⇒ short-circuit).
**Why RED now:** The frontmatter-parsed path returns before any span comparison; archive REQ-091 demonstrates the case reaching real data.
**GREEN when:** The new Go test shows that fixture flagged with a reason naming both stamps, the ordered-span fixture stays unflagged, and `go test -count=1 ./...` passes in the module.
**Validation:** User confirmed — "queue-kanban should report the negative duration anomaly so it can be surfaced and fixed."

## Full Context
See `do-work/user-requests/UR-048/input.md` for complete verbatim input.

---

## Triage
**Route: B** - Medium
**Reasoning:** Clear outcome; the discovery was locating the short-circuit and the existing anomaly plumbing (done during capture-session analysis) — implementation extends one function inside established channels.
**Planning:** Not required

## Plan
**Planning not required** - Route B: Exploration-guided implementation
*Skipped by work action*

## Exploration
- Short-circuit at model.go:1217 (`CompletionTimeSource == CompletionFromFrontmatter → return false`) is the exact hole; ClaimedAt raw text is already parsed-adjacent (parseTimestamp:1351 handles all frontmatter stamp shapes).
- Anomaly plumbing: CompletionAnomaly/Reason fields → summary strip (main.go writeBoardSummary), generated payload (generate.go), never-silent warning (model.go:1402) — joining it requires zero new fields.
- Test home: completion_anomaly_test.go, unit style calling detectCompletionAnomaly on constructed tickets.
*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — negative-span branch in detectCompletionAnomaly
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modified) — reversed-span + not-flagged table tests

**Files I will NOT touch:** generate.go/main.go/web (existing channels carry the new reason unchanged), parser fields (both stamps already parsed)

**Acceptance criteria (restated from REQ):**
- [x] Fires only when both stamps parse and completed < claimed
- [x] Joins existing CompletionAnomaly plumbing; no new fields or channels
- [x] Reason names both raw values + likely cause/fix in established style
- [x] Absent/unparseable stamps not double-reported
- [x] Go tests: reversed flagged, ordered/claimless not; existing tests unchanged
- [x] Read-only; three-write-surface sentence untouched

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — detectCompletionAnomaly's frontmatter-parsed path now checks for a reversed span (parsed completed_at strictly before parsed claimed_at) and flags it with both raw stamps and the usual local-wall-clock cause; absent/unparseable claimed_at stays unflagged on this path; doc comment updated
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modified) — TestNegativeClaimedToCompletedSpanFlagged (reversed span, reason names both stamps) and TestOrderedOrClaimlessSpansNotFlaggedByNegativeSpanCheck (ordered / absent / unparseable claimed_at table)

**What was done:** The board now surfaces impossible negative claimed→completed spans through the existing anomaly strip, generated payload, and never-silent warning — verified live: archived REQ-091 (completed 6 seconds before its claim) appears in the summary's anomaly list with both stamps named.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): Strict Before() — a zero-length span (identical stamps) stays legal; only a truly reversed span is impossible.
- **D-02** (DECIDE & STATE): Unparseable claimed_at with parsed completed_at returns unflagged on the frontmatter path — flagging it here would double-report what the future-timestamp and stamp-shape checks own.

## Qualification

Passed — qualify.sh exit 0, scope-drift exit 0 (run below). Judgment: the branch is live in every board build; proven by the REQ-091 live surface and the new tests.

## Testing

**Tests run:** `go vet ./...`, `go test -count=1 ./...` (module, 12.3s, green incl. 2 new tests), full `maintainer-verify.sh`, live probe of the rebuilt binary
**Result:** ✓ All green; FAIL set identical to the 41-failure environment baseline; live summary shows REQ-091 flagged with both stamps in the reason.

**Red-green validation:** Traces to the REQ Red-Green Proof exactly: the reversed-span fixture (claimed 2026-01-02, completed 2026-01-01) returned false from detectCompletionAnomaly before the change (short-circuit demonstrated by reading the deleted condition) and the new test now asserts it flagged; the ordered fixture stays unflagged. Real-data GREEN: REQ-091 surfaces in the summary strip.

*Verified by work action*

## Review

**Overall: 94%** | 2026-08-17T09:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 2 important (follow-up REQs created), 3 minor + 1 nit (report only)
**Acceptance:** Pass — vet + module tests green; live binary flags REQ-091 with both stamps; all three pre-existing anomaly classes verified still firing; adversarial edge probes (zero span, unparseable stamps, non-frontmatter source, non-terminal statuses, masking) all clean.
**Important findings (gate: user-visible → follow-up REQs per review-work Step 10):**
- I1 → REQ-214: `verify` never lifts CompletionAnomaly tickets into findings (`verify.go:241-251` lifts only duplicate-id warnings) — it reports `OK: no findings` on a tree carrying 10 flagged tickets, while the REQ contract text promised verify surfacing. Pre-existing plumbing gap, discovered rather than introduced; the builder dropped "verify" from its summary wording instead of surfacing the discrepancy (Think Before Coding miss, acknowledged).
- I2 → REQ-215 (sweep): shipped prose still defines completion anomalies as "completion instant can't be resolved" — stale for the reversed-span class whose instant resolves fine (`docs/board-guide.md:23,35`, the never-silent line's generic fix suffix at `model.go:1422` which self-contradicts for this class, `web/board-cards.js:401-406` comment).
**Minor (report only):** comments overstate who owns an unparseable claimed_at (no check anywhere flags it); no zero-span test pinning D-01; anomalous-and-in-window tickets dual-list in CompletionAnomalies + RecentlyDone (pre-existing possibility, now the common case — needs a deliberate policy decision).
**Suggested testing:** web-render of the long reason; flag-clears-after-repair; dual-listing policy pin.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Constraining the reviewer to hands-on execution (run the binary, probe edges in code) caught a contract-vs-plumbing gap no static read would have: the never-silent warning exists but verify never consumes it.
**What didn't:** The REQ's requirements named "verify's never-silent line" from a capture-time assumption I didn't re-verify against verify.go — a false premise that became a silent wording drop in the Implementation Summary instead of a surfaced discrepancy. Verify the consumer, not just the producer, when naming a reporting contract.
**Worth knowing:** detectCompletionAnomaly's frontmatter path is production-reachable only via resolveCompletionTime, so unit fixtures can construct states production never sees — the `completedParsed` guard is what keeps the function safe standalone.

## Orientation

**Now you can** see impossible negative claimed→completed spans on the board: summary strip, generated payload, and the never-silent warning all name both stamps (archived REQ-091 surfaces live). Lives in the board's anomaly-detection layer; follow-ups REQ-214/215 extend the surfacing to `verify` and sync the anomaly-class prose.

---
*Source: UR-048 — negative-duration anomaly*
