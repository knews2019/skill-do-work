# REQ-237 Hand-Back — Backfill the Durations Label Rows When the Longest Spans Cluster

**Branch:** `worktree-agent-REQ-237-durations-label-backfill`
**Commit:** `8608225` `[REQ-237] Backfill the Durations label rows when the longest spans cluster`
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-237-durations-label-backfill`
**maintainer-verify.sh exit code:** `0`

## File Manifest

| File | Verb | What |
|---|---|---|
| `skills/do-work-board/tools/queue-kanban/durations.go` | modified | Collapsed candidate selection + placement into one descending-magnitude pass; row occupancy became interval lists; deleted `durationsLabelTopCount` |
| `skills/do-work-board/tools/queue-kanban/durations_test.go` | modified | Generalized `TestOverflowLabelsGoToTheLongestSpans` to the priority invariant; added `TestClusteredOverflowLabelsFillBothLabelRows`; extracted `placedDurationLabelsByRow` so label geometry has one definition in the tests |

`git diff --stat HEAD~1 HEAD`:

```
 .../do-work-board/tools/queue-kanban/durations.go  | 142 ++++++++--------
 .../tools/queue-kanban/durations_test.go           | 182 +++++++++++++++++----
 2 files changed, 221 insertions(+), 103 deletions(-)
```

Nothing outside the declared write set. Working tree clean at hand-back.

## What Changed

`selectDurationLabelCandidates` froze the six longest spans per band **before** placement ran, so placement had no way to ask for a replacement — a candidate that collided was dropped and its space went to nobody. The two steps are now one: `durationLabelMagnitudeOrder` lists the band's samples longest-first, and `placeDurationLabelBand` offers each one a row until nothing more fits. Backfill falls out of the ordering rather than being bolted on, which is the shape the brief suggested and it checked out.

Two consequences worth flagging:

1. **Row occupancy had to become interval lists.** The old packer summarized a row as one `occupiedTo` float, which is only valid because a completion-ordered walk visits x monotonically. A magnitude-ordered walk does not, so `durationLabelSpanIsBlocked` now consults every interval on the row. Small n (≤ ~30 per row), linear scan, no structure worth more than that.
2. **`durationsLabelTopCount` is deleted**, and with it the *second* reason a sample could go unlabelled. `HiddenCount` means "did not fit" again, exactly as it did before REQ-231 — so REQ-231's own restatement-sweep lesson applied in reverse, and I corrected the `DurationLabelBand` / `LabelRow` / `planDurationLabels` comments that named the now-gone selection cause. Delete-before-you-add: this change removes a constant and a function rather than adding a backfill mechanism.

The two-pass remainder reservation is preserved unchanged — pass one at full width, and only a pass that actually dropped something is redone with the last row's reserve held back.

## Requirements Traced

| Requirement | Where |
|---|---|
| Colliding bands carry as many of the band's longest spans as physically fit | Descending walk continues past a rejected span; measured 2 → 21 labels on the clustered fixture |
| Every drawn label is still one of the band's *longer* spans; no left-edge first-fit | Magnitude order is the only selection rule there is; asserted by `assertDurationLabelPriority` |
| Labelled + hidden = the band's sample count | Asserted in `TestOverflowLabelsGoToTheLongestSpans`; verified in the rendered payload (21 + 39 = 60) |
| No change to payload shape (`labelRow` / `labelAnchor` / hidden counts) | No struct field added or removed; `board-data.js` carries the same keys before and after |

## RED / GREEN

### RED run 1 — behavioural, against unchanged production code

The new assertions compile against existing symbols, so this first RED is already a real assertion failure rather than a compile error:

```
=== RUN   TestOverflowLabelsGoToTheLongestSpans
    durations_test.go:592: overflow band: REQ-500 (65 min) carries no label, but row 0 anchor "start" is free at [9, 90] — the rows stopped short of what they hold
--- FAIL: TestOverflowLabelsGoToTheLongestSpans (0.00s)
=== RUN   TestClusteredOverflowLabelsFillBothLabelRows
    durations_test.go:632: 2 labels drawn across 2 rows, but one row alone holds 10 at this fixture's pitch — collided spans are not being backfilled
--- FAIL: TestClusteredOverflowLabelsFillBothLabelRows (0.00s)
FAIL
```

### RED run 2 — code present, one behaviour deliberately wrong

With the full implementation in place I neutered **only** the magnitude sort in `durationLabelMagnitudeOrder` (leaving interval packing intact), i.e. backfill without priority — which is precisely first-fit sampling:

```
=== RUN   TestOverflowLabelsGoToTheLongestSpans
    durations_test.go:592: overflow band: REQ-504 (93 min) was passed over on row 0 anchor "end", blocked only by shorter labels (longest blocker REQ-503 at 86 min) — labels are not going to the longest spans
--- FAIL: TestOverflowLabelsGoToTheLongestSpans (0.00s)
=== RUN   TestClusteredOverflowLabelsFillBothLabelRows
--- PASS: TestClusteredOverflowLabelsFillBothLabelRows (0.00s)
FAIL
```

The capacity test still passing here is the useful part: it isolates the two causes cleanly. Interval packing gives you the *count*; magnitude order gives you *which spans*. Each test fails for its own reason and nothing else's.

### GREEN

```
=== RUN   TestDenseOverflowLabelsStayBoundedAndNeverOverlap
--- PASS: TestDenseOverflowLabelsStayBoundedAndNeverOverlap (0.00s)
=== RUN   TestReversedLabelPlacementIsIndependentOfOverflowDensity
--- PASS: TestReversedLabelPlacementIsIndependentOfOverflowDensity (0.00s)
=== RUN   TestDurationsLabelRowsClearTheMarkBands
--- PASS: TestDurationsLabelRowsClearTheMarkBands (0.00s)
=== RUN   TestOverflowLabelsGoToTheLongestSpans
--- PASS: TestOverflowLabelsGoToTheLongestSpans (0.00s)
=== RUN   TestClusteredOverflowLabelsFillBothLabelRows
--- PASS: TestClusteredOverflowLabelsFillBothLabelRows (0.00s)
=== RUN   TestDurationLabelGeometryMatchesTheRenderer
--- PASS: TestDurationLabelGeometryMatchesTheRenderer (0.00s)
=== RUN   TestJavaScriptBehaviorDurationsLabelsFollowTheShippedVerdict
--- PASS: TestJavaScriptBehaviorDurationsLabelsFollowTheShippedVerdict (0.24s)
PASS
```

`bash _dev/tests/maintainer-verify.sh` → **exit 0** (run twice: after implementation, and again after the final comment correction). Not piped through `tail`/`head`; `echo $?` on its own line both times.

## ⚠️ One Brief Instruction I Could Not Satisfy As Written

**`TestOverflowLabelsGoToTheLongestSpans` could not pass unchanged.** The brief said a failure there means first-fit sampling was reintroduced and the fix is wrong. That is not what happened — the test's two numeric assertions *are* the rule REQ-237 replaces:

```go
selectionFloor := magnitudes[durationsLabelTopCount-1]   // literally the 6th-longest span
...
if labelledCount > durationsLabelTopCount { t.Fatalf(...) }  // at most 6 labels
```

Backfill means drawing the 7th-longest span when one of the top six collides. That span is *by definition* below the 6th-longest floor, and the count *by definition* exceeds six. No implementation of the REQ's first requirement can keep both assertions — including a variant capped at six drawn labels, which still trips the floor. The conflict is in the requirement, not in the implementation.

**Resolution:** I kept the test's name and its contract and generalized the assertion to the invariant that outlives both designs — *a span may only be passed over for one at least as long as it*. Concretely, for every unlabelled sample, every row/anchor slot its text could legally occupy must be blocked by a drawn label of at least its own magnitude. That is strictly stronger than the old floor check in the way that matters: it still fails loudly the moment a short span takes a long span's place (RED run 2 proves it does), and it additionally catches rows left half-empty, which the old assertion could not see. `durationsLabelTopCount` no longer exists to assert against.

Flagging rather than deciding quietly, since the brief was emphatic. If you would rather have a hard cap on drawn labels, say so and it is a one-line change — but the REQ text ("the label rows carry as many of the band's longest spans as physically fit"; "the two label rows are a fixed budget either way") reads to me as explicitly wanting the rows filled.

## Render Evidence

Fixture: 60 archived REQs across a 48h window, spans 65 → 478 min growing left to right, so magnitude correlates perfectly with completion time — the clustering shape the REQ measured. Generator: `/tmp/make-237-fixture.sh`, fixture root `/tmp/qk-237-fixture`. Boards built from the pre-change and post-change binaries (`/tmp/qk-237-before`, `/tmp/qk-237-after`).

**Label counts on the same clustered fixture, measured in the live DOM:**

| | Before | After |
|---|---|---|
| Labels drawn in the overflow lane | **2** | **21** (11 on row 0, 10 on row 1) |
| Remainder sentence | `+58 more over 60 min` | `+39 more over 60 min` |
| Labelled + hidden | 2 + 58 = 60 ✓ | 21 + 39 = 60 ✓ |
| Overlapping label/mark pairs | 0 | **0** |
| Same-row label collisions | 0 | **0** (min gap 3.08 user units) |
| Leader ticks | 2 | 21 (one per label) |

I did drive a real browser for this, but **not the shared Playwright one** — see the caveat below. The measurement is `getBoundingClientRect()` / `getBBox()` intersection counting in a headless Chromium I launched and controlled myself over CDP (`/tmp/measure-237-dom.mjs`, a ~25-line zero-dependency driver using Node 22's global `WebSocket`).

Screenshots at `/tmp/durations-237-before.png` and `/tmp/durations-237-after.png` — I looked at both. Before: 60 dots across the full lane, two labels crammed into the far-right corner, the entire text band empty. After: both rows populated edge to edge with legible labels, ascending magnitudes, leader ticks tying each to its dot, remainder clear at the right end of row 1, nothing touching a mark.

**Real archive check:** rendering the worktree's own `do-work/` (3 overflow samples, 0 reversed) gives **3 drawn / 0 hidden both before and after** — byte-identical outcome. This change only bites in the tail, exactly as the REQ predicted.

### ⚠️ A contaminated measurement I caught and discarded

My first DOM measurement ran in the shared Playwright browser and reported "3 labels, 0 overlaps" for the *before* board. It was wrong: a sibling builder shares that browser instance and had navigated the page to `http://127.0.0.1:8242/index.html` (title `worktree-agent-REQ-240-timeline-axis-minute`) between my navigate and my evaluate. I was measuring **their** board. I caught it because the label texts came back as `REQ-043` / `REQ-064` / `REQ-171` — REQ ids that do not exist in my synthetic fixture.

Every number in the table above comes from the isolated browser instead, and each measurement re-checks `location.href` inside the same evaluate before reading anything. **If another builder in this batch reports browser-measured evidence from the shared Playwright instance this session, it is worth asking whether they checked the page identity** — the failure is silent and produces plausible-looking numbers.

## Two Measurements Worth Your Attention (neither is a defect, both are pre-existing)

**1. The label width model under-estimates the real face by ~7%.** `durationsLabelCharacterWidthUnits = 6.2` with a comment calling it "deliberately generous for the 11px sans face". Measured: a 14-character label renders **92.52 user units** wide, i.e. **6.61 units/char**. The 6-unit separation rule absorbs the difference, so nothing collides — the tightest same-row gap in the 21-label render is 3.08 units — but the real margin is about half what the constant claims. This never mattered while only 2–3 labels were ever drawn; with the rows full it consumes half the separation budget. The constant is inside my write set and I deliberately left it alone: retuning it changes label counts across the whole view, which is beyond what REQ-237 asked for. Filed as a discovered task instead.

**2. The two label rows overlap by ~1 unit by the renderer's own declared metrics.** `DURATIONS_LABEL_ROW_HEIGHT = 12`, but the declared text box is `DURATIONS_LABEL_TEXT_ASCENT = 11` + 2 descent = 13 units; the rendered font box measures 12.83 units on that 12-unit pitch. So 19 cross-row bounding-box intersections appear once both rows are densely populated. **This is font line-box padding, not ink** — the screenshot shows two cleanly separated rows, and REQ-231's geometry test only ever asserted rows against the *mark band*, never row against row. Backfill did not cause it; backfill is what makes it reachable. Fixing it would mean `DURATIONS_LABEL_ROW_HEIGHT` in `web/board-durations.js`, outside my write set.

## Integration Seams

**One, in `web/board-durations.js` — a comment, not behaviour.** Lines 338-340:

```js
// Whatever carries no label is stated, never dropped in silence — samples
// selection passed over and samples placement could not fit alike: the
// count is what stops a reader taking the visible labels for all of them.
```

"samples selection passed over" no longer describes anything — selection is gone, and a sample can only be unlabelled because placement could not fit it. This is REQ-231's own lesson running in reverse (it added the second cause; this REQ removes it). The file is outside my write set and a sibling builder is in `web/board-timeline.js`, so I did not touch it. Suggested replacement for that first clause: *"Whatever carries no label is stated, never dropped in silence — the count is what stops a reader taking the visible labels for all of them."*

No behavioural seam: the payload shape is unchanged, and the JS behaviour probes for the Durations labels pass untouched.

## P-A-U

**[PLAN]** Read `CLAUDE.md`, `_dev/primes/prime-kanban-board.md` (including its REQ-231/REQ-235/REQ-236 lesson links), `crew-members/testing.md`, the REQ-237 brief body, and the archived REQ-231 record's Implementation Summary + Lessons Learned. Then measured the current behaviour before designing anything — a throwaway instrumented test printing placed-label counts across three fixture sizes, which is where the 2-of-6 figure and the fact that the *longest* span (REQ-539, 338 min) was itself going unlabelled both came from. Settled on the brief's suggested shape — collapse selection and placement into one descending-magnitude pass — after confirming it needed interval-based row occupancy, since the old single-float occupancy silently assumes a monotonic-x walk. Checked it against `TestOverflowLabelsGoToTheLongestSpans` first, which is where I found that test's assertions cannot survive the change (see the flagged conflict above), and against REQ-231's D-02 anchor lesson, which I kept while correcting its now-false stated rationale.

**[APPLY]** Code stayed inside the write set: `durations.go` and `durations_test.go` only. Two places pushed at that boundary and I did not cross either: the stale remainder-cause comment in `web/board-durations.js` (reported as an integration seam above), and `DURATIONS_LABEL_ROW_HEIGHT` in the same file for the row-pitch measurement. One judgement call *inside* the write set: I corrected a pre-existing stale comment on `packDurationLabelBand` that said the remainder reservation sits at "row 0's right edge" when the code and `durationsLabelRemainderReserveUnits` both say the last row — same function I was rewriting, no behaviour change. All scratch (fixture, binaries, boards, screenshots, CDP scripts) is under `/tmp`; nothing was written to the main tree except this hand-back.

**[UNIFY]** `git diff --stat HEAD~1 HEAD`: `durations.go` 142 lines changed, `durations_test.go` 182 lines changed, 2 files, 221 insertions / 103 deletions. In `durations.go` I checked that the payload shape is untouched, that the two-pass reservation still runs only when something was dropped, that the anchor preference order is unchanged in behaviour, and that every comment naming the deleted selection step was corrected. In `durations_test.go` I checked that the two label-geometry helpers now have one definition (`placedDurationLabelsByRow`), that both new assertions fail for their own distinct reason (RED run 2), and that no test asserts on the internals of the new packer rather than its output. Linters: `gofmt -l` clean (empty output), `go vet ./...` clean, `bash _dev/tests/maintainer-verify.sh` exit 0 including the uncached queue-kanban suite and the strict JavaScript behaviour lane. No debug artifacts in the diff — the instrumented measurement test was written to `/tmp` and copied in/out, never committed; `git status` is clean and `git status --porcelain --untracked-files=all` shows no untracked files in the worktree.

## Decisions

1. **One descending-magnitude pass rather than a select-then-backfill loop.** The brief offered this and it is the smaller change: it deletes a constant and a function instead of adding a replacement-request protocol between two steps. It also fixed a bug nobody had named — on the 40-sample gradient fixture the old code labelled ranks 5 and 6 while dropping the four *longest* spans, because the left-to-right walk reached the shorter ones first and spent the row on them.

2. **No cap on drawn labels.** The REQ says the rows should carry "as many of the band's longest spans as physically fit" and calls the two rows "a fixed budget either way", so I let the walk fill them. On a perfect magnitude/x gradient this reaches fairly deep into the shorter spans (the shortest drawn is 1h47m against a 7h58m longest) — that is what filling a fixed budget in magnitude order means, and every drawn label is still preceded by every longer one having been offered a slot. Easy to change if you want the lane scarcer.

3. **Kept "end" (before the mark) as the preferred anchor, but rewrote its justification.** REQ-231's D-02 lesson is about not losing the preference, and I did not. Its stated *reason* — a left-to-right walk reuses space it has already passed — is no longer true of a magnitude-ordered walk, so the comment now says the order is a consistency choice with the after-the-mark fallback keeping the leftmost sample labellable. Worth noting: by REQ-231's original packing logic, a walk that runs right-to-left (which a gradient fixture's magnitude order effectively is) would pack better preferring "start". I did not change it — it would alter the real board's appearance for a speculative gain, and YAGNI.

4. **Reported the width-model under-estimate rather than fixing it.** In my write set, and tempting, but retuning `durationsLabelCharacterWidthUnits` changes label counts across every board and no measurement showed an actual collision. Scaling the work down is your call, not mine.

## Discovered Tasks

- **[normal] The label width model under-estimates the shipped face.** `durationsLabelCharacterWidthUnits = 6.2` against a measured 6.61 units/char for the 11px sans (a 14-char label renders 92.52 units, not 86.8). The comment claims the estimate is "deliberately generous", and it is not. No collisions today — the 6-unit separation rule absorbs it, leaving 3.08 units of real clearance at worst — but the stated safety margin is roughly half what the code claims, and now that the rows actually fill, that margin is load-bearing.
- **[trivial] The two label rows overlap by ~1 unit by the renderer's own declared metrics.** `DURATIONS_LABEL_ROW_HEIGHT = 12` with a declared text box of 13 units (ascent 11 + descent 2); measured font box 12.83 units. Invisible today because the ink sits well inside the line box, but it means no test can honestly assert row-against-row separation the way `TestDurationsLabelRowsClearTheMarkBands` asserts row-against-mark. Either the row height wants to be 14, or the ascent constant is describing a line box rather than an ascent.
- **[normal] The shared Playwright browser silently invalidates cross-builder DOM evidence.** Two of my measurements were contaminated by a sibling builder navigating the shared page mid-call, and the numbers came back plausible rather than obviously broken. Anything in the batch that used it for render evidence is suspect unless it checked page identity inside the same evaluation. Worth a standing note in `_dev/primes/prime-kanban-board.md` about checking `location.href` in-call, or a one-liner CDP driver in `_dev/` so builders get an isolated browser by default — mine is 25 lines with no dependencies.

## Lessons Learned

**What worked:** Measuring before designing. A throwaway instrumented test across three fixture sizes turned "2 labels out of 6" into something more specific — the old packer was dropping the four *longest* spans on that fixture and labelling ranks 5 and 6, because a left-to-right walk reaches the short left-hand spans first. That reframed the change from "backfill the leftovers" to "the walk was in the wrong order", which is a much smaller fix.

**What didn't:** Trusting a shared browser. The first DOM measurement returned confident, well-formed numbers for a board that was not mine. Nothing about the result looked wrong — the shape was right, the counts were plausible, the overlap count was zero. Only the REQ ids gave it away. A render measurement needs to assert *which page it measured* in the same breath as measuring it; a URL checked before navigation is not the same claim.

**Worth knowing:** When a REQ removes a cause, the same restatement sweep applies as when one adds a cause. REQ-231's lesson was written about adding a second reason a sample could go unlabelled, and it names the comments that became half-truths. Removing that second reason makes exactly the same set of sentences wrong again, in the other direction — including one in a file I could not touch. A brief that pins a test as "must pass unchanged" is worth checking against the requirement text early: here the two were in direct contradiction, and finding that in the first ten minutes rather than after implementing saved rewriting the change to satisfy an impossible constraint.
