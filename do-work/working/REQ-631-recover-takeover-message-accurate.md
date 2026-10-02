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
review_at: 2026-10-02T22:01:33Z
heavy_verified_at: 2026-10-02T22:06:47Z
heavy_verified_revision: 9eed9df5b377811c8b6a008173e1f2e2c3209a93
claimed_at: 2026-10-02T21:57:11Z
commit: 25803a160d4590615fa3f0343ddb20f8809b8005
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

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at d9cf0ff3 (merge 25803a16 plus REQ trail)
**Result:** ✓ All passing — exit 0, gate wall 136s; do-work-cli fast tests 867, slowest file 21.64s < 30s; queue-kanban fast tests 413, slowest file 20.60s < 30s. Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-02-215500/helpers/probe-631.sh` (`go test -count=1 -run TestRecover ./internal/lifecycleadvance/`) → exit 0 (advance probe record satisfied). Whole package ok in 26.6s on the builder branch.

**Red-green validation:** traced to `## Red-Green Proof`:
- TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim (tightened): ✗ `stop reason misstates where the reset sends the claim: "... resets it (requeues it as pending and strips its orchestrator sections)"` → ✓

**Existing tests updated (cross-REQ impact):**
- `recovery_commands_test.go` TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim (from REQ-629): now also requires "returns it to the queue" and rejects "as pending" — intentional

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 6e3d4b81..25803a16
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — do-work-cli source changed
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — do-work-cli source changed
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — do-work-cli source changed

*Verified by work action*

## Review

**Overall: 97%** | 2026-10-02T22:01:33Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None
**Nit findings:**
- `skills/do-work/actions/work-reference.md:310` (Crash Recovery) still says "the reset requeues the claim and strips its orchestrator sections". It does not say "as pending", so it stays true. The new CLI wording also names routing, and this line, `restart-with-parallel-handoff.md:66` and `docs/work-guide.md:130` do not. That is an omission, not a contradiction. — impact-negligible → report only

**Acceptance:** Pass — the new wording matches the recover transition. `state_plan.go:62` moves every recovered claim to `do-work/queue/`, blocked ones included, so "returns it to the queue" holds for all three outcomes. `state_apply.go:579-605` sets `pending` or `pending-answers`, or keeps `blocked`. It always deletes `route`, deletes `write_set` when a Scope section exists, then runs `stripGeneratedRecoverySections`. The code's own comment calls this "the routing decision", so "routing" is a fair word. `next_argv` and the finding code did not change. `go test -count=1 -run TestRecover ./internal/lifecycleadvance/` passed (ok, 7.6s). The REQ is tdd: true, and its Testing section shows the tightened assertion failing on the old "requeues it as pending" text and passing after the change. All three P-A-U boxes are checked.
**Restatement sweep:** I grepped `skills/` for descriptions of what `--take-over` does: `work-reference.md:310`, `restart-with-parallel-handoff.md:66`, `docs/work-guide.md:130,165`, `lessons-do-work-cli.md:81` and `CHANGELOG.md:15-17`. None says the status is always `pending`, so none disagrees with the new wording. The CHANGELOG line describes the old behavior as history.
**Suggested testing:** 0 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*

## Discovered Tasks

From the review, impact-stamped per review-work Step 10; nothing was queued.

- N1: `work-reference.md` Crash Recovery says the takeover "requeues the claim" without mentioning that routing is removed; an omission, not a contradiction. — impact-negligible → report only

## Orientation

Now `recover`'s takeover warning says the reset returns a claim to the queue and strips its routing and orchestrator sections, without naming a single status; lives in the do-work-cli recovery command. Not a map change; `prime_files` is empty and no prime restates the message.

## Heavy Verification Plan

- Base revision: 6e3d4b811e149a128cb6b5ba2e2cde6604bc0502
- Target revision: 25803a160d4590615fa3f0343ddb20f8809b8005 (landed in `commit:`)
- do-work-cli-integrations — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — do-work-cli source changed
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed
- updater — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — do-work-cli source changed
- installer — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — do-work-cli source changed

## Heavy Verification Result

- Target revision: 25803a160d4590615fa3f0343ddb20f8809b8005
- Execution revision: 9eed9df5b377811c8b6a008173e1f2e2c3209a93 (detached drain checkout `.git/work-run-2026-10-02-215500/drain-head`, one run for REQ-630 and REQ-631, QUEUE_KANBAN_BROWSER set to Google Chrome)
- staged-skills: exit 0, executed, 43s
- do-work-cli-integrations: exit 0, executed, 73s
- updater: exit 0, executed, 66s
- installer: exit 0, executed, 29s

Green: every selected lane present, exit 0, none skipped, none reused.
