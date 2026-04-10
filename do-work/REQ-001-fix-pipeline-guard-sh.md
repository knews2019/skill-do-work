---
id: REQ-001
title: Fix pipeline-guard.sh quoting and error handling
status: pending
created_at: 2026-04-10T19:30:00Z
user_request: UR-001
domain: backend
prime_files: []
tdd: false
---

# Fix pipeline-guard.sh quoting and error handling

## What
Fix two shell scripting issues in `hooks/pipeline-guard.sh`: (1) unquoted command substitution on line 27 (`INPUT=$(cat)` should be `INPUT="$(cat)"`), and (2) masked errors on line 53 where `2>/dev/null` on the numeric comparison hides real failures when PENDING is empty or non-numeric — add numeric validation before the comparison instead.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context
Identified in the quick-wins report scan of the repository. These are the only executable code files in the project. The script is a Claude Code stop hook that prevents stopping mid-pipeline — correctness matters for safety.

## Red-Green Proof
**RED prompt/case:** `echo '{}' | INPUT=$(cat) && echo "$INPUT"` — with certain stdin content containing special characters, unquoted substitution can break. Also: `PENDING="" && [ "$PENDING" -gt 0 ] 2>/dev/null` silently succeeds/fails unpredictably.
**Why RED now:** Unquoted command substitution and masked comparison errors are latent bugs.
**GREEN when:** Line 27 uses `INPUT="$(cat)"` with double quotes, and line 53 validates PENDING is numeric before the `-gt` comparison without suppressing stderr.
**Validation:** Inferred during capture

---
*Source: Address quick-wins report findings*
