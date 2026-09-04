# Review: REQ-515 (per-REQ recovery findings never stop the loop)

**Request changes** — the per-record folding in the CLI is correct and well built, but under `do-work run-with-recovery` the same `recover` call that sets a REQ aside then releases that REQ's claim back into `do-work/queue/`, so the REQ the run just excluded is selected again in the same run.
Route B | range `18666d7..28c1460` (builder commit `7309b8a`, seam commit `960d8b3`)

### What's built

- `recover-finalization` now folds journal results one at a time. A refusal owned by one REQ becomes that REQ's record with `FINALIZATION-SET-ASIDE` appended to `reason_codes` and an empty `next_argv`, and the drain continues; the command still returns typed `success`. Verified end to end: two prepared journals, one interrupted, and `next` selects the untouched third REQ.
- `work.md` Step 1, `run-with-recovery.md` Step 1, `work-reference.md` (Step 9 finalization paragraph plus a new exit-summary section 9) and `commit.md` line 51 all read the record contract the same way.
- A new owner contract `_dev/tests/contracts/recovery-set-aside.sh` (10 predicates) plus two Go tests. Focused Go packages, the contract file, `go vet` and ShellCheck all pass on my run.
- What is not built: the exclusion never reaches the selector. `advance` and `next` do not read `finalizations`, so nothing mechanically prevents a set-aside REQ from being selected — under plain `run` it is protected only because the REQ happens to stay in `do-work/working/`, and under `run-with-recovery` that protection is actively removed.

### Decisions / risks for you

- **D-04 (recovery keeps returning typed `success`) has one caller the builder did not check: the Go one.** `internal/lifecycleadvance/recovery_commands.go:57` stops `recover` only on a non-success finalization result. A set-aside is now success, so `recover` walks on into claim recovery, and with `--assume-sole-authority` it recover-claims the set-aside REQ — moving it to `do-work/queue/`, status `pending`, committed. The builder's D-04 reasoned about the half of that caller that would have parked the loop and not the half that now runs. Value of D-04 stands (`OutcomeFindings` would have parked everything); the missing piece is excluding the set-aside REQ from claim recovery.
- **D-02 (ownership plus residue instead of a path test) is weaker than it reads.** Every journal-level refusal is built by `finalizationFailure` (`finalization_apply.go:628-636`), which always stamps `AffectedIDs: [this REQ]`. So the ownership half of the test is always true and the only condition that can still stop the drain is an incomplete rollback. Shared-cause refusals — `FINALIZATION-DIRTY-INDEX` (the one git index) and `FINALIZATION-AMBIGUOUS-SHARED-STATE` (dirty `do-work/**` or release-metadata paths) — are set aside per REQ. The saving grace is that this is fail-closed downstream: any later commit transaction refuses on a dirty index, so nothing foreign gets committed.
- **D-01 (no `next_argv`, verb from the exit summary) works under plain `run` and is thin under `rwr`.** The verb fork in **Stuck Runs Hand Off to Judgment** offers `do-work run-with-recovery` or `do-work cleanup`; when the run already is `run-with-recovery`, the first option is the verb that just ran.

### Findings

**Important:**

- F1 — `recover --assume-sole-authority` releases the set-aside REQ back into the queue, so `run-with-recovery` re-selects it. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go:57-60` and `:78-110`, reached because `internal/finalization/finalization_commands.go:118-131` now returns typed `success`. Reproduced on a copy of the module with the builder's own `seedTwoPlannedFinalizations` fixture: after the set-aside, `do-work/working/REQ-720.md` is still `claimed`, the recover-claim transition applied by `recover` succeeds, `do-work/queue/REQ-720.md` appears, and `next` selects `REQ-720` — the exact REQ the run set aside — while its journal is still on disk at `.git/do-work-finalization/REQ-720.json` at phase `prepared`. A builder then redoes the work, and the finalize tail refuses at `prepareJournal` (`finalization_prepare.go:44`, "an unfinished journal exists ... with a different manifest"). This contradicts the text this REQ added in three places: `run-with-recovery.md:33` ("one REQ this run will not select"), `work.md:127` ("excludes that one REQ from this run's selection"), and `work-reference.md:678` ("the REQ was not claimed and nothing was written to it" — the claim reset is written and committed). — `impact-critical` → needs a queued follow-up REQ (this reviewer does not write under `do-work/`)
- F2 — The new "a finding that owns no REQ stops the run" branch is unreachable, so shared-cause refusals are set aside. `internal/finalization/finalization_commands.go:139-155`: `requestScopedRefusal` can only return false for an empty request id, a nil `Finalization`, no findings, `OutcomeFailure`, or `RollbackIncomplete`. `advanceJournal` never produces the first four, so incomplete rollback is the only live stop. `FINALIZATION-DIRTY-INDEX` and `FINALIZATION-AMBIGUOUS-SHARED-STATE` — the two codes that mean shared-target dirt — are stamped with the journal's REQ and are now per-REQ exclusions. Detailed Requirement 3 is still satisfied in practice only because discovery-level refusals (`finalization_discovery.go:1703`, no `AffectedIDs`) return before the folding loop and were not touched by this diff. — `impact-rule-change` → report only
- F3 — The never-widen lock-in test does not exercise the code this REQ added. `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest` seeds a foreign checkpoint entry, which refuses inside `discoverFinalizationJournals` and returns at `finalization_commands.go:68-74`, before `consumeRecoveryRecord` is ever called. It would pass unchanged if `requestScopedRefusal` returned true for everything. No test covers any false branch of `requestScopedRefusal` — including the one condition (`RollbackIncomplete`) that is actually load-bearing. The REQ's own Batch Constraint asks for a behavior test on the command; the RED test covers the continue path only. — `impact-rule-change` → report only
- F4 — The exclusion never reaches the selector. The REQ's `## What` asks for "an exclusion with its reason code in the selector output"; `advance` and `next` read no finalization records (`internal/lifecycleadvance/advance_commands.go`, `internal/nextselection/`), so the exclusion exists only in the recovery result and in prose telling the agent to honour it — while `work.md:127` also tells the agent that `advance`'s typed result is the sole selection authority. Under plain `run` the two never collide only by accident of the REQ staying in `working/`. — `impact-user-visible` → report only

**Minor:**

- M1 — `namesRequestID` (`finalization_commands.go:172-179`) tests membership, not exclusivity: a finding whose `AffectedIDs` were `[REQ-720, REQ-721]` would be set aside as REQ-720's private exclusion. No producer emits more than one id today, so this is latent; `len(AffectedIDs) == 1` would close it. — `impact-rule-change` → report only
- M2 — Under `run-with-recovery` the exit summary's prescribed resolving verb (`work-reference.md:687`) resolves to `do-work run-with-recovery`, the verb that just ran. REQ-514's rule is about the CLI's `next_argv`, so this is not a violation, but the user-facing advice is circular in the mode where set-asides are most likely. — `impact-user-visible` → report only
- M3 — `work-reference.md:678` says to render each line "from the record's own fields", but the template line contains `[title]`, which is not a field of `FinalizationResult`. A floor agent must open the REQ file for it. — `impact-negligible` → report only
- M4 — Seven of the ten contract predicates are exact-phrase pins (`one record at a time`, `never as one pass/fail gate for the run`, `reason codes, comma-separated`, `recover: <resolving verb>`, `its finding names no REQ`, `Set-aside-by-recovery section`). They redden on a paraphrase with no behavior change. This matches how the other owner contracts in `_dev/tests/contracts/` work and the Batch Constraint is met (the Go RED test is the behavior half), so it is not a sentence pin alone — but only one predicate, the `setAsideReasonCode` extraction from the Go source, is a real cross-artifact check. — `impact-negligible` → report only
- M5 — `_dev/tests/contracts/recovery-set-aside.sh:38`: the missing-action-file branch is guarded by `[ -n "$set_aside_code" ]`, which is unrelated to whether the file exists. With an empty code and a missing file, the run falls through to `grep` on a nonexistent path and prints a misleading "must state" failure plus grep's own stderr. — `impact-negligible` → report only
- M6 — The REQ's `## AI Execution State (P-A-U Loop)` boxes are all unchecked. The hand-back's Verification section shows the equivalent work (diff reviewed file by file, vet, gofmt, `git diff --check`), so this is bookkeeping, not a skipped phase. — `impact-negligible` → report only

**Nit:**

- N1 — One file outside the declared write set (`_dev/tests/contracts/recovery-set-aside.sh`). D-03 flags it and the reasoning (the runner is at a hard 77-line ratchet; predicates live in owner contracts) is right. — `impact-negligible` → report only

### Requirements Checklist

- [~] Step 1 in `work.md` and Step 0.1 in `run-with-recovery.md` iterate recovery records; a refused record excludes that REQ from this run's selection and the loop continues — **partially delivered**. Prose is in both files and the CLI folds per record, but under `run-with-recovery` the exclusion is undone by the same command (F1); under plain `run` it holds only incidentally (F4).
- [x] The composed exit summary lists set-aside REQs with reason codes and resolving verbs — delivered (`work-reference.md:678-687`), pinned by three contract predicates.
- [~] A finding with no owning REQ still stops the run and names a resolving verb per REQ-514 — **partially delivered**. True today, but through the untouched discovery-refusal path; the new gate's unowned branch is unreachable and shared-cause journal refusals are set aside (F2).
- [x] Contract predicates on the per-record wording replace the whole-run pin, and the CLI carries a behavior test for a mixed result (one refused, one clean, selection proceeds) — delivered; `TestRecoverFinalizationSetsAsideRefusedRecordAndFinishesTheRest` is a real RED→GREEN with the failure captured as an assertion, not a compile error.
- [x] Serial and fan-out modes behave the same — delivered. Both read the one Step 1 paragraph in `work.md`; the mode split happens after selection. (The divergence found is `run` versus `run-with-recovery`, a different axis — F1.)
- [~] Constraint "never widen what recovery accepts" — recovery still refuses exactly what it refused; what changed is that two shared-cause refusals no longer stop the run (F2). Downstream is fail-closed: a dirty index makes every later commit transaction refuse, so no foreign bytes can be committed.
- [x] Constraint "keep the floor agent able to follow the loop" — the prescribed argv is `--format json`, the record carries `reason_codes` and an empty `next_argv`, and the prose says what to do with each. Adequate.
- [x] Batch constraint "behavior test or contract predicate, never a sentence pin alone" — both exist.
- [x] Batch constraint "judgment stays prose, mechanics stay in the Go CLI" — held; no new prose walks a shell sequence.

Requirements compliance: 6 of 9 fully delivered, 3 partial.

### Acceptance Testing

**Result: Fail**

- `bash _dev/tests/contracts/recovery-set-aside.sh` → passes, exit 0. `shellcheck` on it → clean.
- `go test -count=1 ./internal/finalization ./internal/lifecycleadvance ./internal/resultmodel ./internal/requeststate` → all pass (finalization 23.5s). `go vet ./internal/finalization` → clean.
- Set-aside behaviour (happy path): reproduced the builder's fixture on a copy of the module. `recover-finalization` returns `success`, record REQ-720 carries `[FINALIZATION-JOURNAL-WRITE FINALIZATION-SET-ASIDE]` at phase `prepared` with empty `next_argv`, record REQ-721 reaches `cleanup_complete`, and selection proceeds. This is the REQ's Red-Green Proof and it holds for `do-work run`.
- Set-aside behaviour under `run-with-recovery` (the failure): continuing from the same state with exactly the claim-recovery step `handleRecover` performs (`TransitionRecover`, `AssumeSoleWriter`, `Commit`, `CheckpointAllEntries`), the transition succeeds, `do-work/queue/REQ-720.md` is created, the stale journal remains, and `next` selects `REQ-720`. The REQ that was set aside is selected again in the same run.
- Never-widen probes requested: a finding naming no REQ → still stops, but via the untouched discovery path, not the new gate. A command-level `failure` outcome → guarded and unreachable from `advanceJournal`. An incomplete rollback → guarded, stops, untested. A finding naming more than one REQ → would be set aside (M1), no producer today.
- Not tested: a live `do-work run-with-recovery` against a real queue (no safe way to do that in this checkout while the maintainer verify suite is running).

### Suggested Additional Testing

1. A Go test in `internal/lifecycleadvance` that runs `recover --assume-sole-authority` over a set-aside fixture and asserts the set-aside REQ's claim is *not* released — this is F1's regression pin.
2. A Go test that drives `requestScopedRefusal` to false through `RollbackIncomplete` and asserts the aggregate stops, so the one live stop condition is pinned.
3. A `FINALIZATION-DIRTY-INDEX` fixture (stage a foreign path before recovery) to decide deliberately whether that shared cause should stop or be set aside, and pin whichever answer is chosen.
4. Manual: run `do-work run` and then `do-work run-with-recovery` against a queue with one stuck finalization tail and read both exit summaries — check the set-aside section renders and that the second run does not re-dispatch the stuck REQ.

### Scores (on the record — not the headline)

**Overall: 50%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 65% | 6 of 9 delivered; the exclusion does not hold under `run-with-recovery` and never reaches the selector |
| Code Quality | 80% | Small, well-named helpers with comments that explain why; blemishes are an unreachable guard and a membership-not-exclusivity ownership test |
| Test Adequacy | 55% | Real RED→GREEN on the continue path; the never-widen lock-in tests untouched code and no false branch of the new gate is covered |
| Scope | 90% | One justified file outside the write set, flagged in D-03; restatement sweep was run and `work.md:588` correctly judged not stale |
| Risk | Critical | The changed `recover` outcome contract has a Go caller that now releases a claim the run declared excluded (F1) |
| Acceptance | Fail | The REQ's stated behavior does not hold under one of the two entry points the REQ names |

Formula: average of 65/80/55/90 = 72.5%; Acceptance = Fail caps at 50%; Risk = Critical caps at 60%. Overall 50%.

### Follow-ups created

None written by this reviewer — the orchestrator owns writes under `do-work/`. F1 is graded `impact-critical` and needs one queued follow-up REQ (suggested shape: `recover` must skip claim recovery for any REQ whose finalization record carries `FINALIZATION-SET-ASIDE`, with the regression test from Suggested Testing item 1). The other 9 findings are report only.

### Restatement Sweep

Run. The diff redefines how `finalizations` records are read after a typed `success`. Grepped every shipped `.md` for `finalizations`, `reason_codes`, `blocked_paths`, "Continue only on typed", and `recover-finalization`. Consumers found: `run-with-recovery.md:33`, `commit.md:51`, `work.md:127`, `work-reference.md:678`, `work-reference.md:1081` — all five updated. `work.md:588` is the `finalize --manifest` paragraph and is correctly left alone: `handleFinalize` calls `advanceJournal` directly and never reaches `consumeRecoveryRecord`, so no set-aside code can appear there. No shell script or hook parses finalization records. One stale statement remains and it is inside the new text itself: `work-reference.md:678` says the set-aside REQ "was not claimed and nothing was written to it", which is false under `run-with-recovery` (F1).

### Self-Validation

Re-checked three assumptions before writing this up. (1) That the rollback really restores the REQ to `working/claimed` rather than leaving it in `archive/` — confirmed by running the fixture and reading `rollbackBeforePrimary`. (2) That the downstream damage from F2 is bounded — confirmed: `gittransaction` refuses every commit on a non-empty index (`git_transaction.go:183-191`), so a shared-dirt set-aside wastes work but cannot commit foreign bytes. (3) That plain `do-work run` is unaffected — confirmed: `work.md` Step 1 invokes `recover` with no authority flag, claims are preserved, and the set-aside REQ stays in `do-work/working/` where queue selection cannot see it.

---

**Verdict: Fail**

Acceptance reasoning: the REQ's first Detailed Requirement names both `actions/work.md` Step 1 and `actions/run-with-recovery.md` Step 0.1, and requires that a refused record "excludes that REQ from this run's selection". Under `run-with-recovery` it does not: the same `recover --assume-sole-authority` invocation that produced the set-aside then recover-claims that REQ back into `do-work/queue/` as `pending`, and `next` selects it — verified by running the builder's own fixture through the claim-recovery step. The run then dispatches a builder onto a REQ whose stale journal will make its finalize tail refuse. The three action paragraphs this REQ added all assert the opposite, so the shipped prose is wrong about the shipped behavior. That is an acceptance failure on a named requirement in a named file, not an edge case, and it caps the score at 50%.

What is right is worth saying: the per-record folding, the REQ-514-conformant set-aside shape (empty `next_argv`, same-command `verification_argv`, own reason codes preserved), the honest RED→GREEN, and the `setAsideReasonCode` cross-artifact predicate are all good work, and the plain `do-work run` path in the REQ's own Red-Green Proof does what it promises. One bounded fix in `internal/lifecycleadvance/recovery_commands.go` — skip claim recovery for REQs carrying `FINALIZATION-SET-ASIDE` — plus a test on the `RollbackIncomplete` stop would move this to Approve.

*Reviewed by review-work action (independent reviewer)*
