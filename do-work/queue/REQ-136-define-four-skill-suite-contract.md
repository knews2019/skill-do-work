---
id: REQ-136
title: "Define the Four-Skill Suite Contract"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-135]
maintenance: true
related: [REQ-135, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [suite/modules.tsv, VERSION]
---

# Define the Four-Skill Suite Contract

## What
Define the architecture and machine-readable distribution contract for a four-skill, atomic-version do-work suite without activating the modular layout yet.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The split needs one stable ownership, packaging, and versioning contract before files or updater behavior move.

## Detailed Requirements
- Add an ADR covering `do-work`, `do-work-board`, `do-work-knowledge`, and `do-work-toolbox`.
- Define source packages under `skills/<name>/` and client destinations under `.claude/skills/<name>/`.
- Define `suite/modules.tsv` as the strict source-to-destination manifest.
- Define root `VERSION` as the single suite version.
- Reject absolute paths, `..` traversal, duplicate sources, duplicate destinations, unknown columns, missing `SKILL.md` files, and destinations outside `.claude/skills`.
- Document that all four modules install and update together.
- Leave the current all-in-one runtime and archive layout active for the bridge release.

## Constraints
- Contract and tests only; do not move active actions or tools in this REQ.
- The contract must work without independent module versions.

## Dependencies
Requires REQ-135 so the suite starts from a warning-clean baseline.

## Builder Guidance
Certainty level: Firm. Prefer one canonical manifest definition over repeated lists.

## Red-Green Proof
**RED prompt/case:** Validate a fixture manifest containing traversal, a duplicate destination, or a module without `SKILL.md`.
**Why RED now:** No suite manifest or validator exists, so the repository cannot distinguish a valid four-skill distribution from unsafe or incomplete layouts.
**GREEN when:** Valid four-module fixtures pass and every malformed/escaping fixture is rejected without writing outside the staging area.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the complete suite boundaries and installation decisions.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
