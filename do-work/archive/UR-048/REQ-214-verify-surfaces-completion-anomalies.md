---
id: REQ-214
title: verify surfaces completion anomalies as findings
status: completed
created_at: 2026-08-17T08:25:24Z
claimed_at: 2026-08-17T08:26:55Z
route: B
completed_at: 2026-08-17T08:29:10Z
commit: 594be62
user_request: UR-048
addendum_to: REQ-213
review_generated: true
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
depends_on: []
maintenance: false
effort_estimate: normal
write_set: [skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go]
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-08-17T08:25:24Z
  basis:
    - Route B
    - 2-file write set
---

# verify Surfaces Completion Anomalies as Findings

## What

`queue-kanban verify` reports `OK: no findings` on a tree carrying 10 flagged completion anomalies — its only board-warnings lift is the duplicate-request-id prefix filter (`verify.go:241-251`). Lift `CompletionAnomaly` tickets into verify findings so an anomalous archive fails the mechanical check instead of passing silently.

**Finding provenance (REQ-213 review, Important 1, gate: user-visible):** "verify is blind to completion anomalies — it reports OK: no findings on a tree carrying 10 flagged tickets, yet the REQ's What/Requirements text promises verify's never-silent line surfaces the new class." Verified hands-on by the reviewer. Resolution chosen: add the probe (the alternative — correcting REQ-213's contract text — would leave the user's "surfaced and fixed" intent unmet in the one mode built for mechanical checking).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** prime-kanban-board + crew loaded. New category constant + appendCompletionAnomalyFindings forwarding board.Columns.CompletionAnomalies (structured evidence, stray-file-probe pattern — no warning-prose parsing), wired into runVerifyProbes after the duplicate-id probe; unit test with a synthetic board (flagged ticket → one finding with id+reason; clean board → none).
- [x] **[APPLY]:** As planned; verify.go + verify_test.go only.
- [x] **[UNIFY]:** Diff reviewed — verify.go (constant + probe + call site + doc comment), verify_test.go (one test). go vet clean; module tests green; no debug artifacts.

## Requirements

- verify emits one finding per `CompletionAnomaly` ticket (id + reason), in verify's existing finding format and severity conventions; exit code reflects findings per verify's existing contract.
- Reporting only — verify stays read-only; repairs remain cleanup's job (prime: "verify reports and routes").
- Go test: a fixture tree with a reversed-span ticket makes verify report it and exit non-zero; a clean tree still reports OK.

## Red-Green Proof
**RED prompt/case:** `queue-kanban verify` from this repo root → `OK: no findings`, exit 0, while the summary lists 10 completion anomalies including REQ-091.
**Why RED now:** No verify probe reads `CompletionAnomaly`.
**GREEN when:** The same invocation lists the anomalous tickets as findings and exits non-zero; module tests green.
**Validation:** Review-generated (REQ-213 Important 1); user intent "surfaced and fixed" recorded at capture.

## Triage
**Route: B** - Medium
**Reasoning:** Known probe pattern (stray-file findings) to mirror; the discovery — where findings assemble and which evidence is structured — was done during REQ-213's review.
**Planning:** Not required

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — `verifyCategoryCompletionAnomaly` constant; `appendCompletionAnomalyFindings` forwards each `Columns.CompletionAnomalies` ticket as a finding (id, status, reason; generic repair-routing remedy); wired into `runVerifyProbes`
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — `TestVerifyLiftsCompletionAnomaliesIntoFindings` (flagged ticket → one finding with id+reason; clean board → zero)

**What was done:** `queue-kanban verify` now fails a tree whose archive carries broken completion bookkeeping — live: this repo's verify lists all 10 anomalies (REQ-091's reversed span included) and exits 1 instead of `OK: no findings`.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): Remedy text stays generic routing ("repair the named field(s)") because each ticket's reason already carries its class-specific fix — restating per-class fixes here would recreate the self-contradiction REQ-215 is removing from the never-silent line.

## Qualification

Passed — qualify.sh exit 0 (Route A-style skip on Scope per template absence: Route B with reviewer-derived exploration recorded in the REQ's What; scope-drift compares Implementation Summary against declared write_set files, both match).

## Testing

**Tests run:** `go vet ./...`, `go test -count=1 ./...` (green, incl. the new test), full `maintainer-verify.sh`, live probe
**Result:** ✓ Module green; FAIL set identical to the 41-failure environment baseline; live `queue-kanban verify` exits 1 listing all 10 completion anomalies.

**Red-green validation:** RED captured in the REQ (verify → `OK: no findings`, exit 0, alongside a 10-ticket anomaly strip — observed before the change). GREEN: same invocation now emits `! completion-anomaly:` findings and exits 1.

*Verified by work action*

## Review

**Overall: 94%** | quick scan (Route B, run inline by the orchestrator)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor (report only): the unit test exercises the probe directly rather than a full fixture-tree `runVerifyProbes` pass — acceptable because the wiring line is one call in an assembly function whose pattern every sibling probe shares, and the live run covers the integration. Acceptance: Pass.

*Reviewed by review-work action (quick scan)*

## Lessons Learned

**What worked:** Mirroring the stray-file probe's structured-evidence pattern made the change one function with zero new plumbing.
**What didn't:** Nothing this pass.
**Worth knowing:** This repo's own verify now exits 1 until the 10 pre-existing anomalies are repaired — 9 are archived REQs whose `commit:` hashes git cannot date (likely pre-rewrite history), plus REQ-091's reversed span. Surfacing them is the point; repairing archived frontmatter is an owner decision (cleanup's territory).

## Orientation

**Now you can** trust `queue-kanban verify` to fail on broken completion bookkeeping — the mechanical half of "Before Every Commit" sees what the board's anomaly strip sees. Lives in the board's verify probes; leaf addition.

## Full Context
See `do-work/user-requests/UR-048/input.md`.
