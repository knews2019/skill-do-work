---
id: REQ-518
title: '[impact-rule-change] Run the full gate once per REQ'
status: claimed
route: C
planning_at: 2026-09-02T21:54:06Z
dispatch_at: 2026-09-02T22:04:18Z
builder_handback_at: 2026-09-02T22:29:30Z
integration_at: 2026-09-02T22:29:30Z
review_at: 2026-09-02T22:49:59Z
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
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go, skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go, skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go, skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go, skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go, skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go, skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/tests/contract-regressions.sh]
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
claimed_at: 2026-09-02T21:47:31Z
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

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A1 of the analysis report's improvements, captured by UR-100.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes persistent gate evidence, adds CLI behavior, and rewrites the work pipeline's attribution path across multiple subsystems with TDD and cross-route regression requirements.

**Planning:** Required

## Plan

1. Add a black-box CLI test first for a record/check lifecycle: exact `HEAD` matches, a project-changing commit invalidates it, and an intervening commit confined to `_dev/gate-runs/` preserves the match. Capture the initial unknown-command assertion failure as RED evidence.
2. Add an `internal/gateevidence` package and public `check-green-gate` / `record-green-gate` handlers. Key Git-private, restart-durable evidence by the canonical JSON encoding of the exact gate argv; bind it to the repository's canonical Git identity and full commit object.
3. Make checks read-only and return typed identity, provenance, state, match basis, stored/current revisions, and the exact baseline revision. Prove log-only ancestry commit-by-commit with NUL-safe Git output; all other ancestry/path shapes are non-matches, while corrupt or unsafe evidence fails closed.
4. Make recording accept only a supplied direct exit status of zero, publish the private record atomically, and ensure every refusal points to a distinct safe next action rather than naming `record-green-gate` as its own fix.
5. Extend the shared result model and text/JSON renderers, then register both commands only at `cmd/do-work-cli/main.go`.
6. Add focused package tests for repository/argv identity, missing and malformed records, ancestry divergence, mixed log/project history, non-regular targets, atomic replacement, zero-effect refusal, and result-rendering parity.
7. Update the CLI prime as a low-noise ownership/index pointer.
8. Update `actions/work.md` Step 5.75 to check typed evidence before a baseline run, fall back to the existing direct baseline on a valid non-match, stop on a failed check, and record every direct green baseline. Keep Step 6.5's direct gate mandatory and record only after its zero exit.
9. Update `actions/work-reference.md` once with the canonical check/record and `baseline_revision` semantics while preserving fingerprinting, attribution, and deferral behavior.
10. Repurpose existing `_dev/tests/contract-regressions.sh` predicates to pin delegation at the two action seams without adding a new sentence predicate or increasing the file's captured 8,417-line ceiling.

The command result must preserve per-record identity (`repository_identity`, exact gate argv and digest, record path), provenance (`direct_zero_exit`/`persisted_green_run`, exit status and revisions), state/outcome (`matches`, `match_basis`, `baseline_revision`, stable state), and canonical findings/changes/recovery argv from one observation.

Testing proceeds RED → GREEN → REFACTOR through the public CLI seam, then focused package tests, race tests, `go vet`, the full CLI module, Go 1.25 compatibility, contract regressions, and the unpiped repository gate at integration.

**Plan validation warning:** The plan has 10 ordered tasks, above the five-task review threshold. They remain one cohesive command-and-consumer contract, but the builder must keep the change inside the declared command, result-model, action, reference, prime, and existing contract-test seams.

*Generated by Plan agent*

## Exploration

- `cmd/do-work-cli/main.go` is the sole registration seam; handlers return typed results and leave rendering/exit behavior to `internal/commandruntime` and `internal/resultmodel`.
- `CommandResult` already uses optional typed domain projections. Gate evidence should follow that pattern, normalize argv to a non-null array, and render text/JSON from the same object.
- Reuse `internal/atomicfile.CreateExclusive` and `ReplaceExisting`; the latter supplies the required regular-file, identity, mode, and atomic-replacement guards.
- Store repository-wide evidence below the canonical Git common directory, not `git rev-parse --git-path`, which is per-worktree in linked worktrees.
- Prove transparent history commit-by-commit: ancestry first, then NUL-delimited `diff-tree` paths for every intervening commit. Only commits whose paths are all under `_dev/gate-runs/` are transparent.
- The CLI can preserve exact argv after `--`; it should not wrap or re-run the gate. The action remains the owner of direct execution and supplies zero-exit evidence to the recorder.
- Existing repository-gate sections in `actions/work.md`, `actions/work-reference.md`, and `contract-regressions.sh` provide the two narrow consumer seams. Fingerprinting, deferral, repair, and late-attribution branches stay unchanged.
- The current contract file is already 8,440 lines, above the captured 8,417 figure; this REQ will stay line-neutral or shrink it, while REQ-519 owns the absolute ratchet.
- Current finalization/release commits after Step 6.5 remain non-log invalidators. REQ-519 owns moving the full release gate; this REQ must not broaden the transparent path class.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go` (new) — public CLI RED/GREEN lifecycle coverage
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — register gate-evidence handlers
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` (new) — Git-private record and history-validation mechanics
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` (new) — record and matching behavior tests
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` (new) — typed public command handlers
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go` (new) — command/refusal/result tests
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — gate-evidence typed projection and text rendering
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — text/JSON parity tests
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — ownership and verification index
- `skills/do-work/actions/work.md` (modify) — cached baseline check and post-green recording calls
- `skills/do-work/actions/work-reference.md` (modify) — canonical baseline evidence semantics
- `_dev/tests/contract-regressions.sh` (modify) — line-neutral delegation seam assertions

**Files I will NOT touch:** `internal/corehelpers`, `internal/atomicfile`, request lifecycle/finalization packages, gate scripts, other action sections, or release/version files during implementation.

**Acceptance criteria (restated from REQ):**
- [ ] A green run can durably record exact gate argv, repository identity, and revision outside hand-edited pipeline state.
- [ ] A check matches only the same repository/argv at exact `HEAD` or across commits confined to `_dev/gate-runs/`.
- [ ] Missing, different, divergent, corrupt, or non-log-descendant evidence never counts as a match; unsafe/launch failures stop safely.
- [ ] Step 5.75 skips only an already-proven identical baseline and otherwise runs the existing direct gate.
- [ ] Every direct green baseline and Step 6.5 gate records evidence; Step 6.5 itself remains mandatory.
- [ ] Existing fingerprint attribution and repository-gate deferral behavior remain unchanged.
- [ ] The public CLI behavior is proven RED then GREEN, including exact match, moved project `HEAD`, and log-only descendant cases.
- [ ] Existing contract predicates are adjusted without adding a sentence predicate or increasing `_dev/tests/contract-regressions.sh` beyond its current 8,440 lines.

## Pre-Flight

**Git:** ⚠ One pre-existing untracked project artifact outside `do-work/` — preserve and exclude from this REQ (`ai-reports/2026-09-03_0010_maintainer-verify-gate-cost-analysis/index.html`). Two unrelated REQ-516 commits also landed after this REQ was claimed; they are part of the baseline, not this implementation.
**Tests baseline:** ⚠ `go test -count=1 ./...` launched and failed in pre-existing `internal/corehelpers` command-rendering/inventory tests; details are recorded in `do-work/working/baseline-failures.txt`.
**Dependencies:** ✓ Go module dependencies available; all other packages completed before the recorded `corehelpers` failures.
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` passed directly at `bab2198d59e891ec23d98a209c3c03187bc1741d`, establishing the green pre-build attribution baseline.

*Checked by work action*

## Decisions

- **D-01:** Store evidence under the canonical Git common directory so the integration checkout and linked worktrees share one repository-bound record.
- **D-02:** Treat only intervening commits whose every changed path is under `_dev/gate-runs/` as transparent; later lifecycle, release, version, and changelog commits remain invalidators until REQ-519 changes the integration boundary.
- **D-03:** Keep gate execution action-owned and accept an explicitly reported direct zero exit in `record-green-gate`; wrapping the gate would break the existing direct-output and fingerprint-attribution seam.
- **D-04:** Prescribe the shipped skill-root launcher with explicit repository root and JSON format in both action owners, avoiding PATH and working-directory assumptions.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added repository- and exact-argv-bound green-gate evidence with atomic Git-private persistence, typed check/record CLI results, and commit-by-commit log-only ancestry validation. The work pipeline now skips only a proven identical pre-build gate, records every direct green run, and still requires the direct Step 6.5 gate; existing line-neutral contract tests pin both consumer seams.

## Qualification

Passed — 12 files verified, 8 acceptance criteria traced, and P-A-U confirmed. The qualifier's four static-reference warnings are expected exceptions: three are Go test files discovered by the Go harness, and `gate_commands.go` contributes package-level handlers registered through `gateevidence.Handlers()` in the public CLI entry point. All modifications are substantive and within Scope; the unrelated AI report remains excluded.

## Testing

**Tests run:** `go test ./cmd/do-work-cli ./internal/gateevidence ./internal/resultmodel -count=1`; `go test -race -count=1 ./internal/gateevidence`; `go vet ./...`; `go test -count=1 ./...`; `bash _dev/tests/do-work-cli-go125-compatibility.sh`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/maintainer-verify.sh`; public `record-green-gate` invocation through the shipped launcher.
**Result:** ✓ All passing. The mandatory final repository gate passed, and `record-green-gate` returned typed success with the exact gate argv and Git-private record path.

**Red-green validation:**
- `TestGreenGateEvidenceLifecycle` in `cmd/do-work-cli/gate_evidence_integration_test.go`: the test file was the first source change; against pre-implementation commit `bab2198d`, it fails at the public seam with `UNKNOWN-COMMAND` for `record-green-gate` → ✓ passes against the implementation, covering exact match, project-commit invalidation, and log-only descendant matching. The original builder did not retain its RED transcript, so the work action reproduced the same assertion-level RED in an isolated detached worktree before accepting this evidence.

**New tests added:**
- Public CLI gate-evidence lifecycle and linked-worktree sharing
- Gate record identity, ancestry, invalidation, malformed/unsafe targets, permissions, and replacement
- Handler argv parsing, non-green refusal/non-mutation, and text/JSON result parity

*Verified by work action*

## Review

**Overall: 82%** | 2026-09-02T22:49:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 95% |
| Test Adequacy | 93% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The current Step 9 lifecycle/release commit occurs after the recorded Step 6.5 gate and therefore invalidates the record before the next REQ; the evidence commands are correct, but the normal two-REQ run reaches one-full-gate-per-REQ only after the integration boundary moves. — impact-user-visible → already owned by dependent REQ-519; its acceptance test must count full-gate invocations across two consecutive REQs.

**Minor findings:** 0
**Acceptance:** Partial — the evidence mechanism and action seams work end to end, while REQ-519 owns moving the full gate after lifecycle/release integration so the next REQ can reuse it.
**Suggested testing:** 2 items — after REQ-519, count full-gate invocations across two consecutive ordinary REQs; repeat with REQ-523's gate-log commit and a session restart.
**Follow-ups created:** None; **sweeps appended to:** None. Existing REQ-519 already owns the same integration-boundary root cause and is dependency-gated on this REQ.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Keeping gate execution action-owned while the Go CLI owns exact-argv, repository-bound evidence preserved the existing output/fingerprint seam.
- Git-common-directory storage plus commit-by-commit path validation gives linked worktrees one durable cache without broadening transparent history.

**What didn't:**
- The first builder prescribed a bare `do-work-cli`; the action contract needs the shipped skill-root launcher and explicit repository root at every executable seam.
- The original RED transcript was not retained. Reproducing the public test against the pre-implementation commit proved the transition, but future TDD builders must save contemporaneous RED output in their handback.

**Worth knowing:** Step 6.5 still precedes lifecycle/release finalization, so its record is intentionally invalidated by the successful REQ commit. REQ-519 owns moving the full release gate and proving the two-REQ invocation count; do not weaken the matcher by treating lifecycle/version/changelog commits as transparent.

## Orientation

[MAP CHANGED] The do-work CLI now owns durable green-gate evidence in a dedicated `internal/gateevidence` subsystem, and the work pipeline consumes its typed result to avoid a duplicate pre-build full gate only when repository, argv, and revision history match exactly. The final per-REQ gate remains mandatory; REQ-519 completes the batch-level integration boundary needed for the next REQ to reuse the record. Both touched primes remain current: the CLI prime indexes the new owner, and the action prime's executable/cross-reference rules still resolve.
