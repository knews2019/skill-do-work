---
id: REQ-133
title: Testing columns honor doneWindow in their empty copy
status: completed
completed_at: 2026-08-07T14:55:51Z
commit: 680be8e
claimed_at: 2026-08-07T14:47:08Z
created_at: 2026-08-07T14:42:57Z
user_request: UR-030
domain: frontend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/web/board.js, tools/queue-kanban/generate_test.go]
route: A
---

# Testing Columns Honor DoneWindow in Their Empty Copy

## What

Fix the queue-kanban Testing-view empty-state regression. When `filterState.doneWindow` hides Testing cards that otherwise exist, Testing columns must display “No matches” instead of “Nothing here,” without allowing that Testing-only filter to affect ordinary Board, Calendar, or By-UR empty states.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Extend the Node-backed Go regression to execute `columnEmptyText()` for ordinary Board copy before and during Testing view, plus the explicit Testing filter decision; prove RED before adding the optional boolean override and passing `hasActiveVisibleFilters()` only from `fillTestingColumn()`.
- [x] **[APPLY]:** Added the three-outcome caller regression and the optional explicit filter override. Root cause: `fillTestingColumn()` used the ordinary no-argument empty-copy path, so its Testing-only `doneWindow` decision was never supplied; ordinary callers remain no-argument and Board-scoped.
- [x] **[UNIFY]:** Reviewed `tools/queue-kanban/web/board.js` for explicit Testing-only filter propagation and preserved no-argument ordinary callers; reviewed `tools/queue-kanban/generate_test.go` for all three runnable copy assertions and production-faithful caller seams; reviewed this REQ for P-A-U/root-cause notes. `gofmt`, focused/full Go tests, `go vet`, `go build`, `git diff --check`, and `_dev/tests/contract-regressions.sh` pass; no debug artifacts or unrelated project-file changes found.

## Context

The current helper split is partially correct: `hasActiveFilters()` covers search, domain, and status, while `hasActiveVisibleFilters()` adds `doneWindow` only when `viewState.view === "testing"`. The regression remains because `fillTestingColumn()` calls the same no-argument `columnEmptyText()` path as ordinary columns, so it never supplies the Testing-only filter decision.

Hidden Board columns can still be rendered while the Testing view is active. Making the shared empty-copy helper read `hasActiveVisibleFilters()` by default would therefore leak the Testing-only filter into those hidden columns and could leave stale “No matches” copy after a view switch.

## Detailed Requirements

- Update `tools/queue-kanban/web/board.js` so Testing columns treat `filterState.doneWindow` as an active filter when choosing their empty-state copy.
- Keep ordinary Board-column empty copy dependent only on search, domain, and status filters through the existing ordinary-filter path.
- Do not make the shared `columnEmptyText()` call `hasActiveVisibleFilters()` globally.
- Pass the relevant filter decision explicitly through the Testing-column empty-state path.
- Give `columnEmptyText()` an optional explicit boolean parameter: use that boolean when supplied, otherwise fall back to `hasActiveFilters()`.
- Ordinary columns must continue calling `columnEmptyText()` without an argument.
- Testing columns must call `columnEmptyText(hasActiveVisibleFilters())`.
- Extend `tools/queue-kanban/generate_test.go` with a regression test that executes the empty-copy decision and proves all three outcomes:
  - with only `doneWindow` set and `viewState.view === "board"`, ordinary Board empty copy is “Nothing here”;
  - after changing `viewState.view` to `"testing"`, ordinary Board empty copy is still “Nothing here,” covering hidden Board columns rendered while Testing is active;
  - the explicit Testing-column path returns “No matches.”
- Format the Go test with `gofmt`.
- Run `go test ./...` in `tools/queue-kanban` and any repository lint/format checks required by this repository.

## Constraints

- Preserve the Testing-only scope of `doneWindow`; it must not affect Board, Calendar, or By-UR empty states.
- Follow the existing ES5-style JavaScript conventions in `board.js`.
- Keep the change surgical to the queue-kanban regression and its proof.
- Do not include the merged alignment-note fix; that belongs to the sentence-aligner application, not this upstream repository.

## Builder Guidance

The implementation direction is firm. Use the optional-boolean shape validated by the user:

```js
function columnEmptyText(filtersActive) {
  var resolvedFiltersActive =
    typeof filtersActive === "boolean"
      ? filtersActive
      : hasActiveFilters();
  return resolvedFiltersActive ? "No matches" : "Nothing here";
}
```

The regression test must assert the resulting copy, not only the helper booleans; the current test already proves `hasActiveFilters()` and `hasActiveVisibleFilters()` diverge correctly but still passes while the UI is wrong.

## Red-Green Proof

**RED prompt/case:** Execute the inlined `columnEmptyText()` behavior with `filterState` containing only `doneWindow: "168"`. Ordinary Board copy is “Nothing here” in Board view and must remain so after `viewState.view` changes to `"testing"`; the explicit Testing-column decision should be “No matches,” but the current no-argument Testing call also returns “Nothing here.”
**Why RED now:** `fillTestingColumn()` does not pass `hasActiveVisibleFilters()` to `columnEmptyText()`, and the existing Go regression checks only filter-helper booleans rather than the actual empty-copy call paths.
**GREEN when:** The runnable regression proves ordinary empty copy is “Nothing here” both before and during Testing view, while the explicit Testing empty-copy path returns “No matches,” and the full queue-kanban and repository checks pass.
**Validation:** User confirmed by requesting capture and execution immediately after accepting the validated feedback and implementation shape.

## Dependencies

None.

## Full Context

See `do-work/user-requests/UR-030/input.md` for the current verbatim command and the complete validated upstream prompt.

---
*Source: UR-030 - "Fix the queue-kanban Testing-view empty-state regression using the validated requirements above"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request identifies the exact production and test files, gives the required helper shape and caller behavior, and defines runnable regression outcomes.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

`fillTestingColumn()` reused the ordinary no-argument `columnEmptyText()` path, whose fallback deliberately checks only search, domain, and status. The Testing-only `doneWindow` decision therefore never reached the empty-copy helper; changing the fallback globally would have leaked that filter into hidden Board columns.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/web/board.js` (modified)
- `tools/queue-kanban/generate_test.go` (modified)

**What was done:** Added an optional explicit filter-state argument to the shared empty-copy helper, passed Testing-visible filter state only from the Testing-column caller, and replaced the helper-only regression with runnable Board, hidden-Board, and Testing empty-copy assertions.

## Qualification

Passed — 2 project files verified, all 8 detailed requirements traced, P-A-U confirmed. The production diff is the requested optional-boolean helper plus one explicit Testing caller; ordinary callers remain no-argument, the regression runs the real caller seams, and no debug artifacts or unrelated changes were found.

## Testing

**Tests run:** `gofmt -w generate_test.go`; `go test ./... -run '^TestTestingDoneWindowIsViewSpecific$' -count=1`; `go test ./...`; `go vet ./...`; `go build ./...` in `tools/queue-kanban`; repository checks `_dev/tests/contract-regressions.sh`, `_dev/tests/record-commit-hash-guards.sh`, `_dev/tests/update-script-behavior.sh`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- `TestTestingDoneWindowIsViewSpecific` in `tools/queue-kanban/generate_test.go`: ✗ before implementation (`Testing empty copy with doneWindow = "Nothing here", want No matches`) → ✓ after implementation

**Existing tests updated (cross-REQ impact):**
- `tools/queue-kanban/generate_test.go`: strengthened the existing view-specific filter regression from helper booleans to production caller-seam copy assertions; intended Board isolation is preserved while the previously untested Testing empty-copy behavior is now covered.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-07T14:55:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — all three caller-seam copy outcomes pass and the Testing-only filter remains isolated from ordinary columns.
**Suggested testing:** 1 optional browser smoke test
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Testing columns now distinguish date-filtered emptiness from a genuinely empty queue without changing Board, Calendar, or By-UR copy; this lives in the queue-kanban frontend rendering subsystem.
