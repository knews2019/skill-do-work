---
id: REQ-376
title: 'Raise the done line''s faint text to readable contrast'
status: completed
status_changed_at: '2026-08-27T11:33:37Z'
created_at: 2026-08-26T14:40:00Z
user_request: UR-074
addendum_to: REQ-374
domain: ui-design
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
claimed_at: '2026-08-27T11:33:37Z'
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T11:33:37Z'
write_set:
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go
completed_at: '2026-08-27T11:45:22Z'
commit: 8dfdb24
kb_status: promoted
kb_entry: REQ-376-raise-the-done-line-s-faint-text-to-read.md
---

# Raise the Done Line's Faint Text to Readable Contrast

## What

The done line's faint companions — `.relative-time` and `.elapsed-duration` — measure roughly 3.3:1 against `<body>` in both themes at 11px, under the 4.5:1 needed for text that size.

## Context

Discovered while working on REQ-374. Pre-existing for the whole line, not introduced by that REQ: measured light `rgb(108,116,128)` at 0.85 opacity on `rgb(245,247,250)`, dark `rgb(107,116,128)` at 0.85 on `rgb(12,14,18)`. REQ-374's new span reading deliberately inherits the same treatment rather than diverging mid-line, which is why the fix belongs to the line as a whole.

Per the prime, measure against `getComputedStyle(document.body).backgroundColor` — the board's SVG and card surfaces are transparent, so a tone judged against a `--surface-*` token is measured against something the reader never sees.

A related but separate gap, reported by REQ-374's review and not covered here: both new markers convey their meaning only through `title`, which a screen reader does not reliably announce and which is unreachable by keyboard on a non-focusable span.

## Red-Green Proof

**RED prompt/case:** compute the contrast ratio of `.relative-time` and `.elapsed-duration` against `getComputedStyle(document.body).backgroundColor` in a rendered board, in both themes. Both read about 3.3:1.
**Why RED now:** the tokens and the 0.85 opacity together land under 4.5:1 for 11px text.
**GREEN when:** both read at least 4.5:1 in both themes, with the done line still visibly quieter than the card title — the point of the treatment is hierarchy, and a fix that flattens it has traded one defect for another.
**Validation:** Inferred during capture.

## Open Questions
- [x] I discovered this out-of-scope task while working on REQ-374: the done line's faint text measures about 3.3:1 against the page background in both themes, under 4.5:1 for 11px text. Should I process this as a new task? → Confirmed: Yes, add to queue.
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - **Answered 2026-08-27:** The maintainer explicitly approved this task during clarify, independently confirming the broader approval recorded by the concurrent work session. Improve contrast in both themes while preserving the quieter completion-line hierarchy. The separate keyboard and screen-reader tooltip issue remains out of scope. This is approval to build, not confirmation of an existing implementation; status remains `pending`. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.

## Triage

**Route: A** — Focused completion-line palette repair with measured contrast in both themes.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Measure both completion companion styles against the actual body background; change only their scoped color/opacity treatment; pin contrast and hierarchy in a browser test.
- [x] **[APPLY]:** Scoped soft-ink/full-opacity treatment to terminal-card completion companions; added a real-browser contrast regression with nonterminal controls.
- [x] **[UNIFY]:** Reviewed both source files and measured screenshots; formatter and focused browser proof passed. Canonical gate recorded in Testing.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — readable completion-companion colors scoped to terminal cards.
- `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go` (new) — rendered contrast, hierarchy, both themes, and unchanged nonterminal controls.

**What was done:** Completion relative-time and elapsed-duration text now clears 4.5:1 against both the actual card and page backgrounds without changing pending or claimed text.

## Decisions

- D-01 (decide and state): Reuse the existing soft ink token at full opacity, scoped by terminal status. The completion-looking class is also used by nonterminal cards, so class-only scoping would change unrelated displays.

## Orientation

Completion-time companions are readable in both themes while remaining smaller and quieter than card titles. Other timestamp displays retain their existing appearance.

## Qualification

Passed: qualify.sh exit0; both substantive files traced to the requested contrast repair. The unreferenced-new-file warning is the normal Go test discovery exception. Parent inspected both actual-card screenshots and reran the browser test independently, exit0.

## Testing

**Red-green validation:** `TestBrowserBehaviorCompletionCompanionsKeepReadableContrast` failed before CSS changed and passed after. The initial expired fixture was a setup failure, not counted as RED. Final browser proof used Chrome for Testing141.0.7390.37 and both themes at320,768,1280px. Parent independently reran it, exit0. Pending/claimed controls retain the original styles.

| Surface | Before | After |
| --- | ---: | ---: |
| Light body/card | 3.35 / 3.28 | 5.69 / 5.54 |
| Dark body/card | 3.26 / 2.97 | 7.41 / 6.55 |

Both classes share these ratios. Source formatter, go vet, and existing done-span payload/render tests passed. Parent and reviewer inspected full completed-card crops in both themes; no clipping and title hierarchy retained.

## Review

**Overall: 98.75%** | Acceptance: Pass. No important findings or follow-ups. Independent reviewer reran the new contrast test and adjacent span tests successfully. Minor: completed-with-issues and cancelled share the selector but lack separate browser fixtures.

## Lessons Learned

The completion-looking CSS class also appears on pending and claimed lines; scope a done-card change by status and keep a nonterminal control. The cards are opaque even though the chart SVG is transparent, so measure the actual card as well as the body. Time-sensitive browser fixtures must remain inside the production Recently Done window.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
