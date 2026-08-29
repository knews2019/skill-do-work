---
id: REQ-411
title: 'Implement dependency-aware queue selection and actionable summaries'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-410]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Implement Dependency-Aware Queue Selection and Actionable Summaries

## What
Move deterministic queue selection and readiness decisions into canonical Go commands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement `do-work-next` and shared selection for explicit REQ/UR targeting, dependency readiness, cycle detection, waves, assignments, negligible-impact filtering, blocked probes, and estimates.
- Compose stable queue summaries from the same typed repository model and return actionable reasons for every skipped or refused item.
- Preserve explicit targeting semantics and successful-dependency gates.
- Provide text/JSON output with exact next argv/Just recipes and verification commands.

## Constraints
- Keep `queue-kanban` as the separate board/UI binary; share behavior through contracts or compatible parsing, not binary consolidation.

## Dependencies
Depends on REQ-410 (doctor and normalized repository findings).

## Builder Guidance
Certainty level: Firm. Lock current selection semantics with fixtures before replacing action-side scans.

## Red-Green Proof
**RED prompt/case:** Select from a fixture containing targeted IDs, a dependency chain, a cycle, an assignment, negligible impact, blocked probes, and wave depth.
**Why RED now:** Selection logic is distributed between action prose and small shell helpers rather than one runnable command.
**GREEN when:** The CLI returns the same eligible set and stable reason for every excluded item in text and JSON, without LLM judgment for deterministic gates.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
