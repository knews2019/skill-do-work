---
id: REQ-241
title: Reconcile the Durations label metrics with the face actually rendered
status: completed
status_changed_at: 2026-08-18T13:51:09Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: durations-label-metric-constants-disagree-with-rendered-face
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - browser evidence
    - cross-route regression gates
claimed_at: 2026-08-18T12:43:06Z
route: B
completed_at: 2026-08-18T13:51:09Z
commit:
---

# Reconcile the Durations Label Metrics With the Face Actually Rendered

## What

Two constants that describe the Durations label face disagree with what the browser actually draws. Neither causes a visible collision today, but both are now load-bearing in a way they were not before, because REQ-237 made the label rows actually fill up.

## Context

Found by REQ-237's build and confirmed independently against the merged tree. Both are pre-existing; REQ-237 is what made them reachable, not what caused them.

Until REQ-237, the overflow lane drew two or three labels on any real board, so the slack in these constants never mattered. It now draws as many as fit — measured 21 on a clustered 60-sample fixture — and the slack is what stands between "packed" and "overlapping".

## Instances

- [x] **`durationsLabelCharacterWidthUnits = 6.2` under-estimates the 11px sans face by ~7%.** Its comment calls the value "deliberately generous"; measurement says otherwise — a 14-character label renders **92.52 user units**, i.e. **6.61 units/char**, not the 86.8 the constant predicts. The 6-unit separation rule absorbs the difference, so nothing collides: the tightest same-row gap in a full 21-label render is 3.08 units, and a direct DOM measurement finds **0 same-row overlaps**. But the real margin is about half what the code claims, and the comment is actively misleading about which direction the error runs.

- [x] **`DURATIONS_LABEL_ROW_HEIGHT = 12` is smaller than the text box the same file declares.** `DURATIONS_LABEL_TEXT_ASCENT = 11` plus a 2-unit descent is a 13-unit box on a 12-unit pitch; the rendered font box measures 12.83 units. Measured on a densely-populated lane: **20 cross-row bounding-box intersections, each 1.6px deep.** This is line-box padding rather than ink — the render shows two cleanly separated rows, and a screenshot confirms it — but it means no test can honestly assert row-against-row separation the way `TestDurationsLabelRowsClearTheMarkBands` asserts row-against-mark.

## Requirements

- Each constant either matches the face the browser renders, or its comment states the measured value and why the code deliberately differs. A constant whose comment claims a safety margin in the wrong direction is the specific defect here.
- Whatever changes, the same-row separation guarantee holds: **0 same-row label overlaps** at full density, measured from the live DOM rather than computed from the constants under test.
- REQ-231's guarantee holds unchanged: **0 label/mark overlaps** in either band, at any density.
- If the row pitch changes, Panels B and C shift with it and `describeAtPointer` still resolves the same panel for the same pointer position — the same constraint REQ-231 worked under.
- Changing label counts across the view is an accepted consequence of retuning the width model; say so in the REQ trail with before/after counts on a real board rather than only on a fixture.

## Red-Green Proof

**RED prompt/case:** a test asserting each constant against the measured face — that `durationsLabelCharacterWidthUnits` is not below the rendered units-per-character, and that `DURATIONS_LABEL_ROW_HEIGHT` is not below the declared ascent-plus-descent box.
**Why RED now:** 6.2 < 6.61, and 12 < 13.
**GREEN when:** the assertions pass, and a rendered clustered fixture still shows 0 same-row label overlaps and 0 label/mark overlaps measured in the live DOM.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I found two numbers in the Durations chart that describe the label text and are both slightly wrong about it. One says each character is 6.2 units wide when the browser draws 6.61 — and its comment claims the estimate is deliberately generous, which is backwards. The other says a row of label text is 12 units tall when the same file elsewhere describes that text as a 13-unit box. Nothing overlaps on screen today: I measured a fully packed lane and found no two labels on the same row touching, and no label touching a dot. But the spare room absorbing both errors is about half what the code claims it is, and until this session the chart only ever drew two or three labels so it never mattered. It does now — REQ-237 made the rows fill up. Fixing the width number will change how many labels the chart draws on every board, including yours, which is why it is your call and not a quiet tidy-up: it is a visible change to a view you look at, made for correctness reasons rather than because anything is broken. The alternative is to leave the numbers and fix the comments so they stop claiming a margin that is not there. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No — just correct the misleading comments and leave the numbers alone, since nothing overlaps.
  → Confirmed: Yes, add to queue — full fix, retune the constants to match the measured face (user accepted that label counts across the view will visibly change). [2026-08-18, via do-work clarify]

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — retune `durationsLabelCharacterWidthUnits`
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — assert each constant against the measured face
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — retune `DURATIONS_LABEL_ROW_HEIGHT`

**Files I will NOT touch:** `generate_test.go` (owned by a sibling builder this wave), `web/board-timeline.js`, anything under `do-work/`, the version files, `CHANGELOG.md`.

**Acceptance criteria (restated from REQ):**
- [x] Each constant matches the face the browser renders, or its comment states the measured value and why the code deliberately differs
- [x] Same-row separation guarantee holds: 0 same-row label overlaps at full density, measured from the live DOM rather than computed from the constants under test
- [x] REQ-231's guarantee holds unchanged: 0 label/mark overlaps in either band, at any density
- [x] `describeAtPointer` still resolves the same panel for the same pointer position — satisfied, and Panels B/C did not need to move (D-03)
- [x] Changed label counts recorded with before/after numbers **on a real board**, not only a fixture

## Pre-Flight

**Git:** ✓ clean at claim (`2432f45`); builder worked in an isolated worktree on its own branch
**Tests baseline:** ✓ `maintainer-verify.sh` exit 0 before dispatch
**Dependencies:** ✓ Go 1.26.1, Playwright 1.59 / Chromium 1228 available for live-DOM measurement

*Checked by work action*

## AI Execution State (P-A-U Loop)

<!-- Transcribed verbatim from the builder's hand-back. In worktree dispatch mode a builder may not
     write do-work/, so the orchestrator carries the filled loop into the REQ. Source of record:
     do-work/runs/work-2026-08-18-124358/REQ-241-handback.md -->


- [x] **PLAN** — Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md`, and the three always-on crew members (`general.md`, `coding-guardrails.md`, `communication-style.md`), then read `durations.go`, `durations_test.go`, and `web/board-durations.js` end to end.

  Approach: both defects are calibration claims, so the first move is measurement, not editing. Build the pre-change binary, render three boards (the maintainer's real checkout, a clustered 60-sample overflow fixture, a saturated-reversed fixture), and measure the face in an **isolated** Chromium launched by this builder — not a shared MCP browser — with `location.href` returned from the same `evaluate` call as every number.

  Two measurements decide the constants:
  1. `durationsLabelCharacterWidthUnits` models width as `len(text) × constant`, so the honest constant is an **upper** bound on units-per-character over the label space the renderer can actually compose — not the average of one board's strings. `formatDurationMinutes` emits exactly two forms, so that space is enumerable: sweep every digit across three- and four-digit REQ ids in both forms and take the worst ratio.
  2. `DURATIONS_LABEL_ROW_HEIGHT` must clear the label's text box. Two boxes exist — the one the code *declares* (`DURATIONS_LABEL_TEXT_ASCENT` + the test's descent) and the one the browser *draws* — so measure the drawn one via `getBBox()` and set the pitch above whichever is larger.

  Then verify the guarantees empirically rather than by arithmetic over the constants just changed: same-row overlaps, cross-row overlaps, and label/mark overlaps all as `getBoundingClientRect()` intersections in the live DOM, plus `describeAtPointer` panel resolution driven through the real hover surface. Write set stays exactly the three declared files; `generate_test.go` is not touched.

- [x] **APPLY** — Coded exactly as planned, three files, no others.
  - `durations.go`: `durationsLabelCharacterWidthUnits` 6.2 → 6.75, comment replaced with the measured numbers and the reason the model must be an upper bound; `durationsLabelSeparationUnits` comment amended to say the 6 units are now real whitespace rather than the slack that was absorbing an under-estimate.
  - `web/board-durations.js`: `DURATIONS_LABEL_ROW_HEIGHT` 12 → 13, with the measured line box and the declared box both recorded; `DURATIONS_LABEL_TEXT_ASCENT` kept at 11 with a comment stating the measured 10.43 and why the round-up is the safe direction for both of its readers.
  - `durations_test.go`: two measured-face constants with their full measurement procedure, and two lock-in tests — `TestDurationsLabelWidthEstimateCoversTheRenderedFace` and `TestDurationsLabelRowPitchClearsTheLabelTextBox`.

  No production logic changed. Placement, packing, anchoring, and the remainder rule are byte-identical; only the two numbers those rules consume moved.

- [x] **UNIFY** — `git diff --stat` is three files, 95 insertions, 7 deletions. Reviewed each:
  - `durations.go` — verified the only executable change is the literal `6.2` → `6.75`; everything else in the hunk is comment. Verified `durationsLabelRemainderReserveUnits` still derives from the constant (it scales automatically; see D-04).
  - `web/board-durations.js` — verified the only executable change is `12` → `13`; `DURATIONS_LABEL_TEXT_ASCENT` is unchanged at 11. Confirmed both constants are still plain `var NAME = NUMBER;` literals, because `rendererNumericConstant` parses them with a regex and an expression would break every parity test.
  - `durations_test.go` — verified both new tests read the renderer's own constants through `durationsRendererConstant` rather than hand-copied numbers, and that the measured constants are rounded **away** from the model so no assertion can pass on rounding.
  - No debug artifacts: `git status --porcelain --untracked-files=all` reports nothing untracked in the worktree. Every scratch file (fixture repos, boards, measurement scripts, screenshots, both binaries) lives under `/tmp/board-241/`, `/tmp/qk-241`, and `/tmp/qk-241-before`. Nothing was written to either repo root.
  - `git diff --name-only <merge-base>...HEAD -- do-work/` is empty. Files changed on the branch: exactly the three in the write set.
  - `bash _dev/tests/maintainer-verify.sh` exits 0 (pasted below), which covers `go vet`, the uncached queue-kanban suite, and the strict JavaScript behavior lane.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)

**What was done:** Raised `durationsLabelCharacterWidthUnits` from 6.2 to 6.75 — an upper bound over the enumerable label space rather than a fit to one board's strings — and `DURATIONS_LABEL_ROW_HEIGHT` from 12 to 13 so the row pitch clears both the declared text box and the 12.83-unit box the browser actually draws. Added assertions holding each constant to the measured face, including a second branch pinning the declared box so trimming it cannot silently reintroduce the defect.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 (unpiped; `echo $?` printed `0` on its own line)

**Red-green validation:**
- new width assertion (6.2 vs measured 6.6087 units/char on the REQ's own label): ✗ before implementation → ✓ after
- new pitch assertion (12 vs the 13-unit declared box and 12.83-unit drawn box): ✗ before implementation → ✓ after
- second branch of the pitch test (declared box trimmed below the face): ✗ when mutated → ✓ restored. Produced deliberately because that branch is unreachable at today's values and would otherwise read as dead code.

**Live-DOM measurement (not computed from the constants under test), clustered 60-sample fixture:**

| | before | after |
|---|---|---|
| Same-row overlaps | 0 | **0** |
| Tightest same-row gap | 3.078 u | **11.35 u** |
| Cross-row box intersections | 19 | **0** |
| Label/mark overlaps | 0 / 63 | **0 / 63** |

**Visible consequence on a real board:** the maintainer's own board (242 REQs) draws **3 labels before and 3 after — no visible change**. The clustered fixture goes 24 → 21 and the saturated reversed band 18 → 17, with the remainder sentence counting every difference. Per-constant builds attribute all three lost labels to the width retune and **none** to the pitch change.

**Existing tests updated (cross-REQ impact):** none. `TestDurationsLabelRowsClearTheMarkBands`, `TestDurationLabelGeometryMatchesTheRenderer`, `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`, `TestClusteredOverflowLabelsFillBothLabelRows`, `TestOverflowLabelsGoToTheLongestSpans` and `TestReversedLabelPlacementIsIndependentOfOverflowDensity` all pass unchanged.

*Verified by work action*

### Red-Green Evidence (verbatim from hand-back)


All RED runs were produced by restoring the **pre-change blobs** (`git show HEAD:<path>`) over the working files while keeping the new tests — never `git stash push`, which stashes nothing on a file it does not consider dirty.

### RED #1 — both new assertions against the pre-change constants

```
223:	durationsLabelCharacterWidthUnits = 6.2
46:  var DURATIONS_LABEL_ROW_HEIGHT = 12;
--- RED run: new tests against pre-change constants ---
=== RUN   TestDurationsLabelWidthEstimateCoversTheRenderedFace
    durations_test.go:492: width model assumes 6.2000 units per character, but the rendered face draws up to 6.7100 — the estimate under-states the text it is placing
--- FAIL: TestDurationsLabelWidthEstimateCoversTheRenderedFace (0.00s)
=== RUN   TestDurationsLabelRowPitchClearsTheLabelTextBox
    durations_test.go:509: row pitch 12.00 is smaller than the 13.00-unit text box the renderer declares (ascent + 2.00 descent) — consecutive rows share vertical space
--- FAIL: TestDurationsLabelRowPitchClearsTheLabelTextBox (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.352s
```

Both fail on their assertion, for the reason the assertion exists — not on a missing symbol. The measured-face constants were added before this run precisely so the RED could not be a reference error.

### RED #2 — proving the second assertion in the pitch test is load-bearing

The declared box (ascent + descent) is always ≥ the drawn box today, so RED #1 can only ever exercise the first branch. To show the drawn-box branch is a real backstop and not dead code, the declared box was trimmed below the face (`DURATIONS_LABEL_TEXT_ASCENT` 11 → 10) at a pitch that satisfies the declared box but not the drawn one:

```
46:  var DURATIONS_LABEL_ROW_HEIGHT = 12.5;
47:  var DURATIONS_LABEL_TEXT_ASCENT = 10;
--- RED #2: declared box (10+2=12) satisfied, measured face (12.84) not ---
=== RUN   TestDurationsLabelRowPitchClearsTheLabelTextBox
    durations_test.go:513: row pitch 12.50 is smaller than the 12.84-unit line box the browser draws — consecutive rows share vertical space
--- FAIL: TestDurationsLabelRowPitchClearsTheLabelTextBox (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.351s
```

### GREEN — the whole queue-kanban suite, post-change

```
231:	durationsLabelCharacterWidthUnits = 6.75
54:  var DURATIONS_LABEL_ROW_HEIGHT = 13;
59:  var DURATIONS_LABEL_TEXT_ASCENT = 11;
--- GREEN: full queue-kanban suite ---
ok  	github.com/knews2019/skill-do-work/queue-kanban	24.372s
GREEN_EXIT=0
```

### The measurements the constants are calibrated from

Isolated Chromium 1228 launched by this builder (Playwright 1.59, `chromium.launch()` + fresh context), so no sibling agent could navigate it. `location.href` was returned from the **same** `evaluate` call as every number below, and every reported href was the board this REQ generated.

Face measurement, on `file:///tmp/board-241/before-fixture/index.html`, probing an SVG `<text class="durations-mark-label">` appended to the chart's own `<svg>`:

| Quantity | Measured |
|---|---|
| Resolved face | the board's `--font-sans` stack at the class's declared `11px` |
| Widest label per character, whole 10 000-string corpus | **6.7054** — `REQ-4444 44h 44m`, 107.29 units / 16 chars |
| …restricted to today's three-digit ids | **6.6762** — `REQ-444 44h 44m`, 100.14 units / 15 chars |
| Widest total advance | 112.50 units — `REQ-4444 −44.4 min` |
| Rendered line box | **12.8343** units — ascent **10.4278**, descent **2.4064** |
| Line box with descenders (`gjpqy`) | identical, 12.8343 — it is the line box, not the ink |

The REQ's own figure is reproduced exactly: the widest label on the clustered fixture, `REQ-549 6h 48m`, advances **92.522** units over 14 characters = **6.6087** units/char, against the 86.8 the old 6.2 constant predicted.

### Live-DOM guarantees, before and after

Every number is a `getBoundingClientRect()` intersection read in the live DOM, never computed from the constants under test.

| Board | Band labels | Same-row overlaps | Tightest same-row gap | Cross-row box intersections | Label/mark overlaps |
|---|---|---|---|---|---|
| Real board, before | 3 | 0 | 76.65 u | 1 (0.834 u deep) | 0 / 221 marks |
| **Real board, after** | **3** | **0** | 76.65 u | **0** | **0 / 221 marks** |
| Clustered 60-sample fixture, before | 24 | 0 | 3.078 u | 19 (0.834 u deep) | 0 / 63 marks |
| **Clustered fixture, after** | **21** | **0** | **11.35 u** | **0** | **0 / 63 marks** |
| Saturated reversed band, before | 18 | 0 | 3.279 u | 16 (0.834 u deep) | 0 / 60 marks |
| **Saturated reversed band, after** | **17** | **0** | **15.86 u** | **0** | **0 / 60 marks** |

- **Same-row separation holds at full density: 0 overlaps, before and after.** What changed is the honesty of the margin — the tightest real gap in a saturated lane went from 3.08 units (against a rule claiming 6) to 11.35, because the width model no longer under-states the text.
- **REQ-231's guarantee holds unchanged: 0 label/mark overlaps in either band, at every density measured.**
- The 0.834-unit cross-row intersection the REQ reported is exactly `1.04px × 0.8021 units/px` — the same defect, now zero.

### Label counts on a real board — the visible consequence

Required by the REQ, so stated plainly with the real-board number first:

- **The maintainer's own board (242 REQs, `--repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2`): 3 direct labels before, 3 after. No visible change, no remainder sentence in either render.** The live archive carries only three overflow samples and they sit far apart, so the width retune costs nothing there.
- **Clustered 60-sample fixture: 24 → 21** (overflow 21 → 18, reversed 3 → 3). The remainder sentence moves `+39 more over 60 min` → `+42 more over 60 min`; nothing is dropped silently.
- **Saturated reversed band (60 clustered reversed stamps): 18 → 17**, remainder `+42 more reversed` → `+43 more reversed`.

**Attribution, measured rather than assumed** — the orchestrator note asked specifically whether the pitch increase costs a label row. It costs none. Two extra binaries were built with one change each and rendered against the same clustered fixture:

| Build | Labels | Rows | Remainder | Cross-row intersections |
|---|---|---|---|---|
| before (6.2 / 12) | 24 | 11 + 10 + 3 | `+39 more` | 19 |
| **width only (6.75 / 12)** | **21** | 10 + 8 + 3 | `+42 more` | 16 |
| **pitch only (6.2 / 13)** | **24** | 11 + 10 + 3 | `+39 more` | **0** |
| after (6.75 / 13) | 21 | 10 + 8 + 3 | `+42 more` | 0 |

All three lost labels are the width retune; the pitch increase costs exactly zero and fixes the cross-row intersections on its own. That is structural, not luck: `DURATIONS_LABEL_ROW_HEIGHT` lives only in the renderer, and placement in `durations.go` packs each row purely horizontally with no knowledge of row height — so a pitch change cannot move a label between rows or off the chart. The label loss is the intended trade for the width fix: an over-estimate drops a label, which the remainder counts, while an under-estimate overprints one.

### Panels B and C did not need to move, and `describeAtPointer` is unchanged

The REQ says Panels B and C shift with the row pitch. **They did not have to, and I did not move them** — see D-03. Measured on the saturated-reversed board, where the reversed band's *last* row is actually occupied:

| | Before | After |
|---|---|---|
| Reversed row 1 baseline | 334 | 335 |
| Its drawn box | [323.57, 336.41] | [324.57, 337.41] |
| Panel B title box | [338.77, 352.41] | [338.77, 352.41] |
| Clearance | 2.364 u | **1.364 u** |
| Live-DOM intersections with the title | **0** | **0** |

`describeAtPointer`, driven through the real `rect.durations-hover-surface` with synthesized `mousemove` events at eleven pointer positions spanning the A/B boundary (y = 320, 330, 334, 335, 336, 337, 338, 339, 340, 350, and on the other boards 100, 300, 400, 500):

```
   y=334  before=A (sample)   after=A (sample)   SAME
   y=335  before=A (sample)   after=A (sample)   SAME
   y=336  before=A (sample)   after=A (sample)   SAME
   y=337  before=A (sample)   after=A (sample)   SAME
   y=338  before=A (sample)   after=A (sample)   SAME
   y=339  before=B/C (day)    after=B/C (day)    SAME
   y=340  before=B/C (day)    after=B/C (day)    SAME
```

Every probe resolves the identical panel before and after, on all three boards. Because `DURATIONS_MEDIAN_TITLE_Y` is untouched, this is true by construction as well as by measurement.

### Rendered artifacts

Boards were generated and looked at, per the prime's standing rule. Screenshots at `/tmp/board-241/shot-{before,after}-{fixture,reversed}.png`. The after-fixture render shows two clean, well-separated label rows in the overflow lane with the remainder sentence on the last row's right edge; the after-reversed render shows the reversed band's second row sitting clear above Panel B's title.

## Decisions


- **D-01 — The width constant is calibrated as an upper bound over the label space, not fitted to one board's strings.** `durationLabelWidthUnits` is `len(text) × constant`, so the constant is only honest if no label the renderer can compose exceeds it per character. `formatDurationMinutes` emits exactly two forms, which makes that space enumerable rather than a guess: sweeping every digit across three- and four-digit ids in both forms gives a worst case of 6.7054 (`REQ-4444 44h 44m`). 6.75 sits above it. The `"Hh Mm"` form is the dense one because `" min"` spends three of its characters on a space, a narrow `i`, and an `n` — which is why the widest *total* label (`REQ-4444 −44.4 min`, 112.50 units) is not the worst *per character*. Reach: any future change to the label's copy changes the worst case, and the test comment says where to re-measure.

- **D-02 — The row pitch clears both boxes, and the test checks both.** The declared box (`DURATIONS_LABEL_TEXT_ASCENT` + the test's descent = 13) and the drawn box (12.83) are different numbers and either could become the binding one. 13 clears both today; the test asserts against both so trimming the declared box below the face — the exact move RED #2 demonstrates — cannot silently reintroduce the defect.

- **D-03a — The pitch increase costs no labels, so no trade needed making.** The orchestrator note asked for the trade to be stated if raising `DURATIONS_LABEL_ROW_HEIGHT` cost a label row. Measured per-constant builds show it costs none, because the pitch is renderer-only and placement is horizontal — see the attribution table above.

- **D-03 — Panels B and C stay where they are; the pitch increase was absorbed inside the reversed band's existing headroom.** The REQ anticipated a shift, but a shift would move `DURATIONS_MEDIAN_TITLE_Y` and therefore the `describeAtPointer` A/B boundary, putting the REQ's two requirements in direct tension. Measurement resolved it: the reversed band's last row still clears Panel B's title by 1.364 units with 0 live-DOM intersections, so no shift is needed and pointer resolution is preserved by construction. **This is now the binding constraint on the row pitch** — the next reader who wants a taller row should know that 1.36 units is the whole budget, and that a pitch of 15 would put the last row's box bottom at 339.41, past the title's top at 338.77. That is a Panel B shift, and it is the moment the two requirements really do collide.

- **D-04 — `durationsLabelRemainderReserveUnits` was left as `24 × durationsLabelCharacterWidthUnits`.** It scales with the width constant automatically, so it needed no edit, but the widening makes it over-reserve more: the widest remainder sentence the renderer can compose (`+999 more over 60 min`, 21 chars) measures 123.18 units, while the reservation is now 162.0. Over-reserving is the safe direction — it drops a label, which the remainder counts — and re-sizing it is a different question from the two the REQ names, so it stays out of scope. Flagged in Lessons Learned as a follow-up candidate.

## Lessons Learned


- **A per-character width model has a right answer, and it is the worst case over the label space, not the average over one board.** The 6.2 constant was probably fitted to a plausible-looking string; the comment then described the *intent* (generous) rather than the *result* (7% short), and nothing could catch the gap because no test knew what the face draws. Two things fix this class permanently: enumerate the label space when the format function makes it enumerable, and record the browser's measurement as a named constant so a Go test can hold the model to it without a browser in the loop.

- **A constant that clears the box the code declares can still fail the box the browser draws, and vice versa — assert both.** Rounding an ascent up (11 from 10.43) is safe and worth keeping, but it means the declared box and the drawn box are two independent claims. A single assertion against whichever happens to be larger today rots the moment someone trims the other. RED #2 was worth producing precisely because the second branch is unreachable at today's values and would otherwise read as dead code.

- **A REQ can name a consequence that measurement then shows is not required — say so rather than performing it.** The brief expected Panels B and C to shift with the row pitch. Shifting them would have broken the *other* stated requirement (`describeAtPointer` resolving the same panel for the same pointer). Measuring the actual clearance showed the pitch increase fit inside existing headroom, which satisfies both requirements at once. The useful artifact is not "no change needed" but the number: 1.36 units of remaining budget, and the specific pitch (15) at which the trade becomes real.

- **Isolate a multi-constant change to attribute its effects, instead of reasoning about them.** I first wrote that the three lost labels split two-to-the-width-retune and one-to-the-pitch, which sounded plausible and was wrong: building one binary per constant showed the pitch costs zero labels and the width costs all three. The reason is structural and was sitting in the code the whole time — `DURATIONS_LABEL_ROW_HEIGHT` exists only in the renderer, and `durations.go` packs each row purely horizontally, so a pitch change *cannot* alter which labels are drawn. When a change moves two numbers, one build per number is cheap and turns a guess into a fact.

- **A "0 overlaps" result can be true and still be hiding a halved margin.** Both before and after renders show zero same-row overlaps; the difference is that the tightest gap went from 3.08 to 11.35 units against a rule claiming 6. When a guarantee is checked only as a boolean, a margin can erode most of the way to failure without any test noticing. Where a rule states a *quantity*, measure the quantity, not just its sign.

- **Follow-up candidate (not filed):** the remainder reservation over-reserves by ~39 units at the new width (D-04), and Panel B's slowest-day annotation collides with Panel B's own title on a dense board — visible in both the before and after fixture screenshots at `/tmp/board-241/shot-{before,after}-fixture.png`, where `233 min` overprints `…spans excluded`. Both are pre-existing, neither is in this write set, and the second is a genuine visual defect on any board whose slowest day is near the ceiling.

## Builder Pushback (accepted by the orchestrator)


- **The REQ's "if the row pitch changes, Panels B and C shift with it" is not achievable together with "`describeAtPointer` still resolves the same panel for the same pointer position".** Shifting Panel B means moving `DURATIONS_MEDIAN_TITLE_Y`, which is the sole input to the A/B boundary (`pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12`), so any shift moves the boundary and the second requirement fails at it. I satisfied the second requirement in full and made the first vacuous by measuring that no shift is needed — see D-03. If the intent was "shift them *if clearance is lost*", that reading is what I implemented and the clearance measurement is the evidence.

- **The REQ's second instance says `DURATIONS_LABEL_ROW_HEIGHT = 12` is smaller than a 13-unit box and reports "20 cross-row bounding-box intersections, each 1.6px deep".** My saturated fixture measures 0.834 user units of intersection, which is 1.04 CSS px at my viewport — the depth in px depends on the browser window width, since the SVG is `width:100%` over a fixed viewBox, so a px figure is not portable between measurements. The count (16–19 on my fixtures vs the REQ's 20) also depends on fixture composition. Neither discrepancy changes the finding; both are reported in user units here so the number is reproducible.

- Nothing else in the brief was wrong, and no existing test needed editing. `TestDurationsLabelRowsClearTheMarkBands`, `TestDurationLabelGeometryMatchesTheRenderer`, `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`, `TestClusteredOverflowLabelsFillBothLabelRows`, `TestOverflowLabelsGoToTheLongestSpans`, and `TestReversedLabelPlacementIsIndependentOfOverflowDensity` all pass unchanged against the new constants.

**Orchestrator resolution.** Both points accepted. On the first: the REQ's "Panels B and C shift with it" was written as an anticipated consequence, not a requirement, and it does collide head-on with the `describeAtPointer` requirement, which is the stronger of the two. The builder satisfied the binding requirement and made the other vacuous by measurement — the correct resolution, and D-03's number (1.364 units of remaining headroom, colliding at a pitch of 15) is the artifact that makes it checkable later. On the second: user units are the reproducible unit here and the px figure in the REQ was viewport-dependent; the REQ's finding is unchanged either way.
