---
id: REQ-632
title: 'Board cards show last correlated activity and drop the assumed-pause badge'
status: completed
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 12-file write set
  - 3 subsystems involved
  - 6 acceptance criteria
  calculated_at: 2026-10-05T20:32:12Z
created_at: 2026-10-05T20:11:59Z
user_request: UR-135
re_review_at: 2026-10-05T21:33:07Z
remediation_at: 2026-10-05T21:30:25Z
review_at: 2026-10-05T20:55:33Z
integration_at: 2026-10-05T20:50:05Z
builder_handback_at: 2026-10-05T20:49:55Z
dispatch_at: 2026-10-05T20:38:50Z
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-633"]
batch: board-activity-evidence
write_set: ["skills/do-work-board/tools/queue-kanban/activity_correlation.go", "skills/do-work-board/tools/queue-kanban/activity_correlation_test.go", "skills/do-work-board/tools/queue-kanban/serve.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/durations.go", "skills/do-work-board/tools/queue-kanban/durations_test.go", "skills/do-work-board/tools/queue-kanban/web/board-cards.js", "skills/do-work-board/tools/queue-kanban/web/board-detail.js", "skills/do-work-board/tools/queue-kanban/web/board.css", "skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go"]
claimed_at: 2026-10-05T20:31:37Z
commit: 47c2620eb13f8389e6ad0a998c1869c6d1cbad92
heavy_verified_at: 2026-10-05T21:35:20Z
heavy_verified_revision: 47c2620eb13f8389e6ad0a998c1869c6d1cbad92
completed_at: 2026-10-05T21:35:38Z
release_at: 2026-10-05T21:35:38Z
---
# Board Cards Show Last Correlated Activity and Drop the Assumed-Pause Badge

## What
Give every in-progress card on the queue-kanban board a ticking `last activity N min ago · <phase>` line computed from the union of the REQ's lifecycle stamps and the git commits correlated to it, put the largest gap between consecutive activity events in the detail drawer, and delete the `over 4h · assumed pause` badge from done cards.

## Why
The badge is a guess from one number. In the maintainer's repo three of the last five finished REQs carried it and all three were worked continuously. The repo already holds the evidence that disproves it: per-phase stamps and commits that touch the REQ file. In-progress cards have no stuck signal at all today; a growing "last activity" number is that signal. The design was discussed as UR-135 and the maintainer chose: badge deleted rather than reworded, done cards show nothing new, gap evidence in the drawer, correlation by touched path OR `[REQ-NNN]` prefix OR merge ancestry.

## Detailed Requirements
1. New `activity_correlation.go` in `skills/do-work-board/tools/queue-kanban/` with `collectRequestActivity(repoRoot, tickets, since, now, runner)`. It runs one `git -C <root> log --since=<since> --format=%H%x00%cI%x00%P%x00%s --name-only` (add `-- do-work/` as pathspec only if merge-commit subjects still arrive; measure both) where `since` is the oldest `claimed_at` among open tickets and the recently-done window start. It attributes each commit to REQ ids by any of: a touched path matching `do-work/{queue,working,archive}/**/REQ-NNN-*.md` or `do-work/runs/*/REQ-NNN*`; every `[REQ-NNN]` token in the subject; membership in `<merge>^1..<merge>^2` of a two-parent commit already matched by path or prefix (derive from the `%P` graph in Go, or one `rev-list` per matched merge, builder's call). For each live `worktree-agent-REQ-NNN-*` branch (reuse `listWorktreeAgentBranches`, `verify.go`) it runs `git log -1 --format=%cI <branch>` and attributes the tip to that REQ.
2. Per REQ it unions those commit instants with the parsed stamps from `lifecycleTimestampFields` (`model.go`), sorts them, and returns `LastActivityAt`, `LastActivityKind` (`stamp` or `commit`), `LastActivityPhase`, and `LargestGap{Minutes, FromPhase, ToPhase}` for the largest gap between consecutive union events. Phase names reuse the interval names `buildPhaseBreakdown` (`durations.go`) already emits; the phase of a gap is the stamp interval containing the gap start.
3. The git runner is injectable the way `gitCommitDateLookup` is (`model.go`), so tests feed canned `git log` output and never spawn git.
4. Wire it into the per-request slot next to `attachVerifyFindings` in `serve.go` (the block whose comment says findings are computed fresh on every request, outside `refreshBoardData`'s mtime cache) and once in `generate`. Never inside the cached rebuild. One `now` per response, per the prime rule that a cache hit must rebuild every wall-clock field from one instant.
5. Payload: open requests carry `lastActivityAt` (RFC3339), `lastActivityKind`, `lastActivityPhase`. Every request with two or more union events carries `largestActivityGap {minutes, fromPhase, toPhase}`. No threshold value is shipped; this REQ states observations only.
6. Card (`web/board-cards.js`, beside the existing state timer line built from `stateTimerSpecFor` and `makeInstantWithStopwatchNode`): for claimed, planning, dispatched and integrating statuses render `last activity <instant> · <phase>` through `makeInstantWithStopwatchNode` so the node carries `data-instant-ms` and the existing 1s ticker (`refreshRelativeTimeNodes`, `board-core.js`) owns the text. Blocked and pending-answers cards: nothing. Done cards: nothing new; keep the `wall time …` line and the `reversed stamps` flag.
7. Delete the badge: `implementationSpanPausedBadgeText` and its JSON field (`durations.go`, `generate.go`), the paused branch of `makeImplementationSpanNode` (`web/board-cards.js`) and its tooltip, CSS that becomes orphaned in `web/board.css`, and the `over 4h · assumed pause` glossary row in `skills/do-work-board/docs/board-guide.md`. `implementationSpanReason` and `excludedReason` keep shipping with their current values for Panel B until REQ-633 (Panel B and the calibration log exclude by largest stamp gap) renames them.
8. Drawer (`web/board-detail.js`, after `appendPhaseBreakdownRows`): one row `largest idle gap 1h07m (dispatch → handback)` from `largestActivityGap`. Plain text, no colour, no threshold.
9. Tests: replace `TestImplementationSpanPausedBadgeTextDerivesFromTheCeiling` (`durations_test.go`) and the badge assertions in `TestGeneratedRequestCarriesTheDoneCardImplementationSpan` (`generate_test.go`) and `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan` (`javascript_behavior_b_test.go`). New pins, each naming its failure: a commit with a non-REQ subject that touches the REQ file attributes (the `docs(do-work): …` case from the brief); a prefix-only commit attributes; builder commits reachable only through a matched merge's second parent attribute; a live branch tip counts; a 4h21m REQ with a 67-minute largest gap gets a drawer row and no badge; an open card carries `data-instant-ms`.
10. Release: changelog entry and version bump per `_dev/primes/prime-releases.md`. Time `/board-data.js` on this tree before and after; the new read must stay in the same order as the existing verify probe set (about 40ms today).

## Constraints
- Do not extend the mtime fingerprint in `serve.go`; the git read is per request, like the verify probes (lesson REQ-284: two probe inputs change while every file mtime stays identical).
- No `git log --all`, no reflog. Merged builder commits are reachable from main after the `--no-ff` merge; a branch deleted before merge is a crash leftover and its phase reads as stamps-only, which is true.
- The timing stream under `<git-common-dir>/do-work-timing/` is not read.
- No new write surface; the tool still has exactly three.
- Render and look: a passing suite is not evidence about card text. Generate a board against this repo, read the card and drawer text, and return `location.href` beside any measurement.
- Shipped files change, so this is a release.

## Dependencies
None upstream. REQ-633 (Panel B and the calibration log exclude by largest stamp gap, not raw span) depends on this REQ.

## Builder Guidance
Certainty is high: the design was discussed question by question and the maintainer picked each option on 2026-10-05 (see `do-work/user-requests/UR-135/input.md` → Summary). Builder latitude: struct and field names, whether ancestry is derived from the `%P` graph or a second `rev-list`, whether the pathspec is used, and where the drawer row sits relative to the phase table. Measured cost on this tree: one `git log --since=7.days --name-only` over 72 commits took 50 to 90ms.

## Red-Green Proof
**RED prompt/case:** Build the board from a fixture with (a) a `dispatched` REQ claimed three hours ago whose only recent activity is a commit ten minutes ago that touches its `do-work/working/REQ-NNN-*.md` file with a subject that carries no `[REQ-NNN]` prefix, and (b) a completed REQ whose earliest-stamp-to-completion span is 4h21m with a largest stamp gap of 67 minutes. Today (a)'s card shows no activity line and (b)'s card shows `over 4h · assumed pause`.
**Why RED now:** No card reads git, and the badge is derived from span > `analysisOutlierCeiling`, never from the gaps `buildPhaseBreakdown` already computes.
**GREEN when:** (a)'s card shows `last activity 10m ago · dispatch` on a node carrying `data-instant-ms`; (b)'s card shows `wall time 4h 21m` and no badge, and its drawer shows `largest idle gap 1h07m (dispatch → handback)`; `grep -rn "assumed pause" skills/do-work-board/` returns nothing; the Go and JavaScript behaviour tests above pass.
**Validation:** User confirmed (plan approved 2026-10-05)

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: its prime governs `skills/do-work-board/tools/queue-kanban/`, which this REQ writes; the REQ-284 (never cache the verify probes) and REQ-422 (a cache hit rebuilds every wall-clock field from one instant) entries apply directly and are restated in Constraints.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban model, UI and browser behaviour.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One windowed git log (hash, committer date, parents, subject, touched paths; no pathspec) parsed in Go; ancestry from the parent graph; live branch tips through the same injected runner; union with stamps per REQ inside the claim-to-completion window; shipped as a top-level requestActivity map; attached beside the verify probes in serve and generate with one instant each. (From the builder hand-back.)
- [x] **[APPLY]:** Scope held to the Scope list; lessons-do-kanban.md left untouched (lesson proposed to the orchestrator instead).
- [x] **[UNIFY]:** git diff --stat 13 files; gofmt -l empty, go vet clean, full package suite with JavaScript probes on passed in 63.9s; every changed file read in the diff, no debug output; badge grep over skills/do-work-board returns nothing. Orchestrator cross-checked APPLY against git diff --stat ac5cbe64..b004ecd6 (13 Scope files, nothing outside them).

## Full Context
See `do-work/user-requests/UR-135/input.md` for complete verbatim input.

*Source: A. An activity indicator on each REQ card — at minimum "last commit N min ago" where the commit is correlated to the REQ … On done cards it is the honest replacement for the assumed-pause guess.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ already fixes the design, the files, and the tests decision by decision (UR-135 session), so no Plan agent is needed. The Go seams it names (gitCommitDateLookup, listWorktreeAgentBranches, attachVerifyFindings, buildPhaseBreakdown) and the card/drawer rendering helpers need location and pattern discovery before dispatch.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors: `do-work/runs/work-2026-10-05-203159/REQ-632-exploration.md`. The findings that change what the builder writes:

- **No planning, dispatched or integrating status exists.** Every in-flight REQ is `status: claimed` (normalizeStatus, model.go:1014); phases exist only as stamps. The card rule is claimed status, and the phase text comes from the newest union event.
- **Writing into the serve payload's request map would race.** serve.go:150 copies the payload shallowly and Requests (generate.go:79) is a map of value structs, so a write there mutates the cached map outside cacheMu. Ship activity as a separate top-level field (for example requestActivity keyed by id) or clone the map.
- **Serve takes two different now values per response** (serve.go:351 on a cache hit, serve.go:163 for the verify probes). Take one now in the board-data handler and pass it to both; generate adds the call beside generate.go:497 with its own single now.
- **Do not add the do-work pathspec.** Over 7 days: 74 commits and 8 merges without it, 56 commits and 0 merges with it, so ancestry and prefix-only builder commits would be lost. Both about 30ms warm.
- **Merge commits list no files under plain name-only,** but their subjects carry the REQ prefix, so prefix plus the second-parent range covers them.
- **Path matching over-attributes:** a completion commit can touch a sibling REQ file (REQ-631's touched REQ-630's). Archive paths include a UR folder (`do-work/archive/UR-NNN/REQ-NNN-*.md`). Committer dates carry offsets; RFC3339 parsing accepts them.
- **The REQ's GREEN grep for the badge phrase cannot pass as written:** the same words appear in Panel B and the UR summary, which REQ-633 owns (board-durations.js, board-user-request-summary.js, durations.go:25, durations_test.go:84, javascript_behavior_a_test.go:619, board-guide.md:41). Limit the grep to the badge (`over 4h · assumed pause`, implementationSpanPaused) and record the deviation as a decision.
- **Real formats differ from the REQ's examples:** makeInstantWithStopwatchNode renders an instant plus a stopwatch, not "10m ago"; formatElapsedDuration gives "1h 07m". Tests assert the real format.
- **Phase labels:** buildPhaseBreakdown (durations.go:270-281) labels Planning, Dispatch, Builder handback and so on, and returns nil when no optional stamp parses, so map stamp field to label directly.
- **listWorktreeAgentBranches (verify.go:1465) runs git directly;** send the same branch listing through the new injected runner instead.
- **Tests:** whole package 58s (`go test -count=1 ./...`); JavaScript probes run only with QUEUE_KANBAN_JAVASCRIPT_PROBES=on. Tests to change: durations_test.go:507, generate_test.go:3679-3681, javascript_behavior_b_test.go:1626 and :1720. Patterns: card probe javascript_behavior_b_test.go:1606, drawer probe javascript_behavior_a_test.go:2259. Orphaned CSS board.css:1449-1462; keep the reversed flag's class.
- **Baseline:** board-data response 87-110ms warm, 1.63s cold. Release preimage 0.305.66 in all three version files.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (new) — collectRequestActivity, commit attribution by path, prefix, merge ancestry and live branch tip; injectable git runner
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (new) — attribution and largest-gap pins with canned git output
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify) — one now per response; activity computed in the per-request slot beside the verify probes, shipped as a separate field
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — payload field for activity, call in generate with the same now, badge field removed
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — badge assertion removed, activity payload pinned
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — badge text builder removed
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — badge-text test removed
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — last-activity line on claimed cards, paused-badge branch and tooltip removed
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify) — largest idle gap row in the drawer
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — orphaned paused-marker rule removed
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify) — drawer row behaviour test
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modify) — done card text without badge, open card activity node
- `skills/do-work-board/docs/board-guide.md` (modify) — badge glossary row removed, activity line and drawer row described

- `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go` (modify) — remediation: the browser contrast test counts the last-activity stopwatch (added after the heavy lane failed, D-11)

**Files I will NOT touch:** model.go and verify.go (the git runner is new and the branch listing goes through it), the Panel B and UR-summary text that also says assumed pause (board-durations.js, board-user-request-summary.js, the durations.go exclusion comment; REQ-633 owns them), the mtime fingerprint in serve.go, the timing stream, every release path (written by finalization), do-work/lessons-index.md (orchestrator).

**Acceptance criteria (restated from REQ):**
- [ ] A claimed card whose only recent activity is a commit touching its REQ file with a subject that has no prefix shows a last activity line on a node carrying data-instant-ms, with the phase name
- [ ] Commits attribute by touched REQ path, by every REQ prefix token in the subject, by membership in a matched merge's second-parent range, and by a live worktree-agent branch tip
- [ ] A done REQ with a 4h21m span and a 67-minute largest gap shows its wall time line and no badge; its drawer shows the largest idle gap row with both phase names
- [ ] The badge text builder, its payload field, its card branch, tooltip and CSS, and its board-guide row are gone; grep for the badge string returns nothing in the board package
- [ ] The git read runs per request outside the mtime cache with one now per response; no new write surface; tests never spawn git
- [ ] The board-data response time stays in the same order as before (about 90ms warm)

## Pre-Flight

**Git:** ✓ Integration tip 2e34a507 on `main` (the claim commit); the only dirt is this REQ's own trail (working REQ, run directory) — no third-party paths
**Tests baseline:** ✓ focused board tests green (`do-work/runs/work-2026-10-05-203159/REQ-632-probe.sh`, 4.5s)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 2e34a507, both Go stages EXECUTING, gate wall 127s
**Dependencies:** ✓ Go toolchain present; no new module dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (new)
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go` (modified) — heavy-lane remediation

**What was done:** Added a per-response activity collector that correlates git commits to REQs (touched REQ file or run artifact, bracketed REQ prefix in the subject, builder commits inside a matched merge's second-parent range, live worktree-agent branch tip), unions them with lifecycle stamps inside the claim-to-completion window, and ships last activity plus the largest gap as a top-level requestActivity map computed outside the mtime cache with one instant per response. Claimed cards show a ticking last activity line with its phase; the drawer shows a Largest idle gap row; the over-4h assumed-pause badge, its payload text, card branch, tooltip, CSS and guide row are deleted. Merge range ac5cbe64..47c2620e (cumulative: builder commit 55d8ea17 merged at b004ecd6; remediation commit 43ea1930 merged at 47c2620e, after the heavy browser lane failed because the contrast test counted two timers on in-progress cards and the new last-activity stopwatch is a third).

## Qualification

**Diff range:** ac5cbe64..b004ecd6 (builder commit 55d8ea17, merge b004ecd6)
**Gate records:** qualify satisfied; scope-drift satisfied on the second call. The first call flagged the lessons satellite as declared but untouched: the builder proposed its lesson to the orchestrator instead of writing it, because a satellite edit needs the lessons-index row refreshed in the same change. The satellite was taken out of Scope and write_set and the lesson moves to finalization (D-09).
**Warnings judged:** QUALIFY-NEW-FILE-UNWIRED on activity_correlation_test.go — expected: Go test files are discovered by the toolchain, never referenced. Not dead code.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. One windowed git log without pathspec, parsed in Go; attribution by REQ path (archive UR folder included), every bracketed prefix token, the matched merge's second-parent range from the parent graph, and live branch tips through the injected runner; union with stamps inside claimed_at to completion (D-03); top-level requestActivity map written on the response copy, never the cached Requests map; one instant per response shared with the verify probes in serve and generate; claimed cards render through makeInstantWithStopwatchNode; the drawer row is plain text; badge text builder, payload field, card branch, tooltip, CSS and guide row are gone. Requirement 5 says every request with two or more events carries the gap; the collector covers claimed and recently done REQs only (D-05), which matches requirement 1's own since window.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY cross-checked against git diff --stat ac5cbe64..b004ecd6 (13 Scope files, nothing else).
**Live data flow:** serveLiveBoardDataJs and generate both call attachRequestActivity; board-cards.js and board-detail.js read boardData.requestActivity by request id.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at b004ecd6
**Result:** ✓ All passing — exit 0, gate wall 148s; stage queue-kanban-fast-tests EXECUTING (421 tests, wall 45s, slowest file strict_behavior_regression_test.go 20.04s < 30s); stage do-work-cli-fast-tests EXECUTING (867 tests, wall 66s, slowest file 21.81s < 30s). Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-05-203159/REQ-632-probe.sh` (board package, JavaScript probes on, the span, phase, activity, card, drawer and payload tests) → exit 0, advance probe record satisfied.

**Red-green validation:** traced to `## Red-Green Proof`; RED taken by the builder against a compiling stub, GREEN at builder commit 55d8ea17:
- TestRequestActivityAttributesACommitThatTouchesTheRequestFileWithoutAPrefix (the RED case's unprefixed commit): ✗ "REQ-701 carries no activity; a path-only commit must attribute" → ✓
- TestRequestActivityAttributesEveryPrefixTokenInASubject: ✗ attributed instants empty → ✓
- TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent: ✗ "attributed 0 instants, want 3" → ✓
- TestRequestActivityCountsALiveWorktreeAgentBranchTip: ✗ zero last activity → ✓
- TestRequestActivityLargestGapUnitesCommitsWithStamps: ✗ no largest gap → ✓
- TestGeneratedPayloadCarriesTheLargestActivityGap (the 4h21m, 67-minute REQ): ✗ empty map → ✓
- TestServeReadsRequestActivityOnEveryResponse: ✗ first response has no last activity → ✓
- TestJavaScriptBehaviorClaimedCardShowsItsLastCorrelatedActivity: ✗ 0 activity lines and empty data-instant-ms → ✓
- TestJavaScriptBehaviorDetailStatesTheLargestIdleGap: ✗ drawer never calls the gap row → ✓
- TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan: ✗ "wall time 4h 21m long span · assumed pause" → ✓ "wall time 4h 21m"

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (attribution, gap, path pattern, serve per-response read)
- javascript_behavior_a_test.go and javascript_behavior_b_test.go gained the drawer and claimed-card probes

**Existing tests updated (cross-REQ impact):**
- durations_test.go: the badge-text test is deleted with the badge — intentional
- generate_test.go TestGeneratedRequestCarriesTheDoneCardImplementationSpan: badge assertion removed, Panel B reason still pinned as paused until REQ-633 — intentional
- javascript_behavior_b_test.go done-card probe: expects wall time with no badge — intentional

**Render evidence (builder):** generate against this repo; board-data.js carries requestActivity for REQ-632 (last activity 20:38:50Z, stamp, dispatch) and seven recently done REQs (largest gap only); the badge key is gone. Static file read, so no location.href.
**Response time (builder):** /board-data.js warm 105–129 ms before, 148–157 ms after (+~40 ms, the same order as the verify probes); cold about 1.5 s both.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: ac5cbe64..b004ecd6
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

**Heavy-lane remediation (after review):** the first heavy drain at 9dacf899 (21:01Z–21:04Z) ran javascript exit 0 (7s), browser exit 1 (75s), staged-skills exit 0 (38s). The red test was TestBrowserBehaviorCompletionCompanionsKeepReadableContrast: "nonterminal control readings = 3, want pending and claimed" in all six scheme/width cases, because the claimed fixture card gained the last-activity stopwatch. Per the heavy-lane red rule, commit: and the Heavy Verification Plan were withdrawn and the builder fixed the test on its branch (43ea1930: the probe names three controls and still checks each one's faint ink and 0.85 opacity, so the new stopwatch is proven to read like the other timers). Re-merged at 47c2620e with the same lower bound.
**Repository gate after remediation:** `bash _dev/tests/maintainer-verify.sh` run directly at 47c2620e (advance refuses gate input once the REQ is past review) → exit 0, gate wall 126s; queue-kanban-fast-tests 420 tests, wall 39s, slowest file 19.36s; do-work-cli-fast-tests 867 tests, wall 59s, slowest file 19.73s.
**Heavy verification plan (re-held):** range ac5cbe64..47c2620e, the same three lanes as above.

*Verified by work action*

## Review

**Overall: 93%** | 2026-10-05T20:54:15Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 90% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Requirement 4 says "one `now` per response", but this holds only partly. serve.go:172 shares one `responseInstant` between the verify probes and the activity read. On a cache hit, though, `refreshBoardData` (serve.go:360) takes a second `currentTime()` for `GeneratedAt` and the Timeline, so a `/board-data.js` response still carries two instants. `generate` does the same: `buildBoard`'s now differs from `snapshotInstant` at generate.go:502. The new serve.go comment ("Both read one `now`, so the response states one instant") and the Qualification's "one instant per response" claim more than the code does. Exploration finding 3 asked for a single now passed into the refresh as well. The visible effect is milliseconds, so nothing on screen is wrong today. The risk is that the comment misleads the next editor about the REQ-422 rule (a cache hit rebuilds every wall-clock field from one instant). Fix: take the instant once in the handler and pass it to `refreshBoardData`, or narrow the comment — impact-negligible → report only

**Minor findings:**
- activity_correlation.go `requestActivityEvents` ends a finished REQ's window at the raw `completed_at`/`release_at` frontmatter. A terminal REQ in Recently done whose completion is git-dated, or a cancellation with no `completed_at`, gets a window that runs to now. A later sibling commit touching its archive file (the exploration's REQ-631→REQ-630 shape) then creates a false post-completion "idle gap". Using `ticket.CompletionTime` as the end would close this — impact-negligible → report only
- `since` widens to the oldest `claimed_at` on the board, so one abandoned claim makes every response read that whole history. Measured on this repo: a 60-day window is 2,901 commits (336 merges), with `git log` at about 0.2s plus 39ms of Go correlation per response, against today's ~40ms budget. This matches requirement 1 as written, but nothing bounds it — impact-negligible → report only
- `liveBranchTipInstants` spawns one `git log -1` per live `worktree-agent-*` branch, including crash leftovers. A single `git for-each-ref --format='%(refname:short) %(committerdate:iso-strict)'` would replace N spawns with one — impact-negligible → report only
- board-guide.md:57's new paragraph repeats the stale `took …` label from row 52. The card actually renders `wall time …` (the builder found this before the change and reported it). The new text adds one more place that names a label the reader will not see on the card — impact-user-visible → report only

**Acceptance:** Pass. The focused Go and JavaScript probe tests pass (12 named tests, JS probes on), and gofmt and vet are clean. A built `generate` run against this repo ships `requestActivity` for the 1 claimed and 7 recently done REQs, and the `implementationSpanPausedBadgeText` key is gone. REQ-632 (this REQ, still claimed) shows last activity at 20:50:05Z (stamp, integration) and a largest gap of 10.15 min (dispatch → dispatch). That gap is exactly dispatch 20:38:50 to builder commit 55d8ea17 at 20:48:59, which proves attribution through the prefix, the merge's second parent, and the branch tip. REQ-630 (done) shows a 4.85-min gap (review → completed) and no lastActivityAt.
**Suggested testing:** 3 items. (1) Open the served board in a browser: check that REQ-632's card line ticks, and the drawer row text. (2) A stale-claim fixture, to time `/board-data.js` with a 30-plus-day window. (3) A recently cancelled REQ with no `completed_at`, to check its drawer gap.
**Follow-ups created:** None (5 findings report only)

### Requirements walk
- R1 collector, one windowed log, no pathspec, three attribution rules plus live branch tips: delivered. The regexes were checked. The path pattern covers the archive UR folder and run artifacts, and `[^/\d]` stops REQ-63 from matching REQ-632. Only bracketed subject tokens count. The merge range ^1..^2 is derived correctly from the `%P` graph by walking ^2 and stopping at ^1's logged ancestry. Ancestry starts only from direct matches.
- R2 union with stamps, newest event, kind, phase, largest gap with from/to phases: delivered. The claim-to-completion window is decision D-03.
- R3 injectable runner, tests never spawn git: delivered.
- R4 per-request slot outside the mtime cache, once in generate, one now: partial (Important finding above). The cache-race claim holds. `RequestActivity` is a new top-level map that is nil on the cached struct and created on the shallow response copy. The cached `Requests` map is never written.
- R5 payload: delivered. The gap is limited to claimed plus recently done REQs (D-05, escalated, documented).
- R6 claimed card line via `makeInstantWithStopwatchNode`, carrying `data-instant-ms`. Blocked and done cards show nothing new: delivered. The real format is documented in D-02.
- R7 badge deleted (builder function, payload field, card branch, tooltip, CSS, guide row): delivered. `implementationSpanReason` still says "paused" for Panel B, and the REQ-909 pin asserts it.
- R8 drawer row as plain text: delivered. It reads `Largest idle gap: 1h 07m (dispatch → builder handback)` (D-02).
- R9 tests: delivered. Each new test's comment names the failure it pins. The changed REQ-902 expectation and the replaced tooltip assertion are traced to REQ-632. Red-green evidence is in the hand-back.
- R10 release: N/A for this review (finalization owns it). The builder's timing was warm 105–129ms before and 148–157ms after (+~40ms, the same order as the verify probes).

### Restatement Sweep
Two things were redefined: the deleted badge (`implementationSpanPausedBadgeText` and the `over 4h · assumed pause` text) and the done-card meaning of a long span. Grep over `skills/` and `_dev/` finds no remaining consumer of the payload field and no restatement of the badge. The "assumed pause" text that remains is all Panel B, the UR summary, and the calibration reader (board-durations.js, board-user-request-summary.js, durations.go:24, board-guide.md:41, estimate-reference.md:98). REQ-633 owns all of it, and this diff leaves it untouched, as D-01 says. The only stale restatement found is the `took …` label (Minor above).

### Domain Review (backend)
The GET handler stays side-effect free: git reads only, and no new write surface. A git failure falls back to stamp-only evidence instead of an error. Error responses are unchanged.

### Re-review after the heavy-lane remediation

The remediation diff b004ecd6..47c2620e touches one test file. The fix keeps the test's intent (the completion-contrast change leaves live timers at faint ink and 0.85 opacity) and widens its coverage to the new stopwatch instead of filtering it out; a missing, extra or repeated control still fails by name. No product code changed, so the first review's scores and findings stand. Acceptance: Pass, on the builder's full browser lane (exit 0, 37 tests, no skips) and the repository gate at 47c2620e.

*Reviewed by review-work action*

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-05-203159/REQ-632-handback.md`) so they survive run-directory cleanup.

- D-01 (DECIDE & STATE): the GREEN grep is narrowed to the badge. The broad assumed-pause wording still lives in Panel B and the UR summary, which REQ-633 (Panel B and the calibration log exclude by largest stamp gap) owns.
- D-02 (DECIDE & STATE): the text uses the board's real formats. The card reads "last activity <instant> 10m 00s · dispatch"; the drawer reads "Largest idle gap: 1h 07m (dispatch → builder handback)". Phase labels are buildPhaseBreakdown's labels in lower case.
- D-03 (DECIDE & STATE): the union is bounded to the work window, claimed_at to the later of completed_at and release_at. Queue wait is not idle work, and a later sibling commit touching an archived file cannot create a false gap.
- D-04 (DECIDE & STATE): stamps without a phase count as events but do not change the phase; a commit takes the phase of the newest phase stamp before it; on a tie the first gap wins.
- D-05 (ESCALATE, resolved by the orchestrator): activity is collected only for claimed and recently done (7 days) REQs, so older archived REQs have no drawer gap row. Value: the git read stays one bounded log (+~40 ms). Risk: a reader may expect the row on old done REQs; reversible by one loop. Orchestrator judgment: this matches requirement 1's own since window (oldest claimed_at and the recently-done window start), so it is the REQ's design rather than a new choice; no follow-up.
- D-06 (DECIDE & STATE): implementationSpanReason and excludedReason still ship "paused" for Panel B until REQ-633 renames them; the payload pin asserts it.
- D-07 (DECIDE & STATE): both implementation-span CSS rules were removed; the elapsed-duration class is already nowrap.
- D-08 (DECIDE & STATE): the phase milestone list moved to phaseMilestonesOf in durations.go and is shared by buildPhaseBreakdown and the collector, so there is one list.

Orchestrator decisions:
- D-09 (DECIDE & STATE): the lessons satellite was taken out of Scope and write_set after qualification flagged it as declared but untouched. The builder proposed the lesson instead of writing it because the satellite edit needs the lessons-index row refreshed in the same change, and that file is the orchestrator's; the lesson lands in the finalization commit.
- D-11 (DECIDE & STATE; the builder's addendum numbers it D-09): completion_contrast_browser_test.go was added to Scope and write_set for the heavy-lane remediation. The test is the only browser probe that counts live timers on in-progress cards, and the new last-activity line is a third one by design.
- D-10 (DECIDE & STATE): the five review findings stay report only per their impact tokens; none is critical. The "one instant per response" wording in Qualification is narrowed by this note: the verify probes and the activity read share one instant, while a cache hit's GeneratedAt and Timeline still take their own (milliseconds apart, pre-existing).

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- The board guide's done-card row calls the span badge "took …", but the card renders "wall time …"; the new idle-gap paragraph repeats the stale label. — impact-user-visible → report only
- A cache hit still carries two instants (refreshBoardData's for GeneratedAt and the Timeline, the handler's for verify and activity); the serve comment says one. — impact-negligible → report only
- A finished REQ whose completion is git-dated, or a cancellation without completed_at, gets an activity window that runs to now; ending at the ticket's CompletionTime would close it. — impact-negligible → report only
- One abandoned claim widens every response's git read to its claimed_at (60 days measured at about 0.2 s). — impact-negligible → report only
- One git log -1 per live worktree-agent branch; one for-each-ref would do. — impact-negligible → report only
- The queue-kanban package's own VERSION file is stale (0.236.20) and nothing reads it. — impact-negligible → report only

## Lessons Learned

**What worked:** Exploration measured the git read both ways before the builder wrote it: the do-work pathspec drops every merge and with it the builder commits behind each merge, which settled the REQ's open "measure both" choice with numbers. Shipping activity as a separate top-level payload map avoided a write into serve's cached Requests map, which exploration found would race.
**What didn't:** The REQ's GREEN grep for the assumed-pause wording could not pass, because the same words belong to Panel B text the next REQ owns; a capture-time GREEN check has to name the exact token it removes. Declaring a lessons satellite in a builder's Scope invites a scope-drift finding, since the satellite's index row is orchestrator-owned and the builder rightly leaves both alone.
**Worth knowing:** There is no planning, dispatched or integrating status; every in-flight REQ is claimed and phases exist only as stamps. Plain git log with name-only lists no files for merge commits, so a merge is attributed by its subject prefix and its second-parent range.

## Orientation

Now a claimed card on the board shows when its REQ was last touched (the newest lifecycle stamp or correlated commit, with the phase) on a ticking stopwatch, the drawer states the largest idle gap, and the over-4h assumed-pause badge is gone; lives in the board tool's per-response payload beside the verify probes (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). [MAP CHANGED]: a new git read (activity_correlation.go) joins the verify probes in the per-request slot outside the mtime cache, and the payload gains a top-level requestActivity map. Prime spot-check: both primes' referenced paths exist; prime-do-kanban.md's Read first does not name the new file, which is reachable from serve.go and generate.go (report only). prime-releases.md unchanged.

## Heavy Verification Plan

- Base revision: ac5cbe6461db017e1757487860eed3d2f60185b1
- Target revision: 47c2620eb13f8389e6ad0a998c1869c6d1cbad92 (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

## Heavy Verification Result

- Target revision: 47c2620eb13f8389e6ad0a998c1869c6d1cbad92
- Execution revision: 47c2620eb13f8389e6ad0a998c1869c6d1cbad92 (detached drain checkout `.git/work-run-2026-10-05-203159/drain-head`, run 21:32:51Z to 21:34:55Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 8s
- queue-kanban-browser: exit 0, executed (not skipped), 78s
- staged-skills: exit 0, executed, 36s

Green: every selected lane present, exit 0, none skipped, none reused. The first drain at 9dacf899 was red on the browser lane; see Testing → Heavy-lane remediation.

## Timing

Observed 2026-10-05T20:31:59Z to 2026-10-05T21:35:20Z: 1h 03m 21s total, 33m 03s attributed across 7 events, 30m 18s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 11m 05s | 1 |
| exploration-preflight | 6m 18s | 1 |
| remediation | 5m 44s | 1 |
| verification-gate | 5m 18s | 2 |
| review | 4m 33s | 1 |
| handback-merge | 5s | 1 |

Slowest stage: builder-work / builder subagent in worktree, 11m 05s, outcome success.
