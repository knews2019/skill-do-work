# Spec: UI Component

> Specification template for frontend UI components — structure, states, accessibility, and tests.

## Output Structure

- **Component file** — markup, logic, props interface
- **Styles** — scoped styles following the project's styling approach (CSS modules, Tailwind, styled-components, etc.)
- **Tests** — rendering tests, interaction tests, accessibility checks
- **Demo/storybook** — if the project uses Storybook or a similar tool, add a story (optional)

## Quality Standards

- Accessible (WCAG 2.1 AA) — keyboard navigable, proper ARIA attributes, sufficient color contrast, semantic HTML
- Responsive — works at 320px, 768px, 1280px minimum
- All states handled: loading, error, empty, populated
- Keyboard navigation works for all interactive elements (focus management, tab order, enter/space activation)
- No hardcoded strings for user-facing text — use the project's i18n system if one exists
- Error boundaries — component failures don't crash the entire page
- `prefers-reduced-motion` respected for animations

## Implementation Checklist

1. Structure — component skeleton, props interface, semantic HTML
2. Base component — core layout and content rendering (populated state)
3. States — loading skeleton, error display, empty state, edge cases
4. Interactions — click handlers, form inputs, hover/focus states
5. Accessibility — ARIA labels, keyboard support, screen reader testing
6. Responsive — mobile-first layout, breakpoint adjustments
7. Tests — render tests, interaction tests, state transitions

## Evolution Path

- **Simple**: Static display component — receives props, renders content, no internal state
- **Medium**: Interactive with state — form inputs, toggles, local state management, event handlers
- **Complex**: Composable with context/hooks — shares state via context, exposes custom hooks, composition API, slot/children patterns

## Common Pitfalls

- Missing keyboard support — mouse works, keyboard doesn't. All clickable elements need `onKeyDown`/`onKeyUp` handlers or use semantic `<button>`/`<a>` elements
- No loading state — component shows blank or flickers while data loads
- Hardcoded strings — makes i18n impossible later
- Missing error boundaries — one broken component takes down the whole page
- Inaccessible custom controls — custom dropdowns, modals, and tabs that don't announce to screen readers
- Layout shift — images and async content without reserved dimensions cause content to jump
