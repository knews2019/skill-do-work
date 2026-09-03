---
id: REQ-541
title: 'The heavy gate is asked for through pending-heavy-testing, never run by the loop'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: [REQ-540]
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/tools/do-work-cli/internal/**/*.go
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/web/
---

# The Heavy Gate Is Asked for Through pending-heavy-testing, Never Run by the Loop

## What

Add one queue status, `pending-heavy-testing`, in the `pending-answers` family: queue-resident, never selected for work, shown by the board in Needs Input, routed by `do-work clarify`. In `actions/work.md` Step 6.5, after the fast gate is green: when this REQ's diff (`git diff --name-only <base>..HEAD`) matches any glob from `maintainer-verify.sh --heavy-surfaces`, the REQ lands its implementation commit as usual, then returns to `do-work/queue/` with `status: pending-heavy-testing` and one Open Question naming the exact `--heavy` command and the revision to run it at, and the loop continues to the next REQ. A REQ whose diff matches no heavy surface completes exactly as today. The loop never runs `--heavy` and never stops for it. The maintainer runs `--heavy` when they choose; answering the question yes completes and archives the REQ, answering no re-queues it as `pending` with the failing lane named.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"the factory should not stop on the heavy tier, those are being asked by placing the REQ into the need input column" (D4), refined on 2026-09-03: "where --heavy is required those tasks go into pending-testing status" and then "pending-testing is wrong, that is allowed, make it pending-heavy-testing". The ask is per REQ, not a standing sweep.

## Context

- The Needs Input column shows `pending-answers` and `blocked` queue residents; `pending-heavy-testing` joins that family so the board needs no new column, only a recognized status.
- The status is read in more than one program. Every reader of the status vocabulary must accept it in the same commit or the REQ is malformed to one of them: the do-work-cli selection, state-apply, schema normalization, doctor, and answer paths (illustrative: `internal/nextselection`, `internal/requeststate`, `internal/schemanormalization`, `internal/doctor`, `internal/publication/answer.go`), the board's classification in `queue-kanban/model.go` (the parser and the board move in lock-step per `prime-kanban-board.md`), and the status lists in `work.md`, `work-reference.md` (Request File Schema) and `capture-reference.md`. Find them by grepping for `pending-answers`; the list here is not a checklist.
- The glob list lives in the gate script (REQ-537) so no prose enumerates surfaces.
- Review and qualification consult recorded evidence since 0.266.8; any surviving sentence that lets them run the gate again is deleted here.

## Detailed Requirements

- One rule in `work-reference.md`, cited from Step 6.5; no new flag on `do-work run`.
- The REQ's implementation commit and its `commit:` field are recorded before the status flip, so the heavy run is the only thing outstanding; the Open Question reads "Run `bash _dev/tests/maintainer-verify.sh --heavy` at `<revision>`; did it exit 0?" with Recommended: yes.
- Selection skips `pending-heavy-testing` the way it skips `pending-answers`; a queue holding only such REQs reports that fact instead of treating them as malformed or as work.
- `do-work clarify` lists them; a yes answer completes and archives through the existing finalize path, a no answer re-queues as `pending` with the failing lane in the note. The loop never flips the status itself.
- The board classifies `pending-heavy-testing` into Needs Input with a label that says which command to run; the static snapshot and the live board agree.
- Delete any sentence in `work.md`, `work-reference.md`, or `review-work.md` that permits review or qualification to relaunch the gate.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** A REQ that edits `skills/do-work-board/tools/queue-kanban/web/board-filters.js` runs through `do-work run` after REQ-537 lands.
**Why RED now:** the loop runs the same gate for every diff and has no heavy ask; nothing appears in Needs Input, and a hand-edited `status: pending-heavy-testing` is rejected by the CLI and rendered as an error by the board.
**GREEN when:** that REQ lands its commit on the fast tier, appears in the board's Needs Input column as `pending-heavy-testing` with the `--heavy` command and revision in its Open Question, the loop proceeds to the next REQ, `do-work clarify` lists it and a yes answer archives it as completed, and a REQ editing only `actions/*.md` completes without the status.
**Validation:** User confirmed (D4, refined 2026-09-03)

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes action routing and a pipeline step contract.

## Addendum (2026-09-03)

User added (2026-09-03 21:29 and 21:34 local, in the test-budget session; 22:05 local, "update the batch REQs per A1-A3 via queued addenda" referring to the velocity report at `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, item A3):

> ```
> each test file should finish under 30 seconds (use the 80% value 20% effort principle until this is obtained)
>
> The rest of the test are accessible only when calling them with the --heavy parameter.
>
> the catch to call the --heavy parameter need to ask user for permission, and it should not block anything, meaning that where --heavy is required those tasks go into pending-testing status.
>
> Also because these tests tend to ballon, make sure to always measure the test duration, and adjust when the limits are reached.
> ```

> ```
> pending-testing is wrong, that is allowed, make it pending-heavy-testing
> ```

- Resolved conflict: the original What, Context, Detailed Requirements and Red-Green Proof specified one standing sweep REQ "Heavy gate runs requested" (`sweep_key: heavy-gate-requested`, appended to, "no new status") → the maintainer's decision is a per-REQ `pending-heavy-testing` status; the sections above were rewritten to that intent and the sweep shape is gone. UR-104's D4 wording ("one file, appended to") is superseded by this addendum; the rest of D4 (never a stopped run, the maintainer runs `--heavy` by hand) is unchanged.
- Title changed to name the status; the file name keeps its original slug.
- `write_set` and `prime_files` widened because the status is read by the CLI and the board, not only by the action files.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.
