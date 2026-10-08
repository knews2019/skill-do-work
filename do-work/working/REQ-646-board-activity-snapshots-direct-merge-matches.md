---
id: REQ-646
title: 'Board activity snapshots direct merge matches before expanding ancestry, so an inner merge of main attributes nothing'
status: claimed
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
claimed_at: 2026-10-08T14:12:17Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review comment "[P2] Freeze direct matches before expanding merge ancestry — activity_correlation.go:142-150", accepted by `do-work-toolbox validate-feedback` on 2026-10-08.*
