---
id: REQ-549
title: '[impact-negligible] Drop the eight dead path tokens from the decision indexes and the lessons prime'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [decisions/topics/_index_skill-architecture.md, decisions/topics/_index_knowledge-base.md, decisions/audits/2026-08-05-shell-logic-in-prose-census.md, decisions/imported-specs/2026-04-12_close-gaps-in-interview.md, _dev/primes/lessons-action-files.md, _dev/tests/audit-lockins.sh]
---

# Drop the eight dead path tokens from the decision indexes and the lessons prime

## What
Eight path tokens cited in `decisions/` and `_dev/primes/` resolve to no tracked file (31 citation lines). Drop the dead entries from the two topic indexes, fix the bare `work-reference.md` in the lessons prime, and mark the census and imported-spec pointers retired the way the ADRs already do. Leave the ADRs untouched.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Agents read the topic indexes and the prime as live routing; a pointer to a deleted file costs a search on every read and teaches the reader that the records are unreliable.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 8, sweep_key `dead-path-pointers-in-records`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag SIMPLE; expected net line delta -8. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `decisions/topics/_index_skill-architecture.md` lists `crew-members/karpathy.md` as a live source: remove the entry.
- `decisions/topics/_index_knowledge-base.md` lists `actions/build-knowledge-base.md` as a live source: remove the entry.
- `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` cites `hooks/pipeline-guard.sh`: mark retired (the hook is gone; `internal/settingshooks/settings_hooks.go` keeps the name only to remove it from consumer settings).
- `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md` names `decisions/imported-specs/expand-skill-do-work-interview.md` as required reading on three lines: mark retired.
- `_dev/primes/lessons-action-files.md` cites a bare `work-reference.md`, unresolvable from `_dev/primes/`: write the full `skills/do-work/actions/work-reference.md` path like every other prime.
- Not touched: the ten ADRs (`adr-001,003,005,006,007,008,009,010,013,014`) are immutable history; `decisions/audits/2026-08-11-defensive-surface.md` is a frozen snapshot by its own header.
- Reproduce at dc8a64e3: `for t in crew-members/karpathy.md actions/build-knowledge-base.md hooks/pipeline-guard.sh actions/pipeline.md actions/pipeline-reference.md _dev/tests/record-commit-hash-guards.sh skills/do-work-knowledge/crew-members/security.md decisions/imported-specs/expand-skill-do-work-interview.md; do git ls-files | grep -q --fixed-strings "$t" || rg -n --fixed-strings "$t" _dev/primes decisions; done`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (create it on first use, executable, invoked from `_dev/tests/contract-regressions.sh` in the fast tier the way `_dev/tests/defensive-surface-audit.sh` is, with the same missing-or-not-executable FAIL line), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Prose only outside the lock-in; no shipped file under `skills/` changes.
- Lock-in limit: dead path tokens cited outside decisions/records: 0 after this REQ (today 8); red when the Reproduce command prints a non-ADR line.

## Dependencies
No dependency. First REQ of the audit batch (`batch: maintainability-audit-2026-09-03`); nothing waits on it.

## Builder Guidance
Firm. Delete or mark retired; do not rewrite surrounding prose.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints 31 lines, five of them outside the ADRs (two topic indexes, the census, the imported spec, the prime).
**GREEN when:** The command prints only ADR lines, and the lock-in asserts the count of dead tokens cited outside `decisions/records/` is 0 (pinned first at today's 8 and lowered to 0 in this REQ).
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for dead-path-pointers-in-records.*
