---
id: REQ-416
title: 'Implement deterministic BKB and Dream commands'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-415]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T04:10:12Z
---

# Implement Deterministic BKB and Dream Commands

## What
Move deterministic knowledge-base and Dream scans into `do-work-cli` while retaining LLM-only judgment phases in actions.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement `bkb-init`, `bkb-status`, and `bkb-lint-structure` as direct commands and flat recipes.
- Implement Dream’s seven deterministic scans behind `dream-scan`.
- Return typed actionable findings in text and JSON with exact next and verification commands.
- Leave contradiction resolution, synthesis, and cluster design to the existing LLM actions, which must consume canonical command results.

## Constraints
- Preserve knowledge file formats and never convert judgment-heavy phases into brittle heuristics.

## Dependencies
Depends on REQ-415 (hook/runtime migration precedes package-specific knowledge commands).

## Builder Guidance
Certainty level: Firm for deterministic scans; explicitly retain action ownership for semantic judgment.

## Red-Green Proof
**RED prompt/case:** Run BKB and all seven Dream scan fixtures through absent CLI commands, including clean, malformed, and finding cases.
**Why RED now:** Deterministic phases are specified in action prose rather than exposed as one stable no-LLM interface.
**GREEN when:** Every scan is directly runnable, text/JSON agree, findings are actionable, and semantic resolution remains delegated to the LLM action.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
