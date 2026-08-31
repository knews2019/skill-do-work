---
id: REQ-426
title: 'Preserve setuid, setgid and sticky bits on managed files instead of stripping them'
status: completed
created_at: 2026-08-30T17:40:00Z
claimed_at: 2026-08-31T10:30:10Z
completed_at: 2026-08-31T10:56:05Z
commit:
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-31T10:30:10Z
  basis:
    - trivial short-circuit
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
addendum_to: REQ-407
review_generated: true
write_set: [skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go, skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go]
---

# Preserve setuid, setgid and sticky bits on managed files instead of stripping them

## What
REQ-407's Go port reads file modes with `info.Mode().Perm()` (mask `0o777`) where the Python it replaced read `stat.S_IMODE(...)` (mask `0o7777`). The three special bits are silently dropped from `Justfile`, `CLAUDE.md` and `.claude/settings.json` on every install and every update.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Added RED regressions for managed-section replacement and real installer paths before changing the two mode readers.
- [x] **[APPLY]:** Applied the special-bit mask only at the two requested production reads and extended their existing tests.
- [x] **[UNIFY]:** Reviewed all four modified files, ran `gofmt`, `git diff --check`, focused and module-wide Go tests, and the canonical maintainer gate; no debug artifacts remained.

## Detailed Requirements
- `permissionsOf` in `managed_section.go` and the settings-mode read in `install_transaction.go` must carry the setuid, setgid and sticky bits through, not just the low nine.
- The regression must be pinned by a test that fails before the fix: at minimum one `managedsection` case on a `0o2644` target and one install case asserting the setgid bit survives a real install.

## Constraints
- Do not widen this into a general mode-handling change. Ordinary permission bits already round-trip correctly and are covered.
- `os.Chmod` and `*os.File.Chmod` already map Go's `ModeSetuid`/`ModeSetgid`/`ModeSticky` back to the syscall bits, so the fix is a mask, not a new syscall path.

## Builder Guidance
Certainty level: Firm. The reviewer's proposed form is `info.Mode().Perm() | (info.Mode() & (os.ModeSetuid|os.ModeSetgid|os.ModeSticky))`.

## Red-Green Proof
**RED prompt/case:** Seed a fixture `Justfile` at `2644`, `CLAUDE.md` at `4644` and `.claude/settings.json` at `1644`, run the installer, and read the modes back.
**Why RED now:** Measured A/B on the merged tree against the pre-REQ shell installer, identical fixtures. The pre-REQ installer returned `2644` / `4644` / `1644` unchanged; the Go installer returned `644` / `644` / `644`. At the unit level, `7644` becomes `644` and `6755` becomes `755` — the execute bits survive and only the special bits are lost.
**GREEN when:** All three files keep their full twelve-bit mode through an install, and the unit cases pass.
**Validation:** Reproduced three ways during REQ-407's review — unit level, A/B against the pre-REQ Python replacer, and end-to-end through both installers.

## Notes On Reach
Adjudicated 2-1 and judged **Minor** rather than Important by its own verifiers, for a reason worth carrying: git records only `100644`/`100755`, `umask` cannot set special bits, and a setgid *directory* propagates the group to new files rather than the setgid bit. So reaching this requires a consumer to have hand-run `chmod g+s` (or `u+s`, or `+t`) on one of these three files. Real, narrow, and cheap to close — but do not let it grow.

---
*Source: REQ-407 review (UR-081).*

---

## Triage

**Route: A** - Simple

**Reasoning:** The root cause, two production read sites, exact mask, and required regression cases are all specified. This is a focused mode-preservation fix in four named Go files.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)

**What was done:** Both managed-section and installer mode reads now preserve setuid, setgid, and sticky bits alongside ordinary permissions. RED-first regressions cover a `02644` managed target and a real reinstall retaining `Justfile` 2644, `CLAUDE.md` 4644, and settings 1644.

## Root Cause

The Go port used `FileMode.Perm()`, which intentionally returns only the low nine permission bits, while the replaced Python used `stat.S_IMODE` and retained the full `07777` mode. Candidate files were therefore published after silently clearing all three special bits.

## Decisions

- **D-01 (DECIDE & STATE):** Used the requested standard-library mode mask instead of adding a syscall path. This matches how Go's chmod functions consume special `FileMode` bits.
- **D-02 (DECIDE & STATE):** Extended the existing broad reinstall fixture instead of adding a parallel installer fixture, keeping the full managed-file proof in one real-install path.

## Discovered Tasks

- Assess whether the permission-preservation contracts at `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go:55` and `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:245` also require setuid, setgid, and sticky-bit preservation; both contain the same `Mode().Perm()` shape outside this REQ's write boundary.

## Qualification

Passed — 4 files verified in merge range `29f07ae0..73bd4a6f`, both stated requirements traced, all changes substantive and in scope, P-A-U confirmed, and no debug artifacts or hollow paths found.

## Testing

**Tests run:** `go test ./internal/managedsection ./internal/suiteinstall`; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing; the canonical gate exited 0 on the merged main tree. Its optional strict browser lane was skipped because no browser was available.

**Red-green validation:**
- Managed-section `02644` replacement: ✗ before (`644`) → ✓ after (`2644`)
- Real installer fixture: ✗ before (`Justfile`, settings, and `CLAUDE.md` each lost special bits) → ✓ after (2644, 1644, and 4644 retained)

**Existing tests updated (cross-REQ impact):**
- `managed_section_test.go` and `install_transaction_test.go` (from REQ-407): extended the original mode-preservation contract to cover the special-bit behavior REQ-426 was created to close.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-31T10:56:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — both captured special-mode regressions pass through the real production seams on the merged tree.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Exact Unix-mode fixtures exposed the semantic difference between Go's `Perm()` mask and Python's `S_IMODE` immediately.
**What didn't:** The earlier ordinary-mode tests could not detect special-bit loss because they never left the low nine bits.
**Worth knowing:** Go represents setuid, setgid, and sticky as high `FileMode` flags; numeric Unix fixtures need explicit conversion in tests before `Chmod`.

## Orientation

Managed configuration replacement now retains the complete Unix mode, including setuid, setgid, and sticky bits, across install and update.
