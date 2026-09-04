# REQ-570 Exploration — delete the `pending-heavy-testing` status from the core skill

Scope: `skills/do-work/` and `_dev/`. Excludes `skills/do-work-board/` (REQ-571) and all CHANGELOG files.
All line numbers are from the working tree at 2026-09-05.

Precondition check (REQ-570 Constraints): no queued or working REQ carries the old status.
`grep -rl '^status: pending-heavy-testing' do-work/queue do-work/working` returns nothing.
`do-work/CHECKPOINT.md:4` carries the value only inside the generated `queue_state:` string
(`0 pending-heavy-testing`), which is rewritten by `advance --checkpoint`.

---

## 1. Readers and writers of the literal `pending-heavy-testing`

### Go — production

- `internal/schemanormalization/schema_normalization.go:30` — **enum entry** in `canonicalValues` for the `status` field.
- `internal/schemanormalization/schema_normalization.go:243` — **reader** inside `DependencySourceReady`; the `(status == "pending-heavy-testing" && commit != "")` disjunct.
- `internal/schemanormalization/schema_normalization.go:236-239` — **prose** (doc comment) explaining the heavy-hold source-ready rule.
- `internal/nextselection/next_selection.go:407-408` — **reader/count**; `summarizeQueue` switch case incrementing `summary.PendingHeavyTesting`. Being a non-`pending` status, it is already never selected.
- `internal/publication/answer.go:255-256` — **reader**; guard requiring the status before a `heavy-testing` answer, refusal `ANSWER-HEAVY-STATUS`.
- `internal/publication/answer.go:262` — **writer**; `status := "pending-heavy-testing"` is the default the answer keeps when neither green nor fully-red.
- `internal/lifecycleadvance/advance_commands.go:110-112` — **reader**; `classifyQueueAdvance` skip case grouping it with `pending-answers`, `blocked-archive-collision`, `blocked-dependency-cycle` into `agent judgment: resolve <status>`.
- `internal/lifecycleadvance/checkpoint_commands.go:129-131` — **count**; `checkpointQueueState` format string and the `counts["pending-heavy-testing"]` lookup.
- `internal/resultmodel/result_model.go:1180-1182` — **prose/count**; the text-format `queue:` line naming the bucket.
- `internal/resultmodel/result_model.go:272` — **count field**; `PendingHeavyTesting int \`json:"pending_heavy_testing"\``.

### Go — tests

- `internal/schemanormalization/schema_normalization_test.go:41` — pins the exact enum warning sentence listing all statuses.
- `internal/schemanormalization/schema_normalization_test.go:162,165` — pins `DependencySourceReady` old behavior.
- `internal/dependencygraph/dependency_graph_test.go:24` — fixture builder adds `commit:` when status is the heavy value.
- `internal/dependencygraph/dependency_graph_test.go:105,124` — fixture rows.
- `internal/nextselection/next_selection_test.go:683,690,751` — fixtures and assertion text.
- `internal/publication/answer_test.go:287,294,327,328,343,351,377,414,415` — heavy answer fixtures and manifests.

### Markdown (action files)

- `actions/work.md:16` — prose; exhaustion report line.
- `actions/work.md:91` — prose; status-flow arrow diagram.
- `actions/work.md:97` — prose; the "Special statuses" bullet for the value.
- `actions/work.md:347` — prose; step 2 of the hold writes the status.
- `actions/work.md:351` — prose; drain trigger condition.
- `actions/work.md:362` — prose; skipped lane leaves it at the status.
- `actions/work.md:487` — prose; Orchestrator Checklist testing-judgment line.
- `actions/work.md:501` — prose; Error Handling table row.
- `actions/work-reference.md:14` — prose; architecture flow "skip pending-answers and pending-heavy-testing".
- `actions/work-reference.md:229` — schema comment on `commit:`.
- `actions/work-reference.md:233-237` — schema block declaring `status: pending-heavy-testing` + `status_changed_at`.
- `actions/work-reference.md:282` — Schema Read Contract table row listing the enum.
- `actions/work-reference.md:319-323` — the **Dependency-source-ready status set** section, which is the normative statement of the rule.
- `actions/work-reference.md:671-681` — Composed Exit Summary item 8, "Held-for-heavy-testing section".
- `actions/clarify.md:11,17,31,35,39,52,60` — When to Use, Step 1 scan, Step 2 gate, Step 2.5 body, Step 3 exclusion note.
- `actions/roadmap.md:54` — source list.
- `actions/roadmap.md:68` — "Needs clarification" bucket definition.
- `actions/cleanup.md:47,59,293` — active-status lists (untouched-items rule, scan list, forbidden-touch rule).
- `actions/restart-with-parallel-handoff.md:64` — clarify hint in the handoff command list.

### Shell / `_dev`

- No match. `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/*.sh` contain the token `heavy` only as the verification-tier name (`contract-regressions.sh:11,33,34,48`; `probe-lanes.sh:40,46,58`). See item 9.

---

## 2. The satellite identifiers

- `heavy_verified_at` / `heavy_verified_revision`
  - `internal/publication/answer.go:277,280` — written on green; `:284,285` — deleted on red.
  - `internal/requestmodel/request_model.go:57,58` — struct fields `HeavyVerifiedAt`, `HeavyVerifiedRevision`; `:315` — parsed from frontmatter.
  - `internal/nextselection/next_selection.go:323,326,330` — read by `matchingHeavyReviewPhase`.
  - `internal/nextselection/next_selection_test.go:718,719` — green and stale fixtures.
  - `internal/publication/answer_test.go:309` — asserts both are written.
  - `actions/work-reference.md:239-244` — schema block; `actions/work.md:97`, `actions/clarify.md:52` — prose.
  - **REQ-570 keeps these two fields** as durable evidence written by the drain. Their current sole writer (`answer.go`) is deleted, so a new writer is needed.
- `ResumePhase` / `resume_phase`
  - `internal/nextselection/next_selection.go:315` — set on every `SelectionRecord`.
  - `internal/resultmodel/result_model.go:228` — JSON field; `:1193-1194` — text renderer.
  - `internal/nextselection/next_selection_test.go:732` — assertion.
  - `actions/work.md:362`, `actions/clarify.md:52`, `actions/work-reference.md:240` — prose.
- `matchingHeavyReviewPhase` — defined `internal/nextselection/next_selection.go:322-338`; called only at `:315`. It validates: all three fields nonblank, `heavy_verified_at` parses as RFC3339, `commit` / `heavy_verified_revision` / `HEAD` all resolve, and ancestry `commit → verified_revision → HEAD` holds. Returns `"review"` or `""`.
- `HeavyTestingEvidence` — `internal/publication/publication_types.go:91-96`; embedded as `AnswerManifest.HeavyTesting` at `:105`. Validated by `validateHeavyTestingEvidence` (`answer.go:610-666`), consumed by `appendHeavyVerificationResult` (`answer.go:684-698`).
- `HeavyLaneResult` — `internal/publication/publication_types.go:81-89`. Referenced in a comment at `internal/resultmodel/result_model.go:99` (the `HeavyLaneExecution` type deliberately mirrors its field names); that comment goes stale if the type is deleted.
- `heavy-testing` answer mode — `internal/publication/answer.go:254` (the `case`), and the mode list in the `ANSWER-MODE-INVALID` message at `answer.go:352` ("answer mode must be clarify, heavy-testing, stakeholder, or verify-repair").
- `ANSWER-HEAVY-*` codes — `answer.go:256` (`-STATUS`), `:292` (`-ARCHIVE-FORBIDDEN`), `:612` (`-EVIDENCE-MISSING`), `:615,619,623,630,634,646,650,653,657` (`-EVIDENCE-INVALID`), `:626` (`-EVIDENCE-MISMATCH`), `:637` (`-EVIDENCE-NON-ANCESTOR`), `:640` (`-EVIDENCE-STALE`), `:661` (`-EVIDENCE-NOT-GREEN`). Test pins at `answer_test.go:356,371,417`.
- `pending_heavy_testing` — `internal/resultmodel/result_model.go:272` (JSON tag only).
- `PendingHeavyTesting` — `result_model.go:272`; incremented `next_selection.go:408`; read `result_model.go:1179,1182`; tests `next_selection_test.go:681,689,749`.

---

## 3. How `recover` classifies working claims

Entry point `handleRecover` — `internal/lifecycleadvance/recovery_commands.go:25-149`.

Order of operations:
1. Parse `--assume-sole-authority` / `--take-over REQ-NNN` (`:168-191`). Mode is `observe`, `sole-authority`, or `take-over`.
2. Run `recover-finalization --discover` first (`:48`). Any non-success returns early (`:57-59`).
3. Collect set-aside REQs from the finalization records via reason code `finalization.SetAsideReasonCode` (`:61`, `:156-166`).
4. Discover the repository, then `recoverableWorkingRequests(snapshot)` (`:66`, defined `:193-209`): a request qualifies when `TreeSection == "working"` **and** (`status == "claimed"` **or** `status == "blocked"` with nonblank `blocked_by`). Sorted by request id then relative path. **This is the single seam where a "held for heavy lanes" classification must branch.**
5. Per request, checkpoint evidence is read from the already-discovered `snapshot.CheckpointClaimsByID[requestID]` (`:82`) and mapped to `[]resultmodel.SelectionClaimEvidence` by `recoveryCheckpointEvidence` (`:211-220`).
6. Decision branches, each producing one `RecoveryClaimResult`:
   - set aside → `Decision: "finalization set aside; claim preserved"` (`:84-88`), no finding.
   - not authorized → `Decision: "takeover available; claim preserved"` plus a warning finding `RECOVERY-TAKEOVER-AVAILABLE` (`:89-102`).
   - authorized → `requeststate.TransitionRecover` plan applied (`:103-113`); on failure `Decision: "recovery refused"` (`:116-121`); on success `Decision: "recovered to queue"`, `Recovered: true`, info finding `RECOVERY-CLAIM-RECOVERED` (`:122-131`), then the snapshot is re-discovered (`:132-138`).
   - no working claims at all → info finding `RECOVERY-NONE` (`:140-146`).

Typed record shape (add the new classification in this shape):

- `resultmodel.RecoveryClaimResult` — `internal/resultmodel/result_model.go:529-535`: `RequestID`, `RequestPath`, `CheckpointEvidence []SelectionClaimEvidence`, `Decision string` (free text), `Recovered bool`.
- `resultmodel.RecoveryResult` — `result_model.go:539-546`: `AuthorityMode`, `TakeOverRequestID`, `FinalizationPassed`, `Claims []RecoveryClaimResult`, `NextArgv`, `VerificationArgv`. Normalized (nil-slice fill) at `result_model.go:855-867`.
- `resultmodel.SelectionClaimEvidence` — `result_model.go:235-242`: `Source`, `ClaimedAt`, `Writer`, `Path`, `SourceLine`, `HeaderText`.
- Reason codes in use today are finding `Code` strings, not a typed enum: `RECOVERY-USAGE`, `RECOVERY-DISCOVERY-FAILED`, `RECOVERY-TAKEOVER-NOT-FOUND`, `RECOVERY-TAKEOVER-AVAILABLE`, `RECOVERY-CLAIM-RECOVERED`, `RECOVERY-NONE`. Helpers: `recoveryFailure` `:237`, `recoveryFailureWithState` `:241`, `recoveryFinding` `:246`.

Live-writer / checkpoint plumbing:
- Checkpoint claims are parsed during repository discovery, not by recovery: `internal/repositorymodel/repository_model.go:355-377`, keyed into `CheckpointClaimsByID` (`:108,164,371`). The claim section is delimited by the literal heading `## In Progress (interrupted)` — `CheckpointClaimBounds`, `repository_model.go:381-385`.
- `CheckpointClaimEvidence` struct — `repository_model.go:90-101`: `RequestID`, `ClaimedAt`, `Writer`, `HasWriter`, `RelativePath`, `SourceLine`, `HeaderText`.
- The `recover-claim` transition validates in `internal/requeststate/state_plan.go:141-165`: working tree section, status `claimed` (or `blocked` + `blocked_by`), `--assume-sole-writer`, commit required, and **exactly one** checkpoint evidence mode among `CheckpointWriter` / `CheckpointUnlabeled` / `CheckpointAbsent` / `CheckpointAllEntries` (refusal `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`). Recovery picks `CheckpointAbsent` when no evidence, else `CheckpointAllEntries` (`recovery_commands.go:107-111`).
- Recovered status is decided by `recoveredStatus` — `internal/requeststate/state_apply.go:719-728`: `blocked` stays `blocked`; a body with an unresolved `- [ ]` under `## Open Questions` becomes `pending-answers`; otherwise `pending`. Checkpoint handling for recover: `state_plan.go:320-336` → `planRecoverCheckpoint`.

There is **no clock reading and no live-writer probe** in this path. "Live writer" is only the checkpoint header's `writer:` label, surfaced as evidence, never compared against a process.

---

## 4. `advance` phase composition for a working request

`handleAdvance` — `internal/lifecycleadvance/advance_commands.go:38-62`. A single `REQ-NNN` argument on a non-queue record routes to `classifyAdvance` (`:71-97`), which dispatches by `TreeSection`: `queue` → `classifyQueueAdvance` (`:99-116`), `working` → `classifyWorkingAdvance` (`:118-257`), `archive` → `classifyArchiveAdvance` (`:259-285`).

Sections are parsed generically: `advanceHeadingPattern` (`:24`) matches every `^## <name>$`; `advanceSections` (`:335-357`) records first-occurrence spans and refuses on any duplicate heading. `hasSection` / `hasAnySection` / `sectionContains` at `:359-375`.

Phase table as implemented in `classifyWorkingAdvance`, in evaluation order. Every step also runs a "later evidence exists before X" guard that returns `ADVANCE-EVIDENCE-MISSING` via `missingBeforeLaterRefusal` (`:318-320`). The guard's later-section list is the hand-written tuple `Plan, Exploration, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation`, repeated at `:130, :140, :147, :154, :169, :182, :190, :197, :209, :221, :229, :237, :244`. **`Heavy Verification Plan` is not in any of these lists and is not read by `advance` at all.**

| Order | Guard | Phase emitted | Kind | Line |
|---|---|---|---|---|
| 0 | status not `claimed`/`completed-with-issues` | refuse `ADVANCE-PHASE-UNKNOWN` | — | 121-123 |
| 1 | `route` empty | `agent judgment: triage and open questions` | judgment | 129-135 |
| 2 | route not A/B/C | refuse `ADVANCE-PHASE-UNKNOWN` | — | 136-138 |
| 3 | no `## Triage` | `agent judgment: triage and open questions` | judgment | 139-145 |
| 4 | `## Open Questions` has `- [ ]` | `agent judgment: triage and open questions` | judgment | 146-152 |
| 5 | `estimate.p50_active_minutes` invalid | `estimate-p50` | mechanical | 153-163 |
| 6 | Route A with B/C-only section | refuse `ADVANCE-PHASE-UNKNOWN` | — | 165-167 |
| 7 | no `## Plan` | `agent judgment: planning` (C) / `record planning not required` | judgment | 168-180 |
| 8 | Route C and `planning_at` blank | `agent judgment: planning` | judgment | 181-187 |
| 9 | not-A and no `## Exploration` | `agent judgment: exploration` | judgment | 189-195 |
| 10 | not-A and no `## Scope` | `agent judgment: scope declaration` | judgment | 196-202 |
| 11 | not-A and empty `write_set` | `agent judgment: scope declaration` | judgment | 203-206 |
| 12 | no `## Implementation Summary` | `preflight` (B/C) / `agent judgment: implementation and summary` (A) | mechanical / judgment | 208-219 |
| 13 | no `## Qualification` | `qualify` | mechanical | 220-227 |
| 14 | no `## Testing` | `test-gate` | mechanical | 228-235 |
| 15 | no `## Review` | `agent judgment: review` | judgment | 236-242 |
| 16 | no `## Lessons Learned` | `agent judgment: lessons and orientation` | judgment | 243-249 |
| 17 | no `## Orientation` | `agent judgment: lessons and orientation` | judgment | 250-253 |
| 18 | all present | `finalize` (`finalization.CommandFinalize`) | mechanical | 254-256 |

Finalize is reached only at step 18; `handleAdvance:55-57` then dispatches to `executeAdvanceFinalization` when extra arguments are supplied, otherwise the projection is returned. Non-finalize mechanical phases go to `executeAdvanceEvidenceGates` (`:58`).

Queue-side classification (`:99-116`): `pending` → `claim`; `blocked` → `blocked-check`; `pending-answers` / `pending-heavy-testing` / `blocked-archive-collision` / `blocked-dependency-cycle` → `agent judgment: resolve <status>`; anything else refuses.

---

## 5. Heavy plan and lane runner surfaces

- Command constants — `internal/heavyverification/heavy_commands.go:15-17`: `plan-heavy-verification`, `plan-heavy-revalidation`, `run-heavy-verification`. Registered at `:20-26`.
- `handlePlanHeavyVerification` `:28-38` → `Plan(...)` (`heavy_verification.go:44`); failure codes `HEAVY-PLAN-USAGE`, `HEAVY-PLAN-UNVERIFIABLE`. Result carried on `CommandResult.HeavyVerification` (`*resultmodel.HeavyVerificationPlan`).
- `handlePlanHeavyRevalidation` `:57-67` → `PlanRevalidation` (`heavy_verification.go:74`); codes `HEAVY-REVALIDATION-USAGE`, `HEAVY-REVALIDATION-UNVERIFIABLE`.
- `handleRunHeavyVerification` `:40-55` → `RunLanes(LaneRunRequest{...})` (`heavy_run.go:89`); refusal code via `LaneRunRefusalCode` (`heavy_run.go:51`). Result on `CommandResult.HeavyVerificationRun`; per-lane findings from `laneUnrecordedFinding` `:352`, `laneSkippedFinding` `:363` (this is `HEAVY-RUN-LANE-SKIPPED`), `laneRedFinding` `:374`. Dirty-tree refusal `refuseDirtyTrackedTree` `:183`.
- Default manifest path `_dev/tests/heavy-lanes.json` (`heavy_commands.go:70`).
- **`## Heavy Verification Plan` has no Go writer or parser.** It is written and re-read entirely by action prose: `actions/work.md:347` (write), `actions/work.md:351` (drift comparison at drain), `actions/clarify.md:39` (manual read).
- **`## Heavy Verification Result`** is written by Go only: `internal/publication/answer.go:697` inside `appendHeavyVerificationResult` (`:684-698`), via `appendSectionEvidence`. Test pin `answer_test.go:309`. Prose references at `actions/clarify.md:52`. No parser exists.

---

## 6. Dependency readiness

- `DependencySourceReady(status, implementationCommit string) bool` — `internal/schemanormalization/schema_normalization.go:240-244`, with the explanatory comment at `:236-239`. Normalizes status first, then `IsTerminalSuccess(...) || (status == "pending-heavy-testing" && commit != "")`.
- Callers (only two, both in `dependencygraph`):
  - `internal/dependencygraph/dependency_graph.go:135` — the per-edge switch; a not-ready dependency appends to `node.UnmetDependencies`, which is what surfaces as `DEPENDENCIES-UNMET`.
  - `internal/dependencygraph/dependency_graph.go:167` — inside `duplicateStatusesSatisfied` (`:154-172`), requiring every duplicate record to be source-ready.
- Node building — `dependency_graph.go:122-148`: after `detectCycles`, each node's `DependenciesSatisfied = len(UnmetDependencies) == 0 && !IsCyclic && !IsAmbiguous`, and `IsReady = RequestStatus == "pending" && DependenciesSatisfied` (`:141-142`). Depth memo at `:145-148`. The node already carries `ImplementationCommit`, so accepting `claimed` + commit needs no new field.
- Sibling predicates in the same file: `IsTerminalSuccess` `:220-224`, `IsTerminalResolved` `:227-230`, `IsStopped` `:233-236`, `DependencySatisfied` `:239`.
- Status enum and alias map: `schema_normalization.go:30` holds `canonicalValues`; the `status` aliases (`complete`/`done`/`finished`/`closed` → `completed`; `canceled`/`abandoned`/`wont-do`/`wontfix` → `cancelled`) live in the same contract table around it. The warning sentence built from `canonicalValues` is pinned verbatim by `schema_normalization_test.go:41`.
- Normative prose: `actions/work-reference.md:319-323` (**Dependency-source-ready status set**).

---

## 7. Test conventions and the tests pinning old behavior

Fixture helpers by package:

- `schemanormalization` — no repository fixtures; `schema_normalization_test.go` calls the pure predicates directly with literal strings.
- `dependencygraph` — `graphFixture` struct (`dependency_graph_test.go:12-17`: `requestID`, `status`, `dependencyKey`, `dependencies`) and `buildFixtureGraph(t, []graphFixture)` (`:19-51`), which writes `do-work/queue/<id>.md` into `t.TempDir()`, calls `repositorymodel.DiscoverRepository`, then `BuildGraph`. Line 24 special-cases the heavy status to inject `commit:`. No git repository is created.
- `nextselection` — `writeCommandRequest(t, root, relativePath, requestID, status, extra)` (`next_commands_test.go:191-201`) and `writeRepositorySelectionFixture(t, root, relativePath, contents)` (`next_selection_test.go:811`). Git repositories are created inline where ancestry matters.
- `publication` — `writeFixture(t, root, path, contents, mode)` (`capture_files_test.go:380`), `canonicalREQFixture(requestID, userRequestID)` (`capture_files_test.go:410`), `initializedGitRepository(t)` (`publication_commands_test.go:225`), `runGitFixture` (`:335`), `runGitFixtureOutput` (`answer_test.go:422`), and the heavy-specific `newHeavyAnswerFixture(t, requestID)` (`answer_test.go:329-347`).
- `lifecycleadvance` — `writeAdvanceRequest(t, root, treeSection, requestID, status, frontmatter, body)` (`advance_commands_test.go:341-347`), `writeAdvanceFile` (`:349`), `runAdvanceGit`, `advanceCLIBinary(t)` (builds the CLI once and runs it as a subprocess), `newAdvanceQueueRepository` (`queue_commands_test.go:237`), `legacyCheckpointRepository` (`recovery_commands_test.go:150`), `seedSetAsideRecoveryFixture` (`recovery_set_aside_test.go:82`).

Tests pinning the OLD behavior — delete or rewrite:

- `internal/schemanormalization/schema_normalization_test.go` — `TestTerminalPredicatesKeepFailureAndCancellationDistinct` (`:149`, asserts at `:162,165`) and `TestNormalizeFieldAppliesAliasesDefaultsAndExactWarnings` (`:8`, enum sentence at `:41`).
- `internal/dependencygraph/dependency_graph_test.go` — `TestPendingHeavyDependencyIsSourceReadyUntilItReturnsToPending` (`:103`) and `TestPendingHeavyDependencyRequiresImplementationCommit` (`:121`); helper branch at `:24`.
- `internal/nextselection/next_selection_test.go` — `TestPendingHeavyTestingIsCountedAndNeverSelected` (`:681`), `TestMatchingHeavyEvidenceResumesAtReviewAndStaleEvidenceDoesNot` (`:695`), `TestPendingHeavyTestingSourceAllowsDependentSelection` (`:749`).
- `internal/publication/answer_test.go` — `TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure` (`:264`), `TestHeavyTestingAnswerRejectsNegativeWallSeconds` (`:348`), `TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane` (`:364`), `TestHeavyTestingAnswerRejectsMismatchedEvidence` (`:401`), plus the helper `newHeavyAnswerFixture` (`:329`).
- `internal/lifecycleadvance/checkpoint_commands_test.go` — no test asserts the `queue_state` string; `TestAdvanceCheckpointChangesOnlyCheckpointAndPreservesLiveEntries` (`:13`) checks paths and preserved-claim counts only. The count-string change is unpinned.
- `internal/lifecycleadvance/advance_commands_test.go` — `TestAdvanceCommandPhaseMatrix` (`:62`) contains no heavy case; the queue skip-case at `advance_commands.go:110` is unpinned by name.

---

## 8. Prose passages, by file

**`skills/do-work/actions/work.md`**
- `:16` — § When to Use. Exit condition: only `pending-answers` or `pending-heavy-testing` remain; routes answers to clarify, says held lanes run at exhaustion.
- `:91` — § Request File Schema. The status-flow arrow line including the `claimed → pending-heavy-testing → pending → claimed` detour.
- `:93` — § Request File Schema, next paragraph. States intermediate phases are tracked by `##` sections, not status. This is the rule REQ-570 cites.
- `:97` — § Request File Schema, Special statuses bullet. Full description of the value, its `commit:` requirement, green/red answers, `heavy_verified_*`, and the exhaustion drain.
- `:344-350` — § Qualification and Testing Judgment, **Heavy-test hold (non-blocking)** paragraph plus its four numbered steps: land commit, set status + append `## Heavy Verification Plan` + append the `## Open Questions` machine line, move the REQ back to `do-work/queue/` and remove its checkpoint claim, recompute selection.
- `:351-372` — **Heavy-lane drain (at queue exhaustion)** paragraph, the `run-heavy-verification` code block, the runner-behavior paragraphs (dirty tree, evidence reuse, dispositions), and `:362` which builds the per-REQ `answer` manifests and names `resume_phase: review`.
- `:487` — § Orchestrator Checklist. Testing-judgment checklist line, "park selected REQs as `pending-heavy-testing`".
- `:501` — § Error Handling table row: held REQs remain after the queue empties → run the drain.

**`skills/do-work/actions/work-reference.md`**
- `:14` — § Architecture ASCII flow, "skip pending-answers and pending-heavy-testing".
- `:229` — § Request File Schema, `commit:` comment adding "also required while pending-heavy-testing".
- `:233-237` — § Request File Schema, the block declaring the status with its two comment lines and `status_changed_at`.
- `:239-244` — § Request File Schema, the `heavy_verified_at` / `heavy_verified_revision` block and its four comment lines describing the green answer and `resume_phase: review`.
- `:282` — § Schema Read Contract, the `status` table row enumerating the vocabulary.
- `:319-323` — § **Dependency-source-ready status set**. The normative rule plus the illustrative-consumers paragraph at `:323`.
- `:671-681` — § Composed Exit Summary item 8, **Held-for-heavy-testing section**, including the rendered block and the browser-engine / clarify remediation sentence.

**`skills/do-work/actions/clarify.md`**
- `:11` — § When to Use bullet: a work run left held REQs whose lanes did not run.
- `:17` — § When to Use, the not-applicable condition listing all three statuses.
- `:31` — Step 1 scan, finds `pending-answers` or `pending-heavy-testing`.
- `:35` — Step 2, the exit sentence and the "skip Steps 2.5-5" branch.
- `:37-52` — Step 2.5 **Run held heavy lanes by hand**, in full: `:39` lists held REQs and refuses malformed ones; `:52` builds the `answer` manifest with `mode: heavy-testing` and all the green/red semantics.
- `:60` — Step 3, note excluding held REQs from the builder-decision renderer.

**`skills/do-work/actions/cleanup.md`**
- `:47` — Rule 4, active statuses left untouched.
- `:59` — the scan list of queue files.
- `:293` — the forbidden-touch rule naming the same active statuses.

**`skills/do-work/actions/roadmap.md`**
- `:54` — source list of queue files by status.
- `:68` — **Needs clarification** bucket; also instructs naming the recorded revision and that clarify owns the permission prompt.

**`skills/do-work/actions/restart-with-parallel-handoff.md`**
- `:64` — the handoff command list; includes `do-work clarify` only if a REQ is at `pending-answers` or `pending-heavy-testing`.

**`skills/do-work/tools/do-work-cli/prime-do-work-cli.md`**
- No occurrence of `pending-heavy-testing` or `resume_phase`. Heavy references are package-boundary prose only: `:27` (the `internal/heavyverification/` package description, ending "Per-lane result evidence still belongs to publication and selector resume validation" — the only sentence tied to the deleted seam), `:33` and `:37` (owned-process import rules), `:97` (Windows cross-compile package list).

---

## 9. Contract tests

`grep -i heavy` over `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/*.sh`:

- `_dev/tests/contract-regressions.sh:11,33,34,48` — the `heavy` **verification tier** name (`fast` vs `heavy`), budget label, and tier comparison. Unrelated to the status.
- `_dev/tests/contracts/probe-lanes.sh:40,46,58` — same tier name, plus the shared heavy-probe CLI build and a SKIP message.

`grep pending-heavy-testing` over the same files: **no matches**. `grep -i 'resume_phase'`: no matches.

`_dev/tests/contracts/core-checks.sh` mentions `skills/do-work/actions/clarify.md` only at `:517` and `:530`, inside a scope-drift fixture's file list (`TestREQ-344` multi-path parsing). Those are filenames in a synthetic REQ body, not predicates about heavy testing, and they do not need to change.

Conclusion for the REQ's Detailed Requirements line "Delete the matching predicates in `_dev/tests/contracts/core-checks.sh` in the same commit": there are no such predicates on disk today. The GREEN criterion "the string `pending-heavy-testing` no longer appears in the core skill's shipped files or contract predicates" is already satisfied on the shell side.
