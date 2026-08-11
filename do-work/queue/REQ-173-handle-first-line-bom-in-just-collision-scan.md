---
id: REQ-173
title: Handle first-line BOM in Just collision scan
status: pending
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-172, REQ-174]
batch: accepted-p2-fixes
write_set: [tools/replace-text-section.sh, skills/do-work/tools/replace-text-section.sh, _dev/tests/install-suite-behavior.sh]
---

# Handle First-Line BOM in Just Collision Scan

## What

Recognize a reserved Just recipe when the first line begins with a UTF-8 BOM, including when `just` is unavailable, without changing the target file's bytes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Just accepts a BOM-prefixed first definition, but the fallback scanner currently misses it. Without `just`, installation can report success while leaving a duplicate reserved recipe.

## Detailed Requirements

- Strip exactly one leading UTF-8 BOM from the first line's classification view only.
- Preserve all target bytes.
- Keep both distributed helper copies byte-identical.
- Reject the collision before confirmation or client mutation when `just` is unavailable.
- Avoid broader Unicode or encoding normalization.

## Constraints

- Preserve the no-Just collision scanner earned by REQ-152.
- Retain the existing multiline-literal behavior from REQ-162.

## Red-Green Proof

**RED prompt/case:** Install into a project whose marker-free Justfile starts with UTF-8 BOM bytes followed by `run-kanban:`, with `just` absent from `PATH`.
**Why RED now:** The helper misses the definition, installation returns success, and the resulting Justfile contains a duplicate recipe.
**GREEN when:** The installer rejects the reserved collision before confirmation or mutation while preserving the BOM-prefixed Justfile byte-for-byte.
**Validation:** User confirmed

## Full Context

See `do-work/user-requests/UR-039/input.md` and the preceding validated-feedback report.

---
*Source: fix accepted*
