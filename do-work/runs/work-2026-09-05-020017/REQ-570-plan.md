# REQ-570 Plan — delete `pending-heavy-testing`; held requests stay `claimed`

Design is fixed by UR-114. The full reader inventory with line numbers is `REQ-570-exploration.md` in this run directory; this file does not repeat it.

## Findings the orchestrator must see before dispatch

- **F1 — dependency REQ-507 is `pending`, not terminal.** `do-work/queue/REQ-507-*.md` has `status: pending`, `commit: ad8bceb7`, and a `## Heavy Verification Result` (red lane). Its code and its Step 8/9 prose are on `main` (`advance` already returns the `finalize` phase in `advance_commands.go`, and `finalization_gate.go` composes `FinalizeBound`). REQ-570 was claimed by explicit naming, which bypasses `depends_on`. Nothing in this plan needs REQ-507's remaining remediation. Record this as a typed finding in the REQ, do not stop.
- **F2 — no queued or working request carries the old status.** The only hits under `do-work/` are REQ-571's filename and title, REQ-570 itself, and the generated `queue_state:` string in `CHECKPOINT.md`, which `advance --checkpoint` rewrites. The "stop and report" constraint is not triggered.
- **F3 — obstacle: a stale `commit:` widens readiness under the new rule.** Neither the `claim` transition nor the `recover-claim` transition strips `commit:` (`requeststate/state_apply.go` lines 516-524 and 526-549). Today `pending`+commit is not ready, so a stale commit is harmless. With the new rule (`claimed` + nonblank `commit:` is ready), re-claiming REQ-506 or REQ-507 makes their dependents source-ready at claim time, against the previous attempt's commit, before any new work lands. The same hole opens on a red drain: the request stays `claimed` with its red `commit:`. **Smallest resolution:** the `claim` transition deletes `commit`, `heavy_verified_at`, and `heavy_verified_revision` (three field names in the existing delete loop shape, one test), and the red branch of the drain prose deletes `commit:` and `## Heavy Verification Plan` from the record before remediation re-dispatch, using the same writer that wrote them at the hold. This adds `internal/requeststate/state_apply.go` and its test to the write set. It is an addition earned by a concrete case (two such records are in the queue now). If the orchestrator declines, record the decision and the risk in `## Decisions`.
- **F4 — no contract predicate names the status.** `grep -n 'heavy\|resume_phase\|Step 2.5' _dev/tests/contracts/core-checks.sh` returns nothing; `contract-regressions.sh` uses `heavy` only as its own tier name. The REQ's "delete the matching predicates in `core-checks.sh`" is satisfied by absence; remove `_dev/tests/contracts/core-checks.sh` from the REQ's `write_set` if it stays untouched and say so in the Implementation Summary.
- **F5 — who writes `heavy_verified_*` and `## Heavy Verification Result` after the `answer` mode is gone.** The drain writes them onto the claimed record the same way Step 7 stamps `review_at` and appends `## Review`: an action-owned edit of the record in `do-work/working/`, before the Step 8 manifest binds the request preimage. No CLI change. The runner's JSON (`run-heavy-verification`) supplies the values.

## Architectural decisions

- **D1 — held is a phase, detected by section plus ancestry.** A claimed request in `do-work/working/` with a `## Heavy Verification Plan` section and a `commit:` that is an ancestor of HEAD is held. No field, no status, no clock.
- **D2 — `DependencySourceReady(status, commit)` keeps its signature.** Only the second disjunct changes to `normalizedStatus == "claimed" && commit != ""`. `dependencygraph` and `duplicateStatusesSatisfied` already pass the node's commit; they need no code change.
- **D3 — the selector loses `ResumePhase` entirely.** A held request is in `working/`, which `Select` never considers (only `TreeSection == "queue"` records are candidates). Nothing replaces `matchingHeavyReviewPhase`; review already ran before the hold.
- **D4 — recovery preserves a held claim instead of releasing it.** The check sits after the authority check in `handleRecover` (`recovery_commands.go` line 92 onward), so an unauthorized session still gets `RECOVERY-TAKEOVER-AVAILABLE`. When authorized and held: no `requeststate` plan runs, `Decision` is `held for heavy lanes; claim preserved`, `Recovered` stays false, a new bool `HeldForHeavyLanes` (`json:"held_for_heavy_lanes"`) on `RecoveryClaimResult` is set, and an info finding `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` carries the section name and the commit as evidence with `NextArgv` `do-work-cli --format json advance REQ-NNN`. The runner argv cannot be derived by the CLI because the lane manifest path is action-owned. Held with a commit that is **not** an ancestor of HEAD falls through to ordinary recovery: the plan section alone is not proof of landed work.
- **D5 — `advance` needs no new phase.** A held record already has Review, Lessons Learned, and Orientation, so `classifyWorkingAdvance` returns `finalize`, which is exactly the green drain's next step. `advanceSections` ignores the unknown `Heavy Verification Plan` heading. Reuse `advanceSections` + `hasSection` from the same package in D4.
- **D6 — the hold and drain prose move to one named subsection between Step 7.5 and Step 8** (`### Step 7.7: Heavy-Test Hold and Drain`). Qualification and Testing Judgment keeps one sentence pointing there. This is a move of two existing paragraphs, not a rewrite; the floor agent needs the position to read in step order.
- **D7 — the ancestor check in recovery is a local three-line helper** (`git merge-base --is-ancestor <commit> HEAD` via `exec.Command`, the shape `finalization_journal.go` line 273 uses). `lifecycleadvance` already shells out to git in `checkpoint_commands.go`; do not import `nextselection` or `publication` for their unexported helpers.

## Files to modify, in compile-safe order

Run `go build ./... && go vet ./...` after every group.

1. **`internal/nextselection/next_selection.go`** — delete `ResumePhase: matchingHeavyReviewPhase(...)` from the `SelectionRecord` literal (line 315), delete `matchingHeavyReviewPhase` (322-337), delete the `case "pending-heavy-testing"` in `summarizeQueue` (407-408). Delete `selectionRevision` / `selectionRevisionIsAncestor` and the `time` import only if nothing else in the package uses them (grep first). **`next_selection_test.go`** — see Testing.
2. **`internal/resultmodel/result_model.go`** — delete `ResumePhase` (228) and `PendingHeavyTesting` (272); delete the two `resume phase` render lines (1193-1195); drop the `PendingHeavyTesting` term from the `queue:` condition and format string (1179-1182); fix the comment at 97-100 that cites publication's `HeavyLaneResult`.
3. **`internal/schemanormalization/schema_normalization.go`** — remove the enum value (30); rewrite `DependencySourceReady` per D2 and its doc comment (236-244).
4. **`internal/dependencygraph/`** — no production change. **`dependency_graph_test.go`** — see Testing.
5. **`internal/publication/answer.go`** — delete the `case "heavy-testing"` block (254-293), `validateHeavyTestingEvidence` (610-666), `appendHeavyVerificationResult` (684-697), and `resolveAnswerCommit` / `answerRevisionIsAncestor` if unused afterwards; change the `ANSWER-MODE-INVALID` message to `answer mode must be clarify, stakeholder, or verify-repair` (352). **`publication_types.go`** — delete `HeavyLaneResult`, `HeavyTestingEvidence` (81-95) and `AnswerManifest.HeavyTesting` (105). **`answer_test.go`** — see Testing.
6. **`internal/lifecycleadvance/advance_commands.go`** — remove the value from the skip case (110). **`checkpoint_commands.go`** — remove the count from the format string and lookup (129-131). **`recovery_commands.go`** — add the held branch per D4 and D7 with a `heldForHeavyLanes(repositoryRoot, request) bool` helper. **`resultmodel/result_model.go`** — add `HeldForHeavyLanes bool` to `RecoveryClaimResult`. **`recovery_commands_test.go`** — see Testing.
7. **(F3, if accepted) `internal/requeststate/state_apply.go`** — in `TransitionClaim` (516-524) delete `commit`, `heavy_verified_at`, `heavy_verified_revision`. Test in the requeststate test file that covers claim bytes.
8. **Prose** — `skills/do-work/actions/work.md`, `work-reference.md`, `clarify.md`, `cleanup.md`, `roadmap.md`, `restart-with-parallel-handoff.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (details below), plus lessons satellites and `do-work/lessons-index.md` rows.
9. **Release** — shipped files change, so finalization carries a version bump and a changelog entry per `_dev/primes/prime-releases.md` (title says what shipped; mirror byte-identical). No CHANGELOG history is rewritten.

## Testing approach (TDD: write RED first, run, then edit production)

Old tests that pin the deleted behavior and must go or be rewritten (cross-REQ test-break rule: intentional; name them in `## Testing`):

- `schema_normalization_test.go:41` warning string lists the value; `TestTerminalPredicatesKeepFailureAndCancellationDistinct` lines 162-167 pin the old disjunct.
- `dependency_graph_test.go`: `TestPendingHeavyDependencyIsSourceReadyUntilItReturnsToPending`, `TestPendingHeavyDependencyRequiresImplementationCommit`, and the `buildFixtureGraph` special case at line 24.
- `next_selection_test.go`: `TestPendingHeavyTestingIsCountedAndNeverSelected`, `TestMatchingHeavyEvidenceResumesAtReviewAndStaleEvidenceDoesNot`, `TestPendingHeavyTestingSourceAllowsDependentSelection`; `runNextSelectionGit` is used 13 times, keep it.
- `answer_test.go`: `TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure`, `TestHeavyTestingAnswerRejectsNegativeWallSeconds`, `TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane`, `TestHeavyTestingAnswerRejectsMismatchedEvidence`, helper `newHeavyAnswerFixture`; `runGitFixtureOutput` has other callers, keep it.

RED tests to write first, one per Red-Green Proof clause:

- **(a)** `schema_normalization_test.go`: `TestDependencySourceReadyAcceptsClaimedWithCommitOnly` — `("claimed","abc123")` true; `("claimed","")` false; `("pending","abc123")` false; `("pending-heavy-testing","abc123")` false; and `NormalizeField("status","pending-heavy-testing").IsRecognized` false. Update the line-41 expected warning string to the shorter list (RED until the enum changes).
- **(b)** `dependency_graph_test.go`: `TestClaimedDependencyWithCommitIsSourceReady` — extend `graphFixture` with a `commit` string and a tree section so a `claimed` fixture is written under `do-work/working/`; assert the dependent of `claimed`+commit `IsReady`, the dependent of `claimed` without commit is not, the dependent of `pending`+commit is not.
- **(c)** `next_selection_test.go`: `TestHeldClaimedSourceAllowsDependentSelection` — REQ-A `claimed` in `do-work/working/` with `commit:` and a `## Heavy Verification Plan` body, REQ-B `pending` with `depends_on: [REQ-A]`; `Select` returns exactly `[REQ-B]`. RED today: REQ-B is `DEPENDENCIES-UNMET`. `ResumePhase` removal is compile-checked; no decorative JSON-absence test.
- **(d)** `answer_test.go`: `TestAnswerRefusesHeavyTestingMode` — a `pending-answers` fixture with `Mode: "heavy-testing"` returns `Refusal.Code == "ANSWER-MODE-INVALID"`. RED today: the code is `ANSWER-HEAVY-STATUS`.
- **(e)** `recovery_commands_test.go`: `TestRecoverHoldsAClaimedRequestWithAHeavyVerificationPlanForTheDrain` — in-process `handleRecover(..., "--assume-sole-authority")` on a git fixture (shape of `recovery_set_aside_test.go`): REQ-740 `claimed` in `working/` with `commit:` = HEAD and a `## Heavy Verification Plan` section; REQ-741 `claimed` without the section. Assert REQ-740 stays in `working/`, its claim has `HeldForHeavyLanes` true and `Recovered` false, a finding with code `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` names it, and REQ-741 is recovered to `queue/`. Second case in the same test: REQ-742 with the section but a `commit:` from an abandoned branch is recovered to `queue/` (D4 fall-through).
- **(F3, if accepted)** requeststate: `TestClaimStripsPriorAttemptCommitAndHeavyEvidence` — a `pending` fixture carrying `commit:`, `heavy_verified_at`, `heavy_verified_revision` is claimed; the claimed bytes carry none of the three.

`advance_commands_test.go` and `checkpoint_commands_test.go` pin neither the skip case's status list nor the `queue_state` string; the final `grep` is their proof.

## Prose edits (one or two sentences each; judgment stays prose, mechanics point at the CLI)

**`actions/work.md`**
- Line 16 (Do NOT use): drop the held clause; keep `pending-answers` → `do-work clarify`.
- Line 91 (Status flow): reduce to `pending → claimed → completed / completed-with-issues / failed`; add that the heavy hold is a phase of `claimed`, marked by `## Heavy Verification Plan` and `commit:`, per the next paragraph's section-tracking rule.
- Lines 97 (Special statuses): delete the `pending-heavy-testing` bullet.
- Line 343 (Qualification and Testing Judgment, last sentence): "Plan affected heavy lanes with the repository's typed planner; selected lanes are recorded in the Testing section and held at Step 7.7 after review."
- Lines 344-362 (Heavy-test hold, drain, answer-manifest paragraph): move to `### Step 7.7: Heavy-Test Hold and Drain` before Step 8. Hold: after Steps 7 and 7.5 pass, land the implementation commit, record its hash in `commit:`, append `## Heavy Verification Plan`; the request stays `claimed` in `do-work/working/`; no status change, no `## Open Questions` line, no move to `queue/`, no checkpoint edit; recompute selection and continue. Drain: trigger is "no claimable pending REQ remains and a held claimed request exists" (also one routed here by recover's `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES`); keep the drift refusal, the union run, the runner paragraph, and the reuse paragraph as they are; replace the answer-manifest paragraph with: green → write `heavy_verified_at`, `heavy_verified_revision`, and `## Heavy Verification Result` (target revision, execution revision, one line per lane) onto the record, then run Steps 8 and 9 in this turn; red → delete `commit:` and the plan section (F3), enter Step 7's remediation path, re-hold after a fresh review; skipped lane → the request stays claimed, `HEAVY-RUN-LANE-SKIPPED` names the lane, the next drain retries, the exit summary names it. Delete the `clarify` route for drift and historical plans: they are typed findings for a human, never a hand edit.
- Line 487 (checklist Testing judgment): drop "park selected REQs as `pending-heavy-testing`"; add `□ Step 7.7: Heavy hold after review (held requests stay claimed); drain at exhaustion, finalize green in the same turn`.
- Line 501 (Error Handling row): condition becomes "Held claimed REQs remain after the queue is empty"; action points at Step 7.7.

**`actions/work-reference.md`**
- Line 14 (diagram): "skip pending-answers".
- Line 229: parenthetical becomes "also recorded at the heavy hold so dependents can build against the landed source".
- Lines 233-244: delete the `status: pending-heavy-testing` / `status_changed_at` lines and comment; keep `heavy_verified_at` / `heavy_verified_revision` with a two-line comment: written by the drain onto the claimed record before finalization; proves which execution revision the lanes checked.
- Line 282 (status row): remove the value; remove "Qualification and Testing Judgment's heavy-test hold" from the read sites.
- Line 321 (Dependency-source-ready set): "terminally successful, or normalized `claimed` with a nonblank `commit:` — a request held for heavy lanes after review; `claimed` without a commit fails closed, and a red drain withdraws `commit:` so the dependency is unmet again." Keep the "scheduling authority only" sentence.
- Line 671 (exit-summary item 8): condition becomes "any claimed REQ in `do-work/working/` still carries `## Heavy Verification Plan` without `heavy_verified_revision` after the drain ran"; in the render block replace the `clarify` sentence with "Plan drift or a stored historical-revalidation plan is a typed finding for a human; the next `do-work run` drain retries once the cause is fixed."
- Line 823 (Testing Section Template): unchanged.

**`actions/clarify.md`** — delete the Step 2.5 subsection (39-52) and the "Use when" bullet at 11; remove the value from lines 17, 31, 35 (report text becomes "No pending questions or blocked REQs — queue is clear"; "skip Steps 3–5"); delete the Step 2.5 sentence at 60.

**`actions/cleanup.md`** — remove the token from the three status lists (47, 59, 293).

**`actions/roadmap.md`** — line 54 drop the value; line 68 the bucket keys on `pending-answers` only and the `clarify` heavy sentence goes.

**`actions/restart-with-parallel-handoff.md`** — line 64: "to answer questions — `do-work clarify`, included only if some REQ is at `pending-answers`".

**`tools/do-work-cli/prime-do-work-cli.md`** — heavyverification bullet's last sentence: "Per-lane result evidence is written onto the claimed record by the work action's drain before finalization." Advance phase table, Tested row: "review → lessons/orientation → heavy hold when lanes were selected → `finalize` through advance". Add the D4 recover behavior to the `lifecycleadvance` bullet in one clause.

**Lessons** — add one `[family: status-as-phase-marker]` bullet to `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and one to `_dev/primes/lessons-action-files.md` (a hold status that a green result flips back to a weaker status makes work less available after verification than during it; track a hold by its section and commit), and refresh both rows in `do-work/lessons-index.md` in the same commit.

## Verification

From `skills/do-work/tools/do-work-cli/`: `go build ./...`, `go vet ./...`, `go test -count=1 ./...`. From the repository root: `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/shipped-package-reference-contract.sh` (changelog mirror), and the sweep, which must print nothing:

```bash
grep -rn 'pending-heavy-testing\|ResumePhase\|resume_phase\|HeavyTestingEvidence\|HeavyLaneResult\|matchingHeavyReviewPhase\|ANSWER-HEAVY' skills/do-work _dev --include='*.go' --include='*.md' --include='*.sh' | grep -v CHANGELOG
```

The REQ's GREEN also names the direct maintainer gate; run it only where the orchestrator's gate step calls for it.

## Tasks

- **T1 — RED.** Write the tests in Testing (a) to (e), delete or rewrite the listed old tests, run `go test ./internal/schemanormalization ./internal/dependencygraph ./internal/nextselection ./internal/publication ./internal/lifecycleadvance` and record the failures.
- **T2 — Go deletions, groups 1-5.** nextselection, resultmodel, schemanormalization, publication. Build and vet after each group; (a) to (d) turn GREEN.
- **T3 — Go recovery and lifecycle, group 6 (plus group 7 if F3 is accepted).** advance skip case, checkpoint count, recovery held branch and result field, claim strip. (e) turns GREEN; full module green.
- **T4 — Prose, group 8.** All action files, the prime, lessons satellites and index rows; run the sweep grep and both contract scripts.
- **T5 — Record.** Write F1 to F5 and the F3 decision into the REQ (`## Decisions`, `## Discovered Tasks` for the REQ-506/507 remediation note), Implementation Summary, Testing with the cross-REQ test list, then finalization with the release entry.
