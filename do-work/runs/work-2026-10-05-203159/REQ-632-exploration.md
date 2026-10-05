# REQ-632 exploration: last correlated activity on cards, drop the assumed-pause badge

All paths are relative to `skills/do-work-board/tools/queue-kanban/` unless they start with `do-work/`, `skills/` or `_dev/`.

## Findings that change what the builder writes

- **F1. No `planning`/`dispatched`/`integrating` statuses exist.** The status vocabulary is `pending, claimed, completed, completed-with-issues, failed, cancelled, pending-answers, blocked, blocked-*` (`skills/do-work/actions/work-reference.md:240`, `normalizeStatus` model.go:1014). All in-flight work is `status: claimed`. "Phases" exist only as stamps. So the card rule in Requirement 6 is effectively `request.status === "claimed"`. The phase text comes from the newest union event.
- **F2. The serve payload copy shares the cached `Requests` map.** serve.go:150 does `liveBoardData := *boardData`, which is a shallow copy. `generatedBoardData.Requests` is `map[string]generatedRequest` (generate.go:79), with value structs. If the builder writes per-request activity into `liveBoardData.Requests[id]`, that write goes into the cached map, outside `cacheMu`. That is a data race, and stale values leak into later responses. Options: (a) clone the map before mutating it, or (b) ship a separate top-level `requestActivity map[id]{...}` field set per response (simpler, and needs no clone). Field naming is left to the builder.
- **F3. Serve already uses two instants per response.** On a cache hit, `refreshBoardData` calls `currentTime()` (serve.go:351). `attachVerifyFindings` then calls it again (serve.go:163). To keep "one `now` per response", take `now := liveServer.currentTime()` once in `serveLiveBoardDataJs` and pass it to both the verify probes and the activity collector. Fully fixing the refresh path's own instant would need a signature change. Calling `currentTime()` once for the per-request block already meets the REQ for the new field.
- **F4. Do NOT add the `-- do-work/` pathspec.** Measured on this repo over 7 days: 74 commits and 8 merges without the pathspec. With the pathspec: 56 commits and **0 merges**. History simplification drops merges that are TREESAME to parent 1, and it also drops builder commits that touch no `do-work/` file. Cost is the same either way: about 30ms warm and 140ms cold (388 lines vs 280 lines).
- **F5. Merge commits list NO files under plain `--name-only`.** Checked on 25803a16, bdb5c43c and 2cbd6858. Merge subjects do carry the prefix (`[REQ-631] merge builder branch worktree-agent-REQ-631-…`), so prefix matching plus `^1..^2` ancestry covers merges. Builder commits checked carry `[REQ-NNN]` too (2596dbd2). `--diff-merges=first-parent` would list the builder's files, but they are not REQ paths, so it is not needed.
- **F6. Path attribution over-attributes.** Completion commit 2fea833a (`[REQ-631] complete: …`) also touches `do-work/archive/UR-134/REQ-630-…md`, so path matching also attributes it to REQ-630. This is acceptable as "touched", but expect it. `do-work/.req-reservations/REQ-632` does not match the required pattern, which is correct. Archive paths look like `do-work/archive/UR-NNN/REQ-NNN-*.md`, so `archive/**` must allow a UR subdirectory. `%cI` carries the local offset (`+03:00`), which `time.Parse(time.RFC3339)` accepts.
- **F7. The GREEN grep is too strong as written.** `grep -rn "assumed pause" skills/do-work-board/` also matches Panel B text, which REQ-633 owns, and the `assumed paused` substring:
  - web/board-durations.js:13, :670, :1267, :1435
  - web/board-user-request-summary.js:106
  - durations.go:25
  - durations_test.go:84
  - javascript_behavior_a_test.go:619 (asserts Panel B's sentence)
  - `docs/board-guide.md:41` (UR summary "an assumed pause")

  Requirement 7 says Panel B keeps shipping until REQ-633. Recommendation: limit the grep to the badge itself (`over 4h · assumed pause`, `implementationSpanPaused`, `PausedBadge`) and record the deviation in the REQ.
- **F8. Format mismatch between the REQ text and the existing formatter.** `makeInstantWithStopwatchNode` renders `<short instant> <duration>`, with the duration from `formatElapsedDuration`: `10m 00s`, `1h 07m` (board-core.js:127-147). It never renders `10m ago`. The drawer's `1h07m` likewise reads `1h 07m` through `formatElapsedDuration(0, ms)`, which is how `appendPhaseBreakdownRows` does it. Assert the real format (`last activity <instant> 10m 00s · <phase>`) and note the wording deviation.
- **F9. Phase names.** `buildPhaseBreakdown` labels are `Claimed, Planning, Dispatch, Builder handback, Integration, Review, Remediation, Re-review, Completed, Release` (durations.go:270-281). The REQ examples are lowercase (`dispatch`, `handback`). Reuse the labels, lowercased or as-is (builder's call). Note that `buildPhaseBreakdown` returns nil when no optional stamp parses (durations.go:284-295). For gap phases, map the stamp field to a label directly rather than calling `buildPhaseBreakdown`. `created_at` / `status_changed_at` / `blocked_at` / `testing_updated_at` have no label there.

## 1. gitCommitDateLookup

- Type: model.go:426 `type gitCommitDateLookup func(repoRoot string, commitHash string) (time.Time, bool)`. The live implementation is `lookupGitCommitDate` (model.go:1539). It checks `gitBinaryAvailable()` and runs `exec.Command("git","-C",repoRoot,"log","-1","--format=%cI",hash)`.
- Threading: `LoadBoard` (model.go:432) passes `lookupGitCommitDate` to `buildBoard(repoRoot, now, recentWindow, gitLookup)` (model.go:441). It is used only in `resolveCompletionTime` (model.go:1456). `main.go:183` calls `LoadBoard`. serve.go:363 hard-codes `lookupGitCommitDate`, and the server struct has no runner field (serve.go:36-51). The injectable-time field `currentTime func() time.Time` (serve.go:40) is the pattern for adding e.g. `activityGitRunner`.
- Test fakes: inline closures (`board_synthetic_test.go:72`, `generate_test.go:352,3662`, `model_test.go:227,501`), named `stubGitLookupNever` (filementions_test.go:21), or `nil`.
- Suggested runner shape: `type gitCommandRunner func(repoRoot string, args ...string) ([]byte, error)`. Feed the `log` output and the branch list/tips through it, so tests spawn no git.

## 2. listWorktreeAgentBranches

- verify.go:1465 `func listWorktreeAgentBranches(repoRoot string) []string`. It calls `exec.Command` directly (`git -C root branch --list worktree-agent-* --format=%(refname:short)`) and is NOT injectable. It returns nil on error.
- REQ id from a branch: `requestIdFromWorktreeName` (verify.go:1079, pattern verify.go:1075, prefix const verify.go:86).
- To keep tests git-free, run the same `branch --list` argv through the injected runner. Alternatively, accept that test temp dirs are not git repos.
- This repo currently has no `worktree-agent-*` branches.

## 3. serve / generate wiring

- serve.go:135-175 `serveLiveBoardDataJs`: `refreshBoardData()`, then the shallow copy (:150), then the per-request comment block (:156-164) with `attachVerifyFindings(&liveBoardData, currentBoard, liveServer.currentTime())`. The new call belongs here, guarded by `currentBoard != nil`. `currentBoard()` is around serve.go:325.
- The mtime cache is `refreshBoardData` (serve.go:334-379). Do not touch it.
- generate: generate.go:497 `attachVerifyFindings(&boardData, board, time.Now())` in `generateStaticSiteWithPublisher`. Add the activity call beside it, with the same `now`. `buildGeneratedBoardData` keeps its board-only signature (comment at generate.go:494-496).
- Payload struct: `generatedRequest` (generate.go:154). The span fields are around generate.go:240-265; the phase breakdown is at :265, and `generatedPhaseBreakdownEntry` is at :280.
- `implementationSpanPausedBadgeText`: field generate.go:76, set at generate.go:754, function durations.go:34-51.
- `ImplementationSpanReason` is set at generate.go:858 from `measureImplementationSpan(...).ExclusionReason`. `ExcludedReason` is set at generate.go:918, with the struct at :363-376. Keep both.
- `attachVerifyFindings` is at generate.go:674.

## 4. Stamps and durations

- `lifecycleTimestampFields` (model.go:1643): created_at, claimed_at, completed_at, planning_at, dispatch_at, builder_handback_at, integration_at, review_at, remediation_at, re_review_at, release_at, status_changed_at, blocked_at, testing_updated_at.
- `parseTimestamp` (model.go:1680) accepts RFC3339, offset-less datetime, `YYYY-MM-DD HH:MM:SS` and a bare date.
- `analysisOutlierCeiling = 4h` (durations.go:32). `dayMedianExclusionReason` is at durations.go:323. Keep both for Panel B and REQ-633.
- `measureImplementationSpan` (durations.go:222) runs from the earliest origin-eligible stamp (`earliestImplementationOrigin` :181, exclusion map :163) to completed_at.
- `buildPhaseBreakdown` (durations.go:257) uses declared pipeline order and does not sort.

## 5. Web

- board-cards.js:66-108 `makeImplementationSpanNode`: the reversed branch is :70-78 (keep it). Delete the paused branch :94-106, including the tooltip at :101.
- The card state timer is at board-cards.js:387-395. The done line is at :397-409. Put the new line right after the state-timer block, using `createElement("div","req-card-completed", "last activity ")` plus the `makeInstantWithStopwatchNode(iso)` node plus a text node `" · " + phase`.
- board-core.js:
  - `formatElapsedDuration` :127
  - `makeElapsedDurationNode` :149 (sets `dataset.instantMs` and `tickFormat="duration"`)
  - `makeInstantWithStopwatchNode` :164
  - `stateTimerSpecFor` :200
  - `refreshRelativeTimeNodes` :256 (selects `[data-instant-ms]`)
- `card.dataset.status = request.status` is at board-cards.js:125.
- board-detail.js:
  - `appendMetaRow(label, valueNode)` :321
  - `appendPhaseBreakdownRows` :461-479, called at :615 in `openRequestDetail`. Add the gap row right after it.
- board.css:1449-1462: the comment and `.implementation-span {white-space:normal}` exist for the "longer assumed-pause marker". `.implementation-span > .status-invalid-flag` is still reached by the reversed flag. Remove or reword the comment and the wrap rule. `.status-invalid-flag` (:1179) is shared, so keep it.

## 6. Badge readers and the tests that assert it

- Go: durations.go:34-51, generate.go:76, generate.go:754.
- JS: board-cards.js:94-106.
- Doc: docs/board-guide.md:53, the row `| \`over 4h · assumed pause\` | … |`. The neighbouring `took …` row at :52 is untouched.
- Tests:
  - durations_test.go:507 `TestImplementationSpanPausedBadgeTextDerivesFromTheCeiling` (delete or replace).
  - generate_test.go:3676 `TestGeneratedRequestCarriesTheDoneCardImplementationSpan`: the badge assert is at :3679-3681. The REQ-902 row at :3691 keeps `"paused"` as the reason, which is still correct.
  - javascript_behavior_b_test.go:1606 `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan`: the `boardData` stub with the badge text is at :1626. The REQ-902 expectation `"wall time 18h 00m over 4h · assumed pause"` is at :1720, and the new value is `"wall time 18h 00m"`. Check the `markerTitle` assertions below :1740 for REQ-902.
- Fixture: `buildImplementationSpanFixturePayload` (generate_test.go:3600). Its const `implementationSpanPausedFixtureSpan` (:3575) is still useful for Panel B.
- Other: user_request_progress_browser_probe_test.go and javascript_behavior_d_test.go mention implementationSpan fields, not the badge.

## 7. JS test pattern

- Helpers:
  - `generateLiveSite(t)` (generate_test.go:368)
  - `sliceBalancedBlockAfter(t, indexHtml, "function name(")` (:2009)
  - `runJavaScriptBehaviorProbe(t, name, probe)` (:283). Node runs over stdin, and the probe skips unless `QUEUE_KANBAN_JAVASCRIPT_PROBES=on`.
  - `mustMarshalJSONString` (:1606)
  - `writeVerifyFixture(t, []verifyFixtureFile{{path, content}})` (verify_test.go:24)
  - `spanFixtureFrontmatter(id,title,status,claimed,completed, extra...)` (generate_test.go:3584)
- Card probe model: javascript_behavior_b_test.go:1606-1700. It fakes the DOM with `makeNode` (with `dataset:{}`), stubs `makeInstantWithStopwatchNode`, and walks `card.childNodes` with `nodeText`. A new open-card probe should slice the REAL `makeInstantWithStopwatchNode`, `makeElapsedDurationNode`, `formatElapsedDuration` and `syncClockSkewTitle`, stub `formatShortInstant`, then assert `dataset.instantMs` on the duration node.
- Drawer probe model: javascript_behavior_a_test.go:2259 `TestJavaScriptBehaviorDetailRendersOnlyObservedPhaseBreakdown`. It slices `appendPhaseBreakdownRows` and collects rows via a fake `appendMetaRow`. Reuse this for the gap row: either slice a new `appendLargestActivityGapRow` function, or extend the test.
- Timings: with probes on, all 16 tests in javascript_behavior_b take 2.5s. The done-card test alone takes 1.6s. The default lane records b/a/d as absent or 0s, because the probes skip. The files nearest the 30s budget are strict_behavior_regression_test.go (19.1s, 2026-10-05), generate_test.go (8.0s) and verify_test.go (4.5s). activity_test.go is 0s.

## 8. Test command

- Run `cd skills/do-work-board/tools/queue-kanban && go test -count=1 ./...`. Measured at 58s wall on this machine.
- The JS probes are off by default. Turn them on with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on`, and make them strict with `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR`.
- Browser probes: `QUEUE_KANBAN_BROWSER_PROBES`, `QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR`, `QUEUE_KANBAN_BROWSER` (browser_probe_test.go:38-44).
- Focused run: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -run 'ImplementationSpan|PhaseBreakdown|Activity' .`

## 9. Git measurements (this repo, 2026-10-05)

| command | lines | warm time |
|---|---|---|
| `git log --since=7.days --format=%H%x00%cI%x00%P%x00%s --name-only` | 388 | 0.03s (0.14s cold) |
| same with `-- do-work/` | 280 | 0.03s |
| commits / merges in 7 days, no pathspec | 74 / 8 | |
| commits / merges in 7 days, with pathspec | 56 / 0 | |

- Output shape: `H\0cI\0P\0subject\n\n<file>\n<file>\n` per commit. Parents are space-separated in `%P`. A merge record has no file lines.
- Baseline `/board-data.js` on this tree, served from a fresh build: cold first request 1.63s, warm 87-110ms (6 samples).

## 10. board-guide.md

- docs/board-guide.md:53 is the badge row. Delete it.
- docs/board-guide.md:41 mentions "an assumed pause" for the UR summary's `N excluded`. That is Panel B / UR-summary semantics, so leave it for REQ-633 or reword it only if the grep is kept literal (see F7).
- There is no glossary row for an in-progress "last activity" line yet. Consider adding one to the Badges table, or to the prose that explains card time lines.

## 11. Release preimage

- `VERSION` = 0.305.66
- `skills/do-work/VERSION` = 0.305.66
- `skills/do-work/actions/version.md:5` `**Current version**: 0.305.66`
- `queue-kanban/VERSION` = 0.236.20 is stale. Nothing references it, so do not treat it as a mirror.
- CHANGELOG entry shape (CHANGELOG.md:14+): `## 0.305.67 — <Descriptive Title> (2026-10-05)`, then a 1-2 sentence plain paragraph on the problem, then `- ` bullets on what changed. A user-visible change with a removed badge is likely a patch per house practice. The finalizer decides the bump size per work-reference Step 9.
