---
id: REQ-462
title: 'Review fix: Preserve deletion precedence across Git XY inventory states'
status: claimed
created_at: 2026-09-01T02:23:45Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-414]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-414, REQ-420]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-414
sweep: true
sweep_key: git-porcelain-xy-deletion-precedence
claimed_at: 2026-09-01T08:42:07Z
---

# Preserve Deletion Precedence Across Git XY Inventory States

## What

Classify every Git porcelain XY state from the path's usable filesystem condition, with deletion taking precedence whenever either side says the path is absent. Done means combined status states cannot be mistaken for readable or associable paths.

The fold-first scan found no eligible pending or pending-answers REQ, sweep or otherwise, in any UR that shares this Git XY classification root cause. REQ-420 covers final whole-suite parity but is dependency-gated and does not yet carry this classifier fix.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `internal/corehelpers/inventory.go`: porcelain `AD` is emitted as `INVENTORY-A` even though the staged-new path was deleted and cannot be read or associated. (found by REQ-414 / UR-081)
- [ ] Git XY classifier: deletion, rename/copy, unmerged, and secret-derived combined states need one differential precedence matrix against the retained inventory helper. (found by REQ-414 / UR-081)

## Requirements

- Define classification over both porcelain status columns and the actual path contract; an absent/deleted path must never be projected as readable.
- Add a differential matrix covering `AD`, staged and worktree deletes, renames, copies, unmerged states, ordinary modifications, and secret-derived variants.
- Preserve ordered typed facts, statuses, and association/quarantine effects for every characterized state.

## Red-Green Proof

**RED prompt/case:** Stage a new file, remove it before commit, and compare the retained inventory helper with `do-work-cli uncommitted-inventory` for the resulting `AD` row.
**Why RED now:** The Go classifier checks `A` before `D`, emits a readable-addition fact, and diverges from the deletion-first retained contract.
**GREEN when:** The direct replay and the complete XY differential matrix agree on status, facts, paths, and downstream associability.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`.

---
*Source: REQ-414 fresh re-review finding 1.*
