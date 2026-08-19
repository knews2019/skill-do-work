---
id: REQ-262
title: Govern the prompt-kit templates' date headers
status: pending
created_at: 2026-08-18T19:30:47Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-055
addendum_to: REQ-253
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work-knowledge/prompts/**
---

# Govern the Prompt-Kit Templates' Date Headers

## What

Three prompt-kit templates in do-work-knowledge carry `Date: [today]` headers that no paragraph of the Timestamp rule governs and that sit outside the citation checker's reach (they are template content, not action prose). Decide whether they join the date-only paragraph's consumer list (UTC, cited) or are declared template-content-out-of-scope like the fenced-block exemption.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Discovered by REQ-253's builder's `[today]` class grep ([low]). The exact template paths are re-derived at claim time by the same grep — the count, not the list, is the contract.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-253: three prompt-kit templates carry ungoverned `Date: [today]` headers. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — prompt templates are consumer-facing content and can stay ungoverned.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.
