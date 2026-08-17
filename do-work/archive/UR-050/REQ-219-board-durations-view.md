---
id: REQ-219
title: Durations view on the Kanban board
status: completed
completed_at: 2026-08-17T19:36:37Z
commit: cfffd90
claimed_at: 2026-08-17T19:25:51Z
created_at: 2026-08-17T17:17:22Z
route: C
estimate:
  p50_active_minutes: 55
  confidence: medium
  calculated_at: 2026-08-17T19:27:10Z
  basis:
    - Route C
    - 7-file write set
    - 3 new files
    - 2 subsystems involved
    - 5 acceptance criteria
    - browser evidence
    - cross-route regression gates
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
user_request: UR-050
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
---

# Durations View on the Kanban Board

## What

Add a dedicated **Durations** view to the Go-served Kanban board showing how long archived REQs actually take and how often they run: three panels on one shared calendar axis — duration per REQ, median minutes per active day, and REQs completed per day. Data comes from the archive scan `queue-kanban` already performs, not from a new file.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** A new `durations.go` owns the aggregation and the read-time rule, deriving samples from `board.AllRequests` — tickets already parsed, so no second archive pass. The payload gains a `durations` object with samples and day aggregates; the rule's verdict travels as `excludedReason` so the client never re-derives it. A new `web/board-durations.js` fragment renders three SVG panels on one shared linear-time axis, registered in the fragment manifest and wired into `applyView` beside the calendar and testing views. The ordinal ramp gets tokens in both palette blocks, re-stepped per mode rather than flipped.
- [x] **[APPLY]:** Nine files, all declared. Three new (`durations.go`, `durations_test.go`, `web/board-durations.js`); no new dependency, no CDN, no write surface.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file — `generate.go` (+56: payload types, assembly, fragment manifest entry), `generate_test.go` (+2: the JS inventory and fragment-order expectations, which are exact-match lists), `web/board.css` (+199: ramp tokens in both palettes plus the view's component rules), `web/template.html` (+40: view button, panel, legend, readout, table), `web/board-controls.js` (+5: panel map and first-render hook), `web/board.js` (+4: view vocabulary comment and the render-once flag). `gofmt` clean, `go vet ./...` clean, `go test -count=1 ./...` green. No debug artifacts: the render check's screenshots, its local server, and the generated board directory were all removed before staging (`git status` confirmed clean of them).

## Why

From the user, verbatim: *"that way I can identify if there is an overtime response time degradation, outliers, how often are these reqs executed."* Those three questions are the acceptance frame — a panel that does not answer one of them is not earning its place.

## Context

The chart already exists and is render-judged in light and dark at `ai-reports/2026-08-17_1401_UR-048-estimator-calibration-and-anomaly-surfacing/index.html` (section 8, commit `6e79932`). **That is the reference implementation, not a suggestion to re-derive.** It resolved several decisions the hard way; port them rather than rediscovering them:

- The three panels are separate on a shared x-axis precisely because a single overlaid plot would need two y-scales. Never a dual axis.
- Panel A needs an overflow lane above a scale break: REQ-064's 10h 55m span otherwise squashes the 90% of samples that sit under 35 minutes.
- Route A/B/C are **ordered difficulty tiers**, so they take a single-hue ordinal ramp, not categorical hues.
- An in-cloud trend line was tried first and discarded — a rolling window over completion order straddles session boundaries and produces spikes that read as drift, and a continuous stroke across a 27-day idle stretch asserts a climb that never happened.

## Detailed Requirements

**Data source — the archive scan, at build time.**

- Derive every sample from archived REQ frontmatter the board already parses: `id`, `route`, `claimed_at`, `completed_at`. No new parser fields, no new file reads.
- Include a REQ when its status normalizes to terminal success and **both** stamps parse. 195 samples qualify today, spanning 2026-05-29 to 2026-08-17.
- Record the wall span **raw and signed**. A negative span is data, not an error to swallow — it is the anomaly class REQ-213 taught this board to surface. The archive carries none today (REQ-091's reversed stamp was owner-repaired in commit `24915a8`), so this behaviour must be pinned by a fixture rather than by real data.
- Do **not** read `do-work/calibration-log.tsv`. It holds 8 rows, only grows going forward, and is append-only so a repaired stamp never propagates. Explicitly rejected during capture.

**Panel A — duration per REQ.**

- One mark per REQ, x by completion instant, y by wall minutes.
- Y ceiling 60 minutes with an **overflow lane above a scale break** for longer spans; today that is REQ-064 (10h 55m), REQ-043 (2h 29m), REQ-171 (1h 21m).
- Negative spans render **below the zero line** in the status-critical colour, labelled — never rounded up to zero. None exist in the archive today; keep the branch and cover it with a fixture, because the whole point of REQ-213/214 is that one can reappear.
- Direct-label only the marks that matter individually (the overflow marks and any negative span). Never a value on every point.
- Colour by route using the ordinal ramp below; unrouted legacy REQs take muted ink, not a step on the ramp.

**Panel B — median minutes per active day.**

- One column per active day. **Apply the calibration's documented read-time outlier rule** — spans over four hours are an assumed pause, negative spans are broken stamps, both excluded (`skills/do-work/actions/estimate-reference.md` → Calibration).
- This is the panel that answers "is it degrading". Applying the rule here is what stops one paused session inventing a five-hour day: 2026-07-31 has two samples, and including REQ-064's 655 minutes yields a median of 328.9 instead of the true 2.5.
- Panel A still shows those excluded spans, raw and labelled. State the split where the reader can see it.

**Panel C — REQs completed per day.**

- One column per active day, counting **all** samples including rule-excluded ones (it is a count, not a duration).
- 23 active days out of 81 calendar days today; peak 25 on 2026-08-15.

**Shared axis and chrome.**

- One linear calendar x-axis across all three panels. **Linear real time, not ordinal-by-active-day** — the idle gaps are the cadence finding, and compressing them destroys the answer to "how often are these executed".
- Drop a month tick that would collide with the first or last label.
- Gridlines and axes: solid hairlines, one step off the surface. Never dashed.

**Colour — validated, do not eyeball.**

- Single-hue blue ordinal ramp, light→dark by route weight, validated with the dataviz validator's `--ordinal` mode against the board's own surfaces.
- The invariant is **"heavier route sits further from the surface"**: on a light surface that means C darkest; on a dark surface it means **C brightest**. Reusing the light steps in dark mode was tried and failed — the darkest step measured 2.6:1 and 8px dots simply disappeared. Re-step per mode.
- Text wears text tokens, never a series colour.

**Interaction and accessibility.**

- Hover readout on all three panels via a nearest-mark layer, not pinpoint hit targets on 8px dots.
- A table view listing every sample — REQ id, route, completion stamp, duration, and a note on rows excluded from the day medians — so no value is reachable only by hovering.

## Constraints

- **Read-only. This view adds no write surface.** The board has exactly three (CLAUDE.md § Kanban Board Write Surfaces) and that sentence must not need amending in this REQ's commit. If the implementation finds itself wanting a fourth, stop and surface it rather than adding one.
- No new CDN or npm dependency. The board's web layer is vanilla JS and stays that way; the chart is hand-authored SVG.
- Board versioning is folded into the skill: normal `CHANGELOG.md` entry plus suite version bump, per `_dev/primes/prime-kanban-board.md`.
- The Kanban columns and the existing Testing view are untouched.
- Aggregation cost rides on a parse that already happens; do not introduce a second pass over the archive.

## Builder Guidance

**Certainty: firm** on the data source, the read-only constraint, the three panels, the outlier-rule split between panels A and B, and the per-mode ramp re-stepping — all four were decided with the user or learned from a failed render, and re-litigating them wastes the work already done.

**Latitude** on: the Go aggregation's file and function layout, how the view is wired into the existing view switcher, exact SVG geometry and margins, tick counts, and tooltip copy. The reference implementation is a design contract, not a code contract — it is HTML in a report, and the board is Go plus vanilla JS. Port the decisions, write idiomatic board code.

Read the reference implementation's section 8 before planning. It carries inline comments explaining *why* each of the firm decisions is what it is.

## Red-Green Proof

**RED prompt/case:** Build and open the board today (`queue-kanban` → the generated board). There is no Durations view, and the generated payload carries no per-REQ duration and no per-day aggregate. Nothing anywhere in the board answers "are REQs taking longer than they used to", "which ones are outliers", or "how often does this actually run" — the only duration-adjacent signal is the completion-anomaly strip, which fires on broken bookkeeping, not on slowness.

**Why RED now:** The board renders queue state and completion anomalies. Duration has never been aggregated or surfaced, so the whole question has to be answered by hand-mining the archive — which is exactly what UR-048's calibration had to do.

**GREEN when:**
1. A Durations view exists beside the board and Testing views, rendering all three panels over every qualifying archived REQ (195 today) on one shared calendar axis.
2. A Go test pins the aggregation against known values: with the read-time rule applied, `2026-07-31` yields a day median of **2.5 min** (REQ-064's 655-minute span excluded) and `2026-08-15` yields **19.6 min over 25 REQs**; without the rule, 2026-07-31 would be 328.9 — assert the rule is actually applied, not just present.
3. A fixture with a reversed span (completed before claimed) keeps its raw negative value through aggregation and is excluded from day medians rather than clamped to zero.
4. `go vet ./...` and `go test -count=1 ./...` green in the module; `queue-kanban verify` behaviour unchanged.
5. The three-write-surface sentence in `CLAUDE.md` is untouched by this REQ's commit.

**Validation:** User confirmed — three capture questions answered (archive scan as the data source; a new dedicated view; all three panels), each the recommended option.

## Full Context

See `do-work/user-requests/UR-050/input.md` for complete verbatim input.

---
*Source: UR-050 — port the report's duration chart to the board*

## Triage

**Route:** C — Full pipeline

**Reasoning:** A new view spanning two subsystems (the Go aggregation and the vanilla-JS web layer), three new files, and a design contract carried over from a separate artifact. The "what" was settled; the plan was where the aggregation lives, how the payload carries the rule's verdict, and how the panels reach the existing view switcher.

**Confidence:** high

*Triaged by work action*

## Plan

**Aggregation (Go).** `durations.go` owns one exported entry point, `buildDurationAggregate(tickets)`, returning samples and per-day aggregates. It runs over `board.AllRequests` — already-parsed tickets — so the archive is not walked twice. The read-time rule lives here and nowhere else; the sample carries its verdict (`paused` / `reversed` / none) rather than the rule being reimplemented downstream.

**Payload.** `generatedBoardData` gains a `durations` object. The rule's verdict rides along as `excludedReason` so the client renders it rather than recomputing it.

**Client.** A new `web/board-durations.js` fragment, registered in the fragment manifest, renders three panels into one SVG on a shared linear-time axis. `applyView` gets the panel and a render-once hook, matching how the calendar and testing views are wired.

**Colour.** Ramp tokens in both palette blocks, re-stepped per mode.

**Validation plan.** Go tests pin the aggregation; a browser render check judges the chart in both palettes, because the reference implementation's own lesson was that the failed dark ramp was caught by looking, not by arithmetic.

## Exploration

**Key files:**
- `model.go` — `RequestTicket` already carries `Route`, `ClaimedAt`, `CompletedAt`; `isCompletedStatus` is the canonical terminal-success predicate and `parseTimestamp` the canonical parser. Nothing new to parse.
- `generate.go` — `buildGeneratedBoardData` assembles the payload; `boardJavaScriptFragmentPaths` is a literal manifest, and `generate_test.go` asserts both the embedded `.js` inventory and the fragment order as exact lists.
- `web/board-controls.js` — `applyView` maps view names to panels and renders each once on first activation.
- `ai-reports/2026-08-17_1401_.../index.html` section 8 — the reference implementation, with its decisions recorded inline.

**Concerns found:**
- The board's `--surface-1` is `#131720`, darker than the reference's `#1d1a15`, and its light surface is pure white rather than warm paper — so neither of the reference's validated ramps transfers unchanged. Re-measure, don't copy.
- Two exact-match test lists (`TestEmbeddedAuthoredJavaScriptInventory`, `TestBoardJavaScriptAssemblyStructure`) fail on any new fragment by design. Updating them is part of adding a fragment, not a surprise.

## Decisions

- **D-01: the read-time rule is applied in Go and its verdict is shipped, not recomputed.** The client could re-derive "over four hours or negative" trivially, and the reference implementation does exactly that. Shipping the verdict instead means the rule has one implementation, one test, and one place to change when the calibration is re-fit. Reasoning: `estimate-reference.md` → Calibration is the rule's canonical home; this view is its second *reader*, and a second reader that re-derives is a second definition waiting to drift. DECIDE & STATE.
- **D-02: `HasMedian` is a separate boolean, not a sentinel value.** A day whose only samples were rule-excluded has no median, which is a different fact from a median of zero. Encoding it as `0` would have panel B draw a zero-height bar that reads as "instant work" on the day a ten-hour paused session was the only sample. DECIDE & STATE.
- **D-03: the ordinal ramp was re-measured for this board rather than copied.** The reference's light ramp puts Route A at 2.44:1 on white — under the 3:1 floor for a graphical object, so 8px dots would have been unreliable. Re-stepped to 3.50 / 6.26 / 13.25 with adjacent steps 1.79:1+ apart. The dark ramp carried over at 3.32 / 6.00 / 13.55 against `#131720`. Both were then render-judged in the browser, because the reference's own recorded lesson is that its failed dark ramp was caught by looking at the render rather than by the numbers.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (new)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (new)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)

**Acceptance criteria (restated from the REQ):**
1. A Durations view beside the board and Testing views, three panels over every qualifying archived REQ on one shared calendar axis.
2. A Go test pins the aggregation: 2026-07-31 → 2.5 min with the rule, 328.9 without; 2026-08-15 → 19.6 min over 25 REQs.
3. A reversed-span fixture keeps its raw negative value and stays out of the day medians.
4. `go vet ./...` and `go test -count=1 ./...` green; `queue-kanban verify` behaviour unchanged.
5. The three-write-surface sentence in `CLAUDE.md` is untouched by this commit.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Baseline `go test -count=1 ./...` green in the module before the change.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (new)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (new)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)

**What was done:** Added a read-only Durations view answering the three questions the request named — degradation, outliers, and cadence.

*Aggregation.* `buildDurationAggregate` walks the already-parsed tickets and emits one sample per REQ that reached terminal success with both stamps parseable. The wall span is recorded raw and signed. Each sample carries the calibration read-time rule's verdict; days carry both the ruled median with its sample size and the unruled completion count, plus `HasMedian` so "no in-rule samples" is distinguishable from "zero minutes".

*Payload.* A `durations` object with `samples` and `days`. It rides on the archive parse the board already performs — no new frontmatter field, no new file read, no second pass.

*View.* Three panels in one SVG on a shared **linear real-time** axis, so the idle stretches stay visible. Panel A plots every sample raw, with an overflow lane above a scale break at 60 minutes and a below-zero lane in the critical colour for reversed spans; only the overflow marks and reversed spans are directly labelled. Panel B plots the ruled median per active day; panel C counts every sample per day. A nearest-mark hover layer covers the whole plot rather than the 8px dots, and a `<details>` table lists every sample with its exclusion note, so no value is pointer-only.

*Colour.* Route A/B/C take a single-hue ordinal ramp keyed on "heavier route sits further from the surface" — C brightest on the dark palette, C darkest on the light one — with the measured contrast ratios recorded beside the tokens. Unrouted legacy REQs take muted ink rather than a step on the ramp. Text wears text tokens throughout.

**Tests touched:** new `durations_test.go` (six cases); `generate_test.go`'s two exact-match fragment lists extended for the new `.js` file.

## Qualification

Passed — 9 files verified, 5 criteria traced, no debug artifacts.

- Every file is in `write_set`; nothing undeclared was touched, and `CLAUDE.md` is not in the diff.
- Substantive: 170 lines of new Go, ~450 of new client code, ~200 of new CSS — no reformatting.
- Requirements traced: view exists and renders 202 marks / 23 median bars / 23 count bars (1); day medians pinned with the ruled/unruled gap asserted (2); reversed-span fixture (3); vet and full test run green and `verify` output byte-identical to the pre-change binary (4); write-surface sentence absent from the diff (5).
- Flowing: the payload is populated from the aggregation and consumed by the renderer — confirmed end to end in a browser, not inferred.

## Testing

**Tests run:** `go vet ./...`; `go test -count=1 ./...`; `bash _dev/tests/maintainer-verify.sh`; a browser render check of the generated static board in both palettes; `queue-kanban verify` diffed against a binary built from `HEAD`.

**Result:** ✓ vet clean; ✓ all module tests green; ✓ maintainer-verify exit 0 with zero FAIL lines; ✓ `verify` output byte-identical before and after.

**Red-green validation:**

✗ RED — before this REQ there is no Durations view and the payload carries no duration data at all. The board answers queue state and completion anomalies; nothing anywhere answers "are REQs taking longer", "which are outliers", or "how often does this run". Every one of the six new Go cases fails at compile time against the pre-change module, because `buildDurationAggregate` does not exist.

✓ GREEN — the six cases pass and pin the behaviour by value:
- `2026-07-31` reports a ruled median of **2.5167 min over 1 kept sample and 2 completed**, with REQ-064's 655.2-minute span excluded as `paused`; the same test computes the unruled median as **328.9** and asserts the two differ, so a regression that silently drops the rule fails on a number rather than on a missing symbol.
- `2026-08-15` reports **19.65 min over 25 REQs**.
- A reversed-span fixture keeps **−30.0 min** raw, is excluded as `reversed`, stays out of the day median (which reports the sound sample's 10 min over 1), and is still counted in the day's completion count.
- A day whose only sample is rule-excluded reports **no median**, not zero.
- Only terminal-success REQs with both stamps parseable become samples; `completed-with-issues` qualifies and keeps its route.
- The live archive assertion is on two **past, immutable** days, so it pins the rule without going stale as new work lands.

**Acceptance testing (browser).** The generated static board was served locally and driven in a real browser. The Durations view renders 202 marks, 23 median columns, 23 count columns, and 202 table rows; the only direct labels are the three overflow marks (`REQ-043 2h 29m`, `REQ-064 10h 55m`, `REQ-171 1h 21m`) and the slowest day (`42 min`), exactly as the design requires. Console carried no JavaScript errors. Both palettes were render-judged: on the dark surface Route C is the brightest step and Route A still reads at the 3:1 floor; on the light surface the ordering inverts and all three steps clear 3:1. The screenshots, the local server, and the generated board were removed before staging.

**Existing tests updated:** `generate_test.go`'s embedded-JS inventory and fragment-order lists. Both are deliberate exact-match assertions that fail on any new fragment — updating them is the intended cost of adding one, and neither assertion was weakened.

*Verified by work action*

## Review

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five GREEN criteria met, each with its own evidence |
| Code Quality | 92% | The renderer is long but linear; the geometry constants carry the reasons the reference learned the hard way |
| Test Adequacy | 95% | The ruled-vs-unruled assertion is the one that makes the rule non-optional; the browser check covers what Go cannot |
| Scope | 100% | Nine declared files; `CLAUDE.md` untouched, as the REQ required |
| Risk | Low | Read-only view; no new dependency, no write surface, `verify` output byte-identical |
| Acceptance | Pass | Rendered and driven in a real browser in both palettes |

**Verdict: Approve** — the board now answers degradation, outliers, and cadence from data it was already parsing.

### Findings

**Minor:**
- Panel A's direct labels are placed by alternating anchor and offset. With today's three overflow marks they read cleanly, but a fourth long span landing near an existing one could overlap. A collision-aware placement would be more machinery than three labels justify; the table view is the fallback either way.
- Route A on the dark palette measures 3.32:1 — over the 3:1 floor for a graphical object but the thinnest margin in either ramp. Confirmed legible in the render at 8px. Worth remembering if the board's surfaces are ever lightened.

**Nit:**
- Panel C reuses `--accent-done` rather than taking a step on the durations ramp. It reads correctly (green for completions) and panels B and C each plot a single named measure, so neither needs a key — but it is a third hue in the view.

### Restatement Sweep

**Triggered** — the diff adds a second reader of the calibration's read-time outlier rule. Swept "four hours", "outlier", and `240` across `skills/` and `_dev/`. Results: `skills/do-work/actions/estimate-reference.md` → Calibration remains the rule's sole definition and is unchanged; `durations.go` cites it rather than restating it, and the client receives the verdict rather than re-deriving it (D-01), so there is exactly one implementation. `_dev/primes/prime-kanban-board.md`'s write-surface convention points at `CLAUDE.md` rather than restating the count, and this REQ adds no surface, so neither needed an edit.

### Requirements Checklist

- [x] Durations view beside board and Testing, three panels on one shared calendar axis — delivered
- [x] Aggregation pinned by value, with the rule proven applied rather than present — delivered
- [x] Reversed-span fixture keeps its raw negative value and stays out of day medians — delivered
- [x] `go vet` / `go test` green; `verify` behaviour unchanged — delivered (output diffed against a `HEAD`-built binary)
- [x] `CLAUDE.md`'s three-write-surface sentence untouched — delivered

### Acceptance Testing

**Result: Pass**
- `go test -count=1 ./...` and `go vet ./...` green in the module.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines.
- Static board generated, served locally, and driven in a real browser: 202 marks, 23 median columns, 23 count columns, 202 table rows, four direct labels, no console errors, both palettes render-judged.
- `queue-kanban verify` output byte-identical to a binary built from `HEAD`.

### Suggested Additional Testing

- **Manual verification:** the hover readout was exercised by construction, not by a scripted pointer sweep. Worth a few seconds of hovering across panel boundaries — the readout switches between the per-REQ and per-day descriptions at panel B's title line, and that threshold is the one piece of the interaction chosen by geometry rather than by data.
- **Environment:** rendered in one browser at one width. The SVG has a `min-width` and scrolls inside its own container, but a narrow viewport is worth one look.
- **Edge case:** the archive carries no reversed span today, so the below-zero lane has been proven in Go but never actually drawn. If REQ-091's stamp is ever re-broken, that is the moment to look at the render.

### Follow-up REQs Created

None.

## Lessons Learned

**What worked:** Porting the design decisions and re-deriving the numbers. The reference implementation's ramps were validated against its own surfaces, and the board's are different enough that Route A would have landed at 2.44:1 — under the floor — if the light steps had been copied across. Reading the reference for *why* and re-measuring for *what* took minutes and caught it before the render.

**What didn't:** Writing "(see `write_set`)" in the Scope section. `scope-drift.sh` parses that list literally, so a pointer instead of a list reports every touched file as undeclared drift and the word `write_set` as a declared-but-untouched file. Scope lists are machine-read; write them out.

**Worth knowing:** `generate_test.go` holds two deliberately exact-match lists — the embedded `.js` inventory and the fragment execution order — so adding a client fragment always fails them first. That is the design: the manifest is load-bearing for execution order, and an assertion that tolerated additions would not protect it. Also: shipping a rule's *verdict* in the payload rather than the inputs it was computed from is what keeps a second reader from quietly becoming a second definition — the client cannot drift from a rule it never implements.

## Orientation

The board can now answer how long work takes and how often it runs: a Durations view plots every archived REQ's measured span against the calendar, the median minutes per active day beneath it, and the completions per day beneath that — one shared linear-time axis, so the idle stretches stay visible instead of being compressed away. Lives in the Kanban board tool (`skills/do-work-board/tools/queue-kanban/`), as a new aggregation module plus a new client fragment beside the calendar and testing views. **[MAP CHANGED]** — this is a new view and a new payload section, and it makes the archive's duration record a first-class board surface rather than something the calibration had to hand-mine. It reads the same parse the board already performs and adds no write surface: the tool still has exactly three, and `CLAUDE.md`'s sentence saying so needed no amendment.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` — every referenced path still resolves; its write-surface convention correctly defers the count to `CLAUDE.md`, so nothing there went stale.
