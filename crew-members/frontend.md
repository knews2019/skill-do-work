# The Renderer — Frontend Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent only when working on frontend-related tasks. Keep rules scoped and concise to minimize token usage. -->

## Implementation Patterns

### Component Structure
- Follow the project's existing component conventions (functional vs class, file organization, naming).
- One component per file unless the project clearly groups related small components.
- Co-locate styles, tests, and types with their component when the project does this.
- Prefer composition over inheritance in component hierarchies.

### State Management
- Use the project's existing state solution — do not introduce a new one.
- Local state for UI-only concerns (open/closed, hover, form input).
- Shared/global state only when multiple unrelated components need the same data.
- Derive computed values rather than storing redundant state.

### Performance Baseline
- Avoid re-renders from unstable references (new objects/arrays in render, inline function definitions in JSX when they cause child re-renders).
- Lazy-load routes and heavy components when the project supports code splitting.
- Images: use appropriate formats, provide dimensions to prevent layout shift.

### Error Handling
- Every data fetch needs loading, success, and error states.
- Form validation: client-side for UX, never trust it for security.
- Display user-friendly error messages — log technical details to console.

## Quality Checks

Before marking UNIFY complete, verify:

| Criterion | What to check |
|-----------|---------------|
| Renders without errors | No console errors/warnings on mount and primary interaction |
| Responsive | Tested or verified at 320px, 768px, 1280px minimum |
| Accessible | Keyboard navigable, semantic HTML, no missing alt/labels |
| No regressions | Existing tests still pass after changes |
| Bundle impact | No unnecessary large dependencies added |

## Scope Discipline

- Do not refactor unrelated components while fixing a bug in one.
- Do not upgrade dependencies unless the REQ explicitly requests it.
- Do not switch styling approaches (e.g., CSS modules to Tailwind) as a side effect.
