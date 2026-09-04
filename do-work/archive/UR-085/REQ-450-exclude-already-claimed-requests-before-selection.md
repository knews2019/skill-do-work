---
id: REQ-450
title: 'Exclude already-claimed requests before selection'
status: completed
route: C
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-01T23:28:40Z
  basis:
    - Route C
    - 6-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
    - full-suite verification
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md]
related: [REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
claimed_at: 2026-09-01T23:28:10Z
planning_at: 2026-09-01T23:31:48Z
dispatch_at: 2026-09-01T23:35:06Z
builder_handback_at: 2026-09-01T23:41:49Z
integration_at: 2026-09-01T23:41:49Z
review_at: 2026-09-01T23:41:49Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
completed_at: 2026-09-01T23:48:27Z
release_at: 2026-09-01T23:49:50Z
commit: 88a6c4f9
---

# Exclude Already-Claimed Requests Before Selection

## What

Exclude a queued request before eligibility when its request record carries `claimed_at` or a foreign writer entry in `do-work/CHECKPOINT.md` still claims it. Return typed already-claimed evidence without inferring lease expiry or writer liveness.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this missing claim-aware selector eligibility root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Project writer-bearing checkpoint claims during the one repository walk, carry exact typed evidence into an early selector veto, preserve all mode provenance, and prove no-rescan behavior across the repository/result/selection seams.
- [x] **[APPLY]:** Implemented the snapshot projection, typed JSON/text claim evidence, early `ALREADY-CLAIMED` exclusion, command expectation, and cross-mode/no-rescan regressions in exactly the seven declared files.
- [x] **[UNIFY]:** Reviewed all seven declared files and the full diff. Verified anchored claim-header parsing, duplicate/source ordering, frontmatter-first evidence, exact text/JSON parity, explicit-target non-override, blocked-probe precedence, unrelated control eligibility, and no debug artifacts. The unrelated concurrent kanban-prime edit is excluded.

## Finding Provenance

- **Finding #1 — P1 — source:** `internal/nextselection/next_selection.go:159-160`

> ````text
> [P1] Exclude claimed requests before selecting — [prj].claude/skills/do-work/tools/do-
> work-cli/internal/nextselection/next_selection.go:159-160
> When a pending queue record has a live claimed_at, or remains named under a writer: entry in do-work/CHECKPOINT.md, this path
> proceeds directly to other filters and can select already-owned work for automatic dispatch. The selector contract requires every
> candidate to be unclaimed (.claude/skills/do-work/actions/work-reference.md:395); load checkpoint claim evidence and reject
> claimed records before selection.
> ````

- **Finding #9 — P1 — source:** `internal/nextselection/next_selection.go:154-160`

> ````text
> [P1] Exclude pending records that still carry a claim — [prj].claude/skills/do-work/tools/do-work-
> cli/internal/nextselection/next_selection.go:154-160
> A queue file with status: pending but a non-empty claimed_at passes this gate and can be selected for another build. Auto-wave's
> documented predicate explicitly requires no live claimed_at (actions/work-reference.md line 395), and the replaced simple selector
> also rejected this state, so inspect record.ClaimedAt before admitting the candidate.
> ````

- **Finding #17 — P2 — source:** `internal/nextselection/next_selection.go:154-157`

> ````text
> [P2] Exclude queued records that already carry a claim stamp — [prj].claude/skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:154-157
> For a queued record whose status is still pending but whose claimed_at is non-empty, such as after an interrupted or manual claim, this status gate accepts it because ClaimedAt is never inspected. Default, fan-out, and simple selection can consequently hand already-
> claimed work to another builder, contrary to the unclaimed selection contract in actions/work.md:35 (.claude/skills/do-work/actions/work.md#L35); emit a typed already-claimed exclusion before eligibility.
> ````

- **Evidence:** `RequestRecord.ClaimedAt` is parsed at `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go:49,270`, while `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:154-166` checks status, assignment, and impact without checking `ClaimedAt`. `RepositorySnapshot` lacks checkpoint claim evidence at `internal/repositorymodel/repository_model.go:84-97`, and discovery ignores `CHECKPOINT.md` at `repository_model.go:187-222`. The selector contract at `skills/do-work/actions/work-reference.md:391-395` requires both no `claimed_at` and no claim by another writer in the synced checkpoint; ADR-018 lines 35-37 and 49-55 records the duplicate-writer incident and governing rule.
- **Surface-cost result:** Earned — the repository documents a synced-checkpoint double-claim incident. A small typed checkpoint projection and claim exclusion cost less than duplicate dispatch; cover pending plus `claimed_at`, a queued REQ in a foreign writer entry, and unrelated checkpoint entries.

## Detailed Requirements

- Reject a pending or queued candidate carrying a non-empty `claimed_at` before any eligibility or dispatch decision.
- Load the typed checkpoint claim evidence required to reject a request named by another writer.
- Emit a typed already-claimed exclusion with enough provenance for callers to explain the result.
- Apply the exclusion consistently to default, fan-out, simple, and targeted selection wherever the unclaimed contract applies.
- Do not invent lease expiry, heartbeat, liveness, or stale-claim heuristics.
- Ignore unrelated checkpoint entries.

## Constraints

- Preserve explicit targeting's documented overrides without turning it into an override for ownership.
- Preserve the selector's no-rescan result-model contract.

## Dependencies

No request prerequisite. Shared selector files with other UR-085 requests do not establish dependency ordering.

## Builder Guidance

Certainty level: Firm. Model stored ownership evidence directly and keep policy out of the projection.

## Red-Green Proof

**RED prompt/case:** Select (1) a pending queue record with non-empty `claimed_at`, (2) a pending queue record named under another writer in `do-work/CHECKPOINT.md`, and (3) a control record with only unrelated checkpoint entries.
**Why RED now:** Selection parses the claim stamp but never applies it, and repository discovery supplies no checkpoint claim evidence.
**GREEN when:** The first two records return typed already-claimed exclusions before eligibility, while the unrelated control remains eligible.
**Validation:** User confirmed after validate-feedback accepted Findings #1/#9/#17.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Findings #1, #9, and #17, captured by UR-085.*

## Triage

**Route: C** - Complex

**Reasoning:** The fix adds checkpoint ownership evidence to the shared one-pass repository model, threads typed provenance into every selector mode, and must preserve explicit-target and no-rescan contracts with cross-mode regressions.

**Planning:** Required

## Plan

1. Extend the one-pass repository snapshot with ordered, policy-free checkpoint claim records keyed by canonical REQ id. Preserve raw claim stamp, writer, checkpoint path, source line, and header text; recognize canonical claim-shaped headers without depending on their section placement.
2. Extend selection exclusion results with deterministic claim evidence and keep JSON/text projections equivalent, including a normalized empty evidence array.
3. Add an early `ALREADY-CLAIMED` selector veto after minimum record validity but before probes, status/override policy, dependency/wave logic, or fan-out. Combine nonblank request `claimed_at` evidence first with every same-id writer-bearing checkpoint claim; explicit targeting never overrides ownership.
4. Prove repository projection, result rendering, and default/fan-out/simple/explicit-REQ/UR-expanded selection behavior, including unrelated checkpoint entries, duplicate evidence, blocked-probe precedence, and discover-once/no-rescan behavior.

**Files:** `internal/repositorymodel/repository_model.go` and tests; `internal/resultmodel/result_model.go` and tests; `internal/nextselection/next_selection.go` and tests.

**Validation warning:** The plan has six file tasks because production and regression files change in lock-step across three typed subsystems. This exceeds the five-task quality heuristic, but splitting would leave the snapshot, result envelope, or selector without its matching contract tests.

*Generated by Plan agent*

## Exploration

`repositorymodel.DiscoverRepository` owns the single walk and already provides contained regular-file reads; checkpoint parsing belongs inside that walk with a claim-header-specific matcher. `nextselection.evaluateCandidate` is the shared policy seam for every selector mode, and `resultmodel.SelectionExclusion` is the JSON/text envelope. The command integration fixture for a queued claimed record also needs its expected exclusion precedence updated from `STATUS-NOT-PENDING` to `ALREADY-CLAIMED`.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go`
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`

**Acceptance criteria:** Claimed frontmatter and writer-bearing checkpoint claims veto selection before all policy/probe stages; exclusions retain exact deterministic ownership evidence and original selection provenance; default, fan-out, simple, explicit-REQ, and UR-expanded modes agree; unrelated/malformed checkpoint entries do not interfere; selection consumes only the discovered snapshot and invents no liveness or lease policy.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/`.
**Tests baseline:** ⚠ Direct cached `go test ./...` reports pre-existing failures in `internal/corehelpers` command-identity fixtures; the repository's canonical uncached maintainer gate passed immediately before this claim. The failing package is outside this REQ's scope, so focused repositorymodel/resultmodel/nextselection tests and the canonical uncached gate remain the attribution checks.
**Dependencies:** ✓ Installed.

## Root Cause

The repository model parsed each request's `claimed_at` but did not project checkpoint ownership, and the selector began policy evaluation without consulting either ownership signal. A stale pending queue record could therefore pass default, simple, fan-out, or targeted selection and be dispatched twice; downstream code could not explain the lost ownership boundary because the result envelope had no claim-evidence type.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified)

**What was done:** The one-pass repository snapshot now retains every canonical writer-bearing checkpoint claim in source order. Selection vetoes any nonblank request claim or same-id checkpoint claim as `ALREADY-CLAIMED` before probes and all policy overrides, returning exact ordered provenance in both JSON and text while continuing to ignore unrelated checkpoint entries.

## Qualification

Passed — seven declared files verified against every detailed requirement and the P-A-U loop. Default, bounded fan-out, simple, explicit-REQ, and UR-expanded paths share the ownership veto; explicit naming never steals ownership; duplicate evidence stays ordered; snapshot mutation after discovery proves there is no second read. The unrelated concurrent `_dev/primes/prime-kanban-board.md` edit is preserved outside scope and will not be staged.

## Testing

**Tests run:** focused uncached repositorymodel/resultmodel/nextselection packages; full uncached CLI module; `go vet ./...`; `git diff --check`; canonical maintainer verification
**Result:** ✓ All passing. Focused packages, full uncached module, vet, diff checks, and canonical maintainer verification pass. The cached pre-flight-only `corehelpers` anomaly did not reproduce under `-count=1`; the optional external-browser lane was unavailable and skipped.

**Red-green validation:**
- Pending `claimed_at` and writer-bearing checkpoint claims now produce typed `ALREADY-CLAIMED` exclusions in every selector mode; unrelated checkpoint evidence leaves the control selectable.
- A claimed blocked record never runs its probe, duplicate claim evidence is retained, and deleting `CHECKPOINT.md` after discovery does not alter snapshot-derived selection.

**New tests added:**
- One-pass checkpoint parsing with duplicate, unrelated, malformed, and out-of-section records
- Typed JSON/text evidence parity and normalized empty arrays
- Cross-mode ownership veto, original provenance, probe precedence, explicit-target refusal, and no-rescan behavior

**Existing tests updated (cross-REQ impact):**
- The mixed command fixture now expects `ALREADY-CLAIMED` before `STATUS-NOT-PENDING` for a queued claimed record.

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-01T23:41:49Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — selection now fails closed on durable ownership evidence, preserves exact per-record provenance, and remains a one-snapshot read in all modes.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by work action*

## Orientation

The canonical queue selector now treats claim ownership as an early, typed eligibility boundary instead of a later status accident. The do-work CLI prime remains current.
