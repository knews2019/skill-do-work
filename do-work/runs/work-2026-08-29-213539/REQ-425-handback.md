# REQ-425 Hand-Back

## Branch

`worktree-agent-REQ-425-trailing-window-bound-assumptions`, in
`/home/user/skill-do-work-worktrees/worktree-agent-REQ-425-trailing-window-bound-assumptions`.

Three commits on top of `5384117`:

| Commit | What |
|---|---|
| `3e9f7f8` | RED only — the two failing tests, before any source change |
| `cbe9de0` | GREEN — both fixes in `web/board-timeline.js` |
| `52f8955` | REQ-390's trailing-window test stops claiming a property its case no longer covers |

Nothing pushed, nothing merged. `VERSION`, `skills/do-work/VERSION`,
`actions/version.md`, `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` are
untouched — serial-only integrator files (CLAUDE.md § Before Every Commit,
scoped to the integrating commit).

## Implementation Summary

**Files changed:** `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`,
`skills/do-work-board/tools/queue-kanban/generate_test.go`. Both are in the REQ's
`write_set`; no third file was touched. 507 insertions, 24 deletions.

**Instance 1 — chips collapse on a drained board.** `timelineTrailingWindow` now
hangs each chip's span off `Math.min(Math.max(nowMs, boundStartMs), boundEndMs)`
rather than off `nowMs`. On a drained queue — where `timelineRange` has stopped
pinning `rangeEnd` to `now` — every chip reads as "the last N days **of the
recorded range**": distinct per chip, non-degenerate, and the only reading with
data in it. The clamp-before-settle discipline (REQ-390 D6) is preserved
unchanged; the anchor sits in front of it, not instead of it.

The function now also returns `isClippedByBounds` beside the two endpoints, and
`renderTrailingWindowControls` reads that flag instead of recomputing the clamp
against `nowMs`. That deletes a second opinion about the same arithmetic — and
with it a branch (`candidate.windowEndMs < nowMs`) that the anchor makes
unreachable, so it would have been a dead guard. This follows the prime's REQ-219
lesson: ship a rule's verdict in the payload so a second reader cannot become a
second definition. The control set is still declared in `template.html` alone and
the lit chip is still derived by re-asking each DOM button.

**Instance 2 — the arrows are not inverses.** New pure function
`timelineSteppedScreenfulWindow`, built on top of `timelinePannedWindow` so there
is no second clamp: it takes the panned window only when the pan moved a **whole**
screenful, and returns the window unchanged otherwise. `steppedWindowFor` now
calls it. `renderControlAvailability` already greys an arrow whose step would not
move the window, so a refused step is visible before it is pressed — no new
availability code. `timelinePannedWindow` itself is unchanged, so the keyboard's
fractional pan and the drag keep the continuous clamp.

No calendar-period arithmetic was reintroduced.

## P-A-U

**[PLAN]** Read `_dev/primes/prime-kanban-board.md` in full plus
`lessons-kanban-board.md`, `CLAUDE.md`, and the `general` / `coding-guardrails` /
`communication-style` / `testing` crew members; read archived REQ-390's Plan,
Decisions, Review and Lessons Learned end to end. Confirmed the root cause in
`timeline.go`'s `timelineRange` (it pins `rangeEnd` to `now` only inside
`if (row.WaitOpen || row.WorkOpen)`). Approach: clamp `now` into the bounds to
get an anchor for the trailing span; ship the clipped verdict with the window so
the readout has one source; add an all-or-nothing screenful step for the arrows
alone, leaving the continuous pan for the keyboard and the drag. Tests first,
both as sweeps rather than as the four ages and five positions that happened to
get measured.

**[APPLY]** Code written as planned, strictly inside the two `write_set` files.
One deviation from the reviewer's sketch, deliberately: the reviewer proposed a
separate anchor expression consumed by both the window function and the clipped
flag; shipping the flag *with* the window removes the second consumer entirely.

**[UNIFY]** `git diff --stat 5384117..HEAD` → 2 files, 507 insertions, 24
deletions; `git diff --name-only` lists exactly `web/board-timeline.js` and
`generate_test.go`.

Reviewed each changed file:

- `web/board-timeline.js` — read the whole diff hunk by hunk. Four hunks:
  `timelineSteppedScreenfulWindow` added after `timelinePannedWindow`; the
  `timelineTrailingWindow` comment and body; the `renderTrailingWindowControls`
  clipped-verdict read; `steppedWindowFor`'s callee. Checked that every remaining
  reference to the old behaviour in the surrounding comments was corrected (the
  "still ends at now" sentence in the clamp-before-settle paragraph now says "at
  the anchor"). Checked no orphaned locals: `askedDayCount` was the only variable
  my change stranded and it is removed. Confirmed no other caller reads
  `timelineTrailingWindow`'s return shape (`grep` gives exactly two call sites)
  and that `nowMs` is still a live closure variable used elsewhere.
- `generate_test.go` — two new test functions plus six changed comment lines in
  `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow`. No existing assertion
  was weakened or deleted.

Linters, each with a direct exit status: `gofmt -l .` prints nothing;
`go vet ./...` exit 0; `node --check web/board-timeline.js` exit 0; the canonical
gate's own shell/contract lanes pass inside `maintainer-verify.sh`.

Debug artifacts: `git diff 5384117..HEAD | grep -E '^\+.*(console\.(log|debug|warn)|debugger|TODO|FIXME|XXX|fmt\.Print|t\.Skip\(|\.only\(|process\.exit)'` returns nothing. The
throwaway measurement script and the drained-board fixture live only in the
session scratchpad and are not in the tree.

## Testing

TDD, `tdd: true`. The failing tests were committed on their own as `3e9f7f8`
before any source change, so the red-green evidence is checkable from git.

### Commands and exit codes

All from the worktree, `GOTOOLCHAIN=go1.26.1+auto`, module
`skills/do-work-board/tools/queue-kanban`. No command's status was read through a
pipe.

| Command | Exit |
|---|---|
| `go test -count=1 -run '^TestJavaScriptBehaviorTimelineTrailingWindowsSurviveADrainedQueue$\|^TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds$' .` (before the fix) | **1 — the RED** |
| same two, after the fix | 0 |
| `go test -count=1 -run '^TestJavaScriptBehaviorTimeline' .` | 0 |
| `go vet ./...` then `go test -count=1 ./...` | 0 |
| `QUEUE_KANBAN_BROWSER=<Chrome 151> go test -count=1 -run '^TestBrowserBehaviorTimeline\|^TestTimeline' .` | 0 — 21 browser probes + the Go timeline unit tests, 38 PASS lines, 0 FAIL |
| `bash _dev/tests/maintainer-verify.sh` (repo root, browser lane in its default skipped state) | **0 — "Maintainer verification passed."** |

Browser build: **Google Chrome for Testing 151.0.7922.174**. The container's
`/opt/pw-browsers` Chromium 141 was not used.
`TestBrowserBehaviorCompletionCompanionsKeepReadableContrast` was not run — it is
the known container-local failure at HEAD and is not in the timeline family.

### RED, verbatim

`TestJavaScriptBehaviorTimelineTrailingWindowsSurviveADrainedQueue`:

```
    generate_test.go:7290: trailing-window chips collapsed onto the one-hour zoom floor on a board with 95 days of range:
        	1d idle: chip 1 settled on a 1.00h window at 2026-07-10T23:00Z
        	...
        	7d idle: chip 7 settled on a 1.00h window at 2026-07-10T23:00Z
    generate_test.go:7297: two trailing-window chips produced the same window, so pressing the second lights the first:
        	7d idle: chips 1 and 7 share 2026-07-10T23:00Z -> 2026-07-11T00:00Z
        	...
    generate_test.go:7304: trailing-window chips landed somewhere other than the last N days of the recorded range:
        	1d idle: chip 1 gives 2026-07-10T23:00Z -> 2026-07-11T00:00Z, want 2026-07-10T00:00Z -> 2026-07-11T00:00Z
        	1d idle: chip 7 gives 2026-07-05T00:00Z -> 2026-07-11T00:00Z, want 2026-07-04T00:00Z -> 2026-07-11T00:00Z
        	...
    generate_test.go:7311: the window's own clipped verdict disagrees with whether the bounds cut it short:
        	-20d idle: chip 90 reports isClippedByBounds=undefined, want true
        	...
    generate_test.go:7316: the four board ages the review measured:
        	3d idle  chip 1  1.00h  2026-07-10T23:00Z -> 2026-07-11T00:00Z
        	10d idle  chip 1  1.00h ...   10d idle  chip 7  1.00h ...
        	40d idle  chip 1/7/30 all 1.00h
        	100d idle chip 1/7/30/90 all 1.00h
--- FAIL: TestJavaScriptBehaviorTimelineTrailingWindowsSurviveADrainedQueue (1.88s)
```

`TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds`:

```
    generate_test.go:7509: timelinePannedWindow moved the window by part of a screenful:
        	a partial screenful from the right bound (step 1, 48.00h of room): moved 48.00h of the window's own 168.00h
        	a partial screenful from the left bound (step -1, 48.00h of room): moved -48.00h of the window's own 168.00h
        	window opening at day 0.5 (step -1, 12.00h of room): moved -12.00h of the window's own 168.00h
        	...
    generate_test.go:7515: timelinePannedWindow is not its own inverse:
        	a partial screenful from the right bound (step 1): press then unpress drifted -120.00h
        	a partial screenful from the left bound (step -1): press then unpress drifted 120.00h
        	...
    generate_test.go:7540: the arrows are wired to timelinePannedWindow; the seven named positions measured:
        	mid-range  step +1  moved=168.00h  drift=0.00h
        	one span from the right bound  step +1  moved=168.00h  drift=0.00h
        	a partial screenful from the right bound  step +1  moved=48.00h  drift=-120.00h
        	flush against the right bound  step +1  moved=0.00h  drift=-168.00h
        	flush against the left bound  step -1  moved=0.00h  drift=168.00h
--- FAIL: TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds (1.83s)
```

Both reproduce the review's numbers exactly: -120.00h a partial screenful from
the right bound, -168.00h flush against it, 0.00h mid-range; and the 3 / 10 / 40 /
100-day chip-collapse table.

### GREEN

```
--- PASS: TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow (1.83s)
--- PASS: TestJavaScriptBehaviorTimelineTrailingWindowsSurviveADrainedQueue (1.80s)
--- PASS: TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds (1.89s)
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.560s
```

### The generalisation sweeps

Neither test asserts the measured cases; both assert a property over a swept
space, and both fatal if the sweep did not sweep.

**Chips — 141 board ages × the shipped chip set.** `idleDays` from **-20**
(a live board whose forecast and padding put the range end past `now`) to
**+120**, against a fixed 95-day recorded range. The chip values are read out of
the generated page with `data-timeline-period="([^"]*)"`, so a chip added to
`template.html` is swept too and nothing in the test restates the control set.
Per sample, per chip: no window at or below `TIMELINE_MIN_SPAN_MS`; the windows
pairwise distinct; the window exactly
`[max(anchor - N days, rangeStart), anchor]` (or the bounds for `all`); and
`isClippedByBounds` equal to whether the range start cut the span short. Vacuity
guards: the sweep must report 141 ages and the page's own chip count, and the
clipped verdict must come out both true and false somewhere in it.

Measured on the same 95-day range, before → after:

```
              BEFORE (5384117)                          AFTER
idle=  0d  distinct=5/5  atFloor=0            distinct=5/5  atFloor=0
idle=  3d  distinct=5/5  atFloor=1            distinct=5/5  atFloor=0
idle= 10d  distinct=4/5  atFloor=2            distinct=5/5  atFloor=0
idle= 40d  distinct=3/5  atFloor=3            distinct=5/5  atFloor=0
idle=100d  distinct=2/5  atFloor=4            distinct=5/5  atFloor=0
idle=120d  distinct=2/5  atFloor=4            distinct=5/5  atFloor=0
```

After the fix the five spans are 24.00h / 168.00h / 720.00h / 2160.00h / 2280.00h
at **every** idle value in the sweep.

**Arrows — 167 window positions × both directions, plus 7 named positions.** A
7-day window walked across a 90-day range in half-day steps, from flush against
the left bound to flush against the right. Per position, per direction: a press
either moves a whole screenful or does not move at all (never a partial one);
wherever it moves, the opposite press returns exactly; wherever a whole screenful
of room exists it **must** move — the half that stops "refuse the partial step"
from being satisfied by refusing everything; the window is never resized; and the
step never lands outside the bounds. Vacuity guard: the sweep must report 167
positions and must contain both presses that moved and presses that refused.

The probe **follows the arrows' own call site** rather than naming a window
function: it slices `steppedWindowFor` out of the shipped bundle, reads which
`timeline*` function it returns, and drives that one (REQ-305 — a probe that
calls the function under test directly cannot hold its call site). Re-pointing
`steppedWindowFor` back at `timelinePannedWindow` turns the test red and the
failure text names the wrong function.

Named positions, before → after (`inverse` is "either the press did not move, or
press-then-unpress returned exactly"):

```
                                       BEFORE          AFTER
mid-range                          moved=168h YES  moved=168h YES
one span from the right bound      moved=168h YES  moved=168h YES
partial screenful from right bound moved= 48h NO   moved=  0h YES   (refused)
flush against right bound          moved=  0h YES  moved=  0h YES   (refused)
partial screenful from left bound  moved=-48h NO   moved=  0h YES   (refused)
flush against left bound           moved=  0h YES  moved=  0h YES   (refused)
```

### Mutation table

Each mutation applied to `web/board-timeline.js`, the two new tests run, the file
restored with `git checkout --`. Seven mutations, seven kills.

| # | Mutation | Killed by |
|---|---|---|
| M1 | `anchorEndMs = nowMs` (revert instance 1) | chips on the floor + shared windows + wrong window |
| M2 | `timelineSteppedScreenfulWindow` always returns the panned window | partial steps + not its own inverse |
| M3 | it refuses every press | vacuity guard: moved on 0 presses |
| M4 | `isClippedByBounds: false` | clipped verdict disagrees |
| M5 | `timelinePannedWindow` stops clamping entirely | vacuity guard: refused 0 of 348 presses |
| M6 | `timelinePannedWindow` clamps only the left bound | **stepped outside the board's range** |
| M7 | `steppedWindowFor` re-pointed at `timelinePannedWindow` | partial steps + not its own inverse, naming `timelinePannedWindow` |

M6 exists because M5 dies on the vacuity guard before reaching the bounds
assertion; without M6 that assertion would be a guard no mutation can break.

### Render evidence

The prime's rule — a chart's correctness is partly a claim about pixels, so
generate a board and look at it. Because the defect only appears on a **drained**
board and this repo's own queue is not drained, this needed a fixture: six
archived REQs, first created 135 days ago, last completed **40 days ago**, no
queue and no working directory, generated with `queue-kanban generate
--repo-root`. Driven in **Chrome for Testing 151.0.7922.174**, headless,
1600×900; `location.href` returned alongside every measurement, and the two runs
used two separate output directories (`drained-board/` and
`drained-board-before/`) so the URLs distinguish them.

BEFORE, at `file:///…/drained-board-before/probe.html`:

```
pressed Last day      lit=Last day  state=part of last day  prev=off next=off  2026-07-23 09:24 → 10:24 UTC
pressed Last 7 days   lit=Last day  state=part of last day  prev=off next=off  2026-07-23 09:24 → 10:24 UTC
pressed Last 30 days  lit=Last day  state=part of last day  prev=off next=off  2026-07-23 09:24 → 10:24 UTC
pressed Last 90 days  lit=Last 90 days  state=part of last 90 days             2026-06-01 → 2026-07-23
Last 7 days → ‹ → ›   lit=Last day, window never moved at any step
```

Three of the five chips land the reader on the **same one-hour empty window**
with both arrows dead and no way out except another chip, and each one lights a
chip they did not press.

AFTER, at `file:///…/drained-board/probe.html`:

```
pressed Last day      lit=Last day      state=last day      2026-07-22 10:24 → 2026-07-23 10:24 UTC
pressed Last 7 days   lit=Last 7 days   state=last 7 days   2026-07-16 10:24 → 2026-07-23 10:24 UTC
pressed Last 30 days  lit=Last 30 days  state=last 30 days  2026-06-23 10:24 → 2026-07-23 10:24 UTC
pressed Last 90 days  lit=Last 90 days  state=last 90 days  2026-04-24 10:24 → 2026-07-23 10:24 UTC
pressed All days      lit=All days      state=all days      2026-04-15 15:12 → 2026-07-23 10:24 UTC
Last 7 days → ‹ → ›   back to 2026-07-16 10:24 → 2026-07-23 10:24 UTC, Last 7 days lit again
```

Every chip lights itself, names itself, and gives a distinct window; the arrow
round trip returns to exactly the window it started from and re-lights the chip.
`›` is correctly dead on every chip (the windows end at the range end) and `‹` is
dead on Last 90 days and All days, whose own width exceeds the room left to their
left — that is the refusal, reported before it is pressed.

### Cross-REQ test changes

One, and no assertion changed: `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow`
(REQ-390) case (4) drives `timelinePannedWindow`, which the arrows no longer
call. The assertions still hold and still earn their place — they pin the
continuous pan the keyboard and the drag rest on, mid-range where its clamp
cannot fire — so only the two comments that called it "the arrows" were
corrected, and they now point at the new bounds test. Commit `52f8955`.

## Integration Seams

**S1 — `web/template.html` ships two sentences that are now half-true on a
drained board, and they are outside this REQ's `write_set`.** Flagged rather
than written, per the brief.

- Line 423, the `.timeline-hint` paragraph: *"Pick Last day, Last 7, 30 or 90
  days for a window ending at now"*. On a drained board the window ends at the
  range end. Suggested wording: *"…for a window ending at the board's most recent
  activity"*.
- Lines 373-376, the HTML comment above the control group: *"each chip ends the
  window at NOW and reaches back the span it names"* — same correction.

The reader-visible state span (`#timeline-period-state`) and the panel's
`aria-label` are both still accurate; only these two need the pass. REQ-390's
review already recorded a batch of surviving "period" wording for *"one pass
whenever the board's prose is next touched"* — these belong with it.

**S2 — no `_dev/primes/prime-kanban-board.md` change is needed.** The prime
documents neither the trailing-window maths nor the step arrows; its `## Traps`
entries all still hold. The drained-board condition is worth a lessons-satellite
line on archive, and REQ-390's own Lessons Learned already states it.

**S3 — the write surfaces count in `CLAUDE.md` is unaffected.** Nothing here
writes anything.

## Decisions

**D-01 — The arrows refuse a step they cannot make in full, rather than
mirroring the clamp. ESCALATE.**

The REQ names two candidate remedies. The second — a back step that moves only
as far as the clamped forward step did — **cannot be written as a pure function
of the window**, and that is not a preference. Every window within a screenful of
the right bound steps to the *same* bound-flush window, so nothing in that window
identifies which of them to return to; mirroring the clamp would need the step to
remember where it came from, which means state in a view whose whole design
(REQ-235) is to derive rather than store. The third option — stop clamping the
step, so it is trivially its own inverse — walks the reader past the end of the
board, and the sweep's `escapedTheBounds` assertion exists to rule it out.

So the choice is really: refuse the partial step, or keep the defect. I refused
it, and put the refusal in a new function used **only by the toolbar arrows**.
`timelinePannedWindow` is untouched, so the keyboard's fractional pan
(`TIMELINE_PAN_FRACTION`, 0.15) and the drag keep the continuous clamp.

*Value:* `‹` and `›` become the reader's undo everywhere both are enabled, which
is the property the deleted calendar-period test held and REQ-390's replacement
lost. The refusal is announced rather than silent — `renderControlAvailability`
already greys an arrow whose step would not move the window, so the arrow goes
dead *before* the press instead of moving the reader somewhere they cannot get
back from. This needed no new availability code, which is itself evidence the
remedy fits the existing model.

*Risk:* low and reversible — a one-line revert of `steppedWindowFor`'s callee.
The cost the reviewer worried about, losing the last partial screenful of range,
is smaller than it looked: that stretch stays reachable by the drag, by the
keyboard's fractional pan, and by Fit all, All days and the date fields. What is
genuinely given up is reaching it *with the arrows*, at the current zoom, in one
press. The behaviour a reader will feel is an arrow that greys out slightly
earlier than before — visible on the drained fixture above, where `‹` is dead on
Last 90 days and All days.

**D-02 — `timelineTrailingWindow` returns its clipped verdict with the window.
DECIDE & STATE.**

The reviewer's sketch had the anchor computed once and consumed twice — by the
window maths and by the readout's clipped flag. Shipping the verdict *with* the
window removes the second consumer, which is the prime's own REQ-219 lesson (ship
a rule's verdict in the payload so a second reader cannot become a second
definition) applied to the two readers of one clamp. It also removed a branch
(`candidate.windowEndMs < nowMs`) that the anchor makes unreachable — a dead
guard, which REQ-323's lesson says not to leave standing.

**D-03 — Both tests sweep rather than pin the measured cases. DECIDE & STATE.**

The review measured four board ages and five window positions. A fix keyed on
those is not a fix — a board is idle for however long it is idle. Both tests
sweep (141 ages, 167 positions), both carry a vacuity guard that fatals if the
sweep did not sweep or if a verdict never came out both ways, and the chip test
reads the control set out of the shipped page rather than restating
`1/7/30/90/all`.

**D-04 — Scope was NOT widened to `web/template.html`. DECIDE & STATE.**

Two shipped sentences describe the old behaviour (S1 above). The fix does not
need them — the behaviour is correct and the derived readout stays honest — so
they are flagged here rather than written, which is what the brief asked for.

## Discovered Tasks

- **The keyboard's ← and → have the same drift as the arrows had.**
  `timelineKeyboardWindow` pans by `TIMELINE_PAN_FRACTION` (0.15) through
  `timelinePannedWindow`, so near a bound a → press moves a partial fraction while
  the ← that follows moves a full one. Same shape as instance 2, deliberately left
  alone: the continuous creep is what keeps the last partial screenful reachable
  after D-01, and a key press has no disabled state to announce a refusal with. If
  it is worth fixing, the remedy is different from D-01's — probably clamping the
  *return* press against the same bound rather than refusing.
- **`web/template.html`'s hint paragraph and toolbar comment still say the chips
  end at now** (S1). Outside this REQ's write set.
- **The runtime empty-window message still says "step to another period"**, and
  the control group's `aria-label` is still "Timeline period" — both carried over
  from REQ-390's Minor findings, both still true, both in `template.html`.
- **On a drained board every chip's window ends at the range end, and nothing on
  screen says the now-line is off to the right.** The readout shows the real
  instants, so nothing lies, but a reader on a long-idle board has no cue that
  "Last 7 days" is not the last seven days. A one-clause addition to the
  `#timeline-period-state` span (e.g. "to the last recorded activity") would say
  it; it is a wording decision on a surface this REQ did not own.
