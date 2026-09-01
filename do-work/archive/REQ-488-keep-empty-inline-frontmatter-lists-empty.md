---
id: REQ-488
title: '[impact-critical] Keep empty inline frontmatter lists empty in request reads'
status: completed
created_at: 2026-09-01T19:45:16Z
user_request: UR-083
addendum_to: REQ-440
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
review_generated: false
claimed_at: 2026-09-01T21:05:56Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
calculated_at: 2026-09-01T21:05:56Z
completed_at: 2026-09-01T21:17:41Z
---

# Keep Empty Inline Frontmatter Lists Empty in Request Reads

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Confirmed the loss happens only when `FieldValue` copies an empty list into a nil destination. Add model and command-level selector regressions first, then preserve the empty allocation in that copy without changing scalar or populated-list decoding.
- [x] **[APPLY]:** Preserved nil for absent/scalar evidence while copying a present empty list as a non-nil empty slice. Added model and command-level regression coverage only.
- [x] **[UNIFY]:** Reviewed the three CLI changes and their diff; verified no debug artifacts, no unrelated source changes, focused tests, full module tests, and the canonical maintainer gate all pass.

## What

A REQ whose frontmatter carries an empty inline list (`depends_on: []`) is read by `do-work-cli` as having one dependency named `[]`. The canonical `next` selector then excludes it with `DEPENDENCY-MISSING: missing dependencies: []`. On 2026-09-01 this excluded 21 of 30 pending REQs from `do-work run`, so the queue only advanced through REQs that happened to omit the field or name a real dependency.

Make an empty inline list read as an empty list everywhere `RequestDocument.FieldValue` is consumed, so `depends_on: []`, `related: []`, and `write_set: []` mean "none", never a single literal item.

## Context

Discovered while resuming REQ-440: `skills/do-work/tools/do-work-cli.sh --format json next` reported `DEPENDENCY-MISSING` with the literal `[]` for every REQ carrying `depends_on: []`. Root cause traced to `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`: `decodeFlowList("[]")` correctly returns an empty non-nil slice, but `FieldValue` copies it with `append([]string(nil), evidence.ListValues...)`, which yields `nil` for an empty input. `listValue` then treats `nil` as "no list" and falls back to the scalar text `[]`.

## Requirements

- `FieldValue` preserves the empty-but-present distinction for `ListValues` (an empty inline list stays a zero-length non-nil slice, or `listValue` checks field presence rather than nil-ness).
- `RequestRecord.DependsOn` is empty for `depends_on: []`; `next` selects such a REQ as dependency-ready.
- No change to how an absent field, a scalar value, or a populated list is read.

## Red-Green Proof
**RED prompt/case:** A request-model test that parses frontmatter containing `depends_on: []` and asserts `DependsOn` is empty; and a `next` selection test where a pending REQ with `depends_on: []` must be selected rather than excluded as `DEPENDENCY-MISSING`.
**Why RED now:** `FieldValue` returns `ListValues == nil` for the empty list, so `listValue` falls back to the scalar `[]` and the graph reports a missing target named `[]`.
**GREEN when:** Both tests pass and `do-work-cli --format json next` against a queue of `depends_on: []` REQs lists them as selected or `FAN-OUT-LIMIT`, never `DEPENDENCY-MISSING`.
**Validation:** Discovered task from REQ-440; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Triage

- **Route:** A — direct bug fix.
- **Why:** The root cause, affected method, expected behavior, and regression cases are all explicit; the smallest fix is localized to request-model evidence copying.

## Plan

Route A: planning skipped; the captured red-green proof defines the implementation and verification sequence.

## Implementation Summary

**What was done:** Request-model reads now retain the distinction between an absent/scalar list and a present empty flow list, so queue selection receives no literal square-bracket dependency.

- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified) — preserve a non-nil empty list when returning copied field evidence.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified) — cover absent, scalar, populated, and empty flow-list projections.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified) — prove the selector chooses a REQ with an empty dependency list.

## Testing

- **Red-green validation:** `TestTypedRecordPreservesEmptyInlineListEvidence` failed before the fix because `DependsOn` was `[]string{"[]"}` for `depends_on: []`; it passes after the fix. `TestNextCommandTreatsEmptyInlineDependencyListAsReady` passes with the REQ selected rather than excluded as `DEPENDENCY-MISSING`.
- **Focused tests:** `go test ./internal/requestmodel -run TestTypedRecordPreservesEmptyInlineListEvidence -count=1`; `go test ./internal/nextselection -run TestNextCommandTreatsEmptyInlineDependencyListAsReady -count=1`.
- **Module regression:** `go test ./...`.
- **Repository gate:** `bash _dev/tests/maintainer-verify.sh` passed; its optional browser lane was skipped because no browser is configured.

## Qualification

Passed — 3 project files verified, 3 requirements traced, and all P-A-U checks confirmed. The canonical qualification command passed after the implementation summary declared its file entries in the required format.

## Review

**Overall: 100%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Empty, scalar, and populated list behavior are covered; the selector no longer treats an empty list as a missing dependency. |
| Code Quality | 100% | The nil-preserving guard is minimal and keeps the existing scalar fallback intact. |
| Test Adequacy | 100% | Parser/model and real command-level selector regressions satisfy the captured red-green proof. |
| Scope | 100% | Only the request-model fix and its tests changed. |
| Risk | Low | This corrects shared frontmatter evidence; focused and full-suite coverage passed. |
| Acceptance | Pass | The actual queue selector now considers empty-list REQs by their real dependency state. |

**Findings:** None. The restatement sweep found no stale CLI consumer; the board parser already documents the same empty-list distinction.

## Lessons Learned

**What worked:** A model-level assertion plus a command-level selector fixture isolated the nil-versus-empty handoff without broadening the parser.

**What didn't:** The claim transaction exposed a separate checkpoint-entry placement defect; it was recorded as a discovered task instead of being folded into this fix.

**Worth knowing:** A copied list must preserve both value contents and presence semantics because downstream scalar fallback relies on nil specifically.

## Orientation

The do-work CLI's request-model subsystem now carries empty frontmatter lists through to dependency selection as empty collections rather than text values.

## Discovered Tasks

- The canonical claim transaction placed this REQ's checkpoint entry immediately before, rather than inside, the `## In Progress (interrupted)` section. Capture this as a separate request; it is unrelated to empty-list projection.

## Builder Guidance

Certainty level: Firm. The nil-versus-empty copy in `FieldValue` is the accepted root cause; keep the fix to the request model and its tests plus one selector-level regression.

---
*Source: discovered task recorded during REQ-440 (work action Step 8).*
