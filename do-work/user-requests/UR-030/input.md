---
id: UR-030
title: Fix the queue-kanban Testing-view empty-state regression
created_at: 2026-08-07T14:42:57Z
requests: [REQ-133]
word_count: 19
---

# Fix the Queue-Kanban Testing-View Empty-State Regression

## Summary

Fix the queue-kanban Testing view so a Testing-only recently-done filter that hides existing cards produces the filtered empty copy, while ordinary Board, Calendar, and By-UR empty states remain unaffected. The validated implementation direction passes the relevant filter state explicitly to Testing columns and adds a runnable regression test for both visible and hidden view paths.

## Extracted Requests

| ID | Request |
|---|---|
| REQ-133 | Make Testing-column empty copy honor `doneWindow` without leaking that filter into ordinary views. |

## Batch Constraints

- Treat this as one behavioral regression with a runnable RED-first test in `tools/queue-kanban/generate_test.go`.
- Keep `doneWindow` Testing-only: it must not influence Board, Calendar, or By-UR empty states, including hidden Board columns rendered while Testing is active.
- Use an explicit Testing-column filter-state argument rather than changing the shared default globally.
- Run `gofmt`, `go test ./...`, and the repository lint/format checks.
- Do not include the already-merged sentence-aligner alignment-note fix; it belongs to a different application.

## Full Verbatim Input

```arduino
do-work capture-request: Fix the queue-kanban Testing-view empty-state regression using the validated requirements above
  do-work run                 Process the captured fix
```

## Referenced Conversation Context

The capture command refers to the following validated upstream request from the immediately preceding conversation:

> Fix the queue-kanban Testing-view empty-state regression.
>
> Problem: when `filterState.doneWindow` hides existing Testing cards, the Testing columns incorrectly display “Nothing here” instead of “No matches.” However, `doneWindow` is Testing-only and must not affect Board, Calendar, or By-UR empty states—even when hidden Board columns are rendered while the Testing view is active.
>
> Requirements:
>
> 1. Update `tools/queue-kanban/web/board.js` so Testing columns treat `doneWindow` as an active filter.
> 2. Keep ordinary Board columns dependent only on search, domain, and status filters.
> 3. Do not simply make the shared `columnEmptyText()` read `hasActiveVisibleFilters()` globally; that leaks the Testing filter into hidden Board columns, which may remain stale after switching views.
> 4. Pass the relevant filter state explicitly to the Testing-column empty-state path.
> 5. Add a regression test in `tools/queue-kanban/generate_test.go` proving:
>    - Board empty copy remains “Nothing here” when only `doneWindow` is set.
>    - This remains true even while `viewState.view === "testing"`.
>    - Testing empty copy becomes “No matches.”
> 6. Run `gofmt`, `go test ./...`, and any repository lint/format checks.
>
> Reference implementation shape:
>
> ```js
> function columnEmptyText(filtersActive) {
>   var resolvedFiltersActive =
>     typeof filtersActive === "boolean"
>       ? filtersActive
>       : hasActiveFilters();
>   return resolvedFiltersActive ? "No matches" : "Nothing here";
> }
> ```
>
> Ordinary columns call `columnEmptyText()`; Testing columns call `columnEmptyText(hasActiveVisibleFilters())`.
>
> Do not include the merged alignment-note fix: that belongs to the sentence-aligner application, not the upstream do-work repository.

The feedback-validation pass accepted this finding after confirming that `hasActiveFilters()` correctly excludes `doneWindow`, `hasActiveVisibleFilters()` correctly recognizes it only while Testing is visible, `fillTestingColumn()` still calls the no-argument `columnEmptyText()`, and the existing test checks helper booleans rather than the rendered empty-copy decisions.

---
*Captured: 2026-08-07T14:42:57Z*
