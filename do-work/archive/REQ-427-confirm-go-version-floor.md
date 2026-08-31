---
id: REQ-427
title: 'Confirm the Go version floor for installing and updating do-work'
status: completed
created_at: 2026-08-30T17:40:00Z
status_changed_at: 2026-08-31T13:52:30Z
claimed_at: 2026-08-31T13:55:10Z
completed_at: 2026-08-31T14:22:00Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-31T10:30:10Z
  basis:
    - trivial short-circuit
user_request: UR-081
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
addendum_to: REQ-407
review_generated: true
---

# Confirm the Go Version Floor for Installing and Updating Do-Work

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Use the exact-toolchain proof to set one truthful Go 1.25.0 core floor, synchronize current restatements, and add a dedicated compatibility lane without changing the optional board or maintainer floor.
- [x] **[APPLY]:** Updated the core module, launcher boundary, current prerequisite text, compatibility-launcher comments, and focused regression fixtures.
- [x] **[UNIFY]:** Reviewed the 18-file merge range, ran the exact Go 1.25 suite, launcher behavior and contract regressions, checked shell with ShellCheck, and found no debug artifacts.

## What
UR-081 specified Go 1.26.1+ as the prerequisite, and REQ-406 and REQ-407 implemented exactly that. Measurement during REQ-407's review found the module does not need it. This asks whether the floor should stay where you put it.

## Why This Is A Question And Not A Fix
The floor is not a builder's choice. `REQ-406` says "require Go 1.26.1+" and `REQ-407` says "document the Go 1.26.1+ prerequisite", both taken from UR-081. A reviewer raised the height of the floor as a defect and it was **refuted 3-0** on exactly that ground: implementing your stated requirement is not a bug. But the requirement now decides who can install do-work at all, so it is worth confirming rather than inheriting.

## Measurement
Taken on the merged tree, on a copy of the module with only the `go` directive in `go.mod` changed:

| `go` directive | Result with the local Go 1.24.7 toolchain |
|---|---|
| 1.26.1, 1.26.0, 1.25.0 | refuses — `go.mod requires go >= …` |
| **1.24.0** | builds |
| **1.23.0** | builds; `go vet ./...` exit 0; `go test -count=1 ./...` **all six packages pass** |

So no language or standard-library feature in `do-work-cli` requires 1.26.1. As it stands, anyone on Go 1.23, 1.24 or 1.25 cannot install or update do-work.

One thing that is genuinely 1.26: `skills/do-work-board/tools/queue-kanban/go.mod` declares `go 1.26`. The board tool is optional; the installer is not. Lowering the installer's floor would not lower the board's.

## Open Questions
- [x] Should the Go floor for installing and updating stay at 1.26.1, or drop to the lowest version the module actually builds and tests clean on? → Confirmed: lower to `1.23.0`
  **Recommended:** lower it to `1.23.0` in `skills/do-work/tools/do-work-cli/go.mod`, the launcher's `minimum_go_version`, `README.md` and `skills/do-work/actions/version.md`, so the prerequisite matches what the code needs. **Also:** keep 1.26.1 deliberately, if you want one toolchain across the whole suite and are content to exclude older installs; or pick an intermediate floor such as 1.24 to match the toolchain most CI images ship today.
  Value: a floor that matches the code stops excluding installs for no technical reason.
  Risk: low and reversible either way — it is four literals and a `go.mod` directive. Choosing to *keep* 1.26.1 costs nothing to implement, since that is what already ships.

  **Answered 2026-08-30** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: lower the installer and updater floor
  to Go 1.23.0, the lowest version on which the module was built, vetted, and tested successfully.
  The optional board tool's Go 1.26 requirement remains unchanged and is outside this REQ's scope.

- [x] Exact Go 1.23 verification disproved that earlier measurement. Which truthful compatibility target should replace it? → Confirmed: lower the core installer and updater floor to `1.25.0`
  **Recommended:** lower the core installer/updater floor to `1.25.0`, add an exact-Go-1.25 test lane, and update every current core prerequisite restatement. The complete core suite passes on the exact Go 1.25 toolchain.
  **Also:** keep `1.26.1` unchanged; or require `1.23.0` through a separate substantive backport that replaces Go 1.24+'s traversal-resistant `os.Root` filesystem boundary without weakening its protections.
  Value: `1.25.0` removes an unnecessary restriction while keeping the existing filesystem-safety design intact.
  Risk: the `1.25.0` change is small and reversible. Go 1.23 support is materially larger because it must replace rooted open/read/write/stat/walk/remove behavior across three packages and prove equivalent symlink-race containment.

  **Reopened 2026-08-31:** the earlier test changed only the `go.mod` directive while still compiling with a newer Go toolchain. Exact-toolchain tests fail on Go 1.23 (`os.OpenRoot` and `os.Root` unavailable) and Go 1.24 (`Root.ReadFile` unavailable), while Go 1.25 passes all 16 package suites. This is a user-owned compatibility-versus-backport choice; silently publishing 1.23 would make installation fail.

  **Answered 2026-08-31** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommended Go 1.25.0 floor via `do-work clarify` because exact-toolchain
  verification proves it is the lowest version that supports the current filesystem-safety design.
  Keeping Go 1.26.1 and undertaking the substantive Go 1.23 filesystem backport are out of scope.

## Notes
If the answer is to lower it, the change is mechanical and belongs in one small REQ: the `go` directive, `minimum_go_version` in `skills/do-work/tools/do-work-cli.sh`, the README prerequisite line, and the version-action prerequisite line, plus a check that the launcher's refusal message quotes the new number.

---
*Source: REQ-407 review (UR-081). Answer with `do-work clarify`.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The user has confirmed one exact compatibility-floor value and the REQ names the four literals that must change plus the launcher refusal check. No design work remains.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Compatibility Verification

The attempted literal-only change was discarded before commit after exact-toolchain testing invalidated its premise.

- `GOTOOLCHAIN=go1.23.0 go test -count=1 ./...` — failed: `os.OpenRoot` and `os.Root` are unavailable.
- `GOTOOLCHAIN=go1.24.0 go test -count=1 ./...` — failed: `(*os.Root).ReadFile` is unavailable.
- `GOTOOLCHAIN=go1.25.0 go test -count=1 ./...` — passed all 16 packages.
- Host Go 1.26.1 suite, launcher behavior check, contract regressions, and `git diff --check` — passed on the attempted literal update before it was discarded.

Go 1.23 support would replace a traversal-resistant filesystem boundary used by `atomicfile`, `repositorymodel`, and `cleanup`, including six `os.OpenRoot` sites and rooted open, create, stat, walk, and remove operations. No implementation commit was created. Go's official `os.Root` announcement places `os.Root`/`os.OpenRoot` in Go 1.24, and the Go 1.25 release notes introduce `Root.ReadFile`, matching the exact-toolchain failures.

If a lower truthful floor is selected, update every current core prerequisite restatement together: `README.md`, `skills/do-work/actions/version.md`, `skills/do-work/tools/do-work-cli.sh`, `skills/do-work/tools/do-work-cli/go.mod`, `skills/do-work/docs/prescribed-shell-primitives.md`, `skills/do-work/tools/prime-do-work-update.md`, and the compatibility-launcher comments in both root and shipped mirrors. The optional board module, `_dev/tests/maintainer-verify.sh`, historical changelogs, and archived reports remain unchanged.

## Implementation Summary

**What was done:** Lowered the core installer and updater floor to Go 1.25.0 across the module, launchers, and current prerequisite documentation. Added an exact-Go-1.25 compatibility lane and moved launcher boundary fixtures to the truthful accepted and rejected versions.

**Files changed:**
- `README.md` (modified) — lowered current install and update prerequisites.
- `_dev/tests/do-work-cli-go125-compatibility.sh` (new) — verifies exact Go 1.25.0 and runs every core package test.
- `_dev/tests/do-work-cli-launcher-behavior.sh` (modified) — tests Go 1.25.0 acceptance and Go 1.24.99 rejection.
- `skills/do-work/actions/version.md` (modified) — lowered the updater prerequisite restatement.
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — lowered the compatibility-launcher prerequisite.
- `skills/do-work/tools/do-work-cli.sh` (modified) — enforces Go 1.25.0 as the minimum.
- `skills/do-work/tools/do-work-cli/go.mod` (modified) — declares Go 1.25.0.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — maps the exact-version verification lane.
- `skills/do-work/tools/do-work-update.sh` (modified) — lowered the updater launcher comment.
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified) — lowered the shipped fetch launcher comment.
- `skills/do-work/tools/install-do-work-suite.sh` (modified) — lowered the shipped installer launcher comment.
- `skills/do-work/tools/prime-do-work-update.md` (modified) — lowered the canonical updater prerequisite.
- `skills/do-work/tools/replace-text-section.sh` (modified) — lowered the shipped managed-section launcher comment.
- `skills/do-work/tools/validate-suite-manifest.sh` (modified) — lowered the shipped validation launcher comment.
- `tools/fetch-upstream-archive.sh` (modified) — lowered the root launcher mirror comment.
- `tools/install-do-work-suite.sh` (modified) — lowered the root launcher mirror comment.
- `tools/replace-text-section.sh` (modified) — lowered the root launcher mirror comment.
- `tools/validate-suite-manifest.sh` (modified) — lowered the root launcher mirror comment.

**Integration range:** `afbac3c0..098936e8`

*Generated by work action from the builder hand-back*

## Qualification

Passed — 18 files verified against merge range `afbac3c0..098936e8`; the clarified Go 1.25 floor is synchronized across the core module, launchers, and current restatements, and the optional board and maintainer floors remain unchanged. P-A-U evidence and debug-artifact checks passed.

## Testing

**Merged-state checks:**
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — PASS (exit 0); exact `go1.25.0` selected and all 16 core packages passed.
- `bash _dev/tests/do-work-cli-launcher-behavior.sh` — PASS (exit 0); exact-floor acceptance and below-floor refusal passed.
- `bash _dev/tests/contract-regressions.sh` — PASS (exit 0); all contract and 118 named shell behavior cases passed.
- `shellcheck --severity=warning -- ...` over the new lane, launcher test, core launcher, and root/shipped compatibility launchers — PASS (exit 0).
- `DO_WORK_DIFF_RANGE="afbac3c0..098936e8" skills/do-work/tools/checks/qualify.sh ...` — PASS (exit 0).

**Regression evidence:** The dedicated exact-toolchain lane proves the complete core suite on Go 1.25.0; launcher fixtures prove 1.25.0 is accepted and 1.24.99 is rejected with the new prerequisite.

## Review

**Overall: 100%** | 2026-08-31T14:22:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — exact Go 1.25 build/run, boundary refusal, mirror lock-step, contract regressions, and the canonical maintainer gate all pass on merged main.
**Suggested testing:** 1 optional disposable-project install/update smoke on a system Go 1.25 environment.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Exact-toolchain selection turned the compatibility floor into executable proof rather than a `go.mod` edit tested by a newer compiler.
**What didn't:** The earlier Go 1.23 claim changed only the module directive, so the host toolchain hid unavailable standard-library APIs.
**Worth knowing:** The core CLI floor and the optional board/maintainer floor are independent and must remain separately stated and tested.

## Orientation

The core installer and updater now support Go 1.25 while preserving the existing rooted-filesystem safety design; exact-toolchain and launcher lanes keep that boundary truthful.
