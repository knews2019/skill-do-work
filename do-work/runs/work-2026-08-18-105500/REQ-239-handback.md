# REQ-239 hand-back — Give the Timeline's rows a real focus ring

**Branch:** `worktree-agent-REQ-239-timeline-row-focus-ring`
**Commit:** `7f9ba1e` — `[REQ-239] Give the Timeline's rows a real focus ring`
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-239-timeline-row-focus-ring`

**Verdict: ring, not tint.** A 2px `--accent-claimed` ring renders cleanly on an 18px
SVG row — it is not the ugly result the Open Question feared. The existing tint stays
as a deliberate complement (reason below, with evidence).

## File manifest

| File | Verb | What |
|---|---|---|
| `skills/do-work-board/tools/queue-kanban/web/board.css` | added | `.timeline-row:focus-visible { outline: 2px solid var(--accent-claimed); outline-offset: -2px; }`, plus a comment on the tint rule stating why it stays and a comment on the inset offset. `.timeline-row { outline: none; }` kept, now with a one-line reason. |
| `skills/do-work-board/tools/queue-kanban/generate_test.go` | added | `TestGenerateGivesTimelineRowsTheBoardsFocusRing`, appended as one contiguous block at the very end of the file. Nothing above it reflowed or reordered. |

`git diff --stat` (against `main`): 2 files changed, 53 insertions(+), 0 deletions(-).

## RED (verbatim, before any CSS change)

```
$ go test -run TestGenerateGivesTimelineRowsTheBoardsFocusRing ./...
--- FAIL: TestGenerateGivesTimelineRowsTheBoardsFocusRing (0.24s)
    generate_test.go:3045: Timeline rows carry no ".timeline-row:focus-visible {" rule: a keyboard-focused row has no ring, only the tint
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.552s
FAIL
EXIT=1
```

**Stated plainly, as asked:** this first RED is the "rule not found" form. Absence *is*
the defect here, so it is the correct failure — but it is a missing-selector error, not
a wrong-value error, and I am not dressing it up as more than that. The test does more
than presence once the rule exists: it parses the outline out of both
`.control-button:focus-visible` and `.timeline-row:focus-visible` and requires the row's
**width and token to equal the board's reference ring**, and the row's `outline-offset`
to be negative. Those two assertions can fail with a real value mismatch later, which is
the drift this test is actually pinned against.

## GREEN

```
$ go test -run 'TestGenerateGivesTimelineRowsTheBoardsFocusRing|TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard' ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.677s
EXIT=0
```

Verbose form of the same run:

```
=== RUN   TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard
--- PASS: TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard (0.27s)
=== RUN   TestGenerateGivesTimelineRowsTheBoardsFocusRing
--- PASS: TestGenerateGivesTimelineRowsTheBoardsFocusRing (0.22s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.820s
```

**Row activation is unchanged and `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard`
passes** — see above. This REQ touched no JavaScript at all.

### Browser GREEN — real `Tab` presses, shipped stylesheet

Board rebuilt from the edited CSS (`/tmp/qk-239 generate --repo-root <worktree> --out /tmp/board-239`),
served on `127.0.0.1:8239`, driven with Playwright keyboard events. Every measurement below
asserts `document.getElementById('req239-probe') === null`, i.e. **no injected stylesheet** —
the rule came from the generated page.

Light theme, third row (`REQ-003`), after three real `Tab` presses:

```
{"probeStylesheetPresent":false,"id":"REQ-003","focusVisible":true,
 "outline":"2px solid rgb(58, 107, 196)","outlineOffset":"-2px",
 "rect":{"x":52,"y":301,"width":1754,"height":18}}
```

Dark theme, same path: `"outline":"2px solid rgb(111, 156, 230)"` — the dark-palette
`--accent-claimed`, `outlineOffset: "-2px"`, `focusVisible: true`.

I hit the REQ-233 trap once on purpose to confirm it still bites and then stayed off it:
every focus *move* in this evidence is a real `Tab` keypress. Programmatic `.focus()`
appears only to park focus on the row *before* the one under test, never on it.

### Re-taken on an isolated browser context, with page identity asserted in-call

The readings above were taken on the **shared** Playwright instance, so on the lead's
warning about cross-builder contamination I re-took every one of them. Each measurement
below comes from a page in its **own browser context** (`browser.newContext()`), created
inside the same call that reads it, served from a **port only this run used** (`18239`,
not the `8239` of the first pass) with a query token, and every returned object carries
`location.href` **from the same `evaluate` expression that produced the numbers** — not
a URL checked before navigating.

Identity carried by the readings:

```
"href":  "http://127.0.0.1:18239/index.html?req239=isolated-verify"
"title": "worktree-agent-REQ-239-timeline-row-focus-ring — do-work queue board"
"cssRuleText": ".timeline-row:focus-visible {\n  outline: 2px solid var(--accent-claimed);\n  outline-offset: -2px;\n}"
"probeStylesheetPresent": false
```

The title is itself a discriminator: the board titles itself from its repo root, so
`worktree-agent-REQ-239-…` can only be my build. `cssRuleText` is pulled out of the page's
own inlined `<style>`, so the reading proves *which stylesheet* it measured, not just which
URL.

**Every re-taken number is identical to the first pass. Nothing differed.**

| Measurement | First pass (shared browser) | Re-take (isolated, identity asserted) |
|---|---|---|
| keyboard focus, light | `2px solid rgb(58, 107, 196)`, offset `-2px`, `focusVisible: true`, `REQ-003`, row 18px | identical |
| keyboard focus, dark | `2px solid rgb(111, 156, 230)`, offset `-2px`, `focusVisible: true` | identical |
| pointer click | `focus: true`, `focusVisible: false`, `outlineStyle: "none"`, `hitFill: rgb(241, 244, 248)` | identical (on `REQ-011` rather than `REQ-050` — different row, same verdict) |
| bottom-edge row | `REQ-056`, `rowBottom 973` vs `containerBottom 973.34`, `2px solid`, `focusVisible: true` | identical |
| 45 real `Tab` presses | 45/45 landed on a row with `:focus-visible` and a solid outline; `lostFocusAt: -1`; `maxScrollTop: 375`; order jump `REQ-042 → REQ-017` | identical, plus `everyReadingFromMyPage: true` for all 45 |

The re-take also runs at `deviceScaleFactor: 1` (the shared instance was at 2), so the
ring is 2 *physical* pixels rather than 4 in these captures — the harsher case, and it is
still crisp and continuous. Screenshots below are from the isolated run.
|---|---|
| **GREEN, light theme, shipped CSS** — `REQ-003` focused mid-list | `/tmp/req239-shots/green-zoom.png` (full viewport: `/tmp/req239-shots/req239-green-light.png`) |
| **GREEN, dark theme, shipped CSS** — same row | `/tmp/req239-shots/greendark-zoom.png` (full viewport: `/tmp/req239-shots/req239-green-dark.png`) |
| **Bottom-edge case** — focused row flush against the scroll container's bottom edge | `/tmp/req239-shots/be-zoom.png` |
| **Why the offset is inward** — the same ring at `outline-offset: 2px`, clipped on three sides | `/tmp/req239-shots/p2-zoom.png` |

Re-taken in the isolated context (identity asserted in-call; these are the authoritative ones):

| What | Path |
|---|---|
| **GREEN, light theme** — `REQ-003` focused mid-list | `/tmp/req239-shots/v2-light-zoom.png` (full viewport: `/tmp/req239-shots/req239-verify-light.png`) |
| **GREEN, dark theme** — same row | `/tmp/req239-shots/v2-dark-zoom.png` (full viewport: `/tmp/req239-shots/req239-verify-dark.png`) |
| **Bottom-edge case** — `REQ-056` flush against the container's bottom edge | `/tmp/req239-shots/v2-be-zoom.png` |
| **Why the offset is inward** — same row, `outline-offset: 2px` | `/tmp/req239-shots/v2-out-zoom.png` |

**Appearance:** a crisp, continuous 2px rectangle tight around the row, all four edges
present. At 18px the ring leaves a 14px interior, which is still taller than the 10px
segment bars, so the bar never touches the ring and the row does not read as "boxed in".
The row's `REQ-0xx` label sits comfortably inside. It reads exactly like the `.req-card`
ring elsewhere on the board, just shorter. It is legible in both themes; on the dark
surface the `#6f9ce6` ring is the strongest thing in the rows area, which is the point.

**Clipping — checked, not assumed:**

- **Top row** (`REQ-001`, flush against the scroll container's top border): all four edges paint.
- **Bottom edge**: focused `REQ-056` at `rect.bottom = 973` against `scrollRect.bottom = 973.34` —
  all four edges paint (`/tmp/req239-shots/be-zoom.png`).
- **Left/right**: the row `<g>` spans `x = 52 … 1806`, which is the rows SVG's full viewport
  width, so an outward ring has nowhere to go. Inward, both vertical edges paint.
- **At `outline-offset: 2px` the ring is clipped on three sides, and *which* three depends on
  where the row is.** The rows SVG's own viewport always eats the left and right edges. On the
  top row the scroll container eats the top edge and only the bottom survives, painting a
  divider *under the next row*; on the bottom row the container eats the bottom edge and only
  the top survives, painting a line that **strikes through the previous row's label**
  (`/tmp/req239-shots/v2-out-zoom.png`, where the surviving line cuts across `REQ-055` while
  `REQ-056` is the focused row). Either way it names the wrong row. That is the evidence for
  the inward offset — the same reasoning REQ-233's D-06 used for `.timeline-scroll`, plus the
  SVG-viewport half, which is new here.

**Survives scrolling and the virtualized rebuild:** 45 consecutive real `Tab` presses walked
the rows from `REQ-001` down; the container scrolled (`scrollTop` 0 → 12 → 375) and
`renderVisibleRows` rebuilt every row node in the process. Focus landed on a `.timeline-row`
on all 45 presses, with `:focus-visible` true and a solid outline on every one. No press
dropped focus to `<body>`.

## Ring vs tint — the decision

**Ring, and the tint stays.** Both, deliberately, because they answer different questions.

The measurement that settled it: with the shipped rule in place, a **real mouse click** on a
row reports

```
{"href":"http://127.0.0.1:18239/index.html?req239=isolated-verify",
 "activeId":"REQ-011","focus":true,"focusVisible":false,
 "outlineStyle":"none","hitFill":"rgb(241, 244, 248)"}
```

`:focus` true, `:focus-visible` **false**, `outline-style: none` — so the pointer path draws
no ring, exactly as required — while `.timeline-row-hit` still carries `--surface-2`. The tint
is therefore the *only* signal the pointer path has, and the same rule also serves `:hover`.
Removing it would have taken hover feedback away from every mouse user to solve a keyboard
problem. So:

- **pointer / hover → tint alone**
- **keyboard → tint + ring** (the tint fills the row, the ring bounds it)

Two signals on the keyboard path is the deliberate part: at 18px the tint alone is one shade
of difference and the ring alone is a thin rectangle, and together the row reads as filled
*and* bounded. The reasoning is written into `board.css` above both rules so the next reader
does not have to re-derive it.

I did **not** switch the tint's selector from `:focus` to `:focus-visible`. Covering pointer
focus is its job, and narrowing it would have re-created the gap this REQ closes, one step down.

**Asked a second time, after re-verifying on an isolated browser, the answer is still ring.**
The Open Question's worry was specific and testable — a ring on a few-pixel row inside a
virtualized scrolling SVG will clip or look bad — and it does not survive the render: all four
edges paint on the top row, mid-list, and on a row flush with the container's bottom edge; it
holds through 45 keyboard steps and the row rebuild that scrolling triggers; and it is legible
at `deviceScaleFactor: 1`, where it is only two physical pixels. If it had clipped I would be
proposing the tint here instead — the `outline-offset: 2px` capture is what a bad version
actually looks like, and it is nothing like the shipped one. What the worry *did* earn is the
offset: `-2px` is not a style preference, it is the only value that survives the two clips.

## P-A-U

**[PLAN]** Read `CLAUDE.md`; `_dev/primes/prime-kanban-board.md` (including its `## Lessons`
list); `skills/do-work/crew-members/general.md`, `coding-guardrails.md`,
`communication-style.md`, `testing.md`; the REQ body at
`do-work/runs/work-2026-08-18-105500/REQ-239-brief-body.md`; the archived
`do-work/archive/REQ-233-give-the-timeline-a-keyboard-path-to-zoom-and-pan.md`
(D-06 and `## Lessons Learned` — the programmatic-`.focus()` trap and the
`renderVisibleRows` dead-node trap); `web/board.css` §Timeline; `web/board-timeline.js`
(read-only — row geometry, `TIMELINE_ROW_HEIGHT = 18`, `TIMELINE_BAR_HEIGHT = 10`, the
row `<g>` with `tabindex="0"`); and the existing CSS-assertion shape in `generate_test.go`.

Because the REQ's Open Question genuinely asked whether a ring is right at all, I answered
the two unknowns in the browser *before* writing anything: does Chrome even paint `outline`
on an SVG `<g>`, and does a 2px ring survive the two clips. Both were answered with an
injected probe stylesheet on the *unmodified* build, so no source changed during that pass.
Chrome does paint it — proved with a deliberately obnoxious 4px red probe
(`/tmp/req239-shots/vp-zoom.png`) after an element screenshot had misled me into thinking it
did not. Only then did I write the test, capture RED, and change the CSS.

**[APPLY]** Both edits stayed inside the write set: `web/board.css` and `generate_test.go`.
The test is appended as one contiguous block at the very end of `generate_test.go`; no line
above it was touched, so the sibling builder's axis-formatter work should merge without a
textual conflict. `web/board-timeline.js` was read but never written — no markup change was
needed, so there is no integration seam.

**[UNIFY]**

```
$ git diff --stat
 .../tools/queue-kanban/generate_test.go            | 39 ++++++++++++++++++++++
 .../do-work-board/tools/queue-kanban/web/board.css | 14 ++++++++
 2 files changed, 53 insertions(+)
```

- `web/board.css` — one new rule plus three comments; verified by rebuilding the board and
  reading the computed style off a really-focused row in both themes.
- `generate_test.go` — one new test; verified RED before and GREEN after, and checked that
  it fails for the right reason (it asserts equality against the reference ring, not just
  presence of a string).
- `gofmt -l .` → no output (clean). `go vet ./...` → clean.
- `bash _dev/tests/maintainer-verify.sh` → **exit code 0**, run from the worktree root, not
  piped through `tail`/`head`.
- No debug artifacts: `git status --porcelain -uall` is empty in the worktree. The probe
  stylesheet and the click-swallowing listener existed only in the live browser page, never
  in a file. All screenshots and the built binary live under `/tmp`. Playwright's output root
  is pinned to the main tree, so its first screenshot landed at
  `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/req239-ring-probe-scroll.png`; I moved
  it to `/tmp/req239-shots/` immediately and routed every later capture through the gitignored
  `.playwright-mcp/` directory, moving each one out as it was taken. The main tree has no stray
  PNG (`ls *.png` → no matches) and no modification outside this hand-back file.

## Decisions

- **D-01 — Ring, not a better tint.** The Open Question's alternative was a stronger tint on
  the grounds that a ring might clip badly on an 18px virtualized SVG row. Rendered, it does
  not: all four edges paint on the top row, on a mid-list row, and on a row flush with the
  scroll container's bottom edge, and it survives 45 keyboard steps through the virtualized
  rebuild. The reason to prefer a tint did not survive contact with the render, so the
  requirement's presumed means is also the right one.
- **D-02 — `outline-offset: -2px`, not `+2px`.** Not copied from REQ-233's `.timeline-scroll`;
  re-derived and screenshot-checked. A row has two clips, not one: the scroll container above
  and the rows SVG's own viewport on the left and right, since the hit rect is `width="100%"`.
  At `+2px` only the bottom edge survives.
- **D-03 — The tint stays, on `:focus` (not `:focus-visible`).** It is the pointer path's only
  feedback and it doubles as the hover rule; a real click measures `:focus` true /
  `:focus-visible` false / `outline-style: none`. Stated in the CSS itself.
- **D-04 — `.timeline-row { outline: none; }` kept.** It now means "no ring for a plain
  `:focus`", which is a live statement rather than dead weight: it keeps the pointer path
  ring-free even in a browser that draws a UA ring for `:focus`. Given a comment saying so.
- **D-05 — The test compares against the board's reference ring instead of hard-coding
  `2px var(--accent-claimed)`.** The requirement is "the same weight and the same token as the
  rest of the board", so the test asserts that relation. Hard-coding the literal would let the
  two drift apart while the test still passed.

## Discovered Tasks

- **Tab order is not monotonic across the Timeline's virtualized rebuild.** During the 45-press
  walk, focus went `… → REQ-041 → REQ-042` at `scrollTop 375`, and the next `Tab` landed on
  **`REQ-017`** with `scrollTop` back to `0`. Focus is never lost and the ring is always drawn,
  so nothing here is broken in the REQ-239 sense — but a keyboard user tabbing down the list
  gets silently teleported backwards once the container scrolls, because `renderVisibleRows`
  rebuilds the row nodes and the new DOM order does not continue the old one. This is
  pre-existing (`web/board-timeline.js`, outside my write set) and is the same
  dead-node/rebuild hazard REQ-233 recorded, surfacing in the sequential-Tab path rather than
  the arrow-key path. Worth its own REQ.
- **Playwright's screenshot root is the main tree, not the worktree.** Any builder taking
  browser evidence will write a PNG into `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/`
  unless they route it through the gitignored `.playwright-mcp/` and move it out. That is very
  likely how the earlier stray PNG got there, and it will keep happening. A line in
  `_dev/primes/prime-kanban-board.md` § Lessons would cost one sentence.
- **Browser evidence in a fan-out needs page identity inside the reading.** The shared
  Playwright instance is one page object several builders navigate, so a DOM measurement can
  silently describe a sibling's board and look entirely reasonable. Two cheap habits fix it and
  both are one line: serve on a port unique to the run, and return `location.href` (plus
  something only your fixture has — the board titles itself from its repo root, so
  `document.title` alone separates worktrees) *from the same `evaluate` that produces the
  numbers*. Isolating further is `browser.newContext()` in the same call. Worth a line in
  `_dev/primes/prime-kanban-board.md` § Lessons next to the `.focus()` / `:focus-visible` trap,
  since both are "the measurement lied and looked fine".

## Integration seams

**None.** No markup change was needed in `web/board-timeline.js`; the row `<g>` and its
`.timeline-row-hit` rect already carry everything the rule needs. Nothing outside the write
set was touched, and `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`
and `CHANGELOG.md` are untouched.

## Merge note

My branch was cut before REQ-240 merged into `main`, so `generate_test.go` will conflict at
the append point as you predicted. `TestGenerateGivesTimelineRowsTheBoardsFocusRing` is one
contiguous block at the very end of the file and nothing above it moved, so taking both sides
and keeping mine last resolves it. It depends only on `generateLiveSite`,
`sliceBalancedBlockAfter` and `regexp`, all of which predate both branches, and it reads only
the generated stylesheet — REQ-240's axis-formatter work in `web/board-timeline.js` cannot
affect it. `web/board.css` should merge clean: my change is additive, inside the Timeline
block, and no sibling in this batch has that file in its write set.
