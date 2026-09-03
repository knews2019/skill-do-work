---
id: REQ-547
title: '[impact-rule-change] Stop finalize refusing a REQ that has no checkpoint entry'
status: pending
priority: now
created_at: 2026-09-03T16:40:00Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-543, REQ-502]
write_set:
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go
---

# Stop Finalize Refusing a REQ That Has No Checkpoint Entry

## What

`finalize` refuses any REQ with no entry in `do-work/CHECKPOINT.md`'s "In Progress
(interrupted)" section, and the refusal cannot be cleared in place. Make the lifecycle
preimage and postimage path sets agree so the checkpoint's presence or absence stops
deciding whether a REQ can be finalized at all, and give an operator a way out of an
already-written journal that refuses on replay.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Hit for real while releasing REQ-543 on 2026-09-03. The mechanism, read straight out of
the refused journal (3 preimage paths against 4 postimage paths):

- `internal/requeststate/state_plan.go:398` — `planTargets` adds the checkpoint to
  `TargetPaths` only when its bytes would change.
- `internal/requeststate/state_apply.go:156-158` — `PlannedPostimages` emits a checkpoint
  postimage unconditionally.
- `internal/finalization/finalization_apply.go:296` — `imageSetState` compares the two sets
  and refuses with `FINALIZATION-LIFECYCLE-CONFLICT: journal image sets have different path
  counts`.

So a REQ whose checkpoint line is missing (or whose removal would be a no-op) plans one
fewer path than it applies, and finalization refuses.

The refusal itself is safe: phase `prepared`, no commit, no release mutation, `VERSION` and
both changelog mirrors untouched. The unrecoverable part is the replay —
`recover-finalization` re-reads the same journal into the same refusal, and there is no
discard verb. The only way past it was to move the Git-private journal aside by hand and
then add the checkpoint entry the Session Checkpoint template says a claimed REQ in
`do-work/working/` should already have carried.

This is a queue-drain blocker, not a cosmetic refusal: any REQ that reaches finalization
without a checkpoint line stops there and needs manual journal surgery to continue.

## Detailed Requirements

- The lifecycle preimage and postimage sets must name the same paths for the same
  transaction. Fix it on whichever side is wrong — the plan omitting an unchanged path, or
  the apply emitting a postimage for a path it never planned — rather than loosening
  `imageSetState`'s count check, which exists to catch genuine journal corruption.
- A checkpoint entry that is present, absent, or byte-unchanged must all finalize. Absence
  is the case that fails today.
- Provide a recovery path for a journal already written in a refusing state: either
  `recover-finalization` recognises this specific conflict and repairs the image set, or
  there is an explicit discard verb. Manual removal of a Git-private journal must stop
  being the only exit.
- Keep the refusal for real corruption. A journal whose sets genuinely disagree about a
  release or lifecycle path must still refuse, and a test must prove that still happens.

## Constraints

- Do not weaken `imageSetState` into a warning. The check caught a real inconsistency here;
  the inconsistency is the defect, not the check.
- Do not make the checkpoint a required precondition of `finalize` as the fix. That trades a
  refusal for a stricter refusal and leaves the same class of REQ stuck.
- Every new test names the specific failure it pins.

## Red-Green Proof
**RED prompt/case:** Claim a REQ, remove its line from `do-work/CHECKPOINT.md`'s "In Progress
(interrupted)" section, then run `finalize` for it.
**Why RED now:** refuses with `FINALIZATION-LIFECYCLE-CONFLICT: journal image sets have
different path counts`, and `recover-finalization` replays into the same refusal.
**GREEN when:** that sequence finalizes; a journal already refusing is recoverable through a
documented verb rather than by moving files; and a deliberately corrupted image set still
refuses.
**Validation:** Reproduced during the REQ-543 release, 2026-09-03.

## Notes

Fold-first scan run over all six pending `sweep: true` REQs. REQ-502's
`checkpoint-section-blind-line-editing` is a different defect (how a line is edited, not
whether its absence blocks finalization), so this is captured separately rather than folded
into it.
