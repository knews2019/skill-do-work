---
id: REQ-174
title: Validate root Markdown fence info
status: completed
completed_at: 2026-08-11T20:23:00Z
commit: bd5ecf6
kb_status: pending
kb_entry: Markdown fence classification must align marker, info-string, and paragraph state
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-172, REQ-173]
batch: accepted-p2-fixes
write_set: [_dev/tests/shipped-package-reference-contract.sh]
claimed_at: 2026-08-11T19:56:32Z
route: A
---

# Validate Root Markdown Fence Info

## What

Make root and list fence classification share the CommonMark rule that a backtick-fence info string cannot itself contain a backtick.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add the Goldmark-differential fixture first, then centralize the CommonMark backtick-info rule across root/list opening and paragraph-state classification while preserving tilde behavior.
- [x] **[APPLY]:** Implemented the shared marker-aware predicate and the root/list/container/tilde fixtures within the one declared test file.
- [x] **[UNIFY]:** Reviewed the complete diff and changed file; ran Bash syntax, ShellCheck, the shipped-reference contract, `git diff --check`, same-primitive search, and debug/untracked checks.

## Why

The current root fence branch accepts an invalid opener and masks a following link that the repository's pinned Goldmark renderer publishes.

## Detailed Requirements

- Apply one opener/info validity rule to root and list fence branches.
- Reject backticks in backtick-fence info strings.
- Preserve tilde-fence behavior.
- Add the reproduced root-level Goldmark-differential fixture.

## Constraints

- Consolidate the existing list validation rather than adding another independent exception.
- Preserve the classifier behavior earned by REQ-150 and REQ-163.

## Red-Green Proof

**RED prompt/case:** Classify `````lang`invalid\n[live](visible.md)\n````` and compare it with the pinned Goldmark renderer.
**Why RED now:** The shell classifier masks every line and returns no target while Goldmark renders `visible.md` as a link.
**GREEN when:** The classifier returns `visible.md`, agrees with Goldmark, and existing root/list/tilde fence fixtures still pass.
**Validation:** User confirmed

## Full Context

See `do-work/user-requests/UR-039/input.md` and the preceding validated-feedback report.

---
*Source: fix accepted*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the affected test surface, supplies an exact reproduction, and requires a focused classifier fix plus a regression fixture.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — DECIDE & STATE:** Use one marker-group-aware predicate for root and list matches. Both regexes expose the marker differently, but the CommonMark rule is identical.
- **D-02 — DECIDE & STATE:** Apply the predicate to paragraph-state detection too. Rejecting only fence state still misclassifies indented lazy continuations.
- **D-03 — DECIDE & STATE:** Preserve tilde behavior with an explicit fixture because CommonMark restricts backticks only in backtick-fence info strings.

## Root Cause

Root fence candidates were accepted solely by the marker regex, while the CommonMark backtick-info prohibition existed only as list-specific suffix logic. The root path entered fence state and masked later links; paragraph-state classification also treated the invalid lookalike as a real block boundary.

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified)

**What was done:** Centralized marker-aware fence info-string validation across root and list openings plus paragraph/container state, and added root/list/tilde differential fixtures that preserve the existing classifier contracts.

## Discovered Tasks

- **[normal]** Assess `skills/do-work-board/tools/queue-kanban/render.go` question-option fence preprocessing against invalid backtick-fence info strings; it currently toggles on any trimmed triple-backtick prefix and sits outside this REQ's one-file scope.

## Qualification

Passed — 1 implementation file verified, all 4 Detailed Requirements traced, changes are substantive and scoped, and P-A-U evidence is complete. The merged diff range `73e30d0..bd5ecf6` contains only the declared classifier/test file.

## Testing

**Tests run:** `bash _dev/tests/shipped-package-reference-contract.sh`; `bash -n _dev/tests/shipped-package-reference-contract.sh`; `shellcheck _dev/tests/shipped-package-reference-contract.sh`; `git diff --check 73e30d0..bd5ecf6`
**Result:** ✓ All passing on the merged tree after temporarily isolating and restoring the pre-existing REQ-173 prime-link baseline failure.

**Red-green validation:**
- Root backtick info-string Goldmark differential: ✗ before implementation (`expected ['visible.md'], got []`) → ✓ after (`shipped package reference contract: PASS`).
- Invalid root/list paragraph-state continuations: ✗ before the paragraph-state correction → ✓ after with the existing list/container suite intact.

**New tests added:**
- Root invalid backtick-info differential, root/list continuation state, and tilde-info preservation fixtures in `_dev/tests/shipped-package-reference-contract.sh`.

*Verified by work action*

## Review

**Verdict:** Approve — 100%, Acceptance Pass.

- Exact reviewed range: `73e30d0..bd5ecf6`.
- All REQ and UR-039 requirements are implemented in the declared one-file scope; TDD evidence reproduces the root failure and passes the same contract after the fix.
- Syntax, ShellCheck, shipped-reference contract, whitespace, and adjacent Go tests pass from the clean builder worktree.
- One Important adjacent issue remains outside this REQ: the board renderer uses a prefix-only fence heuristic. It is already recorded under Discovered Tasks and must become one follow-up rather than block this fix.
- Durable report: `do-work/runs/work-2026-08-11-225637/REQ-174-review.md`.

## Lessons Learned

- Markdown fence recognition is a compound contract: marker kind, info-string validity, and paragraph/container state must change together or a locally correct opener check can still mask rendered content.
- Differential fixtures against the pinned renderer catch classifier drift more reliably than asserting regex structure.

## Orientation

- Start with the Goldmark-backed fixtures in `_dev/tests/shipped-package-reference-contract.sh`; they define the shipped classifier behavior.
- Keep the board renderer follow-up separate because its hard-break preprocessing is a distinct consumer and implementation surface.

## Knowledge Handoff

- `kb_status: pending`
- `kb_entry: Markdown fence classification must align marker, info-string, and paragraph state`
