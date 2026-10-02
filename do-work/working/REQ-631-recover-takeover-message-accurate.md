---
id: REQ-631
title: 'Recover''s takeover message says where a reset claim goes'
status: claimed
created_at: 2026-10-02T21:51:57Z
user_request: UR-134
domain: general
prime_files: []
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-630"]
batch: ur-133-leftovers
write_set: ["skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go", "skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go"]
claimed_at: 2026-10-02T21:57:11Z
---
# Recover's Takeover Message Says Where a Reset Claim Goes

## What
Make the stop reason of `recover`'s `RECOVERY-TAKEOVER-AVAILABLE` finding describe the takeover reset accurately: it returns the claim to the queue (as `pending`, as `pending-answers` when Open Questions has an unchecked item, or still `blocked`) and strips its orchestrator sections.

## Why
The REQ-629 review (finding M2) found the message, added in 0.305.64, says the reset "requeues it as pending", while `recoveredStatus` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` also yields `pending-answers` and keeps `blocked`.

## Detailed Requirements
- Reword the stop reason in `recovery_commands.go` (around line 103) so it does not claim the status is always `pending`. Short wording such as "returns it to the queue and strips its orchestrator sections" is enough; the handoff action already uses "returns it to the queue".
- Check what else the reset drops (the REQ-629 review said `route` and `write_set`) and keep the wording true without listing every field.
- Update the test in `recovery_commands_test.go` that pins the reset wording.
- Do not rename the finding code and do not change `next_argv`.

## Constraints
No behaviour change. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-630.

## Red-Green Proof
**RED prompt/case:** Run plain `recover` with a working claim: the stop reason says the takeover "requeues it as pending"; a claim with an unchecked Open Question would go to `pending-answers` instead.
**Why RED now:** The 0.305.64 wording names one outcome of three.
**GREEN when:** The test that pins the stop reason asserts the new wording, which no longer says "as pending", and passes.
**Validation:** Inferred during capture

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-134/input.md` for complete verbatim input.

*Source: ok, do board intros and the recover message*
