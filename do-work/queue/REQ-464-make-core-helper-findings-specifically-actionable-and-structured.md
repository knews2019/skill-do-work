---
id: REQ-464
title: 'Review fix: Make core helper findings specifically actionable and structured'
status: pending
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
sweep_key: core-helper-specific-structured-projections
---

# Make Core Helper Findings Specifically Actionable and Structured

## What

Project each core-helper finding as its own structured fact with recovery and verification commands that address that exact condition. Done means neither text nor JSON asks an agent to rescan opaque evidence or infer a remedy from a family-wide fallback.

The fold-first scan found no eligible pending or pending-answers REQ, sweep or otherwise, in any UR that shares this lossy-projection root cause. REQ-420 names actionable parity as a final acceptance goal but is dependency-gated and cannot receive this finding under the fold-first conversion rule.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `internal/corehelpers/commands.go`: code-family fallbacks prescribe generic `git status`, `git diff`, or a command rerun even when those argv do not resolve or specifically verify the named finding. (found by REQ-414 / UR-081)
- [ ] `internal/corehelpers/handoff.go`: all dirty porcelain rows are stored in one newline-bearing evidence string, so individual affected paths and states are not typed records. (found by REQ-414 / UR-081)

## Requirements

- Give each non-clean finding next and verification argv whose success specifically resolves or proves that condition; do not substitute family-wide inspection commands.
- Emit one typed handoff record per dirty path/state, including spaces, newlines, renames, missing/prunable worktrees, and multiple simultaneous rows.
- Preserve deterministic text/JSON parity from the same observation set without multiline opaque evidence blobs.
- Add mutation tests that replace a specific action or collapse structured dirty rows and prove the characterization gate fails.

## Red-Green Proof

**RED prompt/case:** Produce an unchecked P-A-U finding and a two-path dirty handoff, then inspect their next/verification argv and JSON path records.
**Why RED now:** The commands prescribe generic rescans and serialize both dirty paths into one opaque evidence value.
**GREEN when:** Every finding has condition-specific recovery/verification and every dirty path/state is independently represented and actionable in both renderers.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`.

---
*Source: REQ-414 fresh re-review finding 3.*
