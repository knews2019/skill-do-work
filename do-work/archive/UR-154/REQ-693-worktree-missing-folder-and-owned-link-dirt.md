---
id: REQ-693
title: 'Addendum: worktree status and cleanup survive a missing folder, and dirt readers ignore owned worktree links'
status: completed
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
write_set: ["skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go", "skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go", "skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go", "skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work/actions/cleanup.md", "skills/do-work/actions/fan-out-reference.md", "skills/do-work-toolbox/actions/tutorial.md", "skills/do-work/docs/roadmap-guide.md"]
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
dispatch_at: 2026-10-10T19:29:01Z
builder_handback_at: 2026-10-10T19:35:16Z
integration_at: 2026-10-10T20:59:25Z
review_at: 2026-10-10T21:13:55Z
kb_status: pending
commit: 8211de3cafba03026b0410be88e2a7c3e8e0b7ed
heavy_verified_at: 2026-10-10T21:15:38Z
heavy_verified_revision: 8211de3cafba03026b0410be88e2a7c3e8e0b7ed
completed_at: 2026-10-10T21:16:28Z
release_at: 2026-10-10T21:16:28Z
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
- [x] **[PLAN]:** I read the brief, the REQ (D-01 to D-06), the crew rules, the do-work-cli prime and the kanban conventions, `lessons-releases.md`, and the three lesson families. Order: (1) add the data-only `Missing` field so the RED test compiles, (2) write both tests and run RED, (3) lift the link helpers and fix status and cleanup, (4) fix only `worktreeClean` and record the intermediate `WORKTREE-REMOVE-FAILED`, (5) add Pass 5 link removal, which turns GREEN, (6) update the handoff and board readers, (7) make the two prose edits, (8) verify.
- [x] **[APPLY]:** Done as planned. The two deviations are D-07 and D-08 below.
- [x] **[UNIFY]:** I read each of the eight files in `git diff fd9a2378` and checked for debug output, stray files and scope. Checks (run from the worktree root):

| Check | Exit | Wall |
|---|---|---|
| `REQ-693-probe.sh` | 0 | 2s |
| `REQ-693-preflight-probe.sh` | 0 | 6s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal skills/do-work-board/tools/queue-kanban` (prints nothing) | 0 | 0s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/cleanup/ ./internal/corehelpers/ ./internal/resultmodel/` | 0 | 0s |
| `go vet -C skills/do-work-board/tools/queue-kanban .` | 0 | 1s |
| `git diff --check` | 0 | 0s |
| Extra: `go test -count=1` on the whole `cleanup`, `corehelpers` and `resultmodel` packages | 0 | 15s |
| Extra: `go test -count=1 -run 'Worktree\|Verify'` in queue-kanban | 0 | 10s |

No wall-time budget failed, so nothing was rerun.

`git diff fd9a2378 --stat`:

```
 skills/do-work-board/tools/queue-kanban/verify.go  | 90 ++++++++++++++++++---
 skills/do-work/actions/cleanup.md                  |  2 +-
 skills/do-work/actions/fan-out-reference.md        |  2 +-
 .../do-work-cli/internal/cleanup/cleanup_git.go    | 36 +++++++--
 .../internal/cleanup/worktree_lifecycle.go         | 91 ++++++++++++++++------
 .../internal/cleanup/worktree_lifecycle_test.go    | 66 ++++++++++++++++
 .../do-work-cli/internal/corehelpers/handoff.go    | 18 ++++-
 .../internal/resultmodel/result_model.go           |  7 +-
 8 files changed, 267 insertions(+), 45 deletions(-)
```

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
- `skills/do-work-toolbox/actions/tutorial.md` (modify) — coordinator addendum: the "Something seems stuck" recipe points at do-work status (REQ-696 review F1)
- `skills/do-work/docs/roadmap-guide.md` (modify) — coordinator addendum: "stuck" routes to do-work status, broken or failed to do-work forensics (REQ-696 review F2)

**Files I will NOT touch:** `skills/do-work/actions/work-reference.md` (Crash Recovery does not restate the dirt rule), `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go`, `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go`, queue-kanban tests, `CHANGELOG.md`, any VERSION file or mirror, anything under `do-work/` except the builder hand-back file.

**Acceptance criteria (restated from REQ):**
- [ ] With a builder worktree folder deleted by hand, `worktree status` exits 0 and lists that row marked missing instead of failing with `WORKTREE-GIT-FAILED`.
- [ ] `worktree cleanup REQ-N` for a missing folder skips the dirt check and `git worktree remove`, still runs `git branch -d` and `git worktree prune`, and succeeds after a merge.
- [ ] Cleanup Pass 5's clean check (`worktreeClean`) returns clean for a freshly linked worktree whose only untracked paths are owned links.
- [ ] The handoff survey ignores owned links when it decides a worktree is dirty.
- [ ] queue-kanban `worktreeHasUncommittedWork` ignores owned links with the identical rule (same file, same symlink-target check).
- [ ] A reader that cannot answer still reports and never deletes (fail toward consent).
- [ ] One focused missing-folder test (status and cleanup) and one owned-link test for the shared do-work-cli reader; no test per reader.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/actions/fan-out-reference.md` (modified; review F1 adds one clause)
- `skills/do-work-toolbox/actions/tutorial.md` (modified, coordinator addendum)
- `skills/do-work/docs/roadmap-guide.md` (modified, coordinator addendum)

**What was done:** The owned-link rule moved from `lifecycleRun` methods to package-level helpers in `internal/cleanup` (`readWorktreeLinkConfig`, `ownedWorktreeLinks`, exported `OwnedWorktreeLinks`, `worktreeDirtyPaths`). `worktree status` marks a worktree whose folder is gone as `missing` (JSON `missing: true`) and exits 0; `worktree cleanup` skips the dirt check, link removal and `git worktree remove` for a missing folder and runs `git worktree prune` before `git branch -d`. Cleanup Pass 5 (`worktreeClean`), the handoff survey and the board's `worktreeHasUncommittedWork` ignore owned links, and Pass 5 deletes those links before its non-forced `git worktree remove`; an unreadable or invalid `do-work/worktree-links` filters nothing. `actions/cleanup.md` and `actions/fan-out-reference.md` describe both. Coordinator addendum: the tutorial recipe and the roadmap guide route "stuck" to `do-work status` and keep `do-work forensics` for broken or failed work. Builder commit `22d7773e`, addendum commit `df527501` (both on the builder branch), merged as `03f9beda`; review fix `dd0f2c31` (F1: the by-hand happy-path cleanup in `actions/fan-out-reference.md` deletes the owned links before `git worktree remove`), merged as `8211de3c`. Tests: `TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder`, `TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink`.

## Decisions

(from the builder hand-back, verbatim)

- **D-07 (DECIDE & STATE): the Pass 5 test uses REQ-141, not REQ-41.** With REQ-41, Pass 5 reports `request_state=malformed`, because `FilenameID` is zero-padded (`REQ-041`) while `requestIDFromWorktree` returns `REQ-41`. So a two-digit fixture can never reach "settled". The test still uses `lifecycleRepository`: it adds and commits `do-work/queue/REQ-141-owned-link.md`, runs `worktree new REQ-141`, then archives the REQ as `completed`. A one-line comment in the test explains why. The RED still fires for the intended reason (`dirty=true merged=true request_state=settled`).
- **D-08 (DECIDE & STATE): helper names and shapes.** The brief allowed keeping thin method callers. I kept the `readLinkConfig` method because it maps errors to the two refusal codes. I deleted the `ownedLinks` and `dirtyPaths` methods instead of keeping one-line wrappers, and the call sites use the package functions directly (delete before add). Names:
  - `readWorktreeLinkConfig(mainRoot) (linkPaths, invalidLine, readError)`: no new error type. The caller tells "unreadable" from "invalid" by which return value is set.
  - `ownedWorktreeLinks(mainRoot, worktreePath, linkPaths)`
  - `OwnedWorktreeLinks(mainRoot, worktreePath)`: the one export. It reads the config itself and returns nothing on an unreadable or invalid config.
  - `worktreeDirtyPaths(ctx, worktreePath, ownedLinks)`: the shared dirt reader, used by `worktreeClean`, `reportStatus` and `removeBuilderWorktree`.
- **D-09 (DECIDE & STATE): the status row shape is `Missing bool` (`json:"missing"`) next to `Dirty`.** A missing row has `Dirty: false` and still fills ahead/behind and the last-commit time. The text renderer prints `missing` in the clean/dirty slot.
- **D-10 (DECIDE & STATE): the board status read switched to `-z`.** The board now has to compare each path exactly against the config paths. Without `-z`, git quotes some paths (for example non-ASCII ones), and an owned link would then read as dirt. A non-empty entry shorter than 4 bytes counts as dirty, so an unreadable entry never reads as clean (`unknown-reads-as-clean`).
- **D-11 (DECIDE & STATE): the main-root derivation in Pass 5 is inlined.** It is the same one-line expression as `handleWorktree`, and I did not add a shared helper for one line.
- **D-12 (DECIDE & STATE): if removing an owned link fails in Pass 5,** it is reported as `WORKTREE-REMOVE-FAILED` (an existing code) and Pass 5 continues with the next leftover. No new finding code.

- **D-13 DECIDE & STATE (integrator, coordinator addendum):** REQ-696's review (F1, F2) found two more lines that send "stuck" questions to `do-work forensics`. Commit `df527501` on the builder branch points them at `do-work status`, matching `forensics.md:10` and `README.md:176`: `skills/do-work-toolbox/actions/tutorial.md` "Something seems stuck" now runs `do-work status`, and a new "Something seems broken or failed" recipe keeps `do-work forensics` (it was the tutorial's only forensics pointer); `skills/do-work/docs/roadmap-guide.md:48` routes broken or failed to forensics and stuck in-flight work to status. `roadmap-guide.md:5` keeps "stuck" as something forensics detects (REQ-696 D-02). Both files were added to `## Scope` and `write_set`. Value: shipped routing agrees in every package. Risk: none; prose only.

- **D-14 DECIDE & STATE (integrator, after review):** fixed review F1 on the builder branch (`dd0f2c31`, merged as `8211de3c`). The by-hand "Cleanup — happy path" in `skills/do-work/actions/fan-out-reference.md` ran `git worktree remove` straight away and read its refusal as uncommitted builder work, which is false for a linked worktree now that owned links are not dirt. One clause tells the integrator to delete the owned `do-work/worktree-links` symlinks (links only) first. Same file as the write set; no other change. Review F2, F3, F4 and the delta nit N1 stay report only. Value: the hand path and `worktree cleanup` agree. Risk: none beyond prose.

<!-- D-XX counter: last used D-14. Next decision: D-15. -->

## Discovered Tasks

(from the builder hand-back, verbatim)

- impact-user-visible: Cleanup Pass 5 can never settle a builder worktree whose name has a REQ number under three digits (for example `worktree-agent-REQ-41-*`). `worktreeRequestState` compares `requestIDFromWorktree` (`REQ-41`) with the zero-padded `FilenameID` (`REQ-041`) and reports `malformed`, so such a leftover always asks for consent (`internal/cleanup/cleanup_git.go` `worktreeRequestState`). Real REQ ids are three digits or more, so the impact is low → report only
- impact-internal: `listWorktreeAgentWorktrees` in queue-kanban still parses `git worktree list --porcelain` without `-z`, so a worktree path containing a newline would be misread. This existed before this change → report only
- Review F5, F7 (`splitZero` sort) and F8 in REQ-660's review were not touched, as the brief says → report only

## Qualification

**Gate records:** `advance --diff-range 54ecd108..03f9beda` returned `qualify` and `scope-drift` both `satisfied` (provenance `merged_range`), no findings. The first call raised `SCOPE-DECLARED-NOT-TOUCHED` for the backticked words `do-work status` and `do-work forensics` in the two Scope lines the integrator added for the addendum; the backticks were removed from those words only.

**Requirement trace** (against `git diff 54ecd108..03f9beda`, 10 files, +271 -46, and the merged files):
- Missing folder, status: `reportStatus` stats `record.Path`; not-exist sets `row.Missing` and skips the dirt read, so no `WORKTREE-GIT-FAILED`. `rev-list` and `log -1` still run in the main root, so ahead/behind and last-commit age fill in. `WorktreeStatusRow.Missing` (`json:"missing"`) is new and the text renderer prints `missing` in the clean/dirty slot.
- Missing folder, cleanup: `removeBuilderWorktree` sets `folderMissing` from the same stat, skips the dirt check, the link removal and `worktree remove`, and runs `worktree prune` before `branch -d` (git refuses `branch -d` while the prunable record exists).
- Pass 5: `ApplyWorktreeRepairs` reads the main root from the first `-z` porcelain record and the link config once; `worktreeClean` now calls the shared `worktreeDirtyPaths` with the owned links. The non-forced path deletes owned links (links only) before `git worktree remove`; a link-removal error is `WORKTREE-REMOVE-FAILED` and the pass continues. The forced discard path is unchanged.
- Handoff survey: `handoff.go` takes the main root from the first porcelain record and drops status rows whose path is in `cleanup.OwnedWorktreeLinks` (no import cycle; `corehelpers` already depended on `cleanup`).
- Board: `worktreeHasUncommittedWork` reads `status --porcelain -z`, keeps `--no-optional-locks`, skips owned links, handles the rename/copy second field, and treats an entry shorter than 4 bytes as dirt. `ownedWorktreeLinkPaths` is a copy of the do-work-cli rule (same file, same absolute/`..` refusal, same symlink-target check), and both copies name each other in comments.
- Fail toward consent: `OwnedWorktreeLinks` and `ownedWorktreeLinkPaths` return nothing for an unreadable or invalid config or an unknown main tree, and Pass 5 sets `linkPaths = nil` in the same case. `worktree new|status|merge|cleanup` keep `WORKTREE-LINK-CONFIG-UNREADABLE` / `-INVALID`.
- Tests: one missing-folder test (`TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder`, status and cleanup) and one owned-link test for the shared reader through Pass 5 (`TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink`). No per-reader tests.
- Prose: `actions/cleanup.md` Pass 5 step 3 and `actions/fan-out-reference.md` describe owned links and the missing state. Crash Recovery in `work-reference.md` does not restate the dirt rule and is unchanged (REQ allows this).
- Addendum (D-13): `tutorial.md` and `roadmap-guide.md` route "stuck" to `do-work status`.

**After review:** cumulative range `54ecd108..8211de3c`, 10 files (the review fix changes one sentence in `actions/fan-out-reference.md`, already in Scope). `advance` was then in the review phase and refused gate input, so the delta was checked by hand: prose only, `git diff --check` clean.

**Hygiene:** no debug output in the added lines, `gofmt -l` on the four touched Go directories prints nothing, `git diff --check 54ecd108..03f9beda` exit 0.

**Scope:** Route B. `## Scope` declares ten files (eight from pre-dispatch, two added for the coordinator addendum); the merge touched exactly those ten. `cleanup_git_test.go`, `handoff_test.go`, queue-kanban tests, changelogs, VERSION files and `do-work/` are untouched by the builder branch.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at the merge `03f9beda` (1-minute load 2.66 at start, wall 122 s, exit 0), recorded through `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-192201/REQ-693-probe.sh`: `green-gate`, `run-blocked-check` and `scope-drift` satisfied, GREEN probe exit 0 (both new tests print `--- PASS`).
Gate again at the review-fix merge `8211de3c` (1-minute load 1.75 at start, wall 117 s, exit 0), recorded with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` (recorded revision `8211de3cafba03026b0410be88e2a7c3e8e0b7ed`) because `advance` was in the review phase; GREEN probe exit 0 at `8211de3c`.
**Result:** ✓ All passing (do-work-cli 919 Go tests, queue-kanban 423, no test file over the 30 s budget).

**Red-green validation:** (traces to `## Red-Green Proof`, from the builder hand-back)
- `TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder`: ✗ before (`status` exit 2, rows nil, `WORKTREE-GIT-FAILED` "cannot change to ... No such file or directory") → ✓ after (0.50 s).
- `TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink`: ✗ before (`WORKTREE-REQUIRES-CONSENT`, `dirty=true merged=true request_state=settled`) → ✗ with only `worktreeClean` fixed (`WORKTREE-REMOVE-FAILED`: git refuses the untracked symlink, which proves D-06) → ✓ after (0.34 s).
- The handoff survey and board readers have no test (REQ constraint); the builder checked both end to end in a scratch fixture: valid config gives `HANDOFF-WORKTREE-CLEAN` and `merged-worktree-leftover [fixable]`, an invalid config gives `HANDOFF-WORKTREE-DIRTY` and `worktree-present-uncommitted-work`.

**New tests added:**
- `TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder` and `TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink` in `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go`.

**Heavy verification plan:**
- Range: 54ecd1082445c1d7277266537d14fbe5cf718277..8211de3cafba03026b0410be88e2a7c3e8e0b7ed (same six lanes as the first plan at 03f9beda)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills/do-work-board/tools/queue-kanban
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — 10 changed files matched, first: skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-10T21:13:55Z (first review 96% at 03f9beda; re-review of the fix delta 03f9beda..8211de3c at 96%)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. A builder worktree whose folder was deleted no longer breaks `worktree status` or `worktree cleanup`, and all three dirt readers ignore the owned `do-work/worktree-links` symlinks while failing toward consent on a bad config. Full report: `do-work/runs/work-2026-10-10-192201/REQ-693-review.md`.

**Important findings:** None

**Minor findings:**
- F1: `actions/fan-out-reference.md` by-hand happy-path cleanup ran `git worktree remove` without deleting the owned links first, so git refused and the refusal read as uncommitted builder work. Fixed in `dd0f2c31` (D-14); the reviewer confirmed in the delta re-review. Resolved.
- F2 (nit): redundant nil guards at `cleanup_git.go` (`linkPaths = nil`) and `OwnedWorktreeLinks` (early return); they document fail toward consent. — impact-negligible → report only
- F3 (nit): the do-work-cli dirt parse goes through the sorting `splitZero` (REQ-660 F7) while the board keeps git order; the clean/dirty answer agrees (tested with a staged rename plus an owned link), only the listed dirty paths can differ. — impact-negligible → report only
- F4 (nit, the builder's discovered task): Pass 5 cannot settle a `worktree-agent-REQ-41-*` worktree. do-work always pads ids (`REQ-%03d`), so only a non-canonical id reaches this path, and consent is the safe answer there. — impact-negligible → report only
- N1 (nit, delta re-review): the by-hand order deletes the links before any refusal, while `worktree cleanup` checks first and keeps them; harmless at Step 9. — impact-negligible → report only

**Acceptance:** Pass. Package tests green at `03f9beda`; 4 of 4 mutations caught (no Pass 5 link removal, `branch -d` before prune, no owned links passed to `worktreeClean`, dirt read in a missing folder); scratch tests for invalid and unreadable config, foreign link, dry run, rename and missing-plus-unmerged; live handoff and board fixture with a non-ASCII link name.
**Restatement sweep:** wave-end sweep for UR-154. Redefined here: the worktree dirt rule, the `missing` worktree status and Pass 5 link removal. Stale: `fan-out-reference.md:101` (F1, fixed). Consistent: `cleanup.md:140`, `fan-out-reference.md:53,105`, `cleanup-guide.md:29`, `restart-with-parallel-handoff.md:40`, `prime-do-kanban.md:34`, `work-reference.md:308-313`. Inherited from the siblings, all consistent: REQ-694 `SECTION-WRITE-FAILED` / `FRONTMATTER-FIELD-OWNED` (`work.md`, `request_writers.go` usage strings); REQ-695 C3 and the `run:` line (`status.md`, `status-guide.md`, `fan-out-reference.md`); REQ-696 fence rule (`clarify.md:106`, `capture-reference.md:196`) and "stuck" routing (`SKILL.md:40`, `forensics.md:10`, `README.md:176`, `tutorial.md`, `roadmap-guide.md`, `roadmap.md`, `stray-check.md`, `stray-check-guide.md`). REQ-696 review F1 and F2 are closed by D-13.
**Anti-bloat:** 0 unnamed additions. The four package-level helpers are D-01's named lift (two replace deleted methods), the board copy is D-03, `Missing` is D-05, the `mainWorktreePath` plumbing is D-02; 2 tests, both named by the REQ; no new options, files or finding codes.
**Suggested testing:** 2 items (live `worktree status` in a consumer repo with a link file and one builder folder deleted by hand; a real Pass 5 on a merged, settled builder worktree with a linked `node_modules`)
**Follow-ups created:** None (remaining findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Lifting the owned-link rule into package-level helpers let Pass 5, the lifecycle command and the handoff survey share one reader; the board's copy names its twin in a comment. Reproducing git's refusal in a scratch repo before dispatch (D-06) put the link removal into the plan instead of into a review finding.
**What didn't:** Filtering owned links in the clean check alone turned a consent finding into `WORKTREE-REMOVE-FAILED`, and the same gap sat in the by-hand cleanup prose until review (F1): every place that runs a non-forced `git worktree remove` needed the link removal, not only the code path.
**Worth knowing:** git refuses a non-forced `worktree remove` when the only untracked path is a symlink, and after `rm -rf` of a worktree folder `git branch -d` refuses while the prunable record exists, so prune first.

## Orientation

Now `do-work-cli` `internal/cleanup` owns the worktree owned-link and dirt rule (`readWorktreeLinkConfig`, `OwnedWorktreeLinks`, `worktreeDirtyPaths`), used by `worktree status|cleanup`, cleanup Pass 5 and the handoff survey; queue-kanban `verify.go` keeps a twin (`ownedWorktreeLinkPaths`). See `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `_dev/primes/prime-kanban-board.md`; primes checked, all four `prime_files` paths still exist.

## Heavy Verification Plan

- Base: 54ecd1082445c1d7277266537d14fbe5cf718277
- Target: 8211de3cafba03026b0410be88e2a7c3e8e0b7ed
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills/do-work-board/tools/queue-kanban
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — 10 changed files matched, first: skills/do-work-board/tools/queue-kanban/verify.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — 5 changed files matched, first: skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target revision: 8211de3cafba03026b0410be88e2a7c3e8e0b7ed
- Execution revision: 8211de3cafba03026b0410be88e2a7c3e8e0b7ed (detached drain checkout, `QUEUE_KANBAN_BROWSER` set, 1-minute load 2.92 at start). An earlier drain at the first merge `03f9beda` executed all six lanes green (8 s, 76 s, 63 s, 35 s, 66 s, 28 s).
- queue-kanban-javascript: exit 0, reused (fingerprint unchanged since the `03f9beda` run)
- queue-kanban-browser: exit 0, executed, 73 s
- do-work-cli-integrations: exit 0, reused (fingerprint unchanged since the `03f9beda` run)
- staged-skills: exit 0, executed, 35 s
- updater: exit 0, reused (fingerprint unchanged since the `03f9beda` run)
- installer: exit 0, reused (fingerprint unchanged since the `03f9beda` run)

## Timing

Observed 2026-10-10T20:59:14Z to 2026-10-10T21:15:26Z: 16m 12s total, 22m 05s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 11m 17s | 4 |
| review | 10m 35s | 1 |
| handback-merge | 13s | 1 |

Slowest stage: review / review, fix on builder branch, delta re-review, 10m 35s, outcome success.

Notes: no builder-work event, because the hand-back had already landed when the integrator started (dispatch 2026-10-10T19:29:01Z, builder commit 2026-10-10T19:35:16Z). The review event has no end flag, so it spans the review, the fix on the builder branch, the second gate and drain, and the delta re-review; the first heavy drain ran while the reviewer worked. Both overlaps are why attributed time exceeds the total.
