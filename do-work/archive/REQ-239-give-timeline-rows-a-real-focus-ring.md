---
id: REQ-239
title: Give the Timeline's rows a real focus ring
status: completed
completed_at: 2026-08-18T12:27:21Z
commit:
claimed_at: 2026-08-18T11:57:10Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-18T11:57:10Z
  basis:
    - Route B
    - 2-file write set
    - 5 acceptance criteria
    - browser evidence
status_changed_at: 2026-08-18T11:56:00Z
created_at: 2026-08-18T11:09:44Z
user_request: UR-051
addendum_to: REQ-233
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Give the Timeline's Rows a Real Focus Ring

## What

`web/board.css` sets `.timeline-row { outline: none; }`, so a focused Timeline row falls back to `.timeline-row:focus .timeline-row-hit { fill: var(--surface-2); }` — a one-step background tint. Every other focusable thing on the board gets a 2px `--accent-claimed` ring. Give the rows the same ring, and key it on `:focus-visible` so a pointer click does not draw one.

## Context

Found by REQ-233's review. REQ-233 added a visible focus ring to the *chart container* because its requirement said "focus is visible on whatever element takes the keyboard interaction — a focus ring that exists only in the default user-agent style is not enough on a dark surface". The rows next door fail the same test for a different reason: not a user-agent default, but an explicit `outline: none` with a weak substitute.

This is a sibling of REQ-233's requirement rather than part of it — rows are not what REQ-233 added, and rows were already keyboard-activatable before it (pinned by `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard`). It is recorded separately so the fix gets its own before/after rather than riding in on another REQ's merge.

The chart container's ring is the model to copy, including its reasoning: it uses `--accent-claimed`, the token every other ring on the board uses, and `outline-offset: -2px` because the container is flush under the axis. A row's correct offset may differ — a row is not clipped the same way — so this is a judgment, not a copy.

## Requirements

- A keyboard-focused Timeline row draws a focus indicator of the same weight as the rest of the board's rings, using the same token.
- It keys on `:focus-visible`, so a pointer click does not draw one — matching `.control-button:focus-visible`, `.req-card:focus-visible`, and `.calendar-chip:focus-visible`.
- The ring is not clipped by the row's own geometry or by the scroll container.
- The existing `.timeline-row:focus .timeline-row-hit` tint either stays as a complement or is removed deliberately, with the choice stated — two overlapping focus signals is a decision, not an accident.
- No change to row activation behaviour; `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard` still passes.

## Red-Green Proof

**RED prompt/case:** an assertion over the generated stylesheet that `.timeline-row` carries a `:focus-visible` rule with a non-`none` outline, in the same shape the check would use for `.control-button:focus-visible`.
**Why RED now:** `.timeline-row` sets `outline: none` and has no `:focus-visible` rule at all.
**GREEN when:** the assertion passes and a real Tab press onto a row draws a visible, unclipped ring — verified in a browser, not by a programmatic `.focus()`, which does not trigger `:focus-visible` and will report a false negative.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-233 I found that the Timeline's rows have their focus outline explicitly switched off, and get only a faint background tint instead — one shade different from the row's normal colour. Everything else on the board that can be focused gets a clear 2px coloured ring, which is what REQ-233 just added to the chart itself. So a keyboard user moving down the rows has a much weaker sense of where they are than anywhere else on the board. Nothing is broken and the rows still work; this is about how visible the current position is. The fix is a few lines of CSS plus one assertion. I am asking rather than doing it because the tint was written deliberately — someone chose to turn the outline off — and I cannot tell from the code whether that was to avoid a clipped or ugly ring on a dense chart, which is a real concern on rows that are only a few pixels tall. If it was, the answer might be a better tint rather than a ring. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the tint is the deliberate choice for dense rows and should stay.

---

## Triage

**Route: B** - Medium

**Reasoning:** A few lines of CSS, but the REQ's real question — whether a ring is right at all on an 18px virtualized SVG row, given the `outline: none` was written deliberately — could only be answered by rendering it and looking.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — row focus ring
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — stylesheet assertion

**Files I will NOT touch:** `web/board-timeline.js` (no markup change; a sibling builder held it), `web/template.html`.

**Acceptance criteria (restated from REQ):**
- [ ] A keyboard-focused row draws an indicator of the same weight as the rest of the board's rings, using the same token
- [ ] It keys on `:focus-visible`, so a pointer click draws nothing
- [ ] The ring is not clipped by the row's geometry or the scroll container
- [ ] The existing tint stays or goes deliberately, with the choice stated
- [ ] Row activation is unchanged

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** `.timeline-row:focus-visible` draws `outline: 2px solid var(--accent-claimed)` at `outline-offset: -2px` — the same width and token as `.control-button:focus-visible`, `.req-card:focus-visible` and `.timeline-scroll:focus-visible`, drawn inward. `.timeline-row { outline: none; }` stays, now carrying a one-line reason, and the existing `:focus` tint stays as a deliberate complement. No JavaScript was touched at all.

## Qualification

Passed — 2 files verified in the merge range `741e12a..1d76ad1`, 5 acceptance criteria traced.

**Merge conflict resolved by the orchestrator, and the first attempt was wrong.** `generate_test.go` conflicted with REQ-240's appended test. Unlike the two earlier conflicts in this batch, this one was *not* two clean appends: both sides ended mid-function, and git had matched their shared closing braces as a common trailing region. Stripping the markers — which worked for REQ-236 and REQ-235 — produced `expected '(', found TestGenerateGivesTimelineRowsTheBoardsFocusRing`, i.e. one function swallowing the other. The correct resolution duplicates the shared tail so each side's function closes: `ours + closing`, blank line, `theirs + closing`. Both sides' tests then pass, which is the check that the resolution is right rather than merely compiling.

Judgment checks: the ring's token and width are asserted against the board's *reference* ring rather than hardcoded, so a later change to one cannot silently diverge from the other. Nothing hollow — the rule is in the generated stylesheet, which the test reads out of the page.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on the merged tree, run unpiped

**Red-green validation:**
- `TestGenerateGivesTimelineRowsTheBoardsFocusRing`: ✗ `Timeline rows carry no ".timeline-row:focus-visible {" rule: a keyboard-focused row has no ring, only the tint` → ✓. The builder stated plainly that this first RED is the "rule not found" form rather than dressing it up: absence *is* the defect here, so it is the correct failure, but it is a missing-selector error. The test does more than presence once the rule exists — it parses the outline out of both `.control-button:focus-visible` and `.timeline-row:focus-visible` and requires the row's width and token to equal the reference ring, plus a negative offset. Those are the assertions that can fail on drift later.
- `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard` and REQ-240's axis test both pass on the merged resolution.

**Render evidence — and it was re-taken after a contamination warning.** The builder's first pass used the shared browser instance. On being warned mid-build that a sibling had silently invalidated another builder's readings there, it re-took **every** measurement in its own browser context, on a port only that run used, with `location.href`, the page title, and the page's own inlined `cssRuleText` all returned from the same `evaluate` expression that produced the numbers. The title is itself a discriminator — the board titles itself from its repo root, so `worktree-agent-REQ-239-…` can only be that build. **Every re-taken number was identical to the first pass**, and the re-take additionally ran at `deviceScaleFactor: 1`, i.e. a 2-physical-pixel ring rather than 4 — the harsher case.

- Light theme, real `Tab` presses: `outline: 2px solid rgb(58, 107, 196)`, offset `-2px`, `focusVisible: true`, row 18px tall.
- Dark theme: `2px solid rgb(111, 156, 230)` — the dark-palette token — same offset, same behaviour.
- Pointer click: `focus: true`, `focusVisible: false`, `outlineStyle: "none"`. The ring does not appear for the mouse.
- Bottom-edge row: `rowBottom 973` against `containerBottom 973.34` — the ring is drawn complete, not clipped.
- 45 consecutive real `Tab` presses: 45/45 landed on a row with `:focus-visible` and a solid outline; focus was never lost across the virtualized re-render (`lostFocusAt: -1`).
- Every measurement asserted no injected stylesheet was present, so the rule came from the generated page rather than from the probe.

*Verified by work action*

## Decisions

- **D-01**: **Ring, not tint.** The REQ explicitly invited the opposite answer — the `outline: none` was deliberate, and an 18px virtualized SVG row is exactly where a ring might look bad. The builder rendered it and looked: at 18px the ring leaves a 14px interior, still taller than the 10px segment bars, so the bar never touches it and the row does not read as boxed in. It matches `.req-card`'s ring. The fear was reasonable and the render did not bear it out. DECIDE & STATE — with the evidence, because the Open Question was written expecting "tint" to be a live possibility.
- **D-02**: **The existing `:focus` tint stays**, as a deliberate complement rather than by omission — that was one of the REQ's five requirements, and "two overlapping focus signals is a decision, not an accident".
- **D-03**: `outline-offset` is **negative**, matching REQ-233's `.timeline-scroll` decision and for the same class of reason: at a positive offset the ring is clipped on three sides — the rows SVG's own viewport takes left and right, the scroll container takes the top — which paints a divider under the next row instead of a ring around this one. The builder captured a screenshot of the outward version to show it.
- **D-04** (orchestrator): the `generate_test.go` conflict was resolved by duplicating the shared closing braces so each side's function terminates, after a naive marker-strip produced a syntax error. Recorded because the two earlier conflicts in this batch *were* clean appends, and the difference is not visible from the marker positions alone.

## Discovered Tasks

- None.

## Review

**Overall: 97%** | 2026-08-18T12:26:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- Test Adequacy sits at 94% because the assertion is over the generated *stylesheet*, not over rendered pixels — it proves the rule ships with the right token, width and offset sign, and cannot prove the ring is unclipped. That part rests on the builder's browser measurements, which are thorough (both themes, bottom edge, 45 tab presses, pointer negative case) but are evidence in a hand-back rather than a check that reruns. That is the honest ceiling for a CSS defect without a visual-regression harness, and worth naming rather than scoring around.

**Restatement sweep:** the diff adds a focus rule; the board's ring convention is stated by the rules themselves rather than in prose, and this test now pins the row's rule *against* the reference rule so the convention has a mechanical statement for the first time. No document restates it. `_dev/primes/prime-kanban-board.md` gained render-evidence conventions this session and is unaffected here. No stale restatement.

**Acceptance:** Pass — the requirement most at risk (that a ring is even appropriate) was answered by rendering rather than assumed, in both themes and at the bottom edge, with the pointer negative case confirmed.

**Suggested testing:** 1 item
- A visual-regression harness would convert the builder's screenshots into a check that reruns. Out of scope here and probably its own decision, but this REQ is the second in the batch whose real evidence is a picture nothing re-verifies.

**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Treating the Open Question as genuinely open. The REQ said the answer might be a better tint rather than a ring, and gave a real reason — 18px rows, deliberate `outline: none`. Rendering it settled the question in one look: 14px interior, taller than the 10px bars, crisp at 2 physical pixels. Had the builder assumed the requirement's presumed means was its goal, it would have shipped without knowing whether it looked right.

**What didn't:** Nothing in this build — but the contamination warning arriving mid-build is the reason its evidence is trustworthy. The builder re-took every measurement rather than reasoning that its own readings were probably fine, and reported that nothing differed. A re-take that confirms the original is worth as much as one that overturns it, and cheaper than the argument about whether it was needed.

**Worth knowing:** Two merge conflicts in this batch were clean appends and this one was not, with nothing in the marker positions to distinguish them. Both sides ended mid-function and git had folded their shared closing braces into the common tail, so stripping markers made one function swallow the other. The check that catches it is compiling — and then running *both sides' tests*, because a resolution can compile and still have dropped an assertion.

## Orientation

A keyboard-focused Timeline row now draws the same 2px `--accent-claimed` ring as every other focusable thing on the board, inset so the rows SVG and the scroll container cannot clip it, with the existing background tint kept as a complement. Pointer clicks still draw nothing. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

Not `[MAP CHANGED]` — one CSS rule and one assertion. It does close the pair REQ-233 opened: that REQ gave the chart container a ring and left the rows inside it on a faint tint, which was the visible asymmetry. The board's focus-ring convention also gains its first mechanical statement, since the new test compares the row's rule against the reference rule rather than against a hardcoded value. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves and the three-write-surface count is unchanged. The prime is not stale.
