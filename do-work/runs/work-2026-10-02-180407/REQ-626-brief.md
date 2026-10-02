# Builder brief — REQ-626: Give every board page and lens its own URL

**Worktree (your working directory, the only tree you write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-626-board-page-urls`
**Branch:** `worktree-agent-REQ-626-board-page-urls`, based on `7917efbe`. Commit on this branch only, message prefix `[REQ-626]`.
**Route:** B. **TDD: yes** (RED before GREEN, evidence required). **Impact:** user-visible. **Estimate P50:** 20 active minutes.
**Hand-back file (the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-180407/REQ-626-handback.md`. Never stage or commit it.

## Read first (absolute paths; the REQ and exploration live in the MAIN tree, read-only)

1. The REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-626-board-page-urls.md` — What, Detailed Requirements, Red-Green Proof, Exploration, Scope. The `## Scope` "Files I will touch" list is your write boundary, with the one widening below.
2. Exploration with file:line anchors: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-180407/REQ-626-exploration.md`.
3. In your worktree: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`. Lessons satellites (read the view-related bullets at least): `_dev/primes/lessons-kanban-board.md` (lines 41, 43, 71, 76), `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`.
4. Crew rules in your worktree under `skills/do-work/crew-members/`: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `frontend.md` (if present), `testing.md`. Naming rule: two words minimum for anything with reach.

## Fixed design (decided by triage and exploration; do not re-open)

- **Fragments** (lowercase, stable): `#board` (Board, flat Columns lens), `#board/by-ur` (By UR lens), `#board/urs-only` (By UR lens + folded cards), `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`. Filters never enter the URL.
- **All code in `web/board-controls.js`.** Do not touch `web/board.js` or `web/template.html`.
- **Parse** with one function (e.g. `boardStateFromFragment(hashText)`) returning `{view, lens, fold}` or null for anything unknown (empty, `#board-main`, heading ids, `#`). **Format** with one function (e.g. `boardFragmentFromState()`) from `viewState.view`, `viewState.lens` and `userRequestCardsFolded`; a non-Board page omits the lens.
- **Read at boot inside `wireControls()`**, which runs before the first `applyView()`: if the hash parses, set `viewState.view` and the lens/fold through the existing selection path (`applyLensSelection` or equivalent, so `renderedOnce.userRequestLens` resets correctly) and re-sync both button groups with the existing active-button helpers (the template hard-codes Board/Columns as active). Do not call `applyView()` from the boot read — board.js does that next, which keeps the lazy first render path unchanged. Unknown/empty: change nothing, write nothing, throw nothing.
- **Guard every `location`/`history` access** (`typeof window.location`/`window.history` checks) so the existing Node behaviour tests that run `applyView`/`wireControls` without those fakes keep passing. Nothing new runs inside `applyView` itself.
- **Write on click only, never on load.** After a page click or a lens click updates state, write the fragment with `history.replaceState(null, "", "#" + fragment)` inside try/catch, falling back to `location.replace("#" + fragment)`. Skip the write when the hash already equals it. `replaceState` adds no history entry, so Back leaves the page as before (record this as the Back decision).
- **`hashchange` listener** (for a hand-edited address bar): if the new hash parses, apply it the same way and call `applyView()`; if it does not parse (skip link, heading anchors), ignore it.
- **Browser probe harness.** Every browser probe check that requires the page URL to end in `probe.html` must accept a trailing fragment. Add ONE shared helper in `browser_probe_test.go` (e.g. `probePagePathWithoutFragment(href string) string` or a predicate) and use it at the transport checks (`browser_probe_test.go` ~366 and ~481, `durations_browser_probe_test.go` ~567) and at every other `HasSuffix(..., "probe.html")`-style check that fails once a view click writes a fragment. Grep: `grep -ln "probe.html\|\.Href" skills/do-work-board/tools/queue-kanban/*_test.go`. **Scope widening granted:** you may make this one-line substitution in any `*_test.go` in the queue-kanban package that checks the probe page URL; list each such file in the manifest. Change nothing else in those files.

## Tests (TDD; each named for the failure it pins)

- In `javascript_behavior_d_test.go` (smallest file; keep it under the 30 s per-file budget), `TestJavaScriptBehavior…` tests driving the REAL functions sliced from the generated page (`sliceBalancedBlockAfter`; no braces inside string literals in new functions), with hand-written `window`/`location`/`history` fakes:
  - each of the eight fragments read at boot yields the right `view`/`lens`/fold and the right active buttons;
  - unknown (`#board-main`, `#lessons-learned`, `#nope`) and empty fragments leave Board/flat, write nothing, throw nothing;
  - a page click and a lens click write the expected fragment through `history.replaceState`, and a throwing `replaceState` falls back to `location.replace`;
  - a `hashchange` to a known fragment switches the view; to an unknown one is ignored.
- New `page_fragment_browser_probe_test.go`: one real-browser probe, following the existing probe shape (`browser_probe_test.go` helpers), that opens `probe.html#timeline` and asserts the Timeline panel is the visible one, then clicks the Activity button and asserts `location.hash === "#activity"`. This is the REQ's RED/GREEN case in a real engine.
- RED: write the tests first against a compiling stub (functions exist, do nothing), run them, keep the failing output verbatim. Do not commit the stub.

## Never touch

Anything under `do-work/` in either tree except your hand-back file; `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`; `web/board.js`, `web/template.html`, any generated build output; `go.mod`/`go.sum`. A file outside Scope (and outside the granted test widening) that you believe you need: stop and report it in the hand-back instead of editing it.

## Docs

`prime-do-kanban.md`: one Traps bullet: page/lens fragments are written by `board-controls.js` on click only; browser probes must compare the page URL without its fragment (name the helper).

## Testing (record each command's own exit line)

From `skills/do-work-board/tools/queue-kanban` in your worktree:

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run TestJavaScriptBehavior .
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...
QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" QUEUE_KANBAN_JAVASCRIPT_PROBES=off go test -count=1 -run 'Browser|Probe|Citation|Clipboard|Contrast|Priority|Progress|Scroll' . 
gofmt -l . ; go vet .
```

The browser run proves both the new probe and that no existing probe broke on the fragment. If the browser binary is missing, say so; do not report a skip as a pass. Also record per-file times if the run prints them (each test file < 30 s).

## Hand-back format (write the file, then reply with ONE line)

```
# Hand-back — REQ-626
**Branch:** worktree-agent-REQ-626-board-page-urls  **Base:** 7917efbe  **Head:** <hash>
## File manifest
- path (new|modified|deleted) — one line each
## Integration seams
None, or the exact line and the file it belongs in
## P-A-U
- [PLAN] <two or three sentences of the approach you took>
- [APPLY] <confirmation scope stayed inside the declared files + granted widening>
- [UNIFY] <git diff --stat output, linters run, files checked>
## Red-green evidence
- <test name>: RED output (verbatim, trimmed) → GREEN at <commit>
## Tests run
- <command> → exit <n>, <n> tests, wall time
## Decisions
- D-01: ... (reasoning; mark DECIDE & STATE or ESCALATE with Value/Risk)
## Discovered Tasks
- ... or None
## Lessons read
- each lessons file and whether whole-satellite or the named lines
## Blockers
- None, or what stopped you
```
