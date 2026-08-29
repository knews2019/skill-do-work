---
id: REQ-410
title: 'Implement doctor, deterministic forensics, and metadata repairs'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-409]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Implement Doctor, Deterministic Forensics, and Metadata Repairs

## What
Create `doctor` and move deterministic forensic checks and safe metadata repairs into Go.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement `do-work-doctor` with deterministic forensics, blanked-record detection, timestamp audit/repair, collision checks, and remediation output.
- Reuse the typed finding/result schema so an LLM caller receives evidence and exact next commands without rescanning.
- Preserve the distinction between read-only diagnosis, provably safe repair, and explicitly destructive recovery.
- Support text/JSON, dry-run where repairing, optional commit, and shared Git guards.

## Constraints
- Missing or failed canonical tooling must stop the operation with actionable output rather than trigger prose mutation.

## Dependencies
Depends on REQ-409 (cleanup safety and repair classifications).

## Builder Guidance
Certainty level: Firm. Characterize every existing deterministic forensic utility before migration.

## Red-Green Proof
**RED prompt/case:** Run doctor over fixtures containing a blanked REQ, wrong timestamp, ID collision, and clean control.
**Why RED now:** These checks are split between shell helpers and LLM-directed action prose and do not share an actionable result schema.
**GREEN when:** Doctor classifies each fixture deterministically, performs only authorized safe repairs, and emits matching text/JSON evidence plus exact remediation/verification commands.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
