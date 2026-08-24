# REQ-366 builder brief

- Base: `4717d18af7b44e47f6b08d75af1d291a1af688e0`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input`
- Branch: `worktree-agent-REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input`
- Request: `do-work/working/REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input.md`
- Declared scope: `model.go`, `model_test.go`, `actions/board.md` at the paths in the request

Read `CLAUDE.md`, the request, `_dev/primes/prime-kanban-board.md`, and relevant model/render tests.
Implement the exact `status: blocked` plus unmet-dependency bucketing rule while preserving all
named status exceptions and scheduling/probe semantics. Trace inherited counts and badges, add
mutation-sensitive model tests, update the column invariant docs/comments, and generate browser or
render evidence with the retained Chromium where appropriate. Do not edit `do-work/` or release
files in the worktree. Run focused Go tests, visual/browser evidence, the canonical maintainer gate,
and exact diff/artifact checks; commit and write the hand-back to
`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-201759/REQ-366-handback.md`
with exploration, P-A-U, decisions, exact scope, tests, mutation proof, hash, and clean status.
