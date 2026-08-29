---
id: REQ-408
title: 'Build shared request, schema, dependency, atomic-file, and repository packages'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-407]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Build Shared Request, Schema, Dependency, Atomic-File, and Repository Packages

## What
Create the reusable repository model required by every request and queue command.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement shared REQ/UR frontmatter parsing and writing, schema aliases and normalization, canonical timestamps, ID allocation, dependency graphs, atomic files, and repository discovery/modeling.
- Preserve unknown fields and bytes where commands are not authorized to rewrite them.
- Support queue, working, root archive, nested archived UR, reservation, and user-request layouts.
- Give downstream commands typed evidence and exact paths without rescanning through ad hoc shell pipelines.

## Constraints
- Use the Go standard library unless an existing dependency is demonstrably necessary.
- Treat this package layer as the single source of truth for later command families.

## Dependencies
Depends on REQ-407 (the Go installation/runtime path must exist before expanding the model).

## Builder Guidance
Certainty level: Firm. Reuse existing queue-kanban schema behavior where compatible, but do not merge the board binary into this module.

## Red-Green Proof
**RED prompt/case:** Parse representative current and legacy REQ/UR fixtures, malformed frontmatter, timestamp variants, reservation races, and dependency cycles with the absent shared packages.
**Why RED now:** Each shell/action path currently reconstructs parts of the repository model independently.
**GREEN when:** Unit and fixture tests expose one normalized typed model, preserve required bytes/fields, allocate collision-free IDs, and return deterministic dependency results.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
