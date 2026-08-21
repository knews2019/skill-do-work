---
id: REQ-237
title: Backfill the Durations label rows when the longest spans cluster
status: completed
completed_at: 2026-08-18T12:10:30Z
commit: 3720ab9
claimed_at: 2026-08-18T11:42:03Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T11:42:03Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - browser evidence
status_changed_at: 2026-08-18T10:46:58Z
created_at: 2026-08-18T10:52:00Z
user_request: UR-051
addendum_to: REQ-231
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
effort_estimate: normal
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Backfill the Durations Label Rows When the Longest Spans Cluster

## What

In the Durations view's overflow lane, the board picks the six longest spans in a band as label candidates, then walks them left to right and gives each the first text row where it does not touch a label already placed. A candidate that fits nowhere is simply dropped and counted. Nothing then offers the freed space to the seventh-longest span, so where several of the longest spans finish close together in time, the lane's two text rows end up mostly empty while the remainder count carries almost everything.

## Context

Found while reviewing REQ-231, which introduced the six-longest selection rule (the maintainer's chosen "Alternative 2"). Measured on two synthetic 60-sample boards:

- **Magnitude correlated with completion time** (every long span crowded at the right edge): **2 labels drawn out of 6 candidates**, 58 in the remainder.
- **Magnitude scattered across the window**: **5 of 6 place**, which is the healthy case.

So this is the tail, not the norm — but the tail is exactly the shape a burst of long REQs produces, which is also when a reader most wants the lane to talk. Before REQ-231 the lane filled both rows on the same dense fixture (27 labels), though half of them were unreadable under the dots, which is the defect REQ-231 fixed. The two label rows are a fixed budget either way; the question is only whether an unusable candidate's slot should pass to the next-longest span.

The change would be local: `selectDurationLabelCandidates` currently returns a fixed set of six before placement runs, so placement has no way to ask for a replacement. Making the two cooperate — placement pulling the next candidate when one is dropped — is a real change in shape, not a constant tweak, which is why it is a question rather than a fix.

## Requirements

- On a band where selected candidates collide, the label rows carry as many of the band's longest spans as physically fit, rather than stopping at the first six by magnitude.
- Every drawn label is still one of the band's longer spans — backfill may not reintroduce the left-edge first-fit sampling REQ-231 removed.
- Labelled + hidden still equals the band's sample count, so nothing is silently dropped.
- No change to the payload's shape (`labelRow` / `labelAnchor` / per-band hidden counts).

## Red-Green Proof

**RED prompt/case:** a test on a magnitude-gradient fixture (long spans crowded at one edge) asserting that the number of drawn labels equals what the two rows can physically hold, not the number of top-N candidates that happened to fit.
**Why RED now:** measured at 2 labels of a possible ~13 row-slots on exactly that fixture.
**GREEN when:** the same test passes, `TestOverflowLabelsGoToTheLongestSpans` still passes unchanged, and re-rendering the gradient fixture shows both rows populated with long spans.
**Validation:** apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-231: the Durations chart's top strip picks the six longest jobs to label, but if several of those six finished at nearly the same time they cannot all fit side by side, and the ones that do not fit are just dropped — nothing offers their space to the seventh-longest job instead. On a test board where the longest jobs all clustered, that left 2 labels drawn where the strip had room for about 13; on a board where the long jobs were spread out, 5 of 6 placed fine. So it only bites when a burst of slow work lands together, which is arguably when you most want to read the strip. Fixing it means letting the placement step ask the selection step for a replacement whenever it drops one, which is a genuine change to how the two cooperate rather than a tuning knob — and you may reasonably prefer the current simpler rule, where the remainder count carries the overflow and the strip stays predictable. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the remainder count already states what is not shown, and a two-step selection is more machinery than the lane is worth.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome was stated and the two files were named, but the fix required understanding why the existing packer's row occupancy was valid only under a completion-ordered walk — which is not visible from the requirement and had to be read out of the code first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — collapse selection and placement into one descending-magnitude pass
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — capacity and priority assertions
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify, **integration seam applied by the orchestrator**) — the remainder comment naming the removed selection cause

**Files I will NOT touch:** `web/board.css` (no styling change), `model.go` (no payload change), `board-timeline.js` / `generate_test.go` (sibling builders hold them).

**Acceptance criteria (restated from REQ):**
- [ ] A colliding band carries as many of its longest spans as physically fit
- [ ] Every drawn label is still one of the band's longer spans — no left-edge first-fit sampling
- [ ] Labelled + hidden equals the band's sample count
- [ ] No change to the payload's shape

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified — integration seam, applied by the orchestrator inside the merge commit)

**What was done:** `selectDurationLabelCandidates` froze the six longest spans per band *before* placement ran, so a candidate that collided was dropped and its space went to nobody. The two steps are now one pass: `durationLabelMagnitudeOrder` lists the band's samples longest-first, and `placeDurationLabelBand` offers each a row until nothing more fits. Backfill falls out of the ordering rather than being bolted on — the change **deletes** `durationsLabelTopCount` and a function rather than adding a replacement-request protocol.

Row occupancy had to change with it. The old packer summarised a row as a single `occupiedTo` float, which is only valid because a completion-ordered walk visits x monotonically; a magnitude-ordered walk does not, so `durationLabelSpanIsBlocked` now consults every interval on the row. Small n, linear scan, no structure worth more than that. The two-pass remainder reservation is preserved: pass one at full width, and only a pass that actually dropped something is redone with the last row's reserve held back.

Deleting the constant also removed the *second* reason a sample could go unlabelled, so `HiddenCount` means "did not fit" again exactly as it did before REQ-231 — and the comments REQ-231 had corrected in the other direction were corrected back, including one in `web/board-durations.js` handed over as an integration seam.

## Qualification

Passed — 3 files verified in the merge range `ca79902..3720ab9`, 4 acceptance criteria traced.

Judgment checks, measured against the merged tree in an **isolated browser tab whose URL was asserted inside the same call that read the DOM** (see the contamination note below):

| Same clustered fixture (60 overflow samples, magnitude ∝ time) | Before REQ-237 | After |
|---|---|---|
| Labels drawn (whole view) | 8 | **27** |
| Overflow lane specifically | 2 | **21** (row 0: 11, row 1: 10) |
| Remainder sentence | `+58 more over 60 min` | `+39 more over 60 min` |
| Labelled + hidden | 2 + 58 = 60 ✓ | 21 + 39 = 60 ✓ |
| Label/mark overlapping pairs | 0 | **0** — REQ-231's guarantee holds |
| Same-row label overlaps | 0 | **0** |

The remainder arithmetic cross-checks the label count independently: 58 − 39 = 19, and 21 − 2 = 19.

**The builder's second finding was confirmed rather than accepted.** It reported ~20 cross-row bounding-box intersections and argued they are line-box padding rather than ink. Measured: **20 cross-row overlaps, every one exactly 1.6px deep**, against row baselines 12 units apart — consistent with a 13.6px line box on a 12-unit pitch, and the render shows two cleanly separated rows. Zero of those overlaps are same-row, which is the property that would have mattered.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on the merged tree, run unpiped

**Red-green validation:**
- `TestOverflowLabelsGoToTheLongestSpans` (generalized) and `TestClusteredOverflowLabelsFillBothLabelRows`: ✗ → ✓. The first RED was already behavioural rather than a compile error — `2 labels drawn across 2 rows, but one row alone holds 10 at this fixture's pitch — collided spans are not being backfilled`, and `REQ-500 (65 min) carries no label, but row 0 anchor "start" is free at [9, 90] — the rows stopped short of what they hold`.
- A second RED neutered **only** the magnitude sort, leaving interval packing intact — i.e. backfill without priority, which is first-fit sampling: `REQ-504 (93 min) was passed over on row 0 anchor "end", blocked only by shorter labels (longest blocker REQ-503 at 86 min)`. The capacity test **passed** during that run, which is the useful part: it isolates the two causes cleanly — interval packing gives the count, magnitude order gives which spans.
- REQ-231's `TestDurationsLabelRowsClearTheMarkBands`, `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`, `TestReversedLabelPlacementIsIndependentOfOverflowDensity`, `TestDurationLabelGeometryMatchesTheRenderer` and the JS behaviour probe all pass unmodified.

**Existing tests updated (cross-REQ impact):**
- `TestOverflowLabelsGoToTheLongestSpans` (from REQ-231) — its two numeric assertions were `selectionFloor := magnitudes[durationsLabelTopCount-1]` and `labelledCount > durationsLabelTopCount`, which *are* the six-label cap this REQ removes. They cannot survive it: backfill means drawing the 7th-longest span when one of the top six collides, which is by definition below that floor and above that count. The test keeps its name and contract and now asserts the invariant that outlives both designs — **a span may only be passed over for one at least as long as it**, checked per candidate row/anchor slot. That is strictly stronger than the floor check: it still fails loudly when a short span takes a long span's place (proved by RED run 2) and additionally catches rows left half-empty, which the old assertion could not see.

*Verified by work action*

## Decisions

- **D-01**: One descending-magnitude pass rather than a select-then-backfill loop. Smaller change — deletes a constant and a function instead of adding a protocol between two steps. It also fixed an unnamed bug: on the 40-sample gradient fixture the old code labelled ranks 5 and 6 while dropping the four *longest* spans, because the left-to-right walk reached the shorter ones first and spent the row on them. DECIDE & STATE.
- **D-02**: No cap on drawn labels. The REQ calls the two rows "a fixed budget either way" and asks for as many of the longest as physically fit, so the walk fills them. On a perfect gradient this reaches into shorter spans (shortest drawn 1h47m against a 7h58m longest) — that is what filling a fixed budget in magnitude order means, and every drawn label is still preceded by every longer one having been offered a slot. DECIDE & STATE.
- **D-03**: Kept "end" (before the mark) as the preferred anchor — REQ-231's D-02 lesson — but rewrote its *justification*, which had stopped being true: the stated reason was that a left-to-right walk reuses space it has already passed, and this walk is not left-to-right. The preference is now a consistency choice with the after-the-mark fallback keeping the leftmost sample labellable. The builder noted a right-to-left walk might pack better preferring "start" and declined to change it on YAGNI grounds. DECIDE & STATE.
- **D-04**: Reported the width-model under-estimate rather than fixing it, despite the constant being inside the write set — retuning it changes label counts across every board and no measurement showed an actual collision. Scaling that work down is the maintainer's call. Queued as REQ-241. DECIDE & STATE.
- **D-05** (orchestrator): applied the handed-back `web/board-durations.js` comment seam inside the merge commit, so it lands in this REQ's merge range rather than as a child of it.

## Discovered Tasks

- [normal] `durationsLabelCharacterWidthUnits = 6.2` under-estimates the rendered face by ~7% (measured 6.61 units/char), while its comment claims the estimate is "deliberately generous" — the error runs the wrong way. No collision today (tightest same-row gap 3.08 units), but the margin is roughly half what the code claims and is now load-bearing. **Queued as REQ-241.**
- [trivial] `DURATIONS_LABEL_ROW_HEIGHT = 12` against a declared 13-unit text box, giving 20 cross-row bounding-box intersections at full density — line-box padding, not ink. **Consolidated into REQ-241**, same root cause: a metric constant disagreeing with the face actually rendered.
- [normal] Panel B's slowest-day annotation is drawn at `y = 355` against a title at `y = 350`, so `209 min` renders through "paused and broken spans excluded". **Pre-existing** — identical `x`/`y` on a board built from the pre-REQ-237 binary, checked side by side — and invisible on this repository only because the slowest day happens to fall clear of the title's width. **Queued as REQ-242.**
- [normal] The shared browser instance silently invalidates cross-builder DOM evidence. Recorded as a durable convention in `_dev/primes/prime-kanban-board.md` rather than as a REQ — it is guidance for the next person gathering render evidence, which is what a prime is for.

## Review

**Overall: 96%** | 2026-08-18T12:09:46Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- **The brief was wrong and the builder said so.** It instructed that `TestOverflowLabelsGoToTheLongestSpans` must pass unchanged, and that a failure there means first-fit was reintroduced. That test's assertions are literally the six-label cap this REQ removes, so no correct implementation could satisfy both. The builder found the contradiction early, flagged it as an instruction it could not satisfy as written, explained why, and generalized the assertion to a stronger invariant rather than deleting or weakening it. Recording it as a finding against **the brief**, not the build: a builder that had quietly edited the test to pass would have looked identical in the diff and been much worse.

**Restatement sweep:** this REQ *removes* a cause — a sample can no longer be unlabelled because selection passed it over — so the same sweep REQ-231 ran when it *added* that cause applies in reverse. Three comments in `durations.go` (`durationsLabelRowUnplaced`, `DurationLabelBand`, `planDurationLabels`) were corrected by the builder, and one in `web/board-durations.js` was outside its write set and handed back as a seam, applied here. `_dev/primes/prime-kanban-board.md` describes write surfaces and conventions, not the selection rule. REQ-231's archived record is dated history and stays. No stale restatement remains.

**Acceptance:** Pass — 2 → 21 labels in the overflow lane on the clustered fixture, remainder arithmetic cross-checking the count, and zero label/mark and zero same-row overlaps, all measured in an isolated browser with page identity asserted in-call.

**Suggested testing:** 2 items
- The real consuming board (677 REQs). This repository's own archive has 3 overflow samples and renders byte-identically before and after, so nothing here exercises the change on real data at scale.
- Whether filling both rows is *wanted* at full density is a taste question no test answers. D-02 fills them because the REQ asked for it; a maintainer looking at 21 labels may prefer a scarcer lane.

**Follow-ups created:** REQ-241, REQ-242; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Measuring before designing. A throwaway instrumented test across three fixture sizes turned "2 labels out of 6" into something sharper — the old packer was dropping the four *longest* spans and labelling ranks 5 and 6, because a left-to-right walk reaches the short left-hand spans first. That reframed the task from "backfill the leftovers" to "the walk was in the wrong order", which is a much smaller fix and the reason this change deletes code instead of adding it.

**What didn't:** Trusting a shared browser. The first DOM measurement returned confident, well-formed numbers for a board that was not the builder's — right shape, plausible counts, zero overlaps. Only the REQ ids gave it away. A render measurement has to assert *which page it measured* in the same breath as measuring it; a URL checked before navigation is a different claim. This is now a convention in `_dev/primes/prime-kanban-board.md`.

**Worth knowing:** A brief that pins a test as "must pass unchanged" is worth checking against the requirement text in the first ten minutes. Here the two were in direct contradiction, and finding it early is what stopped the change being contorted to satisfy an impossible constraint. The general form: when a REQ removes a rule, every test that asserted *that rule's shape* is in scope for generalization — and the honest move is to say so loudly, because a quietly-edited test and a correctly-generalized one look the same in a diff.

## Orientation

The Durations chart's overflow and reversed bands now fill their two label rows with as many of the band's longest spans as physically fit, instead of freezing a top-six list before placement and dropping whatever collided. On a board where several long REQs finish close together the lane goes from two labels to twenty-one; on this repository, which has three overflow samples, the render is byte-identical. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

**[MAP CHANGED]** — selection and placement are no longer two steps. There is no candidate list and no `durationsLabelTopCount`; there is one descending-magnitude walk that offers each span a row, and row occupancy is an interval list rather than a single high-water mark, because a magnitude-ordered walk does not visit x monotonically. Anything that later wants to influence *which* spans get labelled has one place to do it and must not reintroduce a pre-placement filter — that is the shape this REQ exists to remove. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves, the three-write-surface count is unchanged, and the file gained two conventions this session about render evidence. The prime is not stale.
