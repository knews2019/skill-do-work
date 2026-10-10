---
id: REQ-693
title: 'Addendum: worktree status and cleanup survive a missing folder, and dirt readers ignore owned worktree links'
status: claimed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
addendum_to: REQ-660
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-694, REQ-695, REQ-696]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go", "skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go", "skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go", "skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go", "skills/do-work-board/tools/queue-kanban/verify.go"]
claimed_at: 2026-10-10T19:22:00Z
---
# Addendum: Worktree Status and Cleanup Survive a Missing Folder, and Dirt Readers Ignore Owned Worktree Links

## What

Two gaps the REQ-660 (builder worktree lifecycle command) review left as report only. First, `do-work-cli worktree status` and `worktree cleanup REQ-N` exit 2 (`WORKTREE-GIT-FAILED`) when a `worktree-agent-*` folder is gone, so one stale record hides every row (review F4). Second, cleanup Pass 5, the board's verify check, the handoff survey and crash recovery read the symlinks that `worktree new` creates from `do-work/worktree-links` as uncommitted work and ask for consent (review F3, repeated as REQ-658 F7).

## Prior Implementation

REQ-660 shipped `do-work-cli worktree new|status|merge|cleanup` in `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go`, archived at commit `bd9ecb2d`. It reuses `parseWorktrees` and `cleanupGit` from `cleanup_git.go`. Its own `dirtyPaths` filters owned links through `ownedLinks` (a configured path that is a symlink to the same main-tree path), so `status` and `cleanup` already ignore their own links. The other readers do not.

## Detailed Requirements

- Missing folder (F4): when a builder worktree record's directory does not exist, `worktree status` still exits 0 and lists that row marked missing instead of failing. `worktree cleanup REQ-N` for that REQ treats the folder as already gone: it skips the dirt check and `git worktree remove`, and still runs `git branch -d` and `git worktree prune`. Cleanup Pass 5 already prunes first, so it can serve as the model.
- Owned links (F3): every reader that decides whether a builder worktree holds uncommitted work ignores owned links the same way `dirtyPaths` does. At HEAD the readers are `worktreeClean` (`internal/cleanup/cleanup_git.go:345`, used by cleanup Pass 5), the handoff survey (`internal/corehelpers/handoff.go`, the `status --porcelain` read after the missing-path check), and queue-kanban `worktreeHasUncommittedWork` (`skills/do-work-board/tools/queue-kanban/verify.go:1109`). The list is the current instance set, not a closed list: the condition is "a reader that reads a builder worktree's porcelain status to decide dirty or clean". Crash Recovery (`actions/work-reference.md` → Crash Recovery (Step 1)) follows the reader it uses; change its prose only if it restates the dirt rule.
- queue-kanban is a separate Go module and cannot import do-work-cli. It reads `do-work/worktree-links` itself; keep the owned-link rule identical (same file, same symlink-target check).

## Red-Green Proof

**RED prompt/case:** In a fixture repo, run `worktree new REQ-1`, delete the worktree folder by hand (`rm -rf`), then run `worktree status`. Separately, with a `do-work/worktree-links` file listing `node_modules`, run `worktree new REQ-2` and ask Pass 5's clean check (`worktreeClean`) about that worktree.
**Why RED now:** `status` exits 2 with `WORKTREE-GIT-FAILED` and no rows. `worktreeClean` returns false because the `node_modules` symlink shows as untracked.
**GREEN when:** `status` exits 0 and shows the REQ-1 row marked missing, and `cleanup REQ-1` succeeds after a merge. `worktreeClean` returns true for the freshly linked REQ-2 worktree.
**Validation:** Inferred during capture

## Constraints

- One focused test per failure: a missing-folder test in `worktree_lifecycle_test.go` (status and cleanup), and one owned-link test for the shared do-work-cli reader. Do not add a test per reader.
- Fail toward consent stays the rule: a reader that cannot answer still reports, never deletes.
- Keep the change small: reuse `ownedLinks` and `readLinkConfig` logic inside do-work-cli rather than adding a second link parser there.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high on both failures (reproduced by the reviewer, re-checked at HEAD `0e8e0ef9`). Latitude: the builder picks the row shape for "missing" and whether the do-work-cli readers share one helper.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only (`slugged: partial`); matches the do-work-cli cleanup and corehelpers paths (families `ignored-path-is-expendable`, `cross-module-fact-handoff`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 9045 tokens, bare only; matches `verify.go` (family `unknown-reads-as-clean`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-145/REQ-660-worktree-lifecycle-command.md` → `## Review`, `do-work/runs/work-2026-10-10-131527/REQ-660-review.md` (F3, F4), `do-work/runs/work-2026-10-10-131527/REQ-660-integration-report.md`, `do-work/archive/UR-145/REQ-658-finalize-auto-manifest.md` (F7).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: UR-154 R1 — "`status` and `cleanup` exit 2 when a builder worktree folder is gone (review F4). Cleanup Pass 5, the board's verify check, the handoff survey and crash recovery read the worktree command's own links as uncommitted work and ask for consent (review F3; also REQ-658 F7)."*
