---
id: REQ-278
title: Measure the mark-label face where the sans stack resolves off Linux
status: cancelled
completed_at: 2026-08-19T13:45:20Z
created_at: 2026-08-18T23:53:40Z
status_changed_at: 2026-08-18T23:53:40Z
user_request: UR-051
addendum_to: REQ-265
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Measure the Mark-Label Face Where the Sans Stack Resolves Off Linux

## What

The Durations label row pitch is 13 units, and the largest face anyone has measured draws a 12.9631-unit line box — **0.037 units of slack**. Every measurement this repo holds was taken in a Linux container, and `--font-sans` ends in the open `sans-serif` generic, so **nothing bounds what a Mac or Windows machine actually draws.**

The risk is not hypothetical. `ui-sans-serif`/`system-ui` resolves to Segoe UI on Windows, whose hhea metrics (ascent 2210 / descent 514 over 2048 upem) give a ratio near **1.33** — an 11px line box around 14.6 units inside a 13-unit pitch.

**The strong form of the argument, which the builder's escalation missed:** even unrounded, the max sampled ascent (10.4278, chromium-1228) and the max sampled descent (2.7778, Chromium 141/146) sum to **13.2056** — over the pitch. Two faces this repo has actually measured, recombined, already break it.

**What the consequence actually is, stated honestly:** at ratio 1.33 the *boxes* overlap by roughly 1.6 units while the *ink* does not — row 2's cap-height top sits about 5.3 units below row 1's baseline, and row 1's descender bottom about 2.9 below it. REQ-241 recorded exactly this at pitch 12: "intersected by 0.83 units, which was padding rather than ink — the render showed two clean rows". So off Linux the model's clearance claim stops being true; the board does not visibly break. This is a record-correctness problem, not a shipped rendering defect, and it must not be written up as one.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- **Scope is measuring, not geometry surgery.** Render where `--font-sans` resolves to Segoe UI, SF, or Noto and record the ratio, following this repo's convention: return `location.href` from the same `evaluate` call as every number, and name the build beside every recorded value.
- Do **not** re-budget Panel A's or Panel B's vertical geometry on a number nobody has yet.
- Put the cheaper option on the table explicitly: **giving `.durations-mark-label` a closed font stack the board controls** converts an unbounded question into a bounded one, and is a far smaller change than re-budgeting two panels.
- If the measurement shows the pitch is genuinely exceeded off Linux, the fix is a separate REQ with both constraints visible at once — raising the pitch immediately eats the Panel B ceiling, which REQ-265 just narrowed to 0.10 model units.

## Context

Escalated by REQ-265's builder as D-05 and assessed by its independent review, which confirmed the facts, supplied the stronger argument above, corrected the framing (padding overlap, not visible breakage), and recommended narrowing the scope to measurement.

## Open Questions

- [x] The Durations row pitch has 0.037 units of slack against the largest face ever measured here, every measurement was taken on Linux, and the sans stack ends in an open generic — so nothing bounds what a Mac or Windows machine draws, and two faces this repo has already measured sum past the pitch when recombined. The consequence off Linux is that the model's clearance claim stops being true, not that the board looks broken. Should I process this as a new task? → Discarded — the measurement scope is the wrong direction; see the Cancelled section.
  Recommended: Yes, add to queue (will flip to 'pending') — scoped to measuring, with a closed font stack on the table as the cheap fix.
  Also: No, discard it — accept that the clearance model is only proven for the Linux fallback face.

## Cancelled

- **When:** 2026-08-19T13:45:20Z
- **Why:** user rejected the scope, not the finding. Measuring what Segoe UI, SF, or Noto draw on
  machines this project does not have produces a number that goes stale on the next OS font update
  and still leaves the face unbounded. The user asked for a robust solution instead — explicitly
  open to re-engineering the label layout — and for a solutions report with several options before
  any direction is chosen. Verbatim: "this is the wrong direction, need a more robust HTML solution,
  I'm good even in reengineering the layout if we have to... We definetely do not want to manually
  measure fonts on different operating systems."
- **The finding survives this cancellation.** `--font-sans` ending in the open `sans-serif` generic
  leaves both the vertical pitch and the horizontal per-character width model unbounded. The
  successor work is chosen from the solutions report at
  `ai-reports/2026-08-19_1345_durations-label-face-robustness/index.html`, and gets captured as its
  own REQ once the user picks a direction.
- **Decided by:** user, via `do-work clarify`
