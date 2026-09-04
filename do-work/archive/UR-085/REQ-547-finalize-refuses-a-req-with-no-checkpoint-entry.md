---
id: REQ-547
title: '[impact-rule-change] Stop finalize refusing a REQ that has no checkpoint entry'
status: completed
route: B
review_at: 2026-09-04T20:51:48Z
kb_status: pending
builder_handback_at: 2026-09-04T20:45:24Z
integration_at: 2026-09-04T20:45:24Z
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
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_req547_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-04T20:03:07Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - persistence changes
status_changed_at: 2026-09-04T20:43:02Z
claimed_at: 2026-09-04T20:43:09Z
completed_at: 2026-09-04T20:55:11Z
commit: 865dff21e142f6cef38520619811beb6c8768ac6
release_at: 2026-09-04T20:55:11Z
---

# Stop Finalize Refusing a REQ That Has No Checkpoint Entry

## What

`finalize` refuses any REQ with no entry in `do-work/CHECKPOINT.md`'s "In Progress
(interrupted)" section, and the refusal cannot be cleared in place. Make the lifecycle
preimage and postimage path sets agree so the checkpoint's presence or absence stops
deciding whether a REQ can be finalized at all, and give an operator a way out of an
already-written journal that refuses on replay.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---

## Handoff State (session stopped 2026-09-04T20:2xZ)

**Builder branch:** `worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry` at `b3d25c8`, pushed to origin. **Unmerged, but complete.**

The builder finished and wrote a full hand-back at `do-work/runs/work-2026-09-04-200249/REQ-547-handback.md`: file manifest, decisions, verification with revert-and-show-red evidence, discovered tasks. Nothing is missing from the builder side.

**Resume at the Step 6 hand-back merge.** Do not re-dispatch a builder. The remaining pipeline is: merge the branch, capture the merge range, qualify, run the repository gate, independent review, lessons, release judgment, finalize. Nothing after the builder has been done.

**What it delivered**, per the hand-back: planned postimages now emit exactly the plan's declared target paths and error on a declared target that cannot be projected, so a journal's preimage and postimage sets cannot name different paths; a terminal transition removes **every** checkpoint entry for the REQ it archives instead of only one carrying the manifest's writer label; and `recover-finalization --discard-journal REQ-NNN` is the bounded exit from a journal already written in a refusing state. Tests cover the absent-entry finalize, the foreign-writer-label clearing, the discard verb's four cases, and proof that genuinely disagreeing image sets still refuse.

**This fixes a defect this session actually hit.** Three finalization manifests here passed a `writer_label` that did not match the label `advance` wrote into the checkpoint at claim time. `finalize` returned typed success and archived each REQ but silently left the checkpoint entry behind, so the checkpoint claimed three REQs were in flight that were already archived — no refusal, no report. The builder's second change is what closes it: entry removal keys on the REQ, not on the writer label. Verify that specific behaviour during review rather than taking the hand-back's word for it.

## Triage

**Route: B** — resume the completed builder handoff; no new implementation dispatch.

## Plan

Planning not required. Restore the completed branch, derive postimages from declared targets, clear all claims for a terminal REQ, and provide a bounded pre-mutation journal discard.

## Exploration

The builder handback records the existing absent-checkpoint mitigation and the remaining foreign-writer silent skip. Inspected the merged plan/projector and discard paths; the count guard remains intact.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req547_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

**Acceptance criteria:** declared image sets agree; present/absent/unchanged checkpoints finalize; terminal claims clear regardless of writer; prepared untouched journals can be discarded; real corruption still refuses.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req547_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** Projected lifecycle images follow the declared target set. Terminal transitions clear every matching REQ claim. A bounded `recover-finalization --discard-journal REQ-NNN` handles untouched prepared journals while retaining corruption checks.

**Integration range:** `67c9b372..865dff21`.

## Decisions

Accepted builder D-01 through D-06 from the durable handback. The expanded command/types/tests/prime manifest is required by the explicit recovery-verb requirement; the original four-file declaration omitted those necessary surfaces. This supersedes the narrow captured write set.

## Integration Verification

Read every path in the implementation manifest against `67c9b372..865dff21`. The plan/projector enforce one target set; command/types/prime expose the bounded discard; both test files cover those changes without debug artifacts. `git diff --check` and the full unpiped `bash _dev/tests/maintainer-verify.sh` passed, including ShellCheck, gofmt, contracts, vet and both Go modules. Independent focused tests passed for finalization, requeststate and lifecycleadvance.

## Qualification

Passed — canonical advance returned satisfied qualify and scope-drift records for `67c9b372..865dff21`. The only warning is the new Go test file having no static reference; Go discovers `_test.go` functions automatically, independently confirmed by the focused package run. Every changed file is declared and each acceptance condition maps to production code and tests.

## Gate Repair

The staged-skills lane failed twice on the same obsolete assertion: `core runtime must resolve actions/work.md through sibling do-work-board`. The predicate also fails at pre-merge `67c9b372`: an earlier lifecycle migration removed the last informational board citation and delegates execution to core `advance`. Independent review confirmed there is no remaining board runtime dependency. Removed that single assertion and retained the actual sibling-consumer assertions and canonical advance contracts. Added `_dev/tests/staged-skills-contract.sh` to this integration's completion scope.

## Testing

**Tests run:** canonical focused-test gate passed `go test -count=1 ./internal/finalization ./internal/requeststate ./internal/lifecycleadvance`. Full maintainer gate passed after the builder merge.

**Red-green validation:** builder handback records isolated reversions: removing declared-target projection makes the path-set invariant tests fail; restoring writer-label filtering leaves the archived request's claim behind; removing discard dispatch makes all four command tests fail. Independent merged-tree acceptance reran those tests green. The originally captured absent-entry case was already mitigated on the base; the tests preserve it while pinning the remaining defects.

**Heavy verification:** the exact `67c9b372..865dff21` plan selected do-work-cli-integrations, staged-skills, updater and installer. Integrations (54s), updater (53s), installer (50s) passed unskipped. Staged-skills failed twice on the pre-existing assertion described above, then passed unskipped after its one-line correction (32s). Selected lanes were executed during this targeted completion while the other request's builder remediated; no heavy evidence was reused.

**Existing tests updated:** the cancel test now expects all matching REQ claims to clear while unrelated claims survive; the mode test keeps declared targets aligned with fixture mutations. Existing behavior coverage remains.

## Review

**Acceptance: Pass.** Independent reviewer checked `67c9b372..865dff21`, original request/input, every changed file, corruption refusal, absent/foreign-writer checkpoint cases, discard boundaries and cross-file restatements. Focused acceptance tests passed. No actionable regression found. Overall 98.75%; low risk. Separate review of the staged-skills failure confirmed the obsolete dependency on both base and merged tree and approved removing only that assertion.

The pre-existing recovery/set-aside findings on REQ-515 (continuing after one stuck request) remain outside this change; this review does not close them.

## Lessons Learned

**What worked:** use one declared path set for both lifecycle images; verify terminal state against the REQ identity, independent of the writer label.
**What did not:** independent image projections drifted, and the earlier workaround hid stale claims instead of clearing them. A literal citation-presence test also outlived the dependency it was intended to check.
**Worth knowing:** discard is deliberately limited to untouched prepared journals. It is not a way to bypass corruption or undo a committed transaction. The lesson satellite and index were updated together.

## Orientation

The finalization subsystem can complete requests without a checkpoint claim and clears claims labelled by another writer. Its recovery command now offers a bounded exit for untouched prepared journals. The touched prime describes that command and all referenced runtime paths exist.

**Release preparation check:** the gate and its retry rejected a lesson link written before its archive target existed. Deferred that exact lesson/index payload until the finalization transaction creates the archive; no implementation failure was involved.

**Final gate:** `bash _dev/tests/maintainer-verify.sh` passed after correcting the stale dependency assertion and deferring the archive link. All required checks and selected heavy lanes are green.
