---
id: REQ-518
title: '[impact-rule-change] Run the full gate once per REQ'
status: pending
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-519, REQ-520, REQ-521, REQ-522, REQ-523]
batch: cheap-maintainer-gate
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-02T21:34:07Z
  basis:
    - Route B
    - 4-file write set
    - 1 new files
    - 2 subsystems involved
    - 7 acceptance criteria
    - persistence changes
    - cross-route regression gates
status_changed_at: 2026-09-02T23:26:22Z
---

# Run the Full Gate Once per REQ

## What

Stop running the identical canonical repository gate twice per REQ. Record the revision hash after every green gate run; before dispatch, when `HEAD` equals that recorded revision, take the baseline as green without running the gate. The Step 6.5 run after implementation stays mandatory.

The fold-first scan found no pending or pending-answers REQ that owns this: REQ-469 (Replace the unrelated canonical-gate hold with a blocked set-aside) and REQ-470/REQ-471 change what happens when the gate is red, not how often it runs.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add Git-private, exact-argv green-gate evidence behind typed CLI handlers; consume it only in the baseline lane, keep the final gate direct, and prove behavior through the public CLI before changing implementation.
- [x] **[APPLY]:** Added the typed gate-evidence CLI, Git-private persistence and history proof, pipeline consumers, and existing contract mutations within the declared 12-file scope.
- [x] **[UNIFY]:** Reviewed all 12 source/test/prose files, confirmed `gofmt` and `git diff --check`, kept contract regressions at 8,440 lines, ran focused/race/vet/launcher/contract checks, and found no debug artifacts.

## Why

`actions/work.md` runs the same gate argv at the pre-dispatch baseline (line 402) and again at Step 6.5 (line 531). The gate takes 6.5 minutes, and 167 REQ commits landed in 14 days: about 36 hours of gate wall-clock, half of it a re-run of a revision that was already proven green. Capture-time answer Q3: skip the baseline when `HEAD` is the last green revision.

## Context

The baseline exists for attribution: a red Step 6.5 result is compared to the pre-dispatch fingerprint to decide whether the failure is the REQ's own (`actions/work-reference.md` → Repository Gate Deferral and Resumption). A recorded green revision carries the same information: if the baseline revision was green and Step 6.5 is red, the diff is the cause. The recorded revision must be written only by a green run of the exact gate argv, and read only when it equals `HEAD` exactly.

## Detailed Requirements

- Mechanics in the Go CLI, not in prose: one command records the green revision for a gate argv after a zero exit; one command answers whether `HEAD` matches the recorded green revision for that argv. Where the record lives is the builder's choice; it must survive a session restart and must not be a hand-edited pipeline field.
- `actions/work.md` Step 5.75 and `actions/work-reference.md` → Session state and baseline: when the check reports a match, save it as the green baseline and dispatch; otherwise run the gate as today. A launch failure of the check stops safely, as a gate launch failure does today.
- Every green gate run in the pipeline (baseline or Step 6.5) records its revision, so the next REQ in the same run skips its baseline.
- Attribution at Step 6.5 and the deferral lifecycle work unchanged from a recorded baseline; the fingerprint procedure is untouched.
- A fingerprint or record that names a different argv, a different repository, or a revision not in `HEAD` ancestry never counts as a match.
- The record points at the revision after the gate's own log commit (REQ-523, Log and commit every maintainer gate run): the gate records its green revision after it has committed its log, and a commit whose only paths are under `_dev/gate-runs/` never invalidates a recorded green revision. Without this the log commit moves `HEAD` off the record after every run and the skip never fires. (verify-requests, 2026-09-03)
- The self-referential refusal invariant from REQ-514 applies to any new finding.

## Constraints

- Never waive the gate: a REQ still needs a green Step 6.5 run to complete.
- No new sentence predicates in `_dev/tests/contract-regressions.sh`; adjust the existing baseline predicates it already pins, and cover the behavior in a Go test.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

Roots this batch: REQ-519 (Path-scoped fast lane for the maintainer gate) depends on it because both change the same gate lanes in `work.md` and the maintainer asked for this first. Related to REQ-469 through REQ-471 by write-set overlap only.

## Builder Guidance

Certainty level: Firm on the rule (one full gate run per REQ when the revision is already green), latitude on the record's location and command names. Read the CLI prime first.

## Red-Green Proof

**RED prompt/case:** Run `do-work run REQ-NNN` on a revision whose gate was just green, and count gate invocations before the builder is dispatched.
**Why RED now:** The pipeline runs the full gate again at the baseline although `HEAD` has not moved since the green run; there is no record to consult.
**GREEN when:** A Go test records a green revision, sets `HEAD` to it, and asserts the check reports a match; a second test moves `HEAD` and asserts no match; a third adds one commit touching only `_dev/gate-runs/` on top of the record and asserts the match still holds. The work action reads that check and a run on an already-green revision performs exactly one full gate run, at Step 6.5.
**Validation:** User confirmed (capture-time answer Q3, 2026-09-03).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3636 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing pipeline fields and downstream gate readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on adding exact argv command blocks to shipped action prose.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on structured evidence projection and Git-private atomic persistence.

## Addendum (2026-09-03)

User added:

> ```text
> what does a no-op ticket take 24+ minutes to run?
> [board capture: REQ-533 Repair CLAUDE.md write-surface count contract, claimed 24m 33s, impact-critical, route C]
> this is what we need to resolve
> ```

Observed on REQ-533, an already-green repository-gate repair with an empty project diff: the full gate ran three times, about 7 minutes each. Once as the pre-build proof, once more after a peer commit moved `HEAD` and invalidated the recorded green revision, and once more because the no-op review contract says "rerun the JSON-array gate argv directly". The `gateevidence` package this REQ added in 3eae8110 already records a green revision; the no-op branch of the contract is the one consumer that still insists on re-running.

- For a `repository_gate_repair: true` REQ whose pre-build gate exits 0 (`actions/work-reference.md` → Already-green repair no-op completion), the pre-build run is the only gate run. It records its green revision through the gate-evidence command; qualification and independent review verify that record (exact argv, exit 0, the recorded revision, the expected fingerprint, an empty project diff) and never run the gate again. Rewrite the "rerun the JSON argv directly" sentences in `actions/work.md` Step 6.5 and Step 7, `actions/work-reference.md` § Already-green repair no-op completion, and `actions/review-work.md` Step 1 to say so; delete or rewrite their pins in `_dev/tests/contract-regressions.sh` without the file growing.
- A no-op repair has no diff, so a `HEAD` move caused by someone else's commit cannot make it the cause of a red gate; its recorded evidence stays valid for its own claim revision. Ordinary REQs keep the existing rule (a moved `HEAD` re-runs).
- Done means a no-op repair completes in one gate run plus bookkeeping: under ten minutes wall clock on the maintainer's machine, recorded in the Testing section.
- State: the Go side landed in 3eae8110 (`internal/gateevidence`, CLI wiring, result model) before the claim was recovered back to the queue; the builder starts from that code, not from scratch.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A1 of the analysis report's improvements, captured by UR-100.*

---

