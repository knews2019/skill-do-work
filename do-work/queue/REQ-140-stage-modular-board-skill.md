---
id: REQ-140
title: "Stage the Modular Board Skill"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md, tools/prime-do-work-update.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-136, REQ-138]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Stage the Modular Board Skill

## What
Create a staged `skills/do-work-board` package that owns the board action, board documentation, Just template, and complete queue-kanban Go module.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The board is a distinct compiled application with its own UI and regression surface and should not tax or clutter the core router.

## Detailed Requirements
- Add a dedicated `SKILL.md`, board action/help, documentation, Just template, and all `tools/queue-kanban` source/tests/assets.
- Preserve live, static, summary, CLI, Testing-view, and completion-calendar behavior.
- Preserve port validation, foreign-process protection, bounded shutdown, and browser opening.
- Use the managed `do-work:recipes` section.
- Point board recipes at `.claude/skills/do-work-board/tools/queue-kanban`.
- Point `run-do-work-update` at the core updater.
- Do not remove the legacy board files yet.

## Constraints
- The board continues to read core queue data and the shared schema contract.
- The Go tool remains a single source after cutover; no permanent duplicate implementation.

## Dependencies
Requires REQ-136's suite contract and REQ-138's managed recipe mechanism.

## Builder Guidance
Certainty level: Firm. Preserve all validated queue-kanban fixes, including the Testing done-window empty-state behavior.

## Red-Green Proof
**RED prompt/case:** Load and run board commands from `skills/do-work-board` in an isolated staged suite.
**Why RED now:** Board routing, docs, recipes, and compiled source live inside the monolithic core root.
**GREEN when:** The staged board skill builds and passes all Go tests, its recipes resolve the new install path, and legacy active behavior remains green.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the board and installation boundaries.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
