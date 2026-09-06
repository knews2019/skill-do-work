---
id: REQ-604
status: completed
domain: general
created_at: 2026-09-06T08:19:05Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
route: B
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T16:25:54Z
  basis:
    - Route B
    - 3-file write set
    - 4 acceptance criteria
depends_on: [REQ-601]
related: [REQ-597]
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go, skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Make atomic-download refuse an occupied target the same way in dry-run and live, and check its stat'
claimed_at: 2026-09-06T13:15:49Z
commit: cb01b1d0f40027303197a50e6a28c279172019bc
completed_at: 2026-09-06T13:27:05Z
release_at: 2026-09-06T13:27:05Z
---

# Make Atomic-Download Refuse an Occupied Target the Same Way in Dry-Run and Live, and Check Its Stat

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Align occupancy check across dry-run and live modes so live atomic-download refuses an occupied regular file with exit 2 and DOWNLOAD-TARGET-OCCUPIED finding rather than overwriting. Check error on post-rename stat to prevent nil pointer dereference, returning DOWNLOAD-FAILED with exit 2 on stat error. Add TDD red tests for occupancy rule parity and stat error handling. Re-derive atomic-download description in prescribed-shell-primitives.md.)
- [x] **[APPLY]:** (Agent: Implemented unified occupancy check and checked stat error in `commands.go`, added unit tests `TestAtomicDownloadOccupancyRule` and `TestAtomicDownloadStatFailureDoesNotPanic` in `commands_test.go`, and updated `prescribed-shell-primitives.md`.)
- [x] **[UNIFY]:** (Agent: Ran `git diff --stat`, `gofmt -w`, native project tests, and all maintainer guard checks (`audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh`, `gate.sh`).)

## Triage

**Route: B** — Mechanical core helper behavioral alignment with TDD and documentation update.

**Reasoning:**
- Modifies Go implementation (`commands.go`) to unify occupancy checking between dry-run and live runs, and adds error handling for post-rename stat.
- Requires TDD unit tests in `commands_test.go`.
- Modifies documentation in `prescribed-shell-primitives.md` to reflect measured unified behavior.

## Plan

1. **TDD RED Tests**: Add `TestAtomicDownloadOccupancyRule` asserting dry-run and live modes both return exit 2 and `DOWNLOAD-TARGET-OCCUPIED` ("target already exists") when target exists as a regular file, and directory refusal remains exit 1 ("target is a directory").
2. **Unified Occupancy & Checked Stat in `commands.go`**:
   - Check `os.Lstat(targetPath)` before branching on dry-run: if directory, return `OutcomeFindings` (exit 1, `target is a directory`); if regular file/exists, return `OutcomeFailure` (exit 2, `target already exists`).
   - Check error on post-rename `os.Stat(targetPath)` (via hookable `atomicDownloadStat`), returning `OutcomeFailure` (exit 2, `DOWNLOAD-FAILED`) instead of dereferencing nil.
   - Add `TestAtomicDownloadStatFailureDoesNotPanic` verifying stat error does not panic.
3. **Documentation Update**:
   - Update line 125 of `skills/do-work/docs/prescribed-shell-primitives.md` to describe the unified occupancy rule and handled stat error.
4. **Verification & Qualification**:
   - Run unit tests, repository guards, and full maintainer gate.

## Exploration

Explored `handleAtomicDownload` in `commands.go`:
- Dry-run previously checked `os.Lstat(targetPath)` and failed with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`) on any existing file.
- Live mode previously only checked if target was a directory (exit 1), allowing `os.Rename` to silently clobber existing regular files.
- Live mode dereferenced `info.Size()` on `info, _ := os.Stat(targetPath)` without checking the error, causing a nil pointer dereference if stat failed (e.g. concurrent deletion).
- Shipped callers in `skills/do-work-toolbox/actions/install.md` pre-check with `test -s` and do not depend on overwriting regular files.

## Scope

Files in scope for this change:
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`
- `skills/do-work/docs/prescribed-shell-primitives.md`

## What

Measured by REQ-597's guide builder (fixtures A2/A4 under its scratch directory, evidence in
`do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`): `atomic-download --dry-run` refuses any
existing target with exit 2 and `target already exists` (`commands.go:863-865`), while the live run
refuses only a directory (exit 1, `:869`) and then `os.Rename` silently replaces an occupying regular
file, reporting it as `created` with exit 0. A dry run that says "would refuse" followed by a live run
that overwrites is the opposite of what a dry run is for.

Beside it, `commands.go:891` reads `info, _ := os.Stat(targetPath)` and then `info.Size()`: a stat
failure after a successful rename, which a racing delete can produce, dereferences nil.

## Why

The guide's "Atomic download" section describes one occupancy rule; the command has two. The prose was
left alone by REQ-597 because the code, not the sentence, is what is wrong.

## Detailed Requirements

- One occupancy rule for both modes. Say in the record which one: the safe reading is that the live run
  refuses an existing regular file the way its dry run already does, so the two agree and nothing is
  silently replaced. If a caller depends on replacement, find it before choosing (grep every shipped
  prescribed block that runs `atomic-download`).
- The stat after the rename either handles its error or is removed with the size it reported; no
  discarded error before a dereference.
- Tests: dry-run and live against an occupying regular file give the same exit and the same finding;
  live against a directory unchanged; the stat-failure branch, if kept, does not panic.
- The guide's "Atomic download" sentence is re-derived from the new behaviour, measured.

## Constraints

- Shipped Go: a release. Keep the change to the occupancy check and the stat.

## Open Questions

None.

## Implementation Summary

- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**Unified atomic-download occupancy check across dry-run and live, handled post-rename stat errors, and updated canonical guide:**
- `commands.go`: Positioned `os.Lstat(targetPath)` check before dry-run evaluation so both dry-run and live modes refuse an existing directory with exit 1 (`target is a directory`) and refuse an existing regular file with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`, `target already exists`). Handled error from `atomicDownloadStat(targetPath)`, returning `DOWNLOAD-FAILED` (exit 2) if stat fails rather than dereferencing nil.
- `commands_test.go`: Added `TestAtomicDownloadOccupancyRule` verifying identical exit code 2 and findings between dry-run and live runs against an existing regular file, while retaining exit 1 for directories. Added `TestAtomicDownloadStatFailureDoesNotPanic` verifying stat failure safely returns `DOWNLOAD-FAILED` without panicking.
- `prescribed-shell-primitives.md`: Re-derived atomic download publication description to reflect the unified occupancy refusal before fetching and the error-checked stat.

## Decisions

- **D1 Refuse existing regular files in live mode (safe reading):** Aligned live mode to refuse existing regular files the way dry-run already did, preventing silent file replacement.
- **D2 Error-checked stat with package-level hook:** Retained reported published byte count while safely handling any stat error, using package-level hook `atomicDownloadStat` for deterministic test verification.

## Qualification

**Passed.** Read from range `101dc26609762417f2ae99dc1077228dbe8ed688..cb01b1d0f40027303197a50e6a28c279172019bc`, 3 files, 96 insertions, 8 deletions.
Canonical `qualify` and `scope-drift` satisfied.

- `skills/do-work/tools/do-work-cli` unit tests passed (including `TestAtomicDownloadOccupancyRule` and `TestAtomicDownloadStatFailureDoesNotPanic`).
- Repository guards `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, and `quiet-grep-pipeline-audit.sh` all passed with exit 0.
- Full maintainer gate `gate.sh` passed with exit 0 (97s wall time, 801 tests).

## Testing

**Commands executed:**
- `go test -C skills/do-work/tools/do-work-cli -v -run TestAtomicDownload ./internal/corehelpers/...` — passed, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/audit-lockins.sh` — passed, exit 0.
- `bash _dev/tests/action-shell-blocks.sh` — passed, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — passed, exit 0.
- `DO_WORK_GATE_ROOT="$(pwd)" bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh` — `Maintainer verification passed.`, exit 0 (801 tests).

## Review

**Overall: 99%** | 2026-09-06T16:26:00Z | Synthesis of review lenses (code correctness, test coverage, documentation fidelity, maintainer gate)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Pass.** Dry-run and live modes now share the exact same occupancy rule, silent overwrites of existing regular files are prevented, post-rename stat errors are safely handled without panics, and documentation in `prescribed-shell-primitives.md` accurately describes the unified behavior.

## Remediation

None needed.

## Lessons Learned

**What worked:**
- Unifying the occupancy check ahead of the dry-run branch guarantees strict symmetry between preview and execution.
- Testing the post-rename stat failure via a hookable variable allows testing edge-case races cleanly without unstable filesystem manipulation.

**What didn't:**
- Discarding errors on filesystem stat calls under the assumption that a preceding rename guarantees presence introduces nil pointer dereference risks under concurrency.

**Worth knowing:**
- `atomic-download` is designed for atomic initial publication. Callers needing to overwrite must explicitly manage the destination lifecycle rather than relying on silent clobbering.

## Orientation

Unifies `atomic-download` occupancy checking in `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` so both dry-run and live modes refuse existing regular files with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`) and directories with exit 1, and handles post-rename `os.Stat` errors safely. Updates `prescribed-shell-primitives.md` accordingly.
