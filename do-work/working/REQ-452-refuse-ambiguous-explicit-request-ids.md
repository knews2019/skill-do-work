---
id: REQ-452
title: 'Refuse ambiguous explicit request IDs'
status: completed-with-issues
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-450, REQ-451, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
claimed_at: 2026-09-02T11:18:08Z
route: A
dispatch_at: 2026-09-02T11:26:14Z
builder_handback_at: 2026-09-02T11:31:50Z
integration_at: 2026-09-02T11:31:50Z
review_at: 2026-09-02T11:46:20Z
remediation_at: 2026-09-02T11:52:48Z
re_review_at: 2026-09-02T12:03:38Z
completed_at: 2026-09-02T12:03:38Z
kb_status: pending
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
    - trivial short-circuit
  calculated_at: 2026-09-02T11:18:48Z
---

# Refuse Ambiguous Explicit Request IDs

## What

Preserve duplicate queue-record collision evidence when resolving numeric request IDs, and return an ambiguity exclusion when an explicit `REQ-NNN` token cannot identify exactly one file. Explicit targeting may bypass documented dependency, assignment, and impact gates, but never repository identity ambiguity.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this explicit-target duplicate-ID ambiguity root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the listed prime, full satellite lessons, bug-fix spec, and required crew rules. Add a real `Select` regression that reverses request discovery order and requires the same typed ambiguity exclusion with both queue paths, then consult the snapshot's normalized collision evidence before explicit numeric lookup without changing unique explicit-target overrides.
- [x] **[APPLY]:** Added the duplicate explicit-ID replay in `next_targets_test.go` and the minimum collision-evidence guard in `next_targets.go`. Root cause: `requestByNumber` overwrote an earlier normalized ID and explicit provenance then bypassed graph ambiguity; the guard now excludes the ambiguous number before either arbitrary file becomes a candidate.
- [x] **[UNIFY]:** Reviewed `git diff --stat` and every changed line in `next_targets.go`, `next_targets_test.go`, `repository_model.go`, and `repository_model_test.go`; confirmed the claim/checkpoint and REQ record edits are orchestration-owned. Focused and full do-work-cli tests, module vet, gofmt, `git diff --check`, and the canonical maintainer gate passed; the diff contains no debug artifacts.

## Finding Provenance

- **Finding #4 — P2 — source:** `internal/nextselection/next_targets.go:30`

> ````text
> [P2] Reject ambiguous IDs even for explicit targets — [prj].claude/skills/do-work/
> tools/do-work-cli/internal/nextselection/next_targets.go:30-30
> When two queued records normalize to the same REQ number, this assignment silently overwrites the first entry; because explicit
> candidates later bypass graph ambiguity checks, next REQ-NNN selects whichever duplicate was encountered last and returns its
> exact path as runnable. Explicit targeting only overrides the documented dependency, assignment, and impact gates, not an
> unresolved record identity (.claude/skills/do-work/actions/work-reference.md:394-399), so detect the duplicate here and emit an
> ambiguity exclusion.
> ````

- **Finding #8 — P1 — source:** `internal/nextselection/next_selection.go:196-202`

> ````text
> [P1] Refuse ambiguous explicitly targeted request IDs — [prj].claude/skills/do-work/tools/do-work-
> cli/internal/nextselection/next_selection.go:196-202
> When duplicate queue records claim the same REQ id, requestByNumber has already selected one arbitrary file, and explicit
> provenance skips the node.IsAmbiguous guard here. Thus next REQ-NNN can authorize processing whichever duplicate happened to
> overwrite the map entry; explicit targeting should bypass dependency gates, not repository identity ambiguity.
> ````

- **Finding #18 — P2 — source:** `internal/nextselection/next_targets.go:29-30`

> ````text
> [P2] Refuse explicit IDs that resolve to duplicate queue files — [prj].claude/skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:29-30
> When two queued files normalize to the same numeric REQ id, this assignment silently keeps the last one discovered. next REQ-NNN then selects that arbitrary file because explicit provenance bypasses the graph's IsAmbiguous guard, even though the token cannot
> distinguish the records. Preserve collision evidence and return an ambiguity exclusion instead, consistent with prime-do-work-cli.md:14-16 (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L14-L16).
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:24-31` builds a numeric map whose assignment overwrites an earlier duplicate. Explicit provenance then skips ambiguity at `next_selection.go:196-203`, despite collision evidence already being available in `internal/repositorymodel/repository_model.go:78-93,239-268`. The chosen file therefore depends on discovery order.
- **Surface-cost result:** Earned — the repository already computes collision evidence. Reusing it and adding one duplicate-explicit-target replay is cheaper than arbitrary request authorization.

## Detailed Requirements

- Preserve every path contributing to a duplicate normalized numeric request ID.
- Reject explicit targeting when the token resolves to more than one queue record.
- Emit a typed ambiguity exclusion rather than returning either duplicate as runnable.
- Keep explicit targeting's documented dependency, assignment, and impact overrides unchanged.
- Make the result independent of repository discovery order.

## Constraints

- Reuse normalized collision evidence rather than inventing a second identity rule.

## Dependencies

No request prerequisite. Shared selector files with other UR-085 requests do not establish dependency ordering.

## Builder Guidance

Certainty level: Firm. Identity ambiguity is a repository integrity condition, not an eligibility gate.

## Red-Green Proof

**RED prompt/case:** Place two queue files that normalize to the same numeric REQ ID, explicitly request that `REQ-NNN`, and repeat with reversed discovery order.
**Why RED now:** Numeric-map assignment silently keeps one path and the explicit path skips the graph's ambiguity check.
**GREEN when:** Both replays return the same typed ambiguity exclusion containing collision evidence, and neither duplicate is selected.
**Validation:** User confirmed after validate-feedback accepted Findings #4/#8/#18.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---

## Triage

**Route: A** - Simple

**Reasoning:** The request identifies the exact selector seam, the identity rule to preserve, and a deterministic duplicate-ID replay. The change is a focused bug fix with a direct regression test.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2201 tokens; matches selector collision-evidence work but exceeds the 2000-token required-lessons budget and the partially slugged satellite cannot be narrowed safely. It was still read under the touch-conditional Lessons Discipline rule.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)

**What was done:** Repository collision evidence now normalizes numeric filename and frontmatter claims to the same identity used by explicit selection. Explicit numeric REQ targeting consults that evidence before choosing a record, returns a deterministic typed ambiguity exclusion with every colliding queue path, and stays independent of discovery order even when duplicate frontmatter IDs live under unrelated filenames. The documented explicit-target overrides remain intact for unique identities.

**Root cause:** The numeric request map retained only the last record for a normalized number, while repository collision claims were keyed by raw text and explicit provenance bypassed the dependency graph's ambiguity guard. Numeric-equivalent frontmatter IDs could therefore escape collision evidence and let discovery order choose repository identity.

## Qualification

Passed after remediation — 4 files verified, 5 detailed requirements traced, and P-A-U confirmed. The source change is substantive, aligns repository collision evidence with the selector's numeric identity, preserves every contributing path, produces deterministic sorted queue-path evidence, and leaves unique explicit-target gate overrides unchanged. Mechanical qualification, `git diff --check`, focused repositorymodel/nextselection tests, the full module tests, and module vet passed.

## Testing

**Tests run:**
- `go test ./internal/nextselection -run '^TestExplicitREQRefusesAmbiguousQueueIdentityIndependentOfDiscoveryOrder$' -count=1` — RED before implementation (exit 1): forward discovery selected `REQ-452-first.md`; reversed discovery selected `REQ-0452-second.md`.
- The same focused command — GREEN after implementation (exit 0).
- `go test ./internal/nextselection -run '^TestExplicitREQRefusesNormalizedFrontmatterIdentityCollisionIndependentOfDiscoveryOrder$' -count=1` — RED during review/remediation (exit 1): unrelated filenames allowed discovery order to select different normalized frontmatter records; GREEN after remediation.
- `go test ./internal/repositorymodel -run '^TestCollisionEvidenceNormalizesNumericFrontmatterIdentity$' -count=1` — RED with no collision entry before normalization; GREEN after remediation with both paths preserved.
- `go test -count=1 ./internal/repositorymodel ./internal/nextselection` — passed after remediation.
- `go test -count=1 ./...` and `go vet ./...` — passed from the do-work-cli module after remediation.
- `bash _dev/tests/maintainer-verify.sh` — passed with exit code 0 against the final remediated implementation.

**Red-green validation:** The captured duplicate explicit-ID case and the reviewer-found frontmatter-only variant both failed at the real `Select` caller seam before their fixes and passed afterward in both discovery orders. The repository-model regression independently proves normalized collision evidence retains both paths.

**Regression evidence:** The complete nextselection package and the repository's canonical verification gate passed. No existing tests were changed outside the focused selector test file.

## Review

**Overall: 50%** | 2026-09-02T11:46:20Z

| Dimension | Score |
|-----------|-------|
| Requirements | 20% |
| Code Quality | 75% |
| Test Adequacy | 70% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Frontmatter-only numeric-equivalent IDs with different filenames do not enter `CollisionEntries`, so explicit targeting can still choose one record by discovery order — impact-user-visible → direct remediation in REQ-452.

**Minor findings:** 0 (report only)
**Acceptance:** Fail — the captured filename-collision replay passes, but an uncovered normalized-frontmatter collision still violates repository identity ambiguity.
**Suggested testing:** 1 item — add and reverse a frontmatter-only collision fixture whose filenames have different numbers.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Re-Review

**Overall: 50%** | 2026-09-02T12:03:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 75% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings:**
- Reusing suffix-tolerant filename parsing for frontmatter collision claims makes malformed `REQ-452junk` falsely collide with valid `REQ-452`, vetoing a unique explicit target — impact-user-visible → REQ-497 created.

**Minor findings:** 1 — the remediation-added repositorymodel files were absent from the builder's original UNIFY text; the orchestration record now lists the full verified set.
**Acceptance:** Fail — genuine numeric-equivalent IDs are now handled, but strict frontmatter identity remains incomplete.
**Suggested testing:** 1 item — negative repository-model and selector replays for `REQ-452` versus `REQ-452junk`.
**Follow-ups created:** REQ-497; **sweeps appended to:** None

*Reviewed by review-work action after remediation*

## Remediation

- **Initial review failure:** frontmatter-only numeric equivalents under unrelated filenames escaped collision evidence and remained discovery-order-dependent.
- **Remediation:** normalized repository collision claim keys and added repository-model plus caller-seam RED/GREEN replays; all focused/full tests, vet, and the canonical gate passed.
- **Re-review result:** the normalization was broader than strict frontmatter grammar because it reused suffix-tolerant filename parsing. The remaining Important finding is isolated in REQ-497; this REQ completes with issues after its single allowed remediation.

## Lessons Learned

**What worked:**
- Caller-seam discovery-order replays exposed both arbitrary selection and the missing normalized collision evidence.
- Keeping collision paths as repository-model evidence let the selector return one typed, actionable ambiguity exclusion.

**What didn't:**
- The first fixture used filenames that already shared the target number, masking the frontmatter-only gap.
- The remediation reused a suffix-tolerant filename parser for frontmatter, which fixed genuine numeric equivalents but admitted malformed trailing text.

**Worth knowing:** Collision fixtures must decouple filename claims from frontmatter claims. Use unrelated filename numbers when proving frontmatter identity, and include a malformed adjacent value as a negative control.

## Orientation

Explicit REQ selection now refuses genuine numeric-equivalent queue identities through the shared repository collision model instead of choosing by discovery order. The strict malformed-frontmatter boundary remains isolated in REQ-497 before this selector work is ship-ready.

---
*Source: validate-feedback Findings #4, #8, and #18, captured by UR-085.*
