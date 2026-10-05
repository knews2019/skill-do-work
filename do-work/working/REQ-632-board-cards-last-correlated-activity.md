---
id: REQ-632
title: 'Board cards show last correlated activity and drop the assumed-pause badge'
status: claimed
created_at: 2026-10-05T20:11:59Z
user_request: UR-135
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-633"]
batch: board-activity-evidence
write_set: ["skills/do-work-board/tools/queue-kanban/activity_correlation.go", "skills/do-work-board/tools/queue-kanban/activity_correlation_test.go", "skills/do-work-board/tools/queue-kanban/serve.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/durations.go", "skills/do-work-board/tools/queue-kanban/durations_test.go", "skills/do-work-board/tools/queue-kanban/web/board-cards.js", "skills/do-work-board/tools/queue-kanban/web/board-detail.js", "skills/do-work-board/tools/queue-kanban/web/board.css", "skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go", "skills/do-work-board/docs/board-guide.md"]
claimed_at: 2026-10-05T20:31:37Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-135/input.md` for complete verbatim input.

*Source: A. An activity indicator on each REQ card — at minimum "last commit N min ago" where the commit is correlated to the REQ … On done cards it is the honest replacement for the assumed-pause guess.*
