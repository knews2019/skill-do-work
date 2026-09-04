# Independent review: REQ-505

**Approve with follow-ups — ordinary claims and bounded continuation work, but two queue edge cases miss their required disposition.** Acceptance is Partial, so the orchestrator should apply its remediation gate before completing this request.

This is read-only review preparation before claim. Reviewed queued `do-work/queue/REQ-505-move-selection-and-claim-behind-advance.md`, observed `status: pending`, `route: C`, `commit` and `heavy_verified_revision` both `716187b847d1de0402b69587a2fe5cf7e7bd8516`, and `heavy_verified_at: 2026-09-04T21:02:21Z`. The orchestrator must re-resolve identity and resume evidence when claiming. No queue, request, source, or lifecycle state was changed by this review.

## What's built

Queue-mode `advance` delegates readiness to the canonical selector, commits claims separately, holds detected cycles and explicit archive collisions, and returns frozen membership with tokenized continuation. The action's selection/claim procedure is reduced to one command-owned principle. Later evidence-gate and finalization behavior belongs to REQ-506/507 and was not judged as a regression against this saved implementation.

## Review

**Overall: 71.25%**

| Dimension | Score |
|---|---|
| Requirements | 60% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

The four percentage dimensions average 81.25%; Partial acceptance applies the specified ten-point penalty. Three of the five declared acceptance criteria are fully delivered; frozen-member blockers and collision handling are partial.

**Important findings:**

- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go:286` — Default no-argument selection drops an archive-colliding pending request before the hold loop: `freezeQueueMembers` admits only `DEPENDENCY-CYCLE` exclusions on this path, while the canonical selector classifies the queue/archive duplicate as `DEPENDENCY-AMBIGUOUS`. A nested archived twin therefore leaves the queue copy `pending`, with no `archive-collision-hold` phase or commit. Explicit-target mode works because it includes the exclusion, so the existing explicit collision test misses the default-run regression. Users see an unresolved pending request instead of the required durable collision hold. — impact-user-visible → report only
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go:139` — A frozen unconsumed member that disappears from a UR expansion is absent from both per-request maps, so the claim loop breaks without emitting a selection blocker. Replaying after that member is removed returns exit 0, `outcome: success`, empty phases/findings/claims, and the same nonempty continuation with that member still unconsumed. The run can repeat forever or appear drained without the request-specific blocker the transferred REQ-453 contract requires. The same branch is still present in the main source at review time. — impact-user-visible → report only

**Minor findings:** None.
**Acceptance:** Partial — focused public command tests pass, but independent CLI fixtures reproduce both findings at the exact saved revision.
**Suggested testing:** 2 items, detailed below.
**Follow-ups created:** None (2 findings report only).

## Requirements checklist

- [x] Default, explicit, UR, wave and fan-out queue modes select through the canonical readiness authority and commit successful claims independently. Existing public fixtures verify exact committed paths, clean trees, assignment/dependency provenance and partial claim refusal.
- [ ] Frozen membership, provenance, original flags and dispatch bounds survive continuation; later UR members are ignored and projection precedes bounding. These work, but a disappearing unconsumed UR member lacks its required genuine blocker (finding 2).
- [ ] Nested archive collisions, dependency cycles and successful blocked probes reach guarded state transactions while ordinary work continues. Explicit collisions, cycles and successful/failed probes work; default archive collisions never reach the hold transaction (finding 1).
- [x] Working/archive single-target advance remains read-only at this saved revision, and checkpoint mode retains its separate boundary. Later gate/finalization extension is outside this review range.
- [x] Step 1 and Steps 2.0/2 were collapsed to one selection-and-claim principle, and live guide/reference/clarify/prime readers were aligned. Public Go behavior tests replace the executable ownership of the removed procedure.

## Scope and traceability

Reviewed exact range `eb01a94f2dad78bf30f334e0614393d571ae362e..716187b847d1de0402b69587a2fe5cf7e7bd8516`. Its 18 substantive implementation paths match the Implementation Summary and declared Scope; one additional changed working REQ is owner bookkeeping. Five declared test seams were intentionally unused, as Qualification records. P-A-U is fully checked, and Decisions D-01 through D-04 explain the checkpoint-test update, stateless continuation flags, unbounded observation seam, and cycle consumption.

UR-098's four-part migration constraint was checked explicitly. This range has the CLI owner, deleted prose, and new behavior tests, but no predicate-file deletion. The REQ records that its original sentence-predicate RED was stale after REQ-504. Independently confirmed that the preceding REQ-504 range removed 88 lines from `_dev/tests/contracts/request-state.sh` plus its aggregator invocation, and there is no remaining named targeted-ledger/Step-2 selection predicate to delete at this base. Treat this as an already-retired predicate surface, not a missing live owner or a reason to manufacture a test-file edit. It does not excuse the missing public behavior cases above.

The restatement sweep covered the deleted targeting ledger, selection/claim steps, read-only advance statements, and heavy-review resume references across active actions, guide and contract surfaces at the target. The REQ-453 missing-member requirement is visible in the deleted contract and is the behavioral source of finding 2. No additional actionable stale restatement was established. No approach directive was assigned. Naming, standard-library dependency direction and surgical-change guardrails were checked; neither finding involves data loss or unsafe mutation.

## Acceptance evidence

Created a detached worktree under `.git/work-run-20260905/review-505` at exactly `716187b847d1de0402b69587a2fe5cf7e7bd8516` and built its public CLI. No heavy or full repository gate was rerun.

Fresh command: `go test -count=1 ./internal/lifecycleadvance ./internal/nextselection ./internal/requeststate ./internal/resultmodel` from the saved checkout's CLI module. All passed: lifecycleadvance 8.422s, nextselection 3.521s, requeststate 5.666s, resultmodel 1.252s. The lifecycle package runs the public executable matrix for default and explicit claims, UR chains and forks, later-member exclusion, successful/failed probes, nested explicit collisions, cycles, dirty partial refusal, wave/fan-out, hostile tokens, and working/archive/checkpoint boundaries.

Saved exact-revision heavy evidence was independently read from the request: CLI integrations exit 0 / 61s; staged skills exit 0 / 25s; updater exit 0 / 52s; installer exit 0 / 23s. All four are recorded without skips and at the same execution and target revision. The earlier log-header failures were recorded as infrastructure-only attempts followed by successful same-revision reruns; no result sharing across revisions was assumed.

### Reproduce finding 1

1. In a throwaway Git repository, commit `do-work/queue/REQ-901-fixture.md` with canonical `id: REQ-901`, `title: fixture`, `status: pending`; also commit `do-work/archive/UR-900/REQ-901-fixture.md` with the same id and `status: completed`.
2. Run the exact saved binary with `--repo-root <fixture> --format json advance`, with no target tokens.
3. Observed exit 0 and success; `excluded[0].code` is `DEPENDENCY-AMBIGUOUS`; `queue_advance.frozen_members`, `phases` and `claimed` are empty; queue file remains `pending`. Expected a guarded committed `archive-collision-hold`, exact duplicate-path evidence, and `blocked-archive-collision` while retaining the archive file.

### Reproduce finding 2

1. In a separate throwaway Git repository, commit pending `REQ-901` and `REQ-902` queue files, both with `user_request: UR-900`.
2. Invoke the saved binary with `--repo-root <fixture> --format json advance UR-900`. It commits REQ-901's claim and returns a ledger with REQ-901 consumed and REQ-902 unconsumed.
3. Remove only the unconsumed REQ-902 queue fixture and commit this external change. Replay the returned continuation argv exactly, adding the fixture repository root.
4. Observed exit 0 and success; selected/findings/phases/claims are empty; the selector's only exclusion identifies `UR-900`, not REQ-902; the returned continuation still contains unconsumed REQ-902 and is unchanged. Expected a typed non-success blocker with exact frozen REQ-902 identity/path and verification guidance. Both fixture removal and all commands were confined to the reviewer-owned throwaway repository.

## Suggested additional testing

1. Add a public no-target collision case alongside the existing explicit collision case, including an unrelated ready request to prove the hold and continued claim both occur.
2. Add a public UR continuation case whose unconsumed member disappears; require a typed exact-member blocker, non-success exit, and no claim or mutation. Preserve already-consumed members and the frozen set.

Self-validation checked failure paths after the passing matrix instead of treating existing tests as sufficient acceptance. Both independent fixtures failed their intended expectation before any changes; no repair was attempted. Reviewer-owned checkout, binary and fixture repositories were removed after recording this report. No reviewer subprocess or agent work remains pending.
