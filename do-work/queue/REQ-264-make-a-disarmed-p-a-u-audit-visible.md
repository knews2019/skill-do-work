---
id: REQ-264
title: Make a disarmed P-A-U audit visible in qualify
status: pending
created_at: 2026-08-18T19:52:15Z
status_changed_at: 2026-08-18T20:55:14Z
user_request: UR-055
addendum_to: REQ-254
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-300, REQ-263]
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
---

# Make a Disarmed P-A-U Audit Visible in Qualify

## What

Both of qualify's UNIFY-gated FAIL branches key on a checked `[UNIFY]` box in the REQ file — so a REQ with **no** P-A-U section at all sails through Check 4 with the FAIL half silently disarmed. Every review-generated REQ from the previous session (REQ-250 through REQ-254) lacks the section, and REQ-254's own qualification "Passed" that way: its review re-ran the range armed and got FAILs (fixture TODO lines, qualify's own regex, a doc seam line — false positives of the protected class) with no override on the record. qualify should WARN when the REQ file carries no P-A-U section, so a disarmed audit is visible instead of silent.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-254 review, Important finding 2 (gate: rule-change). Created `pending-answers` per the generation-≥2 depth stop. The companion record-keeping (the armed-run override for REQ-254's own range) was written into REQ-254's archived trail by the orchestrator at integration; this REQ is the rule so the next disarmed audit cannot be silent. Worth deciding at build time whether review-created follow-up REQ templates should simply include the P-A-U block (the capture template already does), which removes the class at the source.

## Open Questions

- [ ] REQ-254's review found qualify's box audit silently disarmed for REQs without a P-A-U section. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — orchestrators should notice a missing section themselves.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.
