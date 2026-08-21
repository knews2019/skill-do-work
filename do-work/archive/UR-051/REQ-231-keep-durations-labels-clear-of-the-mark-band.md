---
id: REQ-231
title: Keep Panel A's direct labels clear of the mark band
status: completed
completed_at: 2026-08-18T10:52:30Z
commit: 720f23c
claimed_at: 2026-08-18T10:15:00Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-18T10:35:12Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
domain: general
created_at: 2026-08-18T00:55:10Z
user_request: UR-051
addendum_to: REQ-226
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/web/board.css
---

# Discovered Task: Keep Panel A's Direct Labels Clear of the Mark Band

## What

In the Durations view's overflow lane, the first row of direct labels sits at the same height as the marks themselves, so in a dense lane a label can be crossed by a *neighbouring* mark. REQ-226 stopped labels from overprinting each other; this is the remaining overlap, between a label and a dot that is not its own.

## Context

Found while implementing REQ-226. `DURATIONS_LANE_MARK_Y` is 40 with a mark radius of 5, so marks occupy roughly y 35-45; `DURATIONS_LANE_LABEL_ROW_Y` is 44, so a first-row label's text box occupies roughly y 33-46. A label always clears its *own* mark, because it is drawn 9 units to one side of it, so the overlap only appears where the band is dense enough for other marks to crowd the label.

REQ-226's collision rule is deliberately label-against-label — its Requirement 1 says "skip any that would collide with one already placed". Extending it to label-against-mark would drop labels wherever the band is dense, which is the opposite of what the remainder count was added for. The straightforward fix is instead geometric: give the lane roughly 12 more user units so both label rows sit below the marks, and shift the panels beneath it down to match.

REQ-226 could not do that: its constraints state "Panel A's existing scale-break design is correct and stays. This REQ fixes the labelling on top of it; it does not redesign the panel." Changing the lane's height is that redesign, so it is a separate decision.

Visible in REQ-226's synthetic 60-sample fixture; invisible on this repository's own board, whose lane carries three samples.

## Requirements

- No first-row label in the overflow lane may share vertical space with the mark band, at any density.
- Panel A's scale break, overflow lane, `60+` tick, and two label rows all stay — this is a spacing change, not a redesign of the device.
- Panels B and C shift down by whatever Panel A grows by; `DURATIONS_MEDIAN_TITLE_Y` is the panel-A/B boundary `describeAtPointer` keys on, so the hover readout must still resolve the same panel for the same pointer position.
- The reversed band gets the same treatment or an explicit note saying why it does not need it.

## Red-Green Proof

**RED prompt/case:** A test asserting that the mark band (`DURATIONS_LANE_MARK_Y` ± the mark radius) and every label row's text box (baseline minus ascent, baseline plus descent) do not intersect, read from the renderer's own constants the way `TestDurationLabelGeometryMatchesTheRenderer` already reads them.
**Why RED now:** marks span roughly y 35-45 and row 0's text box roughly y 33-46, so the two intersect over about 10 units.
**GREEN when:** the same test passes, and re-rendering REQ-226's synthetic dense fixture shows the first-row label clear of the mark blob.
**Validation:** Discovered during REQ-226; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-226: in the board's Durations chart, the top strip that holds the very long REQs draws its first line of text at the same height as the dots themselves. Each dot's own label is offset sideways so it never covers its own dot, but where many long REQs finish close together, a *neighbouring* dot can sit on top of a label. REQ-226 stopped the labels covering each other; this is the leftover case of a label being covered by a dot. Fixing it means making that strip about 12 units taller and moving the two charts below it down to match — a small layout change, which is exactly what REQ-226 was told not to do ("Panel A's existing scale-break design is correct and stays"), so it is your call rather than mine. It shows up on a board like the one in your screenshot and not on this repository's own board, which has three long REQs in that strip. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  **Answered 2026-08-18 (in session):** Yes — and with a chosen design. The maintainer submitted a fresh screenshot of the consuming board (51 overflow samples) showing this defect at scale, asked for an audit and mockup alternatives before implementation, and picked **Alternative 2** from the proposal report at `ai-reports/2026-08-18_1000_durations-lane-legibility-alternatives/` (commit 0f8e349): band separation (marks keep the lane's top strip; both label rows move below a divider, with leader ticks) **plus** top-N-by-magnitude label selection, so the lane's visible labels are the longest spans rather than a left-to-right first-fit sample.

## Addendum (2026-08-18, UR-053)

User submitted a fresh screenshot of the consuming board's Durations view (51 overflow spans — this defect at real scale), asked for an audit and several mockup alternatives via the ai-report format before implementing, and chose **Alternative 2 — top-N extremes in the text band** from the proposal report at `ai-reports/2026-08-18_1000_durations-lane-legibility-alternatives/` (commit 0f8e349). This extends the original intent; nothing is contradicted:

- Band separation (the captured spacing fix): marks keep the lane's top strip, both label rows move below a divider, leader ticks tie each label to its mark.
- **Added requirement:** label selection changes from left-to-right first-fit to top-N-by-magnitude (N=6 per band), so every drawn label is one of the band's longest spans; unselected samples join the drawn remainder count.
- Verbatim input and the screenshot: `do-work/user-requests/UR-053/`.

## Decisions

- **D-01**: Scope expanded beyond the captured spacing fix to include top-N-by-magnitude label selection — the maintainer's explicit pick of Alternative 2 over the spacing-only Alternative 1. Selection stays on the Go side and ships in the existing `labelRow`/`labelAnchor` payload fields; no wire change. DECIDE & STATE (maintainer-directed).
- **D-02**: `web/board.css` added to the write set for the two new classes the mockup introduced (lane divider, label leader ticks). The captured write set predates the chosen design. DECIDE & STATE.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (Alternative 2, chosen by the maintainer in the addendum), but the change is geometric and spans two subsystems — the Go label planner and the JS renderer — so the panel constants, their consumers, and the geometry tests that read them had to be discovered before editing.

**Planning:** Not required

**Resume note:** This REQ was claimed and built by an earlier session that ended before writing any pipeline sections. This session found the implementation complete in the working tree, uncommitted, with `maintainer-verify.sh` green. The sections below record that state and re-verify it independently; `claimed_at` is left at the original claim instant so the calibration span measures the build, not the resume.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — top-N-by-magnitude label candidate selection
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — geometry and selection lock-in tests
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — band separation geometry, divider, leader ticks
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — divider and leader-tick classes

**Files I will NOT touch:** `timeline.go` / `web/board-timeline.js` (adjacent chart view, separate REQs), `model.go` (no schema change — selection ships in the existing `labelRow`/`labelAnchor` payload fields).

**Acceptance criteria (restated from REQ):**
- [ ] No first-row label in the overflow lane shares vertical space with the mark band, at any density
- [ ] Panel A's scale break, overflow lane, `60+` tick, and two label rows all survive — spacing change, not redesign
- [ ] Panels B and C shift down by Panel A's growth; `describeAtPointer` still resolves the same panel for the same pointer position
- [ ] The reversed band gets the same treatment, or an explicit note saying why not
- [ ] (Addendum) Label selection is top-N-by-magnitude, N=6 per band; unselected samples join the drawn remainder count

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

**What was done:** Panel A's overflow lane was split into a mark strip and a text band. `DURATIONS_LANE_MARK_Y` moved 40 → 32, a dashed `DURATIONS_LANE_DIVIDER_Y` at 44 separates the two, and both label rows moved below it (`DURATIONS_LANE_LABEL_ROW_Y` 44 → 56), so no label's text box can intersect the mark band at any density. The band mark radius became a named constant (`DURATIONS_BAND_MARK_RADIUS`) and the label face's ascent a second (`DURATIONS_LABEL_TEXT_ASCENT`) so the new geometry test can assert the separation from the renderer's own numbers rather than from copies. A vertical leader tick now ties each label to its mark across the gap, stopping at the text band's top edge so it cannot cross a first-row label on the way to a second-row one. The reversed band got the same treatment (`DURATIONS_REVERSED_LABEL_ROW_Y` 288 → 322 against `DURATIONS_BELOW_ZERO_Y` 284 → 298), and Panels B and C shifted down by Panel A's growth (view height 570 → 604), with `describeAtPointer`'s panel boundary still derived from `DURATIONS_MEDIAN_TITLE_Y` so it moves with them.

On the Go side, `selectDurationLabelCandidates` narrows each band to its `durationsLabelTopCount` (6) longest spans by `math.Abs(WallMinutes)` before the existing greedy left-to-right packer runs, using a stable sort so equal spans keep completion order. Unselected samples increment `HiddenCount` and are carried by the remainder count, the hover readout, and the table. Selection stays server-side and ships in the existing `labelRow`/`labelAnchor` payload fields — no wire change.

## Qualification

Passed — 4 files verified in the diff, 5 acceptance criteria traced, `tools/checks/qualify.sh` clean. No P-A-U section exists on this REQ (captured before the block was standard), so the box audit had nothing to check; the diff carries no debug artifacts.

Judgment checks: all four files are substantive modifications, not placeholders. Requirement 3 traced by reading `describeAtPointer` — its panel boundary is `DURATIONS_MEDIAN_TITLE_Y - 12`, derived from the constant that moved, so the boundary travelled with the panels rather than being left behind at a stale literal. Requirement 4 traced to the reversed band's own constants and to the geometry test's second table row. Nothing hollow: the selection path returns real index sets and the renderer draws from them.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 — full maintainer baseline including `go vet`, the uncached queue-kanban suite, and the strict JavaScript behavior lane

**Red-green validation:**
- `TestDurationsLabelRowsClearTheMarkBands`: ✗ before → ✓ after. RED reproduced by restoring the pre-fix renderer geometry with only the two constants the test reads by name exposed, so the failure is geometric rather than a compile error: `overflow band: row 0's text box [33, 46] intersects the mark band [35, 45] — a neighbouring dot can overprint the label`. That is the REQ's captured "Why RED now" verbatim (marks ~35-45, row 0 box ~33-46).
- `TestOverflowLabelsGoToTheLongestSpans`: ✗ before → ✓ after. RED reproduced by keeping the full geometry fix and neutering only the candidate-selection skip, so the packer fell back to pure first-fit: `REQ-500 (65 min) carries a label but is not among the 6 longest spans (floor 303 min)`. The geometry test passed during this run, which isolates the two failures to their own causes.

**New tests added:**
- `TestDurationsLabelRowsClearTheMarkBands` — both bands' label rows against their band's mark extent, read from the renderer's own constants
- `TestOverflowLabelsGoToTheLongestSpans` — top-N selection, plus the invariant that labelled + hidden equals the sample count so unselected samples cannot be silently dropped
- `variedOverflowTickets` — dense fixture with a magnitude gradient, built so a first-fit walk provably spends both rows on the shortest spans

**Render evidence (REQ-226's lesson applied — a passing suite is not evidence about pixels):** generated a 105-REQ synthetic board (60 overflow samples, 40 sub-ceiling, 5 reversed) with the pre-fix and post-fix binaries and measured `getBoundingClientRect()` intersections between `.durations-mark-label` and `.durations-mark` in the live DOM. **Before: 27 labels, 55 overlapping label/mark pairs. After: 8 labels, 0 overlapping pairs.** The before render shows the whole first label row struck through by dots; the after render shows the dashed divider with marks above it, labels and leader ticks below, and nothing touching.

*Verified by work action*

## Discovered Tasks

- [normal] **Dropped top-N candidates are not backfilled.** Selection takes the 6 longest spans per band, then the collision packer drops whichever of them do not fit. Nothing promotes the 7th-longest into a row a dropped candidate freed, so where the longest spans cluster in time the two label rows go under-filled. Measured on the magnitude-gradient fixture (every long span crowded at the right edge): 2 labels drawn out of 6 candidates. On a scattered-magnitude fixture of the same size, 5 of 6 place, which is why this is a tail case rather than the common one — but the tail is exactly the shape a burst of long REQs produces.

## Review

**Overall: 97%** | 2026-08-18T10:50:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 2 (report only)
- *Restatement sweep hit.* Three comments described the remainder count as "whatever could not fit" / "the collision rule could not place" — accurate before top-N selection, stale after it, since a sample can now be unlabelled for a second reason. `durationsLabelRowUnplaced`, `DurationLabelBand`, and the JS remainder comment were corrected during review; all three sit inside the declared write set, so this completes the diff rather than widening it. The sweep also cleared `generate_test.go`'s JS probe (it reads `DURATIONS_LANE_LABEL_ROW_Y` from the renderer, so it travelled with the constant), the archived REQ-219/REQ-226 records (dated history, correctly left alone), and `CHANGELOG.md` (release history, out of bounds by rule).
- A second-row label's leader tick stops at the *first* row's text top, so it visually falls a row short of the label it serves. Deliberate — the alternative is a tick crossing a first-row label on its way down — and the code comment says so.

**Acceptance:** Pass — verified by DOM measurement on a 105-REQ synthetic board: 55 overlapping label/mark pairs before, 0 after, with the pre-fix render visibly striking its whole first label row through with dots.

**Suggested testing:** 2 items
- Eyes on the real consuming board (the 51-sample screenshot that prompted the addendum) at the maintainer's own zoom and theme — the divider is `--line-soft`, faint by design, and only a human can say whether it reads as a separator or as noise.
- Light theme was the only one rendered this session. The divider and leader-tick strokes are token-driven so dark theme should follow, but nothing measured it.

**Follow-ups created:** REQ-237; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reproducing RED by restoring the pre-fix geometry while exposing only the two constants the new test reads by name. A straight revert gave `undefined: durationsLabelTopCount` — a compile error, which proves nothing about geometry. Exposing the names and leaving the positions wrong produced `row 0's text box [33, 46] intersects the mark band [35, 45]`, which is the defect stated in the test's own words. The same trick isolated the second test: keep the whole geometry fix, neuter only the selection skip, and exactly one of the two tests fails.

**What didn't:**
- Trusting the label *count* as a legibility proxy. The dense fixture drew 2 labels where the pre-fix render drew 27, which looked like a regression until a second fixture with scattered magnitudes drew 5 of 6. The first fixture correlated magnitude with x-position perfectly, so every top-N candidate landed in the same crowded corner — an adversarial input for top-N, not a defect in it. One fixture cannot tell a design's tail case from its common case.
- Reading the geometry test as sufficient. It reads renderer constants, so it proves the *declared* rows clear the *declared* band; it cannot see a leader tick crossing a glyph or a divider drawn in a colour nobody can see. Measuring `getBoundingClientRect()` intersections in the live DOM is the assertion that actually answers "do any two things touch" — 55 pairs before, 0 after, from the rendered document rather than from the source.

**Worth knowing:** Two reasons a sample can now go unlabelled — selection passed it over, or placement could not fit it — where there used to be one. Three comments still said "could not fit", and the restatement sweep is what caught them; the whole suite passed either way, because no test asserts on prose. When a change adds a *second* cause for an existing outcome, every sentence naming the first cause is now a half-truth.

## Orientation

The board's Durations view separates Panel A's overflow lane into a mark strip and a text band: dots sit above a dashed divider, both label rows below it with leader ticks tying each label to its dot, so no dot can overprint a label at any density. Which spans get labelled also changed — each band's six longest, rather than whichever the left-to-right packer reached first — so the lane's scarce text now answers "where are the outliers" instead of sampling the left edge. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), split as that view already was: selection in `durations.go`, geometry in `web/board-durations.js`.

Not `[MAP CHANGED]` — the payload contract is untouched (selection ships in the existing `labelRow`/`labelAnchor` fields), and the Go-owns-the-rule split is REQ-219's shape, applied rather than altered. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves, and the three-write-surface count is unchanged — this REQ adds none. The prime is not stale.
