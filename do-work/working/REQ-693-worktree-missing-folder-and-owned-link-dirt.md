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
write_set: ["skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go", "skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go", "skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go", "skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work/actions/cleanup.md", "skills/do-work/actions/fan-out-reference.md"]
claimed_at: 2026-10-10T19:22:00Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  basis:
  - Route B
  - 8-file write set
  - 3 subsystems involved
  - 7 acceptance criteria
  calculated_at: 2026-10-10T19:25:00Z
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

Claim-time consult (2026-10-10, pre-dispatch): `_dev/primes/lessons-releases.md` (666 tokens) stays and was read. Budget left: 1334 tokens. Every other match is `slugged: partial`, so it can only be taken bare, and none fits:

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only; matches the do-work-cli cleanup, corehelpers and resultmodel paths (families `ignored-path-is-expendable`, `cross-module-fact-handoff`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 9045 tokens, bare only; matches `verify.go` (families `unknown-reads-as-clean`, `paired-predicate-drift`).
- `_dev/primes/lessons-action-files.md` — 8128 tokens, bare only; matches the one-clause prose edits to `actions/cleanup.md` and `actions/fan-out-reference.md` (downstream readers of a changed rule).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-145/REQ-660-worktree-lifecycle-command.md` → `## Review`, `do-work/runs/work-2026-10-10-131527/REQ-660-review.md` (F3, F4), `do-work/runs/work-2026-10-10-131527/REQ-660-integration-report.md`, `do-work/archive/UR-145/REQ-658-finalize-auto-manifest.md` (F7).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: UR-154 R1 — "`status` and `cleanup` exit 2 when a builder worktree folder is gone (review F4). Cleanup Pass 5, the board's verify check, the handoff survey and crash recovery read the worktree command's own links as uncommitted work and ask for consent (review F3; also REQ-658 F7)."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the reader list are clear and named, but the fix spans two Go modules (do-work-cli cleanup and corehelpers, plus the separate queue-kanban module) and needs the existing `ownedLinks`/`readLinkConfig` pattern lifted into a package-level helper. Exploration settles where the shared helper lives and what Pass 5 must do after a linked worktree reads clean.

**Planning:** Not required

No `## Open Questions` section exists, so Step 3.5 is skipped.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Read at main HEAD `fd9a2378` by the pre-dispatch agent (no separate Explore agent: the REQ already names every reader, so the exploration is a direct read of each).

**Missing folder (F4).**
- `internal/cleanup/worktree_lifecycle.go:241-265` `reportStatus` runs `dirtyPaths` (a `git status` inside `record.Path`) for every builder record and stops the whole command with `WORKTREE-GIT-FAILED` on the first error. `rev-list` and `log -1` run in the main root against `record.Head`, so they still work for a missing folder; only the dirt read fails.
- `worktree_lifecycle.go:344-393` `removeBuilderWorktree` finds `worktreePath` from the porcelain records, runs `dirtyPaths` on it, then `worktree remove`, `branch -d`, `worktree prune`.
- Reproduced in a scratch repo: `git worktree list --porcelain` keeps a `prunable` record (with `HEAD` and `branch`) after `rm -rf` of the folder, and `git branch -d` then refuses with "cannot delete branch ... used by worktree at ...". So the missing-folder cleanup must run `git worktree prune` **before** `git branch -d`.
- `resultmodel/result_model.go:631-638` `WorktreeStatusRow` has no missing marker; the text renderer is `result_model.go:1347-1354`.
- `corehelpers/handoff.go:87` already detects a missing folder with `os.Stat` and reports `HANDOFF-WORKTREE-MISSING`; the same check is the model for the lifecycle command.

**Owned links (F3).**
- The owned-link rule lives only as `lifecycleRun` methods: `readLinkConfig` (`worktree_lifecycle.go:435`, reads `do-work/worktree-links` from the main root, refuses absolute or `..` paths), `ownedLinks` (`:463`, a configured path whose `os.Readlink` target equals `filepath.Join(mainRoot, path)`), and `dirtyPaths` (`:477`, porcelain minus owned links). Both need `run.mainRoot`, which `handleWorktree` takes from the first porcelain `worktree` record.
- Reader 1: `internal/cleanup/cleanup_git.go:345` `worktreeClean(ctx, path)` = empty `git status --porcelain=v1 -z --untracked-files=all`. Used once, by Pass 5 `ApplyWorktreeRepairs` (`cleanup_git.go:255`), which already reads `git worktree list --porcelain -z` (first record = main root).
- Reader 2: `internal/corehelpers/handoff.go:91-109` reads porcelain per worktree (including the main tree) through its own `parseHandoffStatus` (typed rows with `origin`), and emits one `HANDOFF-WORKTREE-DIRTY` per row. Its `worktree list --porcelain` read (line 63) also has the main root as the first record.
- Reader 3: `skills/do-work-board/tools/queue-kanban/verify.go:1109` `worktreeHasUncommittedWork(worktreePath)`, called once from `classifyWorktreeLeftover` (`:1172`). The main tree path is the first `worktree ` line that `listWorktreeAgentWorktrees` (`:1449`) already reads and discards. Its comment (`:1097-1108`) defines the question as "would Pass 5's non-forced `git worktree remove` refuse".
- Import direction: `corehelpers` already depends on `cleanup` transitively (through `doctor`, `go list -deps`), so `corehelpers` may import `cleanup` with no cycle. queue-kanban is a separate module and must carry its own copy of the rule.
- Reproduced in a scratch repo: `git worktree remove` (no `--force`) refuses a worktree whose only untracked path is a symlink ("contains modified or untracked files"). So once `worktreeClean` ignores owned links, Pass 5 must also delete the owned links before its non-forced `git worktree remove`, exactly as `removeBuilderWorktree` does at `worktree_lifecycle.go:377-381`; otherwise Pass 5 turns a consent finding into `WORKTREE-REMOVE-FAILED`. This also keeps the board's "would Pass 5 refuse" comment true.
- The F3 RED only fires when the ignore line does not match the link itself (for example `node_modules/`, a directory-only pattern), which `fan-out-reference.md:53` already warns about. The existing fixture `lifecycleRepository` (`worktree_lifecycle_test.go:21`) uses exactly that shape (`shared-deps/`), so it reproduces the RED.

**Prose.** Crash Recovery (`actions/work-reference.md:308-313`) does not restate the dirt rule, so it stays unchanged (REQ allows this). `actions/cleanup.md` Pass 5 step 3 says "the worktree is clean" and then names `git worktree remove <path>`; it gains one clause. `actions/fan-out-reference.md:53` lists what `worktree status` shows; it gains "or missing".

**Tests to model.** `TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind` (`worktree_lifecycle_test.go:78`, fixture plus `runWorktreeCommand`) and `TestMergedCleanBuilderWorktreeIsAutomaticButUnmergedNeedsExactConsent` (`cleanup_git_test.go:141`, archived completed REQ plus `ApplyWorktreeRepairs`).

**Lessons consulted.** `_dev/primes/lessons-releases.md` (no release-payload change here; the builder edits no CHANGELOG or VERSION). Grepped the dropped satellites for the matching families: `paired-predicate-drift` (a second reader of one rule drifts unless both are pinned together) is why D-03 makes the board copy name its twin in a comment.

**Pre-dispatch decisions** (DECIDE & STATE; the builder's first new id is D-07):
- D-01: The shared do-work-cli helper lives in package `cleanup` as package-level functions taking the main root (read config, owned links, dirt minus owned links); the `lifecycleRun` methods become thin callers that keep today's refusal codes. `corehelpers/handoff.go` imports `cleanup` and filters its typed rows through the exported owned-link function. No second link parser in do-work-cli.
- D-02: Every reader takes the main root from the first porcelain `worktree` record it already reads, the same derivation as `handleWorktree`, because the link target `worktree new` writes is built from that path.
- D-03: queue-kanban gets its own owned-link filter in `verify.go` (same file, same symlink-target check), with a comment naming `worktree_lifecycle.go` as the twin. No new board test (REQ constraint: one test for the shared do-work-cli reader only).
- D-04: Fail toward consent outside the lifecycle command: when `do-work/worktree-links` is unreadable or invalid, Pass 5, the handoff survey and the board filter nothing and report dirt as before. `worktree new|status|merge|cleanup` keep their existing refusals.
- D-05: Missing folder means `os.Stat(record.Path)` reports not-exist. `status` keeps the row, sets a new `missing` field, skips the dirt read and still fills ahead/behind and last-commit age. `cleanup` skips the dirt check, link removal and `git worktree remove`, and runs `git worktree prune` then `git branch -d`.
- D-06: Pass 5 deletes owned links (only links, never their targets) right before its non-forced `git worktree remove`, on the non-dry-run path. The forced discard path is unchanged.

*Generated by pre-dispatch agent (direct read)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (modify) — missing-folder status row and cleanup order; lift the link config, owned-link and dirt logic into package-level helpers
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modify) — Pass 5 clean check ignores owned links and removes them before the non-forced remove
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (modify) — one missing-folder test and one owned-link Pass 5 test
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — missing field on the status row and its text rendering
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (modify) — handoff survey drops owned-link rows
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — board dirt probe ignores owned links with the same rule
- `skills/do-work/actions/cleanup.md` (modify) — one clause in Pass 5 step 3
- `skills/do-work/actions/fan-out-reference.md` (modify) — one clause in the worktree status sentence

**Files I will NOT touch:** `skills/do-work/actions/work-reference.md` (Crash Recovery does not restate the dirt rule), `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go`, `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go`, queue-kanban tests, `CHANGELOG.md`, any VERSION file or mirror, anything under `do-work/` except the builder hand-back file.

**Acceptance criteria (restated from REQ):**
- [ ] With a builder worktree folder deleted by hand, `worktree status` exits 0 and lists that row marked missing instead of failing with `WORKTREE-GIT-FAILED`.
- [ ] `worktree cleanup REQ-N` for a missing folder skips the dirt check and `git worktree remove`, still runs `git branch -d` and `git worktree prune`, and succeeds after a merge.
- [ ] Cleanup Pass 5's clean check (`worktreeClean`) returns clean for a freshly linked worktree whose only untracked paths are owned links.
- [ ] The handoff survey ignores owned links when it decides a worktree is dirty.
- [ ] queue-kanban `worktreeHasUncommittedWork` ignores owned links with the identical rule (same file, same symlink-target check).
- [ ] A reader that cannot answer still reports and never deletes (fail toward consent).
- [ ] One focused missing-folder test (status and cleanup) and one owned-link test for the shared do-work-cli reader; no test per reader.
