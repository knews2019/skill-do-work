---
id: REQ-417
title: 'Implement interview and deterministic memory store commands'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-416]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Implement Interview and Deterministic Memory Store Commands

## What
Expose all deterministic interview and memory operations through `do-work-cli`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement interview list, status, export, ingest, reset, and versions commands.
- Implement memory remember, forget, recall, status, bootstrap, and audit commands.
- Preserve store formats, ordering, redaction, deduplication, version semantics, and atomic mutations.
- Provide direct commands, flat Just recipes, text/JSON parity, dry-run where meaningful, optional commit, and actionable findings.

## Constraints
- Natural-language knowledge actions remain aliases and delegate deterministic phases to the CLI.

## Dependencies
Depends on REQ-416 (knowledge command conventions and scan result integration).

## Builder Guidance
Certainty level: Firm. Characterize rendering and store mutations before migration.

## Red-Green Proof
**RED prompt/case:** Exercise each interview and memory operation against representative stores, malformed data, duplicate records, redaction cases, dry-run, and rollback.
**Why RED now:** These operations are not uniformly available as stable direct Go commands.
**GREEN when:** Every operation has fixture-proven deterministic output/effects, matching text/JSON, and its action alias delegates without free-form fallback.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
