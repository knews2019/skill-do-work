# REQ-241 — Reconcile the Durations Label Metrics With the Face Actually Rendered

## Branch

`worktree-agent-REQ-241-reconcile-durations-label-metrics-with-the-rendered-face`

Implementation commit: `4950c81`

## P-A-U

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

## Files Changed

```
 skills/do-work-board/tools/queue-kanban/durations.go            | 24 ++++++--
 skills/do-work-board/tools/queue-kanban/durations_test.go       | 64 ++++++++++++++++++++++
 skills/do-work-board/tools/queue-kanban/web/board-durations.js  | 14 ++++-
 3 files changed, 95 insertions(+), 7 deletions(-)
```

- **`durations.go`** — `durationsLabelCharacterWidthUnits` raised 6.2 → 6.75 so the per-character width model is an upper bound on the face rather than a 7% under-estimate, and its comment now records the measurement instead of claiming a generosity that ran the other way. The separation constant's comment gained the measured evidence that it, not the width model, was holding the labels apart.
- **`web/board-durations.js`** — `DURATIONS_LABEL_ROW_HEIGHT` raised 12 → 13 so the row pitch clears both the declared 13-unit text box and the measured 12.83-unit line box. `DURATIONS_LABEL_TEXT_ASCENT` stays 11; its new comment states the measured 10.43 and why the deliberate over-statement is safe in both places it is read.
- **`durations_test.go`** — added `durationsMeasuredWidestUnitsPerCharacter` (6.71) and `durationsMeasuredLabelBoxHeightUnits` (12.84) with the reproducible measurement procedure, plus the two tests that hold each constant to its measured value. The row-pitch test checks the declared box *and* the drawn box, because the declared one can itself be trimmed below the face (RED #2 below shows that path is live, not decorative).

## Red-Green Evidence

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

## Verification

```
maintainer-verify: checking Go go1.26.1
maintainer-verify: checking ShellCheck 0.11.0
...
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	20.658s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (6.43s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	6.660s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	2.678s
Maintainer verification passed.
0
```

Run from the worktree root as `bash _dev/tests/maintainer-verify.sh`, unpiped, with `echo $?` on its own line.

## Integration Seams

None. The three files are self-contained, no shared registry is touched, no cross-REQ text is involved, and `generate_test.go` was not opened. The integrator still owes this REQ the usual version bump and changelog entry, which builders are forbidden to write.

One note for the merge, not a seam: `durations_test.go` is a file a sibling could plausibly touch this wave. My additions are a contiguous block inserted immediately after `durationsLabelTextDescentUnits` and before `TestDurationsLabelRowsClearTheMarkBands`; nothing existing in that file was edited or deleted.

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

## Pushback

- **The REQ's "if the row pitch changes, Panels B and C shift with it" is not achievable together with "`describeAtPointer` still resolves the same panel for the same pointer position".** Shifting Panel B means moving `DURATIONS_MEDIAN_TITLE_Y`, which is the sole input to the A/B boundary (`pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12`), so any shift moves the boundary and the second requirement fails at it. I satisfied the second requirement in full and made the first vacuous by measuring that no shift is needed — see D-03. If the intent was "shift them *if clearance is lost*", that reading is what I implemented and the clearance measurement is the evidence.

- **The REQ's second instance says `DURATIONS_LABEL_ROW_HEIGHT = 12` is smaller than a 13-unit box and reports "20 cross-row bounding-box intersections, each 1.6px deep".** My saturated fixture measures 0.834 user units of intersection, which is 1.04 CSS px at my viewport — the depth in px depends on the browser window width, since the SVG is `width:100%` over a fixed viewBox, so a px figure is not portable between measurements. The count (16–19 on my fixtures vs the REQ's 20) also depends on fixture composition. Neither discrepancy changes the finding; both are reported in user units here so the number is reproducible.

- Nothing else in the brief was wrong, and no existing test needed editing. `TestDurationsLabelRowsClearTheMarkBands`, `TestDurationLabelGeometryMatchesTheRenderer`, `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`, `TestClusteredOverflowLabelsFillBothLabelRows`, `TestOverflowLabelsGoToTheLongestSpans`, and `TestReversedLabelPlacementIsIndependentOfOverflowDensity` all pass unchanged against the new constants.

---

# Addendum — Review Remediation

**Commit:** `8de5929` (on top of `4950c81`, same branch)

The finding is correct and I have reproduced it. It is also worse than "the sweep was incomplete": **no sweep of any size could have established that bound**, and the fix the review preferred does not close it either. Both of those are measured below, not argued.

## Files Changed

```
 skills/do-work-board/tools/queue-kanban/durations.go       |  36 +++++--
 skills/do-work-board/tools/queue-kanban/durations_test.go  | 104 +++++++++++++++++----
 2 files changed, 114 insertions(+), 26 deletions(-)
```

- **`durations.go`** — `durationsLabelCharacterWidthUnits` 6.75 → **7.15**, above the measured supremum of the whole label space. The comment now carries the completeness argument rather than a sweep's worst case.
- **`durations_test.go`** — `durationsMeasuredWidestUnitsPerCharacter` (6.71, from the incomplete sweep) replaced by `durationsMeasuredLabelWidthSupremumUnits` (7.1441) with the argument that closes it; added `durationsLabelWidthModelSlackUnits`, `durationsMeasuredLabelBoxDescentUnits`, `durationsMeasuredAxisTitleAscentUnits`; the width test now pins **both sides**; added `TestDurationsLastLabelRowClearsPanelBTitle`.
- **`web/board-durations.js` is untouched by this commit.** The row pitch stays 13 and no formatter clamp was added — see D-05.

## The finding, reproduced

All three of the reviewer's strings, measured with `getComputedTextLength()` on a rendered board in an isolated Chromium:

| Label | chars | advance | u/char | vs 6.75 |
|---|---|---|---|---|
| `REQ-444 −44444h 48m` | 19 | 128.506 u | **6.7635** | over |
| `REQ-4444 4444h 48m` | 18 | 121.524 u | **6.7514** | over |
| `REQ-44444 −44444h 20m` | 21 | 142.255 u | **6.7740** | over |

And the mechanism, confirmed monotone — hour count swept at a fixed three-digit id:

```
REQ-444 4h 44m                6.6436      REQ-444 44444h 44m           6.7541  over
REQ-444 44h 44m               6.6762      REQ-444 444444h 44m          6.7747  over
REQ-444 444h 44m              6.7054      REQ-444 44444444h 44m        6.8092  over
REQ-444 4444h 44m             6.7312      REQ-444 44444444444444h 44m  6.8832  over
```

## Why no sweep could have worked, and why the preferred fix does not close it

**The label space has two independent unbounded parameters, not one.** A label is `"REQ-" + id + " " + duration`. Nothing bounds the id's digit count, and nothing bounds the hour count. Per-character width rises with length in both directions, so "the worst case I sampled" is never "the worst case".

**Clamping the hours leaves it unbounded.** This was the review's preferred fix, on the reasoning that it would make the claim "simply true". Measured with hours clamped to `999h+`, sweeping id length:

```
REQ-444 −999h+                6.7126        REQ-444 999h 44m             6.6913
REQ-4444 −999h+               6.7413        REQ-4444 999h 44m            6.7179
REQ-44444 −999h+              6.7657  over  REQ-44444 999h 44m           6.7416
REQ-444444 −999h+             6.7880  over  REQ-444444 999h 44m          6.7628  over
REQ-4444444444 −999h+         6.8552  over  REQ-4444444444 999h 44m      6.8286  over
```

A five-digit REQ id passes 6.75 with the hours already clamped. So clamping bounds one of the two unbounded parameters and buys nothing for the claim, while changing what a reader sees. **I did not clamp** — see D-05, and see Pushback for what I think should happen with the idea instead.

## What actually closes it

Only **digits** can repeat without limit. Every wide fixed character — `R`, `E`, `Q`, `m`, `h`, `−` — appears at most twice per label, so its contribution is amortized away as the label grows, and per-character width is dragged toward, and cannot pass, the per-character width of a pure digit run. That converts an unbounded search into a measurable quantity.

Measured from both directions:

| Evidence | Result |
|---|---|
| Pure 500-digit run, per character, widest digit `4` | **7.1418** |
| Marginal advance of `4` inside a 40-digit run | **7.1441** ← the supremum |
| Exhaustive over the bounded region: every digit, id lengths 1–40, hour counts 1–40, both duration forms, both signs, mixed-digit ids — **280 800 labels** | max **7.0643** |
| Limit probe, both unbounded parameters grown together — 2 010 / 10 010 / 40 010 characters | 7.1384 / 7.1411 / **7.1417**, rising, still under |
| Randomized search, 4 000 mixed-digit labels at 60+60 digits | max 6.8223 — a uniform `4` run is genuinely the worst mix |

Per-digit run widths, which is why `4` is the answer and not `9` or `8`:

```
4: 7.1418   8: 7.0881   6: 7.0667   9: 7.0667   0: 6.9915
3: 6.9593   5: 6.8626   2: 6.7015   7: 6.5400   1: 5.1658
```

**What makes this complete, stated plainly:** the space is parameterized by exactly two unbounded integers plus a finite set of choices (which digits, which duration form, sign present). The bounded parameters are covered exhaustively; the unbounded ones are covered by their limit; and monotonicity in both connects the two, so there is no gap between them for a worse label to hide in. `durationsLabelCharacterWidthUnits = 7.15` sits above the supremum, so it bounds every label the renderer can compose — including ones no board will ever produce.

### One measurement I threw away, because it was wrong

I first tried to bound the space with a **per-glyph width table** — sum a measured advance per character, and the model becomes an upper bound by construction over any string, with no sweep at all. It fails, and the failure is worth recording: summing isolated glyph advances **under-ran** the rendered string by up to 10 units over a 63 000-string corpus (`REQ-77777777777777 77777h 77m`: table 176.797, measured 186.861). Digits in this face are narrower in isolation than inside a run — `7` measures 6.3294 alone but 6.5425 inside a digit run — so a table built from isolated glyphs is not an upper bound on anything. A second attempt using worst-context *marginal* advances did produce a true upper bound over all 63 000 strings, but over-estimated realistic labels by 14–18%, worse than the flat constant it would have replaced. Dropped both.

## Red-Green Evidence

All REDs produced by restoring the pre-change blob (`git show HEAD:<path>`) or by editing the constant under test in place — never `git stash push`.

### RED A — the lower pin, against the shipped 6.75

```
231:	durationsLabelCharacterWidthUnits = 6.75
--- RED A: lower pin, shipped 6.75 vs supremum 7.1441 ---
=== RUN   TestDurationsLabelWidthEstimateCoversTheRenderedFace
    durations_test.go:534: width model assumes 6.7500 units per character, but label text reaches 7.1441 in the rendered face — the estimate under-states the text it is placing
--- FAIL: TestDurationsLabelWidthEstimateCoversTheRenderedFace (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.614s
```

### RED B — the upper pin, which the previous version did not have

```
--- RED B: upper pin, inflated 9.0 vs supremum 7.1441 + 0.25 slack ---
=== RUN   TestDurationsLabelWidthEstimateCoversTheRenderedFace
    durations_test.go:539: width model assumes 9.0000 units per character against a supremum of 7.1441 — over-estimating by more than 0.25 drops labels the rows could carry
--- FAIL: TestDurationsLabelWidthEstimateCoversTheRenderedFace (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.330s
```

### RED C — the new Panel B clearance assertion, at the pitch D-03 predicted

Note which test fails and which passes. This is the exact silent-erosion scenario the assertion exists for: raising the pitch to clear a bigger face satisfies the row-pitch **floor** while eating the **ceiling**.

```
--- RED C: Panel B clearance, row pitch raised 13 -> 15 ---
=== RUN   TestDurationsLabelRowPitchClearsTheLabelTextBox
--- PASS: TestDurationsLabelRowPitchClearsTheLabelTextBox (0.00s)
=== RUN   TestDurationsLastLabelRowClearsPanelBTitle
    durations_test.go:585: the reversed band's last label row ends at 339.41 but Panel B's title starts at 338.76 — the label rows have grown into the title
--- FAIL: TestDurationsLastLabelRowClearsPanelBTitle (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.343s
```

339.41 against 338.76 is the prose prediction from D-03 ("a pitch of 15 would put the last row's box bottom at 339.41, past the title's top at 338.77") now failing a test instead of sitting in a document.

### GREEN

```
247:	durationsLabelCharacterWidthUnits = 7.15
54:  var DURATIONS_LABEL_ROW_HEIGHT = 13;
--- GREEN: full queue-kanban suite ---
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.864s
GREEN_EXIT=0
```

## Before / after renders

Live DOM, isolated Chromium, `location.href` returned from the same `evaluate` call as every measurement.

| Board | Labels | Same-row overlaps | Tightest same-row gap | Cross-row | Label/mark | Remainder |
|---|---|---|---|---|---|---|
| Real board @ 6.75 | 3 | 0 | 76.65 u | 0 | 0 / 221 | — |
| **Real board @ 7.15** | **3** | **0** | 76.65 u | **0** | **0 / 221** | — |
| Clustered fixture @ 6.75 | 21 | 0 | 11.35 u | 0 | 0 / 63 | `+42 more` |
| **Clustered fixture @ 7.15** | **20** | **0** | **22.67 u** | **0** | **0 / 63** | `+43 more` |
| Reversed band @ 6.75 | 17 | 0 | 15.86 u | 0 | 0 / 60 | `+43 more` |
| **Reversed band @ 7.15** | **17** | **0** | 15.86 u | **0** | **0 / 60** | `+43 more` |

**On a real board the remediation costs nothing: 3 labels at 6.75 and 3 at 7.15.** The clustered fixture loses exactly one label; the saturated reversed band loses none.

Panel B clearance and pointer resolution, re-measured on the saturated-reversed board at 7.15: clearance **1.3636 units**, **0** live-DOM intersections with the title, and `describeAtPointer` resolving `A` at y ≤ 338 and `B/C` at y ≥ 339 exactly as before. `web/board-durations.js` is unchanged by this commit, so this is unchanged by construction as well as by measurement.

Screenshots: `/tmp/board-241b/shot-r2-fixture.png`, `/tmp/board-241b/shot-r2-reversed.png`. Looked at both — two cleanly separated label rows, remainder sentence on the last row's right edge, reversed band clear of Panel B's title.

## Verification

```
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	14.380s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (4.62s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	4.832s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.439s
Maintainer verification passed.
0
```

Run unpiped from the worktree root, `echo $?` on its own line.

## Integration Seams

None. Two files, no shared registry, `generate_test.go` never opened. Branch still touches nothing under `do-work/` and no version or changelog file.

I did **not** touch the main tree's `.playwright-mcp/` directory. None of its contents are mine — every file there is timestamped 13:33–15:04, before my session, and I drive my own Playwright process rather than the MCP browser, so all of it belongs to sibling agents.

## Decisions

- **D-05 — No formatter clamp.** The review preferred clamping the hour count so the label space would be finite. Measured, it is not: with hours clamped to `999h+`, `REQ-44444 −999h+` still reaches 6.7657 u/char, because the REQ id is the *other* unbounded parameter. A clamp would therefore change rendered output without making the claim true, which is the same category of error this remediation exists to remove — a mechanism that looks like it bounds the space and does not. The bound comes from the amortization argument instead, which holds with or without a clamp. **This is not a rejection of the clamp as an idea** — see Pushback.

- **D-06 — The width pin is two-sided.** A width model must over-estimate, but unbounded over-estimation is not free: it drops labels the rows could have carried, invisibly except as a larger remainder count. The one-sided assertion also could not tell a correct constant from one merely equal to the last sweep's output, which is precisely how 6.75 passed against a 6.71 reference. `durationsLabelWidthModelSlackUnits = 0.25` is the stated allowance; anything beyond it is a deliberate choice that should have to edit the slack and say why.

- **D-07 — Reach: this is now the pattern for any constant modelling a rendered face.** State the supremum over the space, not the worst case of a sample; say what makes the sweep complete; pin from both sides. The three constants in this file that model the face (`durationsLabelCharacterWidthUnits`, `DURATIONS_LABEL_ROW_HEIGHT`, `DURATIONS_LABEL_TEXT_ASCENT`) now all carry their measured value, and two of the three are pinned by a test.

## Lessons Learned

- **"Complete sweep" is a property of the argument, not the sample size.** My first sweep was 10 000 strings and wrong; my second was 280 800 and still would have been wrong on its own. What made the third one complete was noticing that the space has exactly two unbounded parameters, that per-character width is monotone in both, and that only digits can repeat — so the bounded part can be enumerated and the unbounded part replaced by a limit, with nothing in between. **When a sweep is the evidence, the reviewable claim is why the sweep closes, and that sentence belongs next to the number.**

- **A one-sided pin cannot distinguish "correct" from "equal to the last measurement".** `6.75 >= 6.71` passed for a constant that was wrong, and would have kept passing if the constant had been dropped back to 6.71. Whenever a test exists to stop a constant drifting, ask which direction it can drift and pin both; a lock-in test that only pins the direction you happened to fix last time is half a test.

- **An isolated glyph advance is not the advance that glyph gets in a string.** This face gives digits tabular widths inside runs — `7` is 6.3294 alone and 6.5425 in a run — so a per-glyph width table built from isolated measurements under-ran real strings by up to 10 units. Any font measurement taken out of context is suspect; measure the marginal contribution in the context the string actually presents, and verify the result is an upper bound over a corpus before trusting it.

- **Beware the trailing-space collapse when measuring text incrementally.** `getComputedTextLength()` ignores a trailing space, so `measure("REQ-444 " + c) − measure("REQ-444 ")` charges `c` for the space as well and reported digit advances of 10.3 units against a true 7.14. It looks like a real signal — ordered, reproducible, plausible — and it is an artifact of the probe. Pad the prefix with a sentinel character whose own width you subtract.

- **A budget recorded only in prose is not a constraint.** D-03 measured the Panel B clearance, named the pitch that breaks it, and wrote it in a document — where the next pitch change would have sailed past it while the row-pitch floor test went green. RED C shows exactly that: at pitch 15 the floor passes and the ceiling fails. **If a hand-back says "this is the binding constraint", that sentence is a request for an assertion.**

## Pushback

- **The review's preferred fix does not achieve its stated purpose, and I did not apply it.** Clamping the hour count leaves the REQ id unbounded; measured, `REQ-44444 −999h+` reaches 6.7657 u/char with hours already clamped. The stated reason for preferring it — "then the claim is simply true" — does not hold. I applied the bound that does hold instead. If this reasoning is wrong I would rather be told than have shipped a clamp that reads like a fix.

- **The clamp is still a good idea, for its other reason, as its own REQ.** A label reading `−44444h 20m` is meaningless output from a broken stamp, and rendering it as `−999h+` is better for a reader regardless of anything to do with width. That is an output-quality change to `formatDurationMinutes` in both the renderer and its Go width-model mirror, it changes what appears on screen, and it belongs next to REQ-246 (mechanical repair of fabricated timestamps) rather than buried in a constants REQ. I have deliberately not filed it — queue state is the orchestrator's.

- **One thing I found and did not fix, because it is width-neutral under this model.** `durations.go`'s `formatDurationLabelMinutes` emits an ASCII hyphen for negative spans while the renderer draws U+2212 (`−`). The two glyphs differ by 1.73 units in this face, so the Go side models a narrower string than the browser draws — but both are one character, and the width model counts characters, so nothing is currently wrong. It would become a real under-estimate the moment anyone replaces the flat constant with a per-glyph model, which is exactly what I attempted above. Worth a comment at minimum; out of this REQ's scope.
