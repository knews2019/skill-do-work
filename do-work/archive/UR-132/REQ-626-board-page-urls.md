---
id: REQ-626
title: 'Give every board page and lens its own URL'
status: completed
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 3-file write set
  - 7 acceptance criteria
  calculated_at: 2026-10-02T18:04:23Z
created_at: 2026-10-02T18:02:37Z
user_request: UR-132
domain: frontend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-627"]
batch: board-links-and-disk-space
write_set: ["skills/do-work-board/tools/queue-kanban/web/board-controls.js", "skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go", "skills/do-work-board/tools/queue-kanban/browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go"]
review_at: 2026-10-02T18:41:24Z
integration_at: 2026-10-02T18:18:36Z
builder_handback_at: 2026-10-02T18:18:20Z
dispatch_at: 2026-10-02T18:09:36Z
heavy_verified_at: 2026-10-02T18:50:16Z
heavy_verified_revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086
claimed_at: 2026-10-02T18:03:38Z
status_changed_at: 2026-10-02T18:36:12Z
commit: 505ec74ac0d71f2554b52a8ef69259c12a7a6bcb
completed_at: 2026-10-02T18:50:25Z
release_at: 2026-10-02T18:50:25Z
---

# Give Every Board Page and Lens Its Own URL

## What
Each board page (Board, Activity, Calendar, Timeline, Durations, Testing) and each Board lens (flat columns, By UR, URs only) gets its own URL that opens that page directly, and switching pages or lenses updates the address bar so the link can be copied and shared.

## Why
The board has many pages and the user wants to point people at a specific one by URL. Today `viewState.view` lives only in JavaScript (web/board-controls.js), so every link opens the Board page and the address bar never changes.

## Detailed Requirements
- A URL fragment names the page, and optionally the Board lens, for example `#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`. The exact spelling is the builder's, but it must be stable, lowercase, and readable.
- Opening the board with a fragment shows that page (and lens) directly, including the lazy first render the view switch does today.
- Clicking a page or lens control updates the fragment without reloading (`history.replaceState` or `location.hash`), and the browser Back button is not broken by it.
- An unknown or empty fragment falls back to the current default (Board, flat lens) without an error.
- The static snapshot (`do-work-board static`, a file opened from disk) behaves the same, since fragments work on `file://` URLs.
- Filters (search text, domain, status, recently-done window) stay out of the URL. User decision at capture.
- Any existing stored-view preference (localStorage, if present) yields to an explicit fragment.

## Constraints
No new dependency; plain JavaScript in the existing files. Keep the change in board-controls.js and template.html, plus tests. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-627 (Show free disk space on the Testing page).

## Builder Guidance
High certainty on the behavior; builder latitude on fragment spelling and on whether lens is a path segment or a query-like suffix. A JavaScript behaviour test in the existing `javascript_behavior_*_test.go` harness is the expected proof; keep the test file under the 30 s per-file budget.

## Red-Green Proof
**RED prompt/case:** Open `http://127.0.0.1:8090/#timeline` (or any page URL). The Board page shows, not Timeline; click Timeline and the address bar still reads the bare URL.
**Why RED now:** The selected view is JavaScript state only; nothing reads or writes the URL.
**GREEN when:** `#timeline` opens the Timeline page directly; each of the six pages and each Board lens has a URL that opens it; clicking a page or lens control updates the address bar to that URL; an unknown fragment falls back to the Board page.
**Validation:** User confirmed

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban UI and browser behaviour. Read it anyway if time allows.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: views and view scroll.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: five small functions in the controls file (fragment parse and format, guarded hash read, selection without render, guarded replaceState write); read once at wiring, write after clicks, hashchange applies known names only.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: declared files plus the granted one-line URL-check substitution in four more browser probe test files; board.js and template.html untouched.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: 10 files, 473 insertions, 19 deletions; gofmt empty, go vet exit 0; every substituted URL check reviewed line by line.

## Full Context
See `do-work/user-requests/UR-132/input.md` for complete verbatim input.

*Source: there are a lot of pages in that kanban (BTW: capture-request, I need a link for each page so I can point to them by URL)*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (one fragment per page and lens, read on load, written on switch, unknown falls back), but where `viewState.view` and the lens are set, how the lazy first render runs, whether a stored-view preference exists, and how the `javascript_behavior_*_test.go` harness drives the page all need discovery before dispatch. Two web files plus one test, no architectural change, so no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors: `do-work/runs/work-2026-10-02-180407/REQ-626-exploration.md`. The findings that change what the builder writes:

- **The browser probe harness rejects any fragment.** The load wait and every in-page evaluate in `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (lines 366 and 481) and `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (line 567) fail unless the page URL ends in `/probe.html`. A view click that writes a fragment would fail every later browser measurement. Those two test helpers must accept a trailing fragment; this widens the Scope beyond the REQ's write_set, which the REQ's "plus tests" constraint already allows.
- **Not every hash is a view.** The skip link targets `#board-main` and rendered detail bodies carry heading ids, so a hash listener must ignore names it does not know. "Unknown falls back to Board" applies at page open only.
- **Boot order fits without touching the boot file.** The control wiring runs before the first view render, so the fragment can be read inside the wiring function in the controls file and the first render shows the right page with its lazy first render.
- **"URs only" is not a lens value.** It is the By-UR lens plus the folded-cards flag, set together through the existing lens-selection helper.
- **The template needs no change.** Its hard-coded active Board button is re-synced by the existing active-button helpers when the fragment selects another page.
- **No stored view or lens preference exists** (localStorage holds only the Verify strip, the panel width and the tester profile), so the "yields to an explicit fragment" requirement holds trivially.
- **Tests:** the JavaScript behaviour harness is plain Node with hand-written fakes; `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` is the smallest file (0.5 s). Functions are sliced by brace matching, so new functions must not put braces in string literals. An existing applyView test runs with no location or history fake, so code inside the view render must not need them.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modify) — read the fragment during control wiring, write it on page and lens clicks, react to hash changes that name a known page
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modify) — behaviour tests for reading, writing and falling back
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modify) — the page-URL checks accept a trailing fragment
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (modify) — same check
- `skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go` (new) — real-browser proof that a fragment opens its page and a click writes the fragment
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — one Traps bullet naming the fragment contract and the probe URL check
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — page-URL checks compare the address without its fragment (added from the hand-back, D-06)
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go` (modify) — page-URL checks compare the address without its fragment (added from the hand-back, D-06)
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modify) — page-URL checks compare the address without its fragment (added from the hand-back, D-06)
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go` (modify) — page-URL checks compare the address without its fragment (added from the hand-back, D-06)

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/web/template.html` (the button state is re-synced in JavaScript), `skills/do-work-board/tools/queue-kanban/web/board.js`, generated build output, release paths (VERSION, both changelogs, the version action), anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] Each of the six pages and each Board lens has a stable, lowercase, readable fragment
- [ ] Opening with a fragment shows that page and lens directly, including its lazy first render
- [ ] Clicking a page or lens control updates the fragment without a reload, and Back is not broken
- [ ] An unknown or empty fragment falls back to Board, flat lens, without an error
- [ ] The static snapshot opened from file:// behaves the same
- [ ] Filters stay out of the URL
- [ ] An explicit fragment wins over any stored view preference (none exists today)

## Pre-Flight

**Git:** ✓ Integration tip 7917efbe on `main`; the only dirt is this REQ's own `do-work/` trail (working REQ, run directory) — no third-party paths
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` (DO_WORK_FAST_STAGE_REUSE=off) exit 0 at 18:07Z, both Go stages EXECUTING, gate wall 132s
**Dependencies:** ✓ Go toolchain and node present; no new dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** The board now reads a page or lens fragment once while wiring its controls, so a link such as `#timeline` or `#board/urs-only` opens that page with its normal lazy first render; page and lens clicks write the fragment with replaceState (falling back to location.replace), never on load; a hand-edited fragment that names a page switches to it and any other fragment is ignored. Browser probe URL checks now ignore the fragment through one shared helper. Merge range 498dd997..505ec74a (builder commit 3100dea9, merge 505ec74a).

## Qualification

**Diff range:** 498dd997..505ec74a (builder commit 3100dea9, merge 505ec74a)
**Gate records:** qualify satisfied; scope-drift satisfied after the Scope was synced with the four browser probe test files the brief's granted widening allowed (D-06).
**Warnings judged:** QUALIFY-NEW-FILE-UNWIRED on `skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go` — expected: a Go test file is found by `go test` through its name, so no static reference can exist. Not dead code.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. Eight fragments in one lookup table, null for anything else (own-property lookup, case-sensitive); the boot read runs inside the wiring function before the first view render and only selects state, so the lazy first render path is unchanged; clicks write through replaceState in try/catch with location.replace as the fallback and skip an unchanged hash; nothing is written on load; the hashchange listener ignores unknown names (skip link, heading ids); every location and history access is guarded so the Node probes without a window still run. Filters are not in the fragment. No stored view preference exists, so the explicit fragment wins trivially.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY was cross-checked against `git diff --stat 498dd997..505ec74a` (ten files, all in the synced Scope, nothing under do-work/).
**Live data flow:** the controls file is one of the fragments spliced into the page by the generator for both static output and serve, so the same code runs on file:// and over HTTP.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (DO_WORK_FAST_STAGE_REUSE=off) on the merged tree at 505ec74a
**Result:** ✓ All passing — exit 0, gate wall 144s; stage queue-kanban-fast-tests EXECUTING (409 tests, wall 40s, slowest file 20.93s < 30s); stage do-work-cli-fast-tests EXECUTING (866 tests, wall 73s, slowest file 24.30s < 30s). Green-gate record satisfied by advance.

**Focused tests:** `go test -count=1 -run 'TestJavaScriptBehavior(Board|UnknownBoard)' .` in the board package with JavaScript probes on and browser probes off → exit 0 (advance probe record satisfied).

**Builder runs (from the hand-back, branch 3100dea9):** the JavaScript behaviour lane (77 tests, 7.0s), the whole package with probes off (57s), and the real-Chrome browser lane (70 top-level passes, 0 failures, 184s; every existing probe that clicks a page passed with the fragment in the URL; new probe file 2.4s).

**Red-green validation:** traced to `## Red-Green Proof`; RED was taken against a compiling stub, GREEN at builder commit 3100dea9:
- TestJavaScriptBehaviorBoardFragmentOpensItsPageAtBoot: ✗ `opening #timeline: got view="board" lens="flat"` → ✓
- TestJavaScriptBehaviorBoardClickWritesItsFragment: ✗ `timelineClick wrote [], want ["history.replaceState #timeline"]` → ✓
- TestJavaScriptBehaviorBoardHashChangeSwitchesKnownPagesOnly: ✗ `hashchange to #durations did not switch and render the page` → ✓
- TestBrowserBehaviorPageFragmentOpensAndFollowsThePage (the REQ's RED case in Chrome): ✗ `probe.html#timeline did not open the Timeline page` → ✓
- TestJavaScriptBehaviorUnknownBoardFragmentKeepsTheBoard pins the fallback and passes against the stub by design.

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (four behaviour tests above)
- `skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go` (one browser probe)

**Existing tests updated (cross-REQ impact):**
- Browser probe URL checks in `skills/do-work-board/tools/queue-kanban/browser_probe_test.go`, `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`, `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`, `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go`, `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go`, `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go` now compare the page address without its fragment — intentional, a page click now writes one.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 498dd997..505ec74a
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

## Review

**Overall: 94%** | 2026-10-02T18:41:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: After a hand-edited fragment, Back to an entry with an empty hash leaves the old page on screen under a bare address. Reproduced in Chromium on the generated static board: open `index.html`, set `#durations`, then `history.back()` gives hash `""` with Durations still shown. The cause is that `boardStateFromFragment("")` returns null and the `hashchange` listener in `skills/do-work-board/tools/queue-kanban/web/board-controls.js:290-298` ignores it. The REQ scoped "empty falls back to Board" to page open only, so this is not a requirement miss, but a link copied at that moment opens Board. No test pins this case. — impact-user-visible → report only
- F2: The skip link (`#board-main`, `web/template.html:13`) replaces the page fragment in the address bar while the page stays as it was. After a keyboard user takes the skip link on Timeline, the address reads `#board-main` and a copied link opens Board. Back restores `#timeline`, verified in Chromium. — impact-negligible → report only
- F3: Restatement Sweep. `skills/do-work-board/docs/board-guide.md:25` still describes the page switch as "**Board / Calendar / Testing**". That lists three of the six pages this diff now names as URLs (it is missing Activity, Timeline and Durations). The guide also never says that each page and lens now has a link (`#timeline`, `#board/by-ur`, ...), which was the user's stated purpose. This text was stale before this diff, but the diff makes the page list a public contract. It is not builder scope drift. — impact-user-visible → report only
- F4 (Nit): Page fragments share the `#` namespace with goldmark auto heading ids (`render.go:28` `parser.WithAutoHeadingID()`). Every REQ has a `## Testing` heading, so an open drawer holds `id="testing"`. Hand-editing to `#testing` with a drawer open also scrolls the drawer, and an in-body markdown link to `#testing` or `#timeline` would switch the board page. Rare and harmless. — impact-negligible → report only

**Acceptance:** Pass. The four focused JavaScript behaviour tests and the real-Chrome fragment probe pass. A hands-on Chromium session on a freshly generated static board confirmed:
- `#board/by-ur` and `#board/urs-only` render the UR lens on the first `applyView` (lens shown, columns hidden, `is-folded` set for URs-only, correct buttons pressed).
- Clicks write the fragment without adding history.
- Hand edits switch the page.
- Back steps between hand-edited page fragments correctly.

**Suggested testing:** 3 items
- Open a page link on `file://` in Safari and Firefox, where `replaceState` may throw. That exercises the `location.replace` fallback, which so far is covered only by a fake.
- Share-link check by hand on the served board (`http://127.0.0.1:8090/#durations`, `#testing`). In serve mode, confirm the Testing page's lazy first render loads the tester profiles.
- Optional pin for F1: a hashchange to `""` after a hand edit (decide whether that should mean Board or keep the current page).

**Follow-ups created:** None (4 findings report only)

### Review notes (evidence for the record)

- **Requirements checklist:** all 7 delivered.
  - Stable lowercase fragments: 8 entries in one table, `board-controls.js:154-163`.
  - Fragment opens the page and lens with its lazy first render: the boot read happens in `wireControls` before `board.js:78-79` runs `renderColumns(); applyView()`.
  - Click writes the fragment with `replaceState` and Back is unchanged: the probe asserts `history.length` delta 0.
  - Unknown or empty fragment falls back at open: covered by the 7-case test.
  - `file://` support: the same spliced code runs there, and the browser probe runs on `file://`.
  - Filters stay out of the URL: nothing about filters is written.
  - Stored-preference precedence: holds trivially, since no stored view preference exists.
- **Focused checks requested by the orchestrator:**
  - (a) A lens click on another page cannot happen through the UI, because `#lens-group` is hidden off-Board (`board-controls.js:55`). Even when it is forced, `writeBoardFragment` writes the current non-Board page, which is correct because non-Board fragments carry no lens.
  - (b) The global URs-only fold (`userRequestCardsFolded`) is written only by `applyLensSelection` and the new `applyBoardStateSelection`. Group heads in `board-cards.js` read the fold but never write it, so no other control can leave the fragment stale.
  - (c) Boot with `#board/by-ur` resets `renderedOnce.userRequestLens` and lets the first `applyView` render the lens. Confirmed in a real browser.
  - (d) No bare `HasSuffix(..., "probe.html")` remains in the package tests. The two remaining `strings.Contains(LocationHref, browserProbePageFileName)` checks at `durations_browser_probe_test.go:966,1277` already tolerate a fragment. `citations_test.go:338` compares a server request path, which never carries a fragment.
  - (e) Back: `replaceState` adds no entries. Back between hand-edited page fragments re-applies the page. The empty-hash case is F1.
- **Restatement Sweep:**
  - The probe-URL rule is consistent across all 7 test files and the new prime bullet (`prime-do-kanban.md:41`).
  - The view list is restated stale only in `docs/board-guide.md:25` (F3).
  - `actions/board.md` and `SKILL.md` name only the Testing view and stay accurate.
- **Decisions:** D-01..D-07 are present in the REQ and the hand-back, and they cover every non-obvious choice visible in the diff (the non-rendering selection helper, the fragment-aware session starter, the case-sensitive own-property lookup, the scope widening).
- **Tests:** red-green is recorded for 4 of 5 new tests. The fallback test passes against a stub by design and says so. The cross-REQ test edits are one-line URL-check substitutions traced to REQ-626 in the helper's comment.
- **P-A-U:** all three boxes are ticked. The orchestrator ticked them from the hand-back and cross-checked against the diff stat.
- **Domain review (frontend):**
  - No new console errors. The only console entry was a favicon 404 from the scratch HTTP server.
  - No dependency was added and nothing is animated.
  - Keyboard behaviour is unchanged except F2.
  - Responsive checks do not apply, because there is no layout change.
- **Coding guardrails:** names are two-part and findable (`boardStateFromFragment`, `writeBoardFragment`, `probePageAddressWithoutFragment`). Comments explain why. No debug artifacts.

*Reviewed by review-work action*

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-02-180407/REQ-626-handback.md`). All are DECIDE & STATE.

- D-01: Back behaviour. replaceState adds no history entry, so Back leaves the board as before; the browser probe asserts the history length does not grow on a click.
- D-02: The boot read uses a new selection helper that renders nothing, because the existing lens-selection helper renders the UR lens before the board has rendered its columns; the new helper resets the UR-lens render flag the same way. Lens clicks still use the existing helper.
- D-03: A non-Board fragment carries no lens, so opening Activity keeps whatever lens was selected; the Board fragment selects the flat lens explicitly.
- D-04: A fragment-aware variant of the browser session starter was added (the old starter delegates with an empty fragment), because the new probe must really open the page with a fragment. Inside a Scope file.
- D-05: The fragment lookup is an own-property, case-sensitive table lookup, so names such as the object prototype's keys or an upper-case page name read as unknown.

Orchestrator decision:
- D-06: Scope widened from the hand-back. Exploration found the probe-page URL check in three places; the builder found it in four more browser probe test files, all needing the same one-line substitution once a click writes a fragment. The brief granted that substitution in advance; the Scope and write_set were synced to the actual files before qualification. DECIDE & STATE.
- D-07: After a machine restart the next session ran `recover --take-over` on this claim, which treated it as a crash: it returned the REQ to the queue and stripped every orchestrator-generated section, although the implementation was already merged at 505ec74a and gate-green. The REQ was re-claimed through queue-mode advance and this session's own sections were restored byte-for-byte from the handoff commit 1794c268; only the new `status_changed_at` was kept. No evidence was re-authored. DECIDE & STATE (the prior text is in git; the merge is unchanged).

## Discovered Tasks

From the builder's hand-back (none) and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- F1: Back to an entry with an empty hash after a hand edit leaves the old page shown under a bare address. — impact-user-visible → report only
- F2: The skip link replaces the page fragment in the address bar. — impact-negligible → report only
- F3: `skills/do-work-board/docs/board-guide.md` names only three of the six pages and does not mention page links. Stale before this diff. — impact-user-visible → report only
- F4: Page fragments share the hash namespace with drawer heading ids. — impact-negligible → report only
- The handoff action told the next session to continue claimed REQs, but `recover --take-over` on the same writer returned both claims to the queue and stripped every orchestrator-generated section of a REQ that was already merged and gate-green (D-07). The handoff prompt and the takeover contract disagree about what "resume" does. — impact-user-visible → report only

## Lessons Learned

**What worked:** Granting the one-line probe-URL substitution in the brief before dispatch; the builder found four more copies of the check than exploration did and fixed them without a scope stop. Reading the fragment inside the control wiring, before the first view render, kept every page's lazy first render unchanged.
**What didn't:** The probe-page URL check had been copied into seven test files instead of living in one helper, so a page-level URL change touched all of them. Explore agents are read-only and cannot write their findings file; the orchestrator has to save it. After a machine restart, `recover --take-over` reset both claims to the queue and stripped the run's sections, so they were restored from the handoff commit.
**Worth knowing:** Any page code that writes the address bar must keep the browser probes' page check fragment-tolerant (use the shared helper in `skills/do-work-board/tools/queue-kanban/browser_probe_test.go`). Unknown hash names are ignored on purpose because the skip link and drawer heading ids share the namespace.

## Orientation

Now every board page and Board lens has a shareable link (`#timeline`, `#board/by-ur`, `#board/urs-only`, ...) that opens it directly, and clicking a page or lens updates the address bar; lives in the board tool's web controls (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). Not a map change: the existing view switch gained a read at wiring time and a write after clicks. Prime spot-check: `prime-do-kanban.md` gained the fragment Traps bullet and its referenced paths exist; `prime-kanban-board.md` and `prime-releases.md` unchanged and current.

## Heavy Verification Plan

- Base revision: 498dd9974336b46aed5bb9954e063101c061fcb0
- Target revision: 505ec74ac0d71f2554b52a8ef69259c12a7a6bcb (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

## Heavy Verification Result

- Target revision: 505ec74ac0d71f2554b52a8ef69259c12a7a6bcb
- Execution revision: 172d0e5caaf54729ae1efa0fbbb09f574c2b2086 (detached drain checkout `.git/work-run-2026-10-02-180407/drain-head`, one run for REQ-626 and REQ-627, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 7s
- queue-kanban-browser: exit 0, executed (not skipped), 78s
- staged-skills: exit 0, executed, 40s

Green: every selected lane present, exit 0, none skipped, none reused.

## Timing

Observed 2026-10-02T18:04:07Z to 2026-10-02T18:50:16Z: 46m 09s total, 20m 05s attributed across 5 events, 26m 04s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 8m 53s | 1 |
| exploration-preflight | 5m 33s | 1 |
| verification-gate | 5m 00s | 2 |
| handback-merge | 39s | 1 |

Slowest stage: builder-work / builder implementation in worktree, 8m 53s, outcome success.
