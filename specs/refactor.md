# Spec: Refactor

> Specification template for refactoring tasks — restructuring code without changing behavior.

## Output Structure

- **Changed files** — modified source files with clearer structure, better naming, or improved patterns
- **Migration notes** — if the refactor changes public API surfaces, document what callers need to update (if applicable)
- **Before/after comparison** — summarize the structural change (not a full diff, just the conceptual shift)
- **Test updates** — existing tests updated to reflect new structure, no new behavior tested

## Quality Standards

- No behavior changes unless explicitly intentional and documented
- All existing tests still pass after refactoring
- Measurable improvement — fewer lines, better performance, clearer API, reduced duplication, or improved readability
- No new dependencies introduced unless the refactor specifically requires them
- Call sites updated — no broken imports, no dangling references
- Clean imports — removed unused imports, no circular dependencies introduced

## Implementation Checklist

1. Identify scope — which files and patterns are being refactored, and what the target structure looks like
2. Verify test coverage — ensure tests exist for the code being refactored. If coverage is insufficient, add tests *before* refactoring (not after — post-refactor tests can't catch regressions from the refactor itself)
3. Make changes incrementally — one structural change at a time, not a big-bang rewrite
4. Run tests after each step — catch regressions early, not after 10 files have changed
5. Update call sites — fix all imports, references, and usages of renamed/moved code
6. Clean up imports — remove unused imports introduced by the restructure
7. Verify no regressions — full test run, lint, type check

## Evolution Path

- **Simple**: Rename/move — rename files, functions, or variables for clarity; move files to better locations
- **Medium**: Extract/inline patterns — extract shared logic into utilities, inline over-abstracted helpers, split large files
- **Complex**: Architectural restructure — change module boundaries, restructure state management, migrate patterns (e.g., callbacks to async/await, class components to hooks)

## Common Pitfalls

- Changing behavior during refactor — the definition of refactoring is structure change without behavior change. If you need both, do them in separate commits.
- Insufficient test coverage before starting — if tests don't exist for the code being refactored, you can't verify the refactor didn't break anything
- Scope creep into adjacent code — "while I'm here, I'll also fix this" leads to large, hard-to-review diffs
- Big-bang rewrites — changing everything at once makes it impossible to bisect regressions. Incremental changes are reviewable and revertable.
- Forgetting call sites — renaming a function but missing callers in test files, scripts, or documentation
