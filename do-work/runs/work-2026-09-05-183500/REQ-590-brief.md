# Builder brief — REQ-590

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-590-cap-the-path-list`
- **Your branch (already checked out there):** `worktree-agent-REQ-590-cap-the-path-list`
- **Route:** A
- **Base commit:** 661765dd (main at dispatch)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work-board/tools/queue-kanban/VERSION` — release paths owned by finalization.
- Any board client file (`skills/do-work-board/tools/queue-kanban/web/*`). REQ-589 (rendering the verify findings strip as the M4 slim band) is in flight on those files in a sibling session.
- Any file outside the write set below.
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate) and do not serve the board on port 8090 — a live board owned by the user is running there.

## Write set

- `skills/do-work-board/tools/queue-kanban/verify.go`
- `skills/do-work-board/tools/queue-kanban/verify_test.go`

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (`tdd: true`)

Also read `_dev/primes/prime-kanban-board.md` and, in `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, the bullet `[family: subject-not-restated-in-detail]` — it governs how a verify finding's detail is composed and says to locate a finding by its `Subject`/`Category` field, never by a substring of its prose.

## The REQ

Read it in full (read only): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-590-cap-the-path-list-in-a-verify-finding.md`. Its `## Detailed Requirements` are the contract and its `## Red-Green Proof` is the target.

## Order of work (`tdd: true`)

1. RED first: add the test to `verify_test.go`, run it, and record the failure output.
2. GREEN: add one shared helper in `verify.go` and use it at the three joins (lines 797, 1177, 1191 in the base revision).
3. Run `go test ./...` in `skills/do-work-board/tools/queue-kanban`.

## Hand-back

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-183500/REQ-590-handback.md` (absolute path, never staged or committed) with: branch name, commit hashes, the file manifest with an action verb per file, the RED/GREEN evidence (test name, failure before, pass after), integration seams (expected: none), `## Decisions` if any, `## Discovered Tasks` if any, and the lesson entries you read.
