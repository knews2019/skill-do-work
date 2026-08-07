---
id: REQ-144
title: "Activate the Four-Skill Distribution"
status: blocked
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-137, REQ-143]
maintenance: true
blocked_by: "user confirmation that every existing client updater reports suite-layout-v2"
blocked_at: 2026-08-07T18:58:02Z
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Activate the Four-Skill Distribution

## What
Switch the live distribution to the four staged skills and migrate bridge-enabled client installations atomically.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Staging alone does not reduce the installed router or context; the manifest and live archive must cut over only after clients can understand the new layout.

## Detailed Requirements
- Activate the live four-module manifest.
- Make `skills/do-work`, `skills/do-work-board`, `skills/do-work-knowledge`, and `skills/do-work-toolbox` the only shipped runtime sources.
- Remove duplicated legacy root actions/tools in the same change.
- Update export rules, README installation, root maintainer Just fallback, changelog/version handling, and every runtime cross-skill reference.
- Add one-release core shims for moved commands. Each shim prints the exact new invocation and stops; it must not forward permanently.
- Ensure `do-work update` and `just run-do-work-update` produce byte-equivalent installations from a bridge client.
- Refresh the managed Just section and migrate known memory-hook paths during update.
- Preserve project queue data, KB data, application files, and unrelated configuration.
- Run hermetic bridge-to-modular and fresh-install tests plus ShellCheck, formatting, Go tests/vet, contract tests, and `git diff --check`.

## Constraints
- Do not begin until the user confirms every existing client has the bridge capability.
- Preserve the current feature-rich work orchestrator.
- Keep compatibility shims for exactly one modular release.

## Dependencies
Requires REQ-137 and REQ-143. It is also held on the external bridge-rollout confirmation recorded in frontmatter.

## Builder Guidance
Certainty level: Firm. Treat this as the release cutover; no partial module activation is acceptable.

## Red-Green Proof
**RED prompt/case:** Update a hermetic legacy client that has the bridge updater, legacy Just recipes, queue data, KB data, and an enabled legacy memory hook.
**Why RED now:** The live manifest and installed layout remain monolithic.
**GREEN when:** One update installs exactly four current skill directories, migrates managed configuration, preserves project data, and both update entry points yield identical bytes.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the bridge rollout gate and cutover contract.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
