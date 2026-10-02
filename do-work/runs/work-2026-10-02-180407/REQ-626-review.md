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
