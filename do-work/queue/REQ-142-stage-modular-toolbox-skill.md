---
id: REQ-142
title: "Stage the Modular Toolbox Skill"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
tdd: true
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Stage the Modular Toolbox Skill

## What
Create a staged `skills/do-work-toolbox` package preserving the current optional reviews, discovery, reporting, repository utilities, and companion installers.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The user wants these capabilities retained, but they should not occupy the core queue router and context.

## Detailed Requirements
- Preserve validate-feedback, code-review, ui-review, present-work, ai-report, slop-check, quick-wins, scan-ideas, deep-explore, prime, inspect, note, stray-check, tidy-repo, tutorial, and companion installers not owned by board or knowledge.
- Move-by-copy every required reference, crew file, template, and guide.
- Add a dedicated router/help contract.
- Preserve links to core queue artifacts where actions read or report on URs/REQs.
- Add consistency tests ensuring every toolbox route exists exactly once and every runtime reference resolves in the staged suite.
- Do not remove or deactivate legacy copies yet.

## Constraints
- Retain all toolbox functionality; this program does not run a usage-pruning phase.
- Board setup belongs to board; memory setup belongs to knowledge; core self-update belongs to core.

## Dependencies
Requires REQ-136's suite contract.

## Builder Guidance
Certainty level: Firm. Preserve behavior while moving ownership; avoid redesigning individual actions during this stage.

## Red-Green Proof
**RED prompt/case:** Invoke each toolbox route through an isolated `skills/do-work-toolbox/SKILL.md` fixture and resolve every referenced runtime file.
**Why RED now:** These optional actions are routed and stored inside the monolithic skill.
**GREEN when:** All toolbox routes dispatch from the staged package, all references resolve, and legacy active behavior remains unchanged.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the retained toolbox action set.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
