---
id: REQ-589
title: 'Addendum: render the verify findings strip as the M4 slim band, one line closed and one row per finding open'
status: completed
created_at: 2026-09-05T18:19:11Z
user_request: UR-125
addendum_to: REQ-588
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-588, REQ-579, REQ-578]
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T18:20:44Z
  basis:
  - Route A
  - 4-file write set
  - 6 acceptance criteria
route: A
dispatch_at: 2026-09-05T18:21:42Z
builder_handback_at: 2026-09-05T18:44:37Z
integration_at: 2026-09-05T18:44:37Z
review_at: 2026-09-05T18:56:03Z
commit: 3442203290fa42de78d53cce39974c8add980207
heavy_verified_at: 2026-09-05T18:57:15Z
heavy_verified_revision: 3442203290fa42de78d53cce39974c8add980207
claimed_at: 2026-09-05T18:20:43Z
completed_at: 2026-09-05T18:57:15Z
release_at: 2026-09-05T18:57:15Z
---

# Addendum: Render the Verify Findings Strip as the M4 Slim Band, One Line Closed and One Row per Finding Open

## What

REQ-588 (the M1 rows, release 0.303.2) made each finding a chip, a detail line and a remedy line; the user's verdict was that it is not visually nice and huge. Replace the strip's rendering with mock-up M4 from `ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/`: a slim band that is one line when closed and one row per finding when open, with each remedy behind its row's chevron. The mock-up page `mockups/m4-closed.html` (and `m4-open.html`, `m4-open-remedy.html`) is the specification: its `<style>` block is the CSS to ship and its markup is the structure to render.

## Prior Implementation

REQ-588 shipped in commit 707ffb6c (merge ab251f24). `renderVerifyFindingsStrip` in `web/board-cards.js` groups findings by subject, prints a `.board-findings-subject` heading per group, then one `.board-findings-row` per finding (chip span, then a text span holding detail, optional "cleanup can fix" tag and remedy as blocks); skipped probes are muted rows with a "not checked" chip. `web/template.html` holds the strip section: a header with the title, `#board-findings-count` and a hint, then `#board-findings-rows` wrapping two `display: contents` hosts `#board-findings-cards` and `#board-findings-skipped-list`, which `applyView` in `web/board-controls.js` (REQ-578) reads to decide whether the strip has content on the Activity view. `web/board.css` holds the `.board-findings-*` rules around lines 607–700. The Node lane pins the list shape (`TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList`), the M1 row (`TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail`) and the Activity hide rule (`TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip`).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Template owns the static shell, the renderer fills the subject list and the rows, mock-up `vf-*` classes ship as `board-findings-*`, Node lane RED first, then render proof. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-182000/REQ-589-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`1528ea72`) touching exactly the four write-set files.
- [x] **[UNIFY]:** `git diff --stat` reviewed (4 files, 890+/377-); `gofmt -l .` empty; `go vet ./...` clean; debug-artifact scan over added lines empty; retired class names grepped out; per-file checks listed in the hand-back.

## Why

The user's words: "neither of these is visually nice and they are huge". A remedy is read after deciding to act, so it belongs behind a click; what stays visible must be short enough to scan.

## Requirements

- **D1, the band.** The strip is a slim band: no card border box, a 3 px amber left edge (`--accent-pending`), 8 px radius, `--surface-2` background, 5 px vertical padding. The header line is a small warning glyph plus a quiet uppercase VERIFY label, then the counts ("3 findings · 1 probe not checked"). The hint sentence is no longer rendered.
- **D2, closed state.** The whole strip is a `<details>` whose summary is the header line: label, counts, then every finding's subject in the mono face preceded by its weight dot (amber for an ordinary finding, green for `fixable`, grey for a skipped probe), the skipped probes summarised as "N probe(s) not checked", and a Show button at the right. Closed height is one line (about 34 px).
- **D3, open state.** Opening swaps the subject list for one row per finding, in the producer's order grouped by subject as today: dot, subject (mono, semibold), category, the detail clipped with an ellipsis at the line's end, the "cleanup can fix" tag as a small green pill when `fixable`, and a chevron. Skipped probes are rows too, with a grey dot and the category "not checked". Each row is a `<details>`; opening it shows the remedy under the row in an inset block labelled "What to do:" with a 2 px amber left rule, and lets the detail wrap. Rows have a hover background and a focus-visible ring. The Show button reads Hide with its chevron flipped while open.
- **D4, remembered.** The strip's open/closed state persists per browser in `localStorage` under one key, best-effort like the detail-panel width in `web/board-detail.js`; the default is closed. Row (remedy) state is not persisted.
- **D5, category as words.** The category token renders lowercase with hyphens shown as spaces (a mechanical transform, no list in the client), in faint ink. No uppercase chip remains.
- **D6, nothing else moves.** Producer, payload, grouping by subject (exact match on the payload field), the two host ids and REQ-578's Activity hide rule, and hide-when-empty stay as they are. No `.board-request*` classes.
- Board changes follow `_dev/primes/prime-kanban-board.md`; embedded web assets reach consumers on the next build.

## Red-Green Proof
**RED prompt/case:** Render the board with three findings under three subjects and one skipped probe (the gallery's data) and look at the strip; in the Node lane, inspect `#board-findings`.
**Why RED now:** The strip has no `details` element, no Show control, no stored state; every remedy is always visible; the header carries the hint sentence and the category is an uppercase chip.
**GREEN when:** The Node lane asserts: the strip's rows sit inside a `details` element that is closed by default and whose summary names every subject; each finding row is a `details` whose summary holds dot, subject, category and detail, and whose content holds the remedy; the category text is the token lowercased with spaces; the skipped probe is a row with the "not checked" category; opening the strip and reloading with the stored key set renders it open; REQ-578's Activity hide test and the hide-when-empty test still pass. The board screenshot matches `mockups/m4-closed.html`, `m4-open.html` and `m4-open-remedy.html` in both themes.
**Validation:** User picked M4 from the gallery ("ok, M4 is good") after two rounds of mock-ups; the gallery page is the approved design.

## Assets

- `ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/` (committed at 2a795ba0): the approved M4 pages and their captures. Screenshot 3 (not saved; the attachment cache expires) showed that gallery's M4 section: the live frame closed, then State 1 closed and State 2 open in light and dark.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6477 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4959 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "neither of these is visually nice and they are huge, please provide better mocks, it's fine to be colapsible as well" / "I want all of the options in the mockup, don't make me imagine what would be, also make it beautiful and professional" / "ok, M4 is good"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The approved mock-up pages are the specification (markup and CSS), the request names the four files, and the renderer change is a restructuring of one function plus its template section. Substantive in size, but nothing to explore or decide. `effort_estimate: effort-substantive`.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** The strip is now the M4 slim band. The template holds the static shell: an outer details element whose summary carries the warning glyph, the VERIFY label, the counts, the subject-list host and the Show/Hide labels, with the two unchanged row hosts inside it. The renderer fills the subject list (one entry per finding with its weight dot, skipped probes summed into one grey entry) and one details row per finding and per skipped probe (dot, subject, category as lowercase words, detail clipped closed, the cleanup-can-fix pill, a chevron; the remedy as the row's content under a "What to do:" label). The open/closed state is stored under one localStorage key and restored at load, default closed. The REQ-579/REQ-588 row rules and the hint sentence are deleted. Merge range 0f38e447..4a909573; builder branch head 1528ea72. Builder-authored Decisions (D-01 to D-09), Discovered Tasks and render evidence live in the hand-back at do-work/runs/work-2026-09-05-182000/REQ-589-handback.md.

## Decisions

- **D-10 (orchestrator) — D-03 stands and is put to the user in the hand-back rather than minted as a follow-up.** The builder's D-03 (a skipped-probe row shows the producer's whole sentence, no bold subject, because the skipped payload carries no subject field and the client must not parse prose) is the one visible difference from the approved mock-up. It is reversible in one function or, better, by a subject field on the skipped-probe payload. The user is present in this session, so the question is asked directly in the run's hand-back instead of through a pending-answers REQ. Value: no client-side parsing of producer text. Risk: the skipped row reads quieter than the mock-up until answered.

## Qualification

**Passed.** Read from the cumulative merge range 0f38e447..34422032 (the builder's two commits, merged at 4a909573 and 34422032; the range also sweeps in a sibling session's two-line edit to lessons-do-kanban.md that landed on main between the two merges, which is not this REQ's and is judged as such) and the builder's render evidence (gallery and live data, three states, both themes, plus 320/768/1280 px) inspected by the orchestrator against the approved M4 captures.

- **D1, the band.** The card border box is gone; the strip is the 3 px amber-edged band with the warning glyph, the VERIFY label and the counts, and the hint sentence is deleted from the template.
- **D2, closed.** One line: label, counts, every finding's subject with its weight dot (amber, green for fixable, grey for skipped probes summed into one entry), Show at the right. Measured 36 px at 1600 px against "about 34 px": the 2 px are the band's own borders. Accepted.
- **D3, open.** One 28 px details row per finding and per skipped probe in the grouped order: dot, mono subject, category words, detail clipped closed and wrapping open, the green cleanup-can-fix pill, a chevron; the remedy under the row in the inset block labelled "What to do:". Hover background and focus ring present in the rules a test reads.
- **D4, remembered.** One key, written on the disclosure's toggle event, read at load, guarded like the detail-panel width; default closed; a storage-denied browser falls back to closed (pinned).
- **D5, category as words.** `replace(/-/g, " ")`, no list.
- **D6, nothing else moves.** Producer, payload and `groupFindingsBySubject` untouched; both host ids survive inside the new shell and REQ-578's Activity test passes with every assertion unchanged; hide-when-empty pinned in the slim-band case.
- **One visible deviation, D-03:** the skipped-probe row has no bold subject because the skipped payload is one string per probe and the client must not split it; the row shows the whole sentence. Agreed with the reasoning (the newest lesson family says exactly this) and escalated to the user in the hand-back (D-10). The honest fix is a subject field on the skipped-probe payload, noted in Discovered Tasks.
- **One layout wart found in the live render, not the gallery:** with six subjects on the closed line the counts wrap ("6 / findings") because the label and count are not `flex: none`/`nowrap`. Sent back to the builder as a one-rule fix on the same branch before review (commit aafb7a70: nowrap and flex none on the label, count and toggle, plus a pinned assertion; re-rendered and inspected: "6 findings" on one line, only the subject list wraps), merged at 34422032.

Scope: four files, all in the write set. The two harness lines in REQ-578's test (an SVG factory and four sliced helpers) are the brief's anticipated case. No drift.

*Checked by work action*

## Testing

**Tests run (post-merge, detached worktree at 34422032; the main tree carries a sibling session's uncommitted do-work/ edits):**
- Focused probe do-work/runs/work-2026-09-05-182000/REQ-589-probe.sh, recorded through advance as the test gate: the three slim-band Node cases, REQ-578's hide-on-Activity case, the assembly structure test and the generate verify-payload tests — exit 0, wall 35s.
- Repository gate `bash _dev/tests/maintainer-verify.sh` — exit 0 on the first run, 3m10s; queue-kanban module 396 tests in 24s, do-work-cli module 771 tests in 81s, slowest file 28.30s against the 30s per-file budget.
- Builder, on the branch: Node lane 63 tests exit 0 (8s, then 13s after the fix); full queue-kanban module 400 tests exit 0 in 49s; `gofmt -l .` and `go vet ./...` clean.

**Red-green validation** (traced to `## Red-Green Proof`):
- javascript_behavior_c_test.go `TestJavaScriptBehaviorVerifyFindingsRenderAsTheSlimBand`: ✗ before implementation (no slim-band helpers in the shipped page) → ✓ after. Pins the closed line naming every finding with its dot colour and the skipped probes summed into one grey entry, four details rows in grouped order with a chevron, the summary part order for plain, fixable and skipped rows, the category as lowercase words, the remedy as row content with the producer's text unprefixed, the skipped row without a remedy, the empty re-render clearing all three hosts, and the shell markup with the display-contents hosts.
- `TestJavaScriptBehaviorVerifyFindingsBandRulesHideWhatIsNotBeingRead`: ✗ before (no accent edge, no subject-list/rows swap, no Show/Hide swap) → ✓ after; after the fix it also pins nowrap and flex none on the count, and fails when those two declarations are reverted (checked by the builder).
- `TestJavaScriptBehaviorVerifyFindingsStripRemembersItsOpenState`: ✗ before (no storage key in the shipped client) → ✓ after; default closed, toggle wired, "open"/"closed" written from the element's own state, both restored on load, a storage-denied browser falls back to closed.
- Render fact: builder screenshots at 1600 px in light and dark for closed, open and one-remedy-open on the gallery data and on the live queue, compared against the approved M4 captures; measured 36 px closed and 28 px per row; the six-finding live line re-captured after the fix with "6 findings" on one line. Inspected by the orchestrator.

**New tests added:**
- `TestJavaScriptBehaviorVerifyFindingsRenderAsTheSlimBand`, `TestJavaScriptBehaviorVerifyFindingsBandRulesHideWhatIsNotBeingRead`, `TestJavaScriptBehaviorVerifyFindingsStripRemembersItsOpenState` (javascript_behavior_c_test.go), replacing REQ-579's one-row-list case and REQ-588's remedy-line case, whose structure this REQ retires.

**Existing tests updated (cross-REQ impact):**
- `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (REQ-578): its stub document gained a createElementNS factory and its sliced-function list the four new helpers; no assertion changed.

**Heavy verification plan:**
- Ranges: 0f38e447..4a909573 (first merge) and 32813e4c..34422032 (fix merge), planned separately so the sibling session's commits between them do not inflate the lane set; both plans select the same three lanes and leave no path uncovered.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — every changed path under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same reasons
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same reasons

*Verified by work action*

## Review

**Overall: 94%** | 2026-09-05T18:56:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 90% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `verify.go:1390-1392` still says the board prints Subject as a heading above the rows that share it; this diff deleted subject headings, so the comment's stated reason is false although the behaviour it defends is unchanged — impact-negligible → report only

**Minor findings:**
- `verify.go:99-103` describes `Subject` as existing so several findings about one thing read "as one block"; the heading-plus-rows block it means is gone, though grouping still orders those rows adjacently — impact-negligible → report only
- The closed line prints an identical entry twice when one subject carries two findings, with nothing to tell them apart (seen on the live queue as `worktree-agent-REQ-590-cap-the-path-list`) — impact-user-visible → report only
- On the closed line, ordinary versus `fixable` is carried by dot colour alone; the rows say it in words with the green pill, the closed band does not — impact-user-visible → report only
- The load-time restore is an immediately invoked function but the test slices only its declaration, so deleting the `)()` would leave the storage test green while the strip forgets every reader's choice — impact-negligible → report only
- Nit: `id="board-findings-rows"` in `web/template.html:235` is read by no script and pinned by no test — impact-negligible → report only
- Nit: the class `board-findings-remedy-text` (`web/board-cards.js:811`) has no CSS rule and no test reads it — impact-negligible → report only

**Acceptance:** Pass — focused Node lane green at the merge revision (`ok`, 5.678s), and the renders match the approved M4 mock-up apart from the disclosed D-03 skipped row.
**Suggested testing:** 6 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Building from a mock-up the user had approved in every state: the M4 pages were the specification, so the builder ported a `<style>` block and a markup shape instead of designing, and the review checked the shipped rules against the mock-up declaration by declaration.
- Rendering on the live queue as well as the gallery data: the six-finding closed line exposed the count wrap that the four-item gallery could never show, and it cost one nine-line commit on the same branch before review.

**What didn't:**
- REQ-588 (this morning) shipped from a one-state mock-up and was judged "not visually nice and huge" the same afternoon; two more mock-up rounds were needed before a layout the user would accept. Layout work should not be captured until the user has seen every state.
- The cumulative merge range over a shared checkout sweeps in sibling commits (here a release and two lessons edits); the heavy planner had to be run per merge commit to keep the lane set honest.

**Worth knowing:**
- The strip's shell is static in the template; the renderer fills only `#board-findings-subjects` and the two row hosts. REQ-578's Activity rule still reads the two host ids.
- The closed line names every finding, so a subject with two findings appears twice with nothing to tell them apart, and fixable versus ordinary is colour alone on that line (review, report only). A skipped probe has no subject field, so its row shows the whole sentence (D-03).
- `verify.go` still carries two comments describing subject headings the board no longer prints (review, report only).

## Orientation

Now the board's Verify Findings strip is one line until you open it — subjects with weight dots closed, one row per finding open, the remedy behind each row's chevron, the choice remembered per browser; lives in the queue-kanban board tool (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). Prime spot-check: every path both primes name still exists; neither is made stale. No map change.

## Heavy Verification Plan

- Base revisions: 0f38e447 (first merge's parent) and 32813e4c (fix merge's parent)
- Target revision: 34422032 (the fix merge, recorded in `commit:`; first merge 4a909573)
- Planned with `plan-heavy-verification --manifest _dev/tests/heavy-lanes.json` for 0f38e447..4a909573 and for 32813e4c..34422032; both select the same lanes, no uncovered paths.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — every changed path matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same reasons
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same reasons

## Heavy Verification Result

- Target revision: 3442203290fa42de78d53cce39974c8add980207
- Execution revision: 3442203290fa42de78d53cce39974c8add980207 (detached worktree `.git/work-run-20260905-1820/gate-34422032`, equal to the main tree's HEAD at the time)
- Run: `run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane queue-kanban-javascript --lane queue-kanban-browser --lane staged-skills` with `QUEUE_KANBAN_BROWSER` naming Chrome, 18:53:29Z to 18:56:49Z
- queue-kanban-javascript: exit 0, executed (fingerprint mismatch), 10s
- queue-kanban-browser: exit 0, executed (fingerprint uncertain, browser runtime), 135s
- staged-skills: exit 0, executed (fingerprint mismatch), 51s

## Timing

Observed 2026-09-05T18:21:42Z to 2026-09-05T18:57:15Z: 35m 33s total, 30m 37s attributed across 3 events, 4m 56s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 22m 55s | 1 |
| verification-gate | 7m 42s | 2 |

Slowest stage: builder-work / implementation, 22m 55s, outcome success.
Slowest command: verification-gate / maintainer-verify, 3m 56s, exit 0, .
