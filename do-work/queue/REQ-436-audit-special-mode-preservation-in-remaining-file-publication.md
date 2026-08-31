---
id: REQ-436
title: '[impact-negligible] Audit special-mode preservation in remaining file publication'
status: pending-answers
created_at: 2026-08-31T10:56:05Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
depends_on: [REQ-426]
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
addendum_to: REQ-426
sweep: true
sweep_key: preserve-special-mode-bits-in-file-publication
---

# Audit Special-Mode Preservation in Remaining File Publication

## What

REQ-426 fixed two managed install paths that silently narrowed Unix modes to the low nine permission bits. The same `Mode().Perm()` publication shape remains in atomic replacement and cleanup moves. Audit both contracts, preserve setuid/setgid/sticky where they promise to preserve a source file's mode, and pin the class so it cannot recur.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go:55` narrows the original target mode before publishing its replacement.
- [ ] `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:245` narrows the source mode before publishing a cleanup move destination.

## Red-Green Proof

**RED prompt/case:** Replace and move regular-file fixtures carrying setuid, setgid, and sticky bits through the two named public seams.
**Why RED now:** Both production calls pass `Mode().Perm()`, the same low-nine-bit mask that REQ-426 proved strips all three special bits.
**GREEN when:** Each contract either preserves the complete mode with RED/GREEN regression proof, or explicitly documents and tests why narrowing is intentional; the audit contains no unclassified `Mode().Perm()` publication path.
**Validation:** Discovered during REQ-426 implementation; apply the finding-closure ratchet.

## Open Questions

- [ ] I found two remaining file-publication paths with the same special-mode-bit narrowing pattern fixed by REQ-426. Should I process this low-reach audit as a new task?
  Recommended: Yes, add to queue (will flip to `pending`) so every mode-preservation promise uses the same complete-mode contract.
  Also: No, discard it; these bits are rarely set on the affected files and REQ-426 already closes the reported installer paths.
  Value: prevents the same silent metadata loss in atomic replacement and cleanup moves.
  Risk: low and reversible; the work is a focused two-path audit with regression tests, but it adds queue work for an uncommon filesystem edge case.

---
*Source: discovered while implementing REQ-426 (UR-081).*
