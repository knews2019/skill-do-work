## Review: REQ-646

**Approve** — the ancestry loop now expands only from direct matches copied before it runs, so a builder's inner merge of main no longer gives main's commits to the REQ. The fix is correct and RED was reproduced independently.
Route A | 84f91190..e3ef0109 (merge e3ef0109, builder commit c9c829c3)

### What's built
- `correlateCommitsToRequests` (`skills/do-work-board/tools/queue-kanban/activity_correlation.go:139-163`) first copies the direct REQ ids of every two-parent commit into `directIdsByMergeHash`. The ancestry loop reads only that copy, so range attribution can no longer feed back into what counts as a direct match.
- New lock-in test `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge` uses the REQ's Red-Green Proof fixture unchanged.
- Still missing from this range (expected): the board CHANGELOG release and the lessons bullet. Both are finalization work (requirements 5 and 6).

### Decisions / risks for you
- None.

### Findings

**Important:**
- None.

**Minor:**
- F1: a merge of main into a builder branch whose own subject carries a bracketed `[REQ-NNN]` token is still a direct match, so it still expands main's commits into that REQ (`activity_correlation.go:146-163`). This behavior existed before the change and follows the stated rule ("direct matches expand"). Git's default merge subject has no brackets, so it only fires when a person or agent writes the token by hand. — impact-negligible → report only

### Requirements Checklist
- [x] R1 Snapshot direct matches before the loop. Delivered: `directIdsByMergeHash` is built in a separate pass, and `attribute` calls in the loop write only to `requestIdsByHash`, which the loop no longer reads for direct ids.
- [x] R2 Lock-in test, written first and failing. Delivered and reproduced (see Acceptance).
- [x] R3 Comment is true. The comment at lines 139-142 says the ids are snapshotted because the loop mutates `requestIdsByHash`, and the code does exactly that. The function doc comment (lines 109-113, "a two-parent commit already matched by the first two rules") is now also true.
- [x] R4 Full queue-kanban suite. The integrator's gate ran it (423 tests, exit 0). This review ran the `TestRequestActivity` subset again: all 6 pass.
- [ ] R5 Release. Not in the merge range. Finalization does it.
- [ ] R6 Lessons bullet and index token refresh. Not in the merge range. Finalization does it. The builder's proposed bullet is accurate.
- [x] Constraint: no git command, log window, or log format changed. The diff touches only the in-memory loop in `correlateCommitsToRequests` and one test.
- [x] Constraint: the existing single-merge test is untouched. The test file diff has zero removed lines, and `TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent` passes.

### Correctness of the snapshot (check 5)
- Direct ids from both the subject token and a REQ path: both go into the same per-hash set before the snapshot, and the snapshot copies the whole set. Correct.
- A matched merge inside another matched merge's range (both direct): each merge expands only its own direct ids. The outer merge no longer passes its id through the inner merge. Nothing correct is lost by this. The outer range is ancestry(outer^2) minus ancestry(outer^1), which already contains every commit of the inner merge's second-parent side that is not reachable from outer^1. The only commits the old behavior added were ones reachable from outer^1, which means main's commits. So the snapshot is strictly narrower in exactly the wrong cases.
- Octopus and one-parent commits: skipped when the snapshot is built (`len(parentHashes) != 2`), the same as before. Root commits with no parents are skipped too. Only two-parent hashes reach the loop, so indexing `parentHashes[0]` and `parentHashes[1]` stays safe without the old in-loop length check.
- Duplicate hash records (if log output were ever joined from several logs) would add duplicate ids to a slice. That only repeats `attribute` calls, which write to a set, so it is harmless.

### Acceptance Testing

**Result: Pass**
- GREEN at e3ef0109 (HEAD equals e3ef0109): `go test . -run 'TestRequestActivity' -count=1 -v` in `skills/do-work-board/tools/queue-kanban` gives 6/6 PASS, including the new test and the untouched single-merge test.
- RED reproduced independently. I created a throwaway detached worktree at e3ef0109 under the scratchpad and replaced only `activity_correlation.go` with its 84f91190 version. Result: `--- FAIL: TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge` with "REQ-701 wrongly attributed main-side commit main1" and "attributed 5 instants, want 4". The worktree was restored and removed with `git worktree remove` (no --force).
- Fixture check: main1, main2, and base0 carry no id. merge1, build2, inner1, and build1 carry REQ-701. The test also asserts exactly 4 instants and that REQ-701 is the only id in the result.

### Suggested Additional Testing
- None.

### Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All requirements in the range are met. Release and lesson belong to finalization |
| Code Quality | 95% | Small, clear two-pass fix. Comment matches the code |
| Test Adequacy | 95% | The fixture is the captured RED case. RED and GREEN were both reproduced |
| Scope | 100% | Only the two write_set Go files changed |
| Risk | Low | Display-only consumers (`activityGap`, `LastActivityAt`) |
| Acceptance | Pass | Focused tests green. RED confirmed on the pre-fix source |

### Follow-ups created
- None (1 finding report only)

## Review

**Overall: 96%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1: a builder's merge of main whose own subject carries a hand-written bracketed `[REQ-NNN]` token is still a direct match and still expands main's commits into that REQ (`activity_correlation.go:146-163`). This existed before the change and follows the stated rule. — impact-negligible → report only
**Acceptance:** Pass — 6/6 `TestRequestActivity` tests pass at e3ef0109, and the new test fails on the 84f91190 source exactly as the Red-Green Proof predicts.
**Restatement sweep:** nothing redefined (the diff makes the code match its existing comment and function doc. No prime, lesson, or doc restates a changed definition)
**Suggested testing:** 0 items
**Follow-ups created:** None (1 finding report only)

*Reviewed by review-work action*
