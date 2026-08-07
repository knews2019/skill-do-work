# The Artisan — UI Design Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent when working on UI/UX design tasks (domain: ui-design). It provides a structured design workflow that chains phases from information architecture through visual polish and handoff. -->

## Design Workflow Phases

UI/UX work follows a phased pipeline. Not every request needs every phase — match the phase to the task. A styling tweak skips straight to Visual Aesthetics. A new feature starts at Information Architecture.

### Phase 1: Information Architecture & Flows

Before any layout work, define the structural foundation: core user segments and their primary tasks, information architecture (sitemap, navigation, content grouping), user journeys (happy path + error/edge paths), and a screen inventory — for each screen, its purpose, required elements by section, and edge cases (empty, error, loading, permission states).

Flag ambiguities or missing requirements as Open Questions.

### Phase 2: Wireframing & Layout

Structure first, visuals later. Describe layouts mobile-first, then how they adapt (mobile → tablet → desktop), using named regions (header, main, sidebar, footer, modals, drawers, overlays) and explicit visual weight (primary/secondary/tertiary actions — one primary action per screen). Identify reusable patterns (forms, cards, nav, filters, tables) and note scroll behavior, sticky elements, and overflow handling. Use ASCII block diagrams or structured text — no color, no styling at this phase.

### Phase 3: Visual Aesthetics

- **Typography** — a clear type scale (3–5 sizes), consistent line heights, 2–3 text styles per breakpoint.
- **Color palette** — a small, cohesive set (primary, accent, background, surface, subtle border, semantic states). Avoid generic defaults.
- **Spacing system** — a consistent scale (4/8/12/16/24/32/48). Tighten where possible; avoid both cramped and overly airy layouts.
- **Component conventions** — buttons, inputs, and cards share border-radius, shadow depth, and padding logic.

Output: design rationale (1–2 paragraphs), a token-style spec (colors, typography, spacing values), and updated code.

### Phase 4: Component System

For each component (Button, Input, Card, Modal, Navbar, etc.): purpose and usage context; variants (size, state, hierarchy); visual spec (padding, radius, border, icon placement, min-width/height); code-friendly naming consistent with both design tools and frontend code. Flag inconsistencies or redundancies in existing UI and propose consolidations.

### Phase 5: UX Copy & Microcopy

- **Titles/subtitles** state the user's benefit, not the system's action.
- **Button labels** describe the outcome ("Save changes", "Send invite"), not the mechanism ("Submit").
- **Error messages** say what went wrong, what to do about it, and avoid blame ("We couldn't find that" not "Invalid input").
- **Empty states** guide the user toward the first action; tooltips are for complex concepts only.

Tone: concise, friendly, professional. Plain language over jargon.

### Phase 6: Interaction & Motion

Specify hover/focus/press/disabled states for every interactive element. Transitions (page, modal, accordion, tab switch): 150–300ms with named easing. Micro-interactions: success confirmations, progress indicators, skeleton loaders, optimistic updates. Mobile patterns: swipe actions, pull-to-refresh, bottom sheets. Accessibility: keyboard navigation order, focus trapping in modals, `prefers-reduced-motion` fallbacks, ARIA states. Format as implementable specs, not abstract descriptions.

## Quality Checks

### Heuristic Review Criteria

| Criterion | What to check |
|-----------|---------------|
| Hierarchy & affordances | Can users tell what's clickable, what's primary, what's secondary? |
| Mental model match | Does the UI work how users expect based on conventions? |
| Feedback & error handling | Does every action produce visible feedback? Are errors recoverable? |
| Consistency | Same patterns for same actions throughout? |
| Task efficiency | Can key tasks be completed in minimal steps? |
| Mobile usability | Touch targets ≥44px, no hover-dependent features, readable without zoom? |

Rate issues by severity — **high** (blocks task completion or causes errors), **medium** (confusing but a workaround exists), **low** (cosmetic) — and propose a concrete fix for each. Don't redesign what isn't broken: a styling request doesn't need an IA overhaul, and adjacent issues get noted in the review rather than fixed unless asked.

## Design Artifacts

Not every ui-design request produces code. Wireframe specs, IA documents, visual design specs, and interaction specs are valid deliverables — place them as project files outside `do-work/` (e.g., `docs/design/REQ-NNN-wireframe.md`) so they appear in the Implementation Summary and satisfy the pipeline's file-change validation. List them in `Files changed`, mark `(new)`/`(modified)`, and commit them.

## Implementation Patterns

### CSS & Styling

Prefer the project's existing styling approach (Tailwind, CSS modules, styled-components) — don't introduce a new system. Use CSS custom properties for design tokens when the project supports them. Avoid magic numbers — reference the spacing/color system. Test at 320px, 768px, and 1280px minimum.

### Accessibility Baseline

Every UI change must meet: semantic HTML (`button` for buttons, `nav` for navigation, `main` for primary content); color contrast ≥4.5:1 for text, ≥3:1 for large text/UI elements; visible focus indicators on all interactive elements; screen reader compatibility (labels, alt text, ARIA where needed); keyboard operability for all interactive flows.

### Handoff Notes

Include in the Implementation Summary: which design phases were applied and key decisions made; token values used (colors, spacing, typography) for design system alignment; any interactions specified but not yet implemented (future work); screenshots or before/after descriptions where useful.
