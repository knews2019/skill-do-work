---
id: REQ-646
title: 'Board activity snapshots direct merge matches before expanding ancestry, so an inner merge of main attributes nothing'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-08T14:12:47Z
created_at: 2026-10-08T14:08:54Z
user_request: UR-141
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-645]
batch: review-residuals-2026-10-08
write_set: [skills/do-work-board/tools/queue-kanban/activity_correlation.go, skills/do-work-board/tools/queue-kanban/activity_correlation_test.go, skills/do-work-board/CHANGELOG.md]
route: A
claimed_at: 2026-10-08T14:12:17Z
dispatch_at: 2026-10-08T14:14:45Z
builder_handback_at: 2026-10-08T14:19:19Z
integration_at: 2026-10-08T14:19:43Z
review_at: 2026-10-08T14:25:13Z
kb_status: pending
commit: e3ef010941eba0d8fd22dd9d25db9ef2b978427c
heavy_verified_at: 2026-10-08T14:26:02Z
heavy_verified_revision: e3ef010941eba0d8fd22dd9d25db9ef2b978427c
completed_at: 2026-10-08T14:26:28Z
release_at: 2026-10-08T14:26:28Z
---
# Board Activity Snapshots Direct Merge Matches Before Expanding Ancestry, So an Inner Merge of Main Attributes Nothing
## What
`correlateCommitsToRequests` (`skills/do-work-board/tools/queue-kanban/activity_correlation.go:142-153`) reads `requestIdsByHash[commit.hash]` as `directIds` inside the ancestry loop, after earlier iterations have already attributed range commits. A merge of main that sits inside a tagged hand-back merge's second-parent range therefore inherits the REQ id, and its own second-parent range (main's commits since the branch point) is attributed to that REQ. Snapshot the direct matches before the loop and expand ancestry from the snapshot only.
## Why
Review finding, P2: "If a builder merges main before its tagged hand-back merge, processing the outer merge attributes its request ID to that inner merge. This loop then reads the mutated map as `directIds` and attributes unrelated main-branch commits through the inner merge, falsely shortening the reported activity gap." Triage confirmed it against the code: the comment at lines 139-141 says ancestry runs "off the DIRECT matches only, collected before any range is attributed", and the code does the opposite. The prime's rule (`prime-do-kanban.md:30`, family `git-history-evidence`: commit evidence for a REQ must be commits that REQ's work owns) is violated. The visible damage is a wrong "Largest gap between events" in the drawer and a wrong last-activity instant on a claimed card. On this repo the trigger is rare (builders are never told to merge main); in consumer repos where builders or humans merge main before hand-back it fires routinely.
## Verified Facts (from triage)
- `activity_correlation.go:142-153`: `directIds := requestIdsByHash[commit.hash]` is read per iteration from the map that `attribute(rangeHash, requestId)` mutates.
- `git log` output is newest first, so the outer hand-back merge is processed before the inner merge it contains; the inner merge is already tagged when its turn comes.
- `activity_correlation_test.go:145` (`TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent`) covers one matched merge and one unmatched merge; no test nests a merge inside a matched range.
- `activityGap` and `LastActivityAt` are the only consumers; nothing writes state from them.
## Detailed Requirements
1. Before the ancestry loop, collect the merges with direct matches and a copy of their id sets (or copy the whole hash-to-ids map). The loop attributes range commits from that snapshot only; `attribute` calls inside the loop never change what a later iteration reads as direct matches.
2. Lock-in test, written first and confirmed failing, in `activity_correlation_test.go` next to the existing merge test: an outer `[REQ-NNN]` merge (parents: main tip, builder tip) whose second-parent side contains an inner merge (parents: earlier builder commit, an earlier main commit) with no REQ token or REQ path of its own; one or more main-side commits reachable only through the inner merge's second parent, touching no REQ file. After correlation the builder's own commits and the outer merge carry the REQ id; the inner merge may carry it (it is the builder's commit); the main-side commits carry nothing. Assert through `correlateCommitsToRequests` or the served activity, whichever the sibling test uses.
3. Keep the comment at lines 139-141 true, or rewrite it to say what the code now does.
4. Run the full queue-kanban test suite; the browser lane is not needed (no UI change).
5. Release per `_dev/primes/prime-releases.md` (board package changelog).
6. Lessons: `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` already carries the `git-history-evidence` family (0.305.67). Add one bullet to that family naming this instance (a comment stated the invariant, the loop read the mutated map) and refresh the satellite's token count in `do-work/lessons-index.md` in the same commit.
## Constraints
- No new git command and no change to the log window or format; this is a pure in-memory ordering fix.
- Do not change which commits a *single* matched merge attributes: the existing merge test stays green and untouched.
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-645's change into it.
## Dependencies
None. REQ-645 (recover --take-over keeps an indented generated-name heading) is related only by origin; disjoint modules, may run in parallel.
## Builder Guidance
Certainty is high: the defect is an ordering bug stated by the code's own comment. Latitude is only in the snapshot shape (list of merges with copied id sets, versus a copied map) and in how the test fixture spells the nested merge.
## Red-Green Proof
**RED prompt/case:** Canned log, newest first: `merge1` `[REQ-701] merge builder branch` parents `main2 build2`; `build2` parent `inner1`; `inner1` (merge, no REQ token) parents `build1 main1`; `main1` parent `base0` touching an unrelated path; `main2` parent `main1`; `build1` parent `base0`; `base0`. Run `correlateCommitsToRequests`.
**Why RED now:** `merge1`'s range tags `build2`, `inner1`, `build1`. When the loop reaches `inner1` its `directIds` already contains REQ-701, so `main1` (inner1's second-parent side minus `build1`'s ancestry) is attributed to REQ-701 and its instant shortens the gap.
**GREEN when:** REQ-701's instants are exactly `merge1`, `build2`, `inner1`, `build1`; `main1`, `main2`, and `base0` carry no REQ id. The existing single-merge test stays green.
**Validation:** User confirmed. The maintainer accepted the triage's remedy and RED/GREEN with "capture them and run them".
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7896 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: family `git-history-evidence` is the rule this REQ restores. The builder reads that one family's bullets directly (fixed-string grep on `[family: git-history-evidence]`) before writing.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, `slugged: partial`): its owning prime governs the board tool this REQ edits. No family matches a correlation ordering bug.
## Full Context
See `do-work/user-requests/UR-141/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Write the REQ's Red-Green Proof fixture as a `cannedLogRecord` test through `correlateCommitsToRequests` (same entry point as the sibling merge test), confirm RED, then snapshot direct matches of two-parent commits into a hash-to-id-slice map before the ancestry loop and read `directIds` from it. No new git command, no log format change.
- [x] **[APPLY]:** Test added and run RED; snapshot implemented in the one function; comment rewritten. Only the two write-boundary files changed.
- [x] **[UNIFY]:** `git diff --stat`: activity_correlation.go +16/-6, activity_correlation_test.go +48 (2 files, 64 insertions, 6 deletions). `gofmt -l .` empty, `go vet ./...` clean, `go build ./...` clean, `git diff --check` clean. Checked activity_correlation.go (snapshot only covers two-parent commits, loop body unchanged apart from slice iteration, comment true; function doc comment still accurate) and activity_correlation_test.go (fixture matches the REQ's Red-Green Proof, distinct instants, no debug output).
*Source: review comment "[P2] Freeze direct matches before expanding merge ancestry — activity_correlation.go:142-150", accepted by `do-work-toolbox validate-feedback` on 2026-10-08.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Ordering bug in one function (`correlateCommitsToRequests`), stated by the code's own comment, with the fixture shape for the lock-in test given in the Red-Green Proof. One Go file plus its test; no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (modified)

**What was done:** `correlateCommitsToRequests` now copies the direct REQ ids of every two-parent commit into a separate `directIdsByMergeHash` map before the ancestry loop, and the loop expands merge ranges only from that copy, so a merge that received its id from another merge's range no longer expands its own range. The comment above the loop now says why the snapshot exists, and the new test `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge` nests a merge of main inside a matched hand-back merge and asserts main's commits carry no REQ id.

## Decisions

(from the builder hand-back)

- D-01 DECIDE & STATE: snapshot shape is a map from merge hash to a copied slice of ids, limited to two-parent commits. Value: the loop cannot see range-given ids, and a slice cannot alias the mutated inner maps. Risk: none; iteration order over ids was already irrelevant (set semantics on attribute).
- D-02 DECIDE & STATE: the test also asserts that only REQ-701 appears in the result, to catch any stray attribution. Value: pins the whole output. Risk: none for this fixture.

## Discovered Tasks

(from the builder hand-back)

- None.

## Qualification

**Gate record:** `advance --diff-range 84f91190..e3ef0109` returned the `qualify` gate `satisfied` with provenance `merged_range` and no findings (no debug artifacts, no output primitives, P-A-U boxes ticked).

**Requirement trace against the merged diff (`git diff 84f91190..e3ef0109`, 2 files, +64/-6):**
1. Snapshot before the loop: `activity_correlation.go` builds `directIdsByMergeHash` from `requestIdsByHash` for two-parent commits in a first pass, and the ancestry loop reads `directIds` only from that map. `attribute` still mutates `requestIdsByHash`, which the loop no longer reads. Met.
2. Lock-in test written first: `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge` spells the Red-Green Proof fixture exactly (merge1, main2, build2, inner1, main1, build1, base0) through `correlateCommitsToRequests`, the same entry point as the sibling merge test. The integrator re-ran it against the pre-fix source in a throwaway checkout: it fails with "REQ-701 wrongly attributed main-side commit main1" and "attributed 5 instants, want 4", and passes at the merge. Met.
3. Comment kept true: the rewritten comment says the direct matches are snapshotted because the loop mutates `requestIdsByHash`. That matches the code. Met.
4. Full queue-kanban suite: covered by the repository gate below.
5. Release and 6. lesson: handled at finalization.

**Constraints:** no git command, log window or log format changed (the diff touches only the in-memory correlation function and its test). The existing single-merge test has zero removed lines in the diff and passes before and after the fix.

**Scope (Route A, judged against `write_set`):** touched files are exactly the two Go files in `write_set`. `skills/do-work-board/CHANGELOG.md` is the third `write_set` entry and is written by the finalization release, not the builder. No `do-work/` path is in the merge range.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge e3ef0109 (repository gate, includes the full queue-kanban `go vet` and uncached `go test`), then `advance` with that argv and exit status plus the focused probe `do-work/runs/work-2026-10-08-141302/REQ-646-probe.sh` (`go test . -run 'RequestActivity|Correlat|Ancestry' -count=1`).
**Result:** ✓ Gate exit 0 on the first run (wall 127 s; queue-kanban 423 tests in 41 s, do-work-cli 876 tests in 58 s). `advance` recorded `green-gate` satisfied and `run-blocked-check` satisfied; the probe exited 0.

**Red-green validation:** *(traces the REQ's `## Red-Green Proof`)*
- `activity_correlation_test.go` `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge`: ✗ before implementation ("REQ-701 wrongly attributed main-side commit main1", "attributed 5 instants, want 4") → ✓ after. The builder recorded this, and the integrator reproduced both results by running the test against the pre-fix `activity_correlation.go` and against the merge in a throwaway checkout. The fixture is the captured RED case unchanged.
- `TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent` (existing, untouched): ✓ before and after.

**New tests added:**
- `TestRequestActivityDoesNotExpandAnInnerMergeOfMainReachedThroughAMatchedMerge` in `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go`

**Heavy verification plan:**
- Range: 84f9119028e7b0596aaeeb33df08ae672edf6fd6..e3ef010941eba0d8fd22dd9d25db9ef2b978427c
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — both changed files matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — both changed files matched subtree skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed files matched subtree skills

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-08T14:25:13Z

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

## Lessons Learned

**What worked:** Writing the triage's Red-Green Proof fixture verbatim as a `cannedLogRecord` test made RED reproducible in seconds, and the integrator could re-prove it against the pre-fix file in a throwaway checkout.
**What didn't:** The old comment stated the invariant ("collected before any range is attributed") while the loop read the same map it was writing, so the comment hid the bug instead of preventing it.
**Worth knowing:** `git log` is newest first, so an outer hand-back merge is always processed before any merge inside its range. Any per-commit decision in `correlateCommitsToRequests` that reads `requestIdsByHash` after the range loop starts sees range-given ids. A builder's merge of main whose subject carries a hand-written `[REQ-NNN]` token still counts as a direct match (review F1, pre-existing).

## Orientation

Board request activity (queue-kanban `activity_correlation.go`, see `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`) no longer counts main's commits as a REQ's activity when its builder merged main before hand-back, so the drawer's largest gap and the card's last-activity time stay honest. Leaf change, no map change.

## Heavy Verification Plan

- Base revision: 84f9119028e7b0596aaeeb33df08ae672edf6fd6
- Target revision: e3ef010941eba0d8fd22dd9d25db9ef2b978427c
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/activity_correlation.go and activity_correlation_test.go matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/activity_correlation.go and activity_correlation_test.go matched subtree skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-board/tools/queue-kanban/activity_correlation.go and activity_correlation_test.go matched subtree skills

## Heavy Verification Result

- Target revision: e3ef010941eba0d8fd22dd9d25db9ef2b978427c
- Execution revision: e3ef010941eba0d8fd22dd9d25db9ef2b978427c (detached checkout of the merge, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: executed, exit 0, 8 s
- queue-kanban-browser: executed, exit 0, 79 s (no `HEAVY-RUN-LANE-SKIPPED`)
- staged-skills: executed, exit 0, 39 s

## Timing

Observed 2026-10-08T14:14:45Z to 2026-10-08T14:26:02Z: 11m 17s total, 11m 38s attributed across 5 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 4m 55s | 2 |
| builder-work | 4m 34s | 1 |
| review | 1m 50s | 1 |
| handback-merge | 19s | 1 |

Slowest stage: builder-work / builder worktree build, 4m 34s, outcome success.
Slowest command: verification-gate / repository gate and focused probe, 2m 32s, exit 0, .
