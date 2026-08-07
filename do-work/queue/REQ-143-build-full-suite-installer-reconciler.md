---
id: REQ-143
title: "Build the Full-Suite Installer and Reconciler"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-138, REQ-139, REQ-140, REQ-141, REQ-142]
maintenance: false
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [tools/install-do-work-suite.sh, suite/modules.tsv, README.md, _dev/tests/contract-regressions.sh]
---

# Build the Full-Suite Installer and Reconciler

## What
Build the canonical fresh-install bootstrap and client-configuration reconciler for the complete four-skill suite, without activating the live modular distribution yet.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Clients need one reliable install command and automated migration of the two project-owned configuration surfaces affected by the split.

## Detailed Requirements
- Add `tools/install-do-work-suite.sh` requiring `--project-root` and a Git repository.
- Download one upstream archive and validate `VERSION`, `suite/modules.tsv`, and all four non-empty `SKILL.md` files in staging before writing.
- Always install the full suite; never install a subset.
- Create a complete Justfile from the board template when none exists, or migrate/replace only the managed `do-work:recipes` section in an existing file.
- Enable core Claude hooks by composing with valid existing settings.
- Do not enable memory hooks on fresh installs.
- If known legacy memory-hook command strings already exist, rewrite only those strings to `do-work-knowledge`.
- Use `jq` when available, Python 3 as fallback, and otherwise leave settings unchanged with an exact manual instruction.
- Restore exact temporary originals when Just/settings validation fails.
- Cover fresh install, reinstall, legacy recipes, custom recipes, invalid markers/settings, hook composition, memory-hook migration, spaces, interruption, and exact four-module verification.
- Keep the live repository in bridge mode until REQ-144.

## Constraints
- All four modules share one version and one confirmation boundary.
- Installer-owned reconciliation must not overwrite unrelated client configuration.

## Dependencies
Requires the managed-section utility and all four staged skill packages.

## Builder Guidance
Certainty level: Firm. Reuse the manifest validator and managed-section utility instead of adding parallel implementations.

## Red-Green Proof
**RED prompt/case:** Run the bootstrap against an empty Git project and against a project with custom Just recipes and existing core/memory hooks.
**Why RED now:** There is no full-suite installer, no four-module verification, and no deterministic configuration reconciliation.
**GREEN when:** Both fixtures receive the same four-version suite; custom content survives, the managed block is current, core hooks are enabled, and memory hooks are only migrated when already present.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for bootstrap, Justfile, hook, and Git policies.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
