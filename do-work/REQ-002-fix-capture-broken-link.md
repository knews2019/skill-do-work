---
id: REQ-002
title: Fix broken link in capture.md
status: pending
created_at: 2026-04-10T19:30:00Z
user_request: UR-001
domain: general
prime_files: []
tdd: false
---

# Fix broken link in capture.md

## What
Fix the broken relative link on line 149 of `actions/capture.md`. The current link `[user-requests/UR-NNN/input.md](./user-requests/UR-NNN/input.md)` doesn't resolve because `user-requests/` is under `do-work/`, not alongside the action files in `actions/`. Replace with a plain path reference since this is a template/example, not a real navigable link.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context
Line 149 of capture.md is inside a template example showing the Complex REQ format. The link is illustrative — it shows where verbatim input lives. Since `UR-NNN` is a placeholder, this was never a real navigable link anyway.

## Red-Green Proof
**RED prompt/case:** Reading `actions/capture.md` line 149 shows a relative link `./user-requests/UR-NNN/input.md` that doesn't resolve from the `actions/` directory.
**Why RED now:** The link path is incorrect relative to where the file lives.
**GREEN when:** The reference uses `do-work/user-requests/UR-NNN/input.md` as a descriptive path without a broken relative link.
**Validation:** Inferred during capture

---
*Source: Address quick-wins report findings*
