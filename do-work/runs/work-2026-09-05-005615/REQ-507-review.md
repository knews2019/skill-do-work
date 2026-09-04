# Review: REQ-507

**Approve with follow-ups** — request-bound finalization works through the public command; two noncritical gaps remain in the surrounding workflow and acceptance proof. This review does not clear the saved heavy-test failure.

Route C. Reviewed implementation range `8e3dbf01e0660424965d79acb2e386b6604e4780..ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`. Current restatements were inspected at `8f46295672f45ebcdc7c80d9dd072109635de983`. This was independent, read-only preparation before the orchestrator claims REQ-507; no lifecycle decision or implementation edit was made.

## What's built

`advance` now recognizes finalization as a mechanical phase, accepts one manifest, binds its request ID and path during the finalizer's single decode, and returns the existing transaction's ordered evidence. The finalizer still owns archive, release, exact commits, provenance, rollback and journal cleanup. Steps 8/9 and the two reference procedures are substantially shorter while retaining follow-up and release-content judgment.

## Findings

**Important:**

- **F01 — The blocked-flip branch points at an unavailable transition.** `skills/do-work/actions/work.md:453` tells an external-precondition failure to return to the queue through its canonical lifecycle command. REQ-507 deleted the only explicit mid-run blocked-flip procedure, while `skills/do-work/actions/work-reference.md:977` still points back to that deleted procedure. The registered request-state transitions in `skills/do-work/tools/do-work-cli/internal/requeststate/state_commands.go:20` are claim, recover, unblock, complete, fail and cancel; none accepts a claimed request and authors a new blocked condition. `recover-claim` can preserve an already-blocked working request, but cannot author that state. An operator reaching this branch must invent a mutation or stop, so the floor-agent continuation promise is incomplete. The no-edits/self-resolving judgment itself survives in the reference table, and no unsafe terminal mutation was observed. — impact-user-visible → report only
- **F02 — The fourth terminal test proves an ordinary no-release transaction, not the already-green repair exception.** `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go:94` labels its fourth row “already green without release”, but its shared fixture at lines 246–283 writes a normal Route C request and ordinary section placeholders, with no repository-gate repair marker, durable no-op evidence or matching gate evidence. The real CLI seam is exercised, and release files are correctly unchanged; that establishes the generic no-release path only. The older finalizer test has the same simplification, so it does not close this particular proof gap. — impact-rule-change → report only

**Minor / Nit:** None.

## Requirements checklist

- [x] Preserve Fold-First minting and question audience judgment.
- [x] Preserve discovered-task consolidation and impact stamping.
- [x] Bind finalization to the selected request ID and exact path before observable preparation.
- [x] Keep the existing finalizer as the sole transaction engine.
- [x] Preserve outcome, findings, changes, rollback and ordered result records; expose those records in text as well as JSON.
- [x] Remove the archive/release/staging/commit recipes from the affected Step 8/9 and reference sections.
- [x] Preserve version-source, bump, ownership, lockfile, title and voice judgments.
- [x] Exercise serial, supplied-commit, completed-with-issues, no-release and refusal paths through a real CLI binary. The specialized already-green evidence claim is limited by F02.
- [x] Keep the implementation inside the frozen 12-path Scope and preserve recorded decisions and all checked P-A-U phases.
- [ ] Fully preserve an executable floor-agent continuation for the removed failure tail: F01 leaves the blocked-flip branch incomplete.

The four-part migration was checked against the actual saved base, not the old analysis count. CLI composition, prose removal, the contract-owner change and public behavior tests are present. The contract diff is **40 additions, zero deletions**; the saved base no longer contains the captured 19 old Step 8/9 tail predicates. The new checks constrain ownership, judgment vocabulary and concise shape. The report therefore does not claim those 19 predicates were deleted by this commit, or require invented deletions to match stale capture-time counts.

## Acceptance testing

**Result: Partial.** The core terminal and refusal behavior passed; F01 leaves a peripheral failure continuation incomplete and F02 limits the claimed fourth-path evidence.

Independently executed in a clean detached checkout at the exact saved target:

- `go test -count=1 ./internal/lifecycleadvance -run 'TestAdvanceFinalization|TestAdvanceCommandPhaseMatrix'` — exit 0, package time **5.882s**. Includes the four terminal rows and missing/duplicate/empty/hostile/mismatched-manifest cases. The tests inspect archive bytes, supplied provenance, release files, Git cleanliness, typed identities and unchanged state on refusal.
- `go test -count=1 ./internal/resultmodel -run 'Finalization'` — exit 0, package time **0.307s**. Checks text/JSON parity and ordered multi-record rendering.

The request retains prior RED-before-GREEN, focused, race and repository-gate evidence. I did not independently replay the historical RED or rerun the full gate/heavy lanes. Saved heavy evidence remains five green lanes and one red staged-skills lane at `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`; it cannot be relabelled green from this review.

## Current-state attribution and restatement sweep

The terminal classifier, manifest parser, bound preparation seam, finalization fixtures and ordered text projection were compared with current source. Later finalization recovery/set-aside work is independently owned and was not attributed to REQ-507. In particular, REQ-547 owns lifecycle/checkpoint postimage corrections; this review does not re-report that repaired issue as a new REQ-507 defect.

The exact saved staged-skills failure is a stale requirement for a `do-work-board/` literal in `work.md`. Commit `c5dff3db93137ac90760751b5a26a142c13f0291` already removes that assertion with an explanation that `advance` is core-owned. That is source evidence of a repair, not execution evidence for a fresh heavy plan. The orchestrator must verify the current required lanes before clearing this request.

The restatement sweep covered Step 8/9, the commit/reference procedures, the CLI prime, public result consumers, release ownership and no-op handling, failure classification, cleanup identity references and the staged package contract. F01 records the operative broken reference/ownership chain. Remaining mentions of obsolete Step 8 substep numbers are navigation residue around already-retained branch identity rules, not additional runtime findings.

## Scores

**Overall: 82.5%** = (90 + 95 + 85 + 100) / 4 − 10 for Partial acceptance.

| Dimension | Score | Basis |
|---|---|---|
| Requirements | 90% | Nine of ten requirements delivered; blocked-flip continuation incomplete |
| Code quality | 95% | Small composition seam, one decode, strict input handling; surrounding workflow mismatch |
| Test adequacy | 85% | Real CLI terminal/refusal coverage and output parity; specialized no-op fixture incomplete |
| Scope | 100% | Exact 12 source paths match declared Scope; decisions are recorded |
| Risk | Low | No introduced corruption or authority bypass reproduced; workflow/proof gaps remain |
| Acceptance | Partial | Core terminal behavior passes; F01/F02 limit complete workflow acceptance |

## Suggested additional testing

1. At claim time, run the required current heavy verification plan and retain its actual execution revision, including the already-repaired staged assertion.
2. If F02 is explicitly captured later, use a genuine already-green repair record and matching durable gate evidence through public `advance`, with a negative control for invalid no-op evidence. Do not count a generic no-release fixture as that proof.

## Follow-ups created

None (2 findings report only). No queue, request, source or commit mutation was made. Self-validation checked the real fixture contents, distinguished pre-existing/fixed issues from this diff, and kept historical green evidence separate from the held heavy result.

## Handoff applicability

After claim, re-resolve REQ-507's saved base/target and retained Scope before applying this review. This report can serve as the independent saved-implementation review if those still match and no remediation changes the reviewed source. It does not substitute for required current testing. The detached review worktree was clean and removed without force after both processes completed.
