# Hand-back — REQ-646 (board activity snapshots direct merge matches before expanding ancestry)

- Branch: worktree-agent-REQ-646-board-activity-snapshots-direct-merge-matches
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-646-board-activity-snapshots-direct-merge-matches
- Base: b38928de
- Commits: c9c829c3 (single commit)

## File manifest
- skills/do-work-board/tools/queue-kanban/activity_correlation.go (modified): `correlateCommitsToRequests` builds `directIdsByMergeHash` (two-parent commits with a copied id slice) before the ancestry loop; the loop expands only from that snapshot. Comment rewritten to say the snapshot exists because the loop mutates `requestIdsByHash`.
- skills/do-work-board/tools/queue-kanban/activity_correlation_test.go (modified): new `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge`, placed after the existing single-merge test (which is untouched).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Write the REQ's Red-Green Proof fixture as a `cannedLogRecord` test through `correlateCommitsToRequests` (same entry point as the sibling merge test), confirm RED, then snapshot direct matches of two-parent commits into a hash-to-id-slice map before the ancestry loop and read `directIds` from it. No new git command, no log format change.
- [x] **[APPLY]:** Test added and run RED; snapshot implemented in the one function; comment rewritten. Only the two write-boundary files changed.
- [x] **[UNIFY]:** `git diff --stat`: activity_correlation.go +16/-6, activity_correlation_test.go +48 (2 files, 64 insertions, 6 deletions). `gofmt -l .` empty, `go vet ./...` clean, `go build ./...` clean, `git diff --check` clean. Checked activity_correlation.go (snapshot only covers two-parent commits, loop body unchanged apart from slice iteration, comment true; function doc comment still accurate) and activity_correlation_test.go (fixture matches the REQ's Red-Green Proof, distinct instants, no debug output).

## Red-green evidence
- Test: `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge`
- Before (RED, `go test . -run TestRequestActivityDoesNotExpandAnInnerMergeOfMain -count=1`):
```
--- FAIL: TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge (0.00s)
    activity_correlation_test.go:222: REQ-701 wrongly attributed main-side commit main1 at 2026-10-05 14:00:00 +0000 UTC
    activity_correlation_test.go:226: REQ-701 attributed 5 instants, want 4: [2026-10-05 17:00:00 +0000 UTC 2026-10-05 15:00:00 +0000 UTC 2026-10-05 14:30:00 +0000 UTC 2026-10-05 14:00:00 +0000 UTC 2026-10-05 13:00:00 +0000 UTC]
FAIL
```
- After (GREEN): `--- PASS: TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge`; the existing `TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent` also PASS.

## Test wall times
| Command | Result | Wall |
|---|---|---|
| `go build ./... && go vet ./... && go test ./... -count=1` | ok | 65 s |
| `go test . -run 'RequestActivity\|Correlat\|Ancestry' -count=1` | ok | 3 s |

## Decisions
- D-01 DECIDE & STATE: snapshot shape is a map from merge hash to a copied slice of ids, limited to two-parent commits. Value: the loop cannot see range-given ids, and a slice cannot alias the mutated inner maps. Risk: none; iteration order over ids was already irrelevant (set semantics on attribute).
- D-02 DECIDE & STATE: the test also asserts that only REQ-701 appears in the result, to catch any stray attribution. Value: pins the whole output. Risk: none for this fixture.

## Discovered Tasks
- None.

## Lessons read
- `lessons-do-kanban.md` bullets with `[family: git-history-evidence]`: 0.305.71 (branch tip is evidence only when the branch owns it) and 0.305.67 (merge paths, pathspec drops merges, `^1..^2` range from the logged graph).
- `prime-do-kanban.md` Traps (git-history-evidence, paired-predicate-drift).

## Proposed lesson bullet (integrator writes; top-entry shape)
- [family: git-history-evidence] <version>: **expand merge ancestry only from direct matches snapshotted before the loop that attributes ranges.** REQ-632's correlation attributes a matched merge's `^1..^2` range to its REQ, and a comment said ancestry ran off direct matches only, but the loop read each merge's ids from the same map it was writing. Because `git log` is newest first, a builder's inner merge of main was tagged by the outer `[REQ-NNN]` hand-back merge before its own turn, then expanded, so main's commits counted as that REQ's activity and shortened the drawer's largest gap. REQ-646 copies the direct matches of two-parent commits into a separate map first; a test nests a merge of main inside a matched merge's range and asserts main's commits carry no REQ id.

## Proposed CHANGELOG entry (board package; integrator writes)
- Title: Board activity no longer credits main's commits to a REQ when its builder merged main before hand-back
- Body: The board links builder commits to a REQ through the hand-back merge. If the builder had merged main into its branch first, that inner merge was treated as a REQ match too, and main's commits were counted as the REQ's activity. This made the "Largest gap between events" too small and could move the card's last-activity time. Only commits matched directly by REQ path or `[REQ-NNN]` subject now expand their merge range.

## Integration seams
- None. Two files, no shared state with REQ-645.
