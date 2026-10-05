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

*Reviewed by review-work action*
