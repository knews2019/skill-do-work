# REQ-632 hand-back (board cards show last correlated activity; assumed-pause badge removed)

- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-632-board-cards-last-correlated-activity
- Branch: worktree-agent-REQ-632-board-cards-last-correlated-activity
- Commit: 55d8ea17 (one commit on top of 2e34a507). Not merged, not pushed. No release files touched.

## File manifest (paths under skills/do-work-board/)
- tools/queue-kanban/activity_correlation.go — new. gitCommandRunner + runGitCommand, correlateCommitsToRequests (path / prefix / merge ^1..^2 from the %P graph), liveBranchTipInstants (branch listing through the injected runner), collectRequestActivity, attachRequestActivity, payload types.
- tools/queue-kanban/activity_correlation_test.go — new. 8 tests, canned git output, no git spawned.
- tools/queue-kanban/serve.go — modified. activityGitRunner field; one responseInstant shared by attachVerifyFindings and attachRequestActivity in the per-request block.
- tools/queue-kanban/generate.go — modified. Badge field removed; RequestActivity top-level map; generate takes one snapshotInstant for both attach calls; one stale comment updated.
- tools/queue-kanban/durations.go — modified. implementationSpanPausedBadgeText removed; milestone list moved to package-level phaseMilestonesOf (shared by buildPhaseBreakdown and the collector, no second copy).
- tools/queue-kanban/durations_test.go — modified. Badge-text test deleted.
- tools/queue-kanban/generate_test.go — modified. Badge assertion removed; REQ-909 fixture (4h21m, 67-min gap) added; fixture attaches activity with a failing runner; fixture const renamed implementationSpanPausedFixtureSpan → implementationSpanOverCeilingFixtureSpan.
- tools/queue-kanban/javascript_behavior_a_test.go — modified. New drawer gap-row probe.
- tools/queue-kanban/javascript_behavior_b_test.go — modified. Done-card probe expects no badge (REQ-902, REQ-909); new claimed-card probe.
- tools/queue-kanban/web/board-cards.js — modified. Paused branch and tooltip deleted; last-activity line on claimed cards.
- tools/queue-kanban/web/board-detail.js — modified. appendLargestActivityGapRow, called from openRequestDetail after the phase rows.
- tools/queue-kanban/web/board.css — modified. Orphaned `.implementation-span` wrap rules and their comment deleted (`.status-invalid-flag` kept).
- docs/board-guide.md — modified. Badge row replaced by a `last activity …` row; one paragraph on the Largest idle gap row.
- tools/queue-kanban/lessons-do-kanban.md — NOT touched (see Lessons).

## Red-green evidence
RED was run against a compiling stub (collect returned nil, attach did nothing, the JS had no new code). GREEN is the same tests after the change.
- TestRequestActivityAttributesACommitThatTouchesTheRequestFileWithoutAPrefix — RED: "REQ-701 carries no activity; a path-only commit must attribute" → PASS
- TestRequestActivityAttributesEveryPrefixTokenInASubject — RED: "REQ-701 attributed instants = [], want exactly the prefixed commit" (and REQ-702) → PASS
- TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent — RED: "REQ-701 lacks the matched merge … attributed []", "attributed 0 instants, want 3" → PASS
- TestRequestActivityCountsALiveWorktreeAgentBranchTip — RED: "last activity = 0001-01-01 … want the live branch tip" → PASS
- TestRequestActivityLargestGapUnitesCommitsWithStamps — RED: "REQ-702 has no largest gap" → PASS
- TestGeneratedPayloadCarriesTheLargestActivityGap — RED: "REQ-909 payload carries no largestActivityGap: map[]" → PASS
- TestServeReadsRequestActivityOnEveryResponse — RED: "first response lastActivityAt = \"\"" → PASS
- TestJavaScriptBehaviorClaimedCardShowsItsLastCorrelatedActivity — RED: "claimed card rendered 0 activity lines … want one line \"last activity 2026-10-05T17:50:00Z 10m 00s · dispatch\"" and "data-instant-ms \"\"" → PASS
- TestJavaScriptBehaviorDetailStatesTheLargestIdleGap — RED: "openRequestDetail never calls appendLargestActivityGapRow(" → PASS
- TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan — RED: REQ-902 "wall time 18h 00m long span · assumed pause", REQ-909 "wall time 4h 21m long span · assumed pause" → PASS
- TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts — added after GREEN (archive UR folder, run artifacts, reservation marker excluded); passes.

## P-A-U
- [PLAN] One windowed `git log --since --format=%H%x00%cI%x00%P%x00%s --name-only` (no pathspec), parsed in Go; ancestry derived from the %P graph; branch tips via the same runner. Union with stamps per REQ, bounded to the work window, phases from the shared milestone list. Ship as top-level `requestActivity`. Attach beside verify in serve (one now) and generate (one now).
- [APPLY] Scope held to the REQ's Scope list. lessons-do-kanban.md not written.
- [UNIFY] git diff --stat: 13 files, 352 insertions, 132 deletions plus two new files. `gofmt -l .` empty, `go vet .` clean, `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...` → ok 63.9s. Each changed file was read in the diff; no debug output left behind. Badge grep: `grep -rn "over 4h · assumed pause\|implementationSpanPaused\|PausedBadge" skills/do-work-board/` returns nothing.

## Decisions
- D-01 DECIDE & STATE: the GREEN grep was narrowed to the badge. The broad `assumed pause` text still lives in Panel B and the UR summary (board-durations.js, board-user-request-summary.js, durations.go:24, durations_test.go:84, javascript_behavior_a_test.go:619, board-guide.md:41). REQ-633 owns that text.
- D-02 DECIDE & STATE: the formats are the real ones. The card reads `last activity <short instant> 10m 00s · dispatch`, not `10m ago`. The drawer reads `Largest idle gap: 1h 07m (dispatch → builder handback)`. Phase labels are buildPhaseBreakdown's labels in lowercase, so the REQ's "handback" shows as "builder handback".
- D-03 DECIDE & STATE: the union is bounded to the work window. It starts at claimed_at, because queue wait before a claim is not idle work. For a finished REQ it ends at the later of completed_at and release_at. This also stops the over-attribution in exploration finding F6 (a sibling's completion commit touches an older REQ's archive file) from creating fake gaps. Events past now plus the 2-minute skew allowance are dropped.
- D-04 DECIDE & STATE: stamps without a phase (created_at, status_changed_at, blocked_at, testing_updated_at) count as events but do not change the phase. A commit takes the phase of the newest phase stamp before it. On a tie, the first gap wins. A gap that sits inside one phase reads "dispatch → dispatch".
- D-05 ESCALATE (Value/Risk): activity is collected only for Claimed plus RecentlyDone tickets (7 days). `since` is the earlier of now−7d and the oldest claimed_at in that set. Older archived REQs get no drawer gap row. Value: the git read stays one bounded log (+~40ms). Risk: a reader may expect the row on old done REQs. Easy to undo: a stamps-only gap for every REQ costs one Go loop, but its "idle" claim would ignore commits.
- D-06 DECIDE & STATE: `implementationSpanReason` and `excludedReason` still ship "paused" for Panel B. The REQ-909 payload pin asserts this.
- D-07 DECIDE & STATE: both `.implementation-span` CSS rules were removed, not only the wrap rule. The second rule only kept the value and the marker from wrapping apart, and `.elapsed-duration` is already nowrap. The class names stay on the elements.
- D-08 DECIDE & STATE: the milestone list was moved to `phaseMilestonesOf` in durations.go instead of being copied into the collector, so there is still only one list of labels.

## Discovered Tasks
- impact-low: docs/board-guide.md calls the done-card span badge `took …`, but the card renders `wall time …`. This was already wrong before this change → report only.
- impact-low: tools/queue-kanban/VERSION (0.236.20) is stale and nothing reads it (exploration §11) → report only.

## Lessons
- Read: _dev/primes/lessons-kanban-board.md entries REQ-284 and REQ-422, and the REQ-305 and REQ-337 entries the new probes follow (call-site pin with the opening paren). lessons-do-kanban.md was read only through the prime's Traps.
- Not written: lessons-do-kanban.md. A new satellite bullet needs a do-work/lessons-index.md refresh in the same edit, and that file is outside my write boundary. Proposed bullet for the orchestrator: `[family: git-history-evidence] Plain git log --name-only lists no paths for merge commits, and a do-work/ pathspec drops merges entirely through history simplification — correlate by subject prefix plus the merge's ^1..^2 range, never by pathspec (REQ-632).`

## Render evidence and timing
- Ran `generate --repo-root <main tree> --out <scratchpad>/gen` with the worktree binary and read board-data.js (static file, no browser, so there is no location.href). 610 REQs. The `implementationSpanPausedBadgeText` key is gone. There are 8 requestActivity entries: 1 claimed + 7 recently done.
  - REQ-632 (claimed): lastActivityAt 2026-10-05T20:38:50Z, kind stamp, phase dispatch. Largest gap 7.03 min, claimed → claimed.
  - REQ-631 (done): largestActivityGap 4.85 min review → review. No lastActivityAt.
- /board-data.js through serve on loopback, 1 cold + 7 warm samples, same machine, back to back:
  - before (main-tree binary at HEAD): cold 1.56s, warm 105–129 ms
  - after: cold 1.53s, warm 148–157 ms (+~40 ms, the same order as the verify probes)
- Both servers were killed (pgrep confirms none left). Each one was also started under a 90-second perl alarm.

## Integration seams
None. Tests construct liveBoardServer only through newLiveBoardServer, which now sets activityGitRunner.

## Addendum — remediation (heavy browser lane)

- **What failed:** TestBrowserBehaviorCompletionCompanionsKeepReadableContrast (completion_contrast_browser_test.go) failed with "nonterminal control readings = 3, want pending and claimed" in all six scheme/width cases.
- **Cause (Confirmed):** the fixture tree is not a git repo, so the activity collector fell back to stamps only. The claimed card REQ-903's claimed_at therefore became its last activity, and the card got a second live `.elapsed-duration` (the last-activity stopwatch) next to its state timer.
- **Fix:** the probe now lists three nonterminal controls by card status and line lead: `pending:updated`, `claimed:claimed`, `claimed:last activity`. It fails on a missing, extra, or repeated control, and the message names all three. The existing style loop now also proves that the activity stopwatch keeps faint ink at 0.85 opacity, the same as the other live timers. It does, so the test was not weakened.
- **D-09 (scope extension, authorized by the orchestrator):** completion_contrast_browser_test.go was outside the declared Scope.
- **Other selectors checked:** I searched every `*browser*test.go` for `.elapsed-duration` and `.req-card-completed`. Only this file matches. `instant-ms` selectors appear only in javascript_behavior_b and _d, and those already passed in the JS lane.
- **Commit:** 43ea1930, on top of 55d8ea17. No rebase and no merge of main.
- **Lane result:** `QUEUE_KANBAN_BROWSER=<Chrome> GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` exited 0: wall 83s, 37 tests, slowest file timeline_browser_probe_test.go at 46.66s, no skips in the output. The focused test passed all 6 subtests in real Chrome. I created do-work/test-durations.tsv in the worktree with the header line. git ignores it and it is not committed.
