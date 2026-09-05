# Independent review — REQ-587 (give the Timeline view one scroll surface)

Reviewed from the merge range `93ec7792..8fad73b2` in a detached worktree cut at `8fad73b2`. Every number below was measured by this review, not read from the hand-back. Main tree never touched, nothing committed, nothing under `do-work/` written. Both review worktrees removed.

---

## Review: REQ-587

**Approve with follow-ups** — the Timeline now has one scroll surface, reproduced in a real engine by this review; six report-only findings, one of them a real test gap on the very defect the REQ paid for.
Route C | merge range `93ec7792..8fad73b2` (merge commit `8fad73b2`, builder branch head `564c81a3`)

### What's built

- The Timeline's chart is no longer its own scroll box. `.board-main` is the single scroll surface, the time axis is pinned to the board's top edge with rows passing under it, and the five sites that read or wrote a scroll position now work in the board's coordinate space.
- The board's 24px top padding moved onto whichever child of `.board-main` the reader actually sees first, keyed on "not hidden" rather than on a list of the strips that can precede the view panel.
- A new browser probe pins the result with four guards and eight assertions, including the REQ's own RED expression run verbatim and the anchor behaviour the orchestrator asked to be pinned rather than hand-checked.

### Decisions / risks for you

- **D-06 (the focus ring moved to the sticky axis) — escalated, and it holds up.** I measured the state the plan worried about: with the board scrolled to its maximum the axis is still pinned and visible (offset 0.0px from the board's top edge), and focusing the chart from the bottom of the page scrolls the board back (4543 → 3924) with the axis at 0. There is no measured state where a keyboard user gets no indication. The residual is F6 below.
- **D-07 (entering the Timeline resets the board's scroll to 0) — escalated, and the "cannot fire mid-session" claim is true.** `applyView` has exactly two production callers: the view-button click (`web/board-controls.js:200`) and init (`web/board.js:79`). The reset is gated on `viewState.view === "timeline"` and sits after the panel toggle and before the first render. The cost is real and correctly stated: leaving the Timeline and coming back loses your place, which the inner box used to keep. The trade is honest for a page where the chart is now the tall thing.
- **D-22 (one file beyond the declared Scope).** `timeline_browser_probe_test.go` was disclosed by the builder, verified green at the base revision, accepted by the orchestrator, and both the `## Scope` list and `write_set` were extended before integration. I confirmed the base-green claim myself. This is the process working, not silent drift.

### Findings

**Important:**

- **F1 — the D-15 fix has no regression test, and the whole lane stays green with the defect restored.** `web/board-timeline.js:2082` (`rowsScrollTop`). I restored the clamp the plan specified (`Math.max(0, boardScrollHost.scrollTop - rowsOffsetPx())`) and ran the full timeline + activity browser lane and the Node lane: **`ok`, everything green**. I then measured the behaviour that clamp produces, in the same fixture: pressing the All-days chip while the board sits above the chart scrolls it **0 → 636px**, dragging the view heading from **+443.8px to -192.2px** — off screen above. The shipped, unclamped code measures **0 → 0** with the heading unmoved. So the code is right and the guard is absent: the new probe's anchor scenario always starts from a positive rows-scroll position (`timelineScrollProbeAnchorScrollPixels = 500`), and the negative-position path — the one the builder found by hand and the only one that can round-trip wrong — is never exercised. One extra anchor case starting at `boardMain.scrollTop = 0` would close it. — impact-user-visible → report only

**Minor:**

- **F2 — D-16's recorded reason for the fourth hunk is wrong in its conclusion.** The hand-back says the vacuity guard "compared a number with itself and could never fire again". It would always fire, not never: I reverted that hunk alone and got `timeline_browser_probe_test.go:3784: Fit-all SVG height 14422 does not exceed the 14422px viewport, so the probe never exercised virtual scrolling` — a hard failure. The fix (re-pointing `viewportHeight` at `#board-main.clientHeight`) is correct and restores the guard's meaning; only the written justification is inverted, and a later reader trusting it would look for a silent pass instead of a red probe. — impact-negligible → report only
- **F3 — Restatement Sweep: `web/board-timeline.js:150` still calls `#timeline-scroll` "the scroll container".** The comment reads "This view binds to the scroll container and to window, which BOTH outlive a render". After this change the view binds to three things — the chart (no longer a scroll container), `#board-main` (which is), and window — and the phrase now names the wrong element in the one file where the distinction is the whole point of the REQ. Every other "inner box" / "58vh" mention in `web/` is correctly past tense ("used to be", "used to mean", "never mattered while"). — impact-negligible → report only
- **F4 — Restatement Sweep, outside this REQ's Scope: `_dev/primes/lessons-kanban-board.md:71` states the padding recipe too narrowly.** The `[family: sticky-pins-to-content-box]` entry from REQ-585 says to "move the board's top padding onto the view's first element, scoped to that view (`.board-main:has(> #view-activity.is-active)`)". REQ-587 found that recipe insufficient: three strips and a JS-inserted warnings banner sit outside the view panels and can precede the active one, so the padding has to land on the board's first *unhidden child*. A builder following the satellite as written on the next view would reproduce the bug this REQ avoided. The satellite is the canonical home for that rule, so it should carry the generalised form. Lesson capture belongs to work.md's Lessons-Capture Phase in orchestrated mode, so this is a pointer for that phase, not a defect in the diff. — impact-rule-change → report only

**Nit:**

- **F5 — the "no strips at all" arrangement of D-05 is never measured.** The probe's fixture always produces the warnings banner and the anomalies strip, and asserts `viewPanelMarginTop == "0px"` (rule 3 working with a preceding sibling). The arrangement where `#view-timeline` itself carries the 24px is only reasoned about, never read back. The three rules are simple enough that I verified it by construction — with the strips hidden and the other five panels carrying `hidden` from `applyView`, `#view-timeline` is the only `:not([hidden])` child and rule 2 applies to it alone — but nothing pins it. — impact-negligible → report only
- **F6 — `.timeline-scroll:focus-visible { outline: none }` is unconditional while its replacement depends on `:has()`.** In an engine without `:has()` support the chart is focusable with no visible indication at all, where the same missing support costs the padding rule only 24px. The dependency is not new (`board.css:2138` predates this REQ and REQ-585 added another) and `:has()` is baseline in every engine this tool targets, so this is a note rather than a defect. — impact-negligible → report only

### Requirements Checklist

- [x] Exactly one element scrolls on the Timeline — **delivered.** My own run: `red=[true false]`, `overflowY=visible`, chart `3958/3958`, board `689/4864`.
- [x] The inner box's computed `overflow-y` is `visible`, not merely equal heights — **delivered**, and this is the assertion that carries the weight (see mutation M1).
- [x] The time axis stays visible, pinned to the board's top edge, no rows painted above it — **delivered.** `axisTop=0.0`, `rowsAboveAxis=0` after a 1500px scroll.
- [x] Rows still render at every scroll position, proved by disjoint rendered sets — **delivered.** 39 ids before, 44 after, intersection empty, 34 intersecting the visible rect.
- [x] The row under the pointer stays put when the window chips change — **delivered.** `REQ-7003`, display list 62 → 201, board 1209 → 1893, offset -4.4 → -3.9, drift 0.5px.
- [x] Wheel-zoom, drag-pan and the keyboard path all still work — **delivered.** All five named regression probes green in my run, plus the whole lane.
- [x] Same style as the Activity fix, reusing REQ-585's scoping rather than inventing a second mechanism — **delivered.** Same `:has()` shape, extended to cover the strips; the Activity block is byte-identical and its own probe is green.
- [x] Measured in a real engine, browser and build named, recorded in the hand-back — **delivered.** Chrome 152.0.0.0 in the hand-back; HeadlessChrome/152.0.0.0 in my independent run.

8 of 8 delivered.

### Acceptance Testing

**Result: Pass**

Run in a detached worktree cut at `8fad73b2`, engine named explicitly, output checked for SKIP lines rather than trusting the exit code.

- `QUEUE_KANBAN_BROWSER=<Chrome> QUEUE_KANBAN_BROWSER_PROBES=on go test -count=1 -v -run 'TestBrowserBehavior' .` → **65 RUN, 65 PASS, 0 FAIL, 0 SKIP**, 131.8s. Zero occurrences of `SKIP`/`skipping`/`skipped` anywhere in the output.
- `go test -count=1 ./...` → `ok`, 42.6s. `go vet ./...` clean, `gofmt -l .` empty, `go build ./...` clean.
- The four probes D-16 updated were confirmed green at the **base** revision `93ec7792` in a second worktree — `PressBecomesAPanOnlyAfterMoving`, `PointerAndKeyboardPathsStayAlive`, `PointerCaptureWaitsForThePanEngage`, `ListsRowsBeneathUserRequestHeaders`, all PASS. The claim that they were green before the change is true, so the four hunks are adaptations to an intentional behaviour change, not cover for a pre-existing failure. All four hunks are mechanical: three add one synchronous `dispatchEvent(new Event("scroll"))` on `#board-main` after a `scrollIntoView`, driving the shipped listener with the shipped handler; the fourth re-points a vacuity guard's viewport number. No assertion's meaning changes.

**Mutation testing — what the new probe actually catches.**

| Mutation | Caught? | By what |
|---|---|---|
| **M1** — restore both `overflow` declarations on `.timeline-scroll`, leave `max-height` gone | **Yes** | The computed-overflow assertion alone (`timeline_scroll_browser_probe_test.go:267`). The height comparison passed, and the REQ's own RED expression still returned `[true false]`. D-02's claim that the height comparison cannot see this is exactly right, and the extra assertion is the only thing standing between "fixed" and "looks fixed". |
| **M2** — move the `scroll` listener back to `scrollHost` | **Yes, three times over** | A4 (`:295`, 39 of 39 rows still drawn after the scroll) and the in-viewport assertion (`:299`), plus `TimelinePressBecomesAPanOnlyAfterMoving` and `TimelinePointerCaptureWaitsForThePanEngage` go red. The anchor assertion stays green, correctly — it does not depend on the listener. |
| **M3** — restore the plan's clamp inside `rowsScrollTop` (the D-15 defect) | **No** | Full timeline + activity browser lane `ok`; Node lane `ok`. See F1 — with the defect live and measurable (board 0 → 636, heading +443.8 → -192.2), nothing turns red. |
| **M4** — revert the vacuity-guard re-point (D-16 hunk 4) | **Yes** | `ListsRowsBeneathUserRequestHeaders` fails at `:3784`. This is what proves F2: the guard fires rather than going silent. |

**D-15 audited from the code, at all five consumers.** Three read `rowsScrollTop()` and two use `rowsOffsetPx()` directly. `topVisibleRowAnchor` needs the signed value or the write disagrees with the read (the measured defect). `renderVisibleRows` floors it at its own call site and a negative value can only widen the range upward, so it is a superset. `scrollTimelineRowIntoView` is the case that needs the sign most: with the board at 0 and a 689px offset, `visibleBottomPx` evaluates to 53 in rows coordinates, which is exactly document y 774 — a row genuinely on screen near the bottom correctly triggers no scroll, where a clamped value would have scrolled it. The write in `refreshWindowRows` and the jump-to-open-work fallback both go through `rowsOffsetPx()`/`boardScrollTopForRowTop`, which is clamped and produces a positive target in every case. **I found no consumer where the unclamped value produces a wrong result.** The one case I expected to be wrong — the reader above the chart pressing a chip, where preserving a row's screen position could drag the page — measures 0 → 0 with the heading unmoved, because the anchor there is the display list's first row and a window widening cannot insert anything above it.

**D-05 audited across every arrangement.** `applyView` sets `panel.hidden = !isActiveView` on all six panels (`web/board-controls.js:30`), and all three strips are toggled by the `hidden` property, not a class or `display:none` — `renderAnomaliesStrip` (`web/board-cards.js:612,616`), `renderVerifyFindingsStrip` (`:689,693`), the notes strip, and `applyView`'s own `findingsStrip.hidden = …` (`:73`). The warnings banner is inserted or absent, never hidden. `#board-main` has exactly nine static children plus that banner; no `<script>`, `<style>` or other always-present-but-invisible child exists that could eat the margin. So: no strips → the panel alone carries 24px; one strip → the strip carries it and the panel gets 0 from rule 3; several strips → the first carries it; a hidden strip → skipped by `[hidden]`. Exactly one element carries it in every arrangement, and the probe measured the hardest one — the JS-inserted, id-less warnings banner at 24.0px with `viewPanelMarginTop: 0px`.

**Accessibility, measured rather than reasoned.** With the board at its maximum scroll the axis is still pinned (offset 0.0, visible), because the chart's bottom edge is still inside the viewport and the sticky element cannot be pushed out of its own containing block. Calling `focus()` on the chart from the very bottom of the page scrolled the board 4543 → 3924 and left the axis at offset 0 and visible. The selector works structurally: `#timeline-scroll` and `.timeline-axis` are both direct children of `.timeline-chart` (`web/template.html:461,465`), so `:has(> …)` reaches. The keyboard affordance the hint text promises is preserved.

### Restatement Sweep

**Applied — the diff redefines which element is the Timeline's scroll container, and that meaning is stated in several places.** I swept `_dev/primes/prime-kanban-board.md`, `_dev/primes/lessons-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, `lessons-do-kanban.md`, `skills/do-work-board/actions/board.md`, `skills/do-work-board/docs/board-guide.md`, every comment in `web/`, and the Go test sources, for `timeline-scroll`, `58vh`, "scroll container", "scroll host", "scroll surface", "inner box", `max-height` and `board-main`. Two stale restatements found, both reported above: **F3**, `web/board-timeline.js:150`, which still calls `#timeline-scroll` "the scroll container" in the present tense; and **F4**, `_dev/primes/lessons-kanban-board.md:71`, whose `[family: sticky-pins-to-content-box]` entry states the padding recipe in the narrower form REQ-587 proved insufficient — a file this REQ never declared, reported anyway because that is what the sweep is for. Everything else agrees with the new meaning: the prime's only `#timeline-scroll` mention (`:29`) is about background tokens and is untouched by this change; `board-guide.md` and `board.md` say nothing about the scroll model; the shipped hint paragraph ("Scroll down to go further back inside the window", "Tab to the chart") stays true; and every remaining `58vh`/"inner box" mention in `web/` is deliberately past tense and correct.

### Coding-guardrails and directive notes (informational)

- **Think Before Coding** — D-01 recorded at Step 3.5, D-02–D-13 by the plan, D-14–D-21 by the builder, D-22 by the orchestrator. Every escalation carries value, risk and a stated reversal. Nothing was silently resolved.
- **Simplicity First** — the memo pair reuses the file's existing `measuredPlotWidthPx`/`invalidatePlotWidth` idiom rather than inventing a second caching shape, and the dead `timelineVisibleRowRange` was deleted with the one comment naming it reworded. Delete-before-you-add honoured.
- **Surgical Changes** — every changed line traces to the REQ. The four Node-stub hunks are one string each. The Activity view's CSS block is byte-identical.
- **Goal-Driven Execution** — red-green is real and independently reproduced, and the Node lane's byte-identical row counts and row ids after the host split are the strongest single piece of evidence that the coordinate conversions are no-ops where they should be.
- **Naming for Reach** — every new name with reach is two words or more and survives a plain-text grep: `boardScrollHost`, `rowsOffsetPx`, `stickyAxisHeightPx`, `invalidateBoardScrollGeometry`, `rowsScrollTop`, `boardScrollTopForRowTop`, `timelineScrollProbeRequestCount`, `timelineScrollSurfaceMeasurement`, `assertTimelineAnchorSurvivesAWindowChange`, `intersectingRowIds`. No finding.
- **P-A-U** — all three boxes `[x]`, and the UNIFY box records per-file checks that match what I found in the diff.
- Diff hygiene clean: no `console.log`, `debugger`, `TODO`, `FIXME` or `XXX` anywhere in the added lines. `git diff --stat` over the range confirms 9 files, 937 insertions, 46 deletions — the Qualification section's correction of the hand-back's stale "320 insertions, 51 deletions" is right.

### Suggested Additional Testing

1. **Close F1 with one assertion.** Add a second anchor case to `assertTimelineAnchorSurvivesAWindowChange` starting at `boardMain.scrollTop = 0` and asserting the board does not move. Without it the D-15 fix is unguarded and the next edit to `rowsScrollTop` gets a false green.
2. **The focus ring still has no automated evidence.** D-14's defect passed every assertion and was caught only by a screenshot. That is a documented, deliberate trade (a programmatic `focus()` cannot answer a `:focus-visible` question), but it means any future change to the two focus rules needs a human eye on it. Worth naming in the prime's traps if it recurs.
3. **Manual pass on a real (non-headless) engine at a small viewport.** Below the 760px breakpoint the moved padding is 24px where the board's own was 18px, an accepted 6px difference. Worth one look at 320px and 768px to confirm it reads as intentional.
4. **Re-entry to the Timeline after scrolling deep into it** (D-07's stated cost). Confirm with a human that losing your place on return is the behaviour the maintainer wants, since the alternative is six lines.
5. **A queue with notes and no findings.** The probe's fixture produces the warnings banner and the anomalies strip. The notes-only and nothing-at-all arrangements of D-05 are covered by reasoning, not measurement (F5).

### Scores (on the record — not the headline)

**Overall: 92%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | 8 of 8 delivered, each independently re-measured in a real engine |
| Code Quality | 92% | Reuses the file's own idioms, comments carry the reasoning; one stale comment (F3), one inverted justification (F2) |
| Test Adequacy | 82% | The probe catches M1 and M2 decisively and the Node lane's identical output is real proof; the D-15 fix has no guard at all (F1) |
| Scope | 95% | 9 files against 8 declared; the ninth was disclosed, verified green at base, accepted, and both Scope and write_set extended before integration |
| Risk | Low | Contained to one view; one listener on the shared board, released on every render; the scroll reset is gated on the view and cannot fire mid-session |
| Acceptance | Pass | 65 browser probes PASS, 0 SKIP; full package `ok`; vet, gofmt, build clean |

### Follow-ups created

None (6 findings report only)

---

## Append to REQ File

```markdown
## Review

**Overall: 92%** | 2026-09-05T20:31:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 82% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The D-15 fix has no regression test: restoring the plan's clamp in `rowsScrollTop` (`web/board-timeline.js:2082`) leaves the whole browser lane and the Node lane green, while the reintroduced defect measurably scrolls the board 0 → 636px and drags the view heading from +443.8px to -192.2px on a window-chip press from above the chart. The new probe's anchor case always starts from a positive rows-scroll position, so the negative path is never exercised — impact-user-visible → report only

**Minor findings:**
- D-16's recorded reason for the vacuity-guard hunk is inverted — the guard would always fire, not "never fire again"; reverting that hunk alone produces a hard failure at `timeline_browser_probe_test.go:3784`. The fix is correct; the written justification is not — impact-negligible → report only
- Restatement Sweep: `web/board-timeline.js:150` still calls `#timeline-scroll` "the scroll container" in the present tense, in the one file where that distinction is the point of the REQ — impact-negligible → report only
- Restatement Sweep, outside this REQ's Scope: `_dev/primes/lessons-kanban-board.md:71` states the padding recipe as "the view's first element", the narrower form REQ-587 proved insufficient once strips can precede the view panel; the satellite is the canonical home and should carry the generalised rule — impact-rule-change → report only
- Nit: the "no strips at all" arrangement of D-05 is verified by construction but pinned by no assertion — impact-negligible → report only
- Nit: `.timeline-scroll:focus-visible { outline: none }` is unconditional while its replacement ring depends on `:has()`, so an engine without `:has()` gets no focus indication; the dependency predates this REQ and `:has()` is baseline for this tool's engines — impact-negligible → report only

**Acceptance:** Pass — 65 browser probes run at the merge commit in a detached worktree with the engine named, 65 PASS, 0 SKIP; full package `ok` in 42.6s; vet, gofmt and build clean; the four probes D-16 updated confirmed green at base revision 93ec7792.
**Suggested testing:** 5 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*
```

---

## Note on the repository gate

I agree with the orchestrator's reading and did not spend time on it. Three consecutive failures on per-test-file wall-clock budgets, each on a different file, all in packages this REQ does not touch, with load 22–40 from sibling sessions, is a measurement artefact and not a fact about this diff. My own evidence supports that: the `queue-kanban` package — the only package this REQ touches — runs `ok` in 42.6s in an isolated worktree, and the full 65-probe browser lane runs green in 131.8s. Nothing in this change is near a budget.
