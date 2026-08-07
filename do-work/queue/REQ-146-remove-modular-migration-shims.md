---
id: REQ-146
title: "Remove Modular-Migration Compatibility Shims"
status: blocked
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-144, REQ-145]
maintenance: true
blocked_by: "one modular suite release has shipped and the user has confirmed every client migrated its four skill directories, managed Just section, and enabled memory-hook paths"
blocked_at: 2026-08-07T18:58:02Z
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145]
batch: do-work-four-skill-suite
---

# Remove Modular-Migration Compatibility Shims

## What
Remove transitional core routes and migration rules after one modular release and confirmed client migration.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Permanent shims would preserve the same router complexity the modularization is meant to remove.

## Detailed Requirements
- Remove old core routing shims for board, knowledge, and toolbox commands.
- Remove bridge-only manifest compatibility, stale-path deletion rules, and migration messages that no supported client needs.
- Preserve the full-suite bootstrap, current modular updater, managed-section reconciliation, and current-to-current modular updates.
- Ensure core help lists only core commands plus concise pointers to the three sibling skills.
- Prove each active action exists in exactly one skill and old invocations no longer appear as live core routes.

## Constraints
- Do not start until at least one modular release has shipped.
- Require user confirmation that every client has four skill directories, current managed Just paths, and migrated enabled memory hooks.

## Dependencies
Requires REQ-144 and REQ-145. It is also held on the external stabilization and client-confirmation condition in frontmatter.

## Builder Guidance
Certainty level: Firm. Delete transitional machinery; do not replace it with permanent forwarding aliases.

## Red-Green Proof
**RED prompt/case:** Inspect core routing/help and updater compatibility branches after the first modular release.
**Why RED now:** One-release shims and bridge migration paths remain intentionally active.
**GREEN when:** Transitional paths are absent, current modular installs still update successfully, and core exposes only its own functionality plus discovery pointers.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the one-release compatibility decision.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
