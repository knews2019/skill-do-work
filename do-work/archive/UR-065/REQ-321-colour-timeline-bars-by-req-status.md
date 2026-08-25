---
id: REQ-321
title: "Colour timeline bars by REQ status"
status: completed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-23T01:29:48Z
route: B
completed_at: 2026-08-23T02:12:30Z
commit: d5fc642
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-320]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-23T01:30:20Z
  basis:
    - Route B
    - 4-file write set
    - 7 acceptance criteria
    - dependency depth 3
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/browser_probe_test.go
---

# Colour Timeline Bars by REQ Status

## What

Give every timeline bar its REQ's status colour, using the same semantic tokens the board
cards and the Calendar chips already use.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One custom property per row, `--timeline-status-accent`, set by a
  `[data-status]` block mirroring the calendar chips' exactly. Both segments read it; the
  wait carries `fill-opacity: 0.4` so the phase difference is lightness rather than a second
  hue. The renderer only writes `data-status` and the unrecognized class — the vocabulary
  stays in CSS, and the unrecognized verdict is read from the payload rather than re-derived.
- [x] **[APPLY]:** Confined to the declared write set, plus a new
  `timeline_browser_probe_test.go` (D-02).
- [x] **[UNIFY]:** `git diff --stat` reviewed; `node --check` clean; `go vet ./...` clean;
  `gofmt -l .` empty; no debug artifacts. Files verified — `web/board-timeline.js` (the row
  group's attributes), `web/board.css` (the status block, the two segment rules, the open
  rule's stroke, the legend swatches, one dead token removed), `web/template.html` (the
  two-part legend), `timeline_browser_probe_test.go` (new).

## Why

Today a blocked REQ, a pending REQ and a completed one draw identical grey-and-blue bars.
Status — the thing the rest of the board colours everything by — is invisible here unless
you open a row.

## Context

`board.css` already defines the semantic palette every other view reads: `--accent-pending`,
`--accent-claimed`, `--accent-blocked`, `--accent-done`, `--ink-faint` for cancelled, plus
the matching `--tint-*`. The Calendar's `.calendar-chip[data-status=…]` block is the
reference mapping. The Timeline instead has its own `--timeline-wait` / `--timeline-work` /
`--timeline-projected` trio.

The client already holds each row's status: `requestsById[row.id].status`, which
`renderVisibleRows` reads for the cancelled case and `timelineRowDescription` speaks aloud.
Nothing new is needed in the payload.

`calendarDayBreakdown` spells every status out rather than prefix-matching, so a typo like
`blockd-dependency-cycle` falls through to the unrecognized group instead of being counted
as real blocked work. Same rule here: exact match, and an unrecognized status takes the same
accent the Calendar gives it.

## Detailed Requirements

- **The whole bar takes the status colour** — both segments. Wait vs work is told apart by
  lightness: a pale wash for the wait, the solid accent for the work. This was the user's
  choice at capture, over colouring the work segment alone (which leaves every unclaimed REQ
  colourless) and over a separate per-row status stripe (which leaves the bar status-blind).
- Reuse the existing `--accent-*` / `--tint-*` tokens. Do not mint a second palette for the
  same statuses; a REQ must not read as one colour on a card and another on a bar.
- Map status to class through one pure function, exact-match over the full status
  vocabulary, with an explicit fallback for an unrecognized value. Keep it in step with
  `actions/work-reference.md`'s status vocabulary and `model.go`.
- Cancelled keeps its existing dimmed treatment; broken stamps keep their break marker;
  projected segments keep their hatch — a forecast must never read as measured work.
- The open-span dashed outline still has to mean "running to the now-line" once the fill is
  a status hue.
- The legend is now describing two encodings — hue is status, lightness is wait vs work.
  Rewrite it to say both. It is the view's only colour key.
- The status stays spoken as well as coloured: the row `aria-label` and the table column
  already do this, and must survive.

## Constraints

- Both themes. The dark and light blocks in `board.css` both define the accent tokens;
  check the rendered bars in each rather than trusting the token names.
- Contrast between the pale wait wash and the surface has to survive at a 10px bar height.
  Render it and look, per the prime's rule about pixels.
- Serial with the rest of the `timeline-ux-audit` batch.

## Builder Guidance

**Certainty: Firm on the encoding, exploratory on the values.** The user asked for status
colour "like on the calendar view" and picked the whole-bar encoding from three options at
capture, so which channel carries what is settled. The two lightness values that separate
wait from work are not: pick them, render a mixed-status board in both themes, and adjust
what the render shows. Scope cue: reuse the palette that exists. A REQ must read as the
same colour on its card, its calendar chip and its bar, and a second set of tokens for the
same five statuses would be the thing that later drifts.

## Dependencies

`depends_on: [REQ-320]` — **ordering, not logic.** REQ-321 does not need anything REQ-320
produces; it needs REQ-320 not to be editing `web/board-timeline.js` at the same time. Every
REQ in the `timeline-ux-audit` batch writes that one file, and `write_set` is display-only —
`actions/work.md` computes a `--fan-out` wave from `depends_on` alone and explicitly does not
read `write_set`, `batch`, or the Constraints prose. Without this edge the batch's stated
serial requirement was a sentence nothing enforced, and a `--fan-out` run would have
dispatched four concurrent builders into one 1,100-line file.

**The cost, stated rather than hidden:** a chain gates on terminal *success*, so a `failed`
REQ upstream leaves the rest dependency-blocked until someone edits the chain or resolves the
failure. That is the trade for making the metadata say what the prose says.

## Red-Green Proof

**RED prompt/case:** Generate a board holding at least one REQ in each of pending, claimed,
blocked, completed and cancelled, and open the Timeline tab. Every bar is the same
grey-and-blue pair; the only rows that look different are cancelled ones (dimmed) and
broken-stamp ones (a red break marker). Nothing on the chart distinguishes a blocked REQ
from a completed one.

**Why RED now:** `drawSegment` is called with `timeline-segment-wait` /
`timeline-segment-work` and nothing else; status never reaches a class.

**GREEN when:** each bar carries its REQ's status colour from the shared `--accent-*`
tokens, matching what the same REQ's Calendar chip shows; wait and work are still
distinguishable within one bar; an unrecognized status falls to the same accent the Calendar
gives it rather than to a real status's colour; the legend states hue-is-status and
lightness-is-phase; and the screenshot of a mixed-status board shows the difference at a
glance.

**Validation:** User adjusted — the user asked for "color coded, like on the calendar view"
and chose the whole-bar encoding from the options offered at capture.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — forty uniformly grey rows.

---
*Source: "4. req status should be color coded, like on the calendar view."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The encoding was settled with the user at capture and the tokens already exist,
but where they attach is not obvious: the bars are SVG `fill` on classed rects while the
`--accent-*`/`--tint-*` tokens are used by HTML elements through a `--card-accent` indirection,
and the open-span dashed outline and the projected hatch both already claim visual channels
the status hue has to coexist with. Exploration first, no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Three things decided the shape:

- **The `--tint-*` tokens cannot be the wait colour.** They are 12–14% alpha washes designed
  to tint a card body behind text. At a 10px bar height a 14% wash is not a bar. The wait
  needed to be the same accent at a workable opacity, not a different token.
- **The accents reach elements through an indirection already.** Cards use `--card-accent`,
  calendar chips use `--chip-accent`. A third — `--timeline-status-accent` — keeps the
  mapping in one block per view and means a swatch can read the same property a bar does,
  so a legend swatch cannot show a colour no bar uses.
- **The board already decides what an unrecognized status is.** `generate.go` ships
  `statusUnrecognized` per REQ and `board-calendar.js` consumes it. Re-deriving the
  vocabulary here would have made a second definition of "real status" — REQ-219's lesson,
  recorded in this module's own prime. The renderer consumes the payload's verdict; the
  vocabulary appears once more only as CSS selectors, mirroring the calendar's block.

## Decisions

- **D-01 — Opacity, not a second set of pale tokens.** The REQ asked for wait-vs-work by
  lightness. Eight more variables for the same five statuses is the drift the indirection
  exists to prevent, so both segments read one accent and the wait draws at `fill-opacity:
  0.4`. The open-segment rule keeps its own `0.35` and its dashed outline, ordered after so it
  wins for an open span — the dash is what carries "open" now that both halves share a hue.
- **D-02 — A new file, `timeline_browser_probe_test.go`, outside the captured write set.**
  DECIDE & STATE. What this REQ delivers is a cascade that is entirely CSS: attribute →
  custom property → token → per-theme value → two opacities. A Node probe can asserts the
  class names the renderer writes and nothing about the colours they produce. The browser
  lane already exists for exactly this (`browser_probe_test.go`), and `durations_browser_probe_test.go`
  is the precedent for a per-view file. `write_set` updated in the same edit.
- **D-03 — `--timeline-wait` deleted rather than left in place.** After this change nothing
  referenced it, in either theme block. A token defining a colour no element can show is the
  next reader's false lead.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified) — post-review, F3

**What was done:** Every timeline bar takes its REQ's status colour from the same
`--accent-*` tokens the cards and calendar chips use, resolved through one
`--timeline-status-accent` property per row whose `[data-status]` mapping mirrors the
calendar's block exactly. Both segments read that one accent; the wait draws at 40% opacity,
so hue carries the status and lightness carries the phase. The renderer writes `data-status`
and, from the payload's own verdict, `is-status-unrecognized` — so a typo takes the
unrecognized colour instead of being counted as real blocked work. The legend now states both
encodings in two labelled halves, its status swatches reading the same property the bars do.

## Qualification

Passed — 4 files verified in the diff, 8 requirements traced, P-A-U confirmed.

- Mechanical (`tools/checks/qualify.sh`): OK. One WARN, judged: the new
  `timeline_browser_probe_test.go` has no static reference, which is the script's documented
  exception for test files — Go's test runner discovers it by filename.
- **Substantive:** a new CSS block mapping ten statuses, two rewritten segment rules, a
  two-part legend, and a 260-line browser probe. Not whitespace.
- **Requirements traced:** whole bar takes the status colour → the `[data-status]` block and
  both segment rules, measured; reuse the existing tokens → the probe compares each bar
  against the same status's calendar chip and fails if they differ; one pure mapping,
  exact-match, explicit fallback → the CSS block mirrors the calendar's, with the payload's
  `statusUnrecognized` as the fallback trigger; cancelled keeps its dim → `.timeline-row.is-cancelled`
  untouched; broken keeps its marker → `.timeline-segment-broken` untouched; projected keeps its
  hatch → untouched, and it deliberately does not take a status hue because a forecast must
  never read as measured work; open still reads as open → the dashed stroke now takes the
  status accent and the fill-opacity rule is ordered to win; legend states both encodings →
  the two-part legend; status stays spoken → `timelineRowDescription` and the table column are
  untouched.
- **Flowing:** not applicable — no data-fetch path. The payload gained no field; the
  renderer consumes `request.status` and `request.statusUnrecognized`, both already shipped.
- **Contamination check:** REQ-320 touched `board-timeline.js`, `template.html`, `board.css`,
  `generate_test.go`. This REQ touches the first three plus a new test file — expected overlap
  on the shared surface, and `generate_test.go` was declared at capture but not needed, because
  the probe this REQ needed belongs in the browser lane rather than the Node one (D-02).

## Testing

**Tests run:** `go vet ./...` and `go test -count=1 ./...` in
`skills/do-work-board/tools/queue-kanban`; `bash _dev/tests/maintainer-verify.sh` from the
repo root (`GOTOOLCHAIN=go1.26.1`, `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
**Result:** ✓ All passing — canonical gate exit 0, strict JavaScript and strict browser lanes
included.

**Red-green validation:**
- `TestBrowserBehaviorTimelineBarsCarryTheirStatusColour`: **mutation-verified rather than
  merely passing.** Restoring the pre-REQ fills (`--timeline-wait` / `--timeline-work` on the
  two segments) makes it fail with one line per status —
  `status "pending" paints its bar rgb(47, 111, 192) but its calendar chip #8f5e10; one REQ
  must be one colour` — which is the defect this REQ exists to remove, stated by the test.
- The probe asks a real engine for computed fills over all ten statuses in the board's
  vocabulary plus a deliberate typo, and refuses to pass on an unresolved property: a fill of
  `none`, empty, or fully transparent is a hard failure, because a cascade that resolves to
  nothing would otherwise satisfy every "these two agree" assertion at once (REQ-291's lesson).
- Five assertions, each removing a different way to be wrong: the bar matches its own calendar
  chip; both halves share a hue and differ in opacity; the five statuses a reader must tell
  apart are mutually distinct (without which "everything one colour" passes); grouped statuses
  share their group; and the typo reaches the unrecognized colour **through the unrecognized
  class**, not through a prefix match on its name.
- `TestTimelineStatusRowMarkupMatchesTheProbe` pins the probe's fixture to the shipped
  renderer, so it cannot keep passing against markup the board stopped producing — REQ-305's
  lesson, applied to a probe that necessarily builds its own rows.

**Render evidence, both themes.** Generated a board from this repo's archive and drove it in
headless Chromium —
`file:///tmp/claude-0/-home-user-skill-do-work/32295d3b-538a-57cc-a4d8-2d453777559b/scratchpad/board321/probe-colour.html`,
Chromium 1194 (Playwright build, `/opt/pw-browsers/chromium`, UA `HeadlessChrome/141.0.0.0`),
`location.href` returned with the measurements. The board is **dark-first** — `:root` is the
dark palette and `@media (prefers-color-scheme: light)` overrides it — and
`--force-prefers-color-scheme` is a no-op on this build, so dark was reached with
`--blink-settings=preferredColorScheme=0` and confirmed by `matchMedia` reporting `dark`.

| Status | Light | Dark | Rows in view |
|---|---|---|---|
| completed | `rgb(60, 135, 94)` | `rgb(102, 181, 133)` | 29 |
| pending | `rgb(143, 94, 16)` | `rgb(216, 162, 74)` | 3 |
| cancelled | `rgb(108, 116, 128)` | `rgb(107, 116, 128)` | 2 |
| claimed | `rgb(58, 107, 196)` | `rgb(111, 156, 230)` | 1 |
| blocked | `rgb(189, 81, 56)` | `rgb(217, 122, 89)` | legend only |

Five distinct accents in each theme, and every legend swatch measured identical to the bars
of its status — which is the point of the swatches reading the same custom property.

**New tests added:**
- `TestBrowserBehaviorTimelineBarsCarryTheirStatusColour` — the whole cascade, in an engine.
- `TestTimelineStatusRowMarkupMatchesTheProbe` — the fixture's tie to the renderer.

**Existing tests updated (cross-REQ impact):** none.

## Discovered Tasks

- **The canonical gate can hang indefinitely on `generate-report-image`'s interruption case,
  and the same path would leave a backend spinning on a user's machine.** Found twice
  independently: the REQ-320 review agent sat ~35 minutes in that lane and had to kill the
  stub by hand, and this REQ's gate run stalled 9+ minutes on the identical process before I
  killed it. Both times the live processes were the case script
  (`_dev/tests/prescribed-shell-cases/generate-report-image.sh`), the shipped wrapper
  (`skills/do-work-toolbox/scripts/generate-report-image.sh`), and the fixture's stub
  `imagegen` spinning in `while :; do sleep 0.1; done`. The case TERMs the WRAPPER
  (`kill -TERM "$interrupt_helper_pid"`, line 91) and then `wait`s on it; the wrapper carries
  `trap 'exit 143' TERM` and a backend-signalling helper that kills the backend's process
  group — so the machinery to forward the signal exists and did not fire. Nothing here is
  this REQ's code: it is the interruption path of the report-image script and the probe that
  exercises it.
  **Why it matters beyond the gate:** the same wrapper ships to consumers. An interrupted
  `do-work-toolbox ai-report` would leave the image backend running.
  **Not investigated further on purpose** — characterising exactly which of the two sides
  fails to deliver the signal is its own piece of work, in a subsystem this batch does not
  touch, and stretching REQ-321 to cover it would be the scope creep the guardrails forbid.
  Impact judged `impact-user-visible`: a developer running the repo's only accepted proof
  loses tens of minutes per occurrence, and a consumer gets an orphaned process. Not
  `impact-critical` — no security, data-loss, or production path is involved.
  **Routed at Step 8** to `REQ-325` (`status: pending-answers`), per the non-critical
  discovery flow: the user confirms or discards it through `do-work clarify` rather than it
  auto-queueing on my judgment. The fold-first scan found no home — the queue holds only the
  three remaining timeline REQs and no `sweep: true` file, and this is behaviour rather than
  prose.

## Review

**Acceptance: Partial — Overall 72% — Approve with follow-ups.** Independent review agent,
orchestrated mode, measuring contrast in a real engine in both palettes.

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 85% | 7 of 10 full; mapping shape, legend and the contrast constraint partial |
| Code Quality | 78% | Clean CSS mirroring the calendar; two stale comments, one wrong mechanism claim, an invisible-swatch failure mode |
| Test Adequacy | 65% | Strong on hue (3 of 5 mutations caught); blind to opacity magnitude, contrast, the dark theme, the fallback and the legend |
| Scope | 100% | Exactly the four declared files; the extra test file recorded as D-02 |
| Risk | Low | No security, performance or data surface |
| Acceptance | Partial | Gate exit 0 and the hue works; two measured user-visible defects |

**Its verdict in one line, and it is right:** the status hue is delivered and provable; the
phase encoding it displaced got measurably weaker, and nothing in the suite could see that.

**What it verified rather than accepted:** re-ran the gate itself; measured WCAG contrast for
every status in both palettes; ran five mutations against a scratch tree and reported that
three failed correctly and two passed when they should not; and independently confirmed the
`--timeline-wait` deletion was clean, the cascade does not leak between rows, and the broken
marker, cancelled dim, open dash and projected hatch all still read.

## Review Remediation (pre-archive)

- **F1 — light-theme wait/work fell under the board's own 3:1 floor (2.47:1, from a 3.11:1
  baseline).** The diagnosis was the important part: opacity is one-dimensional and had to buy
  two separations that pull in opposite directions. Both of the review's suggested channels
  are now in:
  - **A full-strength 1px outline on the wait.** It costs nothing against the surface and puts
    a hue-correct boundary between the halves that does not depend on the fills being far
    apart. `.timeline-segment.is-open` replaces it with the dashed variant on specificity, so
    outline-present means wait and dashed means open.
  - **A per-theme `--timeline-wait-alpha`** (0.22 dark, 0.45 light). The two palettes want
    opposite values: these accents are bright on the dark surface, so a low alpha pulls the
    wait off the work while the dark page keeps it visible; on the light surface the same
    accents are dark, so alpha has to stay high or the wait vanishes into the page. One number
    could not serve both, which is why the first attempt could not be fixed by choosing a
    better one.

  **Measured after, against the real body surfaces** (light `rgb(245,247,250)`, dark
  `rgb(12,14,18)` — confirmed by walking up from a bar to the first ancestor that paints, which
  is `<body>`):

  | | wait vs work | wait vs surface | outline vs its own fill |
  |---|---|---|---|
  | light (α 0.45) | 2.68–3.30:1 | 1.51–1.57:1 | 2.68–3.30:1 |
  | dark (α 0.22) | 2.43–3.20:1 | 1.68–2.64:1 | 2.43–3.20:1 |

  Dark's wait-vs-work went from 1.71–1.94:1 at the single 0.45 value to 2.43–3.20:1. Where the
  fills still fall short of 3:1 the outline carries the separation, and the outline reads
  2.43:1 or better against the fill it borders in every status and both themes.

- **F2 — the phase assertion pinned ordering, not readability: it passed at 0.98 and at 0.26.**
  Replaced with real WCAG contrast computed from the composited fills, against two floors that
  mean different things (3:1 for two adjacent marks, 1.3:1 for a fill against the page). **The
  first rewrite re-created the vacuity** — it skipped the contrast check whenever a stroke was
  present, which would have let 0.98 through again. The rule now is: the fills are separable
  **or** the outline is visible *against the wait's own fill*, and since the outline is the
  accent at full strength, a wait approaching full strength loses both clauses at once.
  **Mutation-verified, all four failing:** 0.45→0.98 (`separates its wait from its work by
  1.08:1 of fill`), 0.45→0.10 (`draws its wait at 1.09:1 against the page`), outline removed,
  and the whole `[data-status]` block deleted.
- **F3 — no automated check ever saw the dark palette.** `runBrowserBehaviorProbe` gained a
  flag-accepting sibling, and the status probe now runs as two subtests under
  `--blink-settings=preferredColorScheme=0/1`. It also **asserts the scheme the engine
  actually resolved** rather than the one the flag asked for: a build that ignores the flag
  would otherwise measure one palette twice and report it as two, which is how this would
  quietly rot.
- **F4 — the legend's phase half was painted in the claimed accent (1.02:1 apart in light),
  so the view's only colour key showed one chip twice meaning two different things** — and
  showed `--timeline-work`, a colour no bar paints, one line under a comment claiming that
  could not happen. Repainted in muted ink, which cannot be read as a sixth status. The status
  swatches also gained an ink fallback: an unmapped or renamed status rendered them fully
  transparent, so a key entry disappeared silently.
- **F5 — two comments still described the deleted two-hue scheme.** Both rewritten. Prose-only
  by the Fold-First test, but they are in a file this REQ was already editing and UR-065's
  batch constraint assigns prose a REQ falsifies to that REQ, so they are fixed here rather
  than sent to the prose backlog.
- **Minors fixed:** the row fallback is muted ink rather than the claimed blue (a status-less
  row read as claimed work); the unresolved-property guard parses alpha instead of matching
  the substring `", 0)"`, which also matched any opaque colour with a zero blue channel; and
  the `is-open` comment now says specificity (0,2,0) rather than source order, which is what
  actually decides it.

**Carried, not done:** a second engine and a forced-colors pass, both outside what this
environment offers; and extending `generate_test.go`'s table-driven `.req-card[data-status]`
pin to cover all three copies of the status vocabulary, which is a cross-view concern rather
than this REQ's.

## Lessons Learned

**What worked:** Consuming the payload's `statusUnrecognized` verdict instead of re-deriving
the vocabulary. And routing every swatch through the same custom property the bars read —
which is what made the legend's contradiction (F4) a one-line fix rather than an audit.

**What didn't:**

- **Opacity is one-dimensional, and it was asked to buy two separations at once.** Wait from
  work, and wait from the page. In the light palette those pull in opposite directions, so no
  single alpha exists that satisfies both — the first attempt was not a bad number, it was a
  missing channel. Two channels (an outline, and a per-theme alpha) settled it.
- **A test that asserts an ordering is not a test of a difference.** `waitOpacity <
  workOpacity` passes at 0.98 and at 0.26. The REQ's stated risk was legibility and the probe
  measured everything except legibility. Worse, the mutation I chose to "verify" it — restoring
  the pre-REQ fills — was precisely the one thing assertion 1 already greps for, which is
  REQ-293's lesson arriving one REQ later: **choose the mutation before looking at the
  assertion, or you will choose the one it already catches.**
- **And then the first fix re-created the hole.** Exempting the contrast check whenever a
  stroke was present would have let 0.98 through again. An escape hatch in an assertion needs
  its own condition to be measured — here, that the outline is visible *against the fill it
  borders*.
- **Manual evidence in both themes is not coverage of both themes.** The probe lane launches
  Chromium with no colour-scheme flag, so every existing browser test measures light. The dark
  palette is this board's `:root` base and nothing automated had ever looked at it.

**Worth knowing:** the surface behind a timeline bar is `<body>` — nothing between them
paints — so a contrast measurement must read `document.body`'s background, and the real values
are `rgb(245,247,250)` light and `rgb(12,14,18)` dark, not the `--surface-*` tokens you would
reach for. Computing against the wrong surface gave numbers that disagreed with the review's
by a factor of two, which is how the discrepancy was found.

## Orientation

Status is now visible on the Timeline: every bar carries its REQ's status colour from the
shared accents, with the wait/work phase distinction moved onto lightness plus an outline.
Leaf change in the board's frontend — no new module, no payload field. The board's browser
probe lane gained the ability to run a page under a chosen colour scheme, which is the one
piece other views can reuse. `_dev/primes/prime-kanban-board.md` spot-checked: every path it
references still exists.
