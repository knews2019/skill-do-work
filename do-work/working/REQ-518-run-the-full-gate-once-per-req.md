---
id: REQ-518
title: '[impact-rule-change] Run the full gate once per REQ'
status: claimed
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/tools/do-work-cli/lessons-do-work-cli.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/review-work.md
  - _dev/tests/contract-regressions.sh
  - do-work/lessons-index.md
  - VERSION
  - skills/do-work/VERSION
  - skills/do-work/actions/version.md
  - CHANGELOG.md
  - skills/do-work/CHANGELOG.md
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
claimed_at: 2026-09-03T12:42:27Z
---

# Run the Full Gate Once per REQ

## What

Stop running the identical canonical repository gate twice per REQ. Record the revision hash after every green gate run; before dispatch, when `HEAD` equals that recorded revision, take the baseline as green without running the gate. The Step 6.5 run after implementation stays mandatory.

The fold-first scan found no pending or pending-answers REQ that owns this: REQ-469 (Replace the unrelated canonical-gate hold with a blocked set-aside) and REQ-470/REQ-471 change what happens when the gate is red, not how often it runs.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Bind the already-landed gate-evidence record to an explicit revision so the no-op repair branch verifies the record instead of relaunching the gate.
- [x] **[APPLY]:** Added `--at-revision` to `check-green-gate` with a `target_revision` projection, then rewrote the four no-op re-run sentences and the two contract pins that held them.
- [x] **[UNIFY]:** Ran `gofmt -l` (clean), `git diff --check` (clean), `go vet`, the three focused packages, and the full contract suite; contract-regressions stayed at 8468 lines and no debug artifacts were left behind.

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

Every candidate row is `slugged: partial`, so the cheaper `path#family-slug` form is illegal for all three (Required Lessons Budget Contract); each is bare-or-nothing and each bare cost exceeds the 2000-token budget.

- `_dev/primes/lessons-action-files.md` — 3968 tokens, over budget. Matched on changing pipeline prose that downstream actions read.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget. Matched on prescribed command blocks in shipped action prose.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3124 tokens, over budget. Matched on structured evidence projection in do-work-cli internals.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The rule is firm and the Go side already landed in 3eae8110; what remains is finding the exact no-op re-run sentences across four action/prose files plus their contract pins, and extending the existing `gateevidence` package. Outcome is clear, the precise call sites need discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Where the rule already lives.** `internal/gateevidence` (landed in 3eae8110) records and checks a Git-private, exact-argv green revision. `CheckGreenGate` compares the record against `HEAD` only; `RecordGreenGate` writes the post-run `HEAD`. The pipeline already consumes both: `actions/work.md` line 402 (pre-dispatch baseline) and line 531 (Step 6.5 late attribution), plus `actions/work-reference.md` § Session state and baseline and § Late attribution.

**What the addendum is actually about.** The no-op repair branch is the one consumer that still re-runs. Four prose sites insist on it:

- `actions/work-reference.md:470` — "rerun the JSON argv directly at exit 0" (qualification) and "Independent review reruns/validates the gate".
- `actions/work.md:503` — "rerun the recorded argv at direct exit 0" (Step 6.3 qualification).
- `actions/work.md:531` — Step 6.5's mandatory final gate run, with no carve-out for the reviewed no-op branch.
- `actions/work.md:551` and `actions/review-work.md:46` — "verify the exact recorded argv exits 0 again" / "rerun the JSON-array gate argv directly".

**The missing mechanism.** `check-green-gate` binds the record to `HEAD`, so it cannot answer "was this argv green at *that* revision" once a peer commit moves `HEAD`. A no-op repair has no diff, so its evidence should stay valid at its own recorded revision. The gap is an explicit target revision, not a new record format.

**Contract pins.** `_dev/tests/contract-regressions.sh` pins the re-run sentence twice: the `review direct gate replay` predicate (line 3169) and its mutation tuple (line 3531). The `already_green_noop_defects` helper reads three owner sections (`work.md` Step 6.5, `review-work.md` Step 1, `work-reference.md` § Already-green repair no-op completion), so the qualification sentence in `work.md` Step 6.3 is prose-only with no pin.

**Constraint reading.** The REQ forbids *new* sentence predicates and caps the file's line count (now 8468 lines). Rewriting the two existing pins in place satisfies both.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` (modify) — target-revision-aware check
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` (modify) — behavior tests for the new target
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` (modify) — explicit target revision parsing
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go` (modify) — handler-level tests
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go` (modify) — public CLI proof
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — target revision projection
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — projection test
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — one-line index entry
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify) — lesson satellite entry
- `skills/do-work/actions/work.md` (modify) — Steps 6.3, 6.5, 7 no-op sentences
- `skills/do-work/actions/work-reference.md` (modify) — § Already-green repair no-op completion
- `skills/do-work/actions/review-work.md` (modify) — Step 1 no-op review subject
- `_dev/tests/contract-regressions.sh` (modify) — rewrite the two existing re-run pins in place
- `do-work/lessons-index.md` (modify) — refresh the do-work-cli satellite row
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` (modify) — release bookkeeping

**Files I will NOT touch:** `do-work/working/REQ-533-*.md`, `do-work/working/baseline.json`, anything under `do-work/queue/`, and the gate script `_dev/tests/maintainer-verify.sh` (REQ-523 owns gate-run logging, not this REQ).

**Acceptance criteria (restated from REQ):**
- [ ] One command records the green revision for a gate argv after a zero exit; one command answers whether a revision matches that record. Record survives a session restart and is not a hand-edited pipeline field.
- [ ] The baseline lane consumes the check; the Step 6.5 gate stays direct for ordinary REQs.
- [ ] A record naming a different argv, a different repository, or a revision outside the target's ancestry never counts as a match.
- [ ] A commit whose only paths are under `_dev/gate-runs/` never invalidates a recorded green revision.
- [ ] For a `repository_gate_repair: true` REQ whose pre-build gate exits 0, the pre-build run is the only gate run; qualification and independent review verify the record and never run the gate again.
- [ ] A `HEAD` move by an unrelated commit does not invalidate a no-op repair's recorded evidence for its own recorded revision; ordinary REQs keep the `HEAD`-bound rule.
- [ ] `_dev/tests/contract-regressions.sh` does not grow past 8468 lines and gains no new sentence predicate.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `do-work/lessons-index.md` (modified)
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` (modified)

**What was done:** `check-green-gate` gained an optional `--at-revision <revision>` that moves the match target off `HEAD` and onto a named commit, projected as `target_revision` in both output formats. `CheckGreenGate` is now a thin `HEAD` wrapper over `checkGreenGateAtRevision`, so identity, argv, ancestry and the `_dev/gate-runs/` tolerance are proven once and apply to both targets; an unresolvable target is a typed failure, not a miss. The already-green repair no-op branch was then rewritten in `work-reference.md`, `work.md` Steps 6.3, 6.5 and 7, and `review-work.md` Step 1 so the pre-build run is the branch's only gate run and every later authority verifies its record. The two contract pins that held the old re-run sentence were rewritten in place, leaving `_dev/tests/contract-regressions.sh` at 8468 lines with no new sentence predicate.

## Decisions

- **D-01**: Kept `resolveCommitRevision` and `commitResolvesExactly` as two functions rather than folding one into the other. DECIDE & STATE. They answer different questions with different failure semantics: a record naming a commit that no longer exists is a typed miss (`recorded_revision_missing`), while a caller-supplied target that does not resolve is unverifiable evidence and must surface as an error. Collapsing them would force one of those two outcomes to impersonate the other.
- **D-02**: Added `--at-revision` to the existing `check-green-gate` instead of minting a second command. DECIDE & STATE. The identity, argv, ancestry and `_dev/gate-runs/` tolerance rules are the same in both lanes; a second command would duplicate every one of them and let the two drift.
- **D-03**: Scoped Step 6.5's "this run is mandatory" clause to "every REQ but the single exception named below" rather than leaving the exception to contradict it three sentences later. DECIDE & STATE. The pinned substring the contract test mutates begins at `check-green-gate` is never consulted here, so the qualifier could be added ahead of it without touching the pin.

