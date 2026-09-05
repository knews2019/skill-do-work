# Hand-back — REQ-575 (keep every lifecycle stamp: no transition deletes an existing `*_at` field)

## Branch

- Branch: `worktree-agent-REQ-575-append-only-stamps`
- Head commit: `afc30a9f` — `[REQ-575] keep every lifecycle stamp: no transition deletes an existing *_at field`
- Base: `09a13839` (main). One commit on the branch.

## File manifest

Source (modified):

- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` — added `setLifecycleStampWhenAbsent`, which states the append-only rule once and keys it on the `_at` suffix; `TransitionClaim` writes `claimed_at` through it; `TransitionRecover`'s ten-field strip loop is replaced by a single `DeleteField("route")`.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` — removed `parentDocument.DeleteField("claimed_at")` from the gate-deferral parent edit, with a comment naming the schema rule.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` — **outside the declared write set, see D-03** — `selectionClaimEvidence` now counts frontmatter `claimed_at` as live-claim evidence only while `status` is `claimed`.
- `skills/do-work/actions/work-reference.md` — Request File Schema gains the **Stamps are append-only** paragraph beside the Timestamp rule, and the four current-state exceptions are marked where each field is defined (`claimed_at` comment, phase-stamp comment, `blocked_at`, `status_changed_at`, `completed_at` cancel line, `heavy_verified_at`).

Tests (modified):

- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` — new `TestRecoverAndReclaimPreserveEveryLifecycleStamp` plus the `lifecycleStampsByFieldName` helper; `TestRecoverClaimCommitsCleanOwnershipTransferAndPreservesUnrelatedWork` flipped from "recover removed `claimed_at`/`planning_at`" to "recover kept them" (cross-REQ test break, intentional, commented in place).
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` — `TestDeferGateCreatePublishesOneAtomicDependencyLifecycle` now requires the parent's original `claimed_at` instead of an empty one.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` — **outside the declared write set, see D-03** — `TestClaimEvidenceVetoesEverySelectionModeBeforePolicyOrProbe`'s REQ-701 fixture is now `status: claimed` (the stamp alone no longer vetoes), and new `TestRecoveredRequestKeepingItsClaimStampStaysSelectable` pins the other half of the rule.

Nothing under `do-work/` was staged or committed. No release path (`CHANGELOG.md`, `VERSION`, mirrors) was touched.

## P-A-U

**[PLAN]** — Two deletion sites are named in the request, so no discovery was needed. Claim writes `claimed_at` only when the field is absent; recover deletes `route` (and `write_set` when a `## Scope` exists) and nothing else, so the suffix condition is satisfied by having no stamp list at all rather than by a longer list; gate deferral stops deleting the parent's `claimed_at`; one helper carries the rule and its exceptions in prose next to the transitions; the schema states the rule for the next writer. Test plan: one table-driven test that reads the stamp set out of the fixture by suffix and asserts survival across recover and re-claim, plus the defer-gate parent assertion.

**[APPLY]** — Coded as planned. One addition the plan did not have: `internal/nextselection` had to change too (D-03), because the selector treated any `claimed_at` as a live claim and a recovered request would otherwise never be selectable again.

**[UNIFY]** — `git diff --stat`: 7 files, +186/−14 before the comment polish; final tree committed as `afc30a9f`.

- `gofmt -l .` in `skills/do-work/tools/do-work-cli` — no output.
- `go vet ./...` — no output.
- Added lines scanned for `console.log`, `debugger`, `fmt.Print`, `TODO`, `FIXME`, `XXX` — no matches.
- `state_apply.go` — checked that recover still resets `status`, `status_changed_at`, `route`, `write_set` and the generated sections, and that the claim's `commit` / `heavy_verified_*` withdrawal is untouched.
- `defer_gate.go` — checked that only the one delete line went and the parent still moves to `queue/` with `gate_deferred`, `depends_on` and the history entry.
- `next_selection.go` — checked that checkpoint-writer evidence is unchanged and only the frontmatter branch gained the status condition.
- `work-reference.md` — checked that the new paragraph sits in the schema section beside the Timestamp rule and that every exception it names is also marked at its own field.
- Both test files — checked that each changed assertion states the behavior it now pins and names the request that changed it.

## Test evidence

RED (before the implementation, both sites):

- `go test ./internal/requeststate/ -run TestRecoverAndReclaimPreserveEveryLifecycleStamp` → FAIL. Nine subtests failed, one per stamp the fixture carried: `--- FAIL: .../recover_keeps_claimed_at: recover deleted claimed_at (was 2026-09-04T16:39:30Z)`, and the same line for `planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, `release_at`. The dumped post-recover document carried `status`, `id`, `title` and `status_changed_at` only.
- `go test ./internal/publication/ -run TestDeferGateCreatePublishesOneAtomicDependencyLifecycle` → FAIL at `defer_gate_test.go:39`, parent record printed with `ClaimedAt:""` where the fixture had `2026-09-02T01:00:00Z`.

GREEN (after the implementation):

- `go test ./internal/requeststate/ ./internal/publication/` → `ok requeststate 6.959s`, `ok publication 21.577s` (exit 0). Every subtest of `TestRecoverAndReclaimPreserveEveryLifecycleStamp` passes: all nine stamps survive byte for byte, `status_changed_at` is the recover instant `2026-09-04T21:02:00Z`, `route` and `write_set` are gone, and the re-claim at `2026-09-04T23:00:00Z` leaves `claimed_at: 2026-09-04T16:39:30Z`.
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./internal/requeststate/... ./internal/publication/...` → pass.
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./...` → pass, `wall=59s tests=760 slowest-file=internal/finalization/finalization_recovery_test.go:22.78s`. The **first** run of this command reported three files over the 30s per-file limit (`corehelpers/inventory_test.go` 42.63s, `finalization/finalization_recovery_test.go` 40.25s, `publication/defer_gate_test.go` 32.75s) with **no failing test**; all three are untouched by this change and the same files came in at 22.78s and below on the immediate re-run, so that was machine load from the parallel builders, not this diff. Worth knowing if the orchestrator sees the same overrun in the gate.
- `go test -count=1 ./...` in the module → exit 0, 30 packages `ok`.
- `bash _dev/tests/contracts/core-checks.sh` → exit 0, `core-checks contract probes passed.`
- `bash _dev/tests/shipped-package-reference-contract.sh` → exit 0, `PASS` (run because the change edits shipped action prose).
- `bash _dev/tests/select-simple-reqs-behavior.sh` → exit 0, `all probes passed` (run because the change edits the selector).

The repository gate (`_dev/tests/maintainer-verify.sh`) was not run, per the brief.

## Lesson evidence

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — present, read in the parts that touch this change (the `lifecycle-section-evidence`, `alternate-writer-contract-drift`, `closed-enumeration-for-a-condition` and `opaque-evidence-projection` entries, plus the timestamp-authority entry at REQ-434). Not read end to end: the request lists it as dropped for budget at 7300 tokens. `alternate-writer-contract-drift` (REQ-547/REQ-477: a stored-format contract changed in one writer without sweeping the others leaves behavior split) is why every writer and reader of `claimed_at` in the module was grepped, which is how the selector break was found before the module run.
- `_dev/primes/lessons-action-files.md` — present, read for the entries about schema and status contracts. The 0.221.0 entry ("when a rule changes what a field means or what a status can be, grep for every reader of that field and that status before calling the rule shipped") was applied directly: `claimed_at` was swept across Go non-test code, shipped Markdown, board sources and `_dev` scripts. Results of that sweep: `nextselection` was the only breaking reader; `doctor`, `finalization`, `requeststate`'s calibration check and the board read the stamp for reporting or spans and are correct as they stand; the board's stale-claim verify probe walks `board.Columns.Claimed` only, so a recovered `pending` request carrying a stamp raises no finding.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `_dev/primes/prime-action-files.md`, and the four crew files named in the brief plus `testing.md` — all present and read.
- No listed path was missing.

## Decisions

- **D-01 — `heavy_verified_at` and `blocked_at` keep being deleted. ESCALATE.** A literal reading of "any field whose name ends in `_at` is never deleted" would also stop `TransitionClaim` deleting `heavy_verified_at` (added one day ago by REQ-570, deleting the pending-heavy-testing status, as its D-01) and stop `TransitionUnblock` deleting `blocked_at` (documented in the schema since long before: "an unblock REMOVES `blocked_at`, so `status_changed_at` is the only trace of when that flip happened"). Both were left alone: they are withdrawn together with the state they describe (`heavy_verified_at` with the `commit` it verified, `blocked_at` with `blocked_by`), and deleting the first is an active guard that stops dependents building against work a remediation withdrew. The request's own Context section lists only two deletion sites, which suggests it was captured before REQ-570 landed. **Value:** the rule stays true to the goal (durable timing evidence) without reverting a one-day-old dependency-safety guard. **Risk:** if the maintainer wants zero exceptions, both call sites and their tests must change, plus the four exception notes in `work-reference.md`; that is a small, fully reversible edit. Also in this class and left alone: `completed_at` is re-stamped on the `failed` → `cancelled` path, which the schema already documented.
- **D-02 — the schema paragraph names its exceptions instead of claiming an absolute rule. DECIDE & STATE.** A schema sentence saying transitions never remove a stamp, while the same package removes three, would be worse than useless to the next writer. The paragraph states the condition (the suffix), says a stamp added later is covered without editing it, and points at the four fields that carry current state — each of which is also marked at its own definition, so a reader who lands on the field sees it.
- **D-03 — the selector was fixed, and it is outside the declared write set. ESCALATE (scope).** `internal/nextselection/next_selection.go` excluded any candidate carrying `claimed_at` with `ALREADY-CLAIMED`, so once recovery kept the stamp, every recovered request became permanently unselectable — `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` in `lifecycleadvance` caught it. The request's Builder Guidance says "If the builder finds a reader that breaks when a recovered request keeps its old stamps, fix the reader, not the writer", so I proceeded and am flagging it here as the brief requires. The fix is one condition: the stamp is live-claim evidence only while `status` is `claimed`; checkpoint-writer evidence is untouched and still vetoes on its own. **Value:** recovery keeps working, and the claim authority is the status, which is what it always meant. **Risk:** a request whose status is an unrecognized value (a legacy `in-progress`, say) while carrying a stamp is no longer vetoed by the stamp alone; its checkpoint entry still vetoes, and such a request is not a queue candidate in the first place. Reversible in one line if the orchestrator wants the veto keyed differently.
- **D-04 — recover deletes `route` with no loop. DECIDE & STATE.** Replacing the ten-name list with a single `DeleteField("route")` is what makes the suffix condition structural: there is no list left that a future stamp could be added to by mistake. `write_set` deletion stays exactly as it was, gated on the presence of a `## Scope` section.
- **D-05 — the stamp table in the new test is read out of the fixture, not written down. DECIDE & STATE.** `lifecycleStampsByFieldName` collects every frontmatter field ending in `_at`, so the assertion covers a stamp the schema gains later without anyone editing the test. Nine stamps are covered today.

## Discovered Tasks

- `skills/do-work/actions/work-reference.md` line 112 (the `write_set` schema line) points at "**Crash Recovery (Step 1)**, substep 1" for the clearing rule, but that section no longer has numbered substeps — the mechanics moved into the CLI's `recover` command. Stale cross-reference, impact-negligible → report only.
- `internal/requestmodel/request_model.go` still projects `heavy_verified_at` / `heavy_verified_revision` into the typed record with no remaining Go reader; already reported as F3 on REQ-570 (deleting the pending-heavy-testing status) and unchanged by this work — impact-negligible → report only.

## Integration seams

None. Everything this change needs is inside the commit; no line belongs in a file I could not write.

Two notes for the orchestrator, neither of them a line to apply:

- The behavior change is user-visible in the board: after a recover-and-re-claim, the calibration span and the drawer's Claimed row now measure from the **first** claim. The request states that as the intended reading.
- If a release entry is written for this, the honest summary is: lifecycle transitions no longer delete timestamps, so an interrupted request keeps the record of when its work actually started.
