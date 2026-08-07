---
id: REQ-139
title: "Stage the Modular Core Skill"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-137, REQ-138]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Stage the Modular Core Skill

## What
Create a self-contained staged `skills/do-work` package while the repository-root all-in-one distribution remains active.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Core needs an independently discoverable context boundary without weakening the orchestration behavior the user relies on.

## Detailed Requirements
- Stage capture, run, verify-requests, review-work, clarify, abandon, cleanup, commit, roadmap, forensics, version/update/recap, help, required references, queue schema, crew guardrails, specs, core hooks, checks, and updater.
- Preserve the current feature-rich work orchestrator.
- Keep pipeline temporarily for behavioral parity; REQ-145 removes it after cutover.
- Resolve eventual queue-kanban calls from sibling `do-work-board` and knowledge handoffs from sibling `do-work-knowledge`.
- Add package-boundary tests proving every staged runtime reference resolves and no shipped file cites repository-root maintainer instructions.
- Do not activate this package or alter the active root router yet.

## Constraints
- This is a move-by-copy staging step. Legacy active files remain until cutover.
- Avoid permanent duplicated sources; REQ-144 removes the old active copies.

## Dependencies
Requires the bridge updater and managed-section foundation in REQ-137 and REQ-138.

## Builder Guidance
Certainty level: Firm. Preserve behavior first; simplification is limited to functional ownership boundaries.

## Red-Green Proof
**RED prompt/case:** Point an isolated skill loader at `skills/do-work/SKILL.md` and validate all runtime references in a staged four-skill fixture.
**Why RED now:** No independently loadable core package exists and current references assume one monolithic skill root.
**GREEN when:** Core loads from its staged package, its required references resolve across the staged suite, and the active legacy installation remains unchanged.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the core behavior and module-boundary decisions.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
