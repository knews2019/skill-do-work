---
id: REQ-133
title: Testing columns honor doneWindow in their empty copy
status: pending
created_at: 2026-08-07T14:42:57Z
user_request: UR-030
domain: frontend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/web/board.js, tools/queue-kanban/generate_test.go]
---

# Testing Columns Honor DoneWindow in Their Empty Copy

## What

Fix the queue-kanban Testing-view empty-state regression. When `filterState.doneWindow` hides Testing cards that otherwise exist, Testing columns must display “No matches” instead of “Nothing here,” without allowing that Testing-only filter to affect ordinary Board, Calendar, or By-UR empty states.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
