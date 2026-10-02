---
id: REQ-631
title: 'Recover''s takeover message says where a reset claim goes'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-02T21:57:20Z
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
required_lessons: ["skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#destructive-next-argv"]
write_set: ["skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go", "skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go"]
dispatch_at: 2026-10-02T21:58:14Z
builder_handback_at: 2026-10-02T21:58:14Z
integration_at: 2026-10-02T21:58:14Z
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

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18027 tokens, over the 2000 budget; `slugged: partial`). Matched: the next step a finding suggests. Narrowed at claim time to its `destructive-next-argv` family, now in `required_lessons`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: tighten the existing stop-reason assertion first (must say "returns it to the queue", must not say "as pending"), see it fail, then reword the one string to match what the recover transition does (status may be pending, pending-answers or blocked; route, write_set and generated sections are removed).
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: only the two declared files changed; next_argv and the finding code untouched.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: `git diff --stat` 2 files, 7 insertions, 1 deletion; gofmt clean, go vet clean; `go test -count=1 ./internal/lifecycleadvance/` ok in 26.6s (< 30s); no other shipped file quotes the old wording (grep).

## Full Context
See `do-work/user-requests/UR-134/input.md` for complete verbatim input.

*Source: ok, do board intros and the recover message*

---

## Triage

**Route: A** - Simple

**Reasoning:** One string in a named file and the test that pins it; the reset's real outcomes are already known from `recoveredStatus` and the recover transition.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modified)

**What was done:** The takeover finding's stop reason now says `recover --take-over REQ-NNN` "resets it (returns it to the queue and strips its routing and orchestrator sections)" instead of "requeues it as pending", which named one of the three statuses the reset can leave. The existing test now requires the new phrase and rejects "as pending". Merge range 6e3d4b81..25803a16 (builder commit 2596dbd2, merge 25803a16).

## Qualification

**Diff range:** 6e3d4b81..25803a16 (builder commit 2596dbd2, merge 25803a16)
**Gate records:** qualify satisfied. Route A, so no Scope comparison.
**Warnings judged:** none.
**Orchestrator read of the diff:** one string and one assertion. The new wording matches the recover transition in `requeststate/state_apply.go` (status pending, pending-answers or kept blocked; `route` deleted, `write_set` deleted when Scope exists, generated sections stripped). `next_argv` and the finding code are unchanged, so there is no behaviour change.
**P-A-U honesty:** the orchestrator played the builder role on branch `worktree-agent-REQ-631-recover-message`; APPLY cross-checked against `git diff --stat 6e3d4b81..25803a16` (two files).
