---
id: REQ-723
title: '[impact-negligible] Archive-collision helper skips assets folders'
status: completed
created_at: 2026-11-14T22:31:56Z
user_request: UR-174
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-11-14T23:29:33Z
  basis:
    - Route A
    - 2-file write set
    - 3 acceptance criteria
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go]
claimed_at: 2026-11-14T23:28:11Z
dispatch_at: 2026-11-14T23:32:28Z
builder_handback_at: 2026-11-15T01:53:29Z
integration_at: 2026-11-15T01:53:44Z
review_at: 2026-11-15T01:58:57Z
commit: 858c070867e9f4929cf0306f88bf2c1cf9aef3f3
heavy_verified_at: 2026-11-15T02:06:47Z
heavy_verified_revision: 858c070867e9f4929cf0306f88bf2c1cf9aef3f3
kb_status: pending
completed_at: 2026-11-15T02:07:12Z
release_at: 2026-11-15T02:07:12Z
---
# Archive-Collision Helper Skips Assets Folders
## What
The standalone `archive-collision REQ-NNN` command (`handleArchiveCollision`, `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:158-193`) walks `do-work/archive` and lists every file named `REQ-NNN.md` or `REQ-NNN-*.md`, including evidence files inside `assets/` folders. Skip folders named `assets` in that walk, the same rule REQ-722 (Cleanup Pass 6 skips assets folders when scanning for blanked records) added to cleanup and `repositorymodel` already uses.
## Context
- Source: review finding F1 of REQ-722, recorded there as report only and promoted by the user with this capture. Original verdict: impact-negligible.
- Evidence: the walk at `commands.go:164-181` filters on the file name only and never returns `fs.SkipDir`. The reference rule is `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:192` (`directoryEntry.Name() == "assets"`) and the REQ-722 fix in `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:115`.
- Reach: the command is launched by `skills/do-work/tools/checks/archive-collision.sh`. No action prescribes it. The `advance` archive-collision hold reads the `repositorymodel` snapshot and is not affected. That is why impact is negligible.
- No queued REQ in any UR shares this root cause (the queue was empty at capture).
- Existing tests live in `skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go`.
## Red-Green Proof
**RED prompt/case:** A repo whose archive holds only `do-work/archive/UR-001/assets/REQ-001-evidence.md`. Run `do-work-cli archive-collision REQ-001`.
**Why RED now:** The command prints that evidence path and reports findings, as if REQ-001 had an archived twin.
**GREEN when:** A new Go test in `archive_collision_test.go` runs the command on that fixture and gets success with empty output, and a real archived record in the same UR folder (`do-work/archive/UR-001/REQ-001-done.md`) is still listed when present. The two existing archive-collision tests still pass.
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget, and its row is `slugged: partial` so no targeted entry is legal. Matched (claim-time consult, 2026-11-15) on the row clause "a record walker of its own that must copy repositorymodel's folder skips" and family `guard-on-one-entry-path`, where REQ-722 recorded the sibling fix. The pre-dispatch agent read that one bullet (line 5) by hand; the builder brief quotes it.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Copy the `cleanup_git.go:115` assets skip (comment and branch) into the `handleArchiveCollision` WalkDir callback; RED test first, both directions.
- [x] **[APPLY]:** (from the builder hand-back) One test function (with a small local `outputOf` closure so failure messages print the output text, not a pointer), one 4-line branch; RED commit `3cbdc83e`, fix `1be00b67`.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat df1bf423..HEAD` = 2 files, +28: archive_collision_test.go (+24), commands.go (+4); no debug output, gofmt/vet clean, missing-archive / not-a-directory / symlink paths untouched (skip placed after the root checks). Integrator re-check after the merge: see Testing.
*Source: do-work capture-request: the standalone archive-collision helper (internal/corehelpers/commands.go:164-176) does not skip assets/ folders, the same gap REQ-722 fixed in cleanup Pass 6. Source: REQ-722 review finding F1.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both files, the exact walk (`commands.go:164-181`) and the reference rule to copy (`repository_model.go:192`, `cleanup_git.go:115`); the change is one `fs.SkipDir` branch plus one Go test.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions
- **D-01**: Copy the skip rule into the walk (return `fs.SkipDir` for a directory named `assets`), not switch the command to the `repositorymodel` snapshot. Reasoning: the REQ asks for the skip rule; it is one branch and keeps the command's output format and its symlink and missing-archive behavior unchanged. Value: smallest diff, same rule as `cleanup_git.go:115`. Risk: low; a later folder-skip added to `repositorymodel` would again need copying here (the lesson family `guard-on-one-entry-path` already names that drift), reversible.
- **D-02**: Skip only `assets`, not the other `repositorymodel` skips (top-level `deliverables`, dot-folders). Reasoning: the walk starts at `do-work/archive`, so `deliverables` can never appear as a top section there, and the REQ scopes the change to `assets`. Value: no behavior change beyond the REQ. Risk: a dot-folder inside the archive would still be scanned, as today; reversible.
- **D-03** (from the builder hand-back) DECIDE & STATE: failure messages print `outcome` and the output text via a local closure instead of `%+v` of the whole `CommandResult` (which prints ~40 fields and a pointer address). Kept local to the test; no shared helper added.
<!-- D-XX counter: last used D-03. Next decision: D-04. -->

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go` (modified)

**What was done:** The `handleArchiveCollision` walk now returns `fs.SkipDir` for any folder named `assets`, placed after the archive-root checks, so evidence files under `do-work/archive/**/assets/` are no longer listed as archived twins. A new test `TestArchiveCollisionSkipsAssetsFolders` pins both directions: an assets-only archive gives success with empty output, and a real record beside the assets folder is still the only line printed. Hand-back merge range 2e5b6f70..858c0708, no conflicts.

## Discovered Tasks

(from the builder hand-back) None.

## Qualification

**Gate records:** `advance --diff-range 2e5b6f70..858c0708` (qualify) returned `qualify satisfied success` with no findings: no scope drift, no debug artifacts (the diff adds one comment, one `if` branch and one test function; no TODO, FIXME or print lines).

**Requirement trace (against `git diff 2e5b6f70..858c0708`, 2 files, +28/-0):**
- Skip folders named `assets` in the `handleArchiveCollision` walk: `commands.go` returns `fs.SkipDir` when `entry.IsDir() && entry.Name() == "assets"`, the same rule as `repository_model.go:192` (`directoryEntry.Name() == "assets"`, any depth) and `cleanup_git.go:115`. ✓
- RED case, an archive holding only `do-work/archive/UR-001/assets/REQ-001-evidence.md` now gives success with empty output: first half of `TestArchiveCollisionSkipsAssetsFolders`. ✓
- A real archived record in the same UR folder (`do-work/archive/UR-001/REQ-001-done.md`) is still listed: second half of the same test asserts `OutcomeFindings` and that path as the only output line. ✓
- The two existing archive-collision tests still pass unmodified: `TestArchiveCollisionFindsNestedRequestsWithoutFollowingDirectoryLinks` and `TestArchiveCollisionDistinguishesMissingArchiveFromUnreadableArchive` pass in the GREEN probe and the package run. ✓
- Missing-archive, not-a-directory and symlink behavior unchanged: the new branch sits after the archive-root checks; the diff touches no other line. ✓

**Scope comparison (Route A):** `write_set` names the 2 files; the merge touched exactly those 2. No undeclared file.

## Testing

**Tests run:**
- `gofmt -l ./internal/corehelpers` (empty), `go vet ./internal/corehelpers/` (clean), `go test -count=1 ./internal/corehelpers/` → `ok` (6.6 s), in `skills/do-work/tools/do-work-cli` on the merged main checkout.
- Repository gate `bash _dev/tests/maintainer-verify.sh` via `helpers/gate.sh` at merge `858c070867e9f4929cf0306f88bf2c1cf9aef3f3` (detached drain-head checkout, stage reuse off), started at 1-minute load 3.47 → `GATE-EXIT=0 WALL=130s`.
- GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-723-probe.sh` through `advance --probe-file` → probe status 0 (the three archive-collision tests PASS); green gate recorded for the merge revision.

**Result:** ✓ All passing

**Red-green validation:**
- `TestArchiveCollisionSkipsAssetsFolders` (`archive_collision_test.go`): ✗ at RED commit `3cbdc83e` (`evidence file in assets/ reported as archived twin: outcome=findings output="do-work/archive/UR-001/assets/REQ-001-evidence.md\n"`) → ✓ at `1be00b67` and at the merge (from the hand-back; re-run here).

**New tests added:**
- `TestArchiveCollisionSkipsAssetsFolders`: an archive holding only `UR-001/assets/REQ-001-evidence.md` → success, empty output; then add `UR-001/REQ-001-done.md` → findings, that path is the only line.

**Heavy verification plan:**
- Range: 2e5b6f70910d79eb9459734d6dd9ebf929146669..858c070867e9f4929cf0306f88bf2c1cf9aef3f3
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — both changed files matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — both changed files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — both changed files matched subtree skills/do-work/tools/do-work-cli

*Verified by work action (integrator, run work-2026-11-15-012821)*

## Review

**Overall: 98%** | 2026-11-15T01:58:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nit N1 (no score weight): `archive_collision_test.go:55-60` local `outputOf` closure differs from the siblings' `%+v` failure format, but it is test-local and prints the output text instead of a pointer, so it earns its place — impact-negligible → report only. Anti-bloat count: 0.
**Acceptance:** Pass — local `go test -count=1 -run TestArchiveCollision ./internal/corehelpers/` at 858c0708: 3/3 PASS; red-green recorded at `3cbdc83e` → `1be00b67`.
**Restatement sweep:** redefined the archive-collision command's scan set (now excludes `assets/` folders); grep of `skills/` finds no restatement of it, nothing stale.
**Suggested testing:** 0 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: 2e5b6f70910d79eb9459734d6dd9ebf929146669
- Target: 858c070867e9f4929cf0306f88bf2c1cf9aef3f3
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — both changed files matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — both changed files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — both changed files matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target: 858c070867e9f4929cf0306f88bf2c1cf9aef3f3
- Execution revision: 858c070867e9f4929cf0306f88bf2c1cf9aef3f3 (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, start load 4.87)
- do-work-cli-integrations: executed, exit 0, 71 s
- staged-skills: executed, exit 0, 39 s
- updater: executed, exit 0, 70 s
- installer: executed, exit 0, 29 s

## Lessons Learned

**What worked:** RED test first, written in both directions (assets-only → empty, real record beside assets → listed alone), so a skip that is missing or too wide both fail.
**What didn't:** Nothing.
**Worth knowing:** None beyond the hand-back. This is a second instance of the `guard-on-one-entry-path` family that REQ-722 recorded: any record walker of its own must copy `repositorymodel`'s folder skips. No new lesson bullet (the builder proposed none).

## Orientation

do-work-cli `corehelpers`: the standalone `archive-collision` walk now skips `assets/` folders, matching `repositorymodel` and cleanup Pass 6.

## Timing

Observed 2026-11-14T23:28:59Z to 2026-11-15T02:06:47Z: 2h 37m 48s total, 2h 30m 49s attributed across 6 events, 6m 59s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 2h 21m 01s | 1 |
| verification-gate | 5m 58s | 2 |
| planning | 1m 58s | 1 |
| review | 1m 32s | 1 |
| handback-merge | 20s | 1 |

Slowest stage: builder-work / builder worktree build, 2h 21m 01s, outcome success.
