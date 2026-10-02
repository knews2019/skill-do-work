---
id: REQ-628
title: 'Board guide names all six pages, their links, and the disk line'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-02T20:48:08Z
created_at: 2026-10-02T19:49:03Z
user_request: UR-133
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-629"]
batch: ur-132-follow-ups
write_set: ["skills/do-work-board/docs/board-guide.md", "skills/do-work-board/actions/board.md"]
dispatch_at: 2026-10-02T20:48:05Z
builder_handback_at: 2026-10-02T20:48:45Z
integration_at: 2026-10-02T20:48:50Z
review_at: 2026-10-02T20:52:52Z
heavy_verified_at: 2026-10-02T21:14:49Z
heavy_verified_revision: 069f922e64c1ed14ca31211642e395a784561f57
claimed_at: 2026-10-02T20:47:38Z
commit: 4c8e249acfa058bd75e6ed8a80175172ff31fbd5
completed_at: 2026-10-02T21:15:03Z
release_at: 2026-10-02T21:15:03Z
---
# Board Guide Names All Six Pages, Their Links, and the Disk Line

## What
Update the board's user-facing docs so they match what shipped in 0.305.61 and 0.305.62: the page switcher has six pages, every page and Board lens has its own link, and the Testing toolbar shows the repo root's free disk space.

## Why
The REQ-626 review (Restatement Sweep, finding F3) found `skills/do-work-board/docs/board-guide.md` line 25 still calls the switcher "**Board / Calendar / Testing**", three of the six pages, and the guide never mentions the page links that were the user's stated purpose. The REQ-627 review (finding M1) found the guide's Testing section and `skills/do-work-board/actions/board.md` Step 6 describe the Testing toolbar without the new disk line. The user asked to capture this.

## Detailed Requirements
- `board-guide.md`: the toolbar paragraph names the switcher by what it covers (all six pages: Board, Activity, Calendar, Timeline, Durations, Testing), not by a partial list.
- `board-guide.md`: one short passage says each page and Board lens has its own link and lists the fragments (`#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`), that clicking a page updates the address bar without adding Back steps, that filters are not part of the link, and that links work on a static snapshot opened from disk.
- `board-guide.md` Testing view section: one sentence on the disk line (free of total, amber below 10 GiB, red below 3 GiB, "(at generation)" on a static snapshot, "not measured" when the platform cannot measure; the live board shows the value from page load until reload).
- `board.md` Step 6: one clause naming the disk line in the Testing toolbar, only if that step describes the toolbar's contents.
- Do not restate behaviour the code does not have; check each sentence against `web/board-controls.js` and `web/board-testing.js`.

## Constraints
Docs only. Shipped files change, so this is a release. Follow the lesson "describe the board by what its switcher covers, not by naming some of its tabs" (`_dev/primes/lessons-kanban-board.md`).

## Dependencies
None. Independent of REQ-629 (handoff resume must not reset claims).

## Builder Guidance
High certainty; wording is the builder's. Keep each addition short; this guide is read by people, not agents.

## Red-Green Proof
**RED prompt/case:** Read `skills/do-work-board/docs/board-guide.md`: the switcher is described as "Board / Calendar / Testing", nothing says a page can be opened by link, and the Testing section does not mention the disk line.
**Why RED now:** The docs predate 0.305.61 and 0.305.62.
**GREEN when:** The guide names all six pages, documents the page and lens links with their fragments, and the Testing section describes the disk line; `grep -n "Board / Calendar / Testing" skills/do-work-board/docs/board-guide.md` returns nothing.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: views and static output. The bullet this REQ cites (REQ-232, describe the board by what its switcher covers) carries no family marker, so no narrower entry exists; the REQ's Constraints name that bullet, and the builder reads it directly.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: three insertions in the board guide (switcher sentence, a links paragraph before the lens section, a disk-line sentence in the Testing view) and one clause in board.md Step 6, each checked against `web/board-controls.js` and `web/board-testing.js`.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: only `board-guide.md` and `board.md` changed.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: `git diff --stat` 2 files, 5 insertions, 3 deletions; no linter covers Markdown here; re-read both files: fragment list matches `boardStateFromFragment`, replaceState with the location.replace fallback for file:// matches `writeBoardFragment`, thresholds match `verify.go` (strict below 10 and 3 GiB), "(at generation)" and skip text match `diskSpaceLineFor`; the RED grep returns nothing.

## Full Context
See `do-work/user-requests/UR-133/input.md` for complete verbatim input.

*Source: capture both as REQs (finding F3 of the REQ-626 review and finding M1 of the REQ-627 review)*

---

## Triage

**Route: A** - Simple

**Reasoning:** Docs-only change that names its two files and lists every sentence to add; the only discovery is checking each sentence against `web/board-controls.js` and `web/board-testing.js`, which the builder does inline.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)

**What was done:** The board guide's toolbar paragraph now names the page switcher by what it covers, all six pages. A new paragraph lists each page and Board lens link (`#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`), says a click updates the address bar without Back steps, that filters are not in the link, and that links work on a static snapshot from disk. The Testing view section and board.md Step 6 now mention the free-disk-space line with its thresholds; the guide also gives the static label, the unmeasured case, and the page-load value in serve mode. Merge range dda564c3..4c8e249a (builder commit 4e87f8cd, merge 4c8e249a).

## Qualification

**Diff range:** dda564c3..4c8e249a (builder commit 4e87f8cd, merge 4c8e249a)
**Gate records:** qualify satisfied. Route A, so no Scope comparison.
**Warnings judged:** none.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. The toolbar sentence names all six switcher pages in the template's button order. The links paragraph lists exactly the eight fragments in `boardStateFromFragment`; "without adding Back steps" matches `writeBoardFragment` (replaceState, with `location.replace` as the file:// fallback, both of which add no history entry); filters are not in the fragment table. The disk sentence matches `diskSpaceLineFor` and the strict thresholds in `verify.go`. board.md Step 6 describes the toolbar's contents, so it gets one clause. Nothing describes behaviour the code lacks.
**P-A-U honesty:** the three boxes were ticked by the orchestrator, who also played the builder role on branch `worktree-agent-REQ-628-board-guide`; APPLY cross-checked against `git diff --stat dda564c3..4c8e249a` (two files, nothing under do-work/).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 4c8e249a
**Result:** ✓ All passing — exit 0, gate wall 128s; stages queue-kanban-fast-tests and do-work-cli-fast-tests EXECUTING (no_prior_evidence). Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-02-204757/helpers/probe-628.sh` (the Red-Green Proof as a check: the old switcher phrase is gone, all eight fragments are named, both files mention the disk line) → exit 0 (advance probe record satisfied). Non-behavioral docs change, so no harness test; the probe is regression proof, and it fails on the pre-change tree at dda564c3 because the old phrase is present.

**Red-green validation:** traced to `## Red-Green Proof`: RED `grep -n "Board / Calendar / Testing"` matched line 25 at dda564c3 → GREEN returns nothing at 4c8e249a; the guide names six pages, eight fragments and the disk line.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: dda564c3..4c8e249a
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-02T20:52:52Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A (docs-only; probe-628.sh is the regression proof) |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

Every sentence the diff adds was checked against the code. Switcher: `web/template.html:82-97` has six `data-view-target` buttons in the order the guide lists. Links: the eight fragments match `boardStateFromFragment` (`web/board-controls.js:154-167`). "No Back steps" matches `writeBoardFragment` (`replaceState`, with `location.replace` as the file:// fallback, `:205-218`). The fragment is read on load and on `hashchange` (`:286-300`). Filters are not in the fragment table. Disk line: the text shape matches `diskSpaceLineFor` (`web/board-testing.js:114-126`) and `formatGibibytes` (`verify.go:351`, `%.1f GiB`). The thresholds are strict `<` at 10 GiB and 3 GiB (`verify.go:227-228, 341-343`). The colours are `--accent-pending` #d8a24a / #8f5e10 (amber) and `--accent-blocked` #d97a59 / #bd5138 (a coral or brick red), so "amber" and "red" are fair names. "Measured on page load until reload" matches `serve.go:135-163`: the data is rebuilt on every board-data.js request, and testing POSTs update the local data without reloading it. Restatement sweep: no shipped doc still says "Board / Calendar / Testing". `board.md:3,14` already use the switcher-covers wording.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: Partial page lists remain in shipped intros, against the REQ-232 lesson ("describe the board by what its switcher covers"). These are `skills/do-work-board/docs/board-guide.md:3` ("a Kanban board … plus a queue activity calendar and a testing track") and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:3` ("Kanban board + queue activity calendar"). Both are older than this diff and outside the REQ's requirements. Line 25 now lists all six pages, so a reader is not misled — impact-negligible → report only
- N1 (nit): The guide says `not measured` appears "when the platform cannot measure it". The code shows `disk: not measured on <os>` for an unsupported platform, `disk: not measured: <error>` when a measurement fails on a supported platform, and `disk: not measured` when there is no reading (`verify.go:290-293`, `board-testing.js:115-119`). The failure case is not covered, but the wording is close enough for a person reading the guide — impact-negligible → report only

**Acceptance:** Pass. `probe-628.sh` re-run exits 0. The RED phrase is gone, all eight fragments are named, and both files mention the disk line.
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Discovered Tasks

From the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- M1: the intros of `skills/do-work-board/docs/board-guide.md` (line 3) and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (line 3) still name only some of the board's pages. — impact-negligible → report only
- N1: the guide's "not measured" wording covers an unsupported platform; a failed measurement on a supported one reads `not measured: <error>`. — impact-negligible → report only

## Lessons Learned

**What worked:** Checking each new sentence against the fragment table and the line helper kept the guide from describing behaviour the code lacks.
**What didn't:** Nothing failed.
**Worth knowing:** Two more partial page lists remain in intros (review M1); the switcher-coverage lesson applies to every intro, not only the toolbar paragraph.

## Orientation

Now the board guide describes all six pages, their links, and the Testing page's disk line; lives in the board package's user docs (`skills/do-work-board/docs/board-guide.md`, `skills/do-work-board/actions/board.md`). Not a map change. Prime spot-check: `prime-kanban-board.md`, `prime-do-kanban.md` and `prime-releases.md` paths still exist; none was made stale by this change.

## Heavy Verification Plan

- Base revision: dda564c357ebb69a2261661b8b10fdb43946998d
- Target revision: 4c8e249acfa058bd75e6ed8a80175172ff31fbd5 (landed in `commit:`)
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

## Heavy Verification Result

- Target revision: 4c8e249acfa058bd75e6ed8a80175172ff31fbd5
- Execution revision: 069f922e64c1ed14ca31211642e395a784561f57 (detached drain checkout `.git/work-run-2026-10-02-204757/drain-head`, one run for REQ-628 and REQ-629 at the integration tip containing both merges)
- staged-skills: exit 0, executed, 33s

Green: every selected lane present, exit 0, none skipped, none reused. An earlier run at 8f9f3099 (before REQ-629's D-05 prose fixes) was also green on all four lanes.
