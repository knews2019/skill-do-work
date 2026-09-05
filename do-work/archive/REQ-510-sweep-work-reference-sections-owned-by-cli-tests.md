---
id: REQ-510
title: '[impact-rule-change] Sweep work-reference sections whose contract is now a CLI behavior test'
status: completed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-509]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md, _dev/tests/contract-regressions.sh, skills/do-work/docs/]
claimed_at: 2026-09-05T00:33:24Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
completed_at: 2026-09-05T08:02:54Z
commit: 87226175
release_at: 2026-09-05T08:02:54Z
---

# Sweep work-reference Sections Whose Contract Is Now a CLI Behavior Test

## What
Delete every `work-reference.md` section whose contract moved into a Go behavior test during REQ-503 to REQ-509, and fix the cross-references in `work.md` and docs. Keep the Execution Model, the schema read contract, the Fold-First and Ratchet homes, and the minimal templates.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both primes and both lesson satellites, the coding/anti-slop crew rules and `CLAUDE.md`, then re-derived the plan against the live files — the prepared plan's per-area line budget was measured against an assumption that rewriting a paragraph shrinks the line count; in this file a paragraph *is* one line (161,545 bytes over 1,064 lines), so only structural deletions move the count. Evidence: the plan's "frontmatter template −55" was measurable at −11; its "templates −90" at ≈ −15.
- [x] **[APPLY]:** Only the five declared paths changed. Evidence: `git diff --stat` on the commit reads `_dev/tests/contracts/core-checks.sh | 2 +-`, `work-reference.md | 375 ++--`, `work.md | 2 +-`, `command-line-guide.md | 2 +-`, `work-guide.md | 2 +-`; 71 insertions, 312 deletions.
- [x] **[UNIFY]:** Read the whole diff, re-ran both suites and the named Go owner tests, and linted the one shell file at the repo's own severity. Evidence: `shellcheck -x --severity=warning _dev/tests/contracts/core-checks.sh` → exit 0 (the 29 `info`-level SC2016/SC2030/SC2031 findings are pre-existing and identical on `HEAD`); no debug artifacts, no `do-work/` path staged.

## Why
`work-reference.md` is 1,250 lines with 66 inbound references from `work.md`; after the chain, many sections restate what a command now guarantees.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- For each candidate section, name the Go test that now owns the contract before deleting; a section with no owning test stays.
- Fix every inbound reference in the same commit; `shipped-package-reference-contract.sh` must pass.
- Delete the matching predicates.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-509; last in the chain.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete the Commit and Metadata-Commit Procedure section and run the contract and shipped-reference suites.
**Why RED now:** Predicates and reference checks fail; inbound links dangle.
**GREEN when:** Both suites pass; `work-reference.md` is under 700 lines; every remaining section is judgment, schema, or a minimal template.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Addendum (2026-09-03, 23:45 local)

User added (23:35 local, "yes, on the nine cancelation", applying the velocity report's triage table):

- REQ-471 (flow and reader consistency plus documentation for the gate-blocked set-aside) was cancelled into this sweep. When deleting `work-reference.md` sections here, also sweep any surviving sentence that says an unrelated canonical-gate failure must preserve a claim and stop the session; the shipped behaviour is the deferral lifecycle (REQ-491 to REQ-494) plus the retry-once rule REQ-559 adds, and queue summaries must distinguish blocked work, pending user decisions and dependency-gated work. One sweep, no new section.
- Coherence check: no contradiction; this widens the sweep's search condition by one sentence family.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is stated exactly — delete every work-reference section whose contract moved into a Go behavior test — but which sections qualify, and which test owns each one, has to be discovered by reading both the prose and the test suite. That discovery is the whole job, so exploration is required and a plan on top of it would only restate the request.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
## Exploration

Both revisions of `work-reference.md` were split on `^## ` outside fences and compared section by section, and every candidate section was traced to the Go behavior test that would own its contract after REQ-503 through REQ-509. The rule the request sets is strict and decided every call: a section with no owning test stays. That is why Worktree Dispatch Mode survives at 92 lines — the CLI implements no worktree dispatch verb at all, so not one of those lines has an owner.

The prepared plan in `do-work/runs/work-2026-09-05-005615/REQ-510-plan.md` was re-validated rather than followed. Its line budget assumed rewriting prose shrinks line counts; in this file a paragraph is one line, so only structural deletions move the number. Measured after the sweep: bytes fell 10.7% while lines fell 22.6%, and 152 of the remaining lines are still over 300 characters.

Every predicate pinning this file across `core-checks.sh`, `recovery-set-aside.sh` and `prescribed-shell-canonicalization.sh` was read before anything was deleted. All of them pin judgment that survived, so none needed deleting — the one required change was an awk boundary in `core-checks.sh` that used a deleted heading as its end marker and would otherwise have captured to end-of-file.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/work.md`
- `_dev/tests/contracts/core-checks.sh`
- `skills/do-work/docs/work-guide.md`
- `skills/do-work/docs/command-line-guide.md`

`work.md` and the two docs are cross-reference repair only; this request is explicitly not a rewrite of `work.md`. The declared `write_set` names `_dev/tests/contract-regressions.sh`, and the change landed in `_dev/tests/contracts/core-checks.sh`, which is the same suite — the aggregate delegates to it.

**Acceptance criteria:** the contract suite and the shipped-package reference check both pass; every deletion names the Go test that now owns its contract; every inbound reference is repaired in the same commit.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `_dev/tests/contracts/core-checks.sh` (modified)

**What was done:** Deleted or shrank every work-reference.md section whose contract is now owned by a named Go behavior test, and repaired the inbound citations that the deletions touched. The file went from 1,064 lines to 823 in the first commit, then to 849 after the review round restored two sections that should not have been condensed.

Five sections were deleted outright: the repository-gate deferral transaction, the Targeted Run Ledger, the Session Checkpoint Principle, auto-wave's five-condition ready-set computation, and the defer-gate fold-topology paragraph. Five more were shrunk to the judgment that survives, with the heading kept: Composed Exit Summary, In-Progress Record, Commit and Metadata-Commit Procedure, the already-green repair no-op, and Failure Classification. Every deletion is backed in the hand-back by the Go test that now owns the contract, and each of those tests was run and passed. One deletion also removed a live contradiction: auto-wave's condition 2 still said a dependency is satisfied only by completed or completed-with-issues, which contradicted the dependency-source-ready status set defined three sections above it.

Nothing was deleted without an owner, which is why five headings now carry only two or three sentences. Sections with no owning Go test stayed: Execution Model, Folder Structure, the request file schema, the schema read contract, the stuck-runs hand-off, crash recovery, worktree dispatch, the repository gate deferral rules, and every minimal template.

The review round restored two sections verbatim from the pre-sweep revision. Folder Structure had lost the archive/legacy node while cleanup.md and abandon.md still route work into it, and the Progress Reporting Example had lost ten lines including the only Route A illustration. The exit-summary table got its nine exact warning headline strings back at zero extra lines, one stale internal citation was rewritten, and the closed four-item ready-set enumeration was closed at four sites rather than three. The builder then rebuilt the section ledger by splitting both revisions on heading lines and comparing, which surfaced three changed sections the first ledger had omitted.

Reference repair was done by grep on heading names, not by link, because the shipped-reference checker does not resolve prose citations of the arrow-and-bold form. Four inbound sites were updated: the architecture citation in work.md, the checkpoint citation inside work-reference.md, the awk end boundary in core-checks.sh, and the two docs sentences that had become stale or wrong once the canonical text moved.

One scope note. The request's write set names _dev/tests/contract-regressions.sh, while the one-line predicate fix landed in _dev/tests/contracts/core-checks.sh, the probe file that suite runs.

The merged range carries work from other requests as well. Of the 27 entries in the merged-range diff stat, the five listed above are this request's; the rest belong to other work merged in the same range: the Kanban board tool, the publication answer package, three request reservations and one user-request capture. The two commits for this request stat to exactly these five files.

**Implementation range:** `974713ac..87226175`. Builder commit `0e599236269b45e693deb55e3a48d8cc5752b605`, on top of the first commit 6eeeee8de18b50e6864f9de5af5e3e73e130ee9d.

## Decisions

- **D-01 — Kept every heading cited from a file outside the write set, even where the body was gutted:** the commit and metadata-commit procedure, crash recovery, the in-progress record, worktree dispatch and its sub-anchors all have callers in run-with-recovery.md, forensics.md, cleanup.md, review-work.md, board.md and the board guide. The request's constraints say scope must be updated before another path is touched, and a builder cannot edit the request, so renaming a heading would have traded a deleted section for a dangling citation the builder was not allowed to repair. The builder names this the close call of the whole request.
- **D-02 — Condensed the 48-line architecture diagram instead of deleting it, and said so rather than inventing an owner:** no Go test owns an ASCII drawing, because there is no contract in it; the pipeline order is work.md's own step list. This deviates from the request's strict rule that a section with no owning test stays, and the builder recorded it as a condensation and flagged it as the place to push back. After the review round it is the only deliberate condensation of unowned material that survives; reverting it is one command away and costs 43 lines.
- **D-03 — Turned the composed exit summary's nine repeated render blocks into one table rather than deleting the section, then corrected the claim:** the per-category triggers are advance's typed exclusion codes and recover's records, but the headline choice, section order, remedy sentences and resolving-verb judgment are the action's and nothing tests them. The original claim that the table kept everything was wrong: the nine exact warning headline strings lived only in the fenced blocks. The builder put them back inside the existing table cells at zero line cost and states plainly that the trade bought nothing. What the table did drop, and the builder still stands behind, is nine repeated "applies if any REQ has status X" clauses that were a second queue scan advance already answers.
- **D-04 — Kept the crash recovery section even though recover is fully tested:** the six surviving lines are the authority decision, and three other actions cite this heading as the home for takeover judgment. Deleting it would have saved six lines and cost four repairs the builder could not make.
- **D-05 — Restored one sentence that had been over-cut:** the first pass dropped "a successful repair causes a fresh selection", which is what lets suppressed parents resume. Caught on diff review and put back. Recorded because the sweep's real risk is exactly this shape: a deletion that looks like duplication and is a guarantee.
- **D-06 — Repaired the two docs sentences the sweep implicated, and no more:** command-line-guide.md claimed the whole CLI never executes gates or schedules queue work, which contradicted the new ownership map, so it is now scoped to defer-gate. The auto-wave gloss in work-guide.md became the last closed enumeration of the ready set once the canonical one was deleted, so it gained "and not dropped by a filter you passed".
- **D-07 — Reported the under-700 miss instead of meeting it:** work-reference.md is 849 lines, not under 700. This deviates from the request's stated GREEN criterion. Reaching 699 needs 150 lines out of a 323-line group of judgment with no CLI owner, whose largest block is worktree dispatch at 91 lines for a verb the CLI does not have. The request, the dispatch and the prepared plan all forbid deleting an unowned guarantee, so the builder stopped and reported the conflict.
- **D-08 — Restored by revert, not by rewrite:** for the two lost sections the builder took the exact pre-sweep text rather than writing improved replacements. A section that should not have been touched comes back as it was; re-authoring it would have made the loss harder to audit and put new unowned prose in its place.

## Qualification

Passed the request-bound advance qualify and scope-drift gates for `974713ac..87226175` against the merged range, both satisfied. Five declared files changed across the request's two commits, 71 insertions and 286 deletions when restricted to them; the range stat is larger only because integration was serial and interleaved with six other requests.

Independent review verified each of the six disclosed deletions names a Go test that genuinely pins the deleted contract, reading the named tests rather than accepting the ledger. It found two condensations the first ledger listed as kept, and both were restored verbatim in the continuation.

The under-700 acceptance clause is not met at 849 lines, and the review agrees with the refusal after re-deriving every number: the explicit keep-list alone is 497 lines and the remaining 323 are judgment with no CLI owner, led by 91 worktree lines the CLI implements no verb for. Reaching 699 would mean deleting an unowned guarantee, which the plan forbids.

The P-A-U boxes were reconciled from the builder hand-back, which is where worktree dispatch puts them.
## Testing

**Red-green validation:** The request's named RED case was run first — delete only `## Commit & Metadata-Commit Procedure (Step 9)` (lines 994-1001 of the 1,064-line file), leave callers and predicates untouched, then run both suites. `DO_WORK_TEST_FILE_BUDGET_SECONDS=300 bash _dev/tests/contract-regressions.sh` exited 1 with three failures, quoted from the hand-back:

```
FAIL: work-reference.md commit procedure must hand one judged manifest to advance and consume ordered evidence.
FAIL: actions/work-reference.md must state that finalization records are read per REQ — missing: one record at a time
FAIL: actions/work-reference.md must state that the surviving whole-run stop is the finding no REQ owns — missing: its finding names no REQ
```

In the same RED state `bash _dev/tests/shipped-package-reference-contract.sh` exited 0 and printed `shipped package reference contract: PASS`, with the heading gone and two live citations still naming it. That PASS is itself a finding: the checker resolves Markdown links and cross-package path citations, not `→ **Named Heading**` prose. The dangling callers had to be found by grep, in `skills/do-work/actions/run-with-recovery.md:33` and `skills/do-work/actions/work.md:126`. The RED deletion was reverted before the real edit.

GREEN at the final commit `0e599236`: `bash _dev/tests/shipped-package-reference-contract.sh` prints `shipped package reference contract: PASS` and exits 0, and `DO_WORK_TEST_FILE_BUDGET_SECONDS=300 bash _dev/tests/contract-regressions.sh` prints `Contract regression checks passed.` and exits 0 with zero FAIL lines. The Go tests that own the deleted contracts were run by name with `-count=1` and pass; the ones the ledger leans on hardest are `TestAdvanceCheckpointChangesOnlyCheckpointAndPreservesLiveEntries`, `TestAdvanceCheckpointPreservesLegacyClaimDiscovery` and `TestWorkingAdvanceRemainsReadOnlyAfterCheckpointMode` for the deleted checkpoint principle; `TestAdvanceQueueTargetLedgerProjectsBeforeBounding`, `TestAdvanceQueueClaimsURChainAndForkAcrossStatelessContinuation`, `TestAdvanceQueueReportsCommittedPartialClaimBeforeDirtyRefusal` and `TestTargetedURReplayObservesFrozenReadyWorkBeforeApplyingSavedFanOut` for the deleted targeted run ledger; `TestDeferGateCreatePublishesOneAtomicDependencyLifecycle`, `TestDeferGateFoldsSharedFingerprintWithoutOverwritingPriorParent`, `TestDeferGateClassifiesTrackedDirtyTrackedCleanAndUntrackedPreimagesIndependently`, `TestDeferGateTrackedDirtyRepairFoldRollsBackExactPreimagesAfterEveryMutation` and `TestDeferGateRefusesUnsafeStaleCollidingAndStagedInputs` for the deleted deferral transaction; `TestExplicitREQOverridesDependencyAssignmentAndNegligibleFilters`, `TestClaimEvidenceVetoesEverySelectionModeBeforePolicyOrProbe`, `TestHeldClaimedSourceAllowsDependentSelection`, `TestWaveDepthAndFanOutAreSeparateSelectionAxes`, `TestDependencySatisfactionUsesOnlyTerminalSuccess`, `TestClaimedDependencyWithCommitIsSourceReady` and `TestAdvanceQueueWaveFanOutAndHostileTokens` for the deleted auto-wave ready set; and `TestGreenGateEvidenceIdentityAndHistory`, `TestGreenGateEvidenceStaysVerifiableAtRecordedRevisionAfterUnrelatedHeadMove`, `TestGreenGateEvidenceTargetToleratesGateLogCommits`, `TestGreenGateEvidenceRejectsDivergentAndForeignRecords` and `TestGreenGateEvidenceFailsClosedForInvalidTargets` for the deleted green-gate record internals.

**Controls preserved:** `_dev/tests/contract-regressions.sh` was re-run at both budgets; it protects the predicates that pin this file — the qualification anti-rationalization table and the finding-closure ratchet must exist, seven release-judgment tokens in the changelog section, the two surviving commit-section strings ("Pass that one manifest to the current `advance` continuation." and "Consume ordered `finalizations`"), and the `tail_recipe` must-not-match that keeps a hand tail out of the changelog and commit sections. `_dev/tests/recovery-set-aside.sh` protects REQ-515's per-record guard through six literal strings (`FINALIZATION-SET-ASIDE`, `one record at a time`, `Set-aside-by-recovery section`, `reason codes, comma-separated`, `recover: <resolving verb>`, `its finding names no REQ`), all retained verbatim including inside the rewritten exit-summary table, plus the absence of the retired whole-run gate phrase. `_dev/tests/shipped-package-reference-contract.sh` protects Markdown links and cross-package path citations. `_dev/tests/prescribed-shell-canonicalization.sh` protects the pointer from work-reference.md to `../docs/prescribed-shell-primitives.md`. No predicate was deleted: every predicate that pins this file pins judgment that survived. One predicate needed repair rather than deletion — `reference_commit`'s awk range in `_dev/tests/contracts/core-checks.sh` used the now-deleted `## Session Checkpoint Principle (Step 10)` heading as its end boundary and would have captured to end of file, so it now ends at `## Progress Reporting Example`.

**Module verification:** Commands and results as recorded in the hand-back.

```
$ DO_WORK_TEST_FILE_BUDGET_SECONDS=300 bash _dev/tests/contract-regressions.sh
Contract regression checks passed.                       (exit 0, zero FAIL lines)

$ bash _dev/tests/shipped-package-reference-contract.sh
shipped package reference contract: PASS                 (exit 0)

$ shellcheck -x --severity=warning _dev/tests/contracts/core-checks.sh
(no output; exit 0)

$ wc -l skills/do-work/actions/work-reference.md
849
```

Go owner tests, run with `-count=1` from `skills/do-work/tools/do-work-cli` at the first commit:

```
ok  internal/lifecycleadvance   5.765s   (finalization matrix, queue ledger, checkpoint, recovery, archive-collision)
ok  internal/publication        5.323s   (defer-gate create/fold/rollback/refusal, archived review follow-up)
ok  internal/repairvalidation   0.981s   (whole package)
ok  internal/nextselection      0.021s   (gate priority, explicit-REQ override, claim veto, wave axes, held source)
ok  internal/requeststate       0.435s   (six transitions, error types, claim footprint, checkpoint removal, target resolution)
ok  internal/gateevidence       2.300s   (whole package)
ok  internal/dependencygraph    0.008s   (whole package)
ok  internal/finalization       1.762s   (set-aside, shared dirt, supplied provenance, already-green no-release, UR move)
ok  internal/requestmodel       0.005s   (alias precedence, read-alias selection)
```

The Go packages were not re-run in the review round because that round changed only Markdown. The shellcheck run is from the first round for the same reason: no shell file changed in the second round.

Two reported failures were diagnosed as not belonging to this change. At the stock 30-second per-test-file budget, `bash _dev/tests/contract-regressions.sh` exits 1:

```
SessionStart hook behavior probes passed.
test-file duration: session-start-hook-behavior.sh 35s (limit <30s)
FAIL: … session-start-hook-behavior.sh took 35s; each test file must finish under 30s
FAIL: SessionStart hook behavior probes failed (see the fixture FAIL lines above).
```

Both FAIL lines are wall-clock only — the line immediately above them says the content probes passed. The same fixtures fail the same way on the unmodified baseline (`session-start-hook-behavior.sh` at 66s and `prescribed-shell-canonicalization.sh` at 43s on a more loaded run). Filtering the stock run for everything that is not a duration line, a SKIP, or a passed/PASS line leaves exactly those wall-clock FAIL lines and nothing else. Neither fixture is in this request's write set and neither was changed; the 300-second override lived on the command line only.

## Discovered Tasks

- The shipped reference checker cannot see a `→ **Named Heading**` citation, so deleting a section leaves prose callers dangling while the suite stays green; a checker that extracted those targets and resolved them against the cited file's headings would close it. `_dev/tests/shipped-package-reference-contract.sh`. The builder marked it report only; the review captured it as REQ-582 (teaching the shipped reference contract to resolve arrow-form section citations). → queue as follow-up
- A condensation of material no test owns has no detector at all, so the only guard is a reviewer reading both revisions. A cheap partial guard: compare section-body line counts between the pre- and post-revision of any shipped action file and require a ledger row for every section whose count moved. That is close to the script the builder used to rebuild the corrected ledger. → report only
- `skills/do-work/actions/cleanup.md:48` cites the in-progress record as "(Step 2)"; the heading is and always was "(Step 1)". Pre-existing, outside this request's write set, and the section does state the three things that sentence claims. → report only
- `_dev/tests/contracts/core-checks.sh` carries 29 info-level ShellCheck findings (SC2016, SC2030, SC2031), all pre-existing; the repo lints at `--severity=warning`, where the file is clean. → report only
- Six template fences in `skills/do-work/actions/work-reference.md` open with a `## …` line whose predecessor is the fence marker (Plan, Scope, Pre-Flight, Testing, Implementation Summary, Repository Gate Repair No-Op). Harmless inside fences and identical on the baseline, but the raw-line awk extractions in core-checks.sh would see them if a range boundary were ever placed near one. → report only

## Review

**Overall: 71% at first review; the findings were then closed.**
**Acceptance: Partial at review, remediated after.** The six disclosed deletions were each verified genuinely owned — the reviewer read the named tests and confirmed each pins what the prose promised rather than merely touching the subject.

Four findings were content losses or stale references, and all were fixed in the continuation. Folder Structure had lost the `archive/legacy/` node while cleanup.md and abandon.md still route work into it; Progress Reporting Example had lost ten lines including the only Route A illustration; an internal citation pointed at a rule this sweep deleted; and the closed ready-set enumeration survived in two files. The continuation restored the first two verbatim, closed the references at four sites rather than three, put the nine composed-summary headline strings back at zero line cost, and rebuilt the ledger from the diff rather than from memory.

The under-700 acceptance line is not met and stays not met: 849 lines. The reviewer re-derived every number in the shortfall and confirmed reaching 699 would require deleting an unowned guarantee, which the plan forbids. It agrees with the refusal.

The most valuable finding is outside this request: `shipped-package-reference-contract.sh` cannot see the arrow form of a section citation. 82 such citations ship and two are dangling on the current tree while the suite exits 0. Captured as REQ-582.

## Lessons Learned

Three rules a future sweep of shipped prose can apply.

First, in a file that puts one paragraph on one line, only structural deletions move the line count. Diagrams, fences, tables, comment blocks and list items are the only levers; rewriting prose buys nothing. Line-count budgets written before reading the file will be wrong by a large factor — here a planned 55-line saving was worth 11, and a planned 90 was worth about 15.

Second, when the rule is "delete only what a test owns", a condensation of unowned material is the one move with no detector. The suite cannot fail on it, so it must be disclosed in the ledger or it will be found by a reviewer reading both revisions, or not at all. Build that ledger from the diff, not from memory: doing so here surfaced three changed sections the memory-written ledger had omitted.

Third, a reference checker that resolves links does not see prose citations. Before deleting a heading from shipped Markdown, grep for the heading text itself, and treat a green reference suite as no evidence at all that inbound callers survived.

## Orientation

`work-reference.md` is now 849 lines instead of 1,064, and every section that remains is judgment, schema, or a minimal template, with each deletion traceable to the Go test that owns the contract it described. The gap that made this request risky is now named and queued: the shipped reference checker is blind to arrow-form section citations, so 82 shipped citations of that form are unguarded and two are dangling on the current tree.
