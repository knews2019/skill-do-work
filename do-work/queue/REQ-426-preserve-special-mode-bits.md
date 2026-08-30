---
id: REQ-426
title: 'Preserve setuid, setgid and sticky bits on managed files instead of stripping them'
status: pending
created_at: 2026-08-30T17:40:00Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
