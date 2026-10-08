# Builder brief — REQ-646 (Board activity snapshots direct merge matches before expanding ancestry, so an inner merge of main attributes nothing)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-646-board-activity-snapshots-direct-merge-matches
- Branch: worktree-agent-REQ-646-board-activity-snapshots-direct-merge-matches, created from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-646-board-activity-snapshots-direct-merge-matches.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-6, Constraints, Builder Guidance, Red-Green Proof. Requirements 5 and 6's release and index refresh belong to the integrator; requirement 6's lesson bullet is yours to PROPOSE in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-141/input.md.
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-08-141302/REQ-646-handback.md
- Route A, tdd: true, impact-user-visible, effort-mechanical, domain backend. Go only, no UI change.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, backend.md (if present), testing.md (tdd: true) — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (read both before editing; versioning is folded into the skill and the release step owns VERSION/CHANGELOG, not you). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md read every bullet carrying `[family: git-history-evidence]` (0.305.67 and 0.305.71 are this code's own history), plus the prime's Traps.

## TDD order (mandatory; the integrator verifies red-then-green evidence)
1. In `activity_correlation_test.go`, next to `TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent` (line ~145), add a test with the REQ's Red-Green Proof fixture using `cannedLogRecord` (newest first): `merge1` `[REQ-701] merge builder branch` parents `main2 build2`; `main2` parent `main1`; `build2` parent `inner1`; `inner1` (a merge with no REQ token and no REQ path) parents `build1 main1`; `main1` parent `base0` touching an unrelated path; `build1` parent `base0`; `base0`. Give the commits distinct instants. Assert through `correlateCommitsToRequests` (or the served activity, whichever the sibling test uses) that REQ-701's instants are exactly those of `merge1`, `build2`, `inner1`, `build1`, and that `main1`, `main2`, `base0` carry no REQ id. Run it, confirm it FAILS (today `main1` is attributed through `inner1`), record the exact failure text.
2. Implement; run again (GREEN). Record test name, failure-before, pass-after in the hand-back.

## The change (decided at capture, do not reopen)
- `activity_correlation.go` `correlateCommitsToRequests` (lines ~139-153): before the ancestry loop, snapshot the direct matches (the merges that have direct ids, with a copy of each id set; or a copy of the whole hash-to-ids map). The loop expands ancestry from that snapshot only, so `attribute` calls inside it never change what a later iteration reads as direct matches. Keep the comment at lines 139-141 true or rewrite it to say what the code now does.
- Do not change which commits a single matched merge attributes: the existing merge test stays green and untouched. No new git command; no change to the log window or format.
- Do NOT touch CHANGELOG.md, VERSION, any mirror, `web/`, `do-work/` (except the hand-back file), `_dev/`, or `lessons-do-kanban.md` (propose the bullet; the integrator writes it).

## Write boundary
Exactly `skills/do-work-board/tools/queue-kanban/activity_correlation.go` and `activity_correlation_test.go`. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `cd skills/do-work-board/tools/queue-kanban && go build ./... && go vet ./... && go test ./... -count=1` (whole package; ~60 s).
- `gofmt -l .` in that directory prints nothing.
- Focused probe body: `go test . -run 'RequestActivity|Correlat|Ancestry' -count=1` (must stay well under 30 s).
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: b38928de.
- File manifest: each file with (modified)/(new) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, gofmt/vet, and each file checked.
- Red-green evidence: test name, exact failure before, pass after.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for lessons-do-kanban.md in the file's top-entry shape (`- [family: git-history-evidence] <version>: **bold one-liner.** body`, self-contained, no relative archive link) and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words). The integrator writes both.
- Integration seams (none expected). Test wall times.
