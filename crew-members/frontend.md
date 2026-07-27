# The Renderer — Frontend Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent only when working on frontend-related tasks. Keep rules scoped and concise to minimize token usage. -->

## Security

XSS sanitization, token storage, bundled secrets, and tabnabbing are owned by `crew-members/security.md` (see its Framework-Specific Patterns section) — don't restate them here.

## Opinions

- Animate only `transform` and `opacity` — these run on the compositor thread. Animating `width`, `height`, `top`, `left`, `margin`, or `padding` triggers layout recalculation.
- Use `will-change` sparingly, only on elements that will actually animate, and remove it once the animation completes.
- Respect `prefers-reduced-motion` — disable non-essential animations when the user requests it.
- Virtualize long lists (>100 visible items) instead of rendering everything.
- Distinguish 4xx (don't retry) from 5xx/network errors (retry with backoff) in fetch error handling.

## Quality Checks

Before marking UNIFY complete, verify:

| Criterion | What to check |
|-----------|---------------|
| Renders without errors | No console errors/warnings on mount and primary interaction |
| Responsive | Tested or verified at 320px, 768px, 1280px minimum |
| Accessible | Keyboard navigable, semantic HTML, no missing alt/labels |
| No regressions | Existing tests still pass after changes |
| Bundle impact | No unnecessary large dependencies added |
