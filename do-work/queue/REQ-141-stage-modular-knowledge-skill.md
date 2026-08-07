---
id: REQ-141
title: "Stage the Modular Knowledge Skill"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
tdd: true
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Stage the Modular Knowledge Skill

## What
Create a staged `skills/do-work-knowledge` package for BKB, dream, memory, interview, prompts, knowledge assets, and knowledge hooks.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Knowledge management and session memory are self-contained concerns that accumulated far outside the core task-queue mission.

## Detailed Requirements
- Add a dedicated `SKILL.md` and move-by-copy BKB, dream, memory, interview, prompts, references, templates, docs, hooks, and interviewer guidance.
- Make this package own memory-module setup and hook configuration.
- Keep memory capture disabled on a fresh full-suite install.
- Allow explicit setup to enable memory hooks.
- Define deterministic migration targets from `.claude/skills/do-work/hooks/memory-*.sh` to `.claude/skills/do-work-knowledge/hooks/`.
- Preserve memory privacy, machine-local raw stores, bootstrap semantics, KB behavior, and optional core lessons handoff.
- Do not delete or deactivate legacy copies yet.

## Constraints
- Existing knowledge behavior must remain unchanged through staging.
- Do not enable memory transcript capture merely because the suite is installed.

## Dependencies
Requires REQ-136's suite contract.

## Builder Guidance
Certainty level: Firm. The package is broad but cohesive around retained knowledge and memory functionality.

## Red-Green Proof
**RED prompt/case:** Load BKB, memory, dream, interview, and prompt actions from an isolated `skills/do-work-knowledge` package.
**Why RED now:** These actions and their hooks/templates are coupled to the monolithic skill root.
**GREEN when:** Every knowledge route resolves within the staged package, existing behavioral tests pass, and fresh-install fixtures prove memory hooks remain disabled.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for knowledge ownership and hook policy.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
