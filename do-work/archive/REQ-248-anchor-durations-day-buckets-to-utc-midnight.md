---
id: REQ-248
title: Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas
status: completed
created_at: 2026-08-18T13:54:59Z
claimed_at: 2026-08-18T16:09:27Z
completed_at: 2026-08-18T18:24:42Z
kb_status: pending
route: B
status_changed_at: 2026-08-18T13:54:59Z
user_request: UR-051
addendum_to: REQ-242
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/durations.go
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
---

# Anchor the Durations Day Buckets to UTC Midnight So Panel B Stays on Canvas

## What

Panel B's bars are placed with `xOfEpoch`, which maps each day bucket's **midnight**, while `timeStart` is the **first completion instant**. The two disagree by however far into its first day the earliest sample falls, so the leftmost bar is drawn to the left of the plot area — and on a board with one or two active days the disagreement dominates the whole span and the panel renders off-canvas entirely.

## Context

Found by REQ-242's builder as an unrelated pre-existing quirk, then confirmed and extended by REQ-242's independent review. This is not cosmetic at low day counts.

## Instances

- [ ] **Leftmost bar sits in the axis gutter on the real board.** `x=37.1 width=12` spans 37.1–49.1, entirely left of `DURATIONS_MARGIN_LEFT` (54), and the render shows it struck through by the "0" axis tick. Visible on this repository's own board today.
- [ ] **One active day: Panel B renders empty.** `timeSpan` collapses to the intra-day sample span (about 3 hours), so `xOfEpoch(midnight)` maps to roughly minus three plot-widths — measured annotation at `x=-3330`, bar at `x=-3342`. Both completely off-canvas.
- [ ] **Two active days: same failure, smaller magnitude.** Annotation at `x=-336.5`, bar at `x=-348.5`.

## Requirements

- Every Panel B bar renders inside the plot area at every day count, including one and two active days.
- The slowest-day annotation renders on canvas at every day count — it exists to state a value a clipped bar cannot, and cannot do that from off-screen.
- No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary.
- REQ-241's and REQ-242's guarantees hold unchanged: 0 same-row label overlaps, 0 label/mark overlaps, the annotation clear of every neighbour in its strip.

## Builder Guidance

The suggested root fix from the review is to floor `timeStart` to its UTC midnight and ceil `timeEnd` to the following midnight before computing `timeSpan`, so the axis domain and the day buckets share one origin. Verify that against the other panels before adopting it — Panels A and C read the same domain.

**Generate a board and look at it**, at one, two and many active days. Measure in the live DOM.

## Red-Green Proof

**RED prompt/case:** a test asserting every Panel B bar's x-range and the annotation's x-range fall inside the plot area, evaluated on one-day, two-day and many-day fixtures.
**Why RED now:** measured `x=-3330` on a one-day board and `x=37.1` against a left margin of 54 on the real board.
**GREEN when:** the assertion passes at every day count and a render at one and two days shows Panel B populated.
**Validation:** Review finding on REQ-242; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect is reproduced and the suggested root fix is named, but it has to be checked against Panels A and C which read the same domain, and the day-count evidence needs live measurement — the what is clear, the blast radius is not.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — anchor the axis domain and the day buckets to one origin
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the day-count lock-in probe
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — label planner floors/ceils to the same UTC-day domain (added mid-build via D-01; orchestrator-accepted, see Decisions)

**Files I will NOT touch:**
- `durations_test.go` — REQ-252 owns the measured-face provenance work and is gated behind this REQ. (`durations.go` moved to the touch list via D-01 — the original exclusion's collision reason was empty since REQ-252 is gated behind this REQ.)
- `DURATIONS_MEDIAN_TITLE_Y` and `describeAtPointer`'s A/B boundary — named unchanged by the REQ.

**Acceptance criteria (restated from REQ):**
- [ ] Every Panel B bar renders inside the plot area at every day count, including one and two active days.
- [ ] The slowest-day annotation renders on canvas at every day count.
- [ ] No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary.
- [ ] REQ-241's and REQ-242's guarantees hold unchanged: 0 same-row label overlaps, 0 label/mark overlaps, the annotation clear of every neighbour in its strip.

## Implementation Summary

**What was done:** The Durations axis domain is anchored to whole UTC days — first completion floored to its UTC midnight, ending at the midnight after the last — with day buckets (Panel B/C bars, the slowest-day annotation, the hover day-nearest rule) centred on each day's noon, and the outermost bars clamped inside the plot for the >280-active-day case where the 4-unit minimum bar width exceeds the day slot. The Go label planner (`durationLabelTimeRange`) floors/ceils identically, so renderer and planner share one domain; a new mark-position agreement assertion fails if either side's domain ever drifts alone. Axis end label and aria-label keep naming the last *active* day. Verified live in Chromium 141.0.7390.37 headless at 1/2/14 active days: zero bars or annotations outside the plot area [54, 1182] (before: bar x=−5184.4 at one day).

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — day-anchored domain, `DURATIONS_DAY_MS`, `durationsDayCentreX`/`durationsBarLeftX` helpers, noon-centred clamped bars; annotation and hover read the same centres
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — `TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot` (real `renderDurationsView` over a DOM stub at 1/2/14/400 active days: bars inside plot, one annotation inside plot, every Panel A mark at the Go-planned x) + fixtures; retitled the annotation sweep's off-plot rationale
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified) — `durationLabelTimeRange` floors/ceils to the same UTC-day domain (D-01)

*Integrated by orchestrator from builder hand-back; merge range `695b420..1cb897f`.*

## Decisions

Transcribed from the builder hand-back; D-01's resolution is the orchestrator's.

- **D-01 (ESCALATE → resolved by orchestrator: ACCEPTED): `durations.go` edited outside the declared write set.** The scope excluded it to avoid colliding with REQ-252 — but REQ-252 is gated *behind* this REQ, so no live collision exists, and criterion 4 (REQ-241/242 guarantees hold in the render) is unsatisfiable without moving the planner's domain in lock-step: a one-day board compresses rendered marks ~8x against planned label positions. **Value:** the fix is actually complete — planner and renderer share one domain, pinned by the mark-agreement assertion. **Risk:** none live; reverting commit 61e6c1a alone turns the new assertion red rather than silently reintroducing the defect. The builder followed CLAUDE.md's don't-stall rule (challenge + resolution recorded, isolated loud commit), which the orchestrator endorses; Scope and `write_set` extended accordingly.
- **D-02 (DECIDE & STATE): day buckets centre on noon, not midnight** — a bar at its floored midnight still straddles the previous slot and day one's bar is half off-plot; bar, annotation and hover read one helper so they cannot disagree.
- **D-03 (DECIDE & STATE): outermost bars clamped** past ~280 active days where the 4-unit minimum bar width exceeds the day slot — the criterion says *every* day count; the clamp moves only first/last-day bars, by at most half a bar, where adjacent bars already overlap by design.
- **D-04 (DECIDE & STATE): axis end label and aria-label name the last active day**, not the domain's exclusive end — "to 19 Aug" over a board whose last work was 18 Aug would be a small lie introduced by the fix.
- **D-05 (DECIDE & STATE): the probe drives the full real `renderDurationsView`** over a DOM stub rather than sliced functions — the defect lived in the interaction of domain, bucket placement and bar sizing, which slicing hides. Precedent: `timelineForecastDomStub`.

## Qualification

Passed — 3 files verified in merge range `695b420..1cb897f` (+263/−20), all four acceptance criteria traced (floor/ceil visible on both the JS renderer, `board-durations.js:261-262`, and the Go planner, `durationLabelTimeRange`; noon-centring via `xOfEpoch(dayEpochMs + DAY/2)`; clamp via `durationsBarLeftX`; no diff lines on `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s boundary), P-A-U audited (no debug artifacts; probe output uses the established `process.stdout.write` convention). D-01 scope extension accepted and recorded; Scope + `write_set` extended.

## Review

**Overall: 96%** | 2026-08-18T18:22:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Requirements check, per criterion** (all coordinates in SVG user units against plot area [54, 1182]; every measured number taken on Chromium 141.0.7390.37 headless via Playwright 1.56.1 — the same build the builder recorded — with `location.href` returned in the same evaluate call per the board prime's convention):

- [x] **Every Panel B bar inside the plot at every day count, incl. 1 and 2 active days** — *reproduced by execution* on eleven boards: 1 day (bars [606, 630]), 2 days ([324, 912]), 3 days, 4 days across a month boundary and across a year boundary (both [183, 1053]), a single sample, all-samples-within-one-hour, a completion at exactly 00:00:00Z ([230, 1006] — JS `floor(end/day)*day + day` and Go's unconditional `+24h` agree), 283 and 400 active days (both exactly [54, 1182] — the D-03 clamp engaged), and the REAL repository board (228 samples, 24 active days: bars [55, 1181]; the Instance-1 bar formerly at x=37.1 now starts ≥54). Zero bars outside the plot on any board, Panels B and C both.
- [x] **Slowest-day annotation on canvas at every day count** — *reproduced by execution*: exactly one mid-anchored annotation per board, anchor inside [54, 1182] on all eleven, zero bounding-rect intersections anywhere.
- [x] **No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary** — *verified by diff* (zero lines touch either); the only `describeAtPointer` change is the day-nearest rule reading noon centres — verified live at slot boundaries on the 3-day board; Panel A region still describes a sample.
- [x] **REQ-241/242 guarantees hold unchanged** — *reproduced by execution*: all eight named durations label/annotation tests pass (exit 0, un-piped), full package suite green; live confirmation on the real board: 0 same-strip label overlaps, 0 label/mark overlaps, 0 annotation collisions. The REQ-242 sweep's bounds were kept, only their rationale retitled, traced in the comment and its own commit.

**Class-closure attack — the mechanism held in all three drift directions,** each *reproduced by execution* in a scratch worktree at the merge commit:
- Reverting the JS floor → the new probe FAILS (annotation x=−660.4 off-plot, plus the mark-agreement message on every mark).
- Reverting the Go planner's floor/ceil alone → FAILS with the agreement message in the opposite direction.
- Removing the D-03 clamp → the 400-day subtest FAILS (bar edges outside [54, 1182]) — the clamp is load-bearing and pinned.

**Important findings:** None

**Minor findings:** 2 (report only)
1. Stale restatement, `durations.go:344-348`: `durationLabelPlotX`'s comment still cites the renderer's `timeSpan || 1` guard this diff removed (the `domainSeconds <= 0` branch is now unreachable). Found by the Step-6 restatement sweep; the rest of the corpus (both primes, board-guide, test comments) carries no other statement of the old first-instant domain.
2. D-04 (end label/aria-label name the last *active* day) is delivered — verified live on three boards — but no test pins it; a future edit could silently revert it.

**Nit:** at very high day counts with the slowest day at the domain's left edge, the annotation's centred text box extends into the axis gutter (anchor inside the plot, zero collisions measured); the aria-label reads `firstCompletionMs` where the visible start label reads `timeStart` — same rendered day, two expressions of one value.

**Acceptance:** Pass — eleven scratchpad-fixture boards plus the real repository board generated, rendered and *looked at* (screenshots inspected), geometry measured in the live DOM, hover exercised at slot boundaries, month ticks verified across month- and year-boundary domains.

**Suggested testing:** (1) regenerate the maintainer's large reporting board (560 REQs / 42 active days) and eyeball Panel B's left edge — the environment this repo cannot reproduce; (2) a one-line probe assertion on the end label's text would pin D-04.

**Scope:** `scope-drift.sh` exits 0; the D-01 extension is orchestrator-accepted and recorded; no other file drifted; all new names are two-plus words and greppable.

**Follow-ups created:** None · **Sweeps appended to:** None

**Reviewer-recommended disposition: Approve** (Acceptance Pass, 96% ≥ 75%; both Minors are comment/lock-in-hardening level).

*Reviewed by review-work action (independent adversarial review, orchestrated mode; merge range `695b420..1cb897f`)*

## Lessons Learned

**What worked:** The mark-position agreement assertion is the class-closure this board's geometry work had been missing — it fails in *both* drift directions (JS-only revert, Go-only revert), so renderer and planner cannot silently become two definitions of the domain again. Sweeping day counts to 400 caught the hole behind the instance (the 4-unit minimum bar width overhanging past ~280 days), which the floor/ceil fix alone did not close. Driving the full real `renderDurationsView` over a DOM stub found what sliced-function tests hide.

**What didn't:** The captured write set excluded `durations.go` on a collision theory that was empty (REQ-252 is gated behind this REQ), and criterion 4 was unsatisfiable without it — the D-01 escalation cost a scope negotiation mid-build. When a REQ's own guidance says "check the other readers of this domain", the files encoding those readers belong in the write set from the start.

**Worth knowing:** Day buckets centre on noon, not midnight — a bar at its floored midnight straddles the previous slot. The axis end label deliberately names the last *active* day, not the domain's exclusive end (D-04, currently unpinned by any test). `durationLabelPlotX`'s zero-width-domain guard is now unreachable and its comment stale (review Minor 1).

## Orientation

Now Panel B stays on canvas at every day count — the Durations axis domain is anchored to whole UTC days, shared verbatim between the JS renderer and the Go label planner, and pinned by a both-directions agreement assertion. Lives in the board's Durations view (`skills/do-work-board/tools/queue-kanban/`: `web/board-durations.js`, `durations.go`). [MAP CHANGED] — the Durations domain is now a single day-anchored definition with a fourth pinned geometric guarantee (mark-position agreement, structural), joining the three kinds catalogued in the previous session's architectural note. Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` paths still resolve; not stale.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
