# REQ-626 Exploration — board page and lens URLs

Explore agent findings (read-only agent; saved by the orchestrator). Paths are under `skills/do-work-board/tools/queue-kanban/` unless they start with `_dev/`.

## Critical for the builder

- **Writing a fragment breaks the browser probe harness.** `browser_probe_test.go:366` (load wait) and `:481` (`evaluateInPage`, every measurement) fail unless `location.href` ends in `/probe.html`; `durations_browser_probe_test.go:567` has the same check. Once a view click writes `#timeline`, every later evaluate is fatal (activity_scroll_*, timeline_scroll_*, timeline_browser_probe_*, durations_*, `citations_test.go:189`, `generate_test.go:1041-1077`). These run only with `QUEUE_KANBAN_BROWSER_PROBES` on (the queue-kanban-browser heavy lane). Never write a fragment on load when the hash is empty; the href checks must accept a trailing fragment.
- **Some `#...` URLs are not views.** Skip link `href="#board-main"` (`template.html:13`); rendered detail bodies get heading ids (`render.go:28`, `WithAutoHeadingID`), so an author's `[x](#lessons-learned)` link changes the hash. Ticket links use `href="#"` but `board-controls.js:300-306` preventDefaults them. Unknown fragments fall back only at page open; a `hashchange` listener must ignore names it does not know.
- **`history.replaceState` on file://** should be wrapped in try/catch with `location.replace("#…")` as fallback. Neither reloads. `replaceState` adds no history entry, so Back leaves the page; assigning `location.hash` adds one entry per click.

## 1. View and lens state

- `viewState` at `web/board.js:25-33`: `view` ∈ `board|calendar|durations|timeline|activity|testing`, `lens` ∈ `flat|user-request`. Lazy-render flags `renderedOnce` at `board.js:48-55`.
- Boot: `wireControls()` at `board.js:70`, then `renderColumns()`, then `applyView()` at `:79`. `wireControls` runs before the first `applyView`, so it is the place to read the hash (board.js itself stays untouched).
- `applyView` at `board-controls.js:19-98` shows the `#view-*` panel, hides inapplicable controls, resets `#board-main.scrollTop` for Timeline (`:43-45`), lazy first renders at `:73-92`, Board calls `applyLens()` (`:93-97`).
- View click handler `:196-202`: set `viewState.view`, `setActiveButton("[aria-label='Board views and lenses']", "data-view-target", …)`, `applyView()`.
- "URs only" is `user-request` + module variable `userRequestCardsFolded` (`:3-8`). Selection via `applyLensSelection(lens, fold)` (`:115-123`, resets `renderedOnce.userRequestLens`); `setActiveLensButton` `:103-111`; `applyLens` `:125-139`; lens click handler `:204-208`.
- `template.html`: nav `aria-label="Board views and lenses"` (`:31`); page buttons `data-view-target` (`:82-99`), Board hard-coded `is-active`/`aria-pressed="true"`; `#lens-group` (`:118-136`) with `data-lens-target="flat"`, `"user-request"`, `"user-request"` + `data-ur-cards="folded"`; panels `#view-board` (active) and `#view-*` (hidden) at `:243-540`. Init code must re-sync both button groups.

## 2. Storage, hash, history today

No `location`, `hash`, `history`, or `sessionStorage` use in `web/*.js`. `localStorage` only for Verify strip state (`board-cards.js:866,874`), detail panel width (`board-detail.js:762,771`), tester profile (`board-testing.js:362,377,431`). No stored view or lens preference.

## 3. serve and static output

No self-reload (`serve.go:161`); the only timer rewrites relative-time text (`board.js:68`). `generate.go:23` embeds `web/`; `boardJavaScriptFragmentPaths` (`:43-58`); `assembleBoardJavaScript` (`:1060`) splices fragments, each must end in exactly one LF; `assembleStaticPage` (`:1110`). Data comes from `board-data.js`, so file:// works.

## 4. JavaScript test harness

| File | Size | Tests | Time |
|---|---|---|---|
| `javascript_behavior_a_test.go` | 114 KB | 16 | 0.79s |
| `javascript_behavior_b_test.go` | 86 KB | 16 | 0.81s |
| `javascript_behavior_c_test.go` | 195 KB | 25 | 2.61s |
| `javascript_behavior_d_test.go` | 69 KB | 10 | 0.5s |

- `d` is the natural home. `generateLiveSite(t)` (`generate_test.go:368`); `sliceBalancedBlockAfter(t, html, "function X(")` (`generate_test.go:2009`) cuts functions by brace matching — no braces inside string literals in new functions.
- Plain Node via `runJavaScriptBehaviorProbe` (`generate_test.go:283`); tests hand-write `document`, `window`, `viewState`, `renderedOnce`, stub renderers. No test fakes `location` or `history` yet.
- Gate: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on` + node (`generate_test.go:271-281`); test name prefix `TestJavaScriptBehavior`.
- Closest example: `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (`c_test.go:2564-2705`) drives the real `applyView` with no location/history stubs — new code inside `applyView` must not need them.

## 5. Prime rules

`web/` is embedded; never edit `build/`, generated `index.html` or `board-data.js`; static page must run from file://; `TestBoardJavaScriptAssemblyStructure` locks the fragment manifest; the release is a suite version bump + changelog entry; no parser lock-step concern.

## 6. Tests that could break

`c_test.go:3670-3675` (view button order), `c_test.go:2564` (applyView without location stubs), `generate_test.go:798,801` (control attributes), the browser probes above.

## 7. Lesson bullets

- `_dev/primes/lessons-kanban-board.md:41`: a view that binds to persistent nodes owns its listener teardown.
- `:43`: describe the board by what its switcher covers.
- `:71`, `:76`: views share one scroll surface.
- `lessons-do-kanban.md`: nothing about views.

*Generated by Explore agent*
